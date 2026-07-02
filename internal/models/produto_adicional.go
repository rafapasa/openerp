package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: ProdutoAdicional
// ============================================================

type ProdutoAdicional struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ProdutoID int     `gorm:"column:pro_id;primaryKey" json:"produto_id"`
	Item      int     `gorm:"column:proad_item;primaryKey" json:"item"`
	Descricao string  `gorm:"column:proad_descricao;type:varchar(255);not null" json:"descricao"`
	Valor     float64 `gorm:"column:proad_valor;type:decimal(15,2);not null;default:0.00" json:"valor"`
	Padrao    *int    `gorm:"column:proad_padrao;default:0" json:"padrao,omitempty"`
	Nivel     *int    `gorm:"column:proad_nivel" json:"nivel,omitempty"`

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

func (ProdutoAdicional) TableName() string {
	return "produto_adicional"
}

func (m *ProdutoAdicional) BeforeCreate(tx *gorm.DB) error {
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

func (m *ProdutoAdicional) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

func (m *ProdutoAdicional) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *ProdutoAdicional) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}
