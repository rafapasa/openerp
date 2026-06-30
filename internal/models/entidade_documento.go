package models

import (
	"time"
)

// ============================================================
// MODEL: EntidadeDocumento
// ============================================================

type EntidadeDocumento struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	EntidadeID   int       `gorm:"column:ent_id;primaryKey" json:"entidade_id"`
	Item         int       `gorm:"column:edoc_item;primaryKey" json:"item"`
	Descricao    *string   `gorm:"column:edoc_descricao;type:varchar(255)" json:"descricao,omitempty"`
	DataInclusao time.Time `gorm:"column:edoc_datainclusao;type:date;not null" json:"data_inclusao"`
	Arquivo      []byte    `gorm:"column:edoc_arquivo;type:longblob;not null" json:"-"` // Ocultar no JSON
	Tipo         *string   `gorm:"column:edoc_tipo;type:varchar(10)" json:"tipo,omitempty"`

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
	Entidade *Entidade `gorm:"foreignKey:EntidadeID;references:ent_id" json:"entidade,omitempty"`
}

func (EntidadeDocumento) TableName() string {
	return "entidade_documento"
}

func (m *EntidadeDocumento) BeforeCreate() error {
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

func (m *EntidadeDocumento) BeforeUpdate() error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *EntidadeDocumento) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *EntidadeDocumento) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}
