// internal/repository/tabela_preco_repository.go
package repository

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// INTERFACE
// ============================================================

// TabelaPrecoRepository define o contrato para operações de banco
type TabelaPrecoRepository interface {
	// CRUD Básico
	Create(tabela *models.TabelaPreco) error
	Update(id int, tabela *models.TabelaPreco) error
	Delete(id int) error
	FindByID(id int) (*models.TabelaPreco, error)
	GetByID(id int) (*models.TabelaPreco, error)

	// Buscas Específicas
	FindByDescricao(descricao string, limit int) ([]models.TabelaPreco, error)
	FindActive() ([]models.TabelaPreco, error)
	FindActiveByDate(data time.Time) ([]models.TabelaPreco, error)
	FindByTipo(tipo int) ([]models.TabelaPreco, error)

	// Listagem com Filtros
	List(limit, offset int, filters map[string]interface{}) ([]models.TabelaPreco, int64, error)
	ListWithProdutos(limit, offset int, filters map[string]interface{}) ([]models.TabelaPreco, int64, error)
	FindAll() ([]models.TabelaPreco, error)

	// Consultas de Validação (APENAS CONSULTAS)
	ExistsByDescricao(descricao string, excludeID int) (bool, error)
	ExistsByID(id int) (bool, error)
	Count() (int64, error)
	CountProdutosByTabela(tabelaID int) (int64, error)

	// Operações em Lote
	BulkDelete(ids []int) error

	// Consultas de Dependências
	HasDependentRecords(id int) (bool, error)
	CountDependentRecords(id int) (map[string]int64, error)
}

// ============================================================
// IMPLEMENTAÇÃO CONCRETA (privada)
// ============================================================

type tabelaPrecoRepository struct {
	db *gorm.DB
}

// NewTabelaPrecoRepository cria uma nova instância (retorna a interface)
func NewTabelaPrecoRepository(db *gorm.DB) TabelaPrecoRepository {
	return &tabelaPrecoRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD (APENAS PERSISTÊNCIA)
// ============================================================

// Create salva uma nova tabela de preço
func (r *tabelaPrecoRepository) Create(tabela *models.TabelaPreco) error {
	err := r.db.Create(tabela).Error
	if err != nil {
		return apperrors.NewInternalError("Erro ao criar tabela de preço", err)
	}
	return nil
}

// Update atualiza uma tabela de preço existente
func (r *tabelaPrecoRepository) Update(id int, tabela *models.TabelaPreco) error {
	err := r.db.
		Omit("Produtos", "created_at", "deleted_at").
		Model(&models.TabelaPreco{}).
		Where("tbp_id = ?", id).
		Updates(tabela).Error
	if err != nil {
		return apperrors.NewInternalError("Erro ao atualizar tabela de preço", err)
	}
	return nil
}

// Delete realiza exclusão lógica
func (r *tabelaPrecoRepository) Delete(id int) error {
	err := r.db.
		Model(&models.TabelaPreco{}).
		Where("tbp_id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
	if err != nil {
		return apperrors.NewInternalError("Erro ao excluir tabela de preço", err)
	}
	return nil
}

// FindByID busca uma tabela de preço pelo ID com relacionamentos
func (r *tabelaPrecoRepository) FindByID(id int) (*models.TabelaPreco, error) {
	var tabela models.TabelaPreco
	err := r.db.
		Preload("Produtos", "deleted_at IS NULL").
		Preload("Produtos.Produto").
		Where("tbp_id = ? AND deleted_at IS NULL", id).
		First(&tabela).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("tabela de preço com ID %d não encontrada", id))
		}
		return nil, apperrors.NewInternalError("Erro buscando Tabela de Preço", err)
	}
	return &tabela, nil
}

// GetByID busca uma tabela de preço pelo ID (sem relacionamentos)
func (r *tabelaPrecoRepository) GetByID(id int) (*models.TabelaPreco, error) {
	var tabela models.TabelaPreco
	err := r.db.
		Where("tbp_id = ? AND deleted_at IS NULL", id).
		First(&tabela).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("tabela de preço com ID %d não encontrada", id))
		}
		return nil, apperrors.NewInternalError("Erro buscando Tabela de Preço", err)
	}
	return &tabela, nil
}

