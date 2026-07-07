package dto

import (
	"time"

	"github.com/openerp/backend/internal/models"
)

// ============================================================
// REQUEST
// ============================================================

// GrupoEntidadeRequest é o DTO para criar/atualizar um grupo de entidade
type GrupoEntidadeRequest struct {
	Descricao string `json:"descricao" binding:"required,max=100"`
	CreatedBy *int   `json:"-"`
	UpdatedBy *int   `json:"-"`
}

// ============================================================
// RESPONSE
// ============================================================

// GrupoEntidadeResponse é o DTO para retornar um grupo de entidade
type GrupoEntidadeResponse struct {
	ID        int       `json:"id"`
	Descricao string    `json:"descricao"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FromModel converte o model para o DTO de resposta
func (r *GrupoEntidadeResponse) FromModel(m *models.GrupoEntidade) {
	r.ID = m.ID
	r.Descricao = m.Descricao
	r.CreatedAt = m.CreatedAt
	r.UpdatedAt = m.UpdatedAt
}

// ============================================================
// LIST RESPONSE
// ============================================================

// GrupoEntidadeListResponse é o DTO para retornar uma lista paginada
type GrupoEntidadeListResponse struct {
	Items      []GrupoEntidadeResponse `json:"items"`
	Total      int64                   `json:"total"`
	Page       int                     `json:"page"`
	Limit      int                     `json:"limit"`
	TotalPages int                     `json:"total_pages"`
}
