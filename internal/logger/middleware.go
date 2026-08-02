package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware é o middleware HTTP para logging
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = c.GetString("request_id") // Vem do middleware RequestID
		}

		// Logger com request ID
		logger := Log.With(
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
		)

		// Adiciona ao contexto para uso nos handlers
		c.Set("logger", logger)

		// Log da requisição
		logger.Debug("Request started")

		c.Next()

		// Log da resposta
		latency := time.Since(start)
		status := c.Writer.Status()
		size := c.Writer.Size()

		logger.Info("Request completed",
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.Int("size", size),
		)
	}
}

// GetLogger recupera o logger do contexto
func GetLogger(c *gin.Context) *zap.Logger {
	if logger, exists := c.Get("logger"); exists {
		return logger.(*zap.Logger)
	}
	return Log
}