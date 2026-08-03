package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"github.com/openerp/backend/internal/appwire"
	"github.com/openerp/backend/internal/config"
	"github.com/openerp/backend/internal/database"
	"github.com/openerp/backend/internal/handler"
	"github.com/openerp/backend/internal/logger"
	"github.com/openerp/backend/internal/metrics"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/routes"
	"github.com/openerp/backend/internal/tracing"

	_ "github.com/openerp/backend/docs"
)

func main() {
	// 1. Carregar configurações
	cfg := config.LoadConfig()

	// 2. Inicializar Logger
	logger.Init(cfg.APIEnv)
	defer logger.Sync()

	logger.Log.Info("🚀 Iniciando OpenERP API",
		zap.String("ambiente", cfg.APIEnv),
	)

	// 3. Inicializar Tracing (se habilitado)
	if cfg.TracingEnabled {
		if err := tracing.InitTracer("openerp-api", cfg.JaegerEndpoint); err != nil {
			logger.Log.Fatal("❌ Erro ao inicializar tracer", zap.Error(err))
		}
		logger.Log.Info("✅ Tracing inicializado")
	}

	// 4. Inicializar Métricas
	metrics.Init()
	logger.Log.Info("✅ Métricas inicializadas")

	// 5. Conectar ao banco de dados
	dbMySQL, err := database.NewMySQL(cfg)
	if err != nil {
		logger.Log.Fatal("❌ Erro ao conectar ao MySQL", zap.Error(err))
	}
	defer dbMySQL.Close()
	database.RegisterAuditCallbacks(dbMySQL.GetDB())
	logger.Log.Info("✅ Conectado ao MySQL")

	dbRedis, err := database.NewRedis(cfg)
	if err != nil {
		logger.Log.Fatal("❌ Erro ao conectar ao Redis", zap.Error(err))
	}
	defer dbRedis.Close()
	logger.Log.Info("✅ Conectado ao Redis")

	// 6. Wire - Container de dependências
	container := appwire.InitializeContainer(dbMySQL.GetDB(), dbRedis)

	// 7. Configurar router
	router := setupRouter(cfg)

	// 8. Middlewares
	setupMiddlewares(router, cfg, dbMySQL)

	// 9. Rotas
	setupPublicRoutes(router, container.AuthHandler)
	setupProtectedRoutes(router, container, cfg)

	// 10. Servidor HTTP
	srv := &http.Server{
		Addr:         ":" + cfg.APIPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 11. Iniciar servidor
	go func() {
		logger.Log.Info("🌐 Servidor iniciado",
			zap.String("porta", cfg.APIPort),
			zap.String("swagger", "http://localhost:"+cfg.APIPort+"/swagger/index.html"),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("❌ Erro ao iniciar servidor", zap.Error(err))
		}
	}()

	// 12. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("🛑 Desligando servidor...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Error("❌ Erro no shutdown", zap.Error(err))
	}
	logger.Log.Info("✅ Servidor desligado")
}

func setupRouter(cfg *config.Config) *gin.Engine {
	if cfg.APIEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	return router
}

func setupMiddlewares(router *gin.Engine, cfg *config.Config, dbMySQL *database.MySQL) {
	// Ordem importa!
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(logger.LoggerMiddleware())
	router.Use(metrics.PrometheusMiddleware())
	router.Use(middleware.GormContextMiddleware(dbMySQL.GetDB())) // Passe a instância do GORM para o middleware

	if cfg.TracingEnabled {
		router.Use(tracing.TracingMiddleware("openerp-api"))
	}

	if cfg.RateLimitEnabled {
		router.Use(middleware.RateLimitMiddleware(cfg.RateLimitPerSecond))
	}

	// Health checks
	router.GET("/health", healthCheck)
	router.GET("/ready", readinessCheck)
	router.GET("/metrics", metrics.MetricsHandler())
}

func setupPublicRoutes(router *gin.Engine, authHandler *handler.AuthHandler) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/api/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
	}
}

func setupProtectedRoutes(router *gin.Engine, container *appwire.Container, cfg *config.Config) {
	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	routes.RegisterAllRoutes(api, container)
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "openerp-api",
	})
}

func readinessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
