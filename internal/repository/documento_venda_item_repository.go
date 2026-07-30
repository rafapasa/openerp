// internal/repository/documento_venda_item_repository.go
package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/models"
)

// ============================================================
// INTERFACE - Define o contrato
// ============================================================

// DocumentoVendaItemRepository define o contrato para operações de banco
type DocumentoVendaItemRepository interface {
	// CRUD Básico
	Create(item *models.DocumentoVendaItem) error
	Update(ddvId, dviItem int, item *models.DocumentoVendaItem) error
	Delete(ddvId, dviItem int) error
	FindByID(ddvId, dviItem int) (*models.DocumentoVendaItem, error)
	GetByID(ddvId, dviItem int) (*models.DocumentoVendaItem, error)

	// Listagem com Filtros
	ListByDocumentoVendaID(limit, offset int, ddvId int) ([]models.DocumentoVendaItem, int64, error)

	// Consultas de Validação (APENAS CONSULTAS)
	GetNextItemNumber(ddvId int) (int, error)
	CountByDocumentoVendaID(ddvId int) (int64, error)
	ExistsByDocumentoVendaID(ddvId int) (bool, error)
	SumItensByDocumentoVendaID(ddvId int) (float64, error)
}

// ============================================================
// IMPLEMENTAÇÃO CONCRETA (privada)
// ============================================================

type documentoVendaItemRepository struct {
	db *gorm.DB
}

// NewDocumentoVendaItemRepository cria uma nova instância do repositório
func NewDocumentoVendaItemRepository(db *gorm.DB) DocumentoVendaItemRepository {
	return &documentoVendaItemRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD (APENAS PERSISTÊNCIA)
// ============================================================

// Create cria um novo item com item automático
func (r *documentoVendaItemRepository) Create(item *models.DocumentoVendaItem) error {
	return r.db.Create(item).Error
}

// Update atualiza um item existente
func (r *documentoVendaItemRepository) Update(ddvId, dviItem int, item *models.DocumentoVendaItem) error {
	return r.db.Model(&models.DocumentoVendaItem{}).
		Omit("produto", "operacao_fiscal", "cst_icms", "cst_ipi", "cst_pis_cofins", "created_at", "deleted_at").
		Where("ddv_id = ? AND dvi_item = ?", ddvId, dviItem).
		Updates(item).Error
}

// Delete realiza soft delete de um item
func (r *documentoVendaItemRepository) Delete(ddvId, dviItem int) error {
	return r.db.Model(&models.DocumentoVendaItem{}).
		Where("ddv_id = ? AND dvi_item = ?", ddvId, dviItem).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// FindByID busca um item por ID composto (ddv_id + dvi_item)
func (r *documentoVendaItemRepository) FindByID(ddvId, dviItem int) (*models.DocumentoVendaItem, error) {
	var item models.DocumentoVendaItem
	err := r.db.
		Preload("produto").
		Preload("operacaofiscal").
		Preload("cst_icms").
		Preload("cst_ipi").
		Preload("cst_pis_cofins").
		Where("ddv_id = ? AND dvi_item = ? AND deleted_at IS NULL", ddvId, dviItem).
		First(&item).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item não encontrado")
		}
		return nil, err
	}
	return &item, nil
}

// GetByID busca um item por ID composto (sem relacionamentos)
func (r *documentoVendaItemRepository) GetByID(ddvId, dviItem int) (*models.DocumentoVendaItem, error) {
	var item models.DocumentoVendaItem
	err := r.db.
		Where("ddv_id = ? AND dvi_item = ? AND deleted_at IS NULL", ddvId, dviItem).
		First(&item).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item não encontrado")
		}
		return nil, err
	}
	return &item, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// ListByDocumentoVendaID lista todos os itens de um documento de venda com paginação
func (r *documentoVendaItemRepository) ListByDocumentoVendaID(limit, offset int, ddvId int) ([]models.DocumentoVendaItem, int64, error) {
	var items []models.DocumentoVendaItem
	var total int64

	query := r.db.Model(&models.DocumentoVendaItem{}).
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
		Limit(limit).
		Offset(offset).
		Order("dvi_item ASC").
		Find(&items).Error

	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// ============================================================
// MÉTODOS DE CONSULTA PARA VALIDAÇÕES (APENAS CONSULTAS)
// ============================================================

// GetNextItemNumber retorna o próximo número de item para um documento
func (r *documentoVendaItemRepository) GetNextItemNumber(ddvId int) (int, error) {
	var maxItem int
	err := r.db.Model(&models.DocumentoVendaItem{}).
		Where("ddv_id = ?", ddvId).
		Select("COALESCE(MAX(dvi_item), 0) + 1").
		Scan(&maxItem).Error
	if err != nil {
		return 0, err
	}
	return maxItem, nil
}

// CountByDocumentoVendaID retorna a quantidade de itens de um documento
func (r *documentoVendaItemRepository) CountByDocumentoVendaID(ddvId int) (int64, error) {
	var count int64
	err := r.db.Model(&models.DocumentoVendaItem{}).
		Where("ddv_id = ? AND deleted_at IS NULL", ddvId).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// ExistsByDocumentoVendaID verifica se um documento possui itens
func (r *documentoVendaItemRepository) ExistsByDocumentoVendaID(ddvId int) (bool, error) {
	var count int64
	err := r.db.Model(&models.DocumentoVendaItem{}).
		Where("ddv_id = ? AND deleted_at IS NULL", ddvId).
		Limit(1).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SumItensByDocumentoVendaID retorna o valor total dos itens de um documento
func (r *documentoVendaItemRepository) SumItensByDocumentoVendaID(ddvId int) (float64, error) {
	var total float64
	err := r.db.Model(&models.DocumentoVendaItem{}).
		Where("ddv_id = ? AND deleted_at IS NULL", ddvId).
		Select("COALESCE(SUM(dvi_valortotal), 0)").
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}