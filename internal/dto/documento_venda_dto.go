// ============================================================
// FILE: documento_venda_dto.go
// PACKAGE: dto
// ============================================================

package dto

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/openerp/backend/internal/constants"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// REQUESTS
// ============================================================

// DocumentoVendaItemRequest representa um item no request de um documento de venda.
type DocumentoVendaItemRequest struct {
	DocumentoVendaID   int      `json:"id,omitempty"`
	ProdutoID          int      `json:"produto_id" binding:"required"`
	Quantidade         float64  `json:"quantidade" binding:"required,gt=0"`
	ValorUnitario      float64  `json:"valor_unitario" binding:"gte=0"`
	PercentualDesconto *float64 `json:"percentual_desconto,omitempty"`
	ValorDesconto      *float64 `json:"valor_desconto,omitempty"`
	ValorFrete         *float64 `json:"valor_frete,omitempty"`
	DescricaoProduto   *string  `json:"descricao_produto,omitempty"`
	PesoBruto          *float64 `json:"peso_bruto,omitempty"`
	PesoLiquido        *float64 `json:"peso_liquido,omitempty"`
	CreatedBy          *int     `json:"created_by,omitempty"`
	UpdatedBy          *int     `json:"updated_by,omitempty"`
}

// DocumentoVendaPagamentoRequest representa um pagamento no request
type DocumentoVendaPagamentoRequest struct {
	DocumentoVendaID int        `json:"id,omitempty"`
	FormaPagamentoID int        `json:"forma_pagamento_id" binding:"required"`
	Valor            float64    `json:"valor" binding:"required,gt=0"`
	DataVencimento   *time.Time `json:"data_vencimento,omitempty"`
	DataPagamento    *time.Time `json:"data_pagamento,omitempty"`
	ValorJuros       *float64   `json:"valor_juros,omitempty"`
	ValorMulta       *float64   `json:"valor_multa,omitempty"`
	ValorDesconto    *float64   `json:"valor_desconto,omitempty"`
	Observacao       *string    `json:"observacao,omitempty"`
	CreatedBy        *int       `json:"created_by,omitempty"`
	UpdatedBy        *int       `json:"updated_by,omitempty"`
}

// DocumentoVendaRequest representa o request para criar/atualizar um documento de venda.
type DocumentoVendaRequest struct {
	// ============================================================
	// DADOS PRINCIPAIS
	// ============================================================
	EmpresaFilialID     int  `json:"empresa_filial_id" binding:"required"`
	EntidadeID          *int `json:"entidade_id"`
	CondicaoPagamentoID int  `json:"condicao_pagamento_id" binding:"required"`
	TabelaPrecoID       int  `json:"tabela_preco_id" binding:"required"`
	TipoDocumento       int  `json:"tipo_documento" binding:"required"` // 1-Orçamento, 2-Pedido
	TipoOperacao        int  `json:"tipo_operacao" binding:"required"`  // 0-Entrada, 1-Saída
	Numero              int  `json:"numero,omitempty"`

	// ============================================================
	// DATAS
	// ============================================================
	DataValidade *time.Time `json:"data_validade,omitempty"`
	DataPrevisao *time.Time `json:"data_previsao,omitempty"`

	// ============================================================
	// VALORES
	// ============================================================
	ValorFrete    *float64 `json:"valor_frete,omitempty"`
	ValorDesconto *float64 `json:"valor_desconto,omitempty"`

	// ============================================================
	// OBSERVAÇÕES
	// ============================================================
	ObservacoesInterna *string `json:"observacoes_interna,omitempty"`
	ObservacoesADM     *string `json:"observacoes_adm,omitempty"`
	ObservacoesCliente *string `json:"observacoes_cliente,omitempty"`

	// ============================================================
	// RELACIONAMENTOS
	// ============================================================
	VendedorID       *int    `json:"vendedor_id,omitempty"`
	TransportadoraID *int    `json:"transportadora_id,omitempty"`
	TipoEntrega      *string `json:"tipo_entrega,omitempty"` // RETIRADA, ENTREGA, LOCAL

	// ============================================================
	// LISTAS
	// ============================================================
	Itens      []DocumentoVendaItemRequest      `json:"itens" binding:"required,min=1,dive"`
	Pagamentos []DocumentoVendaPagamentoRequest `json:"pagamentos,omitempty"`

	// ============================================================
	// AUDITORIA
	// ============================================================
	CreatedBy *int `json:"created_by,omitempty"`
	UpdatedBy *int `json:"updated_by,omitempty"`
}

