package auth

import (
	"go.uber.org/fx"

	"github.com/AFK068/skilloom/ingest/internal/store/redis"
)

var Module = fx.Module(
	"auth",

	fx.Provide(
		func(c *redis.Client) RedisStore { return c },
		New,
	),
)
