// ============================================================
// FILE: condicao_pagamento_dto.go
// PACKAGE: dto
// ============================================================

package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/openerp/backend/internal/constants"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// REQUESTS
// ============================================================

// CondicaoPagamentoRequest representa a requisição para criar/atualizar uma condição de pagamento
type CondicaoPagamentoRequest struct {
	// ============================================================
	// DADOS PRINCIPAIS
	// ============================================================
	RotinaContabilID *int     `json:"rotina_contabil_id,omitempty"`
	PortadorID       int      `json:"portador_id" binding:"required"`
	TipoDocumentoID  int      `json:"tipo_documento_id" binding:"required"`
	FormaPagamentoID *int     `json:"forma_pagamento_id,omitempty"`
	Descricao        string   `json:"descricao" binding:"required"`
	TipoCondicao     int      `json:"tipo_condicao" binding:"required"`
	NumeroParcelas   int      `json:"numero_parcelas,omitempty"`
	DiasParcelas     *int     `json:"dias_parcelas,omitempty"`
	DiaPagamento     *int     `json:"dia_pagamento,omitempty"`
	Entrada          int      `json:"entrada,omitempty"`
	Juros            *float64 `json:"juros,omitempty"`
	Comissao         *float64 `json:"comissao,omitempty"`
	Desconto         *float64 `json:"desconto,omitempty"`
	Situacao         int      `json:"situacao,omitempty"`

	// ============================================================
	// USUÁRIO (para auditoria)
	// ============================================================
	CreatedBy *int `json:"created_by,omitempty"`
	UpdatedBy *int `json:"updated_by,omitempty"`
}

// ============================================================
// RESPONSES
// ============================================================

