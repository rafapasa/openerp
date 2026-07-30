// internal/repository/condicao_pagamento_repository.go
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

// CondicaoPagamentoRepository define o contrato para operações de banco
type CondicaoPagamentoRepository interface {
	// CRUD Básico
	Create(condicaoPagamento *models.CondicaoPagamento) error
	Update(id int, condicaoPagamento *models.CondicaoPagamento) error
	Delete(id int) error
	GetByID(id int) (*models.CondicaoPagamento, error)
	FindByID(id int) (*models.CondicaoPagamento, error)

	// Buscas Específicas
	FindByTipoDocumento(id string) (*models.CondicaoPagamento, error)
	FindByDescricao(nome string, limit int) ([]models.CondicaoPagamento, error)

	// Listagem com Filtros
	List(limit, offset int, filters map[string]interface{}) ([]models.CondicaoPagamento, int64, error)

	// Consultas de Validação (APENAS CONSULTAS)
	ExistsByTipoDocumento(id string, excludeID int) (bool, error)
	ExistsByDescricao(descricao string, excludeID int) (bool, error)
	ExistsByID(id int) (bool, error)
	Count() (int64, error)
}

// ============================================================
// IMPLEMENTAÇÃO CONCRETA (privada)
// ============================================================

// condicaoPagamentoRepository é a implementação concreta
type condicaoPagamentoRepository struct {
	db *gorm.DB
}

// NewCondicaoPagamentoRepository cria uma nova instância (retorna a interface)
func NewCondicaoPagamentoRepository(db *gorm.DB) CondicaoPagamentoRepository {
	return &condicaoPagamentoRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD (APENAS PERSISTÊNCIA)
// ============================================================

// Create salva uma nova CondicaoPagamento
func (r *condicaoPagamentoRepository) Create(condicaoPagamento *models.CondicaoPagamento) error {
	return r.db.Create(condicaoPagamento).Error
}

// Update atualiza uma CondicaoPagamento existente
func (r *condicaoPagamentoRepository) Update(id int, condicaoPagamento *models.CondicaoPagamento) error {
	return r.db.
		Omit("Tipo_Documento", "Forma_Pagamento", "Portador", "created_at", "deleted_at").
		Model(&models.CondicaoPagamento{}).
		Where("cdpgt_id = ?", id).
		Updates(condicaoPagamento).Error
}

// Delete realiza exclusão lógica de uma CondicaoPagamento pelo ID
func (r *condicaoPagamentoRepository) Delete(id int) error {
	return r.db.
		Model(&models.CondicaoPagamento{}).
		Where("cdpgt_id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// GetByID busca uma CondicaoPagamento pelo ID (sem relacionamentos)
func (r *condicaoPagamentoRepository) GetByID(id int) (*models.CondicaoPagamento, error) {
	var condicaoPagamento models.CondicaoPagamento
	result := r.db.
		Where("cdpgt_id = ? AND deleted_at IS NULL", id).
		First(&condicaoPagamento)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("CondicaoPagamento com ID %d não encontrada", id))
		}
		return nil, apperrors.NewInternalError("Erro buscando CondicaoPagamento", result.Error)
	}
	return &condicaoPagamento, nil
}

// ============================================================
// MÉTODOS DE BUSCA
// ============================================================

// FindByID busca uma CondicaoPagamento pelo ID com relacionamentos
func (r *condicaoPagamentoRepository) FindByID(id int) (*models.CondicaoPagamento, error) {
	var condicaoPagamento models.CondicaoPagamento
	err := r.db.
		Preload("TipoDocumento").
		Preload("FormaPagamento").
		Preload("Portador").
		Where("cdpgt_id = ? AND deleted_at IS NULL", id).
		First(&condicaoPagamento).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("CondicaoPagamento com ID %d não encontrada", id))
		}
		return nil, err
	}
	return &condicaoPagamento, nil
}

// FindByTipoDocumento busca uma CondicaoPagamento pelo TipoDocumento
func (r *condicaoPagamentoRepository) FindByTipoDocumento(id string) (*models.CondicaoPagamento, error) {
	var condicaoPagamento models.CondicaoPagamento
	err := r.db.
		Preload("TipoDocumento").
		Preload("FormaPagamento").
		Preload("Portador").
		Where("tdoc_id = ? AND deleted_at IS NULL", id).
		First(&condicaoPagamento).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("CondicaoPagamento não encontrada")
		}
		return nil, err
	}
	return &condicaoPagamento, nil
}

// FindByDescricao busca CondicaoPagamento pela descrição (autocomplete)
func (r *condicaoPagamentoRepository) FindByDescricao(nome string, limit int) ([]models.CondicaoPagamento, error) {
	var condicoes []models.CondicaoPagamento
	err := r.db.
		Where("cdpgt_descricao LIKE ? AND deleted_at IS NULL", "%"+nome+"%").
		Limit(limit).
		Order("cdpgt_descricao ASC").
		Find(&condicoes).Error

	if err != nil {
		return nil, err
	}
	return condicoes, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de CondicaoPagamento com paginação e filtros
func (r *condicaoPagamentoRepository) List(limit, offset int, filters map[string]interface{}) ([]models.CondicaoPagamento, int64, error) {
	var condicoes []models.CondicaoPagamento
	var total int64

	query := r.db.Model(&models.CondicaoPagamento{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.CondicaoPagamento{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("TipoDocumento").
		Preload("FormaPagamento").
		Preload("Portador").
		Limit(limit).
		Offset(offset).
		Order("cdpgt_id DESC").
		Find(&condicoes).Error

	if err != nil {
		return nil, 0, err
	}

	return condicoes, total, nil
}

// ============================================================
// MÉTODOS DE CONSULTA PARA VALIDAÇÕES (APENAS CONSULTAS)
// ============================================================

// ExistsByTipoDocumento verifica se já existe uma CondicaoPagamento com o TipoDocumento
func (r *condicaoPagamentoRepository) ExistsByTipoDocumento(id string, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.CondicaoPagamento{}).
		Where("tdoc_id = ? AND deleted_at IS NULL", id)

	if excludeID > 0 {
		query = query.Where("cdpgt_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByDescricao verifica se já existe uma CondicaoPagamento com a descrição
func (r *condicaoPagamentoRepository) ExistsByDescricao(descricao string, excludeID int) (bool, error) {
	if descricao == "" {
		return false, nil
	}

	var count int64
	query := r.db.Model(&models.CondicaoPagamento{}).
		Where("cdpgt_descricao = ? AND deleted_at IS NULL", descricao)

	if excludeID > 0 {
		query = query.Where("cdpgt_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByID verifica se uma CondicaoPagamento existe pelo ID
func (r *condicaoPagamentoRepository) ExistsByID(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.CondicaoPagamento{}).
		Where("cdpgt_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Count retorna o total de CondicoesPagamento
func (r *condicaoPagamentoRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.CondicaoPagamento{}).
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
