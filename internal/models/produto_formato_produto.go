package models

import (
	"time"
)

// ============================================================
// MODEL: ProdutoFormatoProduto
// ============================================================

type ProdutoFormatoProduto struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	FormatoProdutoID int     `gorm:"column:fpro_id;primaryKey" json:"formato_produto_id"`
	ProdutoID        int     `gorm:"column:pro_id;primaryKey" json:"produto_id"`
	ValorAdicional   float64 `gorm:"column:pfp_valoradicional;type:decimal(15,4);not null" json:"valor_adicional"`

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
	FormatoProduto *FormatoProduto `gorm:"foreignKey:FormatoProdutoID;references:fpro_id" json:"formato_produto,omitempty"`
	Produto        *Produto        `gorm:"foreignKey:ProdutoID;references:pro_id" json:"produto,omitempty"`
}

func (ProdutoFormatoProduto) TableName() string {
	return "produto_formato_produto"
}

func (m *ProdutoFormatoProduto) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *ProdutoFormatoProduto) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}
