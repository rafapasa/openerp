// internal/wire/wire.go
//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/openerp/backend/internal/handler"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/service"
)

// ============================================================
// MÓDULOS - ENTIDADE ENDERECO
// ============================================================

var EntidadeEnderecoModule = wire.NewSet(
	// Repositório
	repository.NewEntidadeEnderecoRepository,

	// Service
	service.NewEntidadeEnderecoService,

	// Handler
	handler.NewEntidadeEnderecoHandler,
)

// ============================================================
// MÓDULOS - ENTIDADE CONTATO
// ============================================================

var EntidadeContatoModule = wire.NewSet(
	// Repositório
	repository.NewEntidadeContatoRepository,

	// Service
	service.NewEntidadeContatoService,

	// Handler
	handler.NewEntidadeContatoHandler,
)

// ============================================================
// MÓDULOS - ENTIDADE DOCUMENTO
// ============================================================

var EntidadeDocumentoModule = wire.NewSet(
	// Repositório
	repository.NewEntidadeDocumentoRepository,

	// Service
	service.NewEntidadeDocumentoService,

	// Handler
	handler.NewEntidadeDocumentoHandler,
)

// ============================================================
// MÓDULOS - ENTIDADE REGIME TRIBUTARIO
// ============================================================

var EntidadeRegimeTributarioModule = wire.NewSet(
	// Repositório
	repository.NewEntidadeRegimeTributarioRepository,

	// Service
	service.NewEntidadeRegimeTributarioService,

	// Handler
	handler.NewEntidadeRegimeTributarioHandler,
)

// ============================================================
// MÓDULOS - ENTIDADE LIMITE CREDITO
// ============================================================

var EntidadeLimiteCreditoModule = wire.NewSet(
	// Repositório
	repository.NewEntidadeLimiteCreditoRepository,

	// Service
	service.NewEntidadeLimiteCreditoService,

	// Handler
	handler.NewEntidadeLimiteCreditoHandler,
)

// ============================================================
// MÓDULOS - ENTIDADE (Completo)
// ============================================================

var EntidadeModule = wire.NewSet(
	// Repositórios
	repository.NewEntidadeRepository,
	repository.NewEntidadeEnderecoRepository,         // Já incluso no EntidadeEnderecoModule
	repository.NewEntidadeContatoRepository,          // Já incluso no EntidadeContatoModule
	repository.NewEntidadeDocumentoRepository,        // Já incluso no EntidadeDocumentoModule
	repository.NewEntidadeRegimeTributarioRepository, // Já incluso no EntidadeRegimeTributarioModule
	repository.NewEntidadeLimiteCreditoRepository,    // Já incluso no EntidadeLimiteCreditoModule

	// Services
	service.NewEntidadeService,
	service.NewEntidadeEnderecoService,         // Já incluso no EntidadeEnderecoModule
	service.NewEntidadeContatoService,          // Já incluso no EntidadeContatoModule
	service.NewEntidadeDocumentoService,        // Já incluso no EntidadeDocumentoModule
	service.NewEntidadeRegimeTributarioService, // Já incluso no EntidadeRegimeTributarioModule
	service.NewEntidadeLimiteCreditoService,    // Já incluso no EntidadeLimiteCreditoModule

	// Handlers
	handler.NewEntidadeHandler,
	handler.NewEntidadeEnderecoHandler,         // Já incluso no EntidadeEnderecoModule
	handler.NewEntidadeContatoHandler,          // Já incluso no EntidadeContatoModule
	handler.NewEntidadeDocumentoHandler,        // Já incluso no EntidadeDocumentoModule
	handler.NewEntidadeRegimeTributarioHandler, // Já incluso no EntidadeRegimeTributarioModule
	handler.NewEntidadeLimiteCreditoHandler,    // Já incluso no EntidadeLimiteCreditoModule
)

// ============================================================
// MÓDULOS - DOCUMENTO VENDA
// ============================================================

var DocumentoVendaModule = wire.NewSet(
	// Repositórios
	repository.NewDocumentoVendaRepository,
	repository.NewDocumentoVendaPagamentoRepository,
	// repository.NewDocumentoVendaItemRepository, // Se existir

	// Services
	service.NewDocumentoVendaService,

	// Handlers
	handler.NewDocumentoVendaHandler,
)

