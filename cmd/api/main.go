// cmd/api/main.go
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/openerp/backend/internal/config"
	"github.com/openerp/backend/internal/database"
	"github.com/openerp/backend/internal/handler"

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

	// 2. Carregar .env (opcional, já que o config.LoadConfig() pode fazer isso)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	// 3. Conectar ao banco de dados
	db, err := database.NewMySQL(cfg)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar ao banco: %v", err)
	}
	defer db.Close()
	log.Println("✅ Conectado ao banco de dados!")

	// 4. ✅ Inicializar com Wire - retorna o router pronto!
	router, err := wire.InitializeRouter(db.GetDB(), cfg.JWTSecret)
	if err != nil {
		log.Fatalf("❌ Erro ao inicializar a aplicação com Wire: %v", err)
	}

	// 5. Configurar Swagger e Health Checks (rotas que não dependem do Wire)
	setupPublicRoutes(router, cfg, db)

	// 6. Iniciar servidor
	port := cfg.APIPort
	log.Printf("🌐 Servidor iniciado em http://localhost:%s", port)
	log.Printf("📚 Swagger disponível em http://localhost:%s/swagger/index.html", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Erro ao iniciar servidor: %v", err)
	}
}

// setupPublicRoutes configura rotas públicas que não dependem de autenticação/injeção
func setupPublicRoutes(router *gin.Engine, cfg *config.Config, db *database.MySQL) {
	// Configurar modo do Gin
	if cfg.APIEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Middlewares globais
	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Rota para o Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health checks
	handler.RegisterHealthRoutes(router, cfg, db)
}
