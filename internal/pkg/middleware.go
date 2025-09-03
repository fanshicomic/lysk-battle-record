package pkg

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			defer close(done)
			c.Next()
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			c.Header("Connection", "close")
			c.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{
				"error":   "Request timeout",
				"message": "Request exceeded timeout limit",
			})
			return
		}
	}
}
