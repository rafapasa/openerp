package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/handler"
)

// RegisterTabelaPrecoRoutes registra as rotas para o CRUD de tabelas de preço.
func RegisterTabelaPrecoRoutes(router *gin.RouterGroup, h *handler.TabelaPrecoHandler, itemHandler *handler.TabelaPrecoProdutoHandler) {
	tabelas := router.Group("/tabelas-preco")
	{
		produtos := tabelas.Group("/:id/itens")
		{
			produtos.POST("", itemHandler.Create)
			produtos.GET("", itemHandler.List)
			produtos.GET("/:item", itemHandler.GetByID)
			produtos.PUT("/:item", itemHandler.Update)
			produtos.DELETE("/:item", itemHandler.Delete)
		}

		// Rotas da Tabela de Preço
		tabelas.POST("", h.Create)
		tabelas.GET("", h.List)
		tabelas.GET("/:id", h.GetByID)
		tabelas.PUT("/:id", h.Update)
		tabelas.DELETE("/:id", h.Delete)

		// Rotas dos Produtos dentro da Tabela de Preço

	}
}
