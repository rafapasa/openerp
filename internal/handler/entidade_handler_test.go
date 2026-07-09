// internal/handler/entidade_handler_test.go
package handler_test // Recomendo usar package handler_test (black-box)

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/config"
	"github.com/openerp/backend/internal/database"
	"github.com/openerp/backend/internal/handler" // importe seu handler
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/service"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"gorm.io/gorm"
)

type EntidadeTestSuite struct {
	suite.Suite
	container *mysql.MySQLContainer
	DB        *gorm.DB
	Router    *gin.Engine
	Ctx       context.Context
}

func (s *EntidadeTestSuite) SetupSuite() {
	s.Ctx = context.Background()

	// 1. Sobe o container MySQL 8.4
	mysqlContainer, err := mysql.Run(s.Ctx,
		"mysql:8.4", // A imagem agora é um parâmetro obrigatório
		mysql.WithDatabase("openerp_test"),
		mysql.WithUsername("root"),
		mysql.WithPassword("root123"),
	)
	if err != nil {
		s.T().Fatalf("Falha ao iniciar container MySQL: %v", err)
	}
	s.container = mysqlContainer

	// 2. Pega a string de conexão (ex: root:root123@tcp(localhost:3306)/openerp_test?charset=utf8mb4...)
	connStr, err := s.container.ConnectionString(s.Ctx)
	if err != nil {
		s.T().Fatalf("Falha ao pegar connection string: %v", err)
	}

	// 3. Carrega a config padrão (apenas para pegar o env/logger, etc)
	cfg := config.LoadConfig()

	// 4. ✅ CORREÇÃO AQUI: Passa o connStr como argumento variádico!
	mysqlConn, err := database.NewMySQL(cfg, connStr)
	if err != nil {
		s.T().Fatalf("Falha ao conectar ao DB de teste: %v", err)
	}
	s.DB = mysqlConn.GetDB()

	// 5. Roda o AutoMigrate
	err = s.DB.AutoMigrate(
		&models.Entidade{},
		&models.EntidadeEndereco{},
		&models.EntidadeContato{},
		&models.Municipio{}, // se precisar
		&models.Estado{},    // se precisar
	)
	if err != nil {
		s.T().Fatalf("Falha ao rodar AutoMigrate: %v", err)
	}

	// 6. Configura o Router (Gin) SEM autenticação para facilitar o teste
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Inicializa as camadas (Repository -> Service -> Handler)
	entidadeService := service.NewEntidadeService(s.DB)
	entidadeHandler := handler.NewEntidadeHandler(entidadeService)

	// 🔥 IMPORTANTE: NÃO coloque o middleware de JWT aqui,
	// ou crie um middleware mock que injeta um usuário no contexto.
	api := router.Group("/api/v1")
	{
		api.POST("/entidades", entidadeHandler.Create)
		api.GET("/entidades/:id", entidadeHandler.GetByID)
	}

	s.Router = router
}

func (s *EntidadeTestSuite) TearDownSuite() {
	if s.container != nil {
		s.container.Terminate(s.Ctx)
	}
}

// TearDownTest limpa as tabelas depois de cada teste
func (s *EntidadeTestSuite) TearDownTest() {
	// Limpa na ordem correta (respeitando foreign keys)
	s.DB.Exec("DELETE FROM entidade_contatos")
	s.DB.Exec("DELETE FROM entidade_enderecos")
	s.DB.Exec("DELETE FROM entidades")
}

// --- SEUS TESTES AQUI (TestCreateEntidade_Success, etc.) ---
// (Mantenha os mesmos que você já escreveu)

func TestEntidadeSuite(t *testing.T) {
	suite.Run(t, new(EntidadeTestSuite))
}
