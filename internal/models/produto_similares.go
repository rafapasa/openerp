package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: ProdutoSimilares
// ============================================================

type ProdutoSimilares struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ProdutoID        int `gorm:"column:pro_id;primaryKey" json:"produto_id"`
	ProdutoSimilarID int `gorm:"column:sim_pro_id;primaryKey" json:"produto_similar_id"`

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
	Produto        *Produto `gorm:"foreignKey:ProdutoID;references:pro_id" json:"produto,omitempty"`
	ProdutoSimilar *Produto `gorm:"foreignKey:ProdutoSimilarID;references:pro_id" json:"produto_similar,omitempty"`
}

func (ProdutoSimilares) TableName() string {
	return "produto_similares"
}

func (m *ProdutoSimilares) BeforeCreate(tx *gorm.DB) error {
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

func (m *ProdutoSimilares) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

func (m *ProdutoSimilares) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *ProdutoSimilares) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}
