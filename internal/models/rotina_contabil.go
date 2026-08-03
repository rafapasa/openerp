package models

import (
	"time"

	"github.com/openerp/backend/internal/constants"
)

// ============================================================
// MODEL: RotinaContabil
// ============================================================

type RotinaContabil struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID              int              `gorm:"column:roc_id;primaryKey;autoIncrement" json:"id"`
	EmpresaFilialID int              `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	SetorID         *int             `gorm:"column:set_id" json:"setor_id,omitempty"`
	Descricao       string           `gorm:"column:roc_descricao;type:varchar(255);not null" json:"descricao"`
	Situacao        constants.Status `gorm:"column:roc_situacao;not null;default:1" json:"situacao"`

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
	EmpresaFilial *EmpresaFilial `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
	Setor         *Setor         `gorm:"foreignKey:SetorID;references:set_id" json:"setor,omitempty"`
}

func (RotinaContabil) TableName() string {
	return "rotina_contabil"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *RotinaContabil) IsActive() bool {
	return m.Situacao.IsActive()
}

func (m *RotinaContabil) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *RotinaContabil) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
	m.Situacao = constants.StatusInativo
}
