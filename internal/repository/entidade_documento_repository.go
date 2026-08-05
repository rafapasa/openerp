// internal/repository/entidade_documento_repository.go
package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/models"
)

// ============================================================
// INTERFACE
// ============================================================

// EntidadeDocumentoRepository define o contrato para operações de banco
type EntidadeDocumentoRepository interface {
	// CRUD Básico
	Create(documento *models.EntidadeDocumento) error
	Update(documento *models.EntidadeDocumento) error
	Delete(entidadeID, item int) error
	FindByID(entidadeID, item int) (*models.EntidadeDocumento, error)

	// Buscas Específicas
	FindByEntidadeID(entidadeID int) ([]models.EntidadeDocumento, error)
	FindByEntidadeIDAndTipo(entidadeID int, tipo string) ([]models.EntidadeDocumento, error)

	// Listagem com Filtros
	List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeDocumento, int64, error)

	// Consultas de Validação (APENAS CONSULTAS)
	GetNextItemNumber(entidadeID int) (int, error)
	CountByEntidadeID(entidadeID int) (int64, error)
	ExistsByEntidadeIDAndTipo(entidadeID int, tipo string, excludeItem int) (bool, error)
}

// ============================================================
// IMPLEMENTAÇÃO CONCRETA (privada)
// ============================================================

// entidadeDocumentoRepository é a implementação concreta
type entidadeDocumentoRepository struct {
	db *gorm.DB
}

// NewEntidadeDocumentoRepository cria uma nova instância (retorna a interface)
func NewEntidadeDocumentoRepository(db *gorm.DB) EntidadeDocumentoRepository {
	return &entidadeDocumentoRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD (APENAS PERSISTÊNCIA)
// ============================================================

// Create salva um novo documento com sequencial manual
func (r *entidadeDocumentoRepository) Create(documento *models.EntidadeDocumento) error {
	return r.db.Create(documento).Error
}

// Update atualiza um documento existente
func (r *entidadeDocumentoRepository) Update(documento *models.EntidadeDocumento) error {
	return r.db.
		Omit("created_at", "deleted_at").
		Save(documento).Error
}

// Delete realiza exclusão lógica de um documento
func (r *entidadeDocumentoRepository) Delete(entidadeID, item int) error {
	return r.db.
		Model(&models.EntidadeDocumento{}).
		Where("ent_id = ? AND edoc_item = ?", entidadeID, item).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// ============================================================
// MÉTODOS DE BUSCA
// ============================================================

// FindByID busca um documento pelo ID composto (ent_id + edoc_item)
func (r *entidadeDocumentoRepository) FindByID(entidadeID, item int) (*models.EntidadeDocumento, error) {
	var documento models.EntidadeDocumento
	err := r.db.
		Where("ent_id = ? AND edoc_item = ? AND deleted_at IS NULL", entidadeID, item).
		First(&documento).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("documento para entidade %d item %d não encontrado", entidadeID, item))
		}
		return nil, apperrors.NewInternalError("Erro buscando documento: ", err)
	}
	return &documento, nil
}

// FindByEntidadeID busca todos os documentos de uma entidade
func (r *entidadeDocumentoRepository) FindByEntidadeID(entidadeID int) ([]models.EntidadeDocumento, error) {
	var documentos []models.EntidadeDocumento
	err := r.db.
		Where("ent_id = ? AND deleted_at IS NULL", entidadeID).
		Order("edoc_item ASC").
		Find(&documentos).Error

	if err != nil {
		return nil, err
	}
	return documentos, nil
}

// FindByEntidadeIDAndTipo busca documentos de uma entidade por tipo
func (r *entidadeDocumentoRepository) FindByEntidadeIDAndTipo(entidadeID int, tipo string) ([]models.EntidadeDocumento, error) {
	var documentos []models.EntidadeDocumento
	err := r.db.
		Where("ent_id = ? AND edoc_tipo = ? AND deleted_at IS NULL", entidadeID, tipo).
		Order("edoc_item ASC").
		Find(&documentos).Error

	if err != nil {
		return nil, err
	}
	return documentos, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de documentos com paginação e filtros
func (r *entidadeDocumentoRepository) List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeDocumento, int64, error) {
	var documentos []models.EntidadeDocumento
	var total int64

	query := r.db.Model(&models.EntidadeDocumento{}).Where("deleted_at IS NULL")
	query = r.applyFilters(query, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Order("edoc_item DESC").
		Find(&documentos).Error

	if err != nil {
		return nil, 0, err
	}

	return documentos, total, nil
}

// ============================================================
// MÉTODOS DE CONSULTA PARA VALIDAÇÕES (APENAS CONSULTAS)
// ============================================================

// GetNextItemNumber retorna o próximo número de item para uma entidade
func (r *entidadeDocumentoRepository) GetNextItemNumber(entidadeID int) (int, error) {
	var maxItem int
	err := r.db.Model(&models.EntidadeDocumento{}).
		Where("ent_id = ?", entidadeID).
		Select("COALESCE(MAX(edoc_item), 0) + 1").
		Scan(&maxItem).Error
	if err != nil {
		return 0, err
	}
	return maxItem, nil
}

// CountByEntidadeID retorna a quantidade de documentos de uma entidade
func (r *entidadeDocumentoRepository) CountByEntidadeID(entidadeID int) (int64, error) {
	var count int64
	err := r.db.Model(&models.EntidadeDocumento{}).
		Where("ent_id = ? AND deleted_at IS NULL", entidadeID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// ExistsByEntidadeIDAndTipo verifica se a entidade já possui um documento do tipo especificado
func (r *entidadeDocumentoRepository) ExistsByEntidadeIDAndTipo(entidadeID int, tipo string, excludeItem int) (bool, error) {
	var count int64
	query := r.db.Model(&models.EntidadeDocumento{}).
		Where("ent_id = ? AND edoc_tipo = ? AND deleted_at IS NULL", entidadeID, tipo)

	if excludeItem > 0 {
		query = query.Where("edoc_item != ?", excludeItem)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ============================================================
// MÉTODOS AUXILIARES (FILTROS)
// ============================================================

// applyFilters aplica os filtros à query
func (r *entidadeDocumentoRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for key, value := range filters {
		if value == nil || value == "" {
			continue
		}

		switch key {
		case "entidade_id":
			query = query.Where("ent_id = ?", value)
		case "tipo":
			query = query.Where("edoc_tipo = ?", value)
		case "descricao":
			query = query.Where("edoc_descricao LIKE ?", "%"+value.(string)+"%")
		}
	}
	return query
}
