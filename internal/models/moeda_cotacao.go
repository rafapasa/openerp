package models

import (
	"time"
)

// ============================================================
// MODEL: MoedaCotacao
// ============================================================

type MoedaCotacao struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID           int       `gorm:"column:moecot_iten;primaryKey;autoIncrement" json:"id"`
	MoedaID      int       `gorm:"column:moeda_id;not null" json:"moeda_id"`
	Data         time.Time `gorm:"column:moecot_data;type:date;not null" json:"data"`
	ValorCotacao float64   `gorm:"column:moecot_valorcotacao;type:decimal(15,2);not null" json:"valor_cotacao"`

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
	Moeda *Moeda `gorm:"foreignKey:MoedaID;references:moeda_id" json:"moeda,omitempty"`
}

func (MoedaCotacao) TableName() string {
	return "moeda_cotacao"
}

func (m *MoedaCotacao) BeforeCreate() error {
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

func (m *MoedaCotacao) BeforeUpdate() error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *MoedaCotacao) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *MoedaCotacao) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}
