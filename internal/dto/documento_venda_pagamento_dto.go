// ============================================================
// FILE: documento_venda_pagamento_dto.go
// PACKAGE: dto
// ============================================================

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

// DocumentoVendaPagamentoRequest representa um pagamento no request
type DocumentoVendaPagamentoRequest struct {
	DocumentoVendaID  int     `json:"documento_venda_id,omitempty"` // Changed from "id"
	PortadorID        int     `json:"portador_id" binding:"required"`
	TipoDocumentoID   int     `json:"tipo_documento_id" binding:"required"`
	FormaPagamentoID  *int    `json:"forma_pagamento_id,omitempty"` // Changed to *int
	RotinaContabilID  *int    `json:"rotina_contabil_id,omitempty"`
	DataVencimento    *string `json:"data_vencimento" binding:"required"` // Changed to *string for parsing, and required
	Valor             float64 `json:"valor" binding:"required,gt=0"`
	Documento         *string `json:"documento,omitempty"`
	CodigoAutorizacao *string `json:"codigo_autorizacao,omitempty"`
	COO               *int    `json:"coo,omitempty"`
	BandeiraCartao    *string `json:"bandeira_cartao,omitempty"`
	NSU               *string `json:"nsu,omitempty"`
	Autorizacao       *string `json:"autorizacao,omitempty"`
	CreatedBy         *int    `json:"created_by,omitempty"`
	UpdatedBy         *int    `json:"updated_by,omitempty"`
}

// ============================================================
// RESPONSES
// ============================================================

// DocumentoVendaPagamentoResponse representa a resposta de um pagamento
type DocumentoVendaPagamentoResponse struct {
	DocumentoVendaID   int     `json:"documento_venda_id"`
	Item               int     `json:"item"` // Sequencial por documento
	PortadorID         int     `json:"portador_id"`
	PortadorNome       string  `json:"portador_nome,omitempty"`
	TipoDocumentoID    int     `json:"tipo_documento_id"`
	TipoDocumentoNome  string  `json:"tipo_documento_nome,omitempty"`
	FormaPagamentoID   *int    `json:"forma_pagamento_id,omitempty"`
	FormaPagamentoNome string  `json:"forma_pagamento_nome,omitempty"`
	RotinaContabilID   *int    `json:"rotina_contabil_id,omitempty"`
	DataVencimento     string  `json:"data_vencimento"`
	Valor              float64 `json:"valor"`
	Status             int     `json:"status"`                 // Added from model
	StatusLabel        string  `json:"status_label,omitempty"` // Added for display
	Documento          *string `json:"documento,omitempty"`
	CodigoAutorizacao  *string `json:"codigo_autorizacao,omitempty"`
	COO                *int    `json:"coo,omitempty"`
	BandeiraCartao     *string `json:"bandeira_cartao,omitempty"`
	NSU                *string `json:"nsu,omitempty"`
	Autorizacao        *string `json:"autorizacao,omitempty"`
	CreatedAt          string  `json:"created_at"` // Formatted string
	UpdatedAt          string  `json:"updated_at"` // Formatted string
	CreatedBy          *int    `json:"created_by,omitempty"`
	UpdatedBy          *int    `json:"updated_by,omitempty"`
}

type DocumentoVendaPagamentoListResponse struct {
	Items      []DocumentoVendaPagamentoResponse `json:"items"`
	Total      int64                             `json:"total"`
	Page       int                               `json:"page"`
	Limit      int                               `json:"limit"`
	TotalPages int                               `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE CONVERSÃO
// ============================================================

// ToModel converte DocumentoVendaPagamentoRequest para models.DocumentoVendaPagamento
func (r *DocumentoVendaPagamentoRequest) ToModel() (*models.DocumentoVendaPagamento, error) {
	if r == nil {
		return nil, nil
	}
	pagamento := &models.DocumentoVendaPagamento{}
	if err := utils.MapToModel(r, pagamento); err != nil {
		return nil, err
	}

	// Handle DataVencimento parsing from *string to time.Time
	if r.DataVencimento != nil && *r.DataVencimento != "" {
		parsedTime, err := utils.ParseDate(*r.DataVencimento)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter data_vencimento: %w", err)
		}
		pagamento.DataVencimento = parsedTime
	}
	// Set status from DTO if provided, otherwise default to open

	return pagamento, nil
}

// FromModel converte models.DocumentoVendaPagamento para DocumentoVendaPagamentoResponse
func (r *DocumentoVendaPagamentoResponse) FromModel(pagamento *models.DocumentoVendaPagamento) {
	if pagamento == nil {
		return
	}
	_ = utils.MapToDTO(pagamento, r) // Ignore error for now, assume direct mapping is sufficient

	// Populate relationship names
	if pagamento.Portador != nil {
		r.PortadorNome = pagamento.Portador.Nome
	}
	if pagamento.TipoDocumento != nil {
		r.TipoDocumentoNome = pagamento.TipoDocumento.Nome
	}
	if pagamento.FormaPagamento != nil {
		r.FormaPagamentoNome = pagamento.FormaPagamento.Descricao
	}

	// Format dates
	r.DataVencimento = utils.FormatDate(pagamento.DataVencimento)
	r.CreatedAt = utils.FormatDateTime(pagamento.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(pagamento.UpdatedAt)
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o DocumentoVendaPagamentoRequest
func (r *DocumentoVendaPagamentoRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}

	return nil
}
