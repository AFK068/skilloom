package main

import (
	"log"
	"os"

	"github.com/AFK068/skilloom/gateway/internal/config"
	"github.com/AFK068/skilloom/gateway/internal/server"
)

const (
	configPath = "config.yaml"
)

func main() {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatal("load gateway config: ", err)
	}

	srv := server.New(cfg)
	if err := srv.Start(); err != nil {
		os.Exit(1)
	}
}
