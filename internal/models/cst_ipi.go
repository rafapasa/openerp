package models

import (
	"time"
)

// ============================================================
// MODEL: CSTIPI
// ============================================================

type CSTIPI struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID         int    `gorm:"column:cstipi_id;primaryKey" json:"id"`
	Codigo     string `gorm:"column:cstipi_codigo;type:varchar(2);not null" json:"codigo"`
	Descricao  string `gorm:"column:cstipi_descricao;type:varchar(255);not null" json:"descricao"`
	CalculaIPI int    `gorm:"column:cstipi_calculaipi;default:0" json:"calcula_ipi"`

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
	OperacoesFiscais []OperacaoFiscal `gorm:"foreignKey:CSTIPIID;references:ID" json:"operacoes_fiscais,omitempty"`
}

func (CSTIPI) TableName() string {
	return "cstipi"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *CSTIPI) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *CSTIPI) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}

func (m *CSTIPI) IsTributado() bool {
	return m.CalculaIPI == 1
}
