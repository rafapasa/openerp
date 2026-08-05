// internal/handler/entidade_handler_test.go
package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/config"
	"github.com/openerp/backend/internal/constants"
	"github.com/openerp/backend/internal/database"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"gorm.io/gorm"
)

// mockAuthMiddleware injeta um user_id no contexto para testes
func mockAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", 1) // Mock user ID
		c.Set("login", "testuser")
		c.Set("grupo_id", 1)
		c.Set("empresa_id", 1)
		c.Next()
	}
}

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

	fmt.Printf("container iniciado.")
	fmt.Scanln()
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
	fmt.Printf("Mysql conectado com sucesso no teste: %s", connStr)

	s.DB = mysqlConn.GetDB()

	// 5. Roda o AutoMigrate
	err = s.DB.AutoMigrate(
		&models.Empresa{},
		&models.EmpresaFilial{},
		&models.Usuario{},
		&models.UsuarioFilial{},
		&models.GrupoUsuario{},
		&models.Moeda{},
		&models.MoedaCotacao{},
		&models.Pais{},
		&models.GrupoEntidade{},
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
	entidadeService := service.NewEntidadeService(nil, nil)
	entidadeHandler := NewEntidadeHandler(entidadeService)

	// Usa um middleware mock para simular a autenticação
	api := router.Group("/api/v1")
	api.Use(mockAuthMiddleware()) // Adiciona o middleware mock
	{
		api.POST("/entidades", entidadeHandler.Create)
		api.GET("/entidades/:id", entidadeHandler.GetByID)
	}

	s.Router = router

	s.insertEmpresa()
	s.insertEmpresaFilial()
	s.insertGrupoUsuario()
	s.insertUsuario()
	s.insertUsuarioFilial()
}

func (s *EntidadeTestSuite) getPayLoadEntidadeCNPJ() dto.EntidadeRequest {
	return dto.EntidadeRequest{
		RazaoSocial:      "Empresa Teste Ltda",
		NomeFantasia:     "Teste Empresa",
		InscricaoFederal: "18.996.639/0001-34",
		DataNascimento:   "2013-09-13",
		TipoPessoa:       constants.TipoPessoaJuridica,
		Situacao:         constants.StatusAtivo,
		EmpresaFilialID:  1,
	}
}

func (s *EntidadeTestSuite) getPayLoadEntidadeCPF() dto.EntidadeRequest {
	return dto.EntidadeRequest{
		RazaoSocial:      "Cliente Teste ",
		NomeFantasia:     "Cliente Fisico",
		InscricaoFederal: "036.248.419-80",
		DataNascimento:   "2013-09-13",
		TipoPessoa:       constants.TipoPessoaFisica,
		Situacao:         constants.StatusAtivo,
		EmpresaFilialID:  1,
	}
}

func (s *EntidadeTestSuite) insertEmpresa() {
	empresa := models.Empresa{
		Nome:      "eTools Tecnologia",
		CreatedAt: time.Now(),
		CreatedBy: utils.IntPtr(1),
		UpdatedAt: time.Now(),
		UpdatedBy: utils.IntPtr(1),
	}
	err := s.DB.Create(&empresa).Error
	if err != nil {
		s.T().Fatalf("Falha ao inserir empresa de teste: %v", err)
	}
}

func (s *EntidadeTestSuite) insertEmpresaFilial() {
	empresaFilial := models.EmpresaFilial{
		EmpresaID: 1,
		Nome:      "Filial Teste",
		CreatedAt: time.Now(),
		CreatedBy: utils.IntPtr(1),
		UpdatedAt: time.Now(),
		UpdatedBy: utils.IntPtr(1),
	}
	err := s.DB.Create(&empresaFilial).Error
	if err != nil {
		s.T().Fatalf("Falha ao inserir empresa filial de teste: %v", err)
	}
}

func (s *EntidadeTestSuite) insertUsuario() {
	usuario := models.Usuario{
		ID:             1,
		Nome:           "teste",
		Login:          "Teste",
		Senha:          "123",
		Situacao:       constants.StatusAtivo,
		GrupoUsuarioID: 1,
		CreatedAt:      time.Now(),
		CreatedBy:      utils.IntPtr(1),
		UpdatedAt:      time.Now(),
		UpdatedBy:      utils.IntPtr(1),
	}
	err := s.DB.Create(&usuario).Error
	if err != nil {
		s.T().Fatalf("Falha ao inserir usuário de teste: %v", err)
	}
}

func (s *EntidadeTestSuite) insertUsuarioFilial() {
	usuarioFilial := models.UsuarioFilial{
		UsuarioID:       1,
		EmpresaFilialID: 1,
	}
	err := s.DB.Create(&usuarioFilial).Error
	if err != nil {
		s.T().Fatalf("Falha ao inserir usuário filial de teste: %v", err)
	}
}

func (s *EntidadeTestSuite) insertGrupoUsuario() {
	grupo := models.GrupoUsuario{
		Descricao: "adm",
		Situacao:  constants.StatusAtivo,
	}
	err := s.DB.Create(&grupo).Error
	if err != nil {
		s.T().Fatalf("Falha ao inserir grupo de usuário de teste: %v", err)
	}
}

func (s *EntidadeTestSuite) TearDownSuite() {
	if s.container != nil {
		s.container.Terminate(s.Ctx)
	}
}

// --- AQUI COMEÇAM OS TESTES EM SI ---

