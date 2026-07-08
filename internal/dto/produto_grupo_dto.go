// ============================================================
// FILE: produto_grupo_dto.go
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

// ProdutoGrupoRequest representa a requisição para criar/atualizar um grupo de produto
type ProdutoGrupoRequest struct {
	// ============================================================
	// DADOS PRINCIPAIS
	// ============================================================
	Descricao          string   `json:"descricao" binding:"required"`
	Situacao           int      `json:"situacao,omitempty"`
	ImpressoraID       *int     `json:"impressora_id,omitempty"`
	Cor                *string  `json:"cor,omitempty"`
	PercentualComissao *float64 `json:"percentual_comissao,omitempty"`
	VisivelFrenteCaixa int      `json:"visivel_frente_caixa,omitempty"`
	Agenda             *int     `json:"agenda,omitempty"`
	ControleLote       *int     `json:"controle_lote,omitempty"`

	// ============================================================
	// USUÁRIO (para auditoria)
	// ============================================================
	CreatedBy *int `json:"created_by,omitempty"`
	UpdatedBy *int `json:"updated_by,omitempty"`
}

// ============================================================
// RESPONSES
// ============================================================

// ProdutoGrupoResponse representa a resposta de um grupo de produto
type ProdutoGrupoResponse struct {
	// ============================================================
	// DADOS PRINCIPAIS
	// ============================================================
	ID                      int      `json:"id"`
	Descricao               string   `json:"descricao"`
	Situacao                int      `json:"situacao"`
	SituacaoLabel           string   `json:"situacao_label"`
	ImpressoraID            *int     `json:"impressora_id,omitempty"`
	Cor                     *string  `json:"cor,omitempty"`
	PercentualComissao      *float64 `json:"percentual_comissao,omitempty"`
	VisivelFrenteCaixa      int      `json:"visivel_frente_caixa"`
	VisivelFrenteCaixaLabel string   `json:"visivel_frente_caixa_label,omitempty"`
	Agenda                  *int     `json:"agenda,omitempty"`
	AgendaLabel             string   `json:"agenda_label,omitempty"`
	ControleLote            *int     `json:"controle_lote,omitempty"`
	ControleLoteLabel       string   `json:"controle_lote_label,omitempty"`

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

// ProdutoGrupoListResponse representa a resposta de listagem de grupos de produto
type ProdutoGrupoListResponse struct {
	Items      []ProdutoGrupoResponse `json:"items"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	TotalPages int                    `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE CONVERSÃO (USANDO MAPPER)
// ============================================================

// ToModel converte ProdutoGrupoRequest para models.ProdutoGrupo usando mapper
func (r *ProdutoGrupoRequest) ToModel() (*models.ProdutoGrupo, error) {
	if r == nil {
		return nil, nil
	}

	grupo := &models.ProdutoGrupo{}

	// 1. Usar o mapper para copiar campos (mapeamento automático)
	if err := utils.MapToModel(r, grupo); err != nil {
		return nil, err
	}

	// 2. Tratamentos especiais que o mapper não cobre
	// Se situação não informada, definir como ativo (1)
	if r.Situacao == 0 {
		grupo.Situacao = 1 // StatusAtivo
	}

	// Se visivel_frente_caixa não informado, definir como 0 (Não)
	if r.VisivelFrenteCaixa == 0 {
		grupo.VisivelFrenteCaixa = 0
	}

	return grupo, nil
}

// FromModel converte models.ProdutoGrupo para ProdutoGrupoResponse usando mapper
func (r *ProdutoGrupoResponse) FromModel(grupo *models.ProdutoGrupo) *ProdutoGrupoResponse {
	if grupo == nil {
		return nil
	}

	// 1. Usar o mapper para copiar campos (mapeamento automático)
	if err := utils.MapToDTO(grupo, r); err != nil {
		// Se o mapper falhar, usar fallback manual
		return r.fromModelFallback(grupo)
	}

	// 2. Preencher campos calculados (labels)
	r.SituacaoLabel = constants.Status(grupo.Situacao).String()
	r.VisivelFrenteCaixaLabel = constants.SimNao(grupo.VisivelFrenteCaixa).String()

	// Labels para campos do tipo *int (booleanos)
	if grupo.Agenda != nil && *grupo.Agenda != 0 {
		r.AgendaLabel = constants.SimNao(*grupo.Agenda).String()
	} else {
		r.AgendaLabel = "Não"
	}

	if grupo.ControleLote != nil {
		r.ControleLoteLabel = constants.SimNao(*grupo.ControleLote).String()
	} else {
		r.ControleLoteLabel = "Não"
	}

	// 3. Formatar datas (o mapper não faz isso)
	r.CreatedAt = utils.FormatDateTime(grupo.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(grupo.UpdatedAt)

	return r
}

// ============================================================
// FALLBACK (caso o mapper falhe)
// ============================================================

// fromModelFallback é o fallback manual caso o mapper falhe
func (r *ProdutoGrupoResponse) fromModelFallback(grupo *models.ProdutoGrupo) *ProdutoGrupoResponse {
	// Mapeamento manual campo por campo (seguro)
	r.ID = grupo.ID
	r.Descricao = grupo.Descricao
	r.Situacao = grupo.Situacao
	r.SituacaoLabel = constants.Status(grupo.Situacao).String()
	r.ImpressoraID = grupo.ImpressoraID
	r.Cor = grupo.Cor
	r.PercentualComissao = grupo.PercentualComissao
	r.VisivelFrenteCaixa = grupo.VisivelFrenteCaixa
	r.VisivelFrenteCaixaLabel = constants.SimNao(grupo.VisivelFrenteCaixa).String()
	r.Agenda = grupo.Agenda
	r.ControleLote = grupo.ControleLote
	r.CreatedBy = grupo.CreatedBy
	r.UpdatedBy = grupo.UpdatedBy

	// Labels para campos *int
	if grupo.Agenda != nil {
		r.AgendaLabel = constants.SimNao(*grupo.Agenda).String()
	} else {
		r.AgendaLabel = "Não"
	}

	if grupo.ControleLote != nil {
		r.ControleLoteLabel = constants.SimNao(*grupo.ControleLote).String()
	} else {
		r.ControleLoteLabel = "Não"
	}

	// Datas
	r.CreatedAt = utils.FormatDateTime(grupo.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(grupo.UpdatedAt)

	return r
}

// ============================================================
// FUNÇÕES AUXILIARES
// ============================================================

// getSituacaoLabel retorna o label da situação
func getSituacaoLabel(valor int) string {
	switch valor {
	case 1:
		return "Ativo"
	case 2:
		return "Inativo"
	case 3:
		return "Bloqueado"
	case 9:
		return "Cancelado"
	default:
		return "Desconhecido"
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o ProdutoGrupoRequest
func (r *ProdutoGrupoRequest) Validate() error {
	validate := validator.New()

	// Validações básicas
	if err := validate.Struct(r); err != nil {
		return err
	}

	// Validações customizadas
	if err := constants.Status(r.Situacao).IsValid(); err != nil {
		return err
	}

	if err := constants.SimNao(r.VisivelFrenteCaixa).IsValid(); err != nil {
		return err
	}

	if r.Agenda != nil {
		if err := constants.SimNao(*r.Agenda).IsValid(); err != nil {
			return err
		}
	}

	if r.ControleLote != nil {
		if err := constants.SimNao(*r.ControleLote).IsValid(); err != nil {
			return err
		}
	}
	return nil
}