// CondicaoPagamentoResponse representa a resposta de uma condição de pagamento
type CondicaoPagamentoResponse struct {
	// ============================================================
	// DADOS PRINCIPAIS
	// ============================================================
	ID                 int      `json:"id"`
	RotinaContabilID   *int     `json:"rotina_contabil_id,omitempty"`
	PortadorID         int      `json:"portador_id"`
	PortadorNome       string   `json:"portador_nome,omitempty"`
	TipoDocumentoID    int      `json:"tipo_documento_id"`
	TipoDocumentoNome  string   `json:"tipo_documento_nome,omitempty"`
	FormaPagamentoID   *int     `json:"forma_pagamento_id,omitempty"`
	FormaPagamentoNome string   `json:"forma_pagamento_nome,omitempty"`
	Descricao          string   `json:"descricao"`
	TipoCondicao       int      `json:"tipo_condicao"`
	TipoCondicaoLabel  string   `json:"tipo_condicao_label"`
	NumeroParcelas     int      `json:"numero_parcelas"`
	DiasParcelas       *int     `json:"dias_parcelas,omitempty"`
	DiaPagamento       *int     `json:"dia_pagamento,omitempty"`
	Entrada            int      `json:"entrada"`
	EntradaLabel       string   `json:"entrada_label"`
	Juros              *float64 `json:"juros,omitempty"`
	Comissao           *float64 `json:"comissao,omitempty"`
	Desconto           *float64 `json:"desconto,omitempty"`
	Situacao           int      `json:"situacao"`
	SituacaoLabel      string   `json:"situacao_label"`

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

// CondicaoPagamentoListResponse representa a resposta de listagem de condições de pagamento
type CondicaoPagamentoListResponse struct {
	Items      []CondicaoPagamentoResponse `json:"items"`
	Total      int64                       `json:"total"`
	Page       int                         `json:"page"`
	Limit      int                         `json:"limit"`
	TotalPages int                         `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE CONVERSÃO (TO MODEL)
// ============================================================

// ToModel converte CondicaoPagamentoRequest para models.CondicaoPagamento
func (r *CondicaoPagamentoRequest) ToModel() (*models.CondicaoPagamento, error) {
	if r == nil {
		return nil, nil
	}

	condicao := &models.CondicaoPagamento{}

	// 1. Usar o mapper para copiar campos
	if err := utils.MapToModel(r, condicao); err != nil {
		return nil, err
	}

	// 2. Tratamentos especiais
	// Definir situação padrão (Ativo)
	if r.Situacao == 0 {
		condicao.Situacao = int(constants.SituacaoCondicaoAtivo)
	}

	// Definir entrada padrão (Não)
	if r.Entrada == 0 {
		condicao.Entrada = constants.CondicaoEntradaNao
	}

	// Se for à vista, número de parcelas = 1
	if constants.TipoCondicao(r.TipoCondicao).IsAVista() {
		condicao.NumeroParcelas = 1
	}

	return condicao, nil
}

// ============================================================
// MÉTODOS DE CONVERSÃO (FROM MODEL)
// ============================================================

// FromModel converte models.CondicaoPagamento para CondicaoPagamentoResponse
func (r *CondicaoPagamentoResponse) FromModel(condicao *models.CondicaoPagamento) (*CondicaoPagamentoResponse, error) {
	if condicao == nil {
		return nil, nil
	}

	// 1. Usar o mapper para copiar campos
	if err := utils.MapToDTO(condicao, r); err != nil {
		return nil, err
	}

	// 2. Preencher campos calculados (labels)
	r.TipoCondicaoLabel = constants.TipoCondicao(condicao.TipoCondicao).String()
	r.SituacaoLabel = constants.SituacaoCondicaoPagamento(condicao.Situacao).String()
	r.EntradaLabel = condicao.Entrada.String()

	// 3. Nomes dos relacionamentos (se carregados)
	if condicao.Portador != nil {
		r.PortadorNome = condicao.Portador.Nome
	}

	if condicao.TipoDocumento != nil {
		r.TipoDocumentoNome = condicao.TipoDocumento.Nome
	}

	if condicao.FormaPagamento != nil {
		r.FormaPagamentoNome = condicao.FormaPagamento.Descricao
	}

	// 4. Formatar datas
	r.CreatedAt = utils.FormatDateTime(condicao.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(condicao.UpdatedAt)

	return r, nil
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o CondicaoPagamentoRequest
func (r *CondicaoPagamentoRequest) Validate() error {
	// 1. Validação com go-playground/validator
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}

	// 2. Validar campos obrigatórios usando utils
	if err := utils.ValidateMandatoryFields(r); err != nil {
		return err
	}

	// 3. Validar TipoCondicao usando constantes
	if err := constants.TipoCondicao(r.TipoCondicao).IsValid(); err != nil {
		return err
	}

	// 4. Validar Situacao (se informada)
	if r.Situacao > 0 {
		if err := constants.SituacaoCondicaoPagamento(r.Situacao).IsValid(); err != nil {
			return err
		}
	}

	// 5. Validar Entrada (0 ou 1)
	if r.Entrada != 0 && r.Entrada != 1 {
		return fmt.Errorf("entrada deve ser 0 (Não) ou 1 (Sim), recebido: %d", r.Entrada)
	}

	// 6. Validar NumeroParcelas
	if constants.TipoCondicao(r.TipoCondicao).IsAPrazo() && r.NumeroParcelas < 2 {
		return fmt.Errorf("para condição a prazo, numero_parcelas deve ser maior que 1, recebido: %d", r.NumeroParcelas)
	}

	if constants.TipoCondicao(r.TipoCondicao).IsAVista() && r.NumeroParcelas > 1 {
		return fmt.Errorf("para condição à vista, numero_parcelas deve ser 1, recebido: %d", r.NumeroParcelas)
	}

	// 7. Validar DiasParcelas (se informado)
	if r.DiasParcelas != nil && *r.DiasParcelas < 0 {
		return fmt.Errorf("dias_parcelas não pode ser negativo, recebido: %d", *r.DiasParcelas)
	}

	// 8. Validar DiaPagamento (se informado)
	if r.DiaPagamento != nil && (*r.DiaPagamento < 1 || *r.DiaPagamento > 31) {
		return fmt.Errorf("dia_pagamento deve estar entre 1 e 31, recebido: %d", *r.DiaPagamento)
	}

	// 9. Validar percentuais
	if r.Juros != nil && (*r.Juros < 0 || *r.Juros > 100) {
		return fmt.Errorf("juros deve estar entre 0 e 100, recebido: %.2f", *r.Juros)
	}

	if r.Comissao != nil && (*r.Comissao < 0 || *r.Comissao > 100) {
		return fmt.Errorf("comissao deve estar entre 0 e 100, recebido: %.2f", *r.Comissao)
	}

	if r.Desconto != nil && (*r.Desconto < 0 || *r.Desconto > 100) {
		return fmt.Errorf("desconto deve estar entre 0 e 100, recebido: %.2f", *r.Desconto)
	}

	return nil
}

// ============================================================
// FUNÇÕES AUXILIARES
// ============================================================

// // getSimNaoLabel retorna o label para campos booleanos (0/1)
// func getSimNaoLabel(valor int) string {
// 	switch valor {
// 	case constants.CondicaoEntradaSim:
// 		return "Sim"
// 	case constants.CondicaoEntradaNao:
// 		return "Não"
// 	default:
// 		return "Não"
// 	}
// }
