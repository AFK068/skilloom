package messages

import (
	"go.uber.org/fx"

	"github.com/AFK068/skilloom/ingest/internal/natsclient"
	"github.com/AFK068/skilloom/ingest/internal/store/redis"
)

var Module = fx.Module(
	"messages",

	fx.Provide(
		func(c *redis.Client) RedisStore { return c },
		func(c *natsclient.Client) Publisher { return c },
		New,
	),
)
