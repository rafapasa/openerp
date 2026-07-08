package main

import (
	"log"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/openerp/backend/internal/config"
	"github.com/openerp/backend/internal/database"
	"github.com/openerp/backend/internal/handler"
	"github.com/openerp/backend/internal/routes"
	"github.com/openerp/backend/internal/service"

	_ "github.com/openerp/backend/docs"
)

// @title			OpenERP API
// @version		0.1a
// @description	API do sistema OpenERP.
// @host			localhost:8080
// @BasePath		/api/v1
// @schemes		http https
func main() {
	// 1. Carregar configurações
	cfg := config.LoadConfig()
	log.Printf("🚀 Iniciando OpenERP API - Ambiente: %s", cfg.APIEnv)

	// 2. Conectar ao banco de dados
	db, err := database.NewMySQL(cfg)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar ao banco: %v", err)
	}
	defer db.Close()
	log.Println("✅ Conectado ao banco de dados!")

	// 4. Inicializar serviços
	authService := service.NewAuthService(db.GetDB(), cfg)
	entidadeService := service.NewEntidadeService(db.GetDB())
	entidadeEnderecoService := service.NewEntidadeEnderecoService(db.GetDB())
	entidadeContatoService := service.NewEntidadeContatoService(db.GetDB())
	entidadeDocumentoService := service.NewEntidadeDocumentoService(db.GetDB())
	entidadeRegimeTributarioService := service.NewEntidadeRegimeTributarioService(db.GetDB())
	entidadeLimiteCreditoService := service.NewEntidadeLimiteCreditoService(db.GetDB())
	// Futuros serviços:
	produtoService := service.NewProdutoService(db.GetDB())
	produtoGrupoService := service.NewProdutoGrupoService(db.GetDB())
	produtoSubGrupoService := service.NewProdutoSubgrupoService(db.GetDB())
	// pedidoService := service.NewPedidoService(db.GetDB())

	// 5. Inicializar handlers
	authHandler := handler.NewAuthHandler(authService)
	entidadeHandler := handler.NewEntidadeHandler(entidadeService)
	entidadeEnderecoHandler := handler.NewEntidadeEnderecoHandler(entidadeEnderecoService)
	entidadeContatoHandler := handler.NewEntidadeContatoHandler(entidadeContatoService)
	entidadeDocumentoHandler := handler.NewEntidadeDocumentoHandler(entidadeDocumentoService)
	entidadeRegimeTributarioHandler := handler.NewEntidadeRegimeTributarioHandler(entidadeRegimeTributarioService)
	entidadelimiteCreditoHandler := handler.NewEntidadeLimiteCreditoHandler(entidadeLimiteCreditoService)

	// Futuros handlers:
	produtoHandler := handler.NewProdutoHandler(produtoService)
	produtoGrupoHandler := handler.NewProdutoGrupoHandler(produtoGrupoService)
	produtoSubGrupoHandler := handler.NewProdutoSubgrupoHandler(produtoSubGrupoService)

	// pedidoHandler := handler.NewPedidoHandler(pedidoService)

	// 5. Configurar router
	router := setupRouter(cfg,
		db,
		authHandler,
		entidadeHandler,
		entidadeEnderecoHandler,
		entidadeContatoHandler,
		entidadeDocumentoHandler,
		entidadeRegimeTributarioHandler,
		entidadelimiteCreditoHandler,
		produtoHandler,
		produtoGrupoHandler,
		produtoSubGrupoHandler,
	)

	// 6. Iniciar servidor
	port := cfg.APIPort
	log.Printf("🌐 Servidor iniciado em http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Erro ao iniciar servidor: %v", err)
	}
}

// setupRouter configura o router com todas as rotas
func setupRouter(
	cfg *config.Config,
	db *database.MySQL,
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
) *gin.Engine {
	// Configurar modo do Gin
	if cfg.APIEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Rota para o Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health checks
	routes.RegisterHealthRoutes(router, cfg, db)

	// Registrar todas as rotas
	routes.RegisterRoutes(router, &routes.RouteConfig{
		// Handlers
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
		JWTSecret:                       cfg.JWTSecret,
	})

	return router
}
