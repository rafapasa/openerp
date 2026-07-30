// internal/repository/produto_subgrupo_repository.go
package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// INTERFACE
// ============================================================

// ProdutoSubgrupoRepository define o contrato para operações de banco
type ProdutoSubgrupoRepository interface {
	// CRUD Básico
	Create(subgrupo *models.ProdutoSubgrupo) error
	Update(id int, subgrupo *models.ProdutoSubgrupo) error
	Delete(id int) error
	FindByID(id int) (*models.ProdutoSubgrupo, error)
	GetByID(id int) (*models.ProdutoSubgrupo, error)

	// Buscas Específicas
	FindByDescricao(descricao string, limit int) ([]models.ProdutoSubgrupo, error)
	FindBySituacao(situacao int) ([]models.ProdutoSubgrupo, error)
	FindActive() ([]models.ProdutoSubgrupo, error)

	// Listagem com Filtros
	List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoSubgrupo, int64, error)
	ListWithProdutos(limit, offset int, filters map[string]interface{}) ([]models.ProdutoSubgrupo, int64, error)
	FindAll() ([]models.ProdutoSubgrupo, error)

	// Consultas de Validação (APENAS CONSULTAS)
	ExistsByDescricao(descricao string, excludeID int) (bool, error)
	ExistsByID(id int) (bool, error)
	Count() (int64, error)
	CountProdutosBySubgrupo(subgrupoID int) (int64, error)

	// Operações em Lote
	BulkUpdateStatus(ids []int, situacao int) error
	BulkDelete(ids []int) error

	// Consultas de Dependências
	HasDependentRecords(id int) (bool, error)
	CountDependentRecords(id int) (map[string]int64, error)
}

// ============================================================
// IMPLEMENTAÇÃO CONCRETA (privada)
// ============================================================

type produtoSubgrupoRepository struct {
	db *gorm.DB
}

// NewProdutoSubgrupoRepository cria uma nova instância (retorna a interface)
func NewProdutoSubgrupoRepository(db *gorm.DB) ProdutoSubgrupoRepository {
	return &produtoSubgrupoRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD (APENAS PERSISTÊNCIA)
// ============================================================

// Create salva um novo subgrupo de produto
func (r *produtoSubgrupoRepository) Create(subgrupo *models.ProdutoSubgrupo) error {
	err := r.db.Create(subgrupo).Error
	if err != nil {
		return apperrors.NewInternalError("Erro ao criar subgrupo de produto.", err)
	}
	return nil
}

// Update atualiza um subgrupo de produto existente
func (r *produtoSubgrupoRepository) Update(id int, subgrupo *models.ProdutoSubgrupo) error {
	err := r.db.
		Omit("Produtos", "created_at", "deleted_at").
		Model(&models.ProdutoSubgrupo{}).
		Where("prosg_id = ?", id).
		Updates(subgrupo).Error
	if err != nil {
		return apperrors.NewInternalError("Erro ao atualizar subgrupo de produtos", err)
	}
	return nil
}

// Delete realiza exclusão lógica
func (r *produtoSubgrupoRepository) Delete(id int) error {
	err := r.db.
		Model(&models.ProdutoSubgrupo{}).
		Where("prosg_id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
	if err != nil {
		return apperrors.NewInternalError("Erro ao excluir subgrupo de produtos", err)
	}
	return nil
}

// FindByID busca um subgrupo de produto pelo ID com relacionamentos
func (r *produtoSubgrupoRepository) FindByID(id int) (*models.ProdutoSubgrupo, error) {
	var subgrupo models.ProdutoSubgrupo
	err := r.db.
		Preload("Produtos", "deleted_at IS NULL").
		Where("prosg_id = ? AND deleted_at IS NULL", id).
		First(&subgrupo).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("subgrupo de produto com ID %d não encontrado", id))
		}
		return nil, apperrors.NewInternalError("Erro buscando sub-grupo de produtos", err)
	}
	return &subgrupo, nil
}

