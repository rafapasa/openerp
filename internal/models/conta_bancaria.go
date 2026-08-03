package models

import (
	"time"
)

// ============================================================
// MODEL: ContaBancaria
// ============================================================

type ContaBancaria struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID            int     `gorm:"column:ctb_id;primaryKey;autoIncrement" json:"id"`
	BancoID       int     `gorm:"column:ban_id;not null" json:"banco_id"`
	Agencia       *string `gorm:"column:ctb_agencia;type:varchar(10)" json:"agencia,omitempty"`
	NumeroConta   *string `gorm:"column:ctb_numeroconta;type:varchar(10)" json:"numero_conta,omitempty"`
	CodigoCedente *string `gorm:"column:ctb_codigocedente;type:varchar(20)" json:"codigo_cedente,omitempty"`
	DigitoConta   *string `gorm:"column:ctb_digitoconta;type:varchar(2)" json:"digito_conta,omitempty"`
	DigitoAgencia *string `gorm:"column:ctb_digitoagencia;type:varchar(2)" json:"digito_agencia,omitempty"`

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
	Banco *Banco `gorm:"foreignKey:BancoID;references:ban_id" json:"banco,omitempty"`
}

func (ContaBancaria) TableName() string {
	return "conta_bancaria"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *ContaBancaria) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *ContaBancaria) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}

// IsActive verifica se a conta bancária está ativa (não deletada)
func (m *ContaBancaria) IsActive() bool {
	return m.DeletedAt == nil
}

func (m *ContaBancaria) HasAgencia() bool {
	return m.Agencia != nil && *m.Agencia != ""
}

func (m *ContaBancaria) HasNumeroConta() bool {
	return m.NumeroConta != nil && *m.NumeroConta != ""
}
