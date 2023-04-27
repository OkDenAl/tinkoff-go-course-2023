package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework10/pkg/logger"
	"time"
)

func Logger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		status := c.Writer.Status()

		log.Info("Latency:", latency, "\tMethod:", c.Request.Method, "\tPath:",
			c.Request.URL.Path, "\tStatus:", status)
	}
}
