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

// LimiteCreditoRepository é o repositório para LimiteCredito.
type EntidadeLimiteCreditoRepository struct {
	db *gorm.DB
}

// ============================================================
// CONSTRUCTOR
// ============================================================

// NewEntidadeLimiteCreditoRepository cria uma nova instância.
func NewEntidadeLimiteCreditoRepository(db *gorm.DB) *EntidadeLimiteCreditoRepository {
	return &EntidadeLimiteCreditoRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD
// ============================================================

// Create salva um novo limite de crédito.
func (r *EntidadeLimiteCreditoRepository) Create(limite *models.EntidadeLimiteCredito) error {
	return r.db.Create(limite).Error
}

// Update atualiza um limite de crédito existente.
func (r *EntidadeLimiteCreditoRepository) Update(id int, limite *models.EntidadeLimiteCredito) error {
	// Omit() é usado para evitar a atualização de campos de relacionamentos
	// que não devem ser alterados nesta operação.
	return r.db.Model(limite).Where("elc_id = ?", id).Updates(limite).Error
}

// Delete realiza a exclusão lógica.
func (r *EntidadeLimiteCreditoRepository) Delete(id int) error {
	limite, err := r.FindByID(id)
	if err != nil {
		return err
	}
	if limite.IsDeleted() {
		return errors.New("limite de crédito já foi deletado")
	}
	limite.SoftDelete()
	return r.Update(id, limite)
}

// ============================================================
// MÉTODOS DE BUSCA
// ============================================================

// FindByID busca um limite de crédito pelo ID.
func (r *EntidadeLimiteCreditoRepository) FindByID(id int) (*models.EntidadeLimiteCredito, error) {
	var limite models.EntidadeLimiteCredito
	err := r.db.Where("elc_id = ? AND deleted_at IS NULL", id).First(&limite).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("limite de crédito não encontrado")
		}
		return nil, err
	}
	return &limite, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de limites com paginação e filtros.
func (r *EntidadeLimiteCreditoRepository) List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeLimiteCredito, int64, error) {
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
// MÉTODO ADICIONAL: Verificar duplicidade
// ============================================================

// ExistsByDescricao verifica se já existe um limite com a mesma descrição.
func (r *EntidadeLimiteCreditoRepository) ExistsByDescricao(descricao string, excludeID int) (bool, error) {
	var count int64

	query := r.db.Model(&models.EntidadeLimiteCredito{}).Where("elc_descricao = ? AND deleted_at IS NULL", descricao)

	if excludeID > 0 {
		query = query.Where("elc_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
