package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/openerp/backend/internal/constants"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// REQUESTS
// ============================================================

// ProdutoRequest representa a requisição para criar/atualizar um produto
type ProdutoRequest struct {
	// ============================================================
	// DADOS PRINCIPAIS
	// ============================================================
	EmpresaFilialID   int      `json:"empresa_filial_id" binding:"required"`
	TipoProdutoID     int      `json:"tipo_produto_id" binding:"required"`
	ProdutoGrupoID    *int     `json:"produto_grupo_id,omitempty"`
	ProdutoSubgrupoID *int     `json:"produto_subgrupo_id,omitempty"`
	MarcaID           *int     `json:"marca_id,omitempty"`
	ModeloID          *int     `json:"modelo_id,omitempty"`
	SerieID           *int     `json:"serie_id,omitempty"`
	EspecieID         *int     `json:"especie_id,omitempty"`
	NcmNumero         int      `json:"ncm_numero" binding:"required"`
	Codigo            int      `json:"codigo" binding:"required"`
	Nome              string   `json:"nome" binding:"required"`
	Resumo            string   `json:"resumo" binding:"required"`
	Situacao          int      `json:"situacao,omitempty"` // 1-Ativo, 2-Inativo, 3-Bloqueado
	CustoCompra       *float64 `json:"custo_compra,omitempty"`
	CodigoBarras      string   `json:"codigo_barras,omitempty"`

	// ============================================================
	// DIMENSÕES E PESOS
	// ============================================================
	PesoBruto   *float64 `json:"peso_bruto,omitempty"`
	PesoLiquido *float64 `json:"peso_liquido,omitempty"`
	Altura      *float64 `json:"altura,omitempty"`
	Largura     *float64 `json:"largura,omitempty"`
	Comprimento *float64 `json:"comprimento,omitempty"`

	// ============================================================
	// ESTOQUE
	// ============================================================
	EstoqueMinimo *float64 `json:"estoque_minimo,omitempty"`
	LoteEconomico *float64 `json:"lote_economico,omitempty"`

	// ============================================================
	// REFERÊNCIAS E DESCRIÇÃO
	// ============================================================
	Referencia     string   `json:"referencia,omitempty"`
	Referencia2    string   `json:"referencia2,omitempty"`
	Referencia3    string   `json:"referencia3,omitempty"`
	Referencia4    string   `json:"referencia4,omitempty"`
	Descricao      string   `json:"descricao,omitempty"`
	DescontoMaximo *float64 `json:"desconto_maximo,omitempty"`

	// ============================================================
	// AUDITORIA
	// ============================================================
	CreatedBy *int `json:"created_by,omitempty"`
	UpdatedBy *int `json:"updated_by,omitempty"`
}

// ============================================================
// RESPONSES
// ============================================================

