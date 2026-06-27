package constants

// ============================================================
// CONSTANTES DE SITUAÇÃO (Status)
// ============================================================

// Status representa os possíveis status de um registro
type Status int

const (
	StatusAtivo     Status = 1
	StatusInativo   Status = 2
	StatusBloqueado Status = 3
	StatusCancelado Status = 9
)

// String retorna a representação textual do status
func (s Status) String() string {
	switch s {
	case StatusAtivo:
		return "Ativo"
	case StatusInativo:
		return "Inativo"
	case StatusBloqueado:
		return "Bloqueado"
	case StatusCancelado:
		return "Cancelado"
	default:
		return "Desconhecido"
	}
}

// IsActive verifica se o status é ativo
func (s Status) IsActive() bool {
	return s == StatusAtivo
}

// IsInactive verifica se o status é inativo
func (s Status) IsInactive() bool {
	return s == StatusInativo
}

// IsBlocked verifica se o status é bloqueado
func (s Status) IsBlocked() bool {
	return s == StatusBloqueado
}

// IsCanceled verifica se o status é cancelado
func (s Status) IsCanceled() bool {
	return s == StatusCancelado
}

// ============================================================
// TIPO DE PESSOA
// ============================================================

type TipoPessoa int

const (
	TipoPessoaFisica   TipoPessoa = 1
	TipoPessoaJuridica TipoPessoa = 2
)

func (t TipoPessoa) String() string {
	switch t {
	case TipoPessoaFisica:
		return "Pessoa Física"
	case TipoPessoaJuridica:
		return "Pessoa Jurídica"
	default:
		return "Desconhecido"
	}
}

func (t TipoPessoa) IsFisica() bool {
	return t == TipoPessoaFisica
}

func (t TipoPessoa) IsJuridica() bool {
	return t == TipoPessoaJuridica
}

// ============================================================
// TIPO DE DOCUMENTO
// ============================================================

// TipoDocumento representa os tipos de documento financeiro
type TipoDocumento int

const (
	TipoDocumentoBoleto      TipoDocumento = 1
	TipoDocumentoDuplicata   TipoDocumento = 2
	TipoDocumentoPromissoria TipoDocumento = 3
	TipoDocumentoCheque      TipoDocumento = 4
	TipoDocumentoTitulo      TipoDocumento = 5
)

func (t TipoDocumento) String() string {
	switch t {
	case TipoDocumentoBoleto:
		return "Boleto Bancário"
	case TipoDocumentoDuplicata:
		return "Duplicata"
	case TipoDocumentoPromissoria:
		return "Nota Promissória"
	case TipoDocumentoCheque:
		return "Cheque"
	case TipoDocumentoTitulo:
		return "Título"
	default:
		return "Desconhecido"
	}
}

// ============================================================
// TIPO DE PORTADOR
// ============================================================

type TipoPortador int

const (
	TipoPortadorBanco  TipoPortador = 1
	TipoPortadorOutros TipoPortador = 9
)

func (t TipoPortador) String() string {
	switch t {
	case TipoPortadorBanco:
		return "Banco"
	case TipoPortadorOutros:
		return "Outros"
	default:
		return "Desconhecido"
	}
}

func (t TipoPortador) IsBanco() bool {
	return t == TipoPortadorBanco
}

func (t TipoPortador) IsOutros() bool {
	return t == TipoPortadorOutros
}

// ============================================================
// SITUAÇÃO DO PEDIDO
// ============================================================

type SituacaoPedido int

const (
	SituacaoPedidoAberto      SituacaoPedido = 1
	SituacaoPedidoEmAtividade SituacaoPedido = 2
	SituacaoPedidoFechado     SituacaoPedido = 3
	SituacaoPedidoCancelado   SituacaoPedido = 9
)

func (s SituacaoPedido) String() string {
	switch s {
	case SituacaoPedidoAberto:
		return "Aberto"
	case SituacaoPedidoEmAtividade:
		return "Em Atividade"
	case SituacaoPedidoFechado:
		return "Fechado"
	case SituacaoPedidoCancelado:
		return "Cancelado"
	default:
		return "Desconhecido"
	}
}

func (s SituacaoPedido) IsAberto() bool {
	return s == SituacaoPedidoAberto
}

func (s SituacaoPedido) IsEmAtividade() bool {
	return s == SituacaoPedidoEmAtividade
}

func (s SituacaoPedido) IsFechado() bool {
	return s == SituacaoPedidoFechado
}

func (s SituacaoPedido) IsCancelado() bool {
	return s == SituacaoPedidoCancelado
}

func (s SituacaoPedido) IsActive() bool {
	return s == SituacaoPedidoAberto || s == SituacaoPedidoEmAtividade
}

// ============================================================
// TIPO DE OPERAÇÃO
// ============================================================

type TipoOperacao int

const (
	TipoOperacaoEntrada TipoOperacao = 0
	TipoOperacaoSaida   TipoOperacao = 1
)

func (t TipoOperacao) String() string {
	switch t {
	case TipoOperacaoEntrada:
		return "Entrada"
	case TipoOperacaoSaida:
		return "Saída"
	default:
		return "Desconhecido"
	}
}

