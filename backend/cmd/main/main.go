package main

import (
	"log"
	"net/http"
	"os"

	"feedsystem/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("invalid configuration: %v", err)
		os.Exit(1)
	}

	gin.SetMode(cfg.Server.Mode)

	r := gin.Default()
	r.GET("/ping", ping)

	log.Printf("starting server addr=%s env=%s gin_mode=%s",
		cfg.Addr(), cfg.Env, cfg.Server.Mode)

	if err := r.Run(cfg.Addr()); err != nil {
		log.Printf("server exited: %v", err)
		os.Exit(1)
	}
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
