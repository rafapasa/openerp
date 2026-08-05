// internal/repository/entidade_contato_repository.go
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

// EntidadeContatoRepository define o contrato para operações de banco
type EntidadeContatoRepository interface {
	// CRUD Básico
	Create(contato *models.EntidadeContato) error
	Update(contato *models.EntidadeContato) error
	Delete(entidadeID, item int) error
	FindByID(entidadeID, item int) (*models.EntidadeContato, error)

	// Buscas Específicas
	FindByEntidadeID(entidadeID int) ([]models.EntidadeContato, error)
	FindByEntidadeIDAndTipo(entidadeID, formaContatoID int) ([]models.EntidadeContato, error)

	// Listagem com Filtros
	List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeContato, int64, error)

	// Consultas de Validação (APENAS CONSULTAS)
	GetNextItemNumber(entidadeID int) (int, error)
	CountByEntidadeID(entidadeID int) (int64, error)
	ExistsByEntidadeTipo(entidadeID, formaContatoID int, excludeItem int) (bool, error)
}

// ============================================================
// IMPLEMENTAÇÃO CONCRETA (privada)
// ============================================================

type entidadeContatoRepository struct {
	db *gorm.DB
}

// NewEntidadeContatoRepository cria uma nova instância do repositório
// ✅ Retorna a interface, não a struct concreta
func NewEntidadeContatoRepository(db *gorm.DB) EntidadeContatoRepository {
	return &entidadeContatoRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD (APENAS PERSISTÊNCIA)
// ============================================================

// Create salva um novo contato com sequencial manual
func (r *entidadeContatoRepository) Create(contato *models.EntidadeContato) error {
	err := r.db.Create(contato).Error
	if err != nil {
		return apperrors.NewInternalError("Erro criando contato: ", err)
	}

	return nil
}

// Update atualiza um contato existente
func (r *entidadeContatoRepository) Update(contato *models.EntidadeContato) error {
	err := r.db.
		Omit("FormaContato", "created_at", "deleted_at").
		Model(&models.EntidadeContato{}).
		Where("ent_id = ? AND efc_item = ?", contato.EntidadeID, contato.Item).
		Updates(contato).Error
	if err != nil {
		return apperrors.NewInternalError("Erro atualizando contato: ", err)
	}
	return nil

}

// Delete realiza exclusão lógica de um contato
func (r *entidadeContatoRepository) Delete(entidadeID, item int) error {
	err := r.db.
		Model(&models.EntidadeContato{}).
		Where("ent_id = ? AND efc_item = ?", entidadeID, item).
		Update("deleted_at", gorm.Expr("NOW()")).Error
	if err != nil {
		return apperrors.NewInternalError("Erro deletando contato: ", err)
	}
	return nil
}

// ============================================================
// MÉTODOS DE BUSCA
// ============================================================

// FindByID busca um contato pelo ID composto (ent_id + efc_item)
func (r *entidadeContatoRepository) FindByID(entidadeID, item int) (*models.EntidadeContato, error) {
	var contato models.EntidadeContato
	err := r.db.
		Preload("FormaContato").
		Where("ent_id = ? AND efc_item = ? AND deleted_at IS NULL", entidadeID, item).
		First(&contato).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("Contato para entidade %d item %d não encontrado", entidadeID, item))
		}
		return nil, apperrors.NewInternalError("Erro buscando contato: ", err)
	}
	return &contato, nil
}

// FindByEntidadeID busca todos os contatos de uma entidade
func (r *entidadeContatoRepository) FindByEntidadeID(entidadeID int) ([]models.EntidadeContato, error) {
	var contatos []models.EntidadeContato
	err := r.db.
		Preload("FormaContato").
		Where("ent_id = ? AND deleted_at IS NULL", entidadeID).
		Order("efc_item ASC").
		Find(&contatos).Error

	if err != nil {
		return nil, apperrors.NewInternalError("Erro buscando contato: ", err)
	}
	return contatos, nil
}

// FindByEntidadeIDAndTipo busca contatos de uma entidade por tipo
func (r *entidadeContatoRepository) FindByEntidadeIDAndTipo(entidadeID, formaContatoID int) ([]models.EntidadeContato, error) {
	var contatos []models.EntidadeContato
	err := r.db.
		Preload("FormaContato").
		Where("ent_id = ? AND frc_id = ? AND deleted_at IS NULL", entidadeID, formaContatoID).
		Order("efc_item ASC").
		Find(&contatos).Error

	if err != nil {
		return nil, apperrors.NewInternalError("Erro buscando contato: ", err)
	}
	return contatos, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de contatos com paginação e filtros
func (r *entidadeContatoRepository) List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeContato, int64, error) {
	var contatos []models.EntidadeContato
	var total int64

	query := r.db.Model(&models.EntidadeContato{}).Where("deleted_at IS NULL")
	query = r.applyFilters(query, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("Erro buscando contato: ", err)
	}

	err := query.
		Preload("FormaContato").
		Limit(limit).
		Offset(offset).
		Order("efc_item DESC").
		Find(&contatos).Error

	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro buscando contato: ", err)
	}

	return contatos, total, nil
}

// ============================================================
// MÉTODOS DE CONSULTA PARA VALIDAÇÕES (APENAS CONSULTAS)
// ============================================================

// GetNextItemNumber retorna o próximo número de item para uma entidade
func (r *entidadeContatoRepository) GetNextItemNumber(entidadeID int) (int, error) {
	var maxItem int
	err := r.db.Model(&models.EntidadeContato{}).
		Where("ent_id = ?", entidadeID).
		Select("COALESCE(MAX(efc_item), 0) + 1").
		Scan(&maxItem).Error
	if err != nil {
		return 0, apperrors.NewInternalError("Erro gerando novo proximo item: ", err)
	}
	return maxItem, nil
}

// CountByEntidadeID retorna a quantidade de contatos de uma entidade
func (r *entidadeContatoRepository) CountByEntidadeID(entidadeID int) (int64, error) {
	var count int64
	err := r.db.Model(&models.EntidadeContato{}).
		Where("ent_id = ? AND deleted_at IS NULL", entidadeID).
		Count(&count).Error
	if err != nil {
		return 0, apperrors.NewInternalError("Erro buscando contato: ", err)
	}
	return count, nil
}

// ExistsByEntidadeTipo verifica se a entidade já tem um contato do tipo especificado
func (r *entidadeContatoRepository) ExistsByEntidadeTipo(entidadeID, formaContatoID int, excludeItem int) (bool, error) {
	var count int64
	query := r.db.Model(&models.EntidadeContato{}).
		Where("ent_id = ? AND frc_id = ? AND deleted_at IS NULL", entidadeID, formaContatoID)

	if excludeItem > 0 {
		query = query.Where("efc_item != ?", excludeItem)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, apperrors.NewInternalError("Erro buscando contato: ", err)
	}
	return count > 0, nil
}

// ============================================================
// MÉTODOS AUXILIARES (FILTROS)
// ============================================================

// applyFilters aplica os filtros à query
func (r *entidadeContatoRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for key, value := range filters {
		if value == nil || value == "" {
			continue
		}

		switch key {
		case "entidade_id":
			query = query.Where("ent_id = ?", value)
		case "forma_contato_id":
			query = query.Where("frc_id = ?", value)
		case "informacao":
			query = query.Where("efc_informacao LIKE ?", "%"+value.(string)+"%")
		}
	}
	return query
}
