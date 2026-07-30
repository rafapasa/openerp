package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/openerp/backend/internal/handler"
	"github.com/openerp/backend/internal/middleware"
)

// ============================================================
// TYPES
// ============================================================

// RouteConfig contém todas as dependências para registrar rotas
type RouteConfig struct {
	// Handlers
	AuthHandler                     *handler.AuthHandler
	EntidadeHandler                 *handler.EntidadeHandler
	EntidadeEnderecoHandler         *handler.EntidadeEnderecoHandler
	EntidadeContatoHandler          *handler.EntidadeContatoHandler
	EntidadeDocumentoHandler        *handler.EntidadeDocumentoHandler
	ProdutoHandler                  *handler.ProdutoHandler
	ProdutoGrupoHandler             *handler.ProdutoGrupoHandler
	ProdutoSubgrupoHandler          *handler.ProdutoSubgrupoHandler
	ProdutoMarcaHandler             *handler.ProdutoMarcaHandler
	ProdutoModeloHandler            *handler.ProdutoModeloHandler
	ProdutoVariacaoHandler          *handler.ProdutoVariacaoHandler
	// ProdutoSerieHandler             *handler.ProdutoSerieHandler
	// ProdutoEspecieHandler           *handler.ProdutoEspecieHandler
	// ProdutoCorHandler               *handler.ProdutoCorHandler
	// ProdutoTamanhoHandler           *handler.ProdutoTamanhoHandler	
	TabelaPreco                     *handler.ProdutoModeloHandler
	TabelaPrecoHandler              *handler.TabelaPrecoHandler
	TabelaPrecoProdutoHandler       *handler.TabelaPrecoProdutoHandler
	EntidadeRegimeTributarioHandler *handler.EntidadeRegimeTributarioHandler
	EntidadeLimiteCreditoHandler    *handler.EntidadeLimiteCreditoHandler
	DocumentoVendaHandler           *handler.DocumentoVendaHandler
	JWTSecret                       string
}

// ============================================================
// FUNÇÃO PRINCIPAL
// ============================================================

// RegisterRoutes registra todas as rotas da API
func RegisterRoutes(router *gin.Engine, cfg *RouteConfig) {
	// Rotas públicas (não requerem autenticação)
	registerPublicRoutes(router, cfg)

	// Rotas protegidas (requerem autenticação)
	registerProtectedRoutes(router, cfg)
}

// ============================================================
// ROTAS PÚBLICAS
// ============================================================

func registerPublicRoutes(router *gin.Engine, cfg *RouteConfig) {
	// Rotas de autenticação
	RegisterAuthRoutes(router, cfg.AuthHandler)

	// Rotas públicas adicionais (health, etc)
	// RegisterHealthRoutes(router)
}

// ============================================================
// ROTAS PROTEGIDAS
// ============================================================

func registerProtectedRoutes(router *gin.Engine, cfg *RouteConfig) {
	// Grupo protegido por autenticação
	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// Rotas de entidade
	RegisterEntidadeRoutes(api,
		cfg.EntidadeHandler,
		cfg.EntidadeEnderecoHandler,
		cfg.EntidadeContatoHandler,
		cfg.EntidadeDocumentoHandler,
		cfg.EntidadeRegimeTributarioHandler,
		cfg.EntidadeLimiteCreditoHandler,
	)

	// Futuras rotas:
	RegisterProdutoRoutes(api, cfg.ProdutoHandler)
	RegisterProdutoGrupoRoutes(api, cfg.ProdutoGrupoHandler)
	RegisterProdutoSubgrupoRoutes(api, cfg.ProdutoSubgrupoHandler)
	RegisterProdutoMarcaRoutes(api, cfg.ProdutoMarcaHandler)
	RegisterProdutoVariacaoRoutes(api, cfg.ProdutoVariacaoHandler) // Novo: Rotas de Variação de Produto
	RegisterProdutoModeloRoutes(api, cfg.ProdutoModeloHandler)
	RegisterTabelaPrecoRoutes(api, cfg.TabelaPrecoHandler, cfg.TabelaPrecoProdutoHandler)
	RegisterDocumentoVendaRoutes(api, cfg.DocumentoVendaHandler)
	// RegisterPedidoRoutes(api, cfg.PedidoHandler)
	// RegisterUsuarioRoutes(api, cfg.UsuarioHandler)
}
