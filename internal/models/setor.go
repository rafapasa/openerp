package models

import (
	"time"
)

// ============================================================
// MODEL: Setor
// ============================================================

type Setor struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID        int    `gorm:"column:set_id;primaryKey;autoIncrement" json:"id"`
	Descricao string `gorm:"column:set_descricao;type:varchar(255);not null" json:"descricao"`

	// ============================================================
	// CAMPOS DE AUDITORIA
	// ============================================================
	CreatedAt time.Time  `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime;index" json:"deleted_at,omitempty"`
	CreatedBy *int       `gorm:"column:created_by" json:"created_by,omitempty"`
	UpdatedBy *int       `gorm:"column:updated_by" json:"updated_by,omitempty"`
}

func (Setor) TableName() string {
	return "setor"
}

func (m *Setor) BeforeCreate() error {
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

func (m *Setor) BeforeUpdate() error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *Setor) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *Setor) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}
