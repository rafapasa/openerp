// internal/repository/entidade_limite_credito_repository.go
package repository

import (
	"errors"
	"fmt"

	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
	"gorm.io/gorm"
)

// ============================================================
// INTERFACE (Contrato)
// ============================================================

// EntidadeLimiteCreditoRepository define o contrato para operações de banco
type EntidadeLimiteCreditoRepository interface {
	// CRUD Básico
	Create(limite *models.EntidadeLimiteCredito) error
	Update(id int, limite *models.EntidadeLimiteCredito) error
	Delete(id int) error
	FindByID(id int) (*models.EntidadeLimiteCredito, error)

	// Buscas Específicas
	FindByEntidadeID(entidadeID int) ([]models.EntidadeLimiteCredito, error)
	FindActiveByEntidadeID(entidadeID int) (*models.EntidadeLimiteCredito, error)

	// Listagem com Filtros
	List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeLimiteCredito, int64, error)

	// Consultas de Validação (APENAS CONSULTAS)
	ExistsByDescricao(descricao string, excludeID int) (bool, error)
	ExistsByEntidadeID(entidadeID int) (bool, error)
	CountByEntidadeID(entidadeID int) (int64, error)
	HasActiveLimit(entidadeID int, excludeID int) (bool, error)
}

// ============================================================
// IMPLEMENTAÇÃO CONCRETA (privada)
// ============================================================

// entidadeLimiteCreditoRepository é a implementação concreta
type entidadeLimiteCreditoRepository struct {
	db *gorm.DB
}

// NewEntidadeLimiteCreditoRepository cria uma nova instância (retorna a interface)
func NewEntidadeLimiteCreditoRepository(db *gorm.DB) EntidadeLimiteCreditoRepository {
	return &entidadeLimiteCreditoRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD (APENAS PERSISTÊNCIA)
// ============================================================

// Create salva um novo limite de crédito
func (r *entidadeLimiteCreditoRepository) Create(limite *models.EntidadeLimiteCredito) error {
	return r.db.Create(limite).Error
}

// Update atualiza um limite de crédito existente
func (r *entidadeLimiteCreditoRepository) Update(id int, limite *models.EntidadeLimiteCredito) error {
	return r.db.Model(&models.EntidadeLimiteCredito{}).
		Where("elc_id = ?", id).
		Omit("created_at", "deleted_at").
		Updates(limite).Error
}

// Delete realiza a exclusão lógica
func (r *entidadeLimiteCreditoRepository) Delete(id int) error {
	return r.db.Model(&models.EntidadeLimiteCredito{}).
		Where("elc_id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// ============================================================
// MÉTODOS DE BUSCA
// ============================================================

// FindByID busca um limite de crédito pelo ID
func (r *entidadeLimiteCreditoRepository) FindByID(id int) (*models.EntidadeLimiteCredito, error) {
	var limite models.EntidadeLimiteCredito
	err := r.db.Where("elc_id = ? AND deleted_at IS NULL", id).First(&limite).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("Limite de crédito com ID %d não encontrado.", id))
		}
		return nil, err
	}
	return &limite, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de limites com paginação e filtros
func (r *entidadeLimiteCreditoRepository) List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeLimiteCredito, int64, error) {
	var limites []models.EntidadeLimiteCredito
	var total int64

	query := r.db.Model(&models.EntidadeLimiteCredito{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.EntidadeLimiteCredito{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Order("elc_descricao ASC").Find(&limites).Error
	if err != nil {
		return nil, 0, err
	}

	return limites, total, nil
}

// ============================================================
// MÉTODOS ESPECÍFICOS DO DOMÍNIO
// ============================================================

// FindByEntidadeID busca todos os limites de crédito de uma entidade
func (r *entidadeLimiteCreditoRepository) FindByEntidadeID(entidadeID int) ([]models.EntidadeLimiteCredito, error) {
	var limites []models.EntidadeLimiteCredito
	err := r.db.Where("entidade_id = ? AND deleted_at IS NULL", entidadeID).
		Order("created_at DESC").
		Find(&limites).Error
	if err != nil {
		return nil, err
	}
	return limites, nil
}

// FindActiveByEntidadeID busca o limite ativo de uma entidade
func (r *entidadeLimiteCreditoRepository) FindActiveByEntidadeID(entidadeID int) (*models.EntidadeLimiteCredito, error) {
	var limite models.EntidadeLimiteCredito
	err := r.db.Where("entidade_id = ? AND ativo = true AND deleted_at IS NULL", entidadeID).
		Order("created_at DESC").
		First(&limite).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Nenhum limite ativo encontrado
		}
		return nil, err
	}
	return &limite, nil
}

// ============================================================
// MÉTODOS DE CONSULTA PARA VALIDAÇÕES (APENAS CONSULTAS)
// ============================================================

// ExistsByDescricao verifica se já existe um limite com a mesma descrição
func (r *entidadeLimiteCreditoRepository) ExistsByDescricao(descricao string, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.EntidadeLimiteCredito{}).
		Where("elc_descricao = ? AND deleted_at IS NULL", descricao)

	if excludeID > 0 {
		query = query.Where("elc_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// ExistsByEntidadeID verifica se uma entidade possui algum limite de crédito
func (r *entidadeLimiteCreditoRepository) ExistsByEntidadeID(entidadeID int) (bool, error) {
	var count int64
	err := r.db.Model(&models.EntidadeLimiteCredito{}).
		Where("entidade_id = ? AND deleted_at IS NULL", entidadeID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountByEntidadeID retorna a quantidade de limites de uma entidade
func (r *entidadeLimiteCreditoRepository) CountByEntidadeID(entidadeID int) (int64, error) {
	var count int64
	err := r.db.Model(&models.EntidadeLimiteCredito{}).
		Where("entidade_id = ? AND deleted_at IS NULL", entidadeID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// HasActiveLimit verifica se a entidade já possui um limite ativo
func (r *entidadeLimiteCreditoRepository) HasActiveLimit(entidadeID int, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.EntidadeLimiteCredito{}).
		Where("entidade_id = ? AND ativo = true AND deleted_at IS NULL", entidadeID)

	if excludeID > 0 {
		query = query.Where("elc_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