// GetByID busca um subgrupo de produto pelo ID (sem relacionamentos)
func (r *produtoSubgrupoRepository) GetByID(id int) (*models.ProdutoSubgrupo, error) {
	var subgrupo models.ProdutoSubgrupo
	err := r.db.
		Where("prosg_id = ? AND deleted_at IS NULL", id).
		First(&subgrupo).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("subgrupo de produto com ID %d não encontrado", id))
		}
		return nil, apperrors.NewInternalError("Erro o buscar subgrupo de produtos.", err)
	}
	return &subgrupo, nil
}

// ============================================================
// MÉTODOS DE BUSCA ESPECÍFICOS
// ============================================================

// FindByDescricao busca subgrupos de produto pela descrição (autocomplete)
func (r *produtoSubgrupoRepository) FindByDescricao(descricao string, limit int) ([]models.ProdutoSubgrupo, error) {
	var subgrupos []models.ProdutoSubgrupo
	err := r.db.
		Where("prosg_descricao LIKE ? AND deleted_at IS NULL", "%"+descricao+"%").
		Order("prosg_descricao ASC").
		Limit(limit).
		Find(&subgrupos).Error

	if err != nil {
		return nil, apperrors.NewInternalError("Erro o buscar subgrupo de produtos.", err)
	}
	return subgrupos, nil
}

// FindBySituacao busca subgrupos de produto por situação
func (r *produtoSubgrupoRepository) FindBySituacao(situacao int) ([]models.ProdutoSubgrupo, error) {
	var subgrupos []models.ProdutoSubgrupo
	err := r.db.
		Where("prosg_situacao = ? AND deleted_at IS NULL", situacao).
		Order("prosg_descricao ASC").
		Find(&subgrupos).Error

	if err != nil {
		return nil, apperrors.NewInternalError("Erro o buscar subgrupo de produtos.", err)
	}
	return subgrupos, nil
}

