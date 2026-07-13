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
	ID                 int      `json:"id,omitempty"`
	ProdutoID          int      `json:"produto_id" binding:"required"`
	Quantidade         float64  `json:"quantidade" binding:"required,gt=0"`
	ValorUnitario      float64  `json:"valor_unitario" binding:"gte=0"`
	PercentualDesconto *float64 `json:"percentual_desconto,omitempty"`
	ValorDesconto      *float64 `json:"valor_desconto,omitempty"`
	ValorFrete         *float64 `json:"valor_frete,omitempty"`
	DescricaoProduto   *string  `json:"descricao_produto,omitempty"`
	PesoBruto          *float64 `json:"peso_bruto,omitempty"`
	PesoLiquido        *float64 `json:"peso_liquido,omitempty"`
}

// Validate valida o item do documento
func (r *DocumentoVendaItemRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}

	// Validar percentual de desconto (se informado)
	if r.PercentualDesconto != nil && (*r.PercentualDesconto < 0 || *r.PercentualDesconto > 100) {
		return fmt.Errorf("percentual_desconto deve estar entre 0 e 100, recebido: %.2f", *r.PercentualDesconto)
	}

	// Validar valor de desconto (se informado)
	if r.ValorDesconto != nil && *r.ValorDesconto < 0 {
		return fmt.Errorf("valor_desconto não pode ser negativo, recebido: %.2f", *r.ValorDesconto)
	}

	return nil
}

// DocumentoVendaRequest representa o request para criar/atualizar um documento de venda.
type DocumentoVendaRequest struct {
	EmpresaFilialID     int                          `json:"empresa_filial_id" binding:"required"`
	EntidadeID          *int                         `json:"entidade_id"`
	CondicaoPagamentoID int                          `json:"condicao_pagamento_id" binding:"required"`
	TabelaPrecoID       int                          `json:"tabela_preco_id" binding:"required"`
	TipoDocumento       constants.TipoDocumentoVenda `json:"tipo_documento" binding:"required"` // 1-Orçamento, 2-Pedido
	TipoOperacao        constants.TipoOperacao       `json:"tipo_operacao" binding:"required"`  // 0-Entrada, 1-Saída
	Numero              int                          `json:"numero,omitempty"`
	DataValidade        *time.Time                   `json:"data_validade,omitempty"`
	DataPrevisao        *time.Time                   `json:"data_previsao,omitempty"`
	ValorFrete          *float64                     `json:"valor_frete,omitempty"`
	ValorDesconto       *float64                     `json:"valor_desconto,omitempty"`
	ObservacoesInterna  *string                      `json:"observacoes_interna,omitempty"`
	ObservacoesADM      *string                      `json:"observacoes_adm,omitempty"`
	ObservacoesCliente  *string                      `json:"observacoes_cliente,omitempty"`
	VendedorID          *int                         `json:"vendedor_id,omitempty"`
	TransportadoraID    *int                         `json:"transportadora_id,omitempty"`
	TipoEntrega         *string                      `json:"tipo_entrega,omitempty"` // RETIRADA, ENTREGA, LOCAL

	// Relacionamentos
	Itens []DocumentoVendaItemRequest `json:"itens" binding:"required,min=1,dive"`

	// Auditoria
	CreatedBy *int `json:"created_by,omitempty"`
	UpdatedBy *int `json:"updated_by,omitempty"`
}

// Validate valida o documento de venda
func (r *DocumentoVendaRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}

	// 1. Validar TipoDocumento
	if err := validateTipoDocumentoVenda(r.TipoDocumento); err != nil {
		return err
	}

	// 2. Validar TipoOperacao
	if err := validateTipoOperacao(r.TipoOperacao); err != nil {
		return err
	}

	// 3. Validar TipoEntrega (se informado)
	if r.TipoEntrega != nil {
		if err := validateTipoEntrega(*r.TipoEntrega); err != nil {
			return err
		}
	}

	// 4. Validar itens
	for i, item := range r.Itens {
		if err := item.Validate(); err != nil {
			return fmt.Errorf("item %d: %w", i+1, err)
		}
	}

	// 5. Validar datas
	if r.DataValidade != nil && r.DataPrevisao != nil {
		if r.DataValidade.Before(*r.DataPrevisao) {
			return fmt.Errorf("data_validade deve ser posterior a data_previsao")
		}
	}

	return nil
}

// ============================================================
// RESPONSES
// ============================================================

