// ============================================================
// FILE: produto_serie_dto.go
// PACKAGE: dto
// ============================================================

package dto

import (
	"github.com/go-playground/validator/v10"

	"github.com/openerp/backend/internal/constants"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// REQUESTS
// ============================================================

// ProdutoSerieRequest representa a requisição para criar/atualizar uma série de produto
type ProdutoSerieRequest struct {
	// ============================================================
	// DADOS PRINCIPAIS
	// ============================================================
	Descricao string           `json:"descricao" binding:"required"`
	Situacao  constants.Status `json:"situacao,omitempty"`

	// ============================================================
	// USUÁRIO (para auditoria)
	// ============================================================
	CreatedBy *int `json:"created_by,omitempty"`
	UpdatedBy *int `json:"updated_by,omitempty"`
}

// ============================================================
// RESPONSES
// ============================================================

// ProdutoSerieResponse representa a resposta de uma série de produto
type ProdutoSerieResponse struct {
	// ============================================================
	// DADOS PRINCIPAIS
	// ============================================================
	ID            int              `json:"id"`
	Descricao     string           `json:"descricao"`
	Situacao      constants.Status `json:"situacao"`
	SituacaoLabel string           `json:"situacao_label"`

	// ============================================================
	// AUDITORIA
	// ============================================================
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	CreatedBy *int   `json:"created_by,omitempty"`
	UpdatedBy *int   `json:"updated_by,omitempty"`
}

// ============================================================
// RESPONSE LISTA
// ============================================================

// ProdutoSerieListResponse representa a resposta de listagem de séries de produto
type ProdutoSerieListResponse struct {
	Items      []ProdutoSerieResponse `json:"items"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	TotalPages int                    `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE CONVERSÃO (USANDO MAPPER)
// ============================================================

// ToModel converte ProdutoSerieRequest para models.ProdutoSerie usando mapper
func (r *ProdutoSerieRequest) ToModel() (*models.ProdutoSerie, error) {
	if r == nil {
		return nil, nil
	}

	serie := &models.ProdutoSerie{}

	// 1. Usar o mapper para copiar campos (mapeamento automático)
	if err := utils.MapToModel(r, serie); err != nil {
		return nil, err
	}

	// 2. Tratamentos especiais que o mapper não cobre
	// Converter Situacao (se não informada, definir como ativo = 1)
	if r.Situacao == 0 {
		serie.Situacao = 1 // StatusAtivo
	}

	return serie, nil
}

// FromModel converte models.ProdutoSerie para ProdutoSerieResponse usando mapper
func (r *ProdutoSerieResponse) FromModel(serie *models.ProdutoSerie) *ProdutoSerieResponse {
	if serie == nil {
		return nil
	}

	// 1. Usar o mapper para copiar campos (mapeamento automático)
	if err := utils.MapToDTO(serie, r); err != nil {
		// Se o mapper falhar, usar fallback manual
		return r.fromModelFallback(serie)
	}

	// 2. Preencher campos calculados (labels)
	r.SituacaoLabel = serie.Situacao.String()

	// 3. Formatar datas (o mapper não faz isso)
	r.CreatedAt = utils.FormatDateTime(serie.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(serie.UpdatedAt)

	return r
}

// ============================================================
// FALLBACK (caso o mapper falhe)
// ============================================================

// fromModelFallback é o fallback manual caso o mapper falhe
func (r *ProdutoSerieResponse) fromModelFallback(serie *models.ProdutoSerie) *ProdutoSerieResponse {
	// Mapeamento manual campo por campo (seguro)
	r.ID = serie.ID
	r.Descricao = serie.Descricao
	r.Situacao = serie.Situacao
	r.SituacaoLabel = serie.Situacao.String()
	r.CreatedBy = serie.CreatedBy
	r.UpdatedBy = serie.UpdatedBy
	r.CreatedAt = utils.FormatDateTime(serie.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(serie.UpdatedAt)

	return r
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o ProdutoSerieRequest
func (r *ProdutoSerieRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
