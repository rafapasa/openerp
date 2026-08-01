// ============================================================
// FILE: documento_venda_item_dto.go
// PACKAGE: dto
// ============================================================

package dto

import (
	"github.com/go-playground/validator/v10"
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
	OperacaoFiscalID   *int     `json:"operacao_fiscal_id,omitempty"`
	CstIcmsId          *int     `json:"cst_icms_id,omitempty"`
	CstIpiId           *int     `json:"cst_ipi_id,omitempty"`
	CstPisCofinsId     *int     `json:"cst_pis_cofins_id,omitempty"`
	CreatedBy          *int     `json:"created_by,omitempty"`
	UpdatedBy          *int     `json:"updated_by,omitempty"`
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

// DocumentoVendaItemListResponse representa a resposta de listagem de itens de documento de venda
type DocumentoVendaItemListResponse struct {
	Items      []DocumentoVendaItemResponse `json:"items"`
	Total      int64                        `json:"total"`
	Page       int                          `json:"page"`
	Limit      int                          `json:"limit"`
	TotalPages int                          `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE CONVERSÃO
// ============================================================

// ToModel converte DocumentoVendaItemRequest para models.DocumentoVendaItem
func (r *DocumentoVendaItemRequest) ToModel() (*models.DocumentoVendaItem, error) {
	if r == nil {
		return nil, nil
	}
	item := &models.DocumentoVendaItem{}
	if err := utils.MapToModel(r, item); err != nil {
		return nil, err
	}
	return item, nil
}

// FromModel converte models.DocumentoVendaItem para DocumentoVendaItemResponse
func (r *DocumentoVendaItemResponse) FromModel(item *models.DocumentoVendaItem) {
	if item == nil {
		return
	}
	_ = utils.MapToDTO(item, r) // Ignore error for now, assume direct mapping is sufficient

	if item.Produto != nil {
		r.ProdutoNome = item.Produto.Nome
		r.ProdutoCodigo = item.Produto.Codigo
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o DocumentoVendaItemRequest
func (r *DocumentoVendaItemRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