// ============================================================
// MÓDULOS - PRODUTO
// ============================================================

var ProdutoModule = wire.NewSet(
	// Repositórios
	repository.NewProdutoRepository,
	repository.NewProdutoGrupoRepository,
	repository.NewProdutoSubgrupoRepository,
	repository.NewProdutoMarcaRepository,
	repository.NewProdutoModeloRepository,
	repository.NewProdutoVariacaoRepository,

	// Services
	service.NewProdutoService,
	service.NewProdutoGrupoService,
	service.NewProdutoSubgrupoService,
	service.NewProdutoMarcaService,
	service.NewProdutoModeloService,
	service.NewProdutoVariacaoService,

	// Handlers
	handler.NewProdutoHandler,
	handler.NewProdutoGrupoHandler,
	handler.NewProdutoSubgrupoHandler,
	handler.NewProdutoMarcaHandler,
	handler.NewProdutoModeloHandler,
	handler.NewProdutoVariacaoHandler,
)

// ============================================================
// MÓDULOS - TABELA PREÇO
// ============================================================

var TabelaPrecoModule = wire.NewSet(
	// Repositórios
	repository.NewTabelaPrecoRepository,
	repository.NewTabelaPrecoProdutoRepository,

	// Services
	service.NewTabelaPrecoService,
	service.NewTabelaPrecoProdutoService,

	// Handlers
	handler.NewTabelaPrecoHandler,
	handler.NewTabelaPrecoProdutoHandler,
)

// ============================================================
// MÓDULOS - AUTH
// ============================================================

var AuthModule = wire.NewSet(
	// Repositórios
	// repository.NewUsuarioRepository, // Se existir

	// Services
	service.NewAuthService,

	// Handlers
	handler.NewAuthHandler,
)

// ============================================================
// ROUTER CONFIG
// ============================================================

type RouterConfig struct {
	JWTSecret string

	// Handlers
	AuthHandler                     *handler.AuthHandler
	EntidadeHandler                 *handler.EntidadeHandler
	EntidadeEnderecoHandler         *handler.EntidadeEnderecoHandler
	EntidadeContatoHandler          *handler.EntidadeContatoHandler
	EntidadeDocumentoHandler        *handler.EntidadeDocumentoHandler
	EntidadeRegimeTributarioHandler *handler.EntidadeRegimeTributarioHandler
	EntidadeLimiteCreditoHandler    *handler.EntidadeLimiteCreditoHandler
	ProdutoHandler                  *handler.ProdutoHandler
	ProdutoGrupoHandler             *handler.ProdutoGrupoHandler
	ProdutoSubgrupoHandler          *handler.ProdutoSubgrupoHandler
	ProdutoMarcaHandler             *handler.ProdutoMarcaHandler
	ProdutoModeloHandler            *handler.ProdutoModeloHandler
	ProdutoVariacaoHandler          *handler.ProdutoVariacaoHandler
	TabelaPrecoHandler              *handler.TabelaPrecoHandler
	TabelaPrecoProdutoHandler       *handler.TabelaPrecoProdutoHandler
	DocumentoVendaHandler           *handler.DocumentoVendaHandler
}

