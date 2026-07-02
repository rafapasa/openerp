package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: DocumentoVendaPagamentoCDC
// ============================================================

type DocumentoVendaPagamentoCDC struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	CentroDeCustoID     int      `gorm:"column:cdc_id;primaryKey" json:"centro_de_custo_id"`
	DocumentoVendaID    int      `gorm:"column:ddv_id;primaryKey" json:"documento_venda_id"`
	DocumentoVendaPagamentoItem int `gorm:"column:dvp_item;primaryKey" json:"documento_venda_pagamento_item"`
	Valor               *float64 `gorm:"column:dvpcdc_valor;type:decimal(15,4)" json:"valor,omitempty"`
	Percentual          *float64 `gorm:"column:dvpcdc_percentual;type:decimal(5,2)" json:"percentual,omitempty"`

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
	DocumentoVendaPagamento *DocumentoVendaPagamento `gorm:"foreignKey:DocumentoVendaID,DocumentoVendaPagamentoItem;references:ddv_id,dvp_item" json:"documento_venda_pagamento,omitempty"`
	CentroDeCusto           *CentroDeCusto           `gorm:"foreignKey:CentroDeCustoID;references:cdc_id" json:"centro_de_custo,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (DocumentoVendaPagamentoCDC) TableName() string {
	return "documento_venda_pagamento_cdc"
}

func (d *DocumentoVendaPagamentoCDC) BeforeCreate(tx *gorm.DB) error {
	if d.CreatedBy == nil {
		d.CreatedBy = new(int)
		*d.CreatedBy = 0
	}
	if d.UpdatedBy == nil {
		d.UpdatedBy = new(int)
		*d.UpdatedBy = 0
	}
	return nil
}

func (d *DocumentoVendaPagamentoCDC) BeforeUpdate(tx *gorm.DB) error {
	if d.UpdatedBy == nil {
		d.UpdatedBy = new(int)
		*d.UpdatedBy = 0
	}
	return nil
}