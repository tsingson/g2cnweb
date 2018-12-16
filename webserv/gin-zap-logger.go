package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

// Logger .
func ZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
			return
		}
		if raw != "" {
			path = path + "?" + raw
		}
		logger.Info(path,
			zap.String("ip", c.ClientIP()),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("ua", c.Request.UserAgent()),
		)
	}
}