// ============================================================
// RESPONSES
// ============================================================

// DocumentoVendaItemResponse representa a resposta de um item de documento.
type DocumentoVendaItemResponse struct {
	DocumentoVendaID   int      `json:"documento_venda_id"`
	Item               int      `json:"item"`
	ProdutoID          int      `json:"produto_id"`
	ProdutoNome        string   `json:"produto_nome,omitempty"`
	ProdutoCodigo      int      `json:"produto_codigo,omitempty"`
	Quantidade         float64  `json:"quantidade"`
	ValorUnitario      float64  `json:"valor_unitario"`
	PercentualDesconto *float64 `json:"percentual_desconto,omitempty"`
	ValorDesconto      *float64 `json:"valor_desconto,omitempty"`
	ValorFrete         *float64 `json:"valor_frete,omitempty"`
	TotalItem          float64  `json:"total_item"`
	TotalItemComFrete  float64  `json:"total_item_com_frete,omitempty"`
	PesoBruto          *float64 `json:"peso_bruto,omitempty"`
	PesoLiquido        *float64 `json:"peso_liquido,omitempty"`
}

// DocumentoVendaPagamentoResponse representa a resposta de um pagamento
type DocumentoVendaPagamentoResponse struct {
	DocumentoVendaID   int      `json:"documento_venda_id"`
	Item               int      `json:"item"`
	FormaPagamentoID   int      `json:"forma_pagamento_id"`
	FormaPagamentoNome string   `json:"forma_pagamento_nome,omitempty"`
	Valor              float64  `json:"valor"`
	DataVencimento     string   `json:"data_vencimento,omitempty"`
	DataPagamento      string   `json:"data_pagamento,omitempty"`
	ValorJuros         *float64 `json:"valor_juros,omitempty"`
	ValorMulta         *float64 `json:"valor_multa,omitempty"`
	ValorDesconto      *float64 `json:"valor_desconto,omitempty"`
	Observacao         *string  `json:"observacao,omitempty"`
	// Status             int      `json:"status"` // Not in model
	// StatusLabel        string   `json:"status_label,omitempty"` // Not in model
}

// DocumentoVendaResponse representa a resposta de um documento de venda.
type DocumentoVendaResponse struct {
	// ============================================================
	// DADOS PRINCIPAIS
	// ============================================================
	ID                         int    `json:"id"`
	EmpresaFilialID            int    `json:"empresa_filial_id"`
	EntidadeID                 *int   `json:"entidade_id,omitempty"`
	EntidadeNome               string `json:"entidade_nome,omitempty"`
	EntidadeDocumento          string `json:"entidade_documento,omitempty"`
	CondicaoPagamentoID        int    `json:"condicao_pagamento_id"`
	CondicaoPagamentoDescricao string `json:"condicao_pagamento_descricao,omitempty"`
	TabelaPrecoID              int    `json:"tabela_preco_id"`
	TabelaPrecoDescricao       string `json:"tabela_preco_descricao,omitempty"`
	Numero                     int    `json:"numero"`
	TipoDocumento              int    `json:"tipo_documento"`
	TipoDocumentoLabel         string `json:"tipo_documento_label"`
	TipoOperacao               int    `json:"tipo_operacao"`
	TipoOperacaoLabel          string `json:"tipo_operacao_label"`
	Situacao                   int    `json:"situacao"`
	SituacaoLabel              string `json:"situacao_label"`

	// ============================================================
	// RELACIONAMENTOS
	// ============================================================
	VendedorID         *int    `json:"vendedor_id,omitempty"`
	VendedorNome       string  `json:"vendedor_nome,omitempty"`
	TransportadoraID   *int    `json:"transportadora_id,omitempty"`
	TransportadoraNome string  `json:"transportadora_nome,omitempty"`
	TipoEntrega        *string `json:"tipo_entrega,omitempty"`

	// ============================================================
	// DATAS
	// ============================================================
	DataDocumento string `json:"data_documento"`
	DataValidade  string `json:"data_validade,omitempty"`
	DataPrevisao  string `json:"data_previsao,omitempty"`

	// ============================================================
	// VALORES
	// ============================================================
	TotalProdutos  float64  `json:"total_produtos"`
	TotalDescontos float64  `json:"total_descontos"`
	TotalFrete     float64  `json:"total_frete"`
	TotalLiquido   float64  `json:"total_liquido"`
	ValorFrete     *float64 `json:"valor_frete,omitempty"`
	ValorDesconto  *float64 `json:"valor_desconto,omitempty"`

	// ============================================================
	// OBSERVAÇÕES
	// ============================================================
	ObservacoesInterna *string `json:"observacoes_interna,omitempty"`
	ObservacoesADM     *string `json:"observacoes_adm,omitempty"`
	ObservacoesCliente *string `json:"observacoes_cliente,omitempty"`

	// ============================================================
	// LISTAS
	// ============================================================
	Itens      []DocumentoVendaItemResponse      `json:"itens,omitempty"`
	Pagamentos []DocumentoVendaPagamentoResponse `json:"pagamentos,omitempty"`

	// ============================================================
	// AUDITORIA
	// ============================================================
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	CreatedBy *int   `json:"created_by,omitempty"`
	UpdatedBy *int   `json:"updated_by,omitempty"`
}

