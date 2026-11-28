package natspub

import (
	"go.uber.org/fx"

	"github.com/AFK068/skilloom/ingest/internal/natsclient"
)

var Module = fx.Module(
	"nats-publisher",

	fx.Provide(
		func(c *natsclient.Client) Client { return c },
		New,
	),
)
