package models

import (
	"time"
)

// ============================================================
// MODEL: ProdutoSubgrupoDespesa
// ============================================================

type ProdutoSubgrupoDespesa struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ProdutoSubgrupoID int `gorm:"column:prosg_id;primaryKey" json:"produto_subgrupo_id"`
	DespesaID         int `gorm:"column:desp_id;primaryKey" json:"despesa_id"`

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
	ProdutoSubgrupo *ProdutoSubgrupo `gorm:"foreignKey:ProdutoSubgrupoID;references:prosg_id" json:"produto_subgrupo,omitempty"`
	Despesa         *Despesa         `gorm:"foreignKey:DespesaID;references:desp_id" json:"despesa,omitempty"`
}

func (ProdutoSubgrupoDespesa) TableName() string {
	return "produto_subgrupo_despesa"
}

func (m *ProdutoSubgrupoDespesa) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *ProdutoSubgrupoDespesa) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}
