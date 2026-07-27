package main

import (
	"fmt"
	"log"

	"github.com/openerp/backend/internal/config"
	"github.com/openerp/backend/internal/database"
	"github.com/openerp/backend/internal/models"
)

func main() {
	fmt.Println("🚀 Iniciando migrações do banco de dados...")

	// Carregar configurações
	cfg := config.LoadConfig()

	// Conectar ao banco de dados
	dbConn, err := database.NewMySQL(cfg)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar ao banco de dados para migração: %v", err)
	}
	db := dbConn.GetDB()

	// Lista de modelos para migrar
	// Adicione todos os seus modelos aqui para que o GORM possa criar/atualizar as tabelas
	err = db.AutoMigrate(
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
		&models.Municipio{},
		&models.Estado{},
	)
	if err != nil {
		log.Fatalf("❌ Erro ao executar AutoMigrate: %v", err)
	}

	fmt.Println("✅ Migrações concluídas com sucesso!")
}
