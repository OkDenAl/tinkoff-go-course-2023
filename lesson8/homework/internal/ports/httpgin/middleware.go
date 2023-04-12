package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework8/pkg/logger"
	"time"
)

func Logger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		status := c.Writer.Status()

		log.Info("latency", latency, "method", c.Request.Method, "path", c.Request.URL.Path, "status", status)
	}
}
