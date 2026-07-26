package repository

import (
	"github.com/openerp/backend/internal/models"
	"gorm.io/gorm"
)

type DocumentoVendaItemRepository struct {
	db *gorm.DB
}

func NewDocumentoVendaItemRepository(db *gorm.DB) *DocumentoVendaItemRepository {
	return &DocumentoVendaItemRepository{db: db}
}

func (r *DocumentoVendaItemRepository) Create(item *models.DocumentoVendaItem) error {
	var maxItem int
	err := r.db.Model(&models.DocumentoVendaItem{}).
		Where("ddv_id = ?", item.DocumentoVendaID).
		Select("COALESCE(MAX(dvi_item), 0) + 1").
		Scan(&maxItem).Error
	if err != nil {
		return err
	}

	item.Item = maxItem
	return r.db.Create(item).Error
}

func (r *DocumentoVendaItemRepository) Update(ddvId, dviItem int, item *models.DocumentoVendaItem) error {
	return r.db.Model(&models.DocumentoVendaItem{}).
		Omit("produto", "operacao_fiscal", "cst_icms", "cst_ipi", "cst_pis_cofins").
		Where("ddv_id = ? AND dvi_item = ?", ddvId, dviItem).
		Updates(item).Error
}

func (r *DocumentoVendaItemRepository) Delete(ddvId, dviItem int) error {
	item, err := r.FindByID(ddvId, dviItem)
	if err != nil {
		return err
	}

	if item.IsDeleted() {
		return nil
	}

	item.SoftDelete()
	return r.Update(ddvId, dviItem, item)
}

func (r *DocumentoVendaItemRepository) FindByID(ddvId, dviItem int) (*models.DocumentoVendaItem, error) {
	var item models.DocumentoVendaItem
	err := r.db.Model(&item).
		Preload("produto").
		Preload("operacaofiscal").
		Preload("cst_icms").
		Preload("cst_ipi").
		Preload("cst_pis_cofins").
		Where("ddv_id = ? AND dvi_item = ? AND deleted_at IS NULL", ddvId, dviItem).
		First(&item).Error
	return &item, err
}

func (s *DocumentoVendaItemRepository) GetByID(ddvId, dviItem int) (*models.DocumentoVendaItem, error) {
	var item models.DocumentoVendaItem
	err := s.db.
		Model(&item).
		Where("ddv_id = ? and dvi_item = ? AND deleted_at IS NULL", ddvId, dviItem).
		First(&item).
		Error
	return &item, err

}

func (s *DocumentoVendaItemRepository) ListByDocumentoVendaID(limit, offset int, ddvId int) ([]models.DocumentoVendaItem, int64, error) {
	var items []models.DocumentoVendaItem
	var total int64

	query := s.db.
		Model(&models.DocumentoVendaItem{}).
		Where("ddv_id = ? AND deleted_at IS NULL", ddvId)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("produto").
		Preload("operacaofiscal").
		Preload("cst_icms").
		Preload("cst_ipi").
		Preload("cst_pis_cofins").
		Where("ddv_id = ? AND deleted_at IS NULL", ddvId).
		Limit(limit).
		Offset(offset).
		Order("dvi_item ASC").
		Find(&items).
		Error
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
