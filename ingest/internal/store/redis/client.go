package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/AFK068/skilloom/ingest/internal/config"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/rdb_mock.go -package=mock

type RDB interface {
	Get(ctx context.Context, key string) *goredis.StringCmd
	Ping(ctx context.Context) *goredis.StatusCmd
	Close() error
}

type Client struct {
	RDB  RDB
	addr string
	log  *zap.Logger
}

func New(cfg config.Config, lc fx.Lifecycle, log *zap.Logger) *Client {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	c := &Client{
		RDB:  rdb,
		addr: cfg.Redis.Addr,
		log:  log.Named("redis"),
	}

	lc.Append(fx.Hook{
		OnStart: c.ping,
		OnStop:  c.close,
	})

	return c
}

func (c *Client) ping(ctx context.Context) error {
	if err := c.RDB.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping redis %q: %w", c.addr, err)
	}

	c.log.Info("connected", zap.String("addr", c.addr))
	return nil
}

func (c *Client) close(context.Context) error {
	if err := c.RDB.Close(); err != nil {
		return fmt.Errorf("close redis: %w", err)
	}
	return nil
}
