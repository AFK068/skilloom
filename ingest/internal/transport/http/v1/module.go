package v1

import (
	"go.uber.org/fx"

	httpserver "github.com/AFK068/skilloom/ingest/internal/transport/http"
)

var Module = fx.Module(
	"http-v1",

	fx.Provide(
		fx.Annotate(
			NewMessagesHandler,
			fx.As(new(httpserver.Route)),
			fx.ResultTags(`group:"v1_routes"`),
		),
	),
)
