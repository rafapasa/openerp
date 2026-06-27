package models

import (
	"time"
)

// ============================================================
// MODEL: Banco
// ============================================================

type Banco struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID              int     `gorm:"column:ban_id;primaryKey;autoIncrement" json:"id"`
	EmpresaFilialID int     `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	Codigo          int     `gorm:"column:ban_codigo;not null" json:"codigo"`
	Descricao       string  `gorm:"column:ban_descricao;type:varchar(255);not null" json:"descricao"`
	PastaRemessa    *string `gorm:"column:ban_pastaremessa;type:varchar(255)" json:"pasta_remessa,omitempty"`
	PastaRetorno    *string `gorm:"column:ban_pastaretorno;type:varchar(255)" json:"pasta_retorno,omitempty"`
	Situacao        int     `gorm:"column:ban_situacao;not null;default:1" json:"situacao"`

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
	EmpresaFilial *EmpresaFilial  `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
	Contas        []ContaBancaria `gorm:"foreignKey:BancoID;references:ID" json:"contas,omitempty"`
}

func (Banco) TableName() string {
	return "bancos"
}

func (m *Banco) BeforeCreate() error {
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

func (m *Banco) BeforeUpdate() error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *Banco) IsActive() bool {
	return m.Situacao == 1
}

func (m *Banco) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *Banco) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
	m.Situacao = 2
}

func (m *Banco) HasPastas() bool {
	return (m.PastaRemessa != nil && *m.PastaRemessa != "") ||
		(m.PastaRetorno != nil && *m.PastaRetorno != "")
}
