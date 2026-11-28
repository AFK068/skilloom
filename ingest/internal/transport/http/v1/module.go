package v1

import (
	"go.uber.org/fx"

	"github.com/AFK068/skilloom/ingest/internal/services/messages"
	httpserver "github.com/AFK068/skilloom/ingest/internal/transport/http"
)

var Module = fx.Module(
	"http-v1",

	fx.Provide(
		func(s *messages.Service) MessagesService { return s },
		fx.Annotate(
			NewMessagesHandler,
			fx.As(new(httpserver.Route)),
			fx.ResultTags(`group:"v1_routes"`),
		),
	),
)
