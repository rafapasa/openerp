package models

import (
	"time"
)

// ============================================================
// MODEL: ProdutoCodigoBarras
// ============================================================

type ProdutoCodigoBarras struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ProdutoID           int      `gorm:"column:pro_id;primaryKey" json:"produto_id"`
	CodigoBarras        string   `gorm:"column:pcb_codigobarra;type:varchar(255);primaryKey" json:"codigo_barras"`
	PercentualDesconto  *float64 `gorm:"column:pcb_percdesconto;type:decimal(5,2)" json:"percentual_desconto,omitempty"`
	PercentualAcrescimo *float64 `gorm:"column:pcb_percacrescimo;type:decimal(5,2)" json:"percentual_acrescimo,omitempty"`
	QuantidadeEmbalagem *float64 `gorm:"column:pcb_quantembalagem;type:decimal(15,4)" json:"quantidade_embalagem,omitempty"`

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

func (ProdutoCodigoBarras) TableName() string {
	return "produto_codigobarra"
}

func (m *ProdutoCodigoBarras) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *ProdutoCodigoBarras) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}

func (m *ProdutoCodigoBarras) HasDesconto() bool {
	return m.PercentualDesconto != nil && *m.PercentualDesconto > 0
}

func (m *ProdutoCodigoBarras) HasAcrescimo() bool {
	return m.PercentualAcrescimo != nil && *m.PercentualAcrescimo > 0
}
