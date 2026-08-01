// internal/routes/tabela_preco_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/appwire"
)

func RegisterTabelaPrecoRoutes(router *gin.RouterGroup, container *appwire.Container) {
	// Tabelas de Preço
	tabela := router.Group("/tabelas-preco")
	{
		tabela.GET("", container.TabelaPrecoHandler.List)
		tabela.POST("", container.TabelaPrecoHandler.Create) // Corrected: This was missing
		tabela.GET("/:id", container.TabelaPrecoHandler.GetByID)
		tabela.PUT("/:id", container.TabelaPrecoHandler.Update)
		tabela.DELETE("/:id", container.TabelaPrecoHandler.Delete)
	}

	// Produtos da Tabela
	produto := router.Group("/tabelas-preco/:tabela_id/produtos")
	{
		produto.GET("", container.TabelaPrecoProdutoHandler.List)
		produto.POST("", container.TabelaPrecoProdutoHandler.Create)
		produto.GET("/:item", container.TabelaPrecoProdutoHandler.GetByID)
		produto.PUT("/:item", container.TabelaPrecoProdutoHandler.Update)
		produto.DELETE("/:item", container.TabelaPrecoProdutoHandler.Delete)
	}
}
