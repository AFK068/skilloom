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

type Route interface {
	Register(rg *gin.RouterGroup)
}

type Deps struct {
	fx.In

	Cfg    config.Config
	Log    *zap.Logger
	Routes []Route `group:"v1_routes"`
}

type Server struct {
	Deps
	srv *http.Server
	log *zap.Logger
}

func New(d Deps) *Server {
	router := gin.New()
	router.RedirectTrailingSlash = false
	router.Use(gin.Recovery())

	v1 := router.Group("/v1")
	for _, r := range d.Routes {
		r.Register(v1)
	}

	return &Server{
		Deps: d,
		srv: &http.Server{
			Addr:    d.Cfg.HTTP.Addr,
			Handler: router,
		},
		log: d.Log.Named("http"),
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
