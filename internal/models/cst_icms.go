package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: CSTICMS
// ============================================================

type CSTICMS struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID             int    `gorm:"column:csticms_id;primaryKey" json:"id"`
	Codigo         string `gorm:"column:csticms_codigo;type:varchar(3);not null" json:"codigo"`
	Descricao      string `gorm:"column:csticms_descricao;type:varchar(255);not null" json:"descricao"`
	CalculaICMS    int    `gorm:"column:csticms_calculaicms;not null;default:0" json:"calcula_icms"`
	ReduzBCICMS    int    `gorm:"column:csticms_reduz_bc_icms;not null;default:0" json:"reduz_bc_icms"`
	CalculaICMS_ST int    `gorm:"column:csticms_calculaicms_st;not null;default:0" json:"calcula_icms_st"`
	FCP            int    `gorm:"column:csticms_fcp;not null;default:0" json:"fcp"`

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
	OperacoesFiscais []OperacaoFiscal `gorm:"foreignKey:CSTICMSID;references:ID" json:"operacoes_fiscais,omitempty"`
}

func (CSTICMS) TableName() string {
	return "csticms"
}

func (m *CSTICMS) BeforeCreate(tx *gorm.DB) error {
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

func (m *CSTICMS) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *CSTICMS) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *CSTICMS) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}

func (m *CSTICMS) IsTributado() bool {
	return m.CalculaICMS == 1
}

func (m *CSTICMS) IsST() bool {
	return m.CalculaICMS_ST == 1
}

func (m *CSTICMS) IsFCP() bool {
	return m.FCP == 1
}
