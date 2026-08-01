// cmd/api/main.go
package main

import (
	"log"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/openerp/backend/internal/appwire"
	"github.com/openerp/backend/internal/config"
	"github.com/openerp/backend/internal/database"
	"github.com/openerp/backend/internal/handler"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/routes"

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

	// 3. ✅ Wire - Injeta todas as dependências
	container := appwire.InitializeContainer(db.GetDB())

	// 4. Configurar router
	router := gin.Default()

	// Configurar modo do Gin
	if cfg.APIEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Middlewares globais
	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// 5. ✅ Rotas públicas (não autenticadas)
	setupPublicRoutes(router, container.AuthHandler)

	// 6. ✅ Rotas protegidas (autenticadas)
	setupProtectedRoutes(router, container, cfg)

	// 7. Iniciar servidor
	port := cfg.APIPort
	log.Printf("🌐 Servidor iniciado em http://localhost:%s", port)
	log.Printf("📚 Swagger disponível em http://localhost:%s/swagger/index.html", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Erro ao iniciar servidor: %v", err)
	}
}

// setupPublicRoutes configura rotas públicas que não dependem de autenticação
func setupPublicRoutes(
	router *gin.Engine,
	authHandler *handler.AuthHandler,
) {
	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Rotas de autenticação (públicas)
	auth := router.Group("/api/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
	}
}

// setupProtectedRoutes configura rotas protegidas por autenticação
func setupProtectedRoutes(
	router *gin.Engine,
	container *appwire.Container,
	cfg *config.Config,
) {
	// Grupo protegido por autenticação
	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// ✅ Registra todas as rotas usando o container do Wire
	routes.RegisterAllRoutes(api, container)
}
