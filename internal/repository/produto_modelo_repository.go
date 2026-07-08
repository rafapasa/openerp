package repository

import (
	"errors"
	"fmt"

	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
	"gorm.io/gorm"
)

type ProdutoModeloRepository struct {
	db *gorm.DB
}

func NewProdutoModeloRepository(db *gorm.DB) *ProdutoModeloRepository {
	return &ProdutoModeloRepository{db: db}
}

func (r *ProdutoModeloRepository) Create(modelo *models.ProdutoModelo) error {
	return r.db.Create(modelo).Error
}

func (r *ProdutoModeloRepository) Update(id int, modelo *models.ProdutoModelo) error {
	return r.db.Model(&models.ProdutoModelo{}).Where("prom_id = ?", id).Updates(modelo).Error
}

func (r *ProdutoModeloRepository) Delete(id int) error {
	modelo, err := r.FindByID(id)
	if err != nil {
		return err
	}
	if modelo.IsDeleted() {
		return fmt.Errorf("modelo de produto '%s' já está deletado", modelo.Descricao)
	}
	modelo.SoftDelete()
	return r.Update(id, modelo)
}

func (r *ProdutoModeloRepository) FindByID(id int) (*models.ProdutoModelo, error) {
	var modelo models.ProdutoModelo
	err := r.db.Where("prom_id = ? AND deleted_at IS NULL", id).First(&modelo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("modelo de produto não encontrado")
		}
		return nil, err
	}
	return &modelo, nil
}

func (r *ProdutoModeloRepository) List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoModelo, int64, error) {
	var modelos []models.ProdutoModelo
	var total int64

	query := r.db.Model(&models.ProdutoModelo{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.ProdutoModelo{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Order("prom_id DESC").Find(&modelos).Error
	if err != nil {
		return nil, 0, err
	}

	return modelos, total, nil
}

func (r *ProdutoModeloRepository) ExistsByDescricao(descricao string, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.ProdutoModelo{}).Where("prom_descricao = ? AND deleted_at IS NULL", descricao)
	if excludeID > 0 {
		query = query.Where("prom_id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ProdutoModeloRepository) CountByModelo(id int) (int64, error) {
	var count int64
	err := r.db.Model(&models.Produto{}).Where("prom_id = ? AND deleted_at IS NULL", id).Count(&count).Error
	return count, err
}