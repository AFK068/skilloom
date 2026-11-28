package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/AFK068/skilloom/ingest/internal/config"
	"github.com/AFK068/skilloom/ingest/internal/natsclient"
	"github.com/AFK068/skilloom/ingest/internal/services/messages"
	"github.com/AFK068/skilloom/ingest/internal/store/redis"
	httpserver "github.com/AFK068/skilloom/ingest/internal/transport/http"
	v1 "github.com/AFK068/skilloom/ingest/internal/transport/http/v1"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	fx.New(
		fx.Provide(func() (*zap.Logger, error) { return zap.NewProduction() }),
		fx.WithLogger(func(l *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: l}
		}),
		config.Module,
		natsclient.Module,
		redis.Module,
		messages.Module,
		v1.Module,
		httpserver.Module,
	).Run()
}
