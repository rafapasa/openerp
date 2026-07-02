package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: ProdutoFichaOperacional
// ============================================================

type ProdutoFichaOperacional struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ProdutoID  int      `gorm:"column:pro_id;primaryKey" json:"produto_id"`
	DespesaID  int      `gorm:"column:desp_id;primaryKey" json:"despesa_id"`
	Quantidade float64  `gorm:"column:pfo_quantidade;type:decimal(15,4);not null" json:"quantidade"`
	Custo      *float64 `gorm:"column:pfo_custo;type:decimal(15,4)" json:"custo,omitempty"`

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
	Despesa *Despesa `gorm:"foreignKey:DespesaID;references:desp_id" json:"despesa,omitempty"`
}

func (ProdutoFichaOperacional) TableName() string {
	return "produto_fichaoperacional"
}

func (m *ProdutoFichaOperacional) BeforeCreate(tx *gorm.DB) error {
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

func (m *ProdutoFichaOperacional) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

func (m *ProdutoFichaOperacional) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *ProdutoFichaOperacional) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}
