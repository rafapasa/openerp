// ============================================================
// FILE: tipo_produto_dto.go
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

// TipoProdutoRequest representa a requisição para criar/atualizar um tipo de produto
type TipoProdutoRequest struct {
	// ============================================================
	// DADOS PRINCIPAIS
	// ============================================================
	Codigo            string  `json:"codigo" binding:"required"`
	Descricao         string  `json:"descricao" binding:"required"`
	ValidarFornecedor int     `json:"validar_fornecedor"`
	MovimentaEstoque  int     `json:"movimenta_estoque"`
	Sigla             *string `json:"sigla,omitempty"`
	ReceitaID         *int    `json:"receita_id,omitempty"`
	DespesaID         *int    `json:"despesa_id,omitempty"`
	ProdutoPacote     int     `json:"produto_pacote,omitempty"`
	Combustivel       int     `json:"combustivel,omitempty"`

	// ============================================================
	// USUÁRIO (para auditoria)
	// ============================================================
	CreatedBy *int `json:"created_by,omitempty"`
	UpdatedBy *int `json:"updated_by,omitempty"`
}

// ============================================================
// RESPONSES
// ============================================================

// TipoProdutoResponse representa a resposta de um tipo de produto
type TipoProdutoResponse struct {
	// ============================================================
	// DADOS PRINCIPAIS
	// ============================================================
	ID                int     `json:"id"`
	Codigo            string  `json:"codigo"`
	Descricao         string  `json:"descricao"`
	ValidarFornecedor int     `json:"validar_fornecedor"`
	ValidarFornecedorLabel string `json:"validar_fornecedor_label,omitempty"`
	MovimentaEstoque  int     `json:"movimenta_estoque"`
	MovimentaEstoqueLabel string `json:"movimenta_estoque_label,omitempty"`
	Sigla             *string `json:"sigla,omitempty"`
	ReceitaID         *int    `json:"receita_id,omitempty"`
	DespesaID         *int    `json:"despesa_id,omitempty"`
	ProdutoPacote     int     `json:"produto_pacote"`
	ProdutoPacoteLabel string  `json:"produto_pacote_label,omitempty"`
	Combustivel       int     `json:"combustivel"`
	CombustivelLabel  string  `json:"combustivel_label,omitempty"`

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

// TipoProdutoListResponse representa a resposta de listagem de tipos de produto
type TipoProdutoListResponse struct {
	Items      []TipoProdutoResponse `json:"items"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalPages int                   `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE CONVERSÃO (USANDO MAPPER)
// ============================================================

// ToModel converte TipoProdutoRequest para models.TipoProduto usando mapper
func (r *TipoProdutoRequest) ToModel() (*models.TipoProduto, error) {
	if r == nil {
		return nil, nil
	}

	tipoProduto := &models.TipoProduto{}

	// 1. Usar o mapper para copiar campos (mapeamento automático)
	if err := utils.MapToModel(r, tipoProduto); err != nil {
		return nil, err
	}

	// 2. Tratamentos especiais que o mapper não cobre
	// Valores padrão para campos booleanos (0 = não, 1 = sim)
	if r.ValidarFornecedor == 0 {
		tipoProduto.ValidarFornecedor = 0
	}
	if r.MovimentaEstoque == 0 {
		tipoProduto.MovimentaEstoque = 0
	}
	if r.ProdutoPacote == 0 {
		tipoProduto.ProdutoPacote = 0
	}
	if r.Combustivel == 0 {
		tipoProduto.Combustivel = 0
	}

	return tipoProduto, nil
}

// FromModel converte models.TipoProduto para TipoProdutoResponse usando mapper
func (r *TipoProdutoResponse) FromModel(tipoProduto *models.TipoProduto) *TipoProdutoResponse {
	if tipoProduto == nil {
		return nil
	}

	// 1. Usar o mapper para copiar campos (mapeamento automático)
	if err := utils.MapToDTO(tipoProduto, r); err != nil {
		// Se o mapper falhar, usar fallback manual
		return r.fromModelFallback(tipoProduto)
	}

	// 2. Preencher campos calculados (labels)
	r.ValidarFornecedorLabel = getSimNaoLabel(tipoProduto.ValidarFornecedor)
	r.MovimentaEstoqueLabel = getSimNaoLabel(tipoProduto.MovimentaEstoque)
	r.ProdutoPacoteLabel = getSimNaoLabel(tipoProduto.ProdutoPacote)
	r.CombustivelLabel = getSimNaoLabel(tipoProduto.Combustivel)

	// 3. Formatar datas (o mapper não faz isso)
	r.CreatedAt = utils.FormatDateTime(tipoProduto.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(tipoProduto.UpdatedAt)

	return r
}

// ============================================================
// FALLBACK (caso o mapper falhe)
// ============================================================

// fromModelFallback é o fallback manual caso o mapper falhe
func (r *TipoProdutoResponse) fromModelFallback(tipoProduto *models.TipoProduto) *TipoProdutoResponse {
	// Mapeamento manual campo por campo (seguro)
	r.ID = tipoProduto.ID
	r.Codigo = tipoProduto.Codigo
	r.Descricao = tipoProduto.Descricao
	r.ValidarFornecedor = tipoProduto.ValidarFornecedor
	r.MovimentaEstoque = tipoProduto.MovimentaEstoque
	r.Sigla = tipoProduto.Sigla
	r.ReceitaID = tipoProduto.ReceitaID
	r.DespesaID = tipoProduto.DespesaID
	r.ProdutoPacote = tipoProduto.ProdutoPacote
	r.Combustivel = tipoProduto.Combustivel
	r.CreatedBy = tipoProduto.CreatedBy
	r.UpdatedBy = tipoProduto.UpdatedBy

	// Labels
	r.ValidarFornecedorLabel = getSimNaoLabel(tipoProduto.ValidarFornecedor)
	r.MovimentaEstoqueLabel = getSimNaoLabel(tipoProduto.MovimentaEstoque)
	r.ProdutoPacoteLabel = getSimNaoLabel(tipoProduto.ProdutoPacote)
	r.CombustivelLabel = getSimNaoLabel(tipoProduto.Combustivel)

	// Datas
	r.CreatedAt = utils.FormatDateTime(tipoProduto.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(tipoProduto.UpdatedAt)

	return r
}

// ============================================================
// FUNÇÕES AUXILIARES
// ============================================================

// getSimNaoLabel retorna o label para campos booleanos (0/1)
func getSimNaoLabel(valor int) string {
	switch valor {
	case 1:
		return "Sim"
	case 0:
		return "Não"
	default:
		return "Não"
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o TipoProdutoRequest
func (r *TipoProdutoRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}