package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/handler"
)

// RegisterProdutoSubgrupoRoutes registra as rotas para o CRUD de subgrupo de produtos.
func RegisterProdutoSubgrupoRoutes(router *gin.RouterGroup, h *handler.ProdutoSubgrupoHandler) {
	subgrupos := router.Group("/produto-subgrupos")
	{
		subgrupos.POST("", h.Create)
		subgrupos.GET("", h.List)
		subgrupos.GET("/:id", h.GetByID)
		subgrupos.PUT("/:id", h.Update)
		subgrupos.DELETE("/:id", h.Delete)
	}
}
