package routes

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/openerp/backend/internal/config"
	"github.com/openerp/backend/internal/database"
)

// RegisterHealthRoutes registra as rotas de health check
func RegisterHealthRoutes(router *gin.Engine, cfg *config.Config, db *database.MySQL) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":      "ok",
			"environment": cfg.APIEnv,
			"database":    db.IsConnected(),
			"timestamp":   time.Now().Format(time.RFC3339),
			"service":     "openerp-api",
			"version":     "1.0.0",
		})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"time":    time.Now().Format(time.RFC3339),
		})
	})
}
