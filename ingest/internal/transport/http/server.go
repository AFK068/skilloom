package httpserver

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/AFK068/skilloom/ingest/internal/config"
	"github.com/AFK068/skilloom/ingest/internal/transport/http/respond"
)

type Route interface {
	Register(rg *gin.RouterGroup)
}

type AuthMiddleware interface {
	Handle(c *gin.Context)
}

type Deps struct {
	fx.In

	Cfg    config.Config
	Log    *zap.Logger
	Auth   AuthMiddleware
	Routes []Route `group:"v1_routes"`
}

type Server struct {
	Deps
	srv *http.Server
	log *zap.Logger
}

func New(d Deps) *Server {
	s := &Server{
		Deps: d,
		log:  d.Log.Named("http"),
	}
	s.srv = &http.Server{
		Addr:    d.Cfg.HTTP.Addr,
		Handler: s.newRouter(),
	}

	return s
}

func (s *Server) newRouter() *gin.Engine {
	router := gin.New()
	router.RedirectTrailingSlash = false
	router.Use(gin.Recovery())
	router.NoRoute(noRoute)

	s.registerV1(router)
	return router
}

func (s *Server) registerV1(router *gin.Engine) {
	v1 := router.Group("/v1")
	v1.Use(s.Auth.Handle)
	for _, r := range s.Routes {
		r.Register(v1)
	}
}

func noRoute(c *gin.Context) {
	respond.Error(c, http.StatusNotFound, "route not found")
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