// ProdutoResponse representa a resposta de um produto
type ProdutoResponse struct {
	// ============================================================
	// DADOS PRINCIPAIS
	// ============================================================
	ID                  int      `json:"id"`
	EmpresaFilialID     int      `json:"empresa_filial_id"`
	EmpresaFilialNome   string   `json:"empresa_filial_nome,omitempty"`
	TipoProdutoID       int      `json:"tipo_produto_id"`
	TipoProdutoNome     string   `json:"tipo_produto_nome,omitempty"`
	ProdutoGrupoID      *int     `json:"produto_grupo_id,omitempty"`
	ProdutoGrupoNome    string   `json:"produto_grupo_nome,omitempty"`
	ProdutoSubgrupoID   *int     `json:"produto_subgrupo_id,omitempty"`
	ProdutoSubgrupoNome string   `json:"produto_subgrupo_nome,omitempty"`
	MarcaID             *int     `json:"marca_id,omitempty"`
	MarcaNome           string   `json:"marca_nome,omitempty"`
	ModeloID            *int     `json:"modelo_id,omitempty"`
	ModeloNome          string   `json:"modelo_nome,omitempty"`
	SerieID             *int     `json:"serie_id,omitempty"`
	SerieNome           string   `json:"serie_nome,omitempty"`
	EspecieID           *int     `json:"especie_id,omitempty"`
	EspecieNome         string   `json:"especie_nome,omitempty"`
	NcmNumero           int      `json:"ncm_numero"`
	NcmDescricao        string   `json:"ncm_descricao,omitempty"`
	Codigo              int      `json:"codigo"`
	Nome                string   `json:"nome"`
	Resumo              string   `json:"resumo"`
	Situacao            int      `json:"situacao"`
	SituacaoLabel       string   `json:"situacao_label"`
	CustoCompra         *float64 `json:"custo_compra,omitempty"`
	CodigoBarras        string   `json:"codigo_barras,omitempty"`

	// ============================================================
	// DIMENSÕES E PESOS
	// ============================================================
	PesoBruto   *float64 `json:"peso_bruto,omitempty"`
	PesoLiquido *float64 `json:"peso_liquido,omitempty"`
	Altura      *float64 `json:"altura,omitempty"`
	Largura     *float64 `json:"largura,omitempty"`
	Comprimento *float64 `json:"comprimento,omitempty"`

	// ============================================================
	// ESTOQUE
	// ============================================================
	EstoqueMinimo *float64 `json:"estoque_minimo,omitempty"`
	LoteEconomico *float64 `json:"lote_economico,omitempty"`
	SaldoEstoque  *float64 `json:"saldo_estoque,omitempty"`

	// ============================================================
	// REFERÊNCIAS E DESCRIÇÃO
	// ============================================================
	Referencia     string   `json:"referencia,omitempty"`
	Referencia2    string   `json:"referencia2,omitempty"`
	Referencia3    string   `json:"referencia3,omitempty"`
	Referencia4    string   `json:"referencia4,omitempty"`
	Descricao      string   `json:"descricao,omitempty"`
	DescontoMaximo *float64 `json:"desconto_maximo,omitempty"`
	DataAlteracao  string   `json:"data_alteracao,omitempty"`

	// ============================================================
	// AUDITORIA
	// ============================================================
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	CreatedBy *int   `json:"created_by,omitempty"`
	UpdatedBy *int   `json:"updated_by,omitempty"`
}

// ============================================================
// LIST RESPONSE
// ============================================================

