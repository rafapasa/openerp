package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: CSTPISCOFINS
// ============================================================

type CSTPISCOFINS struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID            int    `gorm:"column:cstpiscofins_id;primaryKey" json:"id"`
	Codigo        string `gorm:"column:cstpiscofins_codigo;type:varchar(2);not null" json:"codigo"`
	Descricao     string `gorm:"column:cstpiscofins_descricao;type:varchar(255);not null" json:"descricao"`
	CalculaPIS    int    `gorm:"column:cstpiscofins_calculapis;not null;default:0" json:"calcula_pis"`
	CalculaCOFINS int    `gorm:"column:cstpiscofins_calculacofins;not null;default:0" json:"calcula_cofins"`

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
	OperacoesFiscais []OperacaoFiscal `gorm:"foreignKey:CSTPISCOFINSID;references:ID" json:"operacoes_fiscais,omitempty"`
}

func (CSTPISCOFINS) TableName() string {
	return "cstpiscofins"
}

func (m *CSTPISCOFINS) BeforeCreate(tx *gorm.DB) error {
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

func (m *CSTPISCOFINS) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *CSTPISCOFINS) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *CSTPISCOFINS) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}

func (m *CSTPISCOFINS) IsTributadoPIS() bool {
	return m.CalculaPIS == 1
}

func (m *CSTPISCOFINS) IsTributadoCOFINS() bool {
	return m.CalculaCOFINS == 1
}
