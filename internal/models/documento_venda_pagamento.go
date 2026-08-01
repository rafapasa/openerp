package models

import (
	"time"

	"gorm.io/gorm"
)

// =====================================================d=======
// MODEL: DocumentoVendaPagamento
// ============================================================

type DocumentoVendaPagamento struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	DocumentoVendaID int  `gorm:"column:ddv_id;primaryKey" json:"documento_venda_id"`
	Item             int  `gorm:"column:dvp_item;primaryKey" json:"item"`
	PortadorID       int  `gorm:"column:por_id;not null" json:"portador_id"`
	TipoDocumentoID  int  `gorm:"column:tdoc_id;not null" json:"tipo_documento_id"`
	FormaPagamentoID *int `gorm:"column:frmpgto_id" json:"forma_pagamento_id,omitempty"` // Changed to omitempty as it's a pointer
	RotinaContabilID *int `gorm:"column:roc_id" json:"rotina_contabil_id,omitempty"`

	// ============================================================
	// CAMPOS DO PAGAMENTO
	// ============================================================
	DataVencimento    time.Time `gorm:"column:dvp_datavencimento;type:date;not null" json:"data_vencimento"`
	Valor             float64   `gorm:"column:dvp_valor;type:decimal(15,2);not null" json:"valor"`
	Documento         *string   `gorm:"column:dvp_documento;type:varchar(30)" json:"documento,omitempty"`
	CodigoAutorizacao *string   `gorm:"column:dvp_codautorizacao;type:varchar(15)" json:"codigo_autorizacao,omitempty"`
	COO               *int      `gorm:"column:dvp_coo" json:"coo,omitempty"`

	// ============================================================
	// CAMPOS DE CARTÃO (para pagamentos com cartão)
	// ============================================================
	BandeiraCartao *string `gorm:"column:dvp_bandeira_cartao;type:varchar(50)" json:"bandeira_cartao,omitempty"`
	NSU            *string `gorm:"column:dvp_nsu;type:varchar(20)" json:"nsu,omitempty"`
	Autorizacao    *string `gorm:"column:dvp_autorizacao;type:varchar(20)" json:"autorizacao,omitempty"`

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
	DocumentoVenda *DocumentoVenda `gorm:"foreignKey:DocumentoVendaID;references:ID" json:"documento_venda,omitempty"`
	Portador       *Portador       `gorm:"foreignKey:PortadorID;references:por_id" json:"portador,omitempty"`
	TipoDocumento  *TipoDocumento  `gorm:"foreignKey:TipoDocumentoID;references:tdoc_id" json:"tipo_documento,omitempty"`
	FormaPagamento *FormaPagamento `gorm:"foreignKey:FormaPagamentoID;references:frmpgto_id" json:"forma_pagamento,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (DocumentoVendaPagamento) TableName() string {
	return "documento_venda_pagamento"
}

func (d *DocumentoVendaPagamento) BeforeCreate(tx *gorm.DB) error {
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

func (d *DocumentoVendaPagamento) BeforeUpdate(tx *gorm.DB) error {
	if d.UpdatedBy == nil {
		d.UpdatedBy = new(int)
		*d.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsParcelado verifica se o pagamento é parcelado
// (identifica pela quantidade de parcelas no documento)
func (d *DocumentoVendaPagamento) IsParcelado() bool {
	// TODO: Implementar lógica para verificar se é parcelado
	// Isso pode ser feito contando os pagamentos do mesmo documento
	return false
}

// GetValorFormatado retorna o valor formatado como string
func (d *DocumentoVendaPagamento) GetValorFormatado() string {
	return string(rune(d.Valor)) // TODO: Formatar corretamente
}

func (d *DocumentoVendaPagamento) IsDeleted() bool {
	return d.DeletedAt != nil
}

func (d *DocumentoVendaPagamento) SoftDelete() {
	d.DeletedAt = new(time.Time)
	*d.DeletedAt = time.Now()

}
