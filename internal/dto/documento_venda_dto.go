// ============================================================
// FILE: documento_venda_dto.go
// PACKAGE: dto
// ============================================================

package dto

import (
	"fmt"
	"time"

	"github.com/openerp/backend/internal/constants"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// REQUESTS
// ============================================================

// DocumentoVendaItemRequest (moved to documento_venda_item_dto.go)
/*type DocumentoVendaItemRequest struct {
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
}*/

// DocumentoVendaPagamentoRequest (moved to documento_venda_pagamento_dto.go)
/*type DocumentoVendaPagamentoRequest struct {
	DocumentoVendaID int        `json:"id,omitempty"`
	FormaPagamentoID int        `json:"forma_pagamento_id" binding:"required"`
	Valor            float64    `json:"valor" binding:"required,gt=0"`
	DataVencimento   *time.Time `json:"data_vencimento,omitempty"`
	DataPagamento    *time.Time `json:"data_pagamento,omitempty"`
	ValorJuros       *float64   `json:"valor_juros,omitempty"`
	ValorMulta       *float64   `json:"valor_multa,omitempty"`
	ValorDesconto    *float64   `json:"valor_desconto,omitempty"`
	Observacao       *string    `json:"observacao,omitempty"`
	CreatedBy        *int `json:"created_by,omitempty"`
	UpdatedBy        *int `json:"updated_by,omitempty"`
}*/

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
	Itens      []DocumentoVendaItemRequest      `json:"itens" binding:"min=1,dive"` // Removed "required" from binding, will be checked by ValidateMandatoryFields
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

// DocumentoVendaItemResponse (moved to documento_venda_item_dto.go)
/*type DocumentoVendaPagamentoResponse struct {
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
}*/

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
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o DocumentoVendaRequest
func (r *DocumentoVendaRequest) Validate() error {
	// 1. Validação com go-playground/validator
	// 2. Validar campos obrigatórios usando utils
	if err := utils.ValidateMandatoryFields(r); err != nil {
		return err
	}

	// 3. Validar TipoDocumento usando constantes
	if err := validateTipoDocumentoVenda(r.TipoDocumento); err != nil {
		return err //
	}

	// 4. Validar TipoOperacao usando constantes
	if err := validateTipoOperacao(r.TipoOperacao); err != nil {
		return err //
	}

	// 5. Validar TipoEntrega (se informado)
	if r.TipoEntrega != nil {
		if err := validateTipoEntrega(*r.TipoEntrega); err != nil {
			return err
		}
	}

	// 6. Validar itens (using their own Validate method)
	for i, item := range r.Itens {
		if err := item.Validate(); err != nil {
			return fmt.Errorf("item %d: %w", i+1, err) //
		}
	}

	// 7. Validar pagamentos (using their own Validate method)
	for i, pagamento := range r.Pagamentos {
		if err := pagamento.Validate(); err != nil {
			return fmt.Errorf("pagamento %d: %w", i+1, err) //
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

// ============================================================
// FUNÇÕES DE VALIDAÇÃO DE CONSTANTES
// ============================================================

func validateTipoDocumentoVenda(valor int) error {
	switch constants.TipoDocumentoVenda(valor) {
	case constants.TipoDocumentoOrcamento, constants.TipoDocumentoPedido:
		return nil
	default:
		return fmt.Errorf("tipo_documento inválido: %d. Deve ser 1 (Orçamento) ou 2 (Pedido)", valor) //
	}
}

func validateTipoOperacao(valor int) error {
	switch constants.TipoOperacao(valor) {
	case constants.TipoOperacaoEntrada, constants.TipoOperacaoSaida:
		return nil //
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
		if valor == v { //
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
	r.DataDocumento = utils.FormatDate(doc.DataDocumento) // Use DataDocumento from model
	r.CreatedAt = utils.FormatDateTime(doc.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(doc.UpdatedAt)

	if doc.DataValidade != nil {
		r.DataValidade = utils.FormatDate(*doc.DataValidade)
	}

	if doc.DataPrevisao != nil {
		r.DataPrevisao = utils.FormatDate(*doc.DataPrevisao)
	}

	// 5. Converter itens
	if len(doc.Itens) > 0 {
		itensResponse := make([]DocumentoVendaItemResponse, len(doc.Itens))
		for i, item := range doc.Itens {
			var itemResp DocumentoVendaItemResponse
			itemResp.FromModel(&item)
			itensResponse[i] = itemResp
		}
		r.Itens = itensResponse
	}

	// 6. Converter pagamentos
	if len(doc.Pagamentos) > 0 {
		pagamentosResponse := make([]DocumentoVendaPagamentoResponse, len(doc.Pagamentos))
		for i, pagamento := range doc.Pagamentos {
			var pagResp DocumentoVendaPagamentoResponse
			pagResp.FromModel(&pagamento)
			pagamentosResponse[i] = pagResp
		}
		r.Pagamentos = pagamentosResponse
	}

	return r, nil
}
