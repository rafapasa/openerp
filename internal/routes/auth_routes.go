package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/openerp/backend/internal/handler"
)

// RegisterAuthRoutes registra as rotas de autenticação
func RegisterAuthRoutes(router *gin.Engine, handler *handler.AuthHandler) {
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/refresh", handler.RefreshToken)
		auth.GET("/me", handler.GetMe)
		auth.POST("/logout", handler.Logout)
	}
}