// ProdutoListResponse representa a resposta de listagem de produtos
type ProdutoListResponse struct {
	Items      []ProdutoResponse `json:"items"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE CONVERSÃO
// ============================================================

// ToModel converte ProdutoRequest para models.Produto
func (r *ProdutoRequest) ToModel() (*models.Produto, error) {
	if r == nil {
		return nil, nil
	}

	produto := &models.Produto{
		EmpresaFilialID:   r.EmpresaFilialID,
		TipoProdutoID:     r.TipoProdutoID,
		ProdutoGrupoID:    r.ProdutoGrupoID,
		ProdutoSubgrupoID: r.ProdutoSubgrupoID,
		MarcaID:           r.MarcaID,
		ModeloID:          r.ModeloID,
		SerieID:           r.SerieID,
		EspecieID:         r.EspecieID,
		NcmNumero:         r.NcmNumero,
		Codigo:            r.Codigo,
		Nome:              r.Nome,
		Resumo:            r.Resumo,
		CustoCompra:       r.CustoCompra,
		CodigoBarras:      utils.StringPtr(r.CodigoBarras),
		PesoBruto:         r.PesoBruto,
		PesoLiquido:       r.PesoLiquido,
		Altura:            r.Altura,
		Largura:           r.Largura,
		Comprimento:       r.Comprimento,
		EstoqueMinimo:     r.EstoqueMinimo,
		LoteEconomico:     r.LoteEconomico,
		Referencia:        utils.StringPtr(r.Referencia),
		Referencia2:       utils.StringPtr(r.Referencia2),
		Referencia3:       utils.StringPtr(r.Referencia3),
		Referencia4:       utils.StringPtr(r.Referencia4),
		Descricao:         utils.StringPtr(r.Descricao),
		DescontoMaximo:    r.DescontoMaximo,
		CreatedBy:         r.CreatedBy,
		UpdatedBy:         r.UpdatedBy,
	}

	// Definir situação
	if r.Situacao == 0 {
		produto.Situacao = constants.StatusAtivo
	} else {
		produto.Situacao = constants.Status(r.Situacao)
	}

	return produto, nil
}

// FromModel converte models.Produto para ProdutoResponse
func (r *ProdutoResponse) FromModel(produto *models.Produto) *ProdutoResponse {
	if produto == nil {
		return nil
	}

	r.ID = produto.ID
	r.EmpresaFilialID = produto.EmpresaFilialID
	r.TipoProdutoID = produto.TipoProdutoID
	r.ProdutoGrupoID = produto.ProdutoGrupoID
	r.ProdutoSubgrupoID = produto.ProdutoSubgrupoID
	r.MarcaID = produto.MarcaID
	r.ModeloID = produto.ModeloID
	r.SerieID = produto.SerieID
	r.EspecieID = produto.EspecieID
	r.NcmNumero = produto.NcmNumero
	r.Codigo = produto.Codigo
	r.Nome = produto.Nome
	r.Resumo = produto.Resumo
	r.Situacao = int(produto.Situacao)
	r.SituacaoLabel = produto.Situacao.String()
	r.CustoCompra = produto.CustoCompra
	r.CodigoBarras = utils.StringValue(produto.CodigoBarras)

	// Dimensões e Pesos
	r.PesoBruto = produto.PesoBruto
	r.PesoLiquido = produto.PesoLiquido
	r.Altura = produto.Altura
	r.Largura = produto.Largura
	r.Comprimento = produto.Comprimento

	// Estoque
	r.EstoqueMinimo = produto.EstoqueMinimo
	r.LoteEconomico = produto.LoteEconomico

	// Referências e Descrição
	r.Referencia = utils.StringValue(produto.Referencia)
	r.Referencia2 = utils.StringValue(produto.Referencia2)
	r.Referencia3 = utils.StringValue(produto.Referencia3)
	r.Referencia4 = utils.StringValue(produto.Referencia4)
	r.Descricao = utils.StringValue(produto.Descricao)
	r.DescontoMaximo = produto.DescontoMaximo

	if produto.DataAlteracao != nil {
		r.DataAlteracao = produto.DataAlteracao.Format("2006-01-02 15:04:05")
	}

	// Auditoria
	r.CreatedAt = produto.CreatedAt.Format("2006-01-02 15:04:05")
	r.UpdatedAt = produto.UpdatedAt.Format("2006-01-02 15:04:05")
	r.CreatedBy = produto.CreatedBy
	r.UpdatedBy = produto.UpdatedBy

	// Preencher relacionamentos (se carregados)
	if produto.TipoProduto != nil {
		r.TipoProdutoNome = produto.TipoProduto.Descricao
	}
	if produto.ProdutoGrupo != nil {
		r.ProdutoGrupoNome = produto.ProdutoGrupo.Descricao
	}
	if produto.ProdutoSubgrupo != nil {
		r.ProdutoSubgrupoNome = produto.ProdutoSubgrupo.Descricao
	}
	if produto.Marca != nil {
		r.MarcaNome = produto.Marca.Descricao
	}
	if produto.Modelo != nil {
		r.ModeloNome = produto.Modelo.Descricao
	}
	if produto.Serie != nil {
		r.SerieNome = produto.Serie.Descricao
	}
	if produto.Especie != nil {
		r.EspecieNome = produto.Especie.Descricao
	}

	return r
}

// ============================================================
// FUNÇÕES AUXILIARES
// ============================================================

// IsValidSituacaoProduto verifica se a situação do produto é válida
func IsValidSituacaoProduto(situacao int) bool {
	switch situacao {
	case int(constants.StatusAtivo):
		return true
	case int(constants.StatusInativo), int(constants.StatusBloqueado):
		return false
	default:
		return false
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o ProdutoRequest
func (r *ProdutoRequest) Validate() error {
	validate := validator.New()

	// Validação customizada para situação
	if r.Situacao > 0 && !IsValidSituacaoProduto(r.Situacao) {
		return fmt.Errorf("situação inválida: %d", r.Situacao)
	}

	return validate.Struct(r)
}
