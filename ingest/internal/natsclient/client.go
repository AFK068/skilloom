package natsclient

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/AFK068/skilloom/ingest/internal/config"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/jetstream_mock.go -package=mock

const (
	streamName     = "INGEST"
	streamSubjects = "ingest.>"
)

type JetStream interface {
	CreateOrUpdateStream(ctx context.Context, cfg jetstream.StreamConfig) (jetstream.Stream, error)
	Publish(
		ctx context.Context,
		subject string,
		payload []byte,
		opts ...jetstream.PublishOpt,
	) (*jetstream.PubAck, error)
}

type Client struct {
	Conn *nats.Conn
	JS   JetStream
	log  *zap.Logger
}

func New(cfg config.Config, lc fx.Lifecycle, log *zap.Logger) (*Client, error) {
	c := &Client{log: log.Named("nats")}

	conn, err := nats.Connect(cfg.NATS.URL, c.options()...)
	if err != nil {
		return nil, fmt.Errorf("connect to nats %q: %w", cfg.NATS.URL, err)
	}

	js, err := jetstream.New(conn)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("init jetstream: %w", err)
	}

	c.Conn = conn
	c.JS = js
	c.log.Info("connected", zap.String("url", conn.ConnectedUrl()))

	lc.Append(fx.Hook{
		OnStart: c.ensureStream,
		OnStop:  c.close,
	})

	return c, nil
}

func (c *Client) ensureStream(ctx context.Context) error {
	_, err := c.JS.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      streamName,
		Subjects:  []string{streamSubjects},
		Storage:   jetstream.FileStorage,
		Retention: jetstream.WorkQueuePolicy,
	})
	if err != nil {
		return fmt.Errorf("ensure stream %q: %w", streamName, err)
	}
	return nil
}

func (c *Client) options() []nats.Option {
	return []nats.Option{
		nats.Name("skilloom-ingest"),
		nats.MaxReconnects(-1),
		nats.DisconnectErrHandler(c.onDisconnect),
		nats.ReconnectHandler(c.onReconnect),
	}
}

func (c *Client) onDisconnect(_ *nats.Conn, err error) {
	c.log.Warn("disconnected", zap.Error(err))
}

func (c *Client) onReconnect(conn *nats.Conn) {
	c.log.Info("reconnected", zap.String("url", conn.ConnectedUrl()))
}

func (c *Client) close(context.Context) error {
	if err := c.Conn.Drain(); err != nil {
		return fmt.Errorf("drain nats connection: %w", err)
	}
	return nil
}
