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
    TipoContatoTelefone   = 1
    TipoContatoCelular    = 2
    TipoContatoEmail      = 3
    TipoContatoWhatsApp   = 4
    TipoContatoSite       = 5
    TipoContatoFacebook   = 6
    TipoContatoInstagram  = 7
)
*/

// ============================================================
// MODEL: EntidadeContato
// ============================================================

type EntidadeContato struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	EntidadeID     int     `gorm:"column:ent_id;primaryKey" json:"entidade_id"`
	Item           int     `gorm:"column:efc_item;primaryKey" json:"item"`
	FormaContatoID int     `gorm:"column:frc_id;not null" json:"forma_contato_id"`
	Informacao     string  `gorm:"column:efc_informacao;type:text;not null" json:"informacao"`
	Descricao      *string `gorm:"column:efc_descricao;type:text" json:"descricao,omitempty"`

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
	Entidade     *Entidade     `gorm:"foreignKey:EntidadeID;references:ent_id" json:"entidade,omitempty"`
	FormaContato *FormaContato `gorm:"foreignKey:FormaContatoID;references:frc_id" json:"forma_contato,omitempty"`
}

func (EntidadeContato) TableName() string {
	return "entidade_formacontato"
}

func (m *EntidadeContato) BeforeCreate(tx *gorm.DB) error {
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

func (m *EntidadeContato) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *EntidadeContato) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *EntidadeContato) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}

// IsActive verifica se o contato está ativo (não deletado)
func (m *EntidadeContato) IsActive() bool {
	return m.DeletedAt == nil
}

func (m *EntidadeContato) IsWhatsApp() bool {
	return m.FormaContatoID == 4
}

func (m *EntidadeContato) IsEmail() bool {
	return m.FormaContatoID == 3
}
