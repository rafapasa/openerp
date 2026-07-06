package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// CONSTANTES
// ============================================================

// Constantes para Tipo de Condição de Pagamento
const (
	TipoCondicaoAVista       = 0
	TipoCondicaoAPrazo       = 1
	TipoCondicaoSemPagamento = 9
)

// ============================================================
// MODEL: CondicaoPagamento
// ============================================================

type CondicaoPagamento struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID               int  `gorm:"column:codpgt_id;primaryKey;autoIncrement" json:"id"`
	RotinaContabilID *int `gorm:"column:roc_id" json:"rotina_contabil_id,omitempty"`
	PortadorID       int  `gorm:"column:por_id;not null" json:"portador_id"`
	TipoDocumentoID  int  `gorm:"column:tdoc_id;not null" json:"tipo_documento_id"`
	FormaPagamentoID *int `gorm:"column:frmpgto_id" json:"forma_pagamento_id,omitempty"`

	Descricao      string   `gorm:"column:codpgt_descricao;type:varchar(255);not null" json:"descricao"`
	TipoCondicao   int      `gorm:"column:codpgt_tipocondicao;not null" json:"tipo_condicao"` // 0-À vista, 1-À prazo, 9-Sem pagamento
	NumeroParcelas int      `gorm:"column:codpgt_numparcelas;not null" json:"numero_parcelas"`
	DiasParcelas   *int     `gorm:"column:codpgt_diasparcelas" json:"dias_parcelas,omitempty"`
	DiaPagamento   *int     `gorm:"column:codpgt_diapagamento" json:"dia_pagamento,omitempty"`
	Entrada        int      `gorm:"column:codpgt_entrada;not null" json:"entrada"` // 0-Não, 1-Sim
	Juros          *float64 `gorm:"column:codpgt_juros;type:decimal(5,2)" json:"juros,omitempty"`
	Comissao       *float64 `gorm:"column:codpgt_comissao;type:decimal(5,2)" json:"comissao,omitempty"`
	Desconto       *float64 `gorm:"column:codpgt_desconto;type:decimal(5,2);default:0.00" json:"desconto,omitempty"`
	Situacao       int      `gorm:"column:codpgt_situacao;not null" json:"situacao"`

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
	Portador       *Portador       `gorm:"foreignKey:PortadorID;references:por_id" json:"portador,omitempty"`
	TipoDocumento  *TipoDocumento  `gorm:"foreignKey:TipoDocumentoID;references:tdoc_id" json:"tipo_documento,omitempty"`
	FormaPagamento *FormaPagamento `gorm:"foreignKey:FormaPagamentoID;references:frmpgto_id" json:"forma_pagamento,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (CondicaoPagamento) TableName() string {
	return "condicao_pagamento"
}

func (c *CondicaoPagamento) BeforeCreate(tx *gorm.DB) error {
	if c.CreatedBy == nil {
		c.CreatedBy = new(int)
		*c.CreatedBy = 0
	}
	if c.UpdatedBy == nil {
		c.UpdatedBy = new(int)
		*c.UpdatedBy = 0
	}
	return nil
}

func (c *CondicaoPagamento) BeforeUpdate(tx *gorm.DB) error {
	if c.UpdatedBy == nil {
		c.UpdatedBy = new(int)
		*c.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsAVista verifica se é à vista
func (c *CondicaoPagamento) IsAVista() bool {
	return c.TipoCondicao == TipoCondicaoAVista
}

// IsAPrazo verifica se é a prazo
func (c *CondicaoPagamento) IsAPrazo() bool {
	return c.TipoCondicao == TipoCondicaoAPrazo
}

// GetNumeroParcelas retorna o número de parcelas
func (c *CondicaoPagamento) GetNumeroParcelas() int {
	if c.IsAVista() {
		return 1
	}
	return c.NumeroParcelas
}
