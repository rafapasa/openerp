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

// EntidadeRegimeTributarioRequest representa a requisição para criar/atualizar um regime tributário da entidade.
type EntidadeRegimeTributarioRequest struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	EntidadeID      int    `json:"entidade_id" binding:"required"`
	Regime          int    `json:"regime" binding:"required"`
	Apuracao        *int   `json:"apuracao,omitempty"`
	Data            string `json:"data" binding:"required"` // Formato: 2006-01-02
	RegimeEspecial  int    `json:"regime_especial"`
	SituacaoTribISS *int   `json:"situacao_trib_iss,omitempty"`
	RegimeMunicipal *int   `json:"regime_municipal,omitempty"`

	// ============================================================
	// AUDITORIA
	// ============================================================
	CreatedBy *int `json:"created_by,omitempty"`
	UpdatedBy *int `json:"updated_by,omitempty"`
}

// ============================================================
// RESPONSES
// ============================================================

// EntidadeRegimeTributarioResponse representa a resposta de um regime tributário da entidade.
type EntidadeRegimeTributarioResponse struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID                   int    `json:"id"`
	EntidadeID           int    `json:"entidade_id"`
	Regime               int    `json:"regime"`
	RegimeLabel          string `json:"regime_label"`
	Apuracao             *int   `json:"apuracao,omitempty"`
	ApuracaoLabel        string `json:"apuracao_label,omitempty"`
	Data                 string `json:"data"`
	RegimeEspecial       int    `json:"regime_especial"`
	RegimeEspecialLabel  string `json:"regime_especial_label"`
	SituacaoTribISS      *int   `json:"situacao_trib_iss,omitempty"`
	SituacaoTribISSLabel string `json:"situacao_trib_iss_label,omitempty"`
	RegimeMunicipal      *int   `json:"regime_municipal,omitempty"`
	RegimeMunicipalLabel string `json:"regime_municipal_label,omitempty"`

	// ============================================================
	// AUDITORIA
	// ============================================================
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	CreatedBy *int   `json:"created_by,omitempty"`
	UpdatedBy *int   `json:"updated_by,omitempty"`
}

// ============================================================
// LIST RESPONSE
// ============================================================

// EntidadeRegimeTributarioListResponse representa a resposta de listagem de regimes.
type EntidadeRegimeTributarioListResponse struct {
	Items      []EntidadeRegimeTributarioResponse `json:"items"`
	Total      int64                              `json:"total"`
	Page       int                                `json:"page"`
	Limit      int                                `json:"limit"`
	TotalPages int                                `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE CONVERSÃO
// ============================================================

// ToModel converte EntidadeRegimeTributarioRequest para models.EntidadeRegimeTributario.
func (r *EntidadeRegimeTributarioRequest) ToModel() (*models.EntidadeRegimeTributario, error) {
	if r == nil {
		return nil, nil
	}

	model := &models.EntidadeRegimeTributario{}
	if err := utils.MapToModel(r, model); err != nil {
		return nil, err
	}

	// Converter data
	if data, err := utils.ParseDate(r.Data); err == nil {
		model.Data = data
	} else {
		return nil, fmt.Errorf("formato de data inválido: %w", err)
	}

	return model, nil
}

// FromModel converte models.EntidadeRegimeTributario para EntidadeRegimeTributarioResponse.
func (r *EntidadeRegimeTributarioResponse) FromModel(model *models.EntidadeRegimeTributario) *EntidadeRegimeTributarioResponse {
	if model == nil {
		return nil
	}

	if err := utils.MapToDTO(model, r); err != nil {
		// Fallback manual em caso de erro no mapper
		r.ID = model.ID
		r.EntidadeID = model.EntidadeID
		r.Regime = model.Regime
		r.Apuracao = model.Apuracao
		r.RegimeEspecial = model.RegimeEspecial
		r.SituacaoTribISS = model.SituacaoTribISS
		r.RegimeMunicipal = model.RegimeMunicipal
		r.CreatedBy = model.CreatedBy
		r.UpdatedBy = model.UpdatedBy
	}

	// Formatar datas
	r.Data = utils.FormatDate(model.Data)
	r.CreatedAt = utils.FormatDateTime(model.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(model.UpdatedAt)

	// Preencher labels
	r.RegimeLabel = getRegimeLabel(model.Regime)
	r.RegimeEspecialLabel = getSimNaoLabel(model.RegimeEspecial)
	if model.Apuracao != nil {
		r.ApuracaoLabel = getApuracaoLabel(*model.Apuracao)
	}
	if model.SituacaoTribISS != nil {
		r.SituacaoTribISSLabel = getSituacaoTribISSLabel(*model.SituacaoTribISS)
	}
	if model.RegimeMunicipal != nil {
		r.RegimeMunicipalLabel = getRegimeMunicipalLabel(*model.RegimeMunicipal)
	}

	return r
}

// ============================================================
// FUNÇÕES AUXILIARES
// ============================================================

func getRegimeLabel(regime int) string {
	switch regime {
	case 1:
		return "Simples Nacional"
	case 2:
		return "Lucro Presumido"
	case 3:
		return "Lucro Real"
	case 4:
		return "MEI"
	case 5:
		return "Isento"
	default:
		return "Desconhecido"
	}
}

func getApuracaoLabel(apuracao int) string {
	// Implementar a lógica para os labels de apuração
	return "Não definido"
}

func getSituacaoTribISSLabel(situacao int) string {
	// Implementar a lógica para os labels de situação do ISS
	return "Não definido"
}

func getRegimeMunicipalLabel(regime int) string {
	// Implementar a lógica para os labels de regime municipal
	return "Não definido"
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o EntidadeRegimeTributarioRequest.
func (r *EntidadeRegimeTributarioRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
