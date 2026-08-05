package dto

import (
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// DTO: DocumentoVendaHistoricoRequest
// ============================================================

type DocumentoVendaHistoricoRequest struct {
	DocumentoVendaID int    `json:"documento_venda_id" binding:"required"`
	Item             int    `json:"item" binding:"required"`
	UsuarioID        int    `json:"usuario_id" binding:"required"`
	FluxoID          int    `json:"fluxo_id" binding:"required"`
	FluxoSequencia   int    `json:"fluxo_sequencia" binding:"required"`
	DataHistorico    string `json:"data_historico" binding:"required"` // Formato: 2006-01-02
	Descricao        string `json:"descricao" binding:"required"`
	Motivo           string `json:"motivo,omitempty"`
	CreatedBy        *int   `json:"created_by,omitempty"`
	UpdatedBy        *int   `json:"updated_by,omitempty"`
}

// ============================================================
// DTO: DocumentoVendaHistoricoResponse
// ============================================================

type DocumentoVendaHistoricoResponse struct {
	DocumentoVendaID int    `json:"documento_venda_id"`
	Item             int    `json:"item"`
	UsuarioID        int    `json:"usuario_id"`
	UsuarioNome      string `json:"usuario_nome,omitempty"`
	FluxoID          int    `json:"fluxo_id"`
	FluxoSequencia   int    `json:"fluxo_sequencia"`
	DataHistorico    string `json:"data_historico"`
	Descricao        string `json:"descricao"`
	Motivo           string `json:"motivo,omitempty"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// ============================================================
// DTO: DocumentoVendaHistoricoListResponse
// ============================================================

type DocumentoVendaHistoricoListResponse struct {
	Items      []DocumentoVendaHistoricoResponse `json:"items"`
	Total      int64                             `json:"total"`
	Page       int                               `json:"page"`
	Limit      int                               `json:"limit"`
	TotalPages int                               `json:"total_pages"`
}

// ============================================================
// DTO: DocumentoVendaHistoricoFilter
// ============================================================

type DocumentoVendaHistoricoFilter struct {
	DocumentoVendaID *int    `form:"documento_venda_id" json:"documento_venda_id,omitempty"`
	UsuarioID        *int    `form:"usuario_id" json:"usuario_id,omitempty"`
	FluxoID          *int    `form:"fluxo_id" json:"fluxo_id,omitempty"`
	DataInicio       *string `form:"data_inicio" json:"data_inicio,omitempty"` // Formato: 2006-01-02
	DataFim          *string `form:"data_fim" json:"data_fim,omitempty"`       // Formato: 2006-01-02
	Page             int     `form:"page" json:"page,omitempty"`
	Limit            int     `form:"limit" json:"limit,omitempty"`
	Sort             string  `form:"sort" json:"sort,omitempty"`
	Order            string  `form:"order" json:"order,omitempty"`
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (r *DocumentoVendaHistoricoRequest) ToModel() (*models.DocumentoVendaHistorico, error) {
	if r == nil {
		return nil, nil
	}

	ddvh := &models.DocumentoVendaHistorico{}

	// 1. Usar o mapper para copiar campos
	if err := utils.MapToModel(r, ddvh); err != nil {
		return nil, err
	}

	return ddvh, nil
}

// FromModel converte models.DocumentoVenda para DocumentoVendaResponse
func (r *DocumentoVendaHistoricoResponse) FromModel(ddvh *models.DocumentoVendaHistorico) (*DocumentoVendaHistoricoResponse, error) {
	if ddvh == nil {
		return nil, nil
	}

	// 1. Usar o mapper para copiar campos
	if err := utils.MapToDTO(ddvh, r); err != nil {
		return nil, err
	}

	return r, nil
}

// GetPageOrDefault retorna a página ou 1 se não definido
func (f *DocumentoVendaHistoricoFilter) GetPageOrDefault() int {
	if f.Page <= 0 {
		return 1
	}
	return f.Page
}

// GetLimitOrDefault retorna o limite ou 20 se não definido
func (f *DocumentoVendaHistoricoFilter) GetLimitOrDefault() int {
	if f.Limit <= 0 {
		return 20
	}
	if f.Limit > 100 {
		return 100
	}
	return f.Limit
}

// GetOffset retorna o offset para paginação
func (f *DocumentoVendaHistoricoFilter) GetOffset() int {
	return (f.GetPageOrDefault() - 1) * f.GetLimitOrDefault()
}

// GetSortOrDefault retorna o campo de ordenação ou "ddvh_datahistorico" se não definido
func (f *DocumentoVendaHistoricoFilter) GetSortOrDefault() string {
	if f.Sort == "" {
		return "ddvh_datahistorico"
	}
	return f.Sort
}

// GetOrderOrDefault retorna a direção de ordenação ou "DESC" se não definido
func (f *DocumentoVendaHistoricoFilter) GetOrderOrDefault() string {
	if f.Order == "" {
		return "DESC"
	}
	if f.Order != "ASC" && f.Order != "DESC" {
		return "DESC"
	}
	return f.Order
}
