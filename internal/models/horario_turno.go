package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: HorarioTurno
// ============================================================

type HorarioTurno struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID        int    `gorm:"column:hort_id;primaryKey;autoIncrement" json:"id"`
	Descricao string `gorm:"column:hort_descricao;type:varchar(255);not null" json:"descricao"`

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
	Pontos []HorarioPonto `gorm:"foreignKey:TurnoID;references:ID" json:"pontos,omitempty"`
}

func (HorarioTurno) TableName() string {
	return "horario_turno"
}

func (m *HorarioTurno) BeforeCreate(tx *gorm.DB) error {
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

func (m *HorarioTurno) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *HorarioTurno) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *HorarioTurno) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}
