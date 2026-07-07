package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// TYPES
// ============================================================

// EntidadeRegimeTributarioRepository é o repositório para regimes da entidade.
type EntidadeRegimeTributarioRepository struct {
	db *gorm.DB
}

// ============================================================
// CONSTRUCTOR
// ============================================================

// NewEntidadeRegimeTributarioRepository cria uma nova instância.
func NewEntidadeRegimeTributarioRepository(db *gorm.DB) *EntidadeRegimeTributarioRepository {
	return &EntidadeRegimeTributarioRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD
// ============================================================

// Create salva um novo regime com sequencial manual.
func (r *EntidadeRegimeTributarioRepository) Create(regime *models.EntidadeRegimeTributario) error {
	// 1. Buscar o próximo número para esta entidade.
	var maxItem int
	err := r.db.Model(&models.EntidadeRegimeTributario{}).
		Where("ent_id = ?", regime.EntidadeID).
		Select("COALESCE(MAX(etr_item), 0) + 1").
		Scan(&maxItem).Error
	if err != nil {
		return err
	}

	// 2. Atribuir o próximo número.
	regime.ID = maxItem

	// 3. Salvar.
	return r.db.Create(regime).Error
}

// Update atualiza um regime existente.
func (r *EntidadeRegimeTributarioRepository) Update(regime *models.EntidadeRegimeTributario) error {
	return r.db.
		Omit("Entidade").
		Model(regime).
		Where("ent_id = ? AND etr_item = ?", regime.EntidadeID, regime.ID).
		Updates(regime).Error
}

// Delete realiza a exclusão lógica.
func (r *EntidadeRegimeTributarioRepository) Delete(entidadeID, item int) error {
	regime, err := r.FindByID(entidadeID, item)
	if err != nil {
		return err
	}

	if regime.IsDeleted() {
		return errors.New("regime tributário já foi deletado")
	}

	regime.SoftDelete()
	return r.Update(regime)
}

// ============================================================
// MÉTODOS DE BUSCA
// ============================================================

// FindByID busca um regime pelo ID composto (ent_id + etr_item).
func (r *EntidadeRegimeTributarioRepository) FindByID(entidadeID, item int) (*models.EntidadeRegimeTributario, error) {
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

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de regimes com paginação e filtros.
func (r *EntidadeRegimeTributarioRepository) List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeRegimeTributario, int64, error) {
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
