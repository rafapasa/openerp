package dto

import (
	"github.com/openerp/backend/internal/constants"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// TabelaPrecoProdutoRequest representa a requisição para criar/atualizar um produto em uma tabela de preço.
type TabelaPrecoProdutoRequest struct {
	TabelaPrecoID       int      `json:"-"` // Virá da URL
	ProdutoID           int      `json:"produto_id" binding:"required"`
	ValorMinimo         float64  `json:"valor_minimo" binding:"required,gte=0"`
	ValorPadrao         float64  `json:"valor_padrao" binding:"required,gte=0"`
	Situacao            int      `json:"situacao,omitempty"`
	QuantAtacado        *float64 `json:"quant_atacado,omitempty"`
	PercDescontoAtacado *float64 `json:"perc_desconto_atacado,omitempty"`
	ValorCusto          *float64 `json:"valor_custo,omitempty"`
	MargemLucro         *float64 `json:"margem_lucro,omitempty"`
	CreatedBy           *int     `json:"-"`
	UpdatedBy           *int     `json:"-"`
}

// TabelaPrecoProdutoResponse representa a resposta de um produto em uma tabela de preço.
type TabelaPrecoProdutoResponse struct {
	TabelaPrecoID       int      `json:"tabela_preco_id"`
	Item                int      `json:"item"`
	ProdutoID           int      `json:"produto_id"`
	ProdutoNome         string   `json:"produto_nome,omitempty"`
	ValorMinimo         float64  `json:"valor_minimo"`
	ValorPadrao         float64  `json:"valor_padrao"`
	Situacao            int      `json:"situacao"`
	SituacaoLabel       string   `json:"situacao_label"`
	QuantAtacado        *float64 `json:"quant_atacado,omitempty"`
	PercDescontoAtacado *float64 `json:"perc_desconto_atacado,omitempty"`
	ValorCusto          *float64 `json:"valor_custo,omitempty"`
	MargemLucro         *float64 `json:"margem_lucro,omitempty"`
	CreatedAt           string   `json:"created_at"`
	UpdatedAt           string   `json:"updated_at"`
	CreatedBy           *int     `json:"created_by,omitempty"`
	UpdatedBy           *int     `json:"updated_by,omitempty"`
}

// TabelaPrecoProdutoListResponse representa a resposta de listagem de produtos em uma tabela.
type TabelaPrecoProdutoListResponse struct {
	Items      []TabelaPrecoProdutoResponse `json:"items"`
	Total      int64                        `json:"total"`
	Page       int                          `json:"page"`
	Limit      int                          `json:"limit"`
	TotalPages int                          `json:"total_pages"`
}

// ToModel converte TabelaPrecoProdutoRequest para models.TabelaPrecoProduto.
func (r *TabelaPrecoProdutoRequest) ToModel() (*models.TabelaPrecoProduto, error) {
	if r == nil {
		return nil, nil
	}
	itemTab := &models.TabelaPrecoProduto{}
	if err := utils.MapToModel(r, itemTab); err != nil {
		return nil, err
	}
	if r.Situacao == 0 {
		itemTab.Situacao = int(constants.StatusAtivo)
	}
	return itemTab, nil
}

// FromModel converte models.TabelaPrecoProduto para TabelaPrecoProdutoResponse.
func (r *TabelaPrecoProdutoResponse) FromModel(item *models.TabelaPrecoProduto) {
	if item == nil {
		return
	}

	if err := utils.MapToDTO(item, r); err != nil {
		// Fallback manual
		r.Item = item.Item
		r.TabelaPrecoID = item.TabelaPrecoID
		r.ProdutoID = item.ProdutoID
		r.ValorMinimo = item.ValorMinimo
		r.ValorPadrao = item.ValorPadrao
		r.Situacao = item.Situacao
		r.QuantAtacado = item.QuantAtacado
		r.PercDescontoAtacado = item.PercDescontoAtacado
		r.ValorCusto = item.ValorCusto
		r.MargemLucro = item.MargemLucro
		r.CreatedBy = item.CreatedBy
		r.UpdatedBy = item.UpdatedBy
	}

	r.SituacaoLabel = constants.Status(item.Situacao).String()
	r.CreatedAt = utils.FormatDateTime(item.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(item.UpdatedAt)

	if item.Produto != nil {
		r.ProdutoNome = item.Produto.Nome
	}
}

// Validate valida o TabelaPrecoProdutoRequest.
func (r *TabelaPrecoProdutoRequest) Validate() error {
	if err := utils.ValidateMandatoryFields(r); err != nil {
		return err
	}

	if err := constants.Status(r.Situacao).IsValid(); err != nil {
		return err
	}

	return nil
}
