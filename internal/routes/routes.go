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
	AuthHandler              *handler.AuthHandler
	EntidadeHandler          *handler.EntidadeHandler
	EntidadeEnderecoHandler  *handler.EntidadeEnderecoHandler
	EntidadeContatoHandler   *handler.EntidadeContatoHandler
	EntidadeDocumentoHandler *handler.EntidadeDocumentoHandler
	// Futuros handlers...
	// ProdutoHandler *handler.ProdutoHandler
	// PedidoHandler  *handler.PedidoHandler

	// Configurações
	JWTSecret string
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
		cfg.EntidadeDocumentoHandler
	)

	// Futuras rotas:
	// RegisterProdutoRoutes(api, cfg.ProdutoHandler)
	// RegisterPedidoRoutes(api, cfg.PedidoHandler)
	// RegisterUsuarioRoutes(api, cfg.UsuarioHandler)
}
