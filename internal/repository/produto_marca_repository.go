package repository

import (
	"errors"
	"fmt"

	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
	"gorm.io/gorm"
)

type ProdutoMarcaRepository struct {
	db *gorm.DB
}

func NewProdutoMarcaRepository(db *gorm.DB) *ProdutoMarcaRepository {
	return &ProdutoMarcaRepository{db: db}
}

func (r *ProdutoMarcaRepository) Create(marca *models.ProdutoMarca) error {
	return r.db.Create(marca).Error
}

func (r *ProdutoMarcaRepository) Update(id int, marca *models.ProdutoMarca) error {
	return r.db.Model(&models.ProdutoMarca{}).Where("promar_id = ?", id).Updates(marca).Error
}

func (r *ProdutoMarcaRepository) Delete(id int) error {
	marca, err := r.FindByID(id)
	if err != nil {
		return err
	}
	if marca.IsDeleted() {
		return fmt.Errorf("marca de produto '%s' já está deletada", marca.Descricao)
	}
	marca.SoftDelete()
	return r.Update(id, marca)
}

func (r *ProdutoMarcaRepository) FindByID(id int) (*models.ProdutoMarca, error) {
	var marca models.ProdutoMarca
	err := r.db.Where("promar_id = ? AND deleted_at IS NULL", id).First(&marca).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("marca de produto não encontrada")
		}
		return nil, err
	}
	return &marca, nil
}

func (r *ProdutoMarcaRepository) List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoMarca, int64, error) {
	var marcas []models.ProdutoMarca
	var total int64

	query := r.db.Model(&models.ProdutoMarca{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.ProdutoMarca{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Order("promar_id DESC").Find(&marcas).Error
	if err != nil {
		return nil, 0, err
	}

	return marcas, total, nil
}

func (r *ProdutoMarcaRepository) ExistsByDescricao(descricao string, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.ProdutoMarca{}).Where("promar_descricao = ? AND deleted_at IS NULL", descricao)
	if excludeID > 0 {
		query = query.Where("promar_id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ProdutoMarcaRepository) CountByMarca(id int) (int64, error) {
	var count int64
	err := r.db.Model(&models.Produto{}).Where("promar_id = ? AND deleted_at IS NULL", id).Count(&count).Error
	return count, err
}