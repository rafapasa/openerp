package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// CONSTANTES
// ============================================================

// Constantes para Tipo de Forma de Pagamento
const (
	TipoFormaPagamentoCheque = 1
	TipoFormaPagamentoOutros = 9
)

// Constantes para Bandeira/Operadora de Cartão
const (
	BandeiraVisa            = 1
	BandeiraMastercard      = 2
	BandeiraAmericanExpress = 3
	BandeiraSorocred        = 4
	BandeiraDinersClub      = 5
	BandeiraElo             = 6
	BandeiraHipercard       = 7
	BandeiraAura            = 8
	BandeiraCabal           = 9
	BandeiraOutros          = 99
)

// ============================================================
// MODEL: FormaPagamento
// ============================================================

type FormaPagamento struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID                      int  `gorm:"column:frmpgto_id;primaryKey;autoIncrement" json:"id"`
	PlanoContasFinanceiroID *int `gorm:"column:pcf_id" json:"plano_contas_financeiro_id,omitempty"`
	EntidadeID              *int `gorm:"column:ent_id" json:"entidade_id,omitempty"`
	RotinaContabilID        *int `gorm:"column:roc_id" json:"rotina_contabil_id,omitempty"`

	Descricao          string  `gorm:"column:frmpgto_descricao;type:varchar(255);not null" json:"descricao"`
	Tipo               *int    `gorm:"column:frmpgto_tipo" json:"tipo,omitempty"` // 1-Cheque, 9-Outros
	VisivelFrenteCaixa int     `gorm:"column:frmpgto_visivelfrentecaixa;not null;default:0" json:"visivel_frente_caixa"`
	BandeiraOperadora  *int    `gorm:"column:frmpgto_bandeiraoperadora" json:"bandeira_operadora,omitempty"`
	TeclaAtalho        *int    `gorm:"column:frmpgto_teclatalho" json:"tecla_atalho,omitempty"`
	Sigla              *string `gorm:"column:frmpgto_sigla;type:varchar(2)" json:"sigla,omitempty"`

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
	PlanoContasFinanceiro *PlanoContasFinanceiro `gorm:"foreignKey:PlanoContasFinanceiroID;references:pcf_id" json:"plano_contas_financeiro,omitempty"`
	Entidade              *Entidade              `gorm:"foreignKey:EntidadeID;references:ent_id" json:"entidade,omitempty"`
	RotinaContabil        *RotinaContabil        `gorm:"foreignKey:RotinaContabilID;references:roc_id" json:"rotina_contabil,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (FormaPagamento) TableName() string {
	return "forma_pagamento"
}

func (f *FormaPagamento) BeforeCreate(tx *gorm.DB) error {
	if f.CreatedBy == nil {
		f.CreatedBy = new(int)
		*f.CreatedBy = 0
	}
	if f.UpdatedBy == nil {
		f.UpdatedBy = new(int)
		*f.UpdatedBy = 0
	}
	return nil
}

func (f *FormaPagamento) BeforeUpdate(tx *gorm.DB) error {
	if f.UpdatedBy == nil {
		f.UpdatedBy = new(int)
		*f.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsCheque verifica se a forma de pagamento é cheque
func (f *FormaPagamento) IsCheque() bool {
	return f.Tipo != nil && *f.Tipo == TipoFormaPagamentoCheque
}

// IsVisivelNoCaixa verifica se está visível no frente de caixa
func (f *FormaPagamento) IsVisivelNoCaixa() bool {
	return f.VisivelFrenteCaixa == 1
}
