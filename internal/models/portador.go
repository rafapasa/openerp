package models

import (
	"time"

	"github.com/openerp/backend/internal/constants"
	"gorm.io/gorm"
)

type Portador struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID                 int                    `gorm:"column:por_id;primaryKey;autoIncrement" json:"id"`
	ConvenioBancarioID *int                   `gorm:"column:cvb_id" json:"convenio_bancario_id,omitempty"`
	EmpresaFilialID    int                    `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	Nome               string                 `gorm:"column:por_descricao;type:varchar(255);not null" json:"nome"`
	Tipo               constants.TipoPortador `gorm:"column:por_tipo;not null" json:"tipo"`
	Situacao           constants.Status       `gorm:"column:por_situacao;not null;default:1" json:"situacao"`

	CreatedAt time.Time  `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime;index" json:"deleted_at,omitempty"`
	CreatedBy *int       `gorm:"column:created_by" json:"created_by,omitempty"`
	UpdatedBy *int       `gorm:"column:updated_by" json:"updated_by,omitempty"`

	EmpresaFilial    *EmpresaFilial    `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
	ConvenioBancario *ConvenioBancario `gorm:"foreignKey:ConvenioBancarioID;references:cvb_id" json:"convenio_bancario,omitempty"`
}

func (Portador) TableName() string {
	return "portador"
}

func (m *Portador) BeforeCreate(tx *gorm.DB) error {
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

func (m *Portador) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *Portador) IsActive() bool {
	return m.Situacao.IsActive()
}

func (m *Portador) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *Portador) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
	m.Situacao = constants.StatusInativo
}

func (m *Portador) IsBanco() bool {
	return m.Tipo.IsBanco()
}

func (m *Portador) IsOutros() bool {
	return m.Tipo.IsOutros()
}
