package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/openerp/backend/internal/config"
	"github.com/openerp/backend/internal/database"
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

	// 3. Testar conexão
	if db.IsConnected() {
		log.Println("✅ Conexão com o banco de dados está ativa")
	} else {
		log.Println("⚠️ Conexão com o banco de dados não está ativa")
	}

	// 4. Configurar o router
	router := setupRouter(cfg, db)

	// 5. Iniciar servidor
	port := cfg.APIPort
	log.Printf("🌐 Servidor iniciado em http://localhost:%s", port)
	log.Printf("📝 Ambiente: %s", cfg.APIEnv)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Erro ao iniciar servidor: %v", err)
	}
}

// setupRouter configura as rotas da API
func setupRouter(cfg *config.Config, db *database.MySQL) *gin.Engine {
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
			"service":     "openerp-api",
			"version":     "1.0.0",
		})
	})

	// Rota de ping (mais simples)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Grupo de rotas da API v1
	v1 := router.Group("/api/v1")
	{
		// Rotas públicas
		v1.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "API is running",
				"time":   time.Now().Format(time.RFC3339),
			})
		})

		// Aqui virão as rotas protegidas
		// auth := v1.Group("/auth")
		// {
		// 	auth.POST("/login", authHandler.Login)
		// 	auth.POST("/register", authHandler.Register)
		// }
	}

	return router
}
