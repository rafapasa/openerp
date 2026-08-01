package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

// DocumentoVendaPagamentoHandler gerencia as requisições HTTP para pagamentos de documento de venda.
type DocumentoVendaPagamentoHandler struct {
	documentoVendaService         service.DocumentoVendaService
	documentoVendaPagamentoService service.DocumentoVendaPagamentoService
}

// NewDocumentoVendaPagamentoHandler cria uma nova instância de DocumentoVendaPagamentoHandler.
func NewDocumentoVendaPagamentoHandler(
	ddvService service.DocumentoVendaService,
	dvpService service.DocumentoVendaPagamentoService,
) *DocumentoVendaPagamentoHandler {
	return &DocumentoVendaPagamentoHandler{
		documentoVendaService:         ddvService,
		documentoVendaPagamentoService: dvpService,
	}
}

// @Summary      Adiciona um novo pagamento a um documento de venda
// @Description  Adiciona um novo pagamento a um documento de venda específico.
// @Tags         Documentos de Venda - Pagamentos
// @Accept       json
// @Produce      json
// @Param        documento_venda_id  path      int                               true  "ID do Documento de Venda"
// @Param        pagamento           body      dto.DocumentoVendaPagamentoRequest  true  "Dados para adicionar o pagamento"
// @Success      201                 {object}  dto.DocumentoVendaResponse
// @Failure      400                 {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      500                 {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /documentos/venda/{documento_venda_id}/pagamentos [post]
func (h *DocumentoVendaPagamentoHandler) Create(c *gin.Context) {
	documentoVendaID, ok := utils.ParseIDParam(c, "documento_venda_id")
	if !ok {
		return
	}

	var req dto.DocumentoVendaPagamentoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	req.DocumentoVendaID = documentoVendaID
	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	// Delega a criação do pagamento ao DocumentoVendaService para garantir o recálculo dos totais
	updatedDoc, err := h.documentoVendaService.AddPagamento(&req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.DocumentoVendaResponse
	resp.FromModel(updatedDoc)
	utils.RespondWithCreated(c, resp)
}

// @Summary      Busca um pagamento de documento de venda por ID
// @Description  Retorna os detalhes de um pagamento específico de um documento de venda.
// @Tags         Documentos de Venda - Pagamentos
// @Accept       json
// @Produce      json
// @Param        documento_venda_id  path      int  true  "ID do Documento de Venda"
// @Param        item                path      int  true  "Número do Item do Pagamento"
// @Success      200                 {object}  dto.DocumentoVendaPagamentoResponse
// @Failure      404                 {object}  utils.ErrorResponse "Pagamento não encontrado"
// @Router       /documentos/venda/{documento_venda_id}/pagamentos/{item} [get]
func (h *DocumentoVendaPagamentoHandler) GetByID(c *gin.Context) {
	documentoVendaID, ok := utils.ParseIDParam(c, "documento_venda_id")
	if !ok {
		return
	}
	itemNum, ok := utils.ParseIDParam(c, "item")
	if !ok {
		return
	}

	pagamento, err := h.documentoVendaPagamentoService.FindByID(documentoVendaID, itemNum)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.DocumentoVendaPagamentoResponse
	resp.FromModel(pagamento)
	utils.RespondWithOK(c, resp)
}

// @Summary      Atualiza um pagamento de documento de venda
// @Description  Atualiza os dados de um pagamento existente de um documento de venda.
// @Tags         Documentos de Venda - Pagamentos
// @Accept       json
// @Produce      json
// @Param        documento_venda_id  path      int                               true  "ID do Documento de Venda"
// @Param        item                path      int                               true  "Número do Item do Pagamento"
// @Param        pagamento           body      dto.DocumentoVendaPagamentoRequest  true  "Dados para atualizar o pagamento"
// @Success      200                 {object}  dto.DocumentoVendaResponse
// @Failure      400                 {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      404                 {object}  utils.ErrorResponse "Pagamento não encontrado"
// @Router       /documentos/venda/{documento_venda_id}/pagamentos/{item} [put]
func (h *DocumentoVendaPagamentoHandler) Update(c *gin.Context) {
	documentoVendaID, ok := utils.ParseIDParam(c, "documento_venda_id")
	if !ok {
		return
	}
	itemNum, ok := utils.ParseIDParam(c, "item")
	if !ok {
		return
	}

	var req dto.DocumentoVendaPagamentoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	req.DocumentoVendaID = documentoVendaID
	userID := middleware.GetUserID(c)
	req.UpdatedBy = &userID

	updatedDoc, err := h.documentoVendaService.EditPagamento(documentoVendaID, itemNum, &req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.DocumentoVendaResponse
	resp.FromModel(updatedDoc)
	utils.RespondWithOK(c, resp)
}

// @Summary      Exclui um pagamento de documento de venda
// @Description  Exclui logicamente um pagamento de um documento de venda.
// @Tags         Documentos de Venda - Pagamentos
// @Accept       json
// @Produce      json
// @Param        documento_venda_id  path      int  true  "ID do Documento de Venda"
// @Param        item                path      int  true  "Número do Item do Pagamento"
// @Success      204                 "Nenhum conteúdo"
// @Failure      400                 {object}  utils.ErrorResponse "Erro ao excluir"
// @Router       /documentos/venda/{documento_venda_id}/pagamentos/{item} [delete]
func (h *DocumentoVendaPagamentoHandler) Delete(c *gin.Context) {
	documentoVendaID, ok := utils.ParseIDParam(c, "documento_venda_id")
	if !ok {
		return
	}
	itemNum, ok := utils.ParseIDParam(c, "item")
	if !ok {
		return
	}

	_, err := h.documentoVendaService.DeletePagamento(documentoVendaID, itemNum)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary      Lista os pagamentos de um documento de venda
// @Description  Retorna uma lista paginada de pagamentos para um documento de venda específico.
// @Tags         Documentos de Venda - Pagamentos
// @Accept       json
// @Produce      json
// @Param        documento_venda_id  path      int  true  "ID do Documento de Venda"
// @Param        limit               query     int  false  "Número de registros por página"
// @Param        offset              query     int  false  "Offset para a paginação"
// @Success      200                 {object}  dto.DocumentoVendaPagamentoListResponse
// @Failure      500                 {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /documentos/venda/{documento_venda_id}/pagamentos [get]
func (h *DocumentoVendaPagamentoHandler) List(c *gin.Context) {
	documentoVendaID, ok := utils.ParseIDParam(c, "documento_venda_id")
	if !ok {
		return
	}

	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)

	payments, total, err := h.documentoVendaPagamentoService.ListByDocumentoVendaID(limit, offset, documentoVendaID)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	respPayments := make([]dto.DocumentoVendaPagamentoResponse, len(payments))
	for i, payment := range payments {
		var respPayment dto.DocumentoVendaPagamentoResponse
		respPayment.FromModel(&payment)
		respPayments[i] = respPayment //
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	utils.RespondWithOK(c, dto.DocumentoVendaPagamentoListResponse{
		Items:      respPayments,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}