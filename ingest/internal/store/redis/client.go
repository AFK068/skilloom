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

type Deps struct {
	fx.In

	Cfg config.Config
	LC  fx.Lifecycle
	Log *zap.Logger
}

type Client struct {
	Deps
	RDB RDB
	log *zap.Logger
}

func New(d Deps) *Client {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:     d.Cfg.Redis.Addr,
		Password: d.Cfg.Redis.Password,
		DB:       d.Cfg.Redis.DB,
	})

	c := &Client{
		Deps: d,
		RDB:  rdb,
		log:  d.Log.Named("redis"),
	}

	d.LC.Append(fx.Hook{
		OnStart: c.ping,
		OnStop:  c.close,
	})

	return c
}

func (c *Client) ping(ctx context.Context) error {
	if err := c.RDB.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping redis %q: %w", c.Cfg.Redis.Addr, err)
	}

	c.log.Info("connected", zap.String("addr", c.Cfg.Redis.Addr))
	return nil
}

func (c *Client) close(context.Context) error {
	if err := c.RDB.Close(); err != nil {
		return fmt.Errorf("close redis: %w", err)
	}
	return nil
}