func NewRouterConfig(
	jwtSecret string,
	authHandler *handler.AuthHandler,
	entidadeHandler *handler.EntidadeHandler,
	entidadeEnderecoHandler *handler.EntidadeEnderecoHandler,
	entidadeContatoHandler *handler.EntidadeContatoHandler,
	entidadeDocumentoHandler *handler.EntidadeDocumentoHandler,
	entidadeRegimeTributarioHandler *handler.EntidadeRegimeTributarioHandler,
	entidadeLimiteCreditoHandler *handler.EntidadeLimiteCreditoHandler,
	produtoHandler *handler.ProdutoHandler,
	produtoGrupoHandler *handler.ProdutoGrupoHandler,
	produtoSubgrupoHandler *handler.ProdutoSubgrupoHandler,
	produtoMarcaHandler *handler.ProdutoMarcaHandler,
	produtoModeloHandler *handler.ProdutoModeloHandler,
	produtoVariacaoHandler *handler.ProdutoVariacaoHandler,
	tabelaPrecoHandler *handler.TabelaPrecoHandler,
	tabelaPrecoProdutoHandler *handler.TabelaPrecoProdutoHandler,
	documentoVendaHandler *handler.DocumentoVendaHandler,
) *RouterConfig {
	return &RouterConfig{
		JWTSecret:                       jwtSecret,
		AuthHandler:                     authHandler,
		EntidadeHandler:                 entidadeHandler,
		EntidadeEnderecoHandler:         entidadeEnderecoHandler,
		EntidadeContatoHandler:          entidadeContatoHandler,
		EntidadeDocumentoHandler:        entidadeDocumentoHandler,
		EntidadeRegimeTributarioHandler: entidadeRegimeTributarioHandler,
		EntidadeLimiteCreditoHandler:    entidadeLimiteCreditoHandler,
		ProdutoHandler:                  produtoHandler,
		ProdutoGrupoHandler:             produtoGrupoHandler,
		ProdutoSubgrupoHandler:          produtoSubgrupoHandler,
		ProdutoMarcaHandler:             produtoMarcaHandler,
		ProdutoModeloHandler:            produtoModeloHandler,
		ProdutoVariacaoHandler:          produtoVariacaoHandler,
		TabelaPrecoHandler:              tabelaPrecoHandler,
		TabelaPrecoProdutoHandler:       tabelaPrecoProdutoHandler,
		DocumentoVendaHandler:           documentoVendaHandler,
	}
}

