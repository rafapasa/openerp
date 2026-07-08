package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/openerp/backend/internal/constants"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ProdutoMarcaRequest representa a requisição para criar/atualizar uma marca de produto.
type ProdutoMarcaRequest struct {
	Descricao string `json:"descricao" binding:"required,max=255"`
	Situacao  int    `json:"situacao,omitempty"`
	CreatedBy *int   `json:"-"`
	UpdatedBy *int   `json:"-"`
}

// ProdutoMarcaResponse representa a resposta de uma marca de produto.
type ProdutoMarcaResponse struct {
	ID            int    `json:"id"`
	Descricao     string `json:"descricao"`
	Situacao      int    `json:"situacao"`
	SituacaoLabel string `json:"situacao_label"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	CreatedBy     *int   `json:"created_by,omitempty"`
	UpdatedBy     *int   `json:"updated_by,omitempty"`
}

// ProdutoMarcaListResponse representa a resposta de listagem de marcas.
type ProdutoMarcaListResponse struct {
	Items      []ProdutoMarcaResponse `json:"items"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	TotalPages int                    `json:"total_pages"`
}

// ToModel converte ProdutoMarcaRequest para models.ProdutoMarca.
func (r *ProdutoMarcaRequest) ToModel() (*models.ProdutoMarca, error) {
	if r == nil {
		return nil, nil
	}
	marca := &models.ProdutoMarca{}
	if err := utils.MapToModel(r, marca); err != nil {
		return nil, err
	}
	if r.Situacao == 0 {
		marca.Situacao = int(constants.StatusAtivo)
	}
	return marca, nil
}

// FromModel converte models.ProdutoMarca para ProdutoMarcaResponse.
func (r *ProdutoMarcaResponse) FromModel(marca *models.ProdutoMarca) {
	if marca == nil {
		return
	}
	r.ID = marca.ID
	r.Descricao = marca.Descricao
	r.Situacao = marca.Situacao
	r.SituacaoLabel = constants.Status(marca.Situacao).String()
	r.CreatedAt = utils.FormatDateTime(marca.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(marca.UpdatedAt)
	r.CreatedBy = marca.CreatedBy
	r.UpdatedBy = marca.UpdatedBy
}

// Validate valida o ProdutoMarcaRequest.
func (r *ProdutoMarcaRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}

	if err := constants.Status(r.Situacao).IsValid(); err != nil {
		return err
	}

	return nil
}
