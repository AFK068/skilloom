package httpserver

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/AFK068/skilloom/ingest/internal/config"
)

type Server struct {
	srv *http.Server
	log *zap.Logger
}

func New(cfg config.Config, log *zap.Logger) *Server {
	router := gin.New()
	router.RedirectTrailingSlash = false
	router.Use(gin.Recovery())

	return &Server{
		srv: &http.Server{
			Addr:    cfg.HTTP.Addr,
			Handler: router,
		},
		log: log.Named("http"),
	}
}

func (s *Server) Start(lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			s.log.Info("server starting", zap.String("addr", s.srv.Addr))
			go func() {
				if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					s.log.Error("http server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: s.srv.Shutdown,
	})
}
