package models

import (
	"time"
)

// ============================================================
// MODEL: Moeda
// ============================================================

type Moeda struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID        int    `gorm:"column:moeda_id;primaryKey;autoIncrement" json:"id"`
	Descricao string `gorm:"column:moeda_descricao;type:varchar(255);not null" json:"descricao"`
	Sigla     string `gorm:"column:moeda_sifra;type:varchar(10);not null" json:"sigla"`

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
	Paises   []Pais         `gorm:"foreignKey:MoedaID;references:ID" json:"paises,omitempty"`
	Cotacoes []MoedaCotacao `gorm:"foreignKey:MoedaID;references:ID" json:"cotacoes,omitempty"`
}

func (Moeda) TableName() string {
	return "moeda"
}

func (m *Moeda) BeforeCreate() error {
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

func (m *Moeda) BeforeUpdate() error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *Moeda) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *Moeda) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}
