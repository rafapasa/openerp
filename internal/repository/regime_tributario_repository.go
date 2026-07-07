package repository

import (
	"errors"

	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
	"gorm.io/gorm"
)

// ============================================================
// TYPES
// ============================================================

// RegimeTributarioRepository é o repositório para RegimeTributario.
type RegimeTributarioRepository struct {
	db *gorm.DB
}

// ============================================================
// CONSTRUCTOR
// ============================================================

// NewRegimeTributarioRepository cria uma nova instância.
func NewRegimeTributarioRepository(db *gorm.DB) *RegimeTributarioRepository {
	return &RegimeTributarioRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD
// ============================================================

// Create salva um novo regime tributário.
func (r *RegimeTributarioRepository) Create(regime *models.EntidadeRegimeTributario) error {
	return r.db.Create(regime).Error
}

// Update atualiza um regime tributário existente.
func (r *RegimeTributarioRepository) Update(id int, regime *models.EntidadeRegimeTributario) error {
	return r.db.
		Model(regime).
		Where("ret_id = ?", id).
		Updates(regime).Error
}

// Delete realiza a exclusão lógica.
func (r *RegimeTributarioRepository) Delete(id int) error {
	regime, err := r.FindByID(id)
	if err != nil {
		return err
	}
	if regime.IsDeleted() {
		return errors.New("regime tributário já foi deletado")
	}
	regime.SoftDelete()
	return r.Update(id, regime)
}

// ============================================================
// MÉTODOS DE BUSCA
// ============================================================

// FindByID busca um regime tributário pelo ID.
func (r *RegimeTributarioRepository) FindByID(id int) (*models.EntidadeRegimeTributario, error) {
	var regime models.EntidadeRegimeTributario
	err := r.db.Where("ret_id = ? AND deleted_at IS NULL", id).First(&regime).Error
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
func (r *RegimeTributarioRepository) List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeRegimeTributario, int64, error) {
	var regimes []models.EntidadeRegimeTributario
	var total int64

	query := r.db.Model(&models.EntidadeRegimeTributario{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.EntidadeRegimeTributario{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Order("etr_descricao ASC").Find(&regimes).Error
	if err != nil {
		return nil, 0, err
	}

	return regimes, total, nil
}


