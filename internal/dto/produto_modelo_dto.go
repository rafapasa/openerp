package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/openerp/backend/internal/constants"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ProdutoModeloRequest representa a requisição para criar/atualizar um modelo de produto.
type ProdutoModeloRequest struct {
	Descricao string `json:"descricao" binding:"required,max=255"`
	Situacao  int    `json:"situacao,omitempty"`
	CreatedBy *int   `json:"-"`
	UpdatedBy *int   `json:"-"`
}

// ProdutoModeloResponse representa a resposta de um modelo de produto.
type ProdutoModeloResponse struct {
	ID            int    `json:"id"`
	Descricao     string `json:"descricao"`
	Situacao      int    `json:"situacao"`
	SituacaoLabel string `json:"situacao_label"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	CreatedBy     *int   `json:"created_by,omitempty"`
	UpdatedBy     *int   `json:"updated_by,omitempty"`
}

// ProdutoModeloListResponse representa a resposta de listagem de modelos.
type ProdutoModeloListResponse struct {
	Items      []ProdutoModeloResponse `json:"items"`
	Total      int64                   `json:"total"`
	Page       int                     `json:"page"`
	Limit      int                     `json:"limit"`
	TotalPages int                     `json:"total_pages"`
}

// ToModel converte ProdutoModeloRequest para models.ProdutoModelo.
func (r *ProdutoModeloRequest) ToModel() (*models.ProdutoModelo, error) {
	if r == nil {
		return nil, nil
	}
	modelo := &models.ProdutoModelo{}
	if err := utils.MapToModel(r, modelo); err != nil {
		return nil, err
	}
	if r.Situacao == 0 {
		modelo.Situacao = int(constants.StatusAtivo)
	}
	return modelo, nil
}

// FromModel converte models.ProdutoModelo para ProdutoModeloResponse.
func (r *ProdutoModeloResponse) FromModel(modelo *models.ProdutoModelo) {
	if modelo == nil {
		return
	}
	r.ID = modelo.ID
	r.Descricao = modelo.Descricao
	r.Situacao = modelo.Situacao
	r.SituacaoLabel = constants.Status(modelo.Situacao).String()
	r.CreatedAt = utils.FormatDateTime(modelo.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(modelo.UpdatedAt)
	r.CreatedBy = modelo.CreatedBy
	r.UpdatedBy = modelo.UpdatedBy
}

// Validate valida o ProdutoModeloRequest.
func (r *ProdutoModeloRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}

	if err := validateSituacao(r.Situacao); err != nil {
		return err
	}

	return nil
}

func validateSituacao(valor int) error {
	if valor == 0 {
		return nil // Opcional, será tratado como Ativo.
	}
	if err := constants.Status(valor).IsValid(); err != nil {
		return err
	}
	return nil
}
