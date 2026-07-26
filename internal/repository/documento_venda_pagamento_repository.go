package repository

import (
	"github.com/openerp/backend/internal/models"
	"gorm.io/gorm"
)

type DocumentoVendaPagamentoRepository struct {
	db *gorm.DB
}

func NewDocumentoVendaPagamentoRepository(db *gorm.DB) *DocumentoVendaPagamentoRepository {
	return &DocumentoVendaPagamentoRepository{db: db}
}

func (r *DocumentoVendaPagamentoRepository) Create(pagamento *models.DocumentoVendaPagamento) error {
	var maxItem int
	err := r.db.Model(&models.DocumentoVendaPagamento{}).
		Where("ddv_id = ?", pagamento.DocumentoVendaID).
		Select("COALESCE(MAX(dvp_item), 0) + 1").
		Scan(&maxItem).Error
	if err != nil {
		return err
	}

	pagamento.Item = maxItem
	return r.db.Create(pagamento).Error
}

func (r *DocumentoVendaPagamentoRepository) Update(ddvId, dvpItem int, pagamento *models.DocumentoVendaPagamento) error {
	return r.db.Model(&models.DocumentoVendaPagamento{}).
		Where("ddv_id = ? AND dvp_item = ?", ddvId, dvpItem).
		Updates(pagamento).Error
}

func (r *DocumentoVendaPagamentoRepository) Delete(ddvId, dvpItem int) error {
	pagamento, err := r.FindByID(ddvId, dvpItem)
	if err != nil {
		return err
	}

	if pagamento.IsDeleted() {
		return nil
	}

	pagamento.SoftDelete()
	return r.Update(ddvId, dvpItem, pagamento)
}

func (r *DocumentoVendaPagamentoRepository) FindByID(ddvId, dvpItem int) (*models.DocumentoVendaPagamento, error) {
	var pagamento models.DocumentoVendaPagamento
	err := r.db.Model(&pagamento).
		Where("ddv_id = ? AND dvp_item = ? AND deleted_at IS NULL", ddvId, dvpItem).
		First(&pagamento).Error
	return &pagamento, err
}

func (r *DocumentoVendaPagamentoRepository) ListByDocumentoVendaID(limit, offset int, ddvId int) ([]models.DocumentoVendaPagamento, int64, error) {
	var pagamentos []models.DocumentoVendaPagamento
	var total int64

	query := r.db.
		Model(&models.DocumentoVendaPagamento{}).
		Where("ddv_id = ? AND deleted_at IS NULL", ddvId)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Where("ddv_id = ? AND deleted_at IS NULL", ddvId).
		Limit(limit).
		Offset(offset).
		Order("dvp_item ASC").
		Find(&pagamentos).
		Error
	if err != nil {
		return nil, 0, err
	}
	return pagamentos, total, nil
}