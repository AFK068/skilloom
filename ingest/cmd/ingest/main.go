package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/AFK068/skilloom/ingest/internal/config"
	"github.com/AFK068/skilloom/ingest/internal/natsclient"
	"github.com/AFK068/skilloom/ingest/internal/store/redis"
	httpserver "github.com/AFK068/skilloom/ingest/internal/transport/http"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	fx.New(
		fx.Provide(
			config.Load,
			func() (*zap.Logger, error) { return zap.NewProduction() },
			natsclient.New,
			redis.New,
			httpserver.New,
		),
		fx.WithLogger(func(l *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: l}
		}),
		fx.Invoke(
			func(*natsclient.Client, *redis.Client) {},
			(*httpserver.Server).Start,
		),
	).Run()
}
