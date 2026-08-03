package models

import (
	"time"
)

// ============================================================
// MODEL: TituloBaixa
// ============================================================

type TituloBaixa struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID               int  `gorm:"column:tbx_item;primaryKey;autoIncrement" json:"id"`
	TituloID         int  `gorm:"column:tit_id;not null" json:"titulo_id"`
	RotinaContabilID *int `gorm:"column:roc_id" json:"rotina_contabil_id,omitempty"`
	FormaPagamentoID int  `gorm:"column:frmpgto_id;not null;default:0" json:"forma_pagamento_id"`
	DespesaID        *int `gorm:"column:desp_id" json:"despesa_id,omitempty"`
	ReceitaID        *int `gorm:"column:rec_id" json:"receita_id,omitempty"`

	// ============================================================
	// CAMPOS DE DATAS
	// ============================================================
	DataPagamento time.Time  `gorm:"column:tbx_datapagto;type:date;not null" json:"data_pagamento"`
	DataLiberacao *time.Time `gorm:"column:tbx_dataliberacao;type:date" json:"data_liberacao,omitempty"`

	// ============================================================
	// CAMPOS DE VALORES
	// ============================================================
	ValorJuros    float64 `gorm:"column:tbx_valorjuros;type:decimal(15,2);not null" json:"valor_juros"`
	ValorDesconto float64 `gorm:"column:tbx_valordesconto;type:decimal(15,2);not null" json:"valor_desconto"`
	ValorPago     float64 `gorm:"column:tbx_valorpago;type:decimal(15,2);not null" json:"valor_pago"`

	// ============================================================
	// CAMPOS DE INFORMAÇÕES
	// ============================================================
	Observacao *string `gorm:"column:tbx_observacao;type:varchar(255)" json:"observacao,omitempty"`

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
	Titulo         *Titulo         `gorm:"foreignKey:TituloID;references:tit_id" json:"titulo,omitempty"`
	FormaPagamento *FormaPagamento `gorm:"foreignKey:FormaPagamentoID;references:frmpgto_id" json:"forma_pagamento,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (TituloBaixa) TableName() string {
	return "titulo_baixas"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// GetTotalPago retorna o total pago (valor + juros - desconto)
func (t *TituloBaixa) GetTotalPago() float64 {
	return t.ValorPago + t.ValorJuros - t.ValorDesconto
}

// IsLiquidado verifica se a baixa liquidou o título
func (t *TituloBaixa) IsLiquidado() bool {
	return t.ValorPago > 0
}
