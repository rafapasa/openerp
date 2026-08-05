// internal/repository/produto_modelo_repository.go
package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// INTERFACE
// ============================================================

// ProdutoModeloRepository define o contrato para operações de banco
type ProdutoModeloRepository interface {
	// CRUD Básico
	Create(modelo *models.ProdutoModelo) error
	Update(id int, modelo *models.ProdutoModelo) error
	Delete(id int) error
	FindByID(id int) (*models.ProdutoModelo, error)
	GetByID(id int) (*models.ProdutoModelo, error)

	// Buscas Específicas
	FindByDescricao(descricao string, limit int) ([]models.ProdutoModelo, error)
	FindBySituacao(situacao int) ([]models.ProdutoModelo, error)
	FindActive() ([]models.ProdutoModelo, error)

	// Listagem com Filtros
	List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoModelo, int64, error)
	ListWithProdutos(limit, offset int, filters map[string]interface{}) ([]models.ProdutoModelo, int64, error)
	FindAll() ([]models.ProdutoModelo, error)

	// Consultas de Validação (APENAS CONSULTAS)
	ExistsByDescricao(descricao string, excludeID int) (bool, error)
	ExistsByID(id int) (bool, error)
	Count() (int64, error)
	CountProdutosByModelo(modeloID int) (int64, error)

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

type produtoModeloRepository struct {
	db *gorm.DB
}

// NewProdutoModeloRepository cria uma nova instância (retorna a interface)
func NewProdutoModeloRepository(db *gorm.DB) ProdutoModeloRepository {
	return &produtoModeloRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD (APENAS PERSISTÊNCIA)
// ============================================================

// Create salva um novo modelo de produto
func (r *produtoModeloRepository) Create(modelo *models.ProdutoModelo) error {
	err := r.db.Create(modelo).Error
	if err != nil {
		return apperrors.NewInternalError("Erro ao criar modelo de produto.", err)
	}
	return nil
}

// Update atualiza um modelo de produto existente
func (r *produtoModeloRepository) Update(id int, modelo *models.ProdutoModelo) error {
	err := r.db.
		Omit("Produtos", "created_at", "deleted_at").
		Model(&models.ProdutoModelo{}).
		Where("prom_id = ?", id).
		Updates(modelo).Error
	if err != nil {
		return apperrors.NewInternalError("Erro ao atualizar modelo de produto.", err)
	}
	return nil

}

// Delete realiza exclusão lógica
func (r *produtoModeloRepository) Delete(id int) error {
	err := r.db.
		Model(&models.ProdutoModelo{}).
		Where("prom_id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
	if err != nil {
		return apperrors.NewInternalError("Erro ao excluir modelo de produto.", err)
	}
	return nil
}

// FindByID busca um modelo de produto pelo ID com relacionamentos
func (r *produtoModeloRepository) FindByID(id int) (*models.ProdutoModelo, error) {
	var modelo models.ProdutoModelo
	err := r.db.
		Preload("Produtos", "deleted_at IS NULL").
		Where("prom_id = ? AND deleted_at IS NULL", id).
		First(&modelo).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("modelo de produto com ID %d não encontrado", id))
		}
		return nil, apperrors.NewInternalError("Erro buscando marca de produtos.", err)
	}
	return &modelo, nil
}

// GetByID busca um modelo de produto pelo ID (sem relacionamentos)
func (r *produtoModeloRepository) GetByID(id int) (*models.ProdutoModelo, error) {
	var modelo models.ProdutoModelo
	err := r.db.
		Where("prom_id = ? AND deleted_at IS NULL", id).
		First(&modelo).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("modelo de produto com ID %d não encontrado", id))
		}
		return nil, apperrors.NewInternalError("Erro buscando marca de produtos.", err)
	}
	return &modelo, nil
}

// ============================================================
// MÉTODOS DE BUSCA ESPECÍFICOS
// ============================================================

// FindByDescricao busca modelos de produto pela descrição (autocomplete)
func (r *produtoModeloRepository) FindByDescricao(descricao string, limit int) ([]models.ProdutoModelo, error) {
	var modelos []models.ProdutoModelo
	err := r.db.
		Where("prom_descricao LIKE ? AND deleted_at IS NULL", "%"+descricao+"%").
		Order("prom_descricao ASC").
		Limit(limit).
		Find(&modelos).Error

	if err != nil {
		return nil, apperrors.NewInternalError("Erro buscando marca de produtos.", err)
	}
	return modelos, nil
}

// FindBySituacao busca modelos de produto por situação
func (r *produtoModeloRepository) FindBySituacao(situacao int) ([]models.ProdutoModelo, error) {
	var modelos []models.ProdutoModelo
	err := r.db.
		Where("prom_situacao = ? AND deleted_at IS NULL", situacao).
		Order("prom_descricao ASC").
		Find(&modelos).Error

	if err != nil {
		return nil, apperrors.NewInternalError("Erro buscando marca de produtos.", err)
	}
	return modelos, nil
}

