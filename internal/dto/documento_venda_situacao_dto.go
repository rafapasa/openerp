package dto

import "time"

// ============================================================
// DTOs para DocumentoVendaSituacao
// ============================================================

// DocumentoVendaSituacaoRequest representa os dados para criar ou atualizar uma situação de documento de venda.
type DocumentoVendaSituacaoRequest struct {
	Descricao string `json:"descricao" binding:"required,max=100"`
	CreatedBy *int   `json:"created_by,omitempty"`
	UpdatedBy *int   `json:"updated_by,omitempty"`
}

// DocumentoVendaSituacaoResponse representa os dados de uma situação de documento de venda retornada pela API.
type DocumentoVendaSituacaoResponse struct {
	ID        int       `json:"id"`
	Descricao string    `json:"descricao"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DocumentoVendaSituacaoListResponse representa a estrutura de resposta para uma lista paginada de situações de documento de venda.
type DocumentoVendaSituacaoListResponse struct {
	Items      []DocumentoVendaSituacaoResponse `json:"items"`
	Total      int64                            `json:"total"`
	Limit      int                              `json:"limit"`
	Offset     int                              `json:"offset"`
	TotalPages int                              `json:"total_pages"`
}

// DocumentoVendaSituacaoFilter representa os parâmetros de filtro para buscar situações de documento de venda.
type DocumentoVendaSituacaoFilter struct {
	Descricao string `form:"descricao"`
}