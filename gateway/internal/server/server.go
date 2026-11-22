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

	for _, route := range cfg.Routes {
		handler := gin.WrapH(proxy.New(route.Upstream, route.Prefix))
		router.Any(route.Prefix, handler)
		router.Any(route.Prefix+"/*path", handler)
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