// ============================================================
// MÉTODOS DE BUSCA ESPECÍFICOS
// ============================================================

// FindByDescricao busca tabelas de preço pela descrição (autocomplete)
func (r *tabelaPrecoRepository) FindByDescricao(descricao string, limit int) ([]models.TabelaPreco, error) {
	var tabelas []models.TabelaPreco
	err := r.db.
		Where("tbp_descricao LIKE ? AND deleted_at IS NULL", "%"+descricao+"%").
		Order("tbp_descricao ASC").
		Limit(limit).
		Find(&tabelas).Error

	if err != nil {
		return nil, apperrors.NewInternalError("Erro buscando Tabela de Preço", err)
	}
	return tabelas, nil
}

// FindActive busca todas as tabelas de preço ativas (baseado na data atual)
func (r *tabelaPrecoRepository) FindActive() ([]models.TabelaPreco, error) {
	var tabelas []models.TabelaPreco
	now := time.Now()

	err := r.db.
		Where("deleted_at IS NULL").
		Where("tbp_datainicio <= ?", now).
		Where("tbp_datafim IS NULL OR tbp_datafim >= ?", now).
		Order("tbp_descricao ASC").
		Find(&tabelas).Error

	if err != nil {
		return nil, err
	}
	return tabelas, apperrors.NewInternalError("Erro buscando Tabela de Preço", err)
}

// FindActiveByDate busca tabelas de preço ativas em uma data específica
func (r *tabelaPrecoRepository) FindActiveByDate(data time.Time) ([]models.TabelaPreco, error) {
	var tabelas []models.TabelaPreco
	err := r.db.
		Where("deleted_at IS NULL").
		Where("tbp_datainicio <= ?", data).
		Where("tbp_datafim IS NULL OR tbp_datafim >= ?", data).
		Order("tbp_descricao ASC").
		Find(&tabelas).Error

	if err != nil {
		return nil, apperrors.NewInternalError("Erro buscando Tabela de Preço", err)
	}
	return tabelas, nil
}

