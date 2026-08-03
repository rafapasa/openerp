package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/handler"
)

// RegisterProdutoVariacaoRoutes registra as rotas para o CRUD de variações de produtos.
func RegisterProdutoVariacaoRoutes(router *gin.RouterGroup, h *handler.ProdutoVariacaoHandler) {
	variacoes := router.Group("/produto-variacoes")
	{
		variacoes.POST("", h.Create)
		variacoes.GET("", h.List)
		variacoes.GET("/:id", h.GetByID)
		variacoes.PUT("/:id", h.Update)
		variacoes.DELETE("/:id", h.Delete)
	}
}
