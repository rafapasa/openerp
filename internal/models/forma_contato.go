package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// CONSTANTES
// ============================================================
// (Adicionar no final do constants.go)

/*
const (
    FormaContatoTipoTelefone   = 1
    FormaContatoTipoCelular    = 2
    FormaContatoTipoEmail      = 3
    FormaContatoTipoWhatsApp   = 4
    FormaContatoTipoSite       = 5
    FormaContatoTipoFacebook   = 6
    FormaContatoTipoInstagram  = 7
)
*/

// ============================================================
// MODEL: FormaContato
// ============================================================

type FormaContato struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID        int    `gorm:"column:frc_id;primaryKey;autoIncrement" json:"id"`
	Descricao string `gorm:"column:frc_descricao;type:varchar(60);not null" json:"descricao"`
	Tipo      int    `gorm:"column:frc_tipo;not null" json:"tipo"`

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
	Contatos []EntidadeContato `gorm:"foreignKey:FormaContatoID;references:ID" json:"contatos,omitempty"`
}

func (FormaContato) TableName() string {
	return "formacontato"
}

func (m *FormaContato) BeforeCreate(tx *gorm.DB) error {
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

func (m *FormaContato) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *FormaContato) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *FormaContato) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}

func (m *FormaContato) IsTelefone() bool {
	return m.Tipo == 1 || m.Tipo == 2
}

func (m *FormaContato) IsDigital() bool {
	return m.Tipo == 3 || m.Tipo == 4 || m.Tipo == 5
}
