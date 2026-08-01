package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// REQUESTS
// ============================================================

// ProdutoTamanhoRequest representa a requisição para criar/atualizar um tamanho de produto.
type ProdutoTamanhoRequest struct {
	ID              int    `json:"id"`
	EmpresaFilialID int    `json:"empresa_filial_id" binding:"required"`
	Descricao       string `json:"descricao" binding:"required,max=255"`
	Sigla           string `json:"sigla,omitempty" binding:"max=10"`
	CreatedBy       *int   `json:"created_by,omitempty"`
	UpdatedBy       *int   `json:"updated_by,omitempty"`
}

// ============================================================
// RESPONSES
// ============================================================

// ProdutoTamanhoResponse representa a resposta de um tamanho de produto.
type ProdutoTamanhoResponse struct {
	ID                int    `json:"id"`
	EmpresaFilialID   int    `json:"empresa_filial_id"`
	EmpresaFilialNome string `json:"empresa_filial_nome,omitempty"`
	Descricao         string `json:"descricao"`
	Sigla             string `json:"sigla,omitempty"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
	CreatedBy         *int   `json:"created_by,omitempty"`
	UpdatedBy         *int   `json:"updated_by,omitempty"`
}

// ProdutoTamanhoListResponse representa a resposta de listagem de tamanhos de produto.
type ProdutoTamanhoListResponse struct {
	Items      []ProdutoTamanhoResponse `json:"items"`
	Total      int64                    `json:"total"`
	Page       int                      `json:"page"`
	Limit      int                      `json:"limit"`
	TotalPages int                      `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE CONVERSÃO
// ============================================================

// ToModel converte ProdutoTamanhoRequest para models.ProdutoTamanho.
func (r *ProdutoTamanhoRequest) ToModel() (*models.ProdutoTamanho, error) {
	if r == nil {
		return nil, nil
	}
	tamanho := &models.ProdutoTamanho{}
	if err := utils.MapToModel(r, tamanho); err != nil {
		return nil, fmt.Errorf("erro ao mapear DTO para modelo de tamanho de produto: %w", err)
	}

	return tamanho, nil
}

// FromModel converte models.ProdutoTamanho para ProdutoTamanhoResponse.
func (r *ProdutoTamanhoResponse) FromModel(tamanho *models.ProdutoTamanho) {
	if tamanho == nil {
		return
	}
	_ = utils.MapToDTO(tamanho, r) // Ignora erro por enquanto, assume que o mapeamento direto é suficiente

	r.CreatedAt = utils.FormatDateTime(tamanho.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(tamanho.UpdatedAt)
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o ProdutoTamanhoRequest.
func (r *ProdutoTamanhoRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