// SetupRouter configura todas as rotas
func SetupRouter(cfg *RouterConfig) *gin.Engine {
	router := gin.Default()

	// Configurar modo do Gin (pode ser sobrescrito pelo main)
	// gin.SetMode(gin.ReleaseMode)

	// Middlewares globais
	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Rotas públicas
	// router.POST("/api/auth/login", cfg.AuthHandler.Login)

	// Rotas protegidas
	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// ============================================================
	// ROTAS - AUTH
	// ============================================================
	if cfg.AuthHandler != nil {
		auth := api.Group("/auth")
		{
			auth.POST("/login", cfg.AuthHandler.Login)
			auth.POST("/refresh", cfg.AuthHandler.Refresh)
			auth.POST("/logout", cfg.AuthHandler.Logout)
		}
	}

	// ============================================================
	// ROTAS - ENTIDADE (Principal)
	// ============================================================
	if cfg.EntidadeHandler != nil {
		entidade := api.Group("/entidades")
		{
			entidade.GET("", cfg.EntidadeHandler.List)
			entidade.POST("", cfg.EntidadeHandler.Create)
			entidade.GET("/:id", cfg.EntidadeHandler.GetByID)
			entidade.PUT("/:id", cfg.EntidadeHandler.Update)
			entidade.DELETE("/:id", cfg.EntidadeHandler.Delete)
		}
	}

	// ============================================================
	// ROTAS - ENTIDADE ENDERECO
	// ============================================================
	if cfg.EntidadeEnderecoHandler != nil {
		endereco := api.Group("/entidades/:entidade_id/enderecos")
		{
			endereco.GET("", cfg.EntidadeEnderecoHandler.List)
			endereco.POST("", cfg.EntidadeEnderecoHandler.Create)
			endereco.GET("/:item", cfg.EntidadeEnderecoHandler.GetByID)
			endereco.PUT("/:item", cfg.EntidadeEnderecoHandler.Update)
			endereco.DELETE("/:item", cfg.EntidadeEnderecoHandler.Delete)
		}
	}

	// ============================================================
	// ROTAS - ENTIDADE CONTATO
	// ============================================================
	if cfg.EntidadeContatoHandler != nil {
		contato := api.Group("/entidades/:entidade_id/contatos")
		{
			contato.GET("", cfg.EntidadeContatoHandler.List)
			contato.POST("", cfg.EntidadeContatoHandler.Create)
			contato.GET("/:item", cfg.EntidadeContatoHandler.GetByID)
			contato.PUT("/:item", cfg.EntidadeContatoHandler.Update)
			contato.DELETE("/:item", cfg.EntidadeContatoHandler.Delete)
		}
	}

	// ============================================================
	// ROTAS - ENTIDADE DOCUMENTO
	// ============================================================
	if cfg.EntidadeDocumentoHandler != nil {
		documento := api.Group("/entidades/:entidade_id/documentos")
		{
			documento.GET("", cfg.EntidadeDocumentoHandler.List)
			documento.POST("", cfg.EntidadeDocumentoHandler.Create)
			documento.GET("/:item", cfg.EntidadeDocumentoHandler.GetByID)
			documento.PUT("/:item", cfg.EntidadeDocumentoHandler.Update)
			documento.DELETE("/:item", cfg.EntidadeDocumentoHandler.Delete)
		}
	}

	// ============================================================
	// ROTAS - ENTIDADE REGIME TRIBUTARIO
	// ============================================================
	if cfg.EntidadeRegimeTributarioHandler != nil {
		regime := api.Group("/entidades/:entidade_id/regimes-tributarios")
		{
			regime.GET("", cfg.EntidadeRegimeTributarioHandler.List)
			regime.POST("", cfg.EntidadeRegimeTributarioHandler.Create)
			regime.GET("/:item", cfg.EntidadeRegimeTributarioHandler.GetByID)
			regime.PUT("/:item", cfg.EntidadeRegimeTributarioHandler.Update)
			regime.DELETE("/:item", cfg.EntidadeRegimeTributarioHandler.Delete)
		}
	}

	// ============================================================
	// ROTAS - ENTIDADE LIMITE CREDITO
	// ============================================================
	if cfg.EntidadeLimiteCreditoHandler != nil {
		limite := api.Group("/entidades/:entidade_id/limites-credito")
		{
			limite.GET("", cfg.EntidadeLimiteCreditoHandler.List)
			limite.POST("", cfg.EntidadeLimiteCreditoHandler.Create)
			limite.GET("/:item", cfg.EntidadeLimiteCreditoHandler.GetByID)
			limite.PUT("/:item", cfg.EntidadeLimiteCreditoHandler.Update)
			limite.DELETE("/:item", cfg.EntidadeLimiteCreditoHandler.Delete)
		}
	}

	// ============================================================
	// ROTAS - PRODUTO
	// ============================================================
	if cfg.ProdutoHandler != nil {
		produto := api.Group("/produtos")
		{
			produto.GET("", cfg.ProdutoHandler.List)
			produto.POST("", cfg.ProdutoHandler.Create)
			produto.GET("/:id", cfg.ProdutoHandler.GetByID)
			produto.PUT("/:id", cfg.ProdutoHandler.Update)
			produto.DELETE("/:id", cfg.ProdutoHandler.Delete)
		}
	}

	if cfg.ProdutoGrupoHandler != nil {
		grupo := api.Group("/produtos/grupos")
		{
			grupo.GET("", cfg.ProdutoGrupoHandler.List)
			grupo.POST("", cfg.ProdutoGrupoHandler.Create)
			grupo.GET("/:id", cfg.ProdutoGrupoHandler.GetByID)
			grupo.PUT("/:id", cfg.ProdutoGrupoHandler.Update)
			grupo.DELETE("/:id", cfg.ProdutoGrupoHandler.Delete)
		}
	}

	if cfg.ProdutoSubgrupoHandler != nil {
		subgrupo := api.Group("/produtos/subgrupos")
		{
			subgrupo.GET("", cfg.ProdutoSubgrupoHandler.List)
			subgrupo.POST("", cfg.ProdutoSubgrupoHandler.Create)
			subgrupo.GET("/:id", cfg.ProdutoSubgrupoHandler.GetByID)
			subgrupo.PUT("/:id", cfg.ProdutoSubgrupoHandler.Update)
			subgrupo.DELETE("/:id", cfg.ProdutoSubgrupoHandler.Delete)
		}
	}

	if cfg.ProdutoMarcaHandler != nil {
		marca := api.Group("/produtos/marcas")
		{
			marca.GET("", cfg.ProdutoMarcaHandler.List)
			marca.POST("", cfg.ProdutoMarcaHandler.Create)
			marca.GET("/:id", cfg.ProdutoMarcaHandler.GetByID)
			marca.PUT("/:id", cfg.ProdutoMarcaHandler.Update)
			marca.DELETE("/:id", cfg.ProdutoMarcaHandler.Delete)
		}
	}

	if cfg.ProdutoModeloHandler != nil {
		modelo := api.Group("/produtos/modelos")
		{
			modelo.GET("", cfg.ProdutoModeloHandler.List)
			modelo.POST("", cfg.ProdutoModeloHandler.Create)
			modelo.GET("/:id", cfg.ProdutoModeloHandler.GetByID)
			modelo.PUT("/:id", cfg.ProdutoModeloHandler.Update)
			modelo.DELETE("/:id", cfg.ProdutoModeloHandler.Delete)
		}
	}

	if cfg.ProdutoVariacaoHandler != nil {
		variacao := api.Group("/produtos/variacoes")
		{
			variacao.GET("", cfg.ProdutoVariacaoHandler.List)
			variacao.POST("", cfg.ProdutoVariacaoHandler.Create)
			variacao.GET("/:id", cfg.ProdutoVariacaoHandler.GetByID)
			variacao.PUT("/:id", cfg.ProdutoVariacaoHandler.Update)
			variacao.DELETE("/:id", cfg.ProdutoVariacaoHandler.Delete)
		}
	}

	// ============================================================
	// ROTAS - TABELA PREÇO
	// ============================================================
	if cfg.TabelaPrecoHandler != nil {
		tabela := api.Group("/tabelas-preco")
		{
			tabela.GET("", cfg.TabelaPrecoHandler.List)
			tabela.POST("", cfg.TabelaPrecoHandler.Create)
			tabela.GET("/:id", cfg.TabelaPrecoHandler.GetByID)
			tabela.PUT("/:id", cfg.TabelaPrecoHandler.Update)
			tabela.DELETE("/:id", cfg.TabelaPrecoHandler.Delete)
		}
	}

	if cfg.TabelaPrecoProdutoHandler != nil {
		tabelaProduto := api.Group("/tabelas-preco/:tabela_id/produtos")
		{
			tabelaProduto.GET("", cfg.TabelaPrecoProdutoHandler.List)
			tabelaProduto.POST("", cfg.TabelaPrecoProdutoHandler.Create)
			tabelaProduto.GET("/:produto_id", cfg.TabelaPrecoProdutoHandler.GetByID)
			tabelaProduto.PUT("/:produto_id", cfg.TabelaPrecoProdutoHandler.Update)
			tabelaProduto.DELETE("/:produto_id", cfg.TabelaPrecoProdutoHandler.Delete)
		}
	}

	// ============================================================
	// ROTAS - DOCUMENTO VENDA
	// ============================================================
	if cfg.DocumentoVendaHandler != nil {
		documento := api.Group("/documentos/venda")
		{
			documento.GET("", cfg.DocumentoVendaHandler.List)
			documento.POST("", cfg.DocumentoVendaHandler.Create)
			documento.GET("/:id", cfg.DocumentoVendaHandler.GetByID)
			documento.PUT("/:id", cfg.DocumentoVendaHandler.Update)
			documento.DELETE("/:id", cfg.DocumentoVendaHandler.Delete)
			documento.POST("/:id/emitir", cfg.DocumentoVendaHandler.Emitir)
			documento.POST("/:id/cancelar", cfg.DocumentoVendaHandler.Cancelar)
		}
	}

	return router
}

// ============================================================
// INICIALIZADOR (GERADO PELO WIRE)
// ============================================================

// InitializeRouter inicializa toda a aplicação com Wire
func InitializeRouter(db *gorm.DB, jwtSecret string) (*gin.Engine, error) {
	wire.Build(
		wire.Value(db),
		wire.Value(jwtSecret),

		// Módulos
		AuthModule,
		EntidadeModule,
		EntidadeEnderecoModule,
		EntidadeContatoModule,
		EntidadeDocumentoModule,
		EntidadeRegimeTributarioModule,
		EntidadeLimiteCreditoModule,
		ProdutoModule,
		TabelaPrecoModule,
		DocumentoVendaModule,

		// Router
		NewRouterConfig,
		SetupRouter,
	)
	return nil, nil // Wire substitui isso pelo código gerado
}
