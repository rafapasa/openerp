// internal/repository/entidade_regime_tributario_repository.go
package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// INTERFACE (Contrato)
// ============================================================

// EntidadeRegimeTributarioRepository define o contrato para operações de banco
type EntidadeRegimeTributarioRepository interface {
	// CRUD Básico
	Create(regime *models.EntidadeRegimeTributario) error
	Update(regime *models.EntidadeRegimeTributario) error
	Delete(entidadeID, item int) error
	FindByID(entidadeID, item int) (*models.EntidadeRegimeTributario, error)

	// Buscas Específicas
	FindByEntidadeID(entidadeID int) ([]models.EntidadeRegimeTributario, error)
	FindDefaultByEntidadeID(entidadeID int) (*models.EntidadeRegimeTributario, error)

	// Listagem com Filtros
	List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeRegimeTributario, int64, error)

	// Consultas de Validação (APENAS CONSULTAS)
	ExistsByEntidadeAndRegime(entidadeID int, regime int, excludeItem int) (bool, error)
	GetNextItemNumber(entidadeID int) (int, error)
	CountByEntidadeID(entidadeID int) (int64, error)
	HasRegimePadrao(entidadeID int, excludeItem int) (bool, error)
}

// ============================================================
// IMPLEMENTAÇÃO CONCRETA (privada)
// ============================================================

// entidadeRegimeTributarioRepository é a implementação concreta
type entidadeRegimeTributarioRepository struct {
	db *gorm.DB
}

// NewEntidadeRegimeTributarioRepository cria uma nova instância (retorna a interface)
func NewEntidadeRegimeTributarioRepository(db *gorm.DB) EntidadeRegimeTributarioRepository {
	return &entidadeRegimeTributarioRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD (APENAS PERSISTÊNCIA)
// ============================================================

// Create salva um novo regime com sequencial manual
func (r *entidadeRegimeTributarioRepository) Create(regime *models.EntidadeRegimeTributario) error {
	return r.db.Create(regime).Error
}

// Update atualiza um regime existente
func (r *entidadeRegimeTributarioRepository) Update(regime *models.EntidadeRegimeTributario) error {
	return r.db.
		Omit("Entidade", "created_at", "deleted_at").
		Model(regime).
		Where("ent_id = ? AND etr_item = ?", regime.EntidadeID, regime.ID).
		Updates(regime).Error
}

// Delete realiza a exclusão lógica
func (r *entidadeRegimeTributarioRepository) Delete(entidadeID, item int) error {
	return r.db.
		Model(&models.EntidadeRegimeTributario{}).
		Where("ent_id = ? AND etr_item = ?", entidadeID, item).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// ============================================================
// MÉTODOS DE BUSCA
// ============================================================

// FindByID busca um regime pelo ID composto (ent_id + etr_item)
func (r *entidadeRegimeTributarioRepository) FindByID(entidadeID, item int) (*models.EntidadeRegimeTributario, error) {
	var regime models.EntidadeRegimeTributario
	err := r.db.
		Where("ent_id = ? AND etr_item = ? AND deleted_at IS NULL", entidadeID, item).
		First(&regime).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("regime tributário não encontrado")
		}
		return nil, err
	}
	return &regime, nil
}

// FindByEntidadeID busca todos os regimes de uma entidade
func (r *entidadeRegimeTributarioRepository) FindByEntidadeID(entidadeID int) ([]models.EntidadeRegimeTributario, error) {
	var regimes []models.EntidadeRegimeTributario
	err := r.db.
		Where("ent_id = ? AND deleted_at IS NULL", entidadeID).
		Order("etr_item ASC").
		Find(&regimes).Error
	if err != nil {
		return nil, err
	}
	return regimes, nil
}

// FindDefaultByEntidadeID busca o regime padrão de uma entidade
func (r *entidadeRegimeTributarioRepository) FindDefaultByEntidadeID(entidadeID int) (*models.EntidadeRegimeTributario, error) {
	var regime models.EntidadeRegimeTributario
	err := r.db.
		Where("ent_id = ? AND etr_padrao = true AND deleted_at IS NULL", entidadeID).
		First(&regime).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Nenhum regime padrão encontrado
		}
		return nil, err
	}
	return &regime, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de regimes com paginação e filtros
func (r *entidadeRegimeTributarioRepository) List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeRegimeTributario, int64, error) {
	var regimes []models.EntidadeRegimeTributario
	var total int64

	query := r.db.Model(&models.EntidadeRegimeTributario{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.EntidadeRegimeTributario{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Order("etr_item ASC").Find(&regimes).Error
	if err != nil {
		return nil, 0, err
	}

	return regimes, total, nil
}

// ============================================================
// MÉTODOS DE CONSULTA PARA VALIDAÇÕES (APENAS CONSULTAS)
// ============================================================

// ExistsByEntidadeAndRegime verifica se já existe um regime com o mesmo regime para a entidade
func (r *entidadeRegimeTributarioRepository) ExistsByEntidadeAndRegime(entidadeID int, regime int, excludeItem int) (bool, error) {
	var count int64
	query := r.db.Model(&models.EntidadeRegimeTributario{}).
		Where("ent_id = ? AND etr_regime = ? AND deleted_at IS NULL", entidadeID, regime)

	if excludeItem > 0 {
		query = query.Where("etr_item != ?", excludeItem)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetNextItemNumber retorna o próximo número de item para uma entidade
func (r *entidadeRegimeTributarioRepository) GetNextItemNumber(entidadeID int) (int, error) {
	var maxItem int
	err := r.db.Model(&models.EntidadeRegimeTributario{}).
		Where("ent_id = ?", entidadeID).
		Select("COALESCE(MAX(etr_item), 0) + 1").
		Scan(&maxItem).Error
	if err != nil {
		return 0, err
	}

	return maxItem, nil
}

// CountByEntidadeID retorna a quantidade de regimes de uma entidade
func (r *entidadeRegimeTributarioRepository) CountByEntidadeID(entidadeID int) (int64, error) {
	var count int64
	err := r.db.Model(&models.EntidadeRegimeTributario{}).
		Where("ent_id = ? AND deleted_at IS NULL", entidadeID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// HasRegimePadrao verifica se a entidade já possui um regime padrão
func (r *entidadeRegimeTributarioRepository) HasRegimePadrao(entidadeID int, excludeItem int) (bool, error) {
	var count int64
	query := r.db.Model(&models.EntidadeRegimeTributario{}).
		Where("ent_id = ? AND etr_padrao = true AND deleted_at IS NULL", entidadeID)

	if excludeItem > 0 {
		query = query.Where("etr_item != ?", excludeItem)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