// ============================================================
// TIPO DE DOCUMENTO DE VENDA
// ============================================================

type TipoDocumentoVenda int

const (
	TipoDocumentoOrcamento TipoDocumentoVenda = 1
	TipoDocumentoPedido    TipoDocumentoVenda = 2
)

func (t TipoDocumentoVenda) String() string {
	switch t {
	case TipoDocumentoOrcamento:
		return "Orçamento"
	case TipoDocumentoPedido:
		return "Pedido"
	default:
		return "Desconhecido"
	}
}

// ============================================================
// CONSTANTES DO PLANO DE CONTAS FINANCEIRO
// ============================================================

const (
	PlanoContasTipoBanco = 1
	PlanoContasTipoOutro = 9

	PlanoContasEspecieSintetica = 1
	PlanoContasEspecieAnalitica = 2

	PlanoContasSituacaoAtivo     = 1
	PlanoContasSituacaoInativo   = 2
	PlanoContasSituacaoBloqueado = 3
)

// ============================================================
// DIAS DA SEMANA
// ============================================================

type DiaSemana int

const (
	DiaSemanaDomingo DiaSemana = 1
	DiaSemanaSegunda DiaSemana = 2
	DiaSemanaTerca   DiaSemana = 3
	DiaSemanaQuarta  DiaSemana = 4
	DiaSemanaQuinta  DiaSemana = 5
	DiaSemanaSexta   DiaSemana = 6
	DiaSemanaSabado  DiaSemana = 7
)

func (d DiaSemana) String() string {
	switch d {
	case DiaSemanaDomingo:
		return "Domingo"
	case DiaSemanaSegunda:
		return "Segunda-feira"
	case DiaSemanaTerca:
		return "Terça-feira"
	case DiaSemanaQuarta:
		return "Quarta-feira"
	case DiaSemanaQuinta:
		return "Quinta-feira"
	case DiaSemanaSexta:
		return "Sexta-feira"
	case DiaSemanaSabado:
		return "Sábado"
	default:
		return "Desconhecido"
	}
}

// ============================================================
// TIPO DE ENDEREÇO
// ============================================================

type TipoEndereco int

const (
	TipoEnderecoCobranca    TipoEndereco = 1
	TipoEnderecoEntrega     TipoEndereco = 2
	TipoEnderecoComercial   TipoEndereco = 3
	TipoEnderecoResidencial TipoEndereco = 4
)

func (t TipoEndereco) String() string {
	switch t {
	case TipoEnderecoCobranca:
		return "Cobrança"
	case TipoEnderecoEntrega:
		return "Entrega"
	case TipoEnderecoComercial:
		return "Comercial"
	case TipoEnderecoResidencial:
		return "Residencial"
	default:
		return "Desconhecido"
	}
}

// ============================================================
// TIPO DE CONTATO
// ============================================================

type TipoContato int

const (
	TipoContatoTelefone  TipoContato = 1
	TipoContatoCelular   TipoContato = 2
	TipoContatoEmail     TipoContato = 3
	TipoContatoWhatsApp  TipoContato = 4
	TipoContatoSite      TipoContato = 5
	TipoContatoFacebook  TipoContato = 6
	TipoContatoInstagram TipoContato = 7
)

func (t TipoContato) String() string {
	switch t {
	case TipoContatoTelefone:
		return "Telefone"
	case TipoContatoCelular:
		return "Celular"
	case TipoContatoEmail:
		return "E-mail"
	case TipoContatoWhatsApp:
		return "WhatsApp"
	case TipoContatoSite:
		return "Site"
	case TipoContatoFacebook:
		return "Facebook"
	case TipoContatoInstagram:
		return "Instagram"
	default:
		return "Desconhecido"
	}
}

// ============================================================
// REGIME TRIBUTÁRIO
// ============================================================

type RegimeTributario int

const (
	RegimeTributarioSimplesNacional RegimeTributario = 1
	RegimeTributarioLucroPresumido  RegimeTributario = 2
	RegimeTributarioLucroReal       RegimeTributario = 3
	RegimeTributarioMEI             RegimeTributario = 4
	RegimeTributarioIsento          RegimeTributario = 5
)

func (r RegimeTributario) String() string {
	switch r {
	case RegimeTributarioSimplesNacional:
		return "Simples Nacional"
	case RegimeTributarioLucroPresumido:
		return "Lucro Presumido"
	case RegimeTributarioLucroReal:
		return "Lucro Real"
	case RegimeTributarioMEI:
		return "MEI"
	case RegimeTributarioIsento:
		return "Isento"
	default:
		return "Desconhecido"
	}
}

// ============================================================
// SEXO
// ============================================================

type Sexo int

const (
	SexoMasculino Sexo = 1
	SexoFeminino  Sexo = 2
)

func (s Sexo) String() string {
	switch s {
	case SexoMasculino:
		return "Masculino"
	case SexoFeminino:
		return "Feminino"
	default:
		return "Desconhecido"
	}
}

func (s Sexo) IsMasculino() bool {
	return s == SexoMasculino
}

func (s Sexo) IsFeminino() bool {
	return s == SexoFeminino
}
