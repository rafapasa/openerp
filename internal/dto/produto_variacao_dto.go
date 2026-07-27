package dto

import (
	"time"
)

// ============================================================
// DTO: ProdutoVariacaoRequest
// Representa os dados de entrada para criação/atualização de uma variação de produto.
// ============================================================

type ProdutoVariacaoRequest struct {
	// CAMPOS PRINCIPAIS
	ID              int     `json:"id"` // Usado para validação em updates
	ProdutoID       int     `json:"produto_id" binding:"required"`
	EmpresaFilialID int     `json:"empresa_filial_id" binding:"required"`
	CorID           *int    `json:"cor_id"`
	TamanhoID       *int    `json:"tamanho_id"`
	SKU             string  `json:"sku" binding:"required,max=50"`
	PrecoAdicional  float64 `json:"preco_adicional"`
	EstoqueAtual    float64 `json:"estoque_atual"` // Para o estoque inicial

	// AUDITORIA
	CreatedBy *int `json:"created_by"`
	UpdatedBy *int `json:"updated_by"`
}

// ============================================================
// DTO: ProdutoVariacaoResponse
// Representa os dados de saída de uma variação de produto.
// ============================================================

type ProdutoVariacaoResponse struct {
	// DADOS PRINCIPAIS
	ID              int     `json:"id"`
	ProdutoID       int     `json:"produto_id"`
	EmpresaFilialID int     `json:"empresa_filial_id"`
	CorID           *int    `json:"cor_id,omitempty"`
	TamanhoID       *int    `json:"tamanho_id,omitempty"`
	SKU             string  `json:"sku"`
	PrecoAdicional  float64 `json:"preco_adicional"`
	EstoqueAtual    float64 `json:"estoque_atual"`

	// RELACIONAMENTOS (apenas IDs e Nomes para evitar recursão)
	ProdutoNome       string `json:"produto_nome,omitempty"`
	EmpresaFilialNome string `json:"empresa_filial_nome,omitempty"`
	CorNome           string `json:"cor_nome,omitempty"`
	TamanhoNome       string `json:"tamanho_nome,omitempty"`
	CorSigla          string `json:"cor_sigla,omitempty"`
	TamanhoSigla      string `json:"tamanho_sigla,omitempty"`

	// AUDITORIA
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	CreatedBy *int       `json:"created_by,omitempty"`
	UpdatedBy *int       `json:"updated_by,omitempty"`
}

// ============================================================
// DTO: ProdutoVariacaoListResponse
// Representa a estrutura de resposta para listagem de variações de produto.
// ============================================================

type ProdutoVariacaoListResponse struct {
	Items      []ProdutoVariacaoResponse `json:"items"`
	Total      int64                     `json:"total"`
	Limit      int                       `json:"limit"`
	Page       int                       `json:"page"`
	TotalPages int                       `json:"total_pages"`
}
