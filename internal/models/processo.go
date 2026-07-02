package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/constants"
)

// ============================================================
// MODEL: Processo
// ============================================================

type Processo struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID               int              `gorm:"column:prc_id;primaryKey;autoIncrement" json:"id"`
	EmpresaFilialID  int              `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	RotinaContabilID *int             `gorm:"column:roc_id" json:"rotina_contabil_id,omitempty"`
	Codigo           int              `gorm:"column:prc_codigo;not null" json:"codigo"`
	Descricao        string           `gorm:"column:prc_descricao;type:varchar(255);not null" json:"descricao"`
	Situacao         constants.Status `gorm:"column:prc_situacao;not null;default:1" json:"situacao"`

	OperacaoFiscalForaEstID   *int `gorm:"column:opf_id_fora_est" json:"operacao_fiscal_fora_est_id,omitempty"`
	OperacaoFiscalNoEstID     *int `gorm:"column:opf_id_no_est" json:"operacao_fiscal_no_est_id,omitempty"`
	OperacaoFiscalForaEstSTID *int `gorm:"column:opf_id_fora_est_st" json:"operacao_fiscal_fora_est_st_id,omitempty"`
	OperacaoFiscalNoEstSTID   *int `gorm:"column:opf_id_no_est_st" json:"operacao_fiscal_no_est_st_id,omitempty"`

	PlanoContasFinanceiroID *int `gorm:"column:pcf_id" json:"plano_contas_financeiro_id,omitempty"`
	NaturezaOperacaoID      *int `gorm:"column:nfno_id" json:"natureza_operacao_id,omitempty"`
	ReceitaID               *int `gorm:"column:rec_id" json:"receita_id,omitempty"`
	DespesaID               *int `gorm:"column:desp_id" json:"despesa_id,omitempty"`

	Devolucao    *int `gorm:"column:prc_devolucao;default:0" json:"devolucao,omitempty"`
	ProcessoNFID *int `gorm:"column:prc_id_nf" json:"processo_nf_id,omitempty"`
	Comodato     *int `gorm:"column:prc_comodato" json:"comodato,omitempty"`
	TipoOperacao *int `gorm:"column:prc_tipooperacao" json:"tipo_operacao,omitempty"` // 0 - entrada, 1 - saida
	Funrural     *int `gorm:"column:prc_funrural" json:"funrural,omitempty"`

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
	EmpresaFilial           *EmpresaFilial              `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
	RotinaContabil          *RotinaContabil             `gorm:"foreignKey:RotinaContabilID;references:roc_id" json:"rotina_contabil,omitempty"`
	OperacaoFiscalForaEst   *OperacaoFiscal             `gorm:"foreignKey:OperacaoFiscalForaEstID;references:opf_id" json:"operacao_fiscal_fora_est,omitempty"`
	OperacaoFiscalNoEst     *OperacaoFiscal             `gorm:"foreignKey:OperacaoFiscalNoEstID;references:opf_id" json:"operacao_fiscal_no_est,omitempty"`
	OperacaoFiscalForaEstST *OperacaoFiscal             `gorm:"foreignKey:OperacaoFiscalForaEstSTID;references:opf_id" json:"operacao_fiscal_fora_est_st,omitempty"`
	OperacaoFiscalNoEstST   *OperacaoFiscal             `gorm:"foreignKey:OperacaoFiscalNoEstSTID;references:opf_id" json:"operacao_fiscal_no_est_st,omitempty"`
	PlanoContasFinanceiro   *PlanoContasFinanceiro      `gorm:"foreignKey:PlanoContasFinanceiroID;references:pcf_id" json:"plano_contas_financeiro,omitempty"`
	NaturezaOperacao        *NotaFiscalNaturezaOperacao `gorm:"foreignKey:NaturezaOperacaoID;references:nfno_id" json:"natureza_operacao,omitempty"`
	Receita                 *Receita                    `gorm:"foreignKey:ReceitaID;references:rec_id" json:"receita,omitempty"`
	Despesa                 *Despesa                    `gorm:"foreignKey:DespesaID;references:desp_id" json:"despesa,omitempty"`
	ProcessoNF              *Processo                   `gorm:"foreignKey:ProcessoNFID;references:prc_id" json:"processo_nf,omitempty"`

	DocumentoVendas []DocumentoVenda   `gorm:"foreignKey:ProcessoID;references:prc_id" json:"documento_vendas,omitempty"`
	NotaFiscais     []NotaFiscal       `gorm:"foreignKey:ProcessoID;references:prc_id" json:"nota_fiscais,omitempty"`
	Operacoes       []ProcessoOperacao `gorm:"foreignKey:ProcessoID;references:prc_id" json:"operacoes,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (Processo) TableName() string {
	return "processo"
}

func (p *Processo) BeforeCreate(tx *gorm.DB) error {
	if p.CreatedBy == nil {
		p.CreatedBy = new(int)
		*p.CreatedBy = 0
	}
	if p.UpdatedBy == nil {
		p.UpdatedBy = new(int)
		*p.UpdatedBy = 0
	}
	return nil
}

func (p *Processo) BeforeUpdate(tx *gorm.DB) error {
	if p.UpdatedBy == nil {
		p.UpdatedBy = new(int)
		*p.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsActive verifica se o processo está ativo
func (p *Processo) IsActive() bool {
	return p.Situacao == constants.StatusAtivo
}

// IsDeleted verifica se o processo foi deletado logicamente
func (p *Processo) IsDeleted() bool {
	return p.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (p *Processo) SoftDelete() {
	now := time.Now()
	p.DeletedAt = &now
	p.Situacao = constants.StatusInativo
}

// IsInactive verifica se o processo está inativo
func (p *Processo) IsInactive() bool {
	return p.Situacao == constants.StatusInativo
}

// IsEntrada verifica se é um processo de entrada
func (p *Processo) IsEntrada() bool {
	return p.TipoOperacao != nil && *p.TipoOperacao == 0
}

// IsSaida verifica se é um processo de saída
func (p *Processo) IsSaida() bool {
	return p.TipoOperacao != nil && *p.TipoOperacao == 1
}

// IsDevolucao verifica se é um processo de devolução
func (p *Processo) IsDevolucao() bool {
	return p.Devolucao != nil && *p.Devolucao == 1
}

// IsComodato verifica se é um processo de comodato
func (p *Processo) IsComodato() bool {
	return p.Comodato != nil && *p.Comodato == 1
}

// HasFunrural verifica se possui Funrural
func (p *Processo) HasFunrural() bool {
	return p.Funrural != nil && *p.Funrural == 1
}

// HasOperacaoFiscalForaEst verifica se possui operação fiscal fora do estado
func (p *Processo) HasOperacaoFiscalForaEst() bool {
	return p.OperacaoFiscalForaEstID != nil && *p.OperacaoFiscalForaEstID > 0
}

// HasOperacaoFiscalNoEst verifica se possui operação fiscal no estado
func (p *Processo) HasOperacaoFiscalNoEst() bool {
	return p.OperacaoFiscalNoEstID != nil && *p.OperacaoFiscalNoEstID > 0
}

// HasDocumentoVendas verifica se o processo possui documentos de venda
func (p *Processo) HasDocumentoVendas() bool {
	return len(p.DocumentoVendas) > 0
}

// HasNotaFiscais verifica se o processo possui notas fiscais
func (p *Processo) HasNotaFiscais() bool {
	return len(p.NotaFiscais) > 0
}

// GetNomeCompleto retorna o nome completo do processo
func (p *Processo) GetNomeCompleto() string {
	return p.Descricao + " (Código: " + string(rune(p.Codigo)) + ")"
}

// SafeToDelete verifica se o processo pode ser excluído
func (p *Processo) SafeToDelete() bool {
	if p.HasDocumentoVendas() {
		return false
	}
	if p.HasNotaFiscais() {
		return false
	}
	return true
}
