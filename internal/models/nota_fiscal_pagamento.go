package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: NotaFiscalPagamento
// ============================================================

type NotaFiscalPagamento struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	NotaFiscalID     int  `gorm:"column:ntf_id;primaryKey" json:"nota_fiscal_id"`
	Item             int  `gorm:"column:nfp_item;primaryKey" json:"item"`
	RotinaContabilID *int `gorm:"column:roc_id" json:"rotina_contabil_id,omitempty"`
	PortadorID       int  `gorm:"column:por_id;not null" json:"portador_id"`
	TipoDocumentoID  int  `gorm:"column:tdoc_id;not null" json:"tipo_documento_id"`

	DataVencimento   time.Time `gorm:"column:nfp_datavencimento;type:date;not null" json:"data_vencimento"`
	Valor            *float64  `gorm:"column:nfp_valor;type:decimal(15,2)" json:"valor,omitempty"`
	FormaPagamentoID *int      `gorm:"column:frmpgto_id" json:"forma_pagamento_id,omitempty"`

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
	NotaFiscal     *NotaFiscal     `gorm:"foreignKey:NotaFiscalID;references:ntf_id" json:"nota_fiscal,omitempty"`
	RotinaContabil *RotinaContabil `gorm:"foreignKey:RotinaContabilID;references:roc_id" json:"rotina_contabil,omitempty"`
	Portador       *Portador       `gorm:"foreignKey:PortadorID;references:por_id" json:"portador,omitempty"`
	TipoDocumento  *TipoDocumento  `gorm:"foreignKey:TipoDocumentoID;references:tdoc_id" json:"tipo_documento,omitempty"`
	FormaPagamento *FormaPagamento `gorm:"foreignKey:FormaPagamentoID;references:frmpgto_id" json:"forma_pagamento,omitempty"`

	CentroCustos []NotaFiscalPagamentoCDC `gorm:"foreignKey:NotaFiscalID,NotaFiscalPagamentoItem;references:ntf_id,nfp_item" json:"centro_custos,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (NotaFiscalPagamento) TableName() string {
	return "nota_fiscal_pagamento"
}

func (n *NotaFiscalPagamento) BeforeCreate(tx *gorm.DB) error {
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

func (n *NotaFiscalPagamento) BeforeUpdate(tx *gorm.DB) error {
	if n.UpdatedBy == nil {
		n.UpdatedBy = new(int)
		*n.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o pagamento foi deletado logicamente
func (n *NotaFiscalPagamento) IsDeleted() bool {
	return n.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (n *NotaFiscalPagamento) SoftDelete() {
	now := time.Now()
	n.DeletedAt = &now
}

// HasFormaPagamento verifica se possui forma de pagamento
func (n *NotaFiscalPagamento) HasFormaPagamento() bool {
	return n.FormaPagamentoID != nil && *n.FormaPagamentoID > 0
}

// HasCentroCustos verifica se possui centros de custo
func (n *NotaFiscalPagamento) HasCentroCustos() bool {
	return len(n.CentroCustos) > 0
}
