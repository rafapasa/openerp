// internal/repository/produto_marca_repository.go
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

// ProdutoMarcaRepository define o contrato para operações de banco
type ProdutoMarcaRepository interface {
	// CRUD Básico
	Create(marca *models.ProdutoMarca) error
	Update(id int, marca *models.ProdutoMarca) error
	Delete(id int) error
	FindByID(id int) (*models.ProdutoMarca, error)
	GetByID(id int) (*models.ProdutoMarca, error)

	// Buscas Específicas
	FindByDescricao(descricao string, limit int) ([]models.ProdutoMarca, error)
	FindBySituacao(situacao int) ([]models.ProdutoMarca, error)
	FindActive() ([]models.ProdutoMarca, error)

	// Listagem com Filtros
	List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoMarca, int64, error)
	ListWithProdutos(limit, offset int, filters map[string]interface{}) ([]models.ProdutoMarca, int64, error)
	FindAll() ([]models.ProdutoMarca, error)

	// Consultas de Validação (APENAS CONSULTAS)
	ExistsByDescricao(descricao string, excludeID int) (bool, error)
	ExistsByID(id int) (bool, error)
	Count() (int64, error)
	CountProdutosByMarca(marcaID int) (int64, error)

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

type produtoMarcaRepository struct {
	db *gorm.DB
}

// NewProdutoMarcaRepository cria uma nova instância (retorna a interface)
func NewProdutoMarcaRepository(db *gorm.DB) ProdutoMarcaRepository {
	return &produtoMarcaRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD (APENAS PERSISTÊNCIA)
// ============================================================

// Create salva uma nova marca de produto
func (r *produtoMarcaRepository) Create(marca *models.ProdutoMarca) error {
	return r.db.Create(marca).Error
}

// Update atualiza uma marca de produto existente
func (r *produtoMarcaRepository) Update(id int, marca *models.ProdutoMarca) error {
	return r.db.
		Omit("Produtos", "created_at", "deleted_at").
		Model(&models.ProdutoMarca{}).
		Where("promar_id = ?", id).
		Updates(marca).Error
}

// Delete realiza exclusão lógica
func (r *produtoMarcaRepository) Delete(id int) error {
	return r.db.
		Model(&models.ProdutoMarca{}).
		Where("promar_id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// FindByID busca uma marca de produto pelo ID com relacionamentos
func (r *produtoMarcaRepository) FindByID(id int) (*models.ProdutoMarca, error) {
	var marca models.ProdutoMarca
	err := r.db.
		Preload("Produtos", "deleted_at IS NULL").
		Where("promar_id = ? AND deleted_at IS NULL", id).
		First(&marca).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("marca de produto com ID %d não encontrada", id))
		}
		return nil, err
	}
	return &marca, nil
}

// GetByID busca uma marca de produto pelo ID (sem relacionamentos)
func (r *produtoMarcaRepository) GetByID(id int) (*models.ProdutoMarca, error) {
	var marca models.ProdutoMarca
	err := r.db.
		Where("promar_id = ? AND deleted_at IS NULL", id).
		First(&marca).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("marca de produto com ID %d não encontrada", id))
		}
		return nil, err
	}
	return &marca, nil
}

// ============================================================
// MÉTODOS DE BUSCA ESPECÍFICOS
// ============================================================

// FindByDescricao busca marcas de produto pela descrição (autocomplete)
func (r *produtoMarcaRepository) FindByDescricao(descricao string, limit int) ([]models.ProdutoMarca, error) {
	var marcas []models.ProdutoMarca
	err := r.db.
		Where("promar_descricao LIKE ? AND deleted_at IS NULL", "%"+descricao+"%").
		Order("promar_descricao ASC").
		Limit(limit).
		Find(&marcas).Error

	if err != nil {
		return nil, err
	}
	return marcas, nil
}

// FindBySituacao busca marcas de produto por situação
func (r *produtoMarcaRepository) FindBySituacao(situacao int) ([]models.ProdutoMarca, error) {
	var marcas []models.ProdutoMarca
	err := r.db.
		Where("promar_situacao = ? AND deleted_at IS NULL", situacao).
		Order("promar_descricao ASC").
		Find(&marcas).Error

	if err != nil {
		return nil, err
	}
	return marcas, nil
}