// ============================================================
// RESPONSE LISTA
// ============================================================

// DocumentoVendaListResponse representa a resposta de listagem de documentos de venda
type DocumentoVendaListResponse struct {
	Items      []DocumentoVendaResponse `json:"items"`
	Total      int64                    `json:"total"`
	Page       int                      `json:"page"`
	Limit      int                      `json:"limit"`
	TotalPages int                      `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE CONVERSÃO (FROM MODEL) para DocumentoVendaPagamentoResponse
// ============================================================

// FromModel converte models.DocumentoVendaPagamento para DocumentoVendaPagamentoResponse
func (r *DocumentoVendaPagamentoResponse) FromModel(pagamento *models.DocumentoVendaPagamento) *DocumentoVendaPagamentoResponse {
	if pagamento == nil {
		return nil
	}

	r.DocumentoVendaID = pagamento.DocumentoVendaID
	r.Item = pagamento.Item
	// FormaPagamentoID is a pointer in model, dereference it if not nil
	if pagamento.FormaPagamentoID != nil {
		r.FormaPagamentoID = *pagamento.FormaPagamentoID
	}
	r.Valor = pagamento.Valor

	if !pagamento.DataVencimento.IsZero() {
		r.DataVencimento = utils.FormatDate(pagamento.DataVencimento)
	}
	if pagamento.DataPagamento != nil && !pagamento.DataPagamento.IsZero() {
		r.DataPagamento = utils.FormatDate(*pagamento.DataPagamento)
	}
	r.ValorJuros = pagamento.ValorJuros
	r.ValorMulta = pagamento.ValorMulta
	r.ValorDesconto = pagamento.ValorDesconto
	r.Observacao = pagamento.Observacao

	if pagamento.FormaPagamento != nil {
		r.FormaPagamentoNome = pagamento.FormaPagamento.Descricao
	}
	return r
}

type DocumentoVendaItemListResponse struct {
	Items      []DocumentoVendaItemResponse `json:"items"`
	Total      int64                        `json:"total"`
	Page       int                          `json:"page"`
	Limit      int                          `json:"limit"`
	TotalPages int                          `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o DocumentoVendaRequest
func (r *DocumentoVendaRequest) Validate() error {
	// 1. Validação com go-playground/validator
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}

	// 2. Validar campos obrigatórios usando utils
	if err := utils.ValidateMandatoryFields(r); err != nil {
		return err
	}

	// 3. Validar TipoDocumento usando constantes
	if err := validateTipoDocumentoVenda(r.TipoDocumento); err != nil {
		return err
	}

	// 4. Validar TipoOperacao usando constantes
	if err := validateTipoOperacao(r.TipoOperacao); err != nil {
		return err
	}

	// 5. Validar TipoEntrega (se informado)
	if r.TipoEntrega != nil {
		if err := validateTipoEntrega(*r.TipoEntrega); err != nil {
			return err
		}
	}

	// 6. Validar itens
	for i, item := range r.Itens {
		if err := validateDocumentoVendaItem(&item); err != nil {
			return fmt.Errorf("item %d: %w", i+1, err)
		}
	}

	// 7. Validar pagamentos
	for i, pagamento := range r.Pagamentos {
		if err := validateDocumentoVendaPagamento(&pagamento); err != nil {
			return fmt.Errorf("pagamento %d: %w", i+1, err)
		}
	}

	// 8. Validar datas
	if r.DataValidade != nil && r.DataPrevisao != nil {
		if r.DataValidade.Before(*r.DataPrevisao) {
			return fmt.Errorf("data_validade deve ser posterior a data_previsao")
		}
	}

	return nil
}

