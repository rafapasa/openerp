package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/handler"
)

// RegisterEntidadeRegimeTributarioRoutes registra as rotas para o CRUD de regimes da entidade.
func RegisterEntidadeRegimeTributarioRoutes(router *gin.RouterGroup, h *handler.EntidadeRegimeTributarioHandler) {
	regimes := router.Group("/:id/regimes-tributarios")
	{
		regimes.POST("/", h.Create)
		regimes.GET("/", h.List)
		regimes.GET("/:item", h.GetByID)
		regimes.PUT("/:item", h.Update)
		regimes.DELETE("/:item", h.Delete)
	}
}
