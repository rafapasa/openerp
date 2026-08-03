package models

import (
	"time"
)

// ============================================================
// MODEL: Pais
// ============================================================

type Pais struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID      int    `gorm:"column:pai_id;primaryKey;autoIncrement" json:"id"`
	MoedaID *int   `gorm:"column:moeda_id" json:"moeda_id,omitempty"`
	Nome    string `gorm:"column:pai_nome;type:varchar(100);not null" json:"nome"`

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
	Moeda     *Moeda             `gorm:"foreignKey:MoedaID;references:moeda_id" json:"moeda,omitempty"`
	Enderecos []EntidadeEndereco `gorm:"foreignKey:PaisID;references:pai_id" json:"enderecos,omitempty"`
	//Documentos []DocumentoVenda   `gorm:"foreignKey:PaisID;references:pai_id" json:"documentos,omitempty"`
}

func (Pais) TableName() string {
	return "pais"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *Pais) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *Pais) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}

func (m *Pais) HasMoeda() bool {
	return m.MoedaID != nil && *m.MoedaID > 0
}