// validateDocumentoVendaItem valida um item
func validateDocumentoVendaItem(item *DocumentoVendaItemRequest) error {
	// Validar percentual de desconto (se informado)
	if item.PercentualDesconto != nil && (*item.PercentualDesconto < 0 || *item.PercentualDesconto > 100) {
		return fmt.Errorf("percentual_desconto deve estar entre 0 e 100, recebido: %.2f", *item.PercentualDesconto)
	}

	// Validar valor de desconto (se informado)
	if item.ValorDesconto != nil && *item.ValorDesconto < 0 {
		return fmt.Errorf("valor_desconto não pode ser negativo, recebido: %.2f", *item.ValorDesconto)
	}

	return nil
}

// validateDocumentoVendaPagamento valida um pagamento
func validateDocumentoVendaPagamento(pagamento *DocumentoVendaPagamentoRequest) error {
	// Validar número de parcelas
	if pagamento.NumeroParcelas > 0 && pagamento.ParcelaAtual > pagamento.NumeroParcelas {
		return fmt.Errorf("parcela_atual (%d) não pode ser maior que numero_parcelas (%d)",
			pagamento.ParcelaAtual, pagamento.NumeroParcelas)
	}

	return nil
}

// ============================================================
// FUNÇÕES DE VALIDAÇÃO DE CONSTANTES
// ============================================================

func validateTipoDocumentoVenda(valor int) error {
	switch constants.TipoDocumentoVenda(valor) {
	case constants.TipoDocumentoOrcamento, constants.TipoDocumentoPedido:
		return nil
	default:
		return fmt.Errorf("tipo_documento inválido: %d. Deve ser 1 (Orçamento) ou 2 (Pedido)", valor)
	}
}

func validateTipoOperacao(valor int) error {
	switch constants.TipoOperacao(valor) {
	case constants.TipoOperacaoEntrada, constants.TipoOperacaoSaida:
		return nil
	default:
		return fmt.Errorf("tipo_operacao inválido: %d. Deve ser 0 (Entrada) ou 1 (Saída)", valor)
	}
}

func validateTipoEntrega(valor string) error {
	validValues := []string{
		constants.TipoEntregaRetirada,
		constants.TipoEntregaEntrega,
		constants.TipoEntregaLocal,
	}

	for _, v := range validValues {
		if valor == v {
			return nil
		}
	}

	return fmt.Errorf("tipo_entrega inválido: %s. Deve ser RETIRADA, ENTREGA ou LOCAL", valor)
}

// ============================================================
// MÉTODOS DE CONVERSÃO (TO MODEL)
// ============================================================

// ToModel converte DocumentoVendaRequest para models.DocumentoVenda
func (r *DocumentoVendaRequest) ToModel() (*models.DocumentoVenda, error) {
	if r == nil {
		return nil, nil
	}

	doc := &models.DocumentoVenda{}

	// 1. Usar o mapper para copiar campos
	if err := utils.MapToModel(r, doc); err != nil {
		return nil, err
	}

	// 2. Tratamentos especiais
	// Converter tipos para as constantes
	doc.TipoDocumento = constants.TipoDocumentoVenda(r.TipoDocumento)
	doc.TipoOperacao = constants.TipoOperacao(r.TipoOperacao)

	// Definir situação padrão (Aberto)
	if doc.Situacao == 0 {
		doc.Situacao = constants.SituacaoPedidoAberto
	}

	return doc, nil
}

// ============================================================
// MÉTODOS DE CONVERSÃO (FROM MODEL)
// ============================================================

