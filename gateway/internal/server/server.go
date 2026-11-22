package server

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/AFK068/skilloom/gateway/internal/config"
)

type Server struct {
	router *gin.Engine
	addr   string
}

func New(cfg config.Config) *Server {
	router := gin.New()
	router.Use(gin.Recovery())

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
