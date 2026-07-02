package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: Manutencao
// ============================================================

type Manutencao struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID        int    `gorm:"column:man_id;primaryKey;autoIncrement" json:"id"`
	Descricao string `gorm:"column:man_descricao;type:varchar(255);not null" json:"descricao"`

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
	VeiculosManutencoes []VeiculoManutencao `gorm:"foreignKey:ManutencaoID;references:man_id" json:"veiculos_manutencoes,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (Manutencao) TableName() string {
	return "manutencao"
}

func (m *Manutencao) BeforeCreate(tx *gorm.DB) error {
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

func (m *Manutencao) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se a manutenção foi deletada logicamente
func (m *Manutencao) IsDeleted() bool {
	return m.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (m *Manutencao) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}

// HasVeiculosManutencoes verifica se a manutenção possui veículos associados
func (m *Manutencao) HasVeiculosManutencoes() bool {
	return len(m.VeiculosManutencoes) > 0
}