// FindActive busca todos os subgrupos de produto ativos
func (r *produtoSubgrupoRepository) FindActive() ([]models.ProdutoSubgrupo, error) {
	var subgrupos []models.ProdutoSubgrupo
	err := r.db.
		Where("prosg_situacao = 1 AND deleted_at IS NULL").
		Order("prosg_descricao ASC").
		Find(&subgrupos).Error

	if err != nil {
		return nil, apperrors.NewInternalError("Erro o buscar subgrupo de produtos.", err)
	}
	return subgrupos, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de subgrupos de produto com paginação e filtros
func (r *produtoSubgrupoRepository) List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoSubgrupo, int64, error) {
	var subgrupos []models.ProdutoSubgrupo
	var total int64

	query := r.db.Model(&models.ProdutoSubgrupo{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.ProdutoSubgrupo{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("Erro o buscar subgrupo de produtos.", err)
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Order("prosg_descricao ASC").
		Find(&subgrupos).Error

	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro o buscar subgrupo de produtos.", err)
	}

	return subgrupos, total, nil
}

// ListWithProdutos retorna subgrupos de produto com contagem de produtos
func (r *produtoSubgrupoRepository) ListWithProdutos(limit, offset int, filters map[string]interface{}) ([]models.ProdutoSubgrupo, int64, error) {
	var subgrupos []models.ProdutoSubgrupo
	var total int64

	query := r.db.Model(&models.ProdutoSubgrupo{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.ProdutoSubgrupo{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("Erro o buscar subgrupo de produtos.", err)
	}

	err := query.
		Select(`
			prosg.*,
			(
				SELECT COUNT(*) 
				FROM produto 
				WHERE prosg_id = produto.prosg_id 
				AND deleted_at IS NULL
			) as total_produtos
		`).
		Limit(limit).
		Offset(offset).
		Order("prosg_descricao ASC").
		Find(&subgrupos).Error

	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro o buscar subgrupo de produtos.", err)
	}

	return subgrupos, total, nil
}

// FindAll busca todos os subgrupos de produto
func (r *produtoSubgrupoRepository) FindAll() ([]models.ProdutoSubgrupo, error) {
	var subgrupos []models.ProdutoSubgrupo
	err := r.db.
		Where("deleted_at IS NULL").
		Order("prosg_descricao ASC").
		Find(&subgrupos).Error
	if err != nil {
		return nil, apperrors.NewInternalError("Erro o buscar subgrupo de produtos.", err)
	}
	return subgrupos, nil
}

// ============================================================
// MÉTODOS DE CONSULTA PARA VALIDAÇÕES (APENAS CONSULTAS)
// ============================================================

// ExistsByDescricao verifica se já existe um subgrupo com a descrição
func (r *produtoSubgrupoRepository) ExistsByDescricao(descricao string, excludeID int) (bool, error) {
	if descricao == "" {
		return false, nil
	}

	var count int64
	query := r.db.Model(&models.ProdutoSubgrupo{}).
		Where("prosg_descricao = ? AND deleted_at IS NULL", descricao)

	if excludeID > 0 {
		query = query.Where("prosg_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, apperrors.NewInternalError("Erro ao verificar descrição.", err)
	}
	return count > 0, nil
}

// ExistsByID verifica se um subgrupo de produto existe pelo ID
func (r *produtoSubgrupoRepository) ExistsByID(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.ProdutoSubgrupo{}).
		Where("prosg_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	if err != nil {
		return false, apperrors.NewInternalError("Erro ao buscando subgrupo de produto.", err)
	}
	return count > 0, nil
}

// Count retorna o total de subgrupos de produto
func (r *produtoSubgrupoRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.ProdutoSubgrupo{}).
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, apperrors.NewInternalError("Erro ao buscando subgrupo de produto.", err)
	}
	return count, nil
}

// CountProdutosBySubgrupo retorna a quantidade de produtos de um subgrupo
func (r *produtoSubgrupoRepository) CountProdutosBySubgrupo(subgrupoID int) (int64, error) {
	var count int64
	err := r.db.Model(&models.Produto{}).
		Where("prosg_id = ? AND deleted_at IS NULL", subgrupoID).
		Count(&count).Error
	if err != nil {
		return 0, apperrors.NewInternalError("Erro ao buscando subgrupo de produto.", err)
	}
	return count, nil
}

// ============================================================
// OPERAÇÕES EM LOTE
// ============================================================

// BulkUpdateStatus atualiza a situação de múltiplos subgrupos
func (r *produtoSubgrupoRepository) BulkUpdateStatus(ids []int, situacao int) error {
	err := r.db.Model(&models.ProdutoSubgrupo{}).
		Where("prosg_id IN ? AND deleted_at IS NULL", ids).
		Update("prosg_situacao", situacao).Error
	if err != nil {
		return apperrors.NewInternalError("Erro ao atualizar subgrupo de produto.", err)
	}
	return nil
}

// BulkDelete realiza exclusão lógica de múltiplos subgrupos
func (r *produtoSubgrupoRepository) BulkDelete(ids []int) error {
	err := r.db.Model(&models.ProdutoSubgrupo{}).
		Where("prosg_id IN ? AND deleted_at IS NULL", ids).
		Update("deleted_at", gorm.Expr("NOW()")).Error
	if err != nil {
		return apperrors.NewInternalError("Erro ao excluir subgrupo de produto.", err)
	}
	return nil
}

// ============================================================
// MÉTODOS DE CONSULTA DE DEPENDÊNCIAS
// ============================================================

// HasDependentRecords verifica se o subgrupo tem registros dependentes
func (r *produtoSubgrupoRepository) HasDependentRecords(id int) (bool, error) {
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
func (r *produtoSubgrupoRepository) CountDependentRecords(id int) (map[string]int64, error) {
	result := make(map[string]int64)

	// Verifica produtos associados
	var countProdutos int64
	if err := r.db.Model(&models.Produto{}).
		Where("prosg_id = ? AND deleted_at IS NULL", id).
		Count(&countProdutos).Error; err != nil {
		return nil, err
	}
	if countProdutos > 0 {
		result["produtos"] = countProdutos
	}

	return result, nil
}
