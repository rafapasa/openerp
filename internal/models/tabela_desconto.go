package models

import (
	"time"
)

// ============================================================
// MODEL: TabelaDesconto
// ============================================================

type TabelaDesconto struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID                int      `gorm:"column:tdesc_id;primaryKey;autoIncrement" json:"id"`
	Descricao         string   `gorm:"column:tdesc_descricao;type:varchar(255);not null" json:"descricao"`
	PercentualProduto float64  `gorm:"column:tdesc_percproduto;type:decimal(5,2);not null;default:0.00" json:"percentual_produto"`
	PercentualServico *float64 `gorm:"column:tdesc_percservico;type:decimal(5,2)" json:"percentual_servico,omitempty"`

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
	Entidades []Entidade `gorm:"foreignKey:TabelaDescontoID;references:ID" json:"entidades,omitempty"`
}

func (TabelaDesconto) TableName() string {
	return "tabela_desconto"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *TabelaDesconto) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *TabelaDesconto) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}

func (m *TabelaDesconto) HasDescontoServico() bool {
	return m.PercentualServico != nil && *m.PercentualServico > 0
}

func (m *TabelaDesconto) GetDescontoMaximo() float64 {
	if m.PercentualProduto > 0 {
		return m.PercentualProduto
	}
	return 0
}
