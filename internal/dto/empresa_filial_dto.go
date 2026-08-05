package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// REQUESTS
// ============================================================

// EmpresaFilialRequest representa a requisição para criar/atualizar uma filial de empresa.
type EmpresaFilialRequest struct {
	EmpresaID        int     `json:"empresa_id" binding:"required"`
	EntidadeID       *int    `json:"entidade_id,omitempty"`
	Numero           int     `json:"numero" binding:"required"`
	Nome             string  `json:"nome" binding:"required,max=100"`
	AliquotaPIS      float64 `json:"aliquota_pis" binding:"required,min=0,max=99.99"`
	AliquotaCOFINS   float64 `json:"aliquota_cofins" binding:"required,min=0,max=99.99"`
	AliquotaISS      *float64 `json:"aliquota_iss,omitempty" binding:"omitempty,min=0,max=99.99"`
	AliquotaFunrural *float64 `json:"aliquota_funrural,omitempty" binding:"omitempty,min=0,max=99.99"`
	CodigoCNAE       *string `json:"codigo_cnae,omitempty" binding:"omitempty,max=20"`
	Mei              int8    `json:"mei,omitempty"` // 0 ou 1
	LogomarcaWeb     *string `json:"logomarca_web,omitempty" binding:"omitempty,max=1000"`
	EnderecoWeb      *string `json:"endereco_web,omitempty" binding:"omitempty,max=1000"`
	CreatedBy        *int    `json:"created_by,omitempty"`
	UpdatedBy        *int    `json:"updated_by,omitempty"`
}

// ============================================================
// RESPONSES
// ============================================================

// EmpresaFilialResponse representa a resposta de uma filial de empresa.
type EmpresaFilialResponse struct {
	ID               int      `json:"id"`
	EmpresaID        int      `json:"empresa_id"`
	EmpresaNome      string   `json:"empresa_nome,omitempty"`
	EntidadeID       *int     `json:"entidade_id,omitempty"`
	Numero           int      `json:"numero"`
	Nome             string   `json:"nome"`
	AliquotaPIS      float64  `json:"aliquota_pis"`
	AliquotaCOFINS   float64  `json:"aliquota_cofins"`
	AliquotaISS      *float64 `json:"aliquota_iss,omitempty"`
	AliquotaFunrural *float64 `json:"aliquota_funrural,omitempty"`
	CodigoCNAE       *string  `json:"codigo_cnae,omitempty"`
	Mei              int8     `json:"mei"`
	MeiLabel         string   `json:"mei_label"`
	LogomarcaWeb     *string  `json:"logomarca_web,omitempty"`
	EnderecoWeb      *string  `json:"endereco_web,omitempty"`
	CreatedAt        string   `json:"created_at"`
	UpdatedAt        string   `json:"updated_at"`
	CreatedBy        *int     `json:"created_by,omitempty"`
	UpdatedBy        *int     `json:"updated_by,omitempty"`
	IsDeleted        bool     `json:"is_deleted"`
}

// ============================================================
// RESPONSE LISTA
// ============================================================

// EmpresaFilialListResponse representa a resposta de listagem de filiais de empresa.
type EmpresaFilialListResponse struct {
	Items      []EmpresaFilialResponse `json:"items"`
	Total      int64                   `json:"total"`
	Page       int                     `json:"page"`
	Limit      int                     `json:"limit"`
	TotalPages int                     `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE CONVERSÃO
// ============================================================

// ToModel converte EmpresaFilialRequest para models.EmpresaFilial.
func (r *EmpresaFilialRequest) ToModel() (*models.EmpresaFilial, error) {
	if r == nil {
		return nil, nil
	}
	filial := &models.EmpresaFilial{}
	if err := utils.MapToModel(r, filial); err != nil {
		return nil, fmt.Errorf("erro ao mapear DTO para modelo de filial de empresa: %w", err)
	}
	return filial, nil
}

// FromModel converte models.EmpresaFilial para EmpresaFilialResponse.
func (r *EmpresaFilialResponse) FromModel(filial *models.EmpresaFilial) {
	if filial == nil {
		return
	}
	// Tenta mapear diretamente, se falhar, usa fallback manual
	if err := utils.MapToDTO(filial, r); err != nil {
		r.fromModelFallback(filial)
	}

	r.MeiLabel = "Não"
	if filial.IsMEI() {
		r.MeiLabel = "Sim"
	}

	if filial.Empresa != nil {
		r.EmpresaNome = filial.Empresa.Nome
	}

	r.CreatedAt = utils.FormatDateTime(filial.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(filial.UpdatedAt)
	r.IsDeleted = filial.IsDeleted()
}

// fromModelFallback é o fallback manual caso o mapper falhe
func (r *EmpresaFilialResponse) fromModelFallback(filial *models.EmpresaFilial) {
	r.ID = filial.ID
	r.EmpresaID = filial.EmpresaID
	r.EntidadeID = filial.EntidadeID
	r.Numero = filial.Numero
	r.Nome = filial.Nome
	r.AliquotaPIS = filial.AliquotaPIS
	r.AliquotaCOFINS = filial.AliquotaCOFINS
	r.AliquotaISS = filial.AliquotaISS
	r.AliquotaFunrural = filial.AliquotaFunrural
	r.CodigoCNAE = filial.CodigoCNAE
	r.Mei = filial.Mei
	r.LogomarcaWeb = filial.LogomarcaWeb
	r.EnderecoWeb = filial.EnderecoWeb
	r.CreatedBy = filial.CreatedBy
	r.UpdatedBy = filial.UpdatedBy
	r.FromModel(filial) // Reutiliza a lógica principal para labels e datas
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o EmpresaFilialRequest.
func (r *EmpresaFilialRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return apperrors.NewValidationError(err.Error())
	}
	return nil
}