package dto

import (
	"time"

	"github.com/openerp/backend/internal/models"
)

// ============================================================
// REQUEST
// ============================================================

// EntidadeLimiteCreditoRequest é o DTO para criar/atualizar um limite de crédito.
type EntidadeLimiteCreditoRequest struct {
	Descricao string  `json:"descricao" binding:"required,max=100"`
	Valor     float64 `json:"valor" binding:"required,gt=0"`
	CreatedBy *int    `json:"-"`
	UpdatedBy *int    `json:"-"`
}
// ============================================================
// RESPONSE
// ============================================================

// EntidadeLimiteCreditoResponse é o DTO para retornar um limite de crédito.
type EntidadeLimiteCreditoResponse struct {
	ID        int       `json:"id"`
	Descricao string    `json:"descricao"`
	Valor     float64   `json:"valor"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FromModel converte o model para o DTO de resposta.
func (r *EntidadeLimiteCreditoResponse) FromModel(m *models.EntidadeLimiteCredito) {
	r.ID = m.EntidadeID	
	r.Descricao = *m.Descricao
	r.Valor = m.Valor
	r.CreatedAt = m.CreatedAt
	r.UpdatedAt = m.UpdatedAt
}

// ============================================================
// LIST RESPONSE
// ============================================================

// EntidadeLimiteCreditoListResponse é o DTO para retornar uma lista paginada.
type EntidadeLimiteCreditoListResponse struct {
	Items      []EntidadeLimiteCreditoResponse `json:"items"`
	Total      int64                           `json:"total"`
	Page       int                             `json:"page"`
	Limit      int                             `json:"limit"`
	TotalPages int                             `json:"total_pages"`
}