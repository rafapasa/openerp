package dto

import (
	"fmt"
	"time"

	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
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

// ToModel converte EntidadeLimiteCreditoRequest para models.EntidadeLimiteCredito.
func (r *EntidadeLimiteCreditoRequest) ToModel() (*models.EntidadeLimiteCredito, error) {
	if r == nil {
		return nil, nil
	}

	limite := &models.EntidadeLimiteCredito{}
	if err := utils.MapToModel(r, limite); err != nil {
		return nil, fmt.Errorf("erro ao mapear DTO para modelo de limite de crédito: %w", err)
	}

	// Definir a data atual se não for fornecida (ou se for um campo obrigatório no modelo)
	if limite.Data.IsZero() {
		limite.Data = time.Now()
	}
	return limite, nil
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
