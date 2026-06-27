package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/openerp/backend/internal/config"
	"github.com/openerp/backend/internal/database"
	"github.com/openerp/backend/internal/handler"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
)

func main() {
	// 1. Carregar configurações
	cfg := config.LoadConfig()
	log.Printf("🚀 Iniciando OpenERP API - Ambiente: %s", cfg.APIEnv)

	// 2. Conectar ao banco de dados
	db, err := database.NewMySQL(cfg)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// 3. Inicializar serviços
	authService := service.NewAuthService(db.GetDB(), cfg)

	// 4. Inicializar handlers
	authHandler := handler.NewAuthHandler(authService)

	// 5. Configurar router
	router := setupRouter(cfg, db, authHandler)

	// 6. Iniciar servidor
	port := cfg.APIPort
	log.Printf("🌐 Servidor iniciado em http://localhost:%s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Erro ao iniciar servidor: %v", err)
	}
}

// setupRouter configura as rotas da API
func setupRouter(cfg *config.Config, db *database.MySQL, authHandler *handler.AuthHandler) *gin.Engine {
	// Configurar modo do Gin
	if cfg.APIEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middlewares globais
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":      "ok",
			"environment": cfg.APIEnv,
			"database":    db.IsConnected(),
			"timestamp":   time.Now().Format(time.RFC3339),
		})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// ============================================================
	// ROTAS PÚBLICAS
	// ============================================================
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// ============================================================
	// ROTAS PROTEGIDAS (requer autenticação)
	// ============================================================
	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// Auth protegido
		api.POST("/auth/logout", authHandler.Logout)
		api.GET("/auth/me", authHandler.GetMe)

		// TODO: Adicionar outras rotas protegidas aqui
		// Entidades
		// api.GET("/entidades", entidadeHandler.List)
		// api.POST("/entidades", entidadeHandler.Create)

		// Produtos
		// api.GET("/produtos", produtoHandler.List)
		// api.POST("/produtos", produtoHandler.Create)

		// Pedidos
		// api.GET("/pedidos", pedidoHandler.List)
		// api.POST("/pedidos", pedidoHandler.Create)
	}

	return router
}
