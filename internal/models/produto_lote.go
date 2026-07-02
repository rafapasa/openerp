package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: ProdutoLote
// ============================================================

type ProdutoLote struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ProdutoID      int        `gorm:"column:pro_id;primaryKey" json:"produto_id"`
	Item           int        `gorm:"column:prol_item;primaryKey" json:"item"`
	Lote           string     `gorm:"column:prol_lote;type:varchar(30);not null" json:"lote"`
	DataValidade   time.Time  `gorm:"column:prol_datavalidade;type:datetime;not null" json:"data_validade"`
	DataCompra     *time.Time `gorm:"column:prol_datacompra;type:datetime" json:"data_compra,omitempty"`
	DataFabricacao *time.Time `gorm:"column:prol_datafabricacao;type:datetime" json:"data_fabricacao,omitempty"`
	QuantCompra    float64    `gorm:"column:prol_quantcompra;type:decimal(15,4);not null" json:"quant_compra"`
	QuantSaldo     float64    `gorm:"column:prol_quantsaldo;type:decimal(15,4);not null" json:"quant_saldo"`

	// ============================================================
	// CAMPOS DE AUDITORIA
	// ============================================================
	CreatedAt time.Time  `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime;index" json:"deleted_at,omitempty"`
	CreatedBy *int       `gorm:"column:created_by" json:"created_by,omitempty"`
	UpdatedBy *int       `gorm:"column:updated_by" json:"updated_by,omitempty"`

	// ============================================================
	// RELACIONAMENTOS
	// ============================================================
	Produto *Produto `gorm:"foreignKey:ProdutoID;references:pro_id" json:"produto,omitempty"`
}

func (ProdutoLote) TableName() string {
	return "produto_lote"
}

func (m *ProdutoLote) BeforeCreate(tx *gorm.DB) error {
	// Buscar próximo item
	var maxItem int
	err := tx.Model(&ProdutoLote{}).
		Where("pro_id = ?", m.ProdutoID).
		Select("COALESCE(MAX(prol_item), 0) + 1").
		Scan(&maxItem).Error
	if err != nil {
		return err
	}
	m.Item = maxItem

	if m.CreatedBy == nil {
		m.CreatedBy = new(int)
		*m.CreatedBy = 0
	}
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

func (m *ProdutoLote) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

func (m *ProdutoLote) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *ProdutoLote) SoftDelete() {
	now := time.Time{}
	m.DeletedAt = &now
}

func (m *ProdutoLote) IsVencido() bool {
	return time.Now().After(m.DataValidade)
}
