package router

import (
	"net/http"

	"feedsystem/internal/middleware"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	r := gin.New()

	r.Use(middleware.RequestID(), gin.Recovery())

	r.GET("/ping", ping)

	return r
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
