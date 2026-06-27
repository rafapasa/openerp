package models

import (
	"time"
)

// ============================================================
// CONSTANTES
// ============================================================

// Constantes para Tipo de Título
const (
	TipoTituloPagamento   = 1
	TipoTituloRecebimento = 2
)

// Constantes para Situação do Título
const (
	SituacaoTituloAberto    = 0
	SituacaoTituloSaldo     = 1
	SituacaoTituloProtecao  = 2
	SituacaoTituloCobranca  = 3
	SituacaoTituloLiquidado = 8
	SituacaoTituloCancelado = 9
)

// Constantes para Origem do Título
const (
	OrigemTituloPedidoCompra = 1
	OrigemTituloNotaFiscal   = 2
	OrigemTituloPedidoVenda  = 3
)

// Constantes para Compromisso do Título
const (
	CompromissoCompromisso = 1
	CompromissoPrevisao    = 2
)

// ============================================================
// MODEL: Titulo
// ============================================================

type Titulo struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID                int  `gorm:"column:tit_id;primaryKey;autoIncrement" json:"id"`
	TipoDocumentoID   int  `gorm:"column:tdoc_id;not null" json:"tipo_documento_id"`
	EntidadeID        *int `gorm:"column:ent_id" json:"entidade_id,omitempty"`
	EmpresaFilialID   int  `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	RotinaContabilID  *int `gorm:"column:roc_id" json:"rotina_contabil_id,omitempty"`
	PortadorID        int  `gorm:"column:por_id;not null" json:"portador_id"`
	RepresentanteID   *int `gorm:"column:ent_id_rep" json:"representante_id,omitempty"`
	RemessaBancariaID *int `gorm:"column:remban_id" json:"remessa_bancaria_id,omitempty"`
	RemessaBaixaID    *int `gorm:"column:remban_id_baixa" json:"remessa_baixa_id,omitempty"`
	DespesaID         *int `gorm:"column:desp_id" json:"despesa_id,omitempty"`
	ReceitaID         *int `gorm:"column:rec_id" json:"receita_id,omitempty"`
	TituloPaiID       *int `gorm:"column:tit_tit_id" json:"titulo_pai_id,omitempty"`

	Documento    string `gorm:"column:tit_documento;type:varchar(50);not null" json:"documento"`
	IDDocGerador *int   `gorm:"column:tit_iddocgerador" json:"id_doc_gerador,omitempty"`
	Origem       *int   `gorm:"column:tit_origem" json:"origem,omitempty"` // 1-PedidoCompra, 2-NotaFiscal, 3-PedidoVenda
	Tipo         int    `gorm:"column:tit_tipo;not null" json:"tipo"`      // 1-Pagamento, 2-Recebimento
	Comissao     int    `gorm:"column:tit_comissoa;not null;default:0" json:"comissao"`
	Compromisso  int    `gorm:"column:tit_compromisso;not null" json:"compromisso"` // 1-Compromisso, 2-Previsão
	Modalidade   *int   `gorm:"column:tit_modalidade" json:"modalidade,omitempty"`
	Ocorrencia   *int   `gorm:"column:tit_ocorrencia" json:"ocorrencia,omitempty"`

	// ============================================================
	// CAMPOS DE DATAS
	// ============================================================
	DataEmissao      time.Time  `gorm:"column:tit_dataemissao;type:date;not null" json:"data_emissao"`
	DataVencimento   time.Time  `gorm:"column:tit_datavencimento;type:date;not null" json:"data_vencimento"`
	DataVencOriginal time.Time  `gorm:"column:tit_datavencoriginal;type:date;not null" json:"data_venc_original"`
	DataAprovacao    *time.Time `gorm:"column:tit_dataaprovacao;type:date" json:"data_aprovacao,omitempty"`
	DiasDesconto     *int       `gorm:"column:tit_diasdesconto" json:"dias_desconto,omitempty"`

	// ============================================================
	// CAMPOS DE VALORES
	// ============================================================
	ValorTitulo   float64  `gorm:"column:tit_valortitulo;type:decimal(15,2);not null" json:"valor_titulo"`
	ValorOriginal float64  `gorm:"column:tit_valororiginal;type:decimal(15,2);not null" json:"valor_original"`
	ValorSaldo    *float64 `gorm:"column:tit_valorsaldo;type:decimal(15,2)" json:"valor_saldo,omitempty"`
	ValorLiberado float64  `gorm:"column:tit_valorliberado;type:decimal(5,2);not null;default:0.00" json:"valor_liberado"`

	// ============================================================
	// CAMPOS DE TAXAS
	// ============================================================
	TaxaJuros          float64  `gorm:"column:tit_taxajuros;type:decimal(5,2);not null;default:0.00" json:"taxa_juros"`
	Multa              float64  `gorm:"column:tit_multa;type:decimal(5,2);not null;default:0.00" json:"multa"`
	TaxaDesconto       float64  `gorm:"column:tit_taxadesconto;type:decimal(5,2);not null;default:0.00" json:"taxa_desconto"`
	PercentualComissao *float64 `gorm:"column:tit_percentualcomissao;type:decimal(5,2)" json:"percentual_comissao,omitempty"`

	// ============================================================
	// CAMPOS DE INFORMAÇÕES
	// ============================================================
	Situacao          int     `gorm:"column:tit_situacao;not null" json:"situacao"` // 0-Aberto, 1-Saldo, 2-Proteção, 3-Cobrança, 8-Liquidado, 9-Cancelado
	Observacao        *string `gorm:"column:tit_observacao;type:text" json:"observacao,omitempty"`
	NossoNumero       *string `gorm:"column:tit_nossonumero;type:varchar(50)" json:"nosso_numero,omitempty"`
	InstrucaoCobranca *string `gorm:"column:tit_instrucaocobranca;type:text" json:"instrucao_cobranca,omitempty"`
	EntRazaoSocial    *string `gorm:"column:tit_ent_razaosocial;type:varchar(255)" json:"ent_razao_social,omitempty"`

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
	TipoDocumento *TipoDocumento `gorm:"foreignKey:TipoDocumentoID;references:tdoc_id" json:"tipo_documento,omitempty"`
	Entidade      *Entidade      `gorm:"foreignKey:EntidadeID;references:ent_id" json:"entidade,omitempty"`
	Portador      *Portador      `gorm:"foreignKey:PortadorID;references:por_id" json:"portador,omitempty"`
	Representante *Entidade      `gorm:"foreignKey:RepresentanteID;references:ent_id" json:"representante,omitempty"`
	TituloPai     *Titulo        `gorm:"foreignKey:TituloPaiID;references:tit_id" json:"titulo_pai,omitempty"`
	Baixas        []TituloBaixa  `gorm:"foreignKey:TituloID;references:ID" json:"baixas,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (Titulo) TableName() string {
	return "titulos"
}

