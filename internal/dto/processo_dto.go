package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/openerp/backend/internal/constants"
	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// REQUESTS
// ============================================================

// ProcessoRequest representa a requisição para criar/atualizar um processo.
type ProcessoRequest struct {
	EmpresaFilialID           int    `json:"empresa_filial_id" binding:"required"`
	Codigo                    int    `json:"codigo" binding:"required"`
	Descricao                 string `json:"descricao" binding:"required,max=255"`
	TipoOperacao              int    `json:"tipo_operacao" binding:"required"` // 0-Entrada, 1-Saída
	Situacao                  int    `json:"situacao,omitempty"`
	RotinaContabilID          *int   `json:"rotina_contabil_id,omitempty"`
	OperacaoFiscalNoEstID     *int   `json:"operacao_fiscal_no_est_id,omitempty"`
	OperacaoFiscalForaEstID   *int   `json:"operacao_fiscal_fora_est_id,omitempty"`
	OperacaoFiscalNoEstSTID   *int   `json:"operacao_fiscal_no_est_st_id,omitempty"`
	OperacaoFiscalForaEstSTID *int   `json:"operacao_fiscal_fora_est_st_id,omitempty"`
	PlanoContasFinanceiroID   *int   `json:"plano_contas_financeiro_id,omitempty"`
	NaturezaOperacaoID        *int   `json:"natureza_operacao_id,omitempty"`
	ReceitaID                 *int   `json:"receita_id,omitempty"`
	DespesaID                 *int   `json:"despesa_id,omitempty"`
	ProcessoNFID              *int   `json:"processo_nf_id,omitempty"`
	CreatedBy                 *int   `json:"created_by,omitempty"`
	UpdatedBy                 *int   `json:"updated_by,omitempty"`
}

// ============================================================
// RESPONSES
// ============================================================

// ProcessoResponse representa a resposta de um processo.
type ProcessoResponse struct {
	ID                        int    `json:"id"`
	EmpresaFilialID           int    `json:"empresa_filial_id"`
	EmpresaFilialNome         string `json:"empresa_filial_nome,omitempty"`
	Codigo                    int    `json:"codigo"`
	Descricao                 string `json:"descricao"`
	TipoOperacao              int    `json:"tipo_operacao"`
	TipoOperacaoLabel         string `json:"tipo_operacao_label"`
	Situacao                  int    `json:"situacao"`
	SituacaoLabel             string `json:"situacao_label"`
	RotinaContabilID          *int   `json:"rotina_contabil_id,omitempty"`
	RotinaContabilDescricao   string `json:"rotina_contabil_descricao,omitempty"`
	OperacaoFiscalNoEstID     *int   `json:"operacao_fiscal_no_est_id,omitempty"`
	OperacaoFiscalForaEstID   *int   `json:"operacao_fiscal_fora_est_id,omitempty"`
	OperacaoFiscalNoEstSTID   *int   `json:"operacao_fiscal_no_est_st_id,omitempty"`
	OperacaoFiscalForaEstSTID *int   `json:"operacao_fiscal_fora_est_st_id,omitempty"`
	PlanoContasFinanceiroID   *int   `json:"plano_contas_financeiro_id,omitempty"`
	NaturezaOperacaoID        *int   `json:"natureza_operacao_id,omitempty"`
	ReceitaID                 *int   `json:"receita_id,omitempty"`
	DespesaID                 *int   `json:"despesa_id,omitempty"`
	ProcessoNFID              *int   `json:"processo_nf_id,omitempty"`
	CreatedAt                 string `json:"created_at"`
	UpdatedAt                 string `json:"updated_at"`
	CreatedBy                 *int   `json:"created_by,omitempty"`
	UpdatedBy                 *int   `json:"updated_by,omitempty"`
}

// ProcessoListResponse representa a resposta de listagem de processos.
type ProcessoListResponse struct {
	Items      []ProcessoResponse `json:"items"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE CONVERSÃO
// ============================================================

// ToModel converte ProcessoRequest para models.Processo.
func (r *ProcessoRequest) ToModel() (*models.Processo, error) {
	if r == nil {
		return nil, nil
	}
	processo := &models.Processo{}
	if err := utils.MapToModel(r, processo); err != nil {
		return nil, fmt.Errorf("erro ao mapear DTO para modelo de processo: %w", err)
	}
	// Definir situação padrão se não informada
	if processo.Situacao == 0 {
		processo.Situacao = constants.StatusAtivo
	}
	return processo, nil
}

// FromModel converte models.Processo para ProcessoResponse.
func (r *ProcessoResponse) FromModel(processo *models.Processo) {
	if processo == nil {
		return
	}
	_ = utils.MapToDTO(processo, r) // Ignora erro por enquanto, assume que o mapeamento direto é suficiente

	r.TipoOperacaoLabel = constants.TipoOperacao(*processo.TipoOperacao).String()
	r.SituacaoLabel = constants.Status(processo.Situacao).String()

	if processo.EmpresaFilial != nil {
		r.EmpresaFilialNome = processo.EmpresaFilial.Nome
	}
	if processo.RotinaContabil != nil {
		r.RotinaContabilDescricao = processo.RotinaContabil.Descricao
	}
	// Outros relacionamentos podem ser populados aqui se necessário

	r.CreatedAt = utils.FormatDateTime(processo.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(processo.UpdatedAt)
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o ProcessoRequest.
func (r *ProcessoRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}
	// Validar TipoOperacao usando constantes
	if constants.TipoOperacao(r.TipoOperacao).IsValid() {
		return apperrors.NewBadRequestError("TipoOperação deve ter valor (0 - Entrada ou 1- Saida)")
	}
	// Validar Situacao (se informada)
	if r.Situacao > 0 {
		if err := constants.Status(r.Situacao).IsValid(); err != nil {
			return err
		}
	}
	return nil
}
