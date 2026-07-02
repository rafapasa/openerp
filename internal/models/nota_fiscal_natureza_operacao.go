package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: NotaFiscalNaturezaOperacao
// ============================================================

type NotaFiscalNaturezaOperacao struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID        int    `gorm:"column:nfno_id;primaryKey;autoIncrement" json:"id"`
	Descricao string `gorm:"column:nfno_descricao;type:varchar(60);not null" json:"descricao"`

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
	Processos   []Processo   `gorm:"foreignKey:NaturezaOperacaoID;references:nfno_id" json:"processos,omitempty"`
	NotaFiscais []NotaFiscal `gorm:"foreignKey:NaturezaOperacaoID;references:nfno_id" json:"nota_fiscais,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (NotaFiscalNaturezaOperacao) TableName() string {
	return "nota_fiscal_natureza_operacao"
}

func (n *NotaFiscalNaturezaOperacao) BeforeCreate(tx *gorm.DB) error {
	if n.CreatedBy == nil {
		n.CreatedBy = new(int)
		*n.CreatedBy = 0
	}
	if n.UpdatedBy == nil {
		n.UpdatedBy = new(int)
		*n.UpdatedBy = 0
	}
	return nil
}

func (n *NotaFiscalNaturezaOperacao) BeforeUpdate(tx *gorm.DB) error {
	if n.UpdatedBy == nil {
		n.UpdatedBy = new(int)
		*n.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se a natureza foi deletada logicamente
func (n *NotaFiscalNaturezaOperacao) IsDeleted() bool {
	return n.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (n *NotaFiscalNaturezaOperacao) SoftDelete() {
	now := time.Now()
	n.DeletedAt = &now
}

// IsActive verifica se a natureza está ativa
func (n *NotaFiscalNaturezaOperacao) IsActive() bool {
	return n.DeletedAt == nil
}

// HasProcessos verifica se a natureza possui processos associados
func (n *NotaFiscalNaturezaOperacao) HasProcessos() bool {
	return len(n.Processos) > 0
}

// HasNotaFiscais verifica se a natureza possui notas fiscais associadas
func (n *NotaFiscalNaturezaOperacao) HasNotaFiscais() bool {
	return len(n.NotaFiscais) > 0
}
