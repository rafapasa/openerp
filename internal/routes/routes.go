// internal/routes/routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/appwire"
)

// RegisterAllRoutes registra todas as rotas da API usando o container
func RegisterAllRoutes(router *gin.RouterGroup, container *appwire.Container) {
	// Rotas de Entidade
	RegisterEntidadeRoutes(router, container)

	// Rotas de Produto
	RegisterProdutoRoutes(router, container)

	// Rotas de Venda
	RegisterVendaRoutes(router, container)

	// Rotas de Fiscal
	RegisterFiscalRoutes(router, container)

	// Rotas de Tabela Preço
	RegisterTabelaPrecoRoutes(router, container)
}
