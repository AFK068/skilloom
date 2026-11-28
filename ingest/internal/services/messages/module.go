package messages

import (
	"go.uber.org/fx"

	"github.com/AFK068/skilloom/ingest/internal/natsclient"
)

var Module = fx.Module(
	"messages",

	fx.Provide(
		func(c *natsclient.Client) Publisher { return c },
		New,
	),
)
