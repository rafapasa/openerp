package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: ProdutoMedia
// ============================================================

type ProdutoMedia struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID              int       `gorm:"column:promed_id;primaryKey;autoIncrement" json:"id"`
	ProdutoID       int       `gorm:"column:pro_id;not null" json:"produto_id"`
	EmpresaFilialID int       `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	Data            time.Time `gorm:"column:promed_data;type:datetime;not null" json:"data"`
	QuantConsumo    float64   `gorm:"column:promed_quantconsumo;type:decimal(15,4);not null" json:"quant_consumo"`
	QuantCompra     float64   `gorm:"column:promed_quantcompra;type:decimal(15,4);not null" json:"quant_compra"`

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
	Produto       *Produto       `gorm:"foreignKey:ProdutoID;references:pro_id" json:"produto,omitempty"`
	EmpresaFilial *EmpresaFilial `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
}

func (ProdutoMedia) TableName() string {
	return "produto_media"
}

func (m *ProdutoMedia) BeforeCreate(tx *gorm.DB) error {
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

func (m *ProdutoMedia) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

func (m *ProdutoMedia) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *ProdutoMedia) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}
