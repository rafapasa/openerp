package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/handler"
)

// RegisterProdutoModeloRoutes registra as rotas para o CRUD de modelos de produtos.
func RegisterProdutoModeloRoutes(router *gin.RouterGroup, h *handler.ProdutoModeloHandler) {
	modelos := router.Group("/produto-modelos")
	{
		modelos.POST("", h.Create)
		modelos.GET("", h.List)
		modelos.GET("/:id", h.GetByID)
		modelos.PUT("/:id", h.Update)
		modelos.DELETE("/:id", h.Delete)
	}
}