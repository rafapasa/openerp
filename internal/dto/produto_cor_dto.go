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

// ProdutoCorRequest representa a requisição para criar/atualizar uma cor de produto.
type ProdutoCorRequest struct {
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

// ProdutoCorResponse representa a resposta de uma cor de produto.
type ProdutoCorResponse struct {
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

// ProdutoCorListResponse representa a resposta de listagem de cores de produto.
type ProdutoCorListResponse struct {
	Items      []ProdutoCorResponse `json:"items"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
	TotalPages int                  `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE CONVERSÃO
// ============================================================

// ToModel converte ProdutoCorRequest para models.ProdutoCor.
func (r *ProdutoCorRequest) ToModel() (*models.ProdutoCor, error) {
	if r == nil {
		return nil, nil
	}
	cor := &models.ProdutoCor{}
	if err := utils.MapToModel(r, cor); err != nil {
		return nil, fmt.Errorf("erro ao mapear DTO para modelo de cor de produto: %w", err)
	}

	return cor, nil
}

// FromModel converte models.ProdutoCor para ProdutoCorResponse.
func (r *ProdutoCorResponse) FromModel(cor *models.ProdutoCor) {
	if cor == nil {
		return
	}
	_ = utils.MapToDTO(cor, r) // Ignora erro por enquanto, assume que o mapeamento direto é suficiente

	r.CreatedAt = utils.FormatDateTime(cor.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(cor.UpdatedAt)
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o ProdutoCorRequest.
func (r *ProdutoCorRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
