package dto

import (
	"strings"
	"time"

	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ProdutoCorRequest representa a estrutura de dados para criar ou atualizar uma cor de produto.
type ProdutoCorRequest struct {
	ID              int    `json:"id"`
	EmpresaFilialID int    `json:"empresa_filial_id" validate:"required,gt=0"`
	Sigla           string `json:"sigla" validate:"required,min=1,max=20"`
	Nome            string `json:"nome" validate:"required,min=3,max=255"`
	CreatedBy       *int   `json:"created_by,omitempty"`
	UpdatedBy       *int   `json:"updated_by,omitempty"`
}

// Validate valida os campos da requisição.
func (r *ProdutoCorRequest) Validate() error {
	if err := utils.ValidateMandatoryFields(r); err != nil {
		return err
	}
	r.Sigla = strings.TrimSpace(r.Sigla)
	r.Nome = strings.TrimSpace(r.Nome)
	return nil
}

// ToModel converte um ProdutoCorRequest para um modelo ProdutoCor.
func (r *ProdutoCorRequest) ToModel() (*models.ProdutoCor, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	return &models.ProdutoCor{
		ID:              r.ID,
		EmpresaFilialID: r.EmpresaFilialID,
		Sigla:           r.Sigla,
		Nome:            r.Nome,
		CreatedBy:       r.CreatedBy,
		UpdatedBy:       r.UpdatedBy,
	}, nil
}

// ProdutoCorResponse representa a estrutura de dados de resposta para uma cor de produto.
type ProdutoCorResponse struct {
	ID                int        `json:"id"`
	EmpresaFilialID   int        `json:"empresa_filial_id"`
	EmpresaFilialNome string     `json:"empresa_filial_nome,omitempty"`
	Sigla             string     `json:"sigla"`
	Nome              string     `json:"nome"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	CreatedBy         *int       `json:"created_by,omitempty"`
	UpdatedBy         *int       `json:"updated_by,omitempty"`
}