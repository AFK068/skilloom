package server

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/AFK068/skilloom/gateway/internal/config"
	"github.com/AFK068/skilloom/gateway/internal/proxy"
)

type Server struct {
	router *gin.Engine
	addr   string
}

func New(cfg config.Config) *Server {
	router := gin.New()
	router.RedirectTrailingSlash = false
	router.UseEscapedPath = true

	// Register routes for each configured upstream.
	for _, route := range cfg.Routes {
		if len(route.Endpoints) == 0 {
			continue
		}

		// Wrap the proxy handler with Gin to handle requests for the route's endpoints.
		handler := gin.WrapH(proxy.New(route.Upstream, route.Prefix))
		for _, endpoint := range route.Endpoints {
			router.Handle(endpoint.Method, route.Prefix+endpoint.Path, handler)
		}
	}

	return &Server{
		router: router,
		addr:   cfg.Listen,
	}
}

func (s *Server) Start() error {
	err := s.router.Run(s.addr)
	if err != nil {
		log.Printf("failed to run Gin server on %q: %v", s.addr, err)
	}

	return err
}