// DocumentoVendaItemResponse representa a resposta de um item de documento.
type DocumentoVendaItemResponse struct {
	ID                 int      `json:"id"`
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

// DocumentoVendaResponse representa a resposta de um documento de venda.
type DocumentoVendaResponse struct {
	// ============================================================
	// DADOS PRINCIPAIS
	// ============================================================
	ID                         int                          `json:"id"`
	EmpresaFilialID            int                          `json:"empresa_filial_id"`
	EntidadeID                 *int                         `json:"entidade_id,omitempty"`
	EntidadeNome               string                       `json:"entidade_nome,omitempty"`
	EntidadeDocumento          string                       `json:"entidade_documento,omitempty"`
	CondicaoPagamentoID        int                          `json:"condicao_pagamento_id"`
	CondicaoPagamentoDescricao string                       `json:"condicao_pagamento_descricao,omitempty"`
	TabelaPrecoID              int                          `json:"tabela_preco_id"`
	TabelaPrecoDescricao       string                       `json:"tabela_preco_descricao,omitempty"`
	Numero                     int                          `json:"numero"`
	TipoDocumento              constants.TipoDocumentoVenda `json:"tipo_documento"`
	TipoDocumentoLabel         string                       `json:"tipo_documento_label"`
	TipoOperacao               constants.TipoOperacao       `json:"tipo_operacao"`
	TipoOperacaoLabel          string                       `json:"tipo_operacao_label"`
	Situacao                   constants.SituacaoPedido     `json:"situacao"`
	SituacaoLabel              string                       `json:"situacao_label"`
	VendedorID                 *int                         `json:"vendedor_id,omitempty"`
	VendedorNome               string                       `json:"vendedor_nome,omitempty"`
	TransportadoraID           *int                         `json:"transportadora_id,omitempty"`
	TransportadoraNome         string                       `json:"transportadora_nome,omitempty"`
	TipoEntrega                *string                      `json:"tipo_entrega,omitempty"`
	TipoEntregaLabel           string                       `json:"tipo_entrega_label,omitempty"`

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
	// RELACIONAMENTOS
	// ============================================================
	Itens []DocumentoVendaItemResponse `json:"itens,omitempty"`

	// ============================================================
	// AUDITORIA
	// ============================================================
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	CreatedBy *int   `json:"created_by,omitempty"`
	UpdatedBy *int   `json:"updated_by,omitempty"`
}

// DocumentoVendaListResponse para listagem paginada.
type DocumentoVendaListResponse struct {
	Documentos []DocumentoVendaResponse `json:"items"`
	Total      int64                    `json:"total"`
	Page       int                      `json:"page"`
	Limit      int                      `json:"limit"`
	TotalPages int                      `json:"total_pages"`
}

// ============================================================
// FUNÇÕES DE CONVERSÃO
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
	doc.TipoDocumento = int(r.TipoDocumento)
	doc.TipoOperacao = int(r.TipoOperacao)

	// Definir situação padrão (Aberto)
	if doc.Situacao == 0 {
		doc.Situacao = constants.SituacaoPedidoAberto
	}

	return doc, nil
}

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
	r.SituacaoLabel = doc.Situacao.String()

	// 3. Label para TipoEntrega
	if doc.TipoEntrega != "" {
		r.TipoEntregaLabel = doc.TipoEntrega
	}

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

	// 5. Calcular totais (se tiver itens)
	if len(doc.Itens) > 0 {
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

	// 6. Converter itens
	if len(doc.Itens) > 0 {
		itensResponse := make([]DocumentoVendaItemResponse, len(doc.Itens))
		for i, item := range doc.Itens {
			itemResp := DocumentoVendaItemResponse{
				ID:            item.ID,
				ProdutoID:     item.ProdutoID,
				Quantidade:    item.Quantidade,
				ValorUnitario: item.ValorUnitario,
				TotalItem:     item.ValorUnitario * item.Quantidade,
			}

			if item.Produto != nil {
				itemResp.ProdutoNome = item.Produto.Nome
				itemResp.ProdutoCodigo = item.Produto.Codigo
			}

			if item.PercentualDesconto != nil {
				itemResp.PercentualDesconto = item.PercentualDesconto
			}

			if item.ValorDesconto != nil {
				itemResp.ValorDesconto = item.ValorDesconto
				itemResp.TotalItem -= *item.ValorDesconto
			}

			if item.ValorFrete != nil {
				itemResp.ValorFrete = item.ValorFrete
				itemResp.TotalItemComFrete = itemResp.TotalItem + *item.ValorFrete
			}

			if item.PesoBruto >= 0 {
				itemResp.PesoBruto = &item.PesoBruto
			}

			if item.PesoLiquido >= 0 {
				itemResp.PesoLiquido = &item.PesoLiquido
			}

			itensResponse[i] = itemResp
		}
		r.Itens = itensResponse
	}

	return r, nil
}

// ============================================================
// FUNÇÕES DE VALIDAÇÃO
// ============================================================

func validateTipoDocumentoVenda(valor constants.TipoDocumentoVenda) error {
	switch constants.TipoDocumentoVenda(valor) {
	case constants.TipoDocumentoOrcamento, constants.TipoDocumentoPedido:
		return nil
	default:
		return fmt.Errorf("tipo_documento inválido: %d. Deve ser 1 (Orçamento) ou 2 (Pedido)", valor)
	}
}

func validateTipoOperacao(valor constants.TipoOperacao) error {
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
// FUNÇÕES AUXILIARES
// ============================================================

func getTipoEntregaLabel(valor string) string {
	switch valor {
	case constants.TipoEntregaRetirada:
		return "Retirada"
	case constants.TipoEntregaEntrega:
		return "Entrega"
	case constants.TipoEntregaLocal:
		return "Local"
	default:
		return "Desconhecido"
	}
}