// Testa a criação de uma entidade com sucesso (POST /api/v1/entidades)
func (s *EntidadeTestSuite) TestCreateEntidade_Success() {
	// 1. Monta o payload (DTO de criação)
	payload := s.getPayLoadEntidadeCNPJ()

	jsonBody, _ := json.Marshal(payload)

	// 2. Cria a requisição HTTP
	req := httptest.NewRequest("POST", "/api/v1/entidades", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	// Se a rota exigir JWT, você precisaria gerar um token aqui.
	// Como é o primeiro teste, vou sugerir remover o middleware de autenticação apenas para este teste,
	// ou criar um token mock. Vou considerar que você tem um middleware de auth, então vamos gerar um token mock.

	// (Opcional) Adicionar token JWT mock se o seu handler exigir:
	// token, _ := utils.GenerateToken(1, "admin") // Gere um token para o header
	// req.Header.Set("Authorization", "Bearer "+token)

	// 3. Executa a requisição no router
	w := httptest.NewRecorder()
	s.Router.ServeHTTP(w, req)

	// 4. Valida o resultado
	assert.Equal(s.T(), http.StatusCreated, w.Code, "Deveria retornar 201 Created. Body: %s", w.Body.String())

	// 5. Decodifica a resposta para verificar se o ID foi gerado
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)

	// // A resposta é o objeto diretamente, não tem "data"
	id, ok := response["id"].(float64)
	assert.True(s.T(), ok, "ID não encontrado ou não é número")
	assert.NotZero(s.T(), id)

	// // Verifica no banco

	// reqGet := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/entidades/%d", int(id)), nil)
	// reqGet.Header.Set("Content-Type", "application/json")
	// wGet := httptest.NewRecorder()
	// s.Router.ServeHTTP(wGet, reqGet)
	// assert.Equal(s.T(), http.StatusOK, wGet.Code, "Deveria retornar 200 OK")

	// var entidade models.Entidade
	// err = json.Unmarshal(wGet.Body.Bytes(), &entidade)
	// assert.NoError(s.T(), err)
	// assert.Equal(s.T(), payload.RazaoSocial, entidade.RazaoSocial)
	// // var entidade models.Entidade
	// err = s.DB.First(&entidade, id).Error
	// assert.NoError(s.T(), err)
	// assert.Equal(s.T(), s.getPayLoadEntidadeCNPJ().RazaoSocial, entidade.RazaoSocial)
}

// Testa erro ao enviar documento duplicado (CNPJ já cadastrado)
func (s *EntidadeTestSuite) TestCreateEntidade_DuplicateDocument() {
	// Cenário: Cria a primeira
	payload1 := s.getPayLoadEntidadeCNPJ()
	json1, _ := json.Marshal(payload1)
	req1 := httptest.NewRequest("POST", "/api/v1/entidades", bytes.NewBuffer(json1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	s.Router.ServeHTTP(w1, req1)
	assert.Equal(s.T(), http.StatusCreated, w1.Code, "Deveria retornar 201 Created para a primeira entidade: %s", w1.Body.String())

	// Tenta criar a segunda com o mesmo documento
	payload2 := s.getPayLoadEntidadeCNPJ()
	json2, _ := json.Marshal(payload2)
	req2 := httptest.NewRequest("POST", "/api/v1/entidades", bytes.NewBuffer(json2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	s.Router.ServeHTTP(w2, req2)

	// Valida: deve dar erro (BAD REQUEST ou CONFLICT)
	assert.Equal(s.T(), http.StatusConflict, w2.Code, "Deveria retornar 409 Conflict para documento duplicado. Body: %s", w2.Body.String())

	assert.Contains(s.T(), w2.Body.String(), "já está cadastrado", "Mensagem de erro deve mencionar 'já está cadastrado'")
}

// Testa erro ao enviar documento duplicado (CNPJ já cadastrado)
func (s *EntidadeTestSuite) TestCreateEntidade_ValiudateCNPJ_CPF() {
	// Cenário: Cria a primeira
	payload1 := s.getPayLoadEntidadeCNPJ()
	payload1.InscricaoFederal = "12.345.678/0001-99"
	json1, _ := json.Marshal(payload1)
	req1 := httptest.NewRequest("POST", "/api/v1/entidades", bytes.NewBuffer(json1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	s.Router.ServeHTTP(w1, req1)
	assert.Equal(s.T(), http.StatusBadRequest, w1.Code, "Validação de inscrição federal CNPJ deveria falhar. Body: %s", w1.Body.String())
	assert.True(s.T(), strings.Contains(w1.Body.String(), "documento inválido"), "Mensagem de erro para CNPJ inválido")

	// Tenta criar a segunda com o mesmo documento
	payload2 := s.getPayLoadEntidadeCPF()
	payload2.InscricaoFederal = "123.456.789-00"
	json2, _ := json.Marshal(payload2)
	req2 := httptest.NewRequest("POST", "/api/v1/entidades", bytes.NewBuffer(json2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	s.Router.ServeHTTP(w2, req2)

	// Valida: deve dar erro (BAD REQUEST ou CONFLICT)
	assert.Equal(s.T(), http.StatusBadRequest, w2.Code, "Validação de inscrição federal CPF deveria falhar. Body: %s", w2.Body.String())
	assert.True(s.T(), strings.Contains(w2.Body.String(), "documento inválido"), "Mensagem de erro para CPF inválido")
}

// TearDownTest roda depois de CADA teste (limpa o banco para o próximo não interferir)
func (s *EntidadeTestSuite) TearDownTest() {
	// Limpa a tabela de entidades após cada teste para isolar os cenários
	s.DB.Exec("DELETE FROM entidade_endereco")
	s.DB.Exec("DELETE FROM entidade_formacontato")
	s.DB.Exec("DELETE FROM entidade")
}

// Função que inicia os testes (entrypoint do 'go test')
func TestEntidadeSuite(t *testing.T) {
	suite.Run(t, new(EntidadeTestSuite))
}