// FromModel converte models.DocumentoVenda para DocumentoVendaResponse
func (r *DocumentoVendaResponse) FromModel(doc *models.DocumentoVenda) (*DocumentoVendaResponse, error) {
	if doc == nil {
		return nil, nil
	}

	// 1. Usar o mapper para copiar campos
	if err := utils.MapToDTO(doc, r); err != nil {
		return nil, err
	}

	// 2. Preencher campos calculados (labels)
	r.TipoDocumentoLabel = constants.TipoDocumentoVenda(doc.TipoDocumento).String()
	r.TipoOperacaoLabel = constants.TipoOperacao(doc.TipoOperacao).String()
	r.SituacaoLabel = constants.SituacaoPedido(doc.Situacao).String()

	// 4. Formatar datas
	r.DataDocumento = utils.FormatDateTime(doc.CreatedAt)
	r.CreatedAt = utils.FormatDateTime(doc.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(doc.UpdatedAt)

	if doc.DataValidade != nil {
		r.DataValidade = utils.FormatDate(*doc.DataValidade)
	}

	if doc.DataPrevisao != nil {
		r.DataPrevisao = utils.FormatDate(*doc.DataPrevisao)
	}

	// 5. Calcular totais
	calculateTotals(r, doc)

	// 6. Converter itens
	convertItems(r, doc)

	// 7. Converter pagamentos
	convertPayments(r, doc)

	return r, nil
}

// ============================================================
// FUNÇÕES AUXILIARES DE CONVERSÃO
// ============================================================

func calculateTotals(r *DocumentoVendaResponse, doc *models.DocumentoVenda) {
	if len(doc.Itens) == 0 {
		return
	}

	var totalProdutos, totalDescontos, totalFrete float64

	for _, item := range doc.Itens {
		totalProdutos += item.ValorUnitario * item.Quantidade
		if item.ValorDesconto != nil {
			totalDescontos += *item.ValorDesconto
		}
		if item.ValorFrete != nil {
			totalFrete += *item.ValorFrete
		}
	}

	r.TotalProdutos = totalProdutos
	r.TotalDescontos = totalDescontos
	r.TotalFrete = totalFrete
	r.TotalLiquido = totalProdutos - totalDescontos + totalFrete
}

func convertItems(r *DocumentoVendaResponse, doc *models.DocumentoVenda) {
	if len(doc.Itens) == 0 {
		return
	}

	itensResponse := make([]DocumentoVendaItemResponse, len(doc.Itens))
	for i, item := range doc.Itens {
		itemResp := DocumentoVendaItemResponse{
			DocumentoVendaID:   item.DocumentoVendaID,
			Item:               item.Item,
			ProdutoID:          item.ProdutoID,
			Quantidade:         item.Quantidade,
			ValorUnitario:      item.ValorUnitario,
			TotalItem:          item.ValorUnitario * item.Quantidade,
			PesoBruto:          utils.Float64Ptr(item.PesoBruto),
			PesoLiquido:        utils.Float64Ptr(item.PesoLiquido),
			PercentualDesconto: item.PercentualDesconto,
			ValorDesconto:      item.ValorDesconto,
			ValorFrete:         item.ValorFrete,
		}

		if item.Produto != nil {
			itemResp.ProdutoNome = item.Produto.Nome
			itemResp.ProdutoCodigo = item.Produto.Codigo
		}
		itensResponse[i] = itemResp
	}
	r.Itens = itensResponse
}

func convertPayments(r *DocumentoVendaResponse, doc *models.DocumentoVenda) {
	if len(doc.Pagamentos) == 0 {
		return
	}

	pagamentosResponse := make([]DocumentoVendaPagamentoResponse, len(doc.Pagamentos))
	for i, pagamento := range doc.Pagamentos {
		pagResp := DocumentoVendaPagamentoResponse{
			DocumentoVendaID: pagamento.DocumentoVendaID,
			Item:             pagamento.Item,
			FormaPagamentoID: *pagamento.FormaPagamentoID,
			Valor:            pagamento.Valor, //
			ValorJuros:       pagamento.ValorJuros, //
			ValorMulta:       pagamento.ValorMulta, //
			ValorDesconto:    pagamento.ValorDesconto, //
			Observacao:       pagamento.Observacao, //
		}

		if pagamento.FormaPagamento != nil {
			pagResp.FormaPagamentoNome = pagamento.FormaPagamento.Descricao //
		}

		if !pagamento.DataVencimento.IsZero() {
			pagResp.DataVencimento = utils.FormatDate(pagamento.DataVencimento) //
		}
		if pagamento.DataPagamento != nil && !pagamento.DataPagamento.IsZero() {
			pagResp.DataPagamento = utils.FormatDate(*pagamento.DataPagamento) //
		}

		// Status label - usando constantes existentes
		// pagResp.StatusLabel = getStatusPagamentoLabel(pagamento.Status)

		pagamentosResponse[i] = pagResp
	}
	r.Pagamentos = pagamentosResponse
}
