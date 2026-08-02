package dto

import "time"

// ============================================================
// DTOs para DocumentoVendaHistorico
// ============================================================

// DocumentoVendaHistoricoRequest representa os dados para registrar uma mudança de situação de documento de venda.
type DocumentoVendaHistoricoRequest struct {
	DocumentoVendaID int    `json:"documento_venda_id" binding:"required"`
	SituacaoAnterior int    `json:"situacao_anterior" binding:"required"`
	SituacaoAtual    int    `json:"situacao_atual" binding:"required"`
	Observacao       string `json:"observacao,omitempty"`
	CreatedBy        *int   `json:"created_by,omitempty"`
}

// DocumentoVendaHistoricoResponse representa os dados de um registro de histórico de situação de documento de venda retornado pela API.
type DocumentoVendaHistoricoResponse struct {
	ID                   int       `json:"id"`
	DocumentoVendaID     int       `json:"documento_venda_id"`
	SituacaoAnterior     int       `json:"situacao_anterior"`
	SituacaoAnteriorDesc string    `json:"situacao_anterior_desc,omitempty"` // Descrição da situação anterior
	SituacaoAtual        int       `json:"situacao_atual"`
	SituacaoAtualDesc    string    `json:"situacao_atual_desc,omitempty"` // Descrição da situação atual
	Observacao           string    `json:"observacao,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	CreatedBy            *int      `json:"created_by,omitempty"`
	CreatedByName        string    `json:"created_by_name,omitempty"` // Nome do usuário que realizou a mudança
}

// DocumentoVendaHistoricoListResponse representa a estrutura de resposta para uma lista paginada de históricos de situação de documento de venda.
type DocumentoVendaHistoricoListResponse struct {
	Items      []DocumentoVendaHistoricoResponse `json:"items"`
	Total      int64                             `json:"total"`
	Limit      int                               `json:"limit"`
	Offset     int                               `json:"offset"`
	TotalPages int                               `json:"total_pages"`
}

// DocumentoVendaHistoricoFilter representa os parâmetros de filtro para buscar históricos de situação de documento de venda.
type DocumentoVendaHistoricoFilter struct {
	DocumentoVendaID *int       `form:"documento_venda_id"`
	SituacaoAnterior *int       `form:"situacao_anterior"`
	SituacaoAtual    *int       `form:"situacao_atual"`
	DataInicio       *time.Time `form:"data_inicio"`
	DataFim          *time.Time `form:"data_fim"`
	Limit            int        `form:"limit,default=10"`
	Offset           int        `form:"offset,default=0"`
}
