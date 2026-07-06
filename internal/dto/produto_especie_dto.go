// ============================================================
// FILE: produto_especie_dto.go
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

// ProdutoEspecieRequest representa a requisição para criar/atualizar uma espécie de produto
type ProdutoEspecieRequest struct {
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

// ProdutoEspecieResponse representa a resposta de uma espécie de produto
type ProdutoEspecieResponse struct {
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

// ProdutoEspecieListResponse representa a resposta de listagem de espécies de produto
type ProdutoEspecieListResponse struct {
	Items      []ProdutoEspecieResponse `json:"items"`
	Total      int64                    `json:"total"`
	Page       int                      `json:"page"`
	Limit      int                      `json:"limit"`
	TotalPages int                      `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE CONVERSÃO (USANDO MAPPER)
// ============================================================

// ToModel converte ProdutoEspecieRequest para models.ProdutoEspecie usando mapper
func (r *ProdutoEspecieRequest) ToModel() (*models.ProdutoEspecie, error) {
	if r == nil {
		return nil, nil
	}

	especie := &models.ProdutoEspecie{}

	// 1. Usar o mapper para copiar campos (mapeamento automático)
	if err := utils.MapToModel(r, especie); err != nil {
		return nil, err
	}

	// 2. Tratamentos especiais que o mapper não cobre
	// Converter Situacao (se não informada, definir como ativo = 1)
	if r.Situacao == 0 {
		especie.Situacao = 1 // StatusAtivo
	}

	return especie, nil
}

// FromModel converte models.ProdutoEspecie para ProdutoEspecieResponse usando mapper
func (r *ProdutoEspecieResponse) FromModel(especie *models.ProdutoEspecie) *ProdutoEspecieResponse {
	if especie == nil {
		return nil
	}

	// 1. Usar o mapper para copiar campos (mapeamento automático)
	if err := utils.MapToDTO(especie, r); err != nil {
		// Se o mapper falhar, usar fallback manual
		return r.fromModelFallback(especie)
	}

	// 2. Preencher campos calculados (labels)
	r.SituacaoLabel = especie.Situacao.String()

	// 3. Formatar datas (o mapper não faz isso)
	r.CreatedAt = utils.FormatDateTime(especie.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(especie.UpdatedAt)

	return r
}

// ============================================================
// FALLBACK (caso o mapper falhe)
// ============================================================

// fromModelFallback é o fallback manual caso o mapper falhe
func (r *ProdutoEspecieResponse) fromModelFallback(especie *models.ProdutoEspecie) *ProdutoEspecieResponse {
	// Mapeamento manual campo por campo (seguro)
	r.ID = especie.ID
	r.Descricao = especie.Descricao
	r.Situacao = especie.Situacao
	r.SituacaoLabel = especie.Situacao.String()
	r.CreatedBy = especie.CreatedBy
	r.UpdatedBy = especie.UpdatedBy
	r.CreatedAt = utils.FormatDateTime(especie.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(especie.UpdatedAt)

	return r
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o ProdutoEspecieRequest
func (r *ProdutoEspecieRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
