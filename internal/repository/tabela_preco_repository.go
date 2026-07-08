package repository

import (
	"errors"
	"fmt"

	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
	"gorm.io/gorm"
)

type TabelaPrecoRepository struct {
	db *gorm.DB
}

func NewTabelaPrecoRepository(db *gorm.DB) *TabelaPrecoRepository {
	return &TabelaPrecoRepository{db: db}
}

func (r *TabelaPrecoRepository) Create(tabela *models.TabelaPreco) error {
	return r.db.Create(tabela).Error
}

func (r *TabelaPrecoRepository) Update(id int, tabela *models.TabelaPreco) error {
	return r.db.Model(&models.TabelaPreco{}).Where("tbp_id = ?", id).Updates(tabela).Error
}

func (r *TabelaPrecoRepository) Delete(id int) error {
	tabela, err := r.FindByID(id)
	if err != nil {
		return err
	}
	if tabela.IsDeleted() {
		return fmt.Errorf("tabela de preço '%s' já está deletada", tabela.Descricao)
	}
	tabela.SoftDelete()
	return r.Update(id, tabela)
}

func (r *TabelaPrecoRepository) FindByID(id int) (*models.TabelaPreco, error) {
	var tabela models.TabelaPreco
	err := r.db.Where("tbp_id = ? AND deleted_at IS NULL", id).First(&tabela).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tabela de preço não encontrada")
		}
		return nil, err
	}
	return &tabela, nil
}

func (r *TabelaPrecoRepository) List(limit, offset int, filters map[string]interface{}) ([]models.TabelaPreco, int64, error) {
	var tabelas []models.TabelaPreco
	var total int64

	query := r.db.Model(&models.TabelaPreco{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.TabelaPreco{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Order("tbp_id DESC").Find(&tabelas).Error
	if err != nil {
		return nil, 0, err
	}

	return tabelas, total, nil
}

func (r *TabelaPrecoRepository) ExistsByDescricao(descricao string, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.TabelaPreco{}).Where("tbp_descricao = ? AND deleted_at IS NULL", descricao)
	if excludeID > 0 {
		query = query.Where("tbp_id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *TabelaPrecoRepository) CountByTabelaPreco(id int) (int64, error) {
	var count int64
	err := r.db.Model(&models.DocumentoVenda{}).Where("tbp_id = ? AND deleted_at IS NULL", id).Count(&count).Error
	return count, err
}