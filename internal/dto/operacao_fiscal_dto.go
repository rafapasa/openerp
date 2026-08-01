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

// OperacaoFiscalRequest representa a requisição para criar/atualizar uma operação fiscal.
type OperacaoFiscalRequest struct {
	EmpresaFilialID int     `json:"empresa_filial_id" binding:"required"`
	CFOP            string  `json:"cfop" binding:"required,len=4"`
	Descricao       string  `json:"descricao" binding:"required,max=255"`
	DataIni         string  `json:"data_ini" binding:"required"` // Formato: YYYY-MM-DD
	DataFim         *string `json:"data_fim,omitempty"`          // Formato: YYYY-MM-DD
	CSTICMSID       int     `json:"cst_icms_id" binding:"required"`
	CSTIPIID        int     `json:"cst_ipi_id" binding:"required"`
	CSTPISCOFINSID  int     `json:"cst_pis_cofins_id" binding:"required"`
	CreatedBy       *int    `json:"created_by,omitempty"`
	UpdatedBy       *int    `json:"updated_by,omitempty"`
}

// ============================================================
// RESPONSES
// ============================================================

// OperacaoFiscalResponse representa a resposta de uma operação fiscal.
type OperacaoFiscalResponse struct {
	ID                 int    `json:"id"`
	EmpresaFilialID    int    `json:"empresa_filial_id"`
	EmpresaFilialNome  string `json:"empresa_filial_nome,omitempty"`
	CFOP               string `json:"cfop"`
	Descricao          string `json:"descricao"`
	DataIni            string `json:"data_ini"`
	DataFim            string `json:"data_fim,omitempty"`
	CSTICMSID          int    `json:"cst_icms_id"`
	CSTICMSCodigo      string `json:"cst_icms_codigo,omitempty"`
	CSTIPIID           int    `json:"cst_ipi_id"`
	CSTIPICodigo       string `json:"cst_ipi_codigo,omitempty"`
	CSTPISCOFINSID     int    `json:"cst_pis_cofins_id"`
	CSTPISCOFINSCodigo string `json:"cst_pis_cofins_codigo,omitempty"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	CreatedBy          *int   `json:"created_by,omitempty"`
	UpdatedBy          *int   `json:"updated_by,omitempty"`
}

// OperacaoFiscalListResponse representa a resposta de listagem de operações fiscais.
type OperacaoFiscalListResponse struct {
	Items      []OperacaoFiscalResponse `json:"items"`
	Total      int64                    `json:"total"`
	Page       int                      `json:"page"`
	Limit      int                      `json:"limit"`
	TotalPages int                      `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE CONVERSÃO
// ============================================================

// ToModel converte OperacaoFiscalRequest para models.OperacaoFiscal.
func (r *OperacaoFiscalRequest) ToModel() (*models.OperacaoFiscal, error) {
	if r == nil {
		return nil, nil
	}
	operacao := &models.OperacaoFiscal{}
	if err := utils.MapToModel(r, operacao); err != nil {
		return nil, fmt.Errorf("erro ao mapear DTO para modelo de operação fiscal: %w", err)
	}

	// Converter datas
	if data, err := utils.ParseDate(r.DataIni); err == nil {
		operacao.DataIni = data
	} else {
		return nil, fmt.Errorf("data de início inválida: %w", err)
	}
	if r.DataFim != nil && *r.DataFim != "" {
		if data, err := utils.ParseDate(*r.DataFim); err == nil {
			operacao.DataFim = &data
		} else {
			return nil, fmt.Errorf("data de fim inválida: %w", err)
		}
	}
	return operacao, nil
}

// FromModel converte models.OperacaoFiscal para OperacaoFiscalResponse.
func (r *OperacaoFiscalResponse) FromModel(operacao *models.OperacaoFiscal) {
	if operacao == nil {
		return
	}
	_ = utils.MapToDTO(operacao, r) // Ignora erro por enquanto, assume que o mapeamento direto é suficiente

	r.DataIni = utils.FormatDate(operacao.DataIni)
	if operacao.DataFim != nil {
		r.DataFim = utils.FormatDate(*operacao.DataFim)
	}
	r.CreatedAt = utils.FormatDateTime(operacao.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(operacao.UpdatedAt)

	if operacao.EmpresaFilial != nil {
		r.EmpresaFilialNome = operacao.EmpresaFilial.Nome
	}
	if operacao.CSTICMS != nil {
		r.CSTICMSCodigo = operacao.CSTICMS.Codigo
	}
	if operacao.CSTIPI != nil {
		r.CSTIPICodigo = operacao.CSTIPI.Codigo
	}
	if operacao.CSTPISCOFINS != nil {
		r.CSTPISCOFINSCodigo = operacao.CSTPISCOFINS.Codigo
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o OperacaoFiscalRequest.
func (r *OperacaoFiscalRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
