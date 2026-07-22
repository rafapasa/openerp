package constants

import (
	"fmt"
	"strconv"
)

// ============================================================
// CONSTANTES DE SITUAÇÃO (Status)
// ============================================================

const (
	TipoTituloPagamento   = 1
	TipoTituloRecebimento = 2
)

const (
	SituacaoTituloAberto    = 0
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

func (s Status) IsValid() error {
	if s == StatusAtivo || s == StatusInativo || s == StatusBloqueado || s == StatusCancelado {
		return nil
	}
	return fmt.Errorf("status inválido: %d. Valores válidos são 1 (Ativo), 2 (Inativo), 3 (Bloqueado), 9 (Cancelado)", s)
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

const (
	// Modelos de Documentos Fiscais
	MDFCodigo01 = "01" // Nota Fiscal (modelo tradicional em papel / formulário)
	MDFCodigo04 = "04" // Nota Fiscal de Produtor
	MDFCodigo06 = "06" // Nota Fiscal / Conta de Energia Elétrica
	MDFCodigo07 = "07" // Nota Fiscal de Serviço de Transporte
	MDFCodigo08 = "08" // Conhecimento de Transporte Rodoviário de Cargas (CTRC)
	MDFCodigo55 = "55" // NF-e (Nota Fiscal Eletrônica)
	MDFCodigo57 = "57" // CT-e (Conhecimento de Transporte Eletrônico)
	MDFCodigo58 = "58" // MDF-e (Manifesto Eletrônico de Documentos Fiscais)
	MDFCodigo59 = "59" // CF-e SAT (Cupom Fiscal Eletrônico)
	MDFCodigo63 = "63" // BP-e (Bilhete de Passagem Eletrônico)
	MDFCodigo65 = "65" // NFC-e (Nota Fiscal de Consumidor Eletrônica)
	MDFCodigo67 = "67" // CT-e OS (Conhecimento de Transporte Eletrônico para Outros Serviços)
)

// ModelosFiscaisEletronicos retorna todos os códigos de modelos eletrônicos
func ModelosFiscaisEletronicos() []string {
	return []string{
		MDFCodigo55,
		MDFCodigo57,
		MDFCodigo58,
		MDFCodigo59,
		MDFCodigo63,
		MDFCodigo65,
		MDFCodigo67,
	}
}

// ModelosFiscaisPapel retorna todos os códigos de modelos em papel
func ModelosFiscaisPapel() []string {
	return []string{
		MDFCodigo01,
		MDFCodigo04,
		MDFCodigo06,
		MDFCodigo07,
		MDFCodigo08,
	}
}

type SimNao int

const (
	SimNaoSim SimNao = 1
	SimNaoNao SimNao = 0
)

func (s SimNao) String() string {
	switch s {
	case SimNaoSim:
		return "Sim"
	case SimNaoNao:
		return "Não"
	default:
		return "Desconhecido"
	}
}

func (s SimNao) IsValid() error {
	if s == SimNaoSim || s == SimNaoNao {
		return nil
	}
	return fmt.Errorf("valor inválido: %d. Valores válidos são 0 (Não) ou 1 (Sim)", s)
}

// Constantes para Tipo de Entrega
const (
	TipoEntregaRetirada = "RETIRADA"
	TipoEntregaEntrega  = "ENTREGA"
	TipoEntregaLocal    = "LOCAL"
)

// ============================================================
// CONSTANTES DE CONDIÇÃO DE PAGAMENTO
// ============================================================

// TipoCondicao representa os tipos de condição de pagamento
type TipoCondicao int

const (
	TipoCondicaoAVista       TipoCondicao = 0
	TipoCondicaoAPrazo       TipoCondicao = 1
	TipoCondicaoSemPagamento TipoCondicao = 9
)

func (t TipoCondicao) String() string {
	switch t {
	case TipoCondicaoAVista:
		return "À Vista"
	case TipoCondicaoAPrazo:
		return "À Prazo"
	case TipoCondicaoSemPagamento:
		return "Sem Pagamento"
	default:
		return "Desconhecido"
	}
}

func (t TipoCondicao) IsAVista() bool {
	return t == TipoCondicaoAVista
}

func (t TipoCondicao) IsAPrazo() bool {
	return t == TipoCondicaoAPrazo
}

func (t TipoCondicao) IsSemPagamento() bool {
	return t == TipoCondicaoSemPagamento
}

func (t TipoCondicao) IsValid() error {
	switch t {
	case TipoCondicaoAVista, TipoCondicaoAPrazo, TipoCondicaoSemPagamento:
		return nil
	default:
		return fmt.Errorf("tipo_condicao inválido: %d. Valores válidos são 0 (À Vista), 1 (À Prazo), 9 (Sem Pagamento)", t)
	}
}

// ============================================================
// SITUAÇÃO DA CONDIÇÃO DE PAGAMENTO
// ============================================================

// SituacaoCondicaoPagamento representa a situação da condição de pagamento
type SituacaoCondicaoPagamento int

const (
	SituacaoCondicaoAtivo   SituacaoCondicaoPagamento = 1
	SituacaoCondicaoInativo SituacaoCondicaoPagamento = 2
)

func (s SituacaoCondicaoPagamento) String() string {
	switch s {
	case SituacaoCondicaoAtivo:
		return "Ativo"
	case SituacaoCondicaoInativo:
		return "Inativo"
	default:
		return "Desconhecido"
	}
}

func (s SituacaoCondicaoPagamento) IsAtivo() bool {
	return s == SituacaoCondicaoAtivo
}

func (s SituacaoCondicaoPagamento) IsInativo() bool {
	return s == SituacaoCondicaoInativo
}

func (s SituacaoCondicaoPagamento) IsValid() error {
	switch s {
	case SituacaoCondicaoAtivo, SituacaoCondicaoInativo:
		return nil
	default:
		return fmt.Errorf("situacao inválida: %d. Valores válidos são 1 (Ativo) ou 2 (Inativo)", s)
	}
}

// ============================================================
// CONSTANTES DE CONDIÇÃO DE PAGAMENTO (VALORES SIM/NAO)
// ============================================================

// Constantes para campo Entrada (0-Não, 1-Sim)
const (
	CondicaoEntradaNao = 0
	CondicaoEntradaSim = 1
)

// ============================================================
// CONSTANTES DE FORMA DE PAGAMENTO (ADICIONAR SE NECESSÁRIO)
// ============================================================

// FormaPagamento representa as formas de pagamento
type FormaPagamento int

const (
	FormaPagamentoDinheiro      FormaPagamento = 1
	FormaPagamentoCartao        FormaPagamento = 2
	FormaPagamentoBoleto        FormaPagamento = 3
	FormaPagamentoCheque        FormaPagamento = 4
	FormaPagamentoTransferencia FormaPagamento = 5
	FormaPagamentoOutros        FormaPagamento = 9
)

func (f FormaPagamento) String() string {
	switch f {
	case FormaPagamentoDinheiro:
		return "Dinheiro"
	case FormaPagamentoCartao:
		return "Cartão"
	case FormaPagamentoBoleto:
		return "Boleto"
	case FormaPagamentoCheque:
		return "Cheque"
	case FormaPagamentoTransferencia:
		return "Transferência"
	case FormaPagamentoOutros:
		return "Outros"
	default:
		return "Desconhecido"
	}
}

// ============================================
// CONSTANTES DE STRING
// ============================================

const (
	USER_FORNEC = "FORNECEDOR"
	PASS_FORNEC = "fornec123"

	TOTVS_USER = "wsIntegracao"
	TOTVS_PASS = "DqUPeBdWjKebKbZ"

	QUOTE_CHAR = "'"

	CAMPO_ROTINA_CONTABIL = "ROC_ID"

	PASTA_TEMP = "./temp/"

	FILE_ERROR_LOG           = "./log/ErrorLog.log"
	FILE_CONNECTION_MySQL    = "./ini/ConnectionStringMySQL.ini"
	FILE_CONNECTION_MSSQL    = "./ini/ConnectionStringMSSQL.ini"
	FILE_CONNECTION_Firebird = "./ini/ConnectionStringFirebird.ini"
	FILE_CONNECTION_Oracle   = "./ini/ConnectionStringOracle.ini"

	VALOR_NULL = "NULL"

	VERSAO_NFE_310 = "310"
	VERSAO_NFE_400 = "400"

	NAO_STR = "0"
	SIM_STR = "1"

	DDV_SITUACAO_ABERTO_STR      = "1"
	DDV_SITUACAO_EM_PROCESSO_STR = "2"
	DDV_SITUACAO_CONCLUIDA_STR   = "3"
	DDV_SITUACAO_FECHADO_STR     = "4"
	DDV_SITUACAO_CANCELADO_STR   = "9"

	NFSIT_ID_ABERTA_STR      = "1"
	NFSIT_ID_FRINALIZADA_STR = "2"
	NFSIT_ID_CANCELADA_STR   = "9"

	LOTE_NFSE = "LOTE_NFSE"

	TITULO_JUSTICA              = "./png/ocorrencia_balanca.png"
	TITULO_PROT_CREDITO         = "./png/ocorrencia_cartao.png"
	TITULO_JUSTICA_PROT_CREDITO = "./png/ocorrencia.png"
)

// ============================================
// EXTENSOES DE ARQUIVO
// ============================================

var EXT_IMAGEM = []string{".JPG", ".JPEG", ".PNG", ".BMP"}
var EXT_PDF = []string{".PDF", ".XPS"}

// ============================================
// CONSTANTES NUMÉRICAS
// ============================================

const (
	ATIVO     = 1
	INATIVO   = 2
	BLOQUEADO = 3

	MERCADORIA_PARA_REVENDA = "00"
	PRODUTO_ACABADO         = "04"

	NAO = 0
	SIM = 1

	DIAS_DO_MES = 30

	CREDITO = 1
	DEBITO  = -1
)

// ============================================
// CONSTANTES DE SEGURANÇA
// ============================================

const (
	SEGRC_PERMISSAO_NEGADO         = 0
	SEGRC_PERMISSAO_ALTERAR_AMAIS  = 1
	SEGRC_PERMISSAO_ALTERAR_AMENOS = -1
	SEGRC_PERMISSAO_ALTERAR        = 2

	SEGRL_PERMISSAO_VISIVEL   = 1
	SEGRL_PERMISSAO_INVISIVEL = 0
)

// ============================================
// MÓDULOS (MOD_SIGLA_*)
// ============================================

const (
	MOD_SIGLA_COMP       = 1  // Compras
	MOD_SIGLA_FIN        = 2  // Financeiro
	MOD_SIGLA_EST        = 3  // Estoque
	MOD_SIGLA_COME       = 4  // Comercial
	MOD_SIGLA_FIS        = 5  // Fiscal
	MOD_SIGLA_CONT       = 6  // Contábil
	MOD_SIGLA_PROD       = 7  // Produção
	MOD_SIGLA_EXP        = 8  // Expedição
	MOD_SIGLA_AMB        = 9  // Ambiente
	MOD_SIGLA_REL        = 10 // Relatorios
	MOD_SIGLA_PDV        = 11 // Frente de Caixa
	MOD_SIGLA_MOB_G      = 12 // Mobile Garçom
	MOD_SIGLA_EMAIL      = 13 // Config.de E-Mail
	MOD_SIGLA_CONFIG_PDV = 14 // Config.de integração PDV
	MOD_SIGLA_TRI        = 15 // Config. de tributação
	MOD_SIGLA_TOTVS      = 16 // Config para integração Totvs
	MOD_SIGLA_DOMINIO    = 17 // Config para integração Dominio
	MOD_SIGLA_LINX       = 18 // Config para integração DTEF - LINX
	MOD_SIGLA_IFOOD      = 19 // Config para integração IFOOD
)

// ============================================
// TIPOS DE PRODUTOS
// ============================================

const (
	TPP_CODIGO_MR  = "00" // Mercadoria Para Revenda
	TPP_CODIGO_MP  = "01" // Matéria-Prima
	TPP_CODIGO_EMB = "02" // Embalagem
	TPP_CODIGO_PP  = "03" // Produto em Processo
	TPP_CODIGO_PA  = "04" // Produto Acabado
	TPP_CODIGO_SP  = "05" // SubProduto
	TPP_CODIGO_PI  = "06" // Produto Intermediario
	TPP_CODIGO_MUC = "07" // Material de Uso e Consumo
	TPP_CODIGO_AI  = "08" // Ativo Imobilizado
	TPP_CODIGO_SER = "09" // Serviços
	TPP_CODIGO_OUI = "10" // Outros Insumos
	TPP_CODIGO_COM = "98" // Comodato
	TPP_CODIGO_OUT = "99" // Outros
)

// ============================================
// CONSTANTES DE ESTOQUE
// ============================================

const (
	PROEST_ORIGEM_PEDIDO_COMPRA = 1
	PROEST_ORIGEM_NOTA_FISCAL   = 2
	PROEST_ORIGEM_PEDIDO_VENDA  = 3
	PROEST_ORIGEM__IMP_PDV      = 4
	PROEST_ORIGEM_PEDIDO_RMA    = 11
	PROEST_ORIGEM_AJUSTE        = 0

	PROEST_TIPO_ENTRADA = 1
	PROEST_TIPO_SAIDA   = -1
)

// ============================================
// IMPOSTOS
// ============================================

const (
	IMP_ID_ICMS     = 1
	IMP_ID_ICMS_ST  = 2
	IMP_ID_IPI      = 3
	IMP_ID_PIS      = 4
	IMP_ID_COFINS   = 5
	IMP_ID_PISST    = 6
	IMP_ID_COFINSST = 7
	IMP_ID_ISS      = 8
	IMP_ID_II       = 9
)

// ============================================
// DOCUMENTO DE VENDA (DDV)
// ============================================

const (
	DDV_TIPODOCUMENTO_ORCAMENTO     = 1
	DDV_TIPODOCUMENTO_PEDIDO        = 2
	DDV_TIPODOCUMENTO_PEDIDO_RAPIDO = 3
	DDV_TIPODOCUMENTO_FRENTE_CAIXA  = 4
	DDV_TIPODOCUMENTO_PEDIDO_BALCAO = 5
	DDV_TIPODOCUMENTO_COMPRA        = 6
	DDV_TIPODOCUMENTO_ENTREGA       = 7
	DDV_TIPODOCUMENTO_IMP_PDV       = 8

	DDV_TIPOOPERACAO_ENTRADA = 0
	DDV_TIPOOPERACAO_SAIDA   = 1

	DDV_ASSOCIA_MESA = 1

	DDV_SITUACAO_ABERTO      = 1
	DDV_SITUACAO_EM_PROCESSO = 2
	DDV_SITUACAO_CONCLUIDA   = 3
	DDV_SITUACAO_FECHADO     = 4
	DDV_SITUACAO_PENDENTE    = 5
	DDV_SITUACAO_CANCELADO   = 9

	DDV_USODESTINATARIO_REVENDA = 1
	DDV_USODESTINATARIO_CONSUMO = 2

	DDV_TIPO_OPERACAO_ENTRADA = 0
	DDV_TIPO_OPERACAO_SAIDA   = 1

	DDV_TIPOFRETE_TERCEIROS    = 0
	DDV_TIPOFRETE_EMITENTE     = 1
	DDV_TIPOFRETE_DESTINATARIO = 2
	DDV_TIPOFRETE_SEM_COBRANCA = 9
)

// ============================================
// CARREGAMENTO (CAR)
// ============================================

const (
	CAR_SITUCAO_ABERTO      = 1
	CAR_SITUCAO_EMATIVIDADE = 2
	CAR_SITUCAO_FINALIZADO  = 8
	CAR_SITUCAO_CANCELADO   = 9
)

// ============================================
// PEDIDO DE COMPRA (PCS)
// ============================================

const (
	PCS_ID_ABERTO     = 1
	PCS_ID_SALDO      = 2
	PCS_ID_FINALIZADO = 3
)

// ============================================
// ENDEREÇO (ETE)
// ============================================

const (
	ETE_TIPO_RESIDENCIAL = 1
	ETE_TIPO_COMERCIAL   = 2
	ETE_TIPO_COBRANCA    = 3
	ETE_TIPO_ENTREGA     = 4
)

// ============================================
// MANIFESTO ELETRONICO (MDFE)
// ============================================

const (
	MDFE_PROCESSOEMISSAO_APP_CONTRIBUITE   = 0
	MDFE_PROCESSOEMISSAO_AVULSA_PELO_FISCO = 1
	MDFE_PROCESSOEMISSAO_AVULSA_PELO_SITE  = 2
	MDFE_PROCESSOEMISSAO_APP_FISCO         = 3

	MDFE_TIPOEMISSAO_NORMAL       = 1
	MDFE_TIPOEMISSAO_CONTINGENCIA = 2
	MDFE_TIPOEMISSAO_SCAN         = 3
	MDFE_TIPOEMISSAO_DPEC         = 4
	MDFE_TIPOEMISSAO_FSDA         = 5
	MDFE_TIPOEMISSAO_SVCAN        = 6
	MDFE_TIPOEMISSAO_SVCRS        = 7
	MDFE_TIPOEMISSAO_SVCSP        = 8
	MDFE_TIPOEMISSAO_OFFLINE      = 9

	MDFE_TIPOAMBIENTE_PRODUCAO    = 1
	MDFE_TIPOAMBIENTE_HOMOLOGACAO = 2

	MDFE_TIPOEMITENTE_TRANSPORTADORA         = 1
	MDFE_TIPOEMITENTE_TRANSP_CARGA_PROPRIA   = 2
	MDFE_TIPOEMITENTE_TRANSP_CTE_GLOBALIZADO = 3

	MDFE_MODAL_RODOVIARIO  = 1
	MDFE_MODAL_AEREO       = 2
	MDFE_MODAL_AQUAVIARIO  = 3
	MDFE_MODAL_FERROVIARIO = 4
)

// ============================================
// NOTA FISCAL (NTF)
// ============================================

const (
	NTF_PROCESSOEMISSAO_APP_CONTRIBUITE   = 0
	NTF_PROCESSOEMISSAO_AVULSA_PELO_FISCO = 1
	NTF_PROCESSOEMISSAO_AVULSA_PELO_SITE  = 2
	NTF_PROCESSOEMISSAO_APP_FISCO         = 3

	NTF_TIPOOPERACAO_ENTRADA = 0
	NTF_TIPOOPERACAO_SAIDA   = 1

	NTF_TIPOEMITENTE_EMISSAO_PROPRIA   = 0
	NTF_TIPOEMITENTE_EMISSAO_TERCEIROS = 1

	NTF_TIPOEMISSAO_NORMAL  = 1
	NTF_TIPOEMISSAO_FS      = 2
	NTF_TIPOEMISSAO_SCAN    = 3
	NTF_TIPOEMISSAO_DPEC    = 4
	NTF_TIPOEMISSAO_FSDA    = 5
	NTF_TIPOEMISSAO_SVCAN   = 6
	NTF_TIPOEMISSAO_SVCRS   = 7
	NTF_TIPOEMISSAO_SVCSP   = 8
	NTF_TIPOEMISSAO_OFFLINE = 9

	NTF_TIPOAMBIENTE_PRODUCAO    = 1
	NTF_TIPOAMBIENTE_HOMOLOGACAO = 2

	NTF_USODESTINATARIO_REVENDA = 1
	NTF_USODESTINATARIO_CONSUMO = 2

	NTF_TIPO_OPERACAO_ENTRADA = 0
	NTF_TIPO_OPERACAO_SAIDA   = 1

	NTF_FINALIDADE_NORMAL       = 1
	NTF_FINALIDADE_COMPLEMENTAR = 2
	NTF_FINALIDADE_AJUSTE       = 3
	NTF_FINALIDADE_DEVOLUCAO    = 4

	NTF_OBSERVACAO_CREDITO_ICMS1 = "Permite o aproveitamento de crédito de ICMS no valor de R$"
	NTF_OBSERVACAO_CREDITO_ICMS2 = ", corresp. à alíq. de 3,11%, nos termos do Art. 23 da LC 123."
	NTF_OBSERVACAO_CREDITO_IPI   = "Doc. emitido por ME/EPP optante pelo SIMPLES NACIONAL, não gera direito a crédito de fiscal de ISS e de IPI."

	NTF_TIPO_IMPRESSAO_SEMGERACAO    = 0
	NTF_TIPO_IMPRESSAO_RETRATO       = 1
	NTF_TIPO_IMPRESSAO_PAISAGEM      = 2
	NTF_TIPO_IMPRESSAO_SIMPLIFICADA  = 3
	NTF_TIPO_IMPRESSAO_AOCONSUMIDOR  = 4
	NTF_TIPO_IMPRESSAO_MSGELETRONICA = 5

	NTF_STATUSCANCELAMENTO_SUCESSO = 135
	NTF_TIPODOCUMENTO_NFE          = 2
	NTF_TIPODOCUMENTO_NFSE         = 1
)

// ============================================
// NOTA FISCAL IV (NFIV)
// ============================================

const (
	NFIV_TIPOOPERACAO_CONCESSIONARIA    = 0
	NFIV_TIPOOPERACAO_FATURAMENTODIRETO = 1
	NFIV_TIPOOPERACAO_VENDADIRETA       = 2
	NFIV_TIPOOPERACAO_OUTROS            = 3
)

// ============================================
// NOTA FISCAL SITUAÇÃO (NFSIT)
// ============================================

const (
	NFSIT_ID_ABERTA         = 1
	NFSIT_ID_FINALIZADA     = 2
	NFSIT_ID_TRANSMITIDA    = 3
	NFSIT_ID_NAOTRANSMITIDA = 8
	NFSIT_ID_CANCELADA      = 9
)

// ============================================
// NOTA FISCAL ITEM (NFI)
// ============================================

const (
	NFI_INDTOTALIZACAO_COMPOEM     = 0
	NFI_INDTOTALIZACAO_NAO_COMPOEM = 1
)

// ============================================
// CT-E
// ============================================

const (
	CT_RETIRA_SIM = 0
	CT_RETIRA_NAO = 1

	CT_TIPOIMPRESSAO_SEM_GERACAO    = 0
	CT_TIPOIMPRESSAO_RETRATO        = 1
	CT_TIPOIMPRESSAO_PAISAGEM       = 2
	CT_TIPOIMPRESSAO_SIMPLIFICADO   = 3
	CT_TIPOIMPRESSAO_NFCe           = 4
	CT_TIPOIMPRESSAO_MSG_ELETRONICA = 5

	CT_TIPOEMISSAO_NORMAL       = 0
	CT_TIPOEMISSAO_CONTINGENCIA = 1
	CT_TIPOEMISSAO_SCAN         = 2
	CT_TIPOEMISSAO_DPEC         = 3
	CT_TIPOEMISSAO_FSDA         = 4
	CT_TIPOEMISSAO_SVCAN        = 5
	CT_TIPOEMISSAO_SVCRS        = 6
	CT_TIPOEMISSAO_SVCSP        = 7
	CT_TIPOEMISSAO_OFF_LINE     = 8

	CT_TIPOAMBIENTE_PRODUCAO    = 0
	CT_TIPOAMBIENTE_HOMOLOCACAO = 1

	CT_TIPOCTE_NORMAL      = 0
	CT_TIPOCTE_COMPLEMENTO = 1
	CT_TIPOCTE_ANULACAO    = 2
	CT_TIPOCTE_SUBSTITUTO  = 3
	CT_TIPOCTE_GTVe        = 4

	CT_TIPOSERVICO_N0RMAL          = 0
	CT_TIPOSERVICO_SUB_CONTRATACAO = 1
	CT_TIPOSERVICO_REDESPACHO      = 2
	CT_TIPOSERVICO_INTERMEDIARIO   = 3
	CT_TIPOSERVICO_MULTI_MODAL     = 4
	CT_TIPOSERVICO_TRANSP_PESSOAS  = 5
	CT_TIPOSERVICO_TRSNP_VALORES   = 6
	CT_TIPOSERVICO_EXESSO_BAGAGEM  = 7
	CT_TIPOSERVICO_GTV             = 8

	CT_IND_IE_TOMADOR_CONTRIBUINTE     = 0
	CT_IND_IE_TOMADOR_ISENTO           = 1
	CT_IND_IE_TOMADOR_NAO_CONTRIBUINTE = 2

	CT_TIPOFRETE_PAGO    = 0
	CT_TIPOFRETE_A_PAGAR = 1
	CT_TIPOFRETE_OUTROS  = 2
)

// ============================================
// ORÇAMENTOS DE COMPRAS (ORS)
// ============================================

const (
	ORS_ID_BLOQUEADO  = 0
	ORS_ID_LIBERADO   = 1
	ORS_ID_SOLICITADO = 2
	ORS_ID_RECEBIDO   = 3
	ORS_ID_APROVADO   = 4
	ORS_ID_FINALIZADO = 5
	ORS_ID_CANCELADO  = 9
)

// ============================================
// ORDEM DE SERVIÇO ITEM (OSI)
// ============================================

const (
	OSI_ID_ABERTO    = 1
	OSI_ID_APROVADO  = 2
	OSI_ID_REPROVADO = 3
	OSI_ID_CANCELADO = 9
)

// ============================================
// ORDEM DE SERVIÇO (OTO)
// ============================================

const (
	OTO_ID_SOLICITACAO = 1
	OTO_ID_RECEPCAO    = 2
)

// ============================================
// CONDIÇÃO DE PAGAMENTO (CODPGT)
// ============================================

const (
	CODPGT_TIPOCONDICAO_AVISTA        = 1
	CODPGT_TIPOCONDICAO_APRAZO        = 2
	CODPGT_TIPOCONDICAO_SEM_PAGAMENTO = 9
)

// ============================================
// FORMAS DE PAGAMENTO (FRMPGTO)
// ============================================

const (
	FRMPGTO_ID_OUTRAS       = 1
	FRMPGTO_ID_DEPOSITO     = 2
	FRMPGTO_ID_CHEQUE       = 3
	FRMPGTO_ID_ADIANTAMENTO = 4
	FRMPGTO_ID_COBBANCARIA  = 5
	FRMPGTO_ID_RENEGOCIACAO = 6

	FRMPGTO_TIPO_CHEQUE        = 1
	FRMPGTO_TIPO_CREDIARIO     = 2
	FRMPGTO_TIPO_CARTAODEBITO  = 3
	FRMPGTO_TIPO_CARTAOCREDITO = 4
	FRMPGTO_TIPO_PIX           = 5
	FRMPGTO_TIPO_CONTRAVALE    = 6
	FRMPGTO_TIPO_ENTREGA       = 7
	FRMPGTO_TIPO_ORCAMENTO     = 8
	FRMPGTO_TIPO_OUTRAS        = 9
	FRMPGTO_TIPO_DEVOLUCAO     = 10
	FRMPGTO_TIPO_VOUCHER       = 11
	FRMPGTO_TIPO_DINHEIRO      = 12
)

// FRMPGTO_TIPO_TEF - Formas de pagamento TEF (usar para verificação)
var FRMPGTO_TIPO_TEF = []int{FRMPGTO_TIPO_CARTAODEBITO, FRMPGTO_TIPO_CARTAOCREDITO, FRMPGTO_TIPO_PIX}

// ============================================
// PLANO DE CONTAS FINANCEIRO (PCF)
// ============================================

const (
	PCF_TIPO_BANCOS   = 1
	PCF_TIPO_CHEQUES  = 2
	PCF_TIPO_DINHEIRO = 3
	PCF_TIPO_CARTAO   = 4
	PCF_TIPO_OUTRAS   = 9
)

// ============================================
// PLANO DE CONTAS CONTÁBIL (PCC)
// ============================================

const (
	PCC_TIPO_SINTETICA = 1
	PCC_TIPO_ANALITICA = 2

	PCC_OPERADORCONTA_CREDITO = 1
	PCC_OPERADORCONTA_DEBITO  = 2

	PCC_ORIGEMCONTA_ATIVO   = 1
	PCC_ORIGEMCONTA_PASSIVO = 2
	PCC_ORIGEMCONTA_DRE     = 3
)

// ============================================
// GCE / TITULOS
// ============================================

const (
	GCE_TIPO_RECEBIMENTO = 2
	GCE_TIPO_PAGAMENTO   = 1

	TIT_TIPO_RECEBIRMENTO = 2
	TIT_TIPO_PAGAMENTO    = 1

	TIT_SITUACAO_ABERTO             = 0
	TIT_SITUACAO_SALDO              = 1
	TIT_SITUACAO_PROETCAO_CREDITO   = 2
	TIT_SITUACAO_COBRANACA_JUDICIAL = 3
	TIT_SITUACAO_LIQUIDADO          = 8
	TIT_SITUACAO_CANCELADO          = 9

	TIT_ORIGEM_PEDIDO_COMPRA = 1
	TIT_ORIGEM_NOTA_FISCAL   = 2
	TIT_ORIGEM_PEDIDO_VENDA  = 3
	TIT_ORIGEM_GER_LOTE      = 4
	TIT_ORIGEM_FRENTE_CAIXA  = 5
	TIT_ORIGEM_RENEGOCIACAO  = 10
	TIT_ORIGEM_CONTRATO      = 12
)

// ============================================
// GCENF
// ============================================

const (
	GCENF_TIPO_NFE    = 1
	GCENF_TIPO_NFSE   = 2
	GCENF_TIPO_PEDIDO = 3
)

// ============================================
// ALTI
// ============================================

const (
	ALTI_SITUACAO_ATIVO    = 1
	ALTI_SITUACAO_SUSPENSO = 2
)

// ============================================
// CAIXA FINANCEIRO (CXFINT)
// ============================================

const (
	CXFINT_SITUACAO_ABERTO  = 1
	CXFINT_SITUACAO_FECHADO = 2
)

// ============================================
// CONFIGURAÇÕES - AMBIENTE
// ============================================

const (
	CONFIG_RELATORIO_MODAL       = 1
	CONFIG_FRENTE_DE_CAIXA       = 7
	CONFIG_RESTAURANTE           = 8
	CONFIG_DIAS_ANIVER           = 9
	CONFIG_MSG_IMP_DDV           = 10
	CONFIG_MSG_IMP_DDV_ANIVER    = 11
	CONFIG_AMBIENTE_FISCAL       = 12
	CONFIG_EXIBIR_REJEICAO       = 13
	CONFIG_TAMANHO_FONTE_GRID    = 14
	CONFIG_PESQUISAR_PELO_INICIO = 16
	CONFIG_CLINICA_VET           = 17
	CONFIG_DIST_BEBIDA           = 18
	CONFIG_EMAIL_CONTADOR        = 19
	CONFIG_NOME_CONTADOR         = 20
	CONFIG_USAR_CDC              = 21
	CONFIG_CEREALISTA            = 22
	CONFIG_GUINCHOS              = 23
	CONFIG_INTEGRACAO_TOTVS      = 24
)

// ============================================
// CONFIGURAÇÕES - COMPRAS
// ============================================

const (
	CONFIG_PRC_COMPRA             = 100
	CONFIG_USAR_DESCONTO_NO_CUSTO = 101
)

// ============================================
// CONFIGURAÇÕES - FINANCEIRO
// ============================================

const (
	CONFIG_TAXA_JUROS               = 200
	CONFIG_TAXA_DESCONTO            = 201
	CONFIG_PERCENTUAL_MULTA         = 202
	CONFIG_DIAS_DECONTO             = 203
	CONFIG_PASTAREMESSA             = 205
	CONFIG_PASTARETORNO             = 206
	CONFIG_PCF_CONTA_PAGAR          = 206
	CONFIG_PCF_CONTA_RECEBER        = 207
	CONFIG_PCF_MASCARA              = 208
	CONFIG_APROVACAO_TITULOS        = 212
	CONFIG_ID_FRMPGTO_COBBANCARIA   = 213
	CONFIG_ID_FRMPGTO_BAIXALOTE     = 214
	CONFIG_PASTA_BACKUP             = 215
	CONFIG_INSTRUC_COB_PADRAO       = 216
	CONFIG_ID_FRMPGTO_ADIANTAMENTO  = 217
	CONFIG_BOLETO_CARNE             = 218
	CONFIG_GCE_PRO_ID_EXTRA         = 219
	CONFIG_ID_FRMPGTO_CONTRATO      = 220
	CONFIG_ID_FRMPGTO_RENEGOCIACAO  = 221
	CONFIG_DIAS_BLOQ_INADIMPLENCIA  = 222
	CONFIG_MARCA_TITULOS_BAIXA_LOTE = 223
	CONFIG_ID_FRMPGTO_DEVOLUCAO     = 224
	CONFIG_USAR_TEF                 = 225
)

// ============================================
// CONFIGURAÇÕES - ESTOQUE
// ============================================

const (
	CONFIG_USAR_FICHA_TECNICA = 300
	CONFIG_VALIDAR_ESTOQUE    = 301
	CONFIG_PTP_ID_GUINCHO     = 302
	CONFIG_UND_ID_HORAS       = 303
)

// ============================================
// CONFIGURAÇÕES - COMERCIAIS
// ============================================

const (
	CONFIG_NFE_AMBIENTE                     = 401
	CONFIG_TABELA_PRECO_PADRAO              = 402
	CONFIG_FLUXO_PADRAO                     = 403
	CONFIG_CONDPGTO_PADRAO                  = 404
	CONFIG_DDV_LOCALIZACAO_PADRAO           = 405
	CONFIG_AUTO_GRAVAR_ITEM                 = 406
	CONFIG_USAR_CODIGOBARRAS                = 407
	CONFIG_CALCULAR_IMPOSTOS_DDV            = 408
	CONFIG_PROCESSO_PADRAO_DDV              = 409
	CONFIG_NFSE_AMBIENTE                    = 410
	CONFIG_CFOP_DENTRO_EST                  = 414
	CONFIG_CFOP_FORA_EST                    = 415
	CONFIG_USAR_CARNE                       = 416
	CONFIG_EXIBIR_ITENS_VENDIDOS_ENTIDADE   = 417
	CONFIG_EXIBIR_VEICULOS_ENTIDADE         = 418
	CONFIG_VALIDAR_VALOR_MINIMO_NA_VENDA    = 419
	CONFIG_VALOR_UNITARIO_VENDA_HABILITADO  = 420
	CONFIG_LIMPAR_TELA_AO_FINALIZAR         = 421
	CONFIG_VALIDAR_INADIMPLENCIA            = 422
	CONFIG_EXIBIR_OPCOES_IMPRESSAO          = 423
	CONFIG_IMPRIMIR_AO_FINALIZAR            = 424
	CONFIG_VALIDAR_DESCONTO                 = 425
	CONFIG_USAR_ADICIONAL_GARCOM            = 426
	CONFIG_PERC_ADICIONAL_GARCOM            = 427
	CONFIG_PRODUTO_ADICIONAL_GARCOM         = 428
	CONFIG_VALOR_TAXA_ENTREGA               = 429
	CONFIG_PRODUTO_TAXA_ENTREGA             = 430
	CONFIG_DDV_TIPO_DOC_PADRAO              = 431
	CONFIG_MARGEM_LUCRO_PADRAO              = 432
	CONFIG_UPDATE_DATA_DDV_AO_FINALIZAR     = 433
	CONFIG_VINCULAR_VEICULO_NO_PEDIDO       = 434
	CONFIG_LIMPAR_PAGTO_CRIAR_PGTO          = 435
	CONFIG_VALIDAR_ESTOQUE_NEGATIVO         = 436
	CONFIG_EXIBIR_CONCLUIDAS                = 437
	CONFIG_EXIBIR_TIT_VENCIDOS              = 438
	CONFIG_EXIBIR_DLG_IMPRESSAO_COMPROVANTE = 439
	CONFIG_DDV_NUMERO_NO_TIT_DOCUMENTO      = 440
	CONFIG_DDV_MOTIVO_CANCELAMENTO          = 441
	CONFIG_ID_TDOC_COMISSAO                 = 442
	CONFIG_USAR_TITULOS_COMISSAO            = 443
	CONFIG_USAR_ANALISE_PRODUTOS            = 444
	CONFIG_SOLICITAR_SENHA_FINALIZAR        = 445
	CONFIG_SOLICITAR_CHECKLIST_INCLUIR      = 446
	CONFIG_SOLICITAR_CHECKLIST_CONCLUIR     = 447
	CONFIG_INFORMAR_PECAS_TERCEIROS         = 448
)

// ============================================
// CONFIGURAÇÕES - FISCAL
// ============================================

const (
	CONFIG_PROCESSO_PADRAO_NTF               = 500
	CONFIG_NATUREZA_PADRAO_NTF               = 501
	CONFIG_TOKEN_NFCE                        = 502
	CONFIG_OPF_ID_NFE_DE_NFCE                = 503
	CONFIG_ID_TOKEN_NFCE                     = 504
	CONFIG_NFSE_TEMPO_ESPERA_RET             = 505
	CONFIG_ATUALIZAR_PRECO_VENDA_NFE_ENTRADA = 506
	CONFIG_EMAIL_PARAO_NFE_NFCE              = 507
	CONFIG_PASTA_XML_NFE                     = 508
	CONFIG_PASTA_XML_NFCE                    = 509
	CONFIG_PASTA_XML_NFSE                    = 510
	CONFIG_IMP_REFERENCIA_NFE                = 511
	CONFIG_TENTATIVA_DUPLICIDADE_NFE         = 512
	CONFIG_VERSAO_NFE                        = 513
	CONFIG_PRC_ID_TRANSF_FILIAL              = 514
	CONFIG_ALIQUOTA_CREDITO_ICMS             = 515
	CONFIG_CNPJ_CONTADOR                     = 516
	CONFIG_CT_TIPOAMBIENTE                   = 517
	CONFIG_MDF_TIPOAMBIENTE                  = 518
	CONFIG_MDF_LAYOUT                        = 519
)

// ============================================
// CONFIGURAÇÕES - CONTÁBIL
// ============================================

const (
	CONFIG_PCC_MASCARA          = 600
	CONFIG_USAR_ROTINA_CONTABIL = 601
)

// ============================================
// CONFIGURAÇÕES - PRODUÇÃO
// ============================================

const (
	CONFIG_RMA_QTD_MAIOR_ORP = 700
)

// ============================================
// CONFIGURAÇÕES - RELATÓRIOS
// ============================================

const (
	CONFIG_REL_IMP_DDV                = 900
	CONFIG_REL_IMP_ETIQUETA           = 901
	CONFIG_REL_IMP_CAR                = 902
	CONFIG_REL_IMP_CAR_ENTREGA        = 903
	CONFIG_REL_IMP_DDV_ENTREGA        = 904
	CONFIG_REL_IMP_COMANDA            = 905
	CONFIG_REL_IMP_DDV_ENTREGA_FUTURA = 906
	CONFIG_REL_IMP_NOTA_PROMISSORIA   = 907
	CONFIG_REL_IMP_RECIBO_BAIXA_REC   = 908
	CONFIG_REL_IMP_RECIBO_BAIXA_PAG   = 911
	CONFIG_REL_IMP_FATURA_REC         = 909
	CONFIG_FORMATO_ETQ                = 910
)

// ============================================
// CONFIGURAÇÕES - FRENTE DE CAIXA
// ============================================

const (
	CONFIG_PRC_ID_PADRAO_PDV            = 1000
	CONFIG_ENT_ID_PADRAO_PDV            = 1001
	CONFIG_CODPGT_ID_PADRAO_PDV         = 1002
	CONFIG_FLU_ID_PADRAO_PDV            = 1003
	CONFIG_TBP_ID_PADRAO_PDV            = 1004
	CONFIG_LOCALIZACAO_PADRAO_PDV       = 1005
	CONFIG_EXIBIR_DLG_IMPRESSORA        = 1006
	CONFIG_SOLICITAR_NUM_COPIAS         = 1007
	CONFIG_SOLICITAR_ENVIO_NFCE         = 1008
	CONFIG_SEMPRE_ENVIAR_NFCE           = 1009
	CONFIG_BALANCA_PESO                 = 1010
	CONFIG_BALANCA_VALOR                = 1011
	CONFIG_TAM_FONTE_ITEM               = 1012
	CONFIG_REL_IMP_NOTA_DEBITO          = 1013
	CONFIG_OPCAO_PADRAO_EMISSAO_NF      = 1014
	CONFIG_CODBARRA_DIVERSOS            = 1015
	CONFIG_USAR_ADICIONAL               = 1016
	CONFIG_PRO_ID_ADICIONAIS            = 1017
	CONFIG_BALANCA_TAM_CODIGO           = 1018
	CONFIG_NUM_MESAS                    = 1019
	CONFIG_USAR_DESC_ITEM               = 1020
	CONFIG_NUM_COPIA_IMPRESSAO          = 1021
	CONFIG_ASSOCIA_MESA_COMANDA         = 1022
	CONFIG_IMPRIME_COMANDA_COZINHA      = 1023
	CONFIG_TAMANHO_BOTOES               = 1024
	CONFIG_TAMANHO_FONTE_BOTOES         = 1025
	CONFIG_TEXTO_BOTOES                 = 1026
	CONFIG_FECHAMENTO_AUTO_VENDA_TOUCH  = 1027
	CONFIG_USAR_COMANDA_COZINHA         = 1028
	CONFIG_PERMITIR_CREDIARIO_ENTREGA   = 1029
	CONFIG_PRC_ID_DEV_VENDA             = 1030
	CONFIG_PRC_ID_VENDA_CONDICIONAL     = 1031
	CONFIG_PRC_ID_DEV_VENDA_CONDICIONAL = 1032
	CONFIG_VISUALIZAR_NFCE              = 1033
)

// ============================================
// CONFIGURAÇÕES - MOBILE GARÇOM
// ============================================

const (
	CONFIG_TABELA_PRECO_GARCOM = 1100
)

// ============================================
// CONFIGURAÇÕES - EMAIL
// ============================================

const (
	CONFIG_SMTP_SERVIDOR = 2
	CONFIG_SMTP_USUARIO  = 3
	CONFIG_SMTP_SENHA    = 4
	CONFIG_SMTP_PORTA    = 5
	CONFIG_SMTP_EMAIL    = 6

	CONFIG_SMTP_SERVIDOR1      = 1200
	CONFIG_SMTP_PORTA1         = 1201
	CONFIG_SMTP_TLS            = 1202
	CONFIG_SMTP_SSL            = 1203
	CONFIG_SMTP_USUARIO_NFC    = 1204
	CONFIG_SMTP_EMAIL_NFC      = 1205
	CONFIG_SMTP_SENHA_NFC      = 1206
	CONFIG_SMTP_USUARIO_BOLETO = 1207
	CONFIG_SMTP_EMAIL_BOLETO   = 1208
	CONFIG_SMTP_SENHA_BOLETO   = 1209
)

// ============================================
// CONFIGURAÇÕES - INTEGRAÇÃO PDV
// ============================================

const (
	CONFIG_PDV_PASTA_INTEGRADOR_EXP = 1400
	CONFIG_PDV_PASTA_INTEGRADOR_IMP = 1401
	CONFIG_UTILIZAR_INTEGRADOR      = 1402
	CONFIG_VERSAO_DJPDV             = 1403
	CONFIG_UTILIZA_COMANDA          = 1404
)

// ============================================
// CONFIGURAÇÕES - TRIBUTAÇÕES
// ============================================

const (
	CONFIG_CST_ICMS_PADRAO            = 411
	CONFIG_CST_IPI_PADRAO             = 412
	CONFIG_CST_PISCOFINS_SAIDA        = 413
	CONFIG_CST_PISCOFINS_ENTRADA      = 1500
	CONFIG_CST_PISCOFINS_ENTRADA_MONO = 1501
)

// ============================================
// CONFIGURAÇÕES - RELATÓRIOS
// ============================================

const (
	CONFIG_REL_DDV_ESPELHO       = 1600
	CONFIG_REL_RECIBO_ADIANT_ENT = 1601
)

// ============================================
// CONFIGURAÇÕES - TOTVS
// ============================================

const (
	CONFIG_TOTVS_PROD_URL   = 1700
	CONFIG_TOTVS_PROD_PORTA = 1701
	CONFIG_TOTVS_PROD_USER  = 1702
	CONFIG_TOTVS_PROD_SENHA = 1703
	CONFIG_TOTVS_HOMO_URL   = 1704
	CONFIG_TOTVS_HOMO_PORTA = 1705
	CONFIG_TOTVS_HOMO_USER  = 1706
	CONFIG_TOTVS_HOMO_SENHA = 1707
	CONFIG_TOTVS_PRODUCAO   = 1708
)

// ============================================
// CONFIGURAÇÕES - DOMINIO
// ============================================

const (
	CONFIG_DOMINIO_USUARIO              = 1800
	CONFIG_DOMINIO_SENHA                = 1801
	CONFIG_DOMINIO_URL_TOKEN            = 1802
	CONFIG_DOMINIO_URL_INFO             = 1803
	CONFIG_DOMINIO_URL_ENABLE           = 1804
	CONFIG_DOMINIO_URL_BATCHES          = 1805
	CONFIG_DOMINIO_AUDIENCE             = 1806
	CONFIG_DOMINIO_COOKIE               = 1807
	CONFIG_DOMINIO_X_INTEGRATIO_KEY     = 1808
	CONFIG_DOMINIO_INTEGRATIO_KEY       = 1809
	CONFIG_DOMINIO_TOKEN                = 1810
	CONFIG_DOMINIO_TOKEN_VALIDADE       = 1811
	CONFIG_DOMINIO_USAR_INTEGRACAO      = 1812
	CONFIG_DOMINIO_TEMPO_SLEEP_CONSULTA = 1813
)

// ============================================
// CONFIGURAÇÕES - LINX (D-TEF)
// ============================================

const (
	CONFIG_LINX_USAR_TEF        = 1900
	CONFIG_LINX_URLCERTIFICADO  = 1901
	CONFIG_LINX_PATHCERTIFICADO = 1902
	CONFIG_LINX_IPPORTASSL      = 1903
)

// ============================================
// CONFIGURAÇÕES - IFOOD
// ============================================

const (
	CONFIG_IFOOD_INTEGRARIFOOD = 2000
	CONFIG_IFOOD_CLIENTID      = 2001
	CONFIG_IFOOD_CLIENTSECRET  = 2002
	CONFIG_IFOOD_TOKEN         = 2003
	CONFIG_IFOOD_VALIDADETOKEN = 2004
	CONFIG_IFOOD_TIPOTOKEN     = 2005
	CONFIG_IFOOD_MERCHANTID    = 2006
	CONFIG_IFOOD_REP_ID        = 2007
)

// ============================================
// HABILITAÇÃO
// ============================================

const (
	HAB_CODIGO_CLIENTE            = 1
	HAB_CODIGO_FORNECEDOR         = 2
	HAB_CODIGO_TRANSPORTADORA     = 3
	HAB_CODIGO_REPRESENTANTE      = 4
	HAB_CODIGO_ASSISTENTE_TECNICO = 5
	HAB_CODIGO_MOTORISTA          = 7
)

// ============================================
// OPERAÇÕES DE PROCESSOS
// ============================================

const (
	OPE_ID_MOVIMENTAESTOQUE = 1
	OPE_ID_GERAFINANCEIRO   = 2
	OPE_ID_LANCACONTABIL    = 3
	OPE_ID_ATUALIZARCUSTO   = 4
)

// ============================================
// ADIANTAMENTO DE ENTIDADE
// ============================================

const (
	ADTENT_TIPOADIANTAMENTO_CREDITO = 1
	ADTENT_TIPOADIANTAMENTO_DEBITO  = -1
)

// ============================================
// ENTIDADE
// ============================================

const (
	ENT_TIPOPESSOA_FISICA   = 1
	ENT_TIPOPESSOA_JURIDICA = 2
)

// ============================================
// CONVÊNIO BANCÁRIO
// ============================================

const (
	CVB_RESPONSAVELEMISSAO_CLIEMITE        = 0
	CVB_RESPONSAVELEMISSAO_BANCOEMITE      = 1
	CVB_RESPONSAVELEMISSAO_BANCOREEMITE    = 2
	CVB_RESPONSAVELEMISSAO_BANCONAOREEMITE = 3

	CVB_LAYOUT_400 = 0
	CVB_LAYOUT_240 = 1
)

// ============================================
// REMESSA BANCÁRIA
// ============================================

const (
	REMBAN_SITUACAO_ABERTA    = 1
	REMBAN_SITUACAO_ENVIADA   = 2
	REMBAN_SITUACAO_CANCELADA = 9
)

// ============================================
// RETORNO BANCÁRIO
// ============================================

const (
	RETBAN_SITUACAO_ABERTO     = 1
	RETBAN_SITUACAO_IMPORTADO  = 2
	RETBAN_SITUACAO_PROCESSADO = 3
	RETBAN_SITUACAO_CANCELADO  = 4

	RETTIT_DEBITOCREDITO_SEM_LANC = 0
	RETTIT_DEBITOCREDITO_DEBITO   = 1
	RETTIT_DEBITOCREDITO_CREDITO  = 2

	RETTIT_OCORRENCIA_LIQUIDACAO     = 0
	RETTIT_OCORRENCIA_PROTESTO       = 1
	RETTIT_OCORRENCIA_INSTRUCAO      = 2
	RETTIT_OCORRENCIA_REG_RECUSADO   = 3
	RETTIT_OCORRENCIA_BAIXA_RECUSADA = 4
	RETTIT_OCORRENCIA_BAIXADO        = 4
	RETTIT_OCORRENCIA_REG_CONFIRMADO = 5
	RETTIT_OCORRENCIA_CANCELAMENTO   = 9
)

// ============================================
// CHEQUES
// ============================================

const (
	CHQREC_SITUACAO_EMCARTEIRA   = 1
	CHQREC_SITUACAO_COMOCORENCIA = 2
	CHQREC_SITUACAO_REJEITADO    = 3
	CHQREC_SITUACAO_COMPENSADO   = 4
	CHQREC_SITUACAO_COBRANCA     = 5
	CHQREC_SITUACAO_REPASSADO    = 6

	OCH_TIPOOCORENCIA_COMPENSADA = 1
	OCH_TIPOOCORENCIA_DEVOLVIDA  = 2
)

// ============================================
// PORTADORES
// ============================================

const (
	POR_TIPO_BANCO_COM_REGISTRO = 1
	POR_TIPO_BANCO_SEM_REGISTRO = 2
	POR_TIPO_OUTROS             = 9
)

// ============================================
// EVENTOS
// ============================================

const (
	EVE_NUMERO_TOTALITEM      = 1
	EVE_NUMERO_VALOR_DESCONTO = 2
	EVE_NUMERO_VALOR_JURO     = 3
	EVE_NUMERO_BC_ICMS        = 11
	EVE_NUMERO_VALOR_ICMS     = 13
	EVE_NUMERO_BC_IPI         = 21
	EVE_NUMERO_VALOR_IPI      = 23
	EVE_NUMERO_BC_PIS         = 31
	EVE_NUMERO_VALOR_PIS      = 33
	EVE_NUMERO_BC_COFINS      = 41
	EVE_NUMERO_VALOR_COFINS   = 43
	EVE_NUMERO_BC_ICMS_ST     = 51
	EVE_NUMERO_VALOR_ICMS_ST  = 53
)

// ============================================
// TABELA DE PREÇO PRODUTO
// ============================================

const (
	TBPP_SITUACAO_BLOQUEADO    = 0
	TBPP_SITUACAO_DESBLOQUEADO = 1
)

// ============================================
// APURAÇÃO TRIBUTÁRIA
// ============================================

const (
	ETR_APURACAO_NENHUM         = 0
	ETR_APURACAO_LUCROPRESUMIDO = 1
	ETR_APURACAO_LUCROREAL      = 2

	ETR_REGIME_SIMPLES         = 1
	ETR_REGIME_SIMPLES_EXCESSO = 2
	ETR_REGIME_NORMAL          = 3

	ETR_REGIMEESPECIAL_NENHUM                 = 0
	ETR_REGIMEESPECIAL_MICRO_EMP_MUNIC        = 1
	ETR_REGIMEESPECIAL_ESTIMATIVA             = 2
	ETR_REGIMEESPECIAL_SOCIEDADE_PROFISIONAIS = 3
	ETR_REGIMEESPECIAL_COOPERATIVA            = 4
	ETR_REGIMEESPECIAL_MICRO_EMP_INDIVIDUAL   = 5
	ETR_REGIMEESPECIAL_MICRO_EMP_PP           = 6
)

// ============================================
// RENEGOCIAÇÃO DE TITULOS
// ============================================

const (
	RENG_SITUACAO_ABERTA     = 1
	RENG_SITUACAO_PROCESSADA = 2

	RENG_TIPO_RECEBIRMENTO = 2
	RENG_TIPO_PAGAMENTO    = 1

	RENGDOC_ORIGEM_TITULO = 5
	RENGDOC_ORIGEM_CHEQUE = 7
)

// ============================================
// ASSISTÊNCIA TÉCNICA AVALIAÇÃO
// ============================================

const (
	ATA_TIPO_FALTA               = 1
	ATA_TIPO_TROCA_COM_DEVOLUCAO = 2
	ATA_TIPO_TROCA_SEM_DEVOLUCAO = 3
	ATA_TIPO_OUTROS              = 9
)

// ============================================
// FORMA DE CONTATO
// ============================================

const (
	FRC_TIPO_EMAIL    = 1
	FRC_TIPO_TELEFONE = 2
	FRC_TIPO_OUTROS   = 9
)

// ============================================
// ORDEM DE PRODUÇÃO
// ============================================

const (
	ORP_TIPO_NORMAL              = 1
	ORP_TIPO_RETRABALHO          = 2
	ORP_TIPO_SERVICOS            = 3
	ORP_TIPO_ASSISTENCIA_TECNICA = 4

	ORP_SITUACAO_ABERTA      = 1
	ORP_SITUACAO_CONFIRMADA  = 2
	ORP_SITUACAO_EM_PRODUCAO = 3
	ORP_SITUACAO_FINALIZADA  = 8
	ORP_SITUACAO_CANCELADA   = 9

	ORP_REALIZACAO_INTERNA = 1
	ORP_REALIZACAO_EXTERNO = 2
)

// ============================================
// REQUISIÇÃO
// ============================================

const (
	REQ_TIPOOPERACAO_ENTRADA = 0
	REQ_TIPOOPERACAO_SAIDA   = 1

	REQ_SITUACAO_ABRTA   = 1
	REQ_SITUACAO_FECHADA = 2
)

// ============================================
// ORDEM DE SERVIÇO
// ============================================

const (
	OS_SITUACAO_ABERTA      = 0
	OS_SITUACAO_EMANDAMENTO = 1
	OS_SITUACAO_CONCLUIDA   = 2
	OS_SITUACAO_CANCELADA   = 9
)

// ============================================
// SEGURANÇA
// ============================================

const (
	SEGRG_NIVEL_NENHUMA = 0
	SEGRG_NIVEL_INSERIR = 1
	SEGRG_NIVEL_ALTERAR = 2
	SEGRG_NIVEL_EXCLUIR = 4
)

// ============================================
// LANCAMENTO FINANCEIRO
// ============================================

const (
	LFINT_SITUACAO_ABERTO  = 1
	LFINT_SITUACAO_FECHADO = 2

	LFIN_TIPOLANCAMENTO_CREDITO = 1
	LFIN_TIPOLANCAMENTO_DEBITO  = -1
)

// ============================================
// CCNFE
// ============================================

const (
	CCNFE_CODIGOSTATUS_REG_VINCULADO = 135
)

// ============================================
// GERAÇÃO DE FATURAS DE CONTRATOS
// ============================================

const (
	GCE_SITUACAO_ABERTA    = 0
	GCE_SITUACAO_GERADA    = 1
	GCE_SITUACAO_EFETIVADA = 2
	GCE_SITUACAO_CANCELADA = 9

	CENT_TIPO_RECEBIRMENTO   = 2
	CENT_TIPO_PAGAMENTO      = 1
	CENT_FORMATO_MENSAL      = 1
	CENT_FORMATO_PERIODOFIXO = 2
)

// ============================================
// DESPESAS
// ============================================

const (
	DESP_FIXA_VARIAVEL_FIIXA    = 1
	DESP_FIXA_VARIAVEL_VARIAVEL = 2
)

// ============================================
// CAIXA FINANCEIRO
// ============================================

const (
	CXFIN_SITUACAO_ABERTO  = 1
	CXFIN_SITUACAO_FECHADO = 2
)

// ============================================
// MODELO DE DOCUMENTO FISCAL
// ============================================

const (
	MDF_CERTIFICADO_MODELO_A1 = "1"
	MDF_CERTIFICADO_MODELO_A3 = "3"
)

// ============================================
// LDESP / LREC
// ============================================

const (
	LDESP_SITUACAO_PREVISTA  = 0
	LDESP_SITUACAO_REALIZADA = 1

	LREC_SITUACAO_PREVISTA  = 0
	LREC_SITUACAO_REALIZADA = 1
)

// ============================================
// PRODUTO_ANALISE
// ============================================

const (
	PROANL_TIPO_PERCENTUAL = 1
	PROANL_TIPO_VALOR      = 2
)

// ============================================
// MENSAGENS DE ALERTA
// ============================================

const (
	MSGALT_SITUACAO_NAO_LIDA = 1
	MSGALT_SITUACAO_LIDA     = 2

	MSGALT_ENVIAR     = 0
	MSGALT_ENVIADO    = 1
	MSGALT_NAO_ENVIAR = 2
)

// ============================================
// TBP
// ============================================

const (
	TBP_TIPO_PERC_ACIMA_VALOR_VENDA  = 1
	TBP_TIPO_PERC_ABAIXO_VALOR_VENBA = 2
	TBP_TIPO_PERC_ACIMA_VALOR_CUSTO  = 3
)

// ============================================
// PYTHON SCRIPT IFOOD
// ============================================

const (
	PY_ID_IFOOD_GERARTOKEN = 1
	PY_ID_IFOOD_MERCHANTS  = 2
)

// ============================================
// TIPOS (ENUMS)
// ============================================

// TmOrigem - Origem das operações
type TmOrigem int

const (
	OrgPedidoDeCompra         TmOrigem = 1
	OrgNotaFiscal             TmOrigem = 2
	OrgPedidoDeVenda          TmOrigem = 3
	OrgGerados                TmOrigem = 4
	OrgTitulos                TmOrigem = 5
	OrgTituloBaixa            TmOrigem = 6
	OrgChequesRecebidos       TmOrigem = 7
	OrgChequeOcorrencia       TmOrigem = 8
	OrgAdiantamentoEnt        TmOrigem = 9
	OrgRenegociacao           TmOrigem = 10
	OrgRequisicaoAlmoxarifado TmOrigem = 11
	OrgContratoEntidade       TmOrigem = 12
	OrgOrdemCompra            TmOrigem = 13
)

// TmUsoDestinatario - Uso do destinatário
type TmUsoDestinatario int

const (
	TudRevenda TmUsoDestinatario = 1
	TudConsumo TmUsoDestinatario = 2
)

// TmTipoOperacao - Tipo de operação
type TmTipoOperacao int

const (
	TopEntrada TmTipoOperacao = 0
	TopSaida   TmTipoOperacao = 1
)

// TDentroForaEstado - Dentro ou fora do estado
type TDentroForaEstado int

const (
	TdfeDentroDoEstado TDentroForaEstado = 0
	TdfeForaEstado     TDentroForaEstado = 1
)

// ============================================
// STRUCTS (RECORDS)
// ============================================

// TRecProdTribCFOP - CFOP por tipo de operação
type TRecProdTribCFOP struct {
	OpfIdNoEst     int
	OpfIdNoEstST   int
	OpfIdForaEst   int
	OpfIdForaEstST int
}

// TRecProdTribCST - CST por tipo de imposto
type TRecProdTribCST struct {
	CstIcmsId       int
	CstIpiId        int
	CstPisConfinsId int
}

// ============================================
// FUNÇÕES AUXILIARES (CONVERSÃO)
// ============================================

// OrigemParaInteger - Converte TmOrigem para int
func OrigemParaInteger(origem TmOrigem) int {
	return int(origem)
}

// UsoDestinatarioParaInteger - Converte TmUsoDestinatario para int
func UsoDestinatarioParaInteger(uso TmUsoDestinatario) int {
	return int(uso)
}

// TipoOperacaoParaInteger - Converte TmTipoOperacao para int
func TipoOperacaoParaInteger(op TmTipoOperacao) int {
	return int(op)
}

// FrmPgtoParaNFeFormaPagamento - Converte forma de pagamento para o formato NFe
// Retorna um inteiro que representa a forma de pagamento no padrão NFe
func FrmPgtoParaNFeFormaPagamento(frmPgto int) int {
	switch frmPgto {
	case FRMPGTO_TIPO_CHEQUE:
		return 1 // fpCheque
	case FRMPGTO_TIPO_CREDIARIO:
		return 2 // fpCreditoLoja
	case FRMPGTO_TIPO_CARTAODEBITO:
		return 3 // fpCartaoDebito
	case FRMPGTO_TIPO_CARTAOCREDITO:
		return 4 // fpCartaoCredito
	case FRMPGTO_TIPO_PIX:
		return 6 // fpPagamentoInstantaneoEstatico
	case FRMPGTO_TIPO_CONTRAVALE:
		return 8 // fpValePresente
	case FRMPGTO_TIPO_ENTREGA, FRMPGTO_TIPO_ORCAMENTO, FRMPGTO_TIPO_DEVOLUCAO:
		return 10 // fpOutro
	case FRMPGTO_TIPO_OUTRAS, FRMPGTO_TIPO_DINHEIRO:
		return 0 // fpDinheiro
	case FRMPGTO_TIPO_VOUCHER:
		return 7 // fpValeRefeicao
	default:
		return 0 // fpDinheiro
	}
}

// IsTipoTEF - Verifica se a forma de pagamento é TEF
func IsTipoTEF(frmPgto int) bool {
	for _, v := range FRMPGTO_TIPO_TEF {
		if v == frmPgto {
			return true
		}
	}
	return false
}

// ============================================
// FUNÇÃO DE OBSERVAÇÃO (TObservacoes)
// ============================================

// GetObservacaoFunRural - Retorna a observação de FunRural
func GetObservacaoFunRural(aliquota, valor float64) string {
	return fmt.Sprintf("Documento com retenção de FunRural, Aliquota: %.2f valor: %.2f", aliquota, valor)
}

// ============================================
// CONSTANTES ADICIONAIS (NÃO ESTAVAM NO ORIGINAL, MAS ÚTEIS)
// ============================================

// StringToInt - Converte string para int com fallback
func StringToInt(s string, defaultVal int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return defaultVal
}

// IntToString - Converte int para string
func IntToString(v int) string {
	return strconv.Itoa(v)
}