func (t *Titulo) BeforeCreate() error {
	if t.CreatedBy == nil {
		t.CreatedBy = new(int)
		*t.CreatedBy = 0
	}
	if t.UpdatedBy == nil {
		t.UpdatedBy = new(int)
		*t.UpdatedBy = 0
	}
	return nil
}

func (t *Titulo) BeforeUpdate() error {
	if t.UpdatedBy == nil {
		t.UpdatedBy = new(int)
		*t.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsAberto verifica se o título está aberto
func (t *Titulo) IsAberto() bool {
	return t.Situacao == SituacaoTituloAberto
}

// IsLiquidado verifica se o título está liquidado
func (t *Titulo) IsLiquidado() bool {
	return t.Situacao == SituacaoTituloLiquidado
}

// IsCancelado verifica se o título está cancelado
func (t *Titulo) IsCancelado() bool {
	return t.Situacao == SituacaoTituloCancelado
}

// IsRecebimento verifica se é título de recebimento
func (t *Titulo) IsRecebimento() bool {
	return t.Tipo == TipoTituloRecebimento
}

// IsPagamento verifica se é título de pagamento
func (t *Titulo) IsPagamento() bool {
	return t.Tipo == TipoTituloPagamento
}

// GetValorSaldo retorna o saldo atual do título
func (t *Titulo) GetValorSaldo() float64 {
	if t.ValorSaldo != nil {
		return *t.ValorSaldo
	}
	return t.ValorTitulo
}

// CalcularValorComJuros calcula o valor com juros aplicados
func (t *Titulo) CalcularValorComJuros(diasAtraso int) float64 {
	if diasAtraso <= 0 {
		return t.GetValorSaldo()
	}
	juros := t.TaxaJuros / 100
	multa := t.Multa / 100
	valorComJuros := t.GetValorSaldo() * (1 + (juros * float64(diasAtraso)))
	valorComMulta := valorComJuros * (1 + multa)
	return valorComMulta
}

// SoftDelete realiza a exclusão lógica
func (t *Titulo) SoftDelete() {
	now := time.Now()
	t.DeletedAt = &now
	t.Situacao = SituacaoTituloCancelado
}
