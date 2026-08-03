package models

import (
	"time"
)

// ============================================================
// MODEL: NotaFiscalPagamentoCDC
// ============================================================

type NotaFiscalPagamentoCDC struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	NotaFiscalID            int      `gorm:"column:ntf_id;primaryKey" json:"nota_fiscal_id"`
	NotaFiscalPagamentoItem int      `gorm:"column:nfp_item;primaryKey" json:"nota_fiscal_pagamento_item"`
	CentroDeCustoID         int      `gorm:"column:cdc_id;primaryKey" json:"centro_de_custo_id"`
	Valor                   *float64 `gorm:"column:nfpcdc_valor;type:decimal(15,4)" json:"valor,omitempty"`
	Percentual              *float64 `gorm:"column:nfpcdc_percentual;type:decimal(5,2)" json:"percentual,omitempty"`

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
	NotaFiscalPagamento *NotaFiscalPagamento `gorm:"foreignKey:NotaFiscalID,NotaFiscalPagamentoItem;references:ntf_id,nfp_item" json:"nota_fiscal_pagamento,omitempty"`
	CentroDeCusto       *CentroDeCusto       `gorm:"foreignKey:CentroDeCustoID;references:cdc_id" json:"centro_de_custo,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (NotaFiscalPagamentoCDC) TableName() string {
	return "nota_fiscal_pagamento_cdc"
}
