package main

import (
	"log/slog"
	"net/http"
	"os"

	"feedsystem/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("invalid configuration", "err", err)
		os.Exit(1)
	}

	setupLogger(cfg.Log)
	slog.Debug("configuration loaded",
		"env", cfg.Env,
		"addr", cfg.Addr(),
		"gin_mode", cfg.Server.Mode,
		"log_level", cfg.Log.Level.String(),
		"log_format", cfg.Log.Format,
	)

	gin.SetMode(cfg.Server.Mode)

	r := gin.Default()
	r.GET("/ping", ping)

	slog.Info("starting server",
		"addr", cfg.Addr(),
		"env", cfg.Env,
		"gin_mode", cfg.Server.Mode,
	)
	if err := r.Run(cfg.Addr()); err != nil {
		slog.Error("server exited", "err", err)
		os.Exit(1)
	}
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func setupLogger(lc config.Log) {
	opts := &slog.HandlerOptions{Level: lc.Level}

	var h slog.Handler
	if lc.Format == "json" {
		h = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		h = slog.NewTextHandler(os.Stdout, opts)
	}
	slog.SetDefault(slog.New(h))
}
