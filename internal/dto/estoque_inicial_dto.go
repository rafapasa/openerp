package dto

import "time"

// ============================================================
// DTOs para Estoque Inicial
// ============================================================

// EstoqueInicialRequest representa os dados para criar ou atualizar um registro de estoque inicial.
type EstoqueInicialRequest struct {
	ProdutoID       int       `json:"produto_id" binding:"required"`
	EmpresaFilialID int       `json:"empresa_filial_id" binding:"required"`
	Quantidade      float64   `json:"quantidade" binding:"required,gt=0"`
	DataLancamento  time.Time `json:"data_lancamento" binding:"required"`
	CreatedBy       *int      `json:"created_by,omitempty"`
	UpdatedBy       *int      `json:"updated_by,omitempty"`
}

// EstoqueInicialResponse representa os dados de um registro de estoque inicial retornado pela API.
type EstoqueInicialResponse struct {
	ID                int       `json:"id"`
	ProdutoID         int       `json:"produto_id"`
	ProdutoNome       string    `json:"produto_nome,omitempty"` // Adicionar nome do produto para facilitar
	EmpresaFilialID   int       `json:"empresa_filial_id"`
	EmpresaFilialNome string    `json:"empresa_filial_nome,omitempty"` // Adicionar nome da filial
	Quantidade        float64   `json:"quantidade"`
	DataLancamento    time.Time `json:"data_lancamento"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// EstoqueInicialListResponse representa a estrutura de resposta para uma lista paginada de estoques iniciais.
type EstoqueInicialListResponse struct {
	Items      []EstoqueInicialResponse `json:"items"`
	Total      int64                    `json:"total"`
	Limit      int                      `json:"limit"`
	Offset     int                      `json:"offset"`
	TotalPages int                      `json:"total_pages"`
}

// EstoqueInicialFilter representa os parâmetros de filtro para buscar estoques iniciais.
type EstoqueInicialFilter struct {
	ProdutoID       *int       `form:"produto_id"`
	EmpresaFilialID *int       `form:"empresa_filial_id"`
	DataInicio      *time.Time `form:"data_inicio"`
	DataFim         *time.Time `form:"data_fim"`
	Limit           int        `form:"limit,default=10"`
	Offset          int        `form:"offset,default=0"`
}
