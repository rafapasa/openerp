package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/openerp/backend/internal/config"
	"github.com/openerp/backend/internal/database"
	"github.com/openerp/backend/internal/handler"
	"github.com/openerp/backend/internal/routes"
	"github.com/openerp/backend/internal/service"
)

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

	// 3. Inicializar serviços
	authService := service.NewAuthService(db.GetDB(), cfg)
	entidadeService := service.NewEntidadeService(db.GetDB())
	// Futuros serviços:
	// produtoService := service.NewProdutoService(db.GetDB())
	// pedidoService := service.NewPedidoService(db.GetDB())

	// 4. Inicializar handlers
	authHandler := handler.NewAuthHandler(authService)
	entidadeHandler := handler.NewEntidadeHandler(entidadeService)
	// Futuros handlers:
	// produtoHandler := handler.NewProdutoHandler(produtoService)
	// pedidoHandler := handler.NewPedidoHandler(pedidoService)

	// 5. Configurar router
	router := setupRouter(cfg, db, authHandler, entidadeHandler)

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
) *gin.Engine {
	// Configurar modo do Gin
	if cfg.APIEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Health checks
	routes.RegisterHealthRoutes(router, cfg, db)

	// Registrar todas as rotas
	routes.RegisterRoutes(router, &routes.RouteConfig{
		// Handlers
		AuthHandler:     authHandler,
		EntidadeHandler: entidadeHandler,
		// Futuros handlers...
		// ProdutoHandler: produtoHandler,
		// PedidoHandler:  pedidoHandler,

		// Configurações
		JWTSecret: cfg.JWTSecret,
	})

	return router
}
