package middleware

import (
	"go.uber.org/fx"

	"github.com/AFK068/skilloom/ingest/internal/services/auth"
	httpserver "github.com/AFK068/skilloom/ingest/internal/transport/http"
)

var Module = fx.Module(
	"http-middleware",

	fx.Provide(
		func(s *auth.Service) Authenticator { return s },
		fx.Annotate(
			NewAuth,
			fx.As(new(httpserver.AuthMiddleware)),
		),
	),
)
