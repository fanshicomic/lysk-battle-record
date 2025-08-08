package usecases

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *LyskServer) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
