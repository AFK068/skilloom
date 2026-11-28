package messages

import (
	"go.uber.org/fx"

	natspub "github.com/AFK068/skilloom/ingest/internal/publisher/nats"
)

var Module = fx.Module(
	"messages",

	fx.Provide(
		func(p *natspub.Publisher) Publisher { return p },
		New,
	),
)