// FindActive busca todas as marcas de produto ativas
func (r *produtoMarcaRepository) FindActive() ([]models.ProdutoMarca, error) {
	var marcas []models.ProdutoMarca
	err := r.db.
		Where("promar_situacao = 1 AND deleted_at IS NULL").
		Order("promar_descricao ASC").
		Find(&marcas).Error

	if err != nil {
		return nil, err
	}
	return marcas, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de marcas de produto com paginação e filtros
func (r *produtoMarcaRepository) List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoMarca, int64, error) {
	var marcas []models.ProdutoMarca
	var total int64

	query := r.db.Model(&models.ProdutoMarca{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.ProdutoMarca{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Order("promar_descricao ASC").
		Find(&marcas).Error

	if err != nil {
		return nil, 0, err
	}

	return marcas, total, nil
}

// ListWithProdutos retorna marcas de produto com contagem de produtos
func (r *produtoMarcaRepository) ListWithProdutos(limit, offset int, filters map[string]interface{}) ([]models.ProdutoMarca, int64, error) {
	var marcas []models.ProdutoMarca
	var total int64

	query := r.db.Model(&models.ProdutoMarca{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.ProdutoMarca{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Select(`
			promar.*,
			(
				SELECT COUNT(*) 
				FROM produto 
				WHERE promar_id = produto.promar_id 
				AND deleted_at IS NULL
			) as total_produtos
		`).
		Limit(limit).
		Offset(offset).
		Order("promar_descricao ASC").
		Find(&marcas).Error

	if err != nil {
		return nil, 0, err
	}

	return marcas, total, nil
}

// FindAll busca todas as marcas de produto
func (r *produtoMarcaRepository) FindAll() ([]models.ProdutoMarca, error) {
	var marcas []models.ProdutoMarca
	err := r.db.
		Where("deleted_at IS NULL").
		Order("promar_descricao ASC").
		Find(&marcas).Error
	if err != nil {
		return nil, err
	}
	return marcas, nil
}

// ============================================================
// MÉTODOS DE CONSULTA PARA VALIDAÇÕES (APENAS CONSULTAS)
// ============================================================

// ExistsByDescricao verifica se já existe uma marca com a descrição
func (r *produtoMarcaRepository) ExistsByDescricao(descricao string, excludeID int) (bool, error) {
	if descricao == "" {
		return false, nil
	}

	var count int64
	query := r.db.Model(&models.ProdutoMarca{}).
		Where("promar_descricao = ? AND deleted_at IS NULL", descricao)

	if excludeID > 0 {
		query = query.Where("promar_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByID verifica se uma marca de produto existe pelo ID
func (r *produtoMarcaRepository) ExistsByID(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.ProdutoMarca{}).
		Where("promar_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Count retorna o total de marcas de produto
func (r *produtoMarcaRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.ProdutoMarca{}).
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CountProdutosByMarca retorna a quantidade de produtos de uma marca
func (r *produtoMarcaRepository) CountProdutosByMarca(marcaID int) (int64, error) {
	var count int64
	err := r.db.Model(&models.Produto{}).
		Where("promar_id = ? AND deleted_at IS NULL", marcaID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// ============================================================
// OPERAÇÕES EM LOTE
// ============================================================

// BulkUpdateStatus atualiza a situação de múltiplas marcas
func (r *produtoMarcaRepository) BulkUpdateStatus(ids []int, situacao int) error {
	return r.db.Model(&models.ProdutoMarca{}).
		Where("promar_id IN ? AND deleted_at IS NULL", ids).
		Update("promar_situacao", situacao).Error
}

// BulkDelete realiza exclusão lógica de múltiplas marcas
func (r *produtoMarcaRepository) BulkDelete(ids []int) error {
	return r.db.Model(&models.ProdutoMarca{}).
		Where("promar_id IN ? AND deleted_at IS NULL", ids).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// ============================================================
// MÉTODOS DE CONSULTA DE DEPENDÊNCIAS
// ============================================================

// HasDependentRecords verifica se a marca tem registros dependentes
func (r *produtoMarcaRepository) HasDependentRecords(id int) (bool, error) {
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
func (r *produtoMarcaRepository) CountDependentRecords(id int) (map[string]int64, error) {
	result := make(map[string]int64)

	// Verifica produtos associados
	var countProdutos int64
	if err := r.db.Model(&models.Produto{}).
		Where("promar_id = ? AND deleted_at IS NULL", id).
		Count(&countProdutos).Error; err != nil {
		return nil, err
	}
	if countProdutos > 0 {
		result["produtos"] = countProdutos
	}

	return result, nil
}
