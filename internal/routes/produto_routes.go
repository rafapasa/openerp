// internal/routes/produto_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/appwire"
)

func RegisterProdutoRoutes(router *gin.RouterGroup, container *appwire.Container) {
	// Produtos
	produto := router.Group("/produtos")
	{
		produto.GET("", container.ProdutoHandler.List)
		produto.POST("", container.ProdutoHandler.Create)
		produto.GET("/:id", container.ProdutoHandler.GetByID)
		produto.PUT("/:id", container.ProdutoHandler.Update)
		produto.DELETE("/:id", container.ProdutoHandler.Delete) // Corrected: This was missing
		// produto.GET("/codigo/:codigo", container.ProdutoHandler.GetByCodigo)
		produto.GET("/barras/:codigo_barras", container.ProdutoHandler.GetByCodigoBarras)
	}

	// Grupos
	grupo := router.Group("/produtos/grupos")
	{
		grupo.GET("", container.ProdutoGrupoHandler.List) // Corrected: This was missing
		grupo.POST("", container.ProdutoGrupoHandler.Create)
		grupo.GET("/:id", container.ProdutoGrupoHandler.GetByID)
		grupo.PUT("/:id", container.ProdutoGrupoHandler.Update)
		grupo.DELETE("/:id", container.ProdutoGrupoHandler.Delete)
	}

	// Subgrupos
	subgrupo := router.Group("/produtos/subgrupos")
	{
		subgrupo.GET("", container.ProdutoSubgrupoHandler.List) // Corrected: This was missing
		subgrupo.POST("", container.ProdutoSubgrupoHandler.Create)
		subgrupo.GET("/:id", container.ProdutoSubgrupoHandler.GetByID)
		subgrupo.PUT("/:id", container.ProdutoSubgrupoHandler.Update)
		subgrupo.DELETE("/:id", container.ProdutoSubgrupoHandler.Delete)
	}

	// Marcas
	marca := router.Group("/produtos/marcas")
	{
		marca.GET("", container.ProdutoMarcaHandler.List) // Corrected: This was missing
		marca.POST("", container.ProdutoMarcaHandler.Create)
		marca.GET("/:id", container.ProdutoMarcaHandler.GetByID)
		marca.PUT("/:id", container.ProdutoMarcaHandler.Update)
		marca.DELETE("/:id", container.ProdutoMarcaHandler.Delete)
	}

	// Modelos
	modelo := router.Group("/produtos/modelos")
	{
		modelo.GET("", container.ProdutoModeloHandler.List) // Corrected: This was missing
		modelo.POST("", container.ProdutoModeloHandler.Create)
		modelo.GET("/:id", container.ProdutoModeloHandler.GetByID)
		modelo.PUT("/:id", container.ProdutoModeloHandler.Update)
		modelo.DELETE("/:id", container.ProdutoModeloHandler.Delete)
	}

	// Cores
	cor := router.Group("/produtos/cores")
	{
		cor.GET("", container.ProdutoCorHandler.List) // Corrected: This was missing
		cor.POST("", container.ProdutoCorHandler.Create)
		cor.GET("/:id", container.ProdutoCorHandler.GetByID)
		cor.PUT("/:id", container.ProdutoCorHandler.Update)
		cor.DELETE("/:id", container.ProdutoCorHandler.Delete)
	}

	// Tamanhos
	tamanho := router.Group("/produtos/tamanhos")
	{
		tamanho.GET("", container.ProdutoTamanhoHandler.List) // Corrected: This was missing
		tamanho.POST("", container.ProdutoTamanhoHandler.Create)
		tamanho.GET("/:id", container.ProdutoTamanhoHandler.GetByID)
		tamanho.PUT("/:id", container.ProdutoTamanhoHandler.Update)
		tamanho.DELETE("/:id", container.ProdutoTamanhoHandler.Delete)
	}

	// Variações
	variacao := router.Group("/produtos/:produto_id/variacoes")
	{
		variacao.GET("", container.ProdutoVariacaoHandler.List)
		variacao.POST("", container.ProdutoVariacaoHandler.Create)
		variacao.GET("/:id", container.ProdutoVariacaoHandler.GetByID)
		variacao.PUT("/:id", container.ProdutoVariacaoHandler.Update)
		variacao.DELETE("/:id", container.ProdutoVariacaoHandler.Delete)
		// variacao.PATCH("/estoque", container.ProdutoVariacaoHandler.UpdateEstoque)
	}
}
