package dto

import (
	"github.com/openerp/backend/internal/constants"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// REQUESTS
// ============================================================

// ProdutoSubgrupoRequest representa a requisição para criar/atualizar um subgrupo de produto.
type ProdutoSubgrupoRequest struct {
	Descricao string `json:"descricao" binding:"required,max=255"`
	Situacao  int    `json:"situacao,omitempty"`
	CreatedBy *int   `json:"-"`
	UpdatedBy *int   `json:"-"`
}

// ============================================================
// RESPONSES
// ============================================================

// ProdutoSubgrupoResponse representa a resposta de um subgrupo de produto.
type ProdutoSubgrupoResponse struct {
	ID            int    `json:"id"`
	Descricao     string `json:"descricao"`
	Situacao      int    `json:"situacao"`
	SituacaoLabel string `json:"situacao_label"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	CreatedBy     *int   `json:"created_by,omitempty"`
	UpdatedBy     *int   `json:"updated_by,omitempty"`
}

// ============================================================
// LIST RESPONSE
// ============================================================

// ProdutoSubgrupoListResponse representa a resposta de listagem de subgrupos.
type ProdutoSubgrupoListResponse struct {
	Items      []ProdutoSubgrupoResponse `json:"items"`
	Total      int64                     `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE CONVERSÃO
// ============================================================

// ToModel converte ProdutoSubgrupoRequest para models.ProdutoSubgrupo.
func (r *ProdutoSubgrupoRequest) ToModel() (*models.ProdutoSubgrupo, error) {
	if r == nil {
		return nil, nil
	}

	subgrupo := &models.ProdutoSubgrupo{}
	if err := utils.MapToModel(r, subgrupo); err != nil {
		return nil, err
	}

	// Se situação não informada, definir como ativo (1).
	if r.Situacao == 0 {
		subgrupo.Situacao = int(constants.StatusAtivo)
	}

	return subgrupo, nil
}

// FromModel converte models.ProdutoSubgrupo para ProdutoSubgrupoResponse.
func (r *ProdutoSubgrupoResponse) FromModel(subgrupo *models.ProdutoSubgrupo) {
	if subgrupo == nil {
		return
	}

	r.ID = subgrupo.ID
	r.Descricao = subgrupo.Descricao
	r.Situacao = subgrupo.Situacao
	r.SituacaoLabel = constants.Status(subgrupo.Situacao).String()
	r.CreatedAt = utils.FormatDateTime(subgrupo.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(subgrupo.UpdatedAt)
	r.CreatedBy = subgrupo.CreatedBy
	r.UpdatedBy = subgrupo.UpdatedBy
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o ProdutoSubgrupoRequest.
func (r *ProdutoSubgrupoRequest) Validate() error {
	if err := utils.ValidateMandatoryFields(r); err != nil {
		return err
	}

	// Validação customizada para situação.
	if err := constants.Status(r.Situacao).IsValid(); err != nil {
		return err
	}

	return nil

}
