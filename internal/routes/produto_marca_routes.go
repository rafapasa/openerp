package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/handler"
)

// RegisterProdutoMarcaRoutes registra as rotas para o CRUD de marcas de produtos.
func RegisterProdutoMarcaRoutes(router *gin.RouterGroup, h *handler.ProdutoMarcaHandler) {
	marcas := router.Group("/produto-marcas")
	{
		marcas.POST("", h.Create)
		marcas.GET("", h.List)
		marcas.GET("/:id", h.GetByID)
		marcas.PUT("/:id", h.Update)
		marcas.DELETE("/:id", h.Delete)
	}
}
