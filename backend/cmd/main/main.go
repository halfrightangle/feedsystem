package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/ping", ping)

	log.Println("starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