// FindByTipo busca tabelas de preço por tipo
func (r *tabelaPrecoRepository) FindByTipo(tipo int) ([]models.TabelaPreco, error) {
	var tabelas []models.TabelaPreco
	err := r.db.
		Where("tbp_tipo = ? AND deleted_at IS NULL", tipo).
		Order("tbp_descricao ASC").
		Find(&tabelas).Error

	if err != nil {
		return nil, apperrors.NewInternalError("Erro buscando Tabela de Preço", err)
	}
	return tabelas, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de tabelas de preço com paginação e filtros
func (r *tabelaPrecoRepository) List(limit, offset int, filters map[string]interface{}) ([]models.TabelaPreco, int64, error) {
	var tabelas []models.TabelaPreco
	var total int64

	query := r.db.Model(&models.TabelaPreco{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.TabelaPreco{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("Erro buscando Tabela de Preço", err)
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Order("tbp_descricao ASC").
		Find(&tabelas).Error

	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro buscando Tabela de Preço", err)
	}

	return tabelas, total, nil
}

// ListWithProdutos retorna tabelas de preço com contagem de produtos
func (r *tabelaPrecoRepository) ListWithProdutos(limit, offset int, filters map[string]interface{}) ([]models.TabelaPreco, int64, error) {
	var tabelas []models.TabelaPreco
	var total int64

	query := r.db.Model(&models.TabelaPreco{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.TabelaPreco{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("Erro buscando Tabela de Preço", err)
	}

	err := query.
		Select(`
			tbp.*,
			(
				SELECT COUNT(*) 
				FROM tabela_preco_produto 
				WHERE tbp_id = tabela_preco.tbp_id 
				AND deleted_at IS NULL
			) as total_produtos
		`).
		Limit(limit).
		Offset(offset).
		Order("tbp_descricao ASC").
		Find(&tabelas).Error

	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro buscando Tabela de Preço", err)
	}

	return tabelas, total, nil
}

// FindAll busca todas as tabelas de preço
func (r *tabelaPrecoRepository) FindAll() ([]models.TabelaPreco, error) {
	var tabelas []models.TabelaPreco
	err := r.db.
		Where("deleted_at IS NULL").
		Order("tbp_descricao ASC").
		Find(&tabelas).Error
	if err != nil {
		return nil, apperrors.NewInternalError("Erro buscando Tabela de Preço", err)
	}
	return tabelas, nil
}

// ============================================================
// MÉTODOS DE CONSULTA PARA VALIDAÇÕES (APENAS CONSULTAS)
// ============================================================

// ExistsByDescricao verifica se já existe uma tabela de preço com a descrição
func (r *tabelaPrecoRepository) ExistsByDescricao(descricao string, excludeID int) (bool, error) {
	if descricao == "" {
		return false, nil
	}

	var count int64
	query := r.db.Model(&models.TabelaPreco{}).
		Where("tbp_descricao = ? AND deleted_at IS NULL", descricao)

	if excludeID > 0 {
		query = query.Where("tbp_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, apperrors.NewInternalError("Erro buscando Tabela de Preço", err)
	}
	return count > 0, nil
}

// ExistsByID verifica se uma tabela de preço existe pelo ID
func (r *tabelaPrecoRepository) ExistsByID(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.TabelaPreco{}).
		Where("tbp_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	if err != nil {
		return false, apperrors.NewInternalError("Erro buscando Tabela de Preço", err)
	}
	return count > 0, nil
}

// Count retorna o total de tabelas de preço
func (r *tabelaPrecoRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.TabelaPreco{}).
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, apperrors.NewInternalError("Erro buscando Tabela de Preço", err)
	}
	return count, nil
}

// CountProdutosByTabela retorna a quantidade de produtos em uma tabela de preço
func (r *tabelaPrecoRepository) CountProdutosByTabela(tabelaID int) (int64, error) {
	var count int64
	err := r.db.Model(&models.TabelaPrecoProduto{}).
		Where("tbp_id = ? AND deleted_at IS NULL", tabelaID).
		Count(&count).Error
	if err != nil {
		return 0, apperrors.NewInternalError("Erro buscando Tabela de Preço", err)
	}
	return count, nil
}

// ============================================================
// OPERAÇÕES EM LOTE
// ============================================================

// BulkDelete realiza exclusão lógica de múltiplas tabelas
func (r *tabelaPrecoRepository) BulkDelete(ids []int) error {
	err := r.db.Model(&models.TabelaPreco{}).
		Where("tbp_id IN ? AND deleted_at IS NULL", ids).
		Update("deleted_at", gorm.Expr("NOW()")).Error
	if err != nil {
		return apperrors.NewInternalError("Erro ao excluir tabelas de preço", err)
	}
	return nil
}

// ============================================================
// MÉTODOS DE CONSULTA DE DEPENDÊNCIAS
// ============================================================

// HasDependentRecords verifica se a tabela tem registros dependentes
func (r *tabelaPrecoRepository) HasDependentRecords(id int) (bool, error) {
	counts, err := r.CountDependentRecords(id)
	if err != nil {
		return false, err
	}

	for _, count := range counts {
		if count > 0 {
			return true, nil
		}
	}
	return false, nil
}

// CountDependentRecords retorna a contagem de registros dependentes por tipo
func (r *tabelaPrecoRepository) CountDependentRecords(id int) (map[string]int64, error) {
	result := make(map[string]int64)

	// Verifica produtos associados
	var countProdutos int64
	if err := r.db.Model(&models.TabelaPrecoProduto{}).
		Where("tbp_id = ? AND deleted_at IS NULL", id).
		Count(&countProdutos).Error; err != nil {
		return nil, err
	}
	if countProdutos > 0 {
		result["produtos"] = countProdutos
	}

	// Verifica entidades que usam esta tabela como padrão
	var countEntidades int64
	if err := r.db.Model(&models.Entidade{}).
		Where("tbp_id = ? AND deleted_at IS NULL", id).
		Count(&countEntidades).Error; err != nil {
		return nil, err
	}
	if countEntidades > 0 {
		result["entidades"] = countEntidades
	}

	return result, nil
}
