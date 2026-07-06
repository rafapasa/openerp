package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: Horario
// ============================================================

type Horario struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID        int    `gorm:"column:hor_id;primaryKey;autoIncrement" json:"id"`
	Descricao string `gorm:"column:hor_descricao;type:varchar(255);not null" json:"descricao"`

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
	Pontos    []HorarioPonto `gorm:"foreignKey:HorarioID;references:ID" json:"pontos,omitempty"`
	Entidades []Entidade     `gorm:"foreignKey:HorarioID;references:ID" json:"entidades,omitempty"`
}

func (Horario) TableName() string {
	return "horario"
}

func (m *Horario) BeforeCreate(tx *gorm.DB) error {
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

func (m *Horario) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *Horario) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *Horario) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}
