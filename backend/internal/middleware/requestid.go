package middleware

import (
	"github.com/gin-gonic/gin"

	"feedsystem/internal/reqid"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = reqid.New()
		}
		c.Request = c.Request.WithContext(reqid.WithContext(c.Request.Context(), id))
		c.Header("X-Request-ID", id)
		c.Next()
	}
}
