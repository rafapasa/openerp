package repository

import (
	"errors"
	"fmt"

	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
	"gorm.io/gorm"
)

// ============================================================
// TYPES
// ============================================================

// ProdutoSubgrupoRepository é o repositório para ProdutoSubgrupo.
type ProdutoSubgrupoRepository struct {
	db *gorm.DB
}

// ============================================================
// CONSTRUCTOR
// ============================================================

// NewProdutoSubgrupoRepository cria uma nova instância.
func NewProdutoSubgrupoRepository(db *gorm.DB) *ProdutoSubgrupoRepository {
	return &ProdutoSubgrupoRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD
// ============================================================

// Create salva um novo subgrupo de produto.
func (r *ProdutoSubgrupoRepository) Create(subgrupo *models.ProdutoSubgrupo) error {
	return r.db.Create(subgrupo).Error
}

// Update atualiza um subgrupo de produto existente.
func (r *ProdutoSubgrupoRepository) Update(id int, subgrupo *models.ProdutoSubgrupo) error {
	return r.db.Model(&models.ProdutoSubgrupo{}).Where("prosg_id = ?", id).Updates(subgrupo).Error
}

// Delete realiza a exclusão lógica.
func (r *ProdutoSubgrupoRepository) Delete(id int) error {
	subgrupo, err := r.FindByID(id)
	if err != nil {
		return err
	}
	if subgrupo.IsDeleted() {
		return fmt.Errorf("subgrupo de produto '%s' já está deletado", subgrupo.Descricao)
	}
	subgrupo.SoftDelete()
	return r.Update(id, subgrupo)
}

// ============================================================
// MÉTODOS DE BUSCA
// ============================================================

// FindByID busca um subgrupo de produto pelo ID.
func (r *ProdutoSubgrupoRepository) FindByID(id int) (*models.ProdutoSubgrupo, error) {
	var subgrupo models.ProdutoSubgrupo
	err := r.db.Where("prosg_id = ? AND deleted_at IS NULL", id).First(&subgrupo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subgrupo de produto não encontrado")
		}
		return nil, err
	}
	return &subgrupo, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de subgrupos com paginação e filtros.
func (r *ProdutoSubgrupoRepository) List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoSubgrupo, int64, error) {
	var subgrupos []models.ProdutoSubgrupo
	var total int64

	query := r.db.Model(&models.ProdutoSubgrupo{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.ProdutoSubgrupo{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Order("prosg_id DESC").Find(&subgrupos).Error
	if err != nil {
		return nil, 0, err
	}

	return subgrupos, total, nil
}

// ============================================================
// MÉTODO ADICIONAL: Verificar duplicidade
// ============================================================

// ExistsByDescricao verifica se já existe um subgrupo com a mesma descrição.
func (r *ProdutoSubgrupoRepository) ExistsByDescricao(descricao string, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.ProdutoSubgrupo{}).Where("prosg_descricao = ? AND deleted_at IS NULL", descricao)
	if excludeID > 0 {
		query = query.Where("prosg_id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