// FindActive busca todos os modelos de produto ativos
func (r *produtoModeloRepository) FindActive() ([]models.ProdutoModelo, error) {
	var modelos []models.ProdutoModelo
	err := r.db.
		Where("prom_situacao = 1 AND deleted_at IS NULL").
		Order("prom_descricao ASC").
		Find(&modelos).Error

	if err != nil {
		return nil, apperrors.NewInternalError("Erro buscando marca de produtos.", err)
	}
	return modelos, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de modelos de produto com paginação e filtros
func (r *produtoModeloRepository) List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoModelo, int64, error) {
	var modelos []models.ProdutoModelo
	var total int64

	query := r.db.Model(&models.ProdutoModelo{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.ProdutoModelo{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("Erro buscando marca de produtos.", err)
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Order("prom_descricao ASC").
		Find(&modelos).Error

	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro buscando marca de produtos.", err)
	}

	return modelos, total, nil
}

// ListWithProdutos retorna modelos de produto com contagem de produtos
func (r *produtoModeloRepository) ListWithProdutos(limit, offset int, filters map[string]interface{}) ([]models.ProdutoModelo, int64, error) {
	var modelos []models.ProdutoModelo
	var total int64

	query := r.db.Model(&models.ProdutoModelo{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.ProdutoModelo{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("Erro buscando marca de produtos.", err)
	}

	err := query.
		Select(`
			prom.*,
			(
				SELECT COUNT(*) 
				FROM produto 
				WHERE prom_id = produto.prom_id 
				AND deleted_at IS NULL
			) as total_produtos
		`).
		Limit(limit).
		Offset(offset).
		Order("prom_descricao ASC").
		Find(&modelos).Error

	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro buscando marca de produtos.", err)
	}

	return modelos, total, nil
}

// FindAll busca todos os modelos de produto
func (r *produtoModeloRepository) FindAll() ([]models.ProdutoModelo, error) {
	var modelos []models.ProdutoModelo
	err := r.db.
		Where("deleted_at IS NULL").
		Order("prom_descricao ASC").
		Find(&modelos).Error
	if err != nil {
		return nil, apperrors.NewInternalError("Erro buscando marca de produtos.", err)
	}
	return modelos, nil
}

// ============================================================
// MÉTODOS DE CONSULTA PARA VALIDAÇÕES (APENAS CONSULTAS)
// ============================================================

// ExistsByDescricao verifica se já existe um modelo com a descrição
func (r *produtoModeloRepository) ExistsByDescricao(descricao string, excludeID int) (bool, error) {
	if descricao == "" {
		return false, nil
	}

	var count int64
	query := r.db.Model(&models.ProdutoModelo{}).
		Where("prom_descricao = ? AND deleted_at IS NULL", descricao)

	if excludeID > 0 {
		query = query.Where("prom_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, apperrors.NewInternalError("Erro buscando marca de produtos.", err)
	}
	return count > 0, nil
}

// ExistsByID verifica se um modelo de produto existe pelo ID
func (r *produtoModeloRepository) ExistsByID(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.ProdutoModelo{}).
		Where("prom_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	if err != nil {
		return false, apperrors.NewInternalError("Erro buscando marca de produtos.", err)
	}
	return count > 0, nil
}

// Count retorna o total de modelos de produto
func (r *produtoModeloRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.ProdutoModelo{}).
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, apperrors.NewInternalError("Erro buscando marca de produtos.", err)
	}
	return count, nil
}

// CountProdutosByModelo retorna a quantidade de produtos de um modelo
func (r *produtoModeloRepository) CountProdutosByModelo(modeloID int) (int64, error) {
	var count int64
	err := r.db.Model(&models.Produto{}).
		Where("prom_id = ? AND deleted_at IS NULL", modeloID).
		Count(&count).Error
	if err != nil {
		return 0, apperrors.NewInternalError("Erro buscando marca de produtos.", err)
	}
	return count, nil
}

// ============================================================
// OPERAÇÕES EM LOTE
// ============================================================

// BulkUpdateStatus atualiza a situação de múltiplos modelos
func (r *produtoModeloRepository) BulkUpdateStatus(ids []int, situacao int) error {
	err := r.db.Model(&models.ProdutoModelo{}).
		Where("prom_id IN ? AND deleted_at IS NULL", ids).
		Update("prom_situacao", situacao).Error
	if err != nil {
		return err
	}
	return nil

}

// BulkDelete realiza exclusão lógica de múltiplos modelos
func (r *produtoModeloRepository) BulkDelete(ids []int) error {
	err := r.db.Model(&models.ProdutoModelo{}).
		Where("prom_id IN ? AND deleted_at IS NULL", ids).
		Update("deleted_at", gorm.Expr("NOW()")).Error
	if err != nil {
		return err
	}
	return nil
}

// ============================================================
// MÉTODOS DE CONSULTA DE DEPENDÊNCIAS
// ============================================================

// HasDependentRecords verifica se o modelo tem registros dependentes
func (r *produtoModeloRepository) HasDependentRecords(id int) (bool, error) {
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
func (r *produtoModeloRepository) CountDependentRecords(id int) (map[string]int64, error) {
	result := make(map[string]int64)

	// Verifica produtos associados
	var countProdutos int64
	if err := r.db.Model(&models.Produto{}).
		Where("prom_id = ? AND deleted_at IS NULL", id).
		Count(&countProdutos).Error; err != nil {
		return nil, err
	}
	if countProdutos > 0 {
		result["produtos"] = countProdutos
	}

	return result, nil
}
