package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/handler"
)

// RegisterEntidadeLimiteCreditoRoutes registra as rotas para o CRUD de Limite de Crédito de uma Entidade.
func RegisterEntidadeLimiteCreditoRoutes(router *gin.RouterGroup, h *handler.EntidadeLimiteCreditoHandler) {
	limites := router.Group("/limites-credito")
	{
		limites.POST("/", h.Create)
		limites.GET("/", h.List)
		limites.GET("/:id", h.GetByID)
		limites.PUT("/:id", h.Update)
		limites.DELETE("/:id", h.Delete)
	}
}
