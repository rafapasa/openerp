package models

import (
	"time"
)

// ============================================================
// MODEL: NotaFiscal
// ============================================================

type NotaFiscal struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID                          int    `gorm:"column:ntf_id;primaryKey;autoIncrement" json:"id"`
	NaturezaOperacaoID          int    `gorm:"column:nfno_id;not null" json:"natureza_operacao_id"`
	EntidadeID                  int    `gorm:"column:ent_id;not null" json:"entidade_id"`
	EntidadeTransportadoraID    *int   `gorm:"column:ent_id_transp" json:"entidade_transportadora_id,omitempty"`
	VeiculoTransportadorID      *int   `gorm:"column:vei_id_trsnp" json:"veiculo_transportador_id,omitempty"`
	SituacaoID                  int    `gorm:"column:nfsit_id;not null" json:"situacao_id"`
	MunicipioID                 *int   `gorm:"column:mun_id" json:"municipio_id,omitempty"`
	PaisID                      *int   `gorm:"column:pai_id" json:"pais_id,omitempty"`
	EstadoID                    *int   `gorm:"column:est_id" json:"estado_id,omitempty"`
	EmpresaFilialID             int    `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	ModeloDocumentoFiscalCodigo string `gorm:"column:mdf_codigo;type:varchar(3);not null" json:"modelo_documento_fiscal_codigo"`
	ProcessoID                  int    `gorm:"column:prc_id;not null" json:"processo_id"`
	CondicaoPagamentoID         *int   `gorm:"column:codpgt_id" json:"condicao_pagamento_id,omitempty"`

	ChaveAcesso        *string    `gorm:"column:ntf_chaveacesso;type:varchar(47)" json:"chave_acesso,omitempty"`
	IndicadorPagamento int        `gorm:"column:ntf_indpgto;not null" json:"indicador_pagamento"`
	Serie              string     `gorm:"column:ntf_serie;type:varchar(3);not null" json:"serie"`
	Numero             int        `gorm:"column:ntf_numero;not null" json:"numero"`
	DataEmissao        time.Time  `gorm:"column:ntf_dataemissao;type:date;not null" json:"data_emissao"`
	DataEntradaSaida   *time.Time `gorm:"column:ntf_dataentradasaida;type:date" json:"data_entrada_saida,omitempty"`
	HoraEntradaSaida   *string    `gorm:"column:ntf_horaentradasaida;type:time" json:"hora_entrada_saida,omitempty"`

	TipoImpressao        int        `gorm:"column:ntf_tipoimpressao;not null" json:"tipo_impressao"`
	TipoOperacao         int        `gorm:"column:ntf_tipooperacao;not null" json:"tipo_operacao"` // 0 - entrada, 1 - saida
	TipoEmitente         int        `gorm:"column:ntf_tipoemitente;not null" json:"tipo_emitente"` // 0 - própria, 1 - terceiros
	TipoEmissao          int        `gorm:"column:ntf_tipoemissao;not null" json:"tipo_emissao"`
	DigitoChave          *int       `gorm:"column:ntf_digitochave" json:"digito_chave,omitempty"`
	TipoAmbiente         int        `gorm:"column:ntf_tipoambiente;not null" json:"tipo_ambiente"`
	Finalidade           int        `gorm:"column:ntf_finalidade;not null" json:"finalidade"`
	ProcessoEmissao      int        `gorm:"column:ntf_processoemissao;not null" json:"processo_emissao"`
	VersaoApp            string     `gorm:"column:ntf_versaoapp;type:varchar(20);not null" json:"versao_app"`
	DataContingencia     *time.Time `gorm:"column:ntf_datacontingencia;type:date" json:"data_contingencia,omitempty"`
	JustificativaConting *string    `gorm:"column:ntf_junstificaticacont;type:varchar(255)" json:"justificativa_contingencia,omitempty"`

	ValorDesconto   *float64 `gorm:"column:ntf_valordesconto;type:decimal(15,2)" json:"valor_desconto,omitempty"`
	ValorAbatimento *float64 `gorm:"column:ntf_valorabat_nt;type:decimal(15,2)" json:"valor_abatimento,omitempty"`
	TipoFrete       int      `gorm:"column:ntf_tipofrete;not null" json:"tipo_frete"`
	ValorFrete      *float64 `gorm:"column:ntf_valorfrete;type:decimal(15,2)" json:"valor_frete,omitempty"`
	ValorSeguro     *float64 `gorm:"column:ntf_valorseguro;type:decimal(15,2)" json:"valor_seguro,omitempty"`
	ValorOutras     *float64 `gorm:"column:ntf_valoroutras;type:decimal(15,2)" json:"valor_outras,omitempty"`

	UsoDestinatario     int     `gorm:"column:ntf_usodestinatario;not null" json:"uso_destinatario"` // 1 - revenda, 2 - consumo
	EnderecoLogradouro  *string `gorm:"column:ntf_end_logradouro;type:varchar(100)" json:"endereco_logradouro,omitempty"`
	EnderecoNumero      *string `gorm:"column:ntf_end_numero;type:varchar(20)" json:"endereco_numero,omitempty"`
	EnderecoComplemento *string `gorm:"column:ntf_end_compl;type:varchar(255)" json:"endereco_complemento,omitempty"`
	EnderecoBairro      *string `gorm:"column:ntf_end_bairro;type:varchar(255)" json:"endereco_bairro,omitempty"`
	EnderecoCEP         *int    `gorm:"column:ntf_end_cep" json:"endereco_cep,omitempty"`

	XML *string `gorm:"column:ntf_xml;type:longtext" json:"xml,omitempty"`

	ValorProdutos float64 `gorm:"column:ntf_valorprodutos;type:decimal(15,2);not null" json:"valor_produtos"`
	ValorLiquido  float64 `gorm:"column:ntf_valorliquido;type:decimal(15,2);not null" json:"valor_liquido"`

	ICMSBase    *float64 `gorm:"column:ntf_icms_base;type:decimal(15,2)" json:"icms_base,omitempty"`
	ICMSValor   *float64 `gorm:"column:ntf_icms_valor;type:decimal(15,2)" json:"icms_valor,omitempty"`
	IPIBase     *float64 `gorm:"column:ntf_ipi_base;type:decimal(15,2)" json:"ipi_base,omitempty"`
	IPIValor    *float64 `gorm:"column:ntf_ipi_valor;type:decimal(15,2)" json:"ipi_valor,omitempty"`
	PISBase     *float64 `gorm:"column:ntf_pis_base;type:decimal(15,2)" json:"pis_base,omitempty"`
	PISValor    *float64 `gorm:"column:ntf_pis_valor;type:decimal(15,2)" json:"pis_valor,omitempty"`
	COFINSBase  *float64 `gorm:"column:ntf_cofins_base;type:decimal(15,2)" json:"cofins_base,omitempty"`
	COFINSValor *float64 `gorm:"column:ntf_cofins_valor;type:decimal(15,2)" json:"cofins_valor,omitempty"`
	ICMSSTBase  *float64 `gorm:"column:ntf_icmsst_base;type:decimal(15,2)" json:"icms_st_base,omitempty"`
	ICMSSTValor *float64 `gorm:"column:ntf_icmsst_valor;type:decimal(15,2)" json:"icms_st_valor,omitempty"`

	Observacao *string `gorm:"column:ntf_observacao;type:text" json:"observacao,omitempty"`

	TransportePlaca *string `gorm:"column:ntf_transp_placa;type:varchar(10)" json:"transporte_placa,omitempty"`
	TransporteUF    *string `gorm:"column:ntf_transp_uf;type:varchar(2)" json:"transporte_uf,omitempty"`
	TransporteRNTC  *string `gorm:"column:ntf_transp_rntc;type:varchar(255)" json:"transporte_rntc,omitempty"`

	DataHoraCancelamento  *time.Time `gorm:"column:ntf_datahoracancelamento;type:datetime" json:"data_hora_cancelamento,omitempty"`
	ProtocoloCancelamento *string    `gorm:"column:ntf_protocolocancelamento;type:varchar(255)" json:"protocolo_cancelamento,omitempty"`
	XMLCancelamento       *string    `gorm:"column:ntf_xmlcancelamento;type:text" json:"xml_cancelamento,omitempty"`
	StatusCancelamento    *int       `gorm:"column:ntf_statuscancelamento" json:"status_cancelamento,omitempty"`
	MotivoCancelamento    *string    `gorm:"column:ntf_motivocancelamento;type:text" json:"motivo_cancelamento,omitempty"`

	NFSeProtocolo         *string `gorm:"column:ntf_nfse_protocolo;type:varchar(255)" json:"nfse_protocolo,omitempty"`
	NFSeCodigoVerificacao *string `gorm:"column:ntf_nfse_codigoverificacao;type:varchar(255)" json:"nfse_codigo_verificacao,omitempty"`
	NFSeNumero            *string `gorm:"column:ntf_nfse_numero;type:varchar(100)" json:"nfse_numero,omitempty"`

	ValorPago  *float64 `gorm:"column:ntf_valorpago;type:decimal(15,2);default:0" json:"valor_pago,omitempty"`
	ValorTroco *float64 `gorm:"column:ntf_valortroco;type:decimal(15,2);default:0" json:"valor_troco,omitempty"`

	TotalPesoBruto   *float64 `gorm:"column:ntf_totalpesobruto;type:decimal(15,4)" json:"total_peso_bruto,omitempty"`
	TotalPesoLiquido *float64 `gorm:"column:ntf_totalpesoliquido;type:decimal(15,4)" json:"total_peso_liquido,omitempty"`

	TransporteNome             *string `gorm:"column:ntf_transp_nome;type:varchar(255)" json:"transporte_nome,omitempty"`
	TransporteInscricaoFederal *string `gorm:"column:ntf_trans_inscricaofederal;type:varchar(20)" json:"transporte_inscricao_federal,omitempty"`
	EntidadeInscricaoFederal   *string `gorm:"column:ntf_ent_inscricaofederal;type:varchar(20)" json:"entidade_inscricao_federal,omitempty"`
	EntidadeInscricaoEstadual  *string `gorm:"column:ntf_ent_inscricaoestadual;type:varchar(20)" json:"entidade_inscricao_estadual,omitempty"`

	IndicadorPresenca *int `gorm:"column:ntf_indpresenca" json:"indicador_presenca,omitempty"` // 0=presencial, 2=internet, etc
	ProdutoMarcaID    *int `gorm:"column:promar_id" json:"produto_marca_id,omitempty"`
	ProdutoEspecieID  *int `gorm:"column:proesp_id" json:"produto_especie_id,omitempty"`

	QuantidadeVolume *float64 `gorm:"column:ntf_qunatvolume;type:decimal(15,4)" json:"quantidade_volume,omitempty"`

	// Totvs Integration
	TotvsNumContrato *string `gorm:"column:ntf_totvs_numcontrato;type:varchar(20)" json:"totvs_num_contrato,omitempty"`
	TotvsSafra       *int    `gorm:"column:ntf_totvs_safra" json:"totvs_safra,omitempty"`
	TotvsTipoSafra   *string `gorm:"column:ntf_totvs_tiposafra;type:varchar(2)" json:"totvs_tipo_safra,omitempty"`
	TotvsCodColigada *int    `gorm:"column:ntf_totvs_codcoligada" json:"totvs_cod_coligada,omitempty"`
	TotvsIDMov       *string `gorm:"column:ntf_totvs_idmov;type:varchar(20)" json:"totvs_id_mov,omitempty"`
	TotvsResposta    *string `gorm:"column:ntf_totvs_resposta;type:text" json:"totvs_resposta,omitempty"`

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
	NaturezaOperacao       *NotaFiscalNaturezaOperacao `gorm:"foreignKey:NaturezaOperacaoID;references:nfno_id" json:"natureza_operacao,omitempty"`
	Entidade               *Entidade                   `gorm:"foreignKey:EntidadeID;references:ent_id" json:"entidade,omitempty"`
	EntidadeTransportadora *Entidade                   `gorm:"foreignKey:EntidadeTransportadoraID;references:ent_id" json:"entidade_transportadora,omitempty"`
	VeiculoTransportador   *Veiculos                   `gorm:"foreignKey:VeiculoTransportadorID;references:vei_id" json:"veiculo_transportador,omitempty"`
	Situacao               *NotaFiscalSituacao         `gorm:"foreignKey:SituacaoID;references:nfsit_id" json:"situacao,omitempty"`
	Municipio              *Municipio                  `gorm:"foreignKey:MunicipioID;references:mun_id" json:"municipio,omitempty"`
	Pais                   *Pais                       `gorm:"foreignKey:PaisID;references:pai_id" json:"pais,omitempty"`
	Estado                 *Estado                     `gorm:"foreignKey:EstadoID;references:est_id" json:"estado,omitempty"`
	EmpresaFilial          *EmpresaFilial              `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
	ModeloDocumentoFiscal  *ModeloDocumentoFiscalSerie `gorm:"foreignKey:ModeloDocumentoFiscalCodigo,EmpresaFilialID;references:mdf_codigo,emf_id" json:"modelo_documento_fiscal,omitempty"`
	Processo               *Processo                   `gorm:"foreignKey:ProcessoID;references:prc_id" json:"processo,omitempty"`
	CondicaoPagamento      *CondicaoPagamento          `gorm:"foreignKey:CondicaoPagamentoID;references:codpgt_id" json:"condicao_pagamento,omitempty"`
	ProdutoMarca           *ProdutoMarca               `gorm:"foreignKey:ProdutoMarcaID;references:promar_id" json:"produto_marca,omitempty"`
	ProdutoEspecie         *ProdutoEspecie             `gorm:"foreignKey:ProdutoEspecieID;references:proesp_id" json:"produto_especie,omitempty"`

	Itens           []NotaFiscalItem           `gorm:"foreignKey:NotaFiscalID;references:ntf_id" json:"itens,omitempty"`
	Pagamentos      []NotaFiscalPagamento      `gorm:"foreignKey:NotaFiscalID;references:ntf_id" json:"pagamentos,omitempty"`
	DocumentosVenda []NotaFiscalDocumentoVenda `gorm:"foreignKey:NotaFiscalID;references:ntf_id" json:"documentos_venda,omitempty"`
	Historico       []NotaFiscalHistorico      `gorm:"foreignKey:NotaFiscalID;references:ntf_id" json:"historico,omitempty"`
	CartasCorrecao  []CartaCorrecaoNFe         `gorm:"foreignKey:NotaFiscalID;references:ntf_id" json:"cartas_correcao,omitempty"`
	Referenciadas   []NotaFiscalReferenciada   `gorm:"foreignKey:NotaFiscalID;references:ntf_id" json:"referenciadas,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (NotaFiscal) TableName() string {
	return "nota_fiscal"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se a nota fiscal foi deletada logicamente
func (n *NotaFiscal) IsDeleted() bool {
	return n.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (n *NotaFiscal) SoftDelete() {
	now := time.Now()
	n.DeletedAt = &now
}

// IsEntrada verifica se é uma nota fiscal de entrada
func (n *NotaFiscal) IsEntrada() bool {
	return n.TipoOperacao == 0
}

// IsSaida verifica se é uma nota fiscal de saída
func (n *NotaFiscal) IsSaida() bool {
	return n.TipoOperacao == 1
}

// IsCancelada verifica se a nota fiscal foi cancelada
func (n *NotaFiscal) IsCancelada() bool {
	return n.StatusCancelamento != nil && *n.StatusCancelamento == 1
}

// HasChaveAcesso verifica se possui chave de acesso
func (n *NotaFiscal) HasChaveAcesso() bool {
	return n.ChaveAcesso != nil && *n.ChaveAcesso != ""
}

// HasXML verifica se possui XML
func (n *NotaFiscal) HasXML() bool {
	return n.XML != nil && *n.XML != ""
}

// GetValorTotal retorna o valor total da nota fiscal
func (n *NotaFiscal) GetValorTotal() float64 {
	return n.ValorLiquido
}

// GetItensCount retorna a quantidade de itens da nota fiscal
func (n *NotaFiscal) GetItensCount() int {
	return len(n.Itens)
}

// GetPagamentosCount retorna a quantidade de pagamentos da nota fiscal
func (n *NotaFiscal) GetPagamentosCount() int {
	return len(n.Pagamentos)
}

// HasTransportadora verifica se possui transportadora
func (n *NotaFiscal) HasTransportadora() bool {
	return n.EntidadeTransportadoraID != nil && *n.EntidadeTransportadoraID > 0
}

// HasVeiculoTransportador verifica se possui veículo transportador
func (n *NotaFiscal) HasVeiculoTransportador() bool {
	return n.VeiculoTransportadorID != nil && *n.VeiculoTransportadorID > 0
}

// IsNFSe verifica se é uma NFSe
func (n *NotaFiscal) IsNFSe() bool {
	return n.ModeloDocumentoFiscalCodigo == "NFS" || n.ModeloDocumentoFiscalCodigo == "NFS-e"
}

// IsNFe verifica se é uma NFe
func (n *NotaFiscal) IsNFe() bool {
	return n.ModeloDocumentoFiscalCodigo == "NFE" || n.ModeloDocumentoFiscalCodigo == "NFe"
}

// IsCTe verifica se é um CTe
func (n *NotaFiscal) IsCTe() bool {
	return n.ModeloDocumentoFiscalCodigo == "CTE" || n.ModeloDocumentoFiscalCodigo == "CT-e"
}

// GetNumeroCompleto retorna o número completo da nota fiscal (série + número)
func (n *NotaFiscal) GetNumeroCompleto() string {
	return n.Serie + "/" + string(rune(n.Numero))
}
