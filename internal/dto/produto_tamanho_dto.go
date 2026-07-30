package dto

import (
	"strings"
	"time"

	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ProdutoTamanhoRequest representa a estrutura de dados para criar ou atualizar um tamanho de produto.
type ProdutoTamanhoRequest struct {
	ID              int    `json:"id"`
	EmpresaFilialID int    `json:"empresa_filial_id" validate:"required,gt=0"`
	Sigla           string `json:"sigla" validate:"required,min=1,max=20"`
	Nome            string `json:"nome" validate:"required,min=3,max=255"`
	CreatedBy       *int   `json:"created_by,omitempty"`
	UpdatedBy       *int   `json:"updated_by,omitempty"`
}

// Validate valida os campos da requisição.
func (r *ProdutoTamanhoRequest) Validate() error {
	if err := utils.ValidateMandatoryFields(r); err != nil {
		return err
	}
	r.Sigla = strings.TrimSpace(r.Sigla)
	r.Nome = strings.TrimSpace(r.Nome)
	return nil
}

// ToModel converte um ProdutoTamanhoRequest para um modelo ProdutoTamanho.
func (r *ProdutoTamanhoRequest) ToModel() (*models.ProdutoTamanho, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	return &models.ProdutoTamanho{
		ID:              r.ID,
		EmpresaFilialID: r.EmpresaFilialID,
		Sigla:           r.Sigla,
		Nome:            r.Nome,
		CreatedBy:       r.CreatedBy,
		UpdatedBy:       r.UpdatedBy,
	}, nil
}

// ProdutoTamanhoResponse representa a estrutura de dados de resposta para um tamanho de produto.
type ProdutoTamanhoResponse struct {
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