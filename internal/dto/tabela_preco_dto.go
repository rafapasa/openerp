package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// TabelaPrecoRequest representa a requisição para criar/atualizar uma tabela de preço.
type TabelaPrecoRequest struct {
	Descricao         string   `json:"descricao" binding:"required,max=100"`
	DataInicio        string   `json:"data_inicio" binding:"required"` // Formato: YYYY-MM-DD
	DataFim           *string  `json:"data_fim,omitempty"`             // Formato: YYYY-MM-DD
	Tipo              *int     `json:"tipo,omitempty"`
	Percentual        *float64 `json:"percentual,omitempty"`
	TipoServico       *int     `json:"tipo_servico,omitempty"`
	PercentualServico *float64 `json:"percentual_servico,omitempty"`
	CreatedBy         *int     `json:"-"`
	UpdatedBy         *int     `json:"-"`
}

// TabelaPrecoResponse representa a resposta de uma tabela de preço.
type TabelaPrecoResponse struct {
	ID                int      `json:"id"`
	Descricao         string   `json:"descricao"`
	DataInicio        string   `json:"data_inicio"`
	DataFim           *string  `json:"data_fim,omitempty"`
	Tipo              *int     `json:"tipo,omitempty"`
	Percentual        *float64 `json:"percentual,omitempty"`
	TipoServico       *int     `json:"tipo_servico,omitempty"`
	PercentualServico *float64 `json:"percentual_servico,omitempty"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
	CreatedBy         *int     `json:"created_by,omitempty"`
	UpdatedBy         *int     `json:"updated_by,omitempty"`
}

// TabelaPrecoListResponse representa a resposta de listagem de tabelas de preço.
type TabelaPrecoListResponse struct {
	Items      []TabelaPrecoResponse `json:"items"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalPages int                   `json:"total_pages"`
}

// ToModel converte TabelaPrecoRequest para models.TabelaPreco.
func (r *TabelaPrecoRequest) ToModel() (*models.TabelaPreco, error) {
	if r == nil {
		return nil, nil
	}
	tabela := &models.TabelaPreco{}
	if err := utils.MapToModel(r, tabela); err != nil {
		return nil, err
	}

	// Converter datas
	if data, err := utils.ParseDate(r.DataInicio); err == nil {
		tabela.DataInicio = data
	}

	if r.DataFim != nil {
		if data, err := utils.ParseDate(*r.DataFim); err == nil {
			tabela.DataFim = &data
		}
	}

	return tabela, nil
}

// FromModel converte models.TabelaPreco para TabelaPrecoResponse.
func (r *TabelaPrecoResponse) FromModel(tabela *models.TabelaPreco) {
	if tabela == nil {
		return
	}

	if err := utils.MapToDTO(tabela, r); err != nil {
		// Fallback manual
		r.ID = tabela.ID
		r.Descricao = tabela.Descricao
		r.Tipo = tabela.Tipo
		r.Percentual = tabela.Percentual
		r.TipoServico = tabela.TipoServico
		r.PercentualServico = tabela.PercentualServico
		r.CreatedBy = tabela.CreatedBy
		r.UpdatedBy = tabela.UpdatedBy
	}

	// Formatar datas
	r.DataInicio = utils.FormatDate(tabela.DataInicio)
	if tabela.DataFim != nil {
		r.DataFim = utils.StringPtr(utils.FormatDate(*tabela.DataFim))
	}
	r.CreatedAt = utils.FormatDateTime(tabela.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(tabela.UpdatedAt)
}

// Validate valida o TabelaPrecoRequest.
func (r *TabelaPrecoRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
