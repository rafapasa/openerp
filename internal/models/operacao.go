package models

import (
	"time"
)

// ============================================================
// MODEL: Operacao
// ============================================================

type Operacao struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID        int    `gorm:"column:ope_id;primaryKey;autoIncrement" json:"id"`
	Descricao string `gorm:"column:ope_descricao;type:varchar(255);not null" json:"descricao"`

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
	Processos []ProcessoOperacao `gorm:"foreignKey:OperacaoID;references:ope_id" json:"processos,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (Operacao) TableName() string {
	return "operacao"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se a operação foi deletada logicamente
func (o *Operacao) IsDeleted() bool {
	return o.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (o *Operacao) SoftDelete() {
	now := time.Now()
	o.DeletedAt = &now
}

// IsActive verifica se a operação está ativa
func (o *Operacao) IsActive() bool {
	return o.DeletedAt == nil
}

// HasProcessos verifica se a operação possui processos associados
func (o *Operacao) HasProcessos() bool {
	return len(o.Processos) > 0
}
