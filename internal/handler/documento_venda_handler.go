package handler

// [RASCUNHO] → [EMITIDO] → [FATURADO] → [ENTREGUE] → [CONCLUÍDO]
//      ↓            ↓            ↓            ↓            ↓
//   Criação    Efetivação   Faturamento   Entrega     Finalização
//   (proposta)   (venda)     (financeiro)  (logística)

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

type DocumentoVendaHandler struct {
	service service.DocumentoVendaService
}

func NewDocumentoVendaHandler(s service.DocumentoVendaService) *DocumentoVendaHandler {
	return &DocumentoVendaHandler{
		service: s,
	}
}

// @Summary      Emite um documento de venda
// @Description  Emite um documento de venda, alterando seu status para emitido.
// @Tags         Documentos de Venda
// @Accept       json
// @Produce      json
// @Param        documento_id   path      int  true  "ID do Documento de Venda"
// @Success      200  {object}  dto.DocumentoVendaResponse
// @Failure      400  {object}  utils.ErrorResponse "Erro ao emitir documento"
// @Failure      404  {object}  utils.ErrorResponse "Documento de venda não encontrado"
// @Failure      422  {object}  utils.ErrorResponse "Documento não pode ser emitido (validação de negócio)"
// @Router       /documentos-venda/{documento_id}/emitir [post]
func (h *DocumentoVendaHandler) Emitir(c *gin.Context) {
	panic("Implementar Emitir") // TODO: Implementar
	// FIXME: Descrição do bug a ser corrigido
	// OPTIMIZE: Descrição da otimização
	// SECURITY: Descrição da melhoria de segurança
}

// @Summary      Cancela um documento de venda
// @Description  Cancela um documento de venda, alterando seu status para cancelado.
// @Tags         Documentos de Venda
// @Accept       json
// @Produce      json
// @Param        documento_id   path      int  true  "ID do Documento de Venda"
// @Success      200  {object}  dto.DocumentoVendaResponse
// @Failure      400  {object}  utils.ErrorResponse "Erro ao cancelar documento"
// @Failure      404  {object}  utils.ErrorResponse "Documento de venda não encontrado"
// @Failure      422  {object}  utils.ErrorResponse "Documento não pode ser cancelado (validação de negócio)"
// @Router       /documentos-venda/{documento_id}/cancelar [post]
func (h *DocumentoVendaHandler) Cancelar(c *gin.Context) {
	panic("Implementar Cancelar") // TODO: Implementar
	// FIXME: Descrição do bug a ser corrigido
	// OPTIMIZE: Descrição da otimização
	// SECURITY: Descrição da melhoria de segurança
}

// @Summary      Cria um novo documento de venda
// @Description  Cria um novo documento de venda com base nos dados fornecidos.
// @Tags         Documentos de Venda
// @Accept       json
// @Produce      json
// @Param        documento  body      dto.DocumentoVendaRequest  true  "Dados para criar o documento de venda"
// @Success      201        {object}  dto.DocumentoVendaResponse
// @Failure      400        {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      500        {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /documentos-venda [post]
func (h *DocumentoVendaHandler) Create(c *gin.Context) {
	var req dto.DocumentoVendaRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	doc, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.DocumentoVendaResponse
	if _, err := resp.FromModel(doc); err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}
	utils.RespondWithCreated(c, resp)
}

// @Summary      Busca um documento de venda por ID
// @Description  Retorna os detalhes de um documento de venda específico.
// @Tags         Documentos de Venda
// @Accept       json
// @Produce      json
// @Param        documento_id   path      int  true  "ID do Documento de Venda"
// @Success      200  {object}  dto.DocumentoVendaResponse
// @Failure      404  {object}  utils.ErrorResponse "Documento de venda não encontrado"
// @Router       /documentos-venda/{documento_id} [get]
func (h *DocumentoVendaHandler) GetByID(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "documento_id")
	if !ok {
		return
	}

	doc, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var response dto.DocumentoVendaResponse
	if _, err := response.FromModel(doc); err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	utils.RespondWithOK(c, response)
}

// @Summary      Atualiza um documento de venda
// @Description  Atualiza os dados de um documento de venda existente.
// @Tags         Documentos de Venda
// @Accept       json
// @Produce      json
// @Param        documento_id   path      int                       true  "ID do Documento de Venda"
// @Param        documento  body      dto.DocumentoVendaRequest  true  "Dados para atualizar o documento de venda"
// @Success      200        {object}  dto.DocumentoVendaResponse
// @Failure      400        {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      404        {object}  utils.ErrorResponse "Documento de venda não encontrado"
// @Router       /documentos-venda/{documento_id} [put]
func (h *DocumentoVendaHandler) Update(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "documento_id")
	if !ok {
		return
	}

	var req dto.DocumentoVendaRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
	}

	doc, err := h.service.Update(id, &req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	utils.RespondWithOK(c, doc)
}

// @Summary      Lista os documentos de venda
// @Description  Retorna uma lista paginada de documentos de venda, com suporte a filtros.
// @Tags         Documentos de Venda
// @Accept       json
// @Produce      json
// @Param        limit   query     int  false  "Número de registros por página"
// @Param        offset  query     int  false  "Offset para a paginação"
// @Param        tipo    query     string  false  "Filtrar por tipo"
// @Param        situacao query     int  false  "Filtrar por situação"
// @Success      200     {object}  dto.DocumentoVendaListResponse
// @Failure      500     {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /documentos-venda [get]
func (h *DocumentoVendaHandler) List(c *gin.Context) {
	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)

	docs, total, err := h.service.List(limit, offset, utils.QueryParamsToFilters(c))
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	respDocs := make([]dto.DocumentoVendaResponse, len(docs))
	for i, doc := range docs {
		var respDoc dto.DocumentoVendaResponse
		respDoc.FromModel(&doc)
		respDocs[i] = respDoc
	}

	response := dto.DocumentoVendaListResponse{
		Items:      respDocs,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: utils.CalculateTotalPages(int(total), limit),
	}

	utils.RespondWithOK(c, response)
}

// @Summary      Exclui um documento de venda
// @Description  Realiza a exclusão lógica de um documento de venda.
// @Tags         Documentos de Venda
// @Accept       json
// @Produce      json
// @Param        documento_id   path      int  true  "ID do Documento de Venda"
// @Success      204  "Nenhum conteúdo"
// @Failure      400  {object}  utils.ErrorResponse "Erro ao excluir"
// @Failure      404  {object}  utils.ErrorResponse "Documento de venda não encontrado"
// @Router       /documentos-venda/{documento_id} [delete]
func (h *DocumentoVendaHandler) Delete(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "documento_id")
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	utils.RespondWithNoContent(c)
}

// @Summary      Lista os itens de um documento de venda
// @Description  Retorna todos os itens de um documento de venda específico.
// @Tags         Documentos de Venda
// @Accept       json
// @Produce      json
// @Param        documento_id   path      int  true  "ID do Documento de Venda"
// @Success      200  {array}   dto.DocumentoVendaItemResponse
// @Failure      404  {object}  utils.ErrorResponse "Documento de venda não encontrado"
// @Router       /documentos-venda/{documento_id}/itens [get]
func (h *DocumentoVendaHandler) ListItem(c *gin.Context) {
	documentoVendaID, ok := utils.ParseIDParam(c, "documento_id")
	if !ok {
		return
	}

	items, err := h.service.ListItem(documentoVendaID)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var respItems []dto.DocumentoVendaItemResponse
	for _, item := range items {
		var resp dto.DocumentoVendaItemResponse
		resp.FromModel(&item)
		respItems = append(respItems, resp)
	}

	utils.RespondWithOK(c, respItems)
}

// @Summary      Adiciona um item ao documento de venda
// @Description  Adiciona um novo item a um documento de venda existente.
// @Tags         Documentos de Venda
// @Accept       json
// @Produce      json
// @Param        documento_id   path      int                           true  "ID do Documento de Venda"
// @Param        item  body      dto.DocumentoVendaItemRequest  true  "Dados do item"
// @Success      201   {object}  dto.DocumentoVendaResponse
// @Failure      400   {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      404   {object}  utils.ErrorResponse "Documento de venda não encontrado"
// @Failure      422   {object}  utils.ErrorResponse "Erro de validação de negócio (ex: estoque insuficiente)"
// @Router       /documentos-venda/{documento_id}/itens [post]
func (h *DocumentoVendaHandler) AddItem(c *gin.Context) {
	documentoVendaID, ok := utils.ParseIDParam(c, "documento_id")
	if !ok {
		return
	}

	var req dto.DocumentoVendaItemRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	req.DocumentoVendaID = documentoVendaID

	ddv, err := h.service.AddItem(&req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.DocumentoVendaResponse
	resp.FromModel(ddv)

	utils.RespondWithCreated(c, resp)
}

// @Summary      Edita um item do documento de venda
// @Description  Atualiza os dados de um item específico de um documento de venda.
// @Tags         Documentos de Venda
// @Accept       json
// @Produce      json
// @Param        documento_id   path      int                           true  "ID do Documento de Venda"
// @Param        item  path      int                           true  "ID do Item"
// @Param        item  body      dto.DocumentoVendaItemRequest  true  "Dados atualizados do item"
// @Success      200   {object}  dto.DocumentoVendaResponse
// @Failure      400   {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      404   {object}  utils.ErrorResponse "Documento de venda ou item não encontrado"
// @Failure      422   {object}  utils.ErrorResponse "Erro de validação de negócio"
// @Router       /documentos-venda/{documento_id}/itens/{item} [put]
func (h *DocumentoVendaHandler) EditItem(c *gin.Context) {
	documentoVendaID, ok := utils.ParseIDParam(c, "documento_id")
	if !ok {
		return
	}
	item, ok := utils.ParseIDParam(c, "item")
	if !ok {
		return
	}

	var req dto.DocumentoVendaItemRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	req.DocumentoVendaID = documentoVendaID

	ddv, err := h.service.EditItem(documentoVendaID, item, &req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.DocumentoVendaResponse
	resp.FromModel(ddv)

	utils.RespondWithOK(c, resp)
}

// @Summary      Remove um item do documento de venda
// @Description  Remove um item específico de um documento de venda.
// @Tags         Documentos de Venda
// @Accept       json
// @Produce      json
// @Param        documento_id   path      int  true  "ID do Documento de Venda"
// @Param        item  path      int  true  "ID do Item"
// @Success      200   {object}  dto.DocumentoVendaResponse
// @Failure      400   {object}  utils.ErrorResponse "Erro ao remover item"
// @Failure      404   {object}  utils.ErrorResponse "Documento de venda ou item não encontrado"
// @Failure      422   {object}  utils.ErrorResponse "Erro de validação de negócio"
// @Router       /documentos-venda/{documento_id}/itens/{item} [delete]
func (h *DocumentoVendaHandler) DeleteItem(c *gin.Context) {
	documentoVendaID, ok := utils.ParseIDParam(c, "documento_id")
	if !ok {
		return
	}
	item, ok := utils.ParseIDParam(c, "item")
	if !ok {
		return
	}
	ddv, err := h.service.DeleteItem(documentoVendaID, item)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	response := dto.DocumentoVendaResponse{}
	response.FromModel(ddv)
	utils.RespondWithOK(c, response)
}

// @Summary      Lista os pagamentos de um documento de venda
// @Description  Retorna todos os pagamentos de um documento de venda específico.
// @Tags         Documentos de Venda
// @Accept       json
// @Produce      json
// @Param        documento_id   path      int  true  "ID do Documento de Venda"
// @Success      200  {array}   dto.DocumentoVendaPagamentoResponse
// @Failure      404  {object}  utils.ErrorResponse "Documento de venda não encontrado"
// @Router       /documentos-venda/{documento_id}/pagamentos [get]
func (h *DocumentoVendaHandler) ListPagamento(c *gin.Context) {
	documentoVendaID, ok := utils.ParseIDParam(c, "documento_id")
	if !ok {
		return
	}

	pagamentos, err := h.service.ListPagamento(documentoVendaID)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var respPagamentos []dto.DocumentoVendaPagamentoResponse
	for _, pagamento := range pagamentos {
		var resp dto.DocumentoVendaPagamentoResponse
		resp.FromModel(&pagamento)
		respPagamentos = append(respPagamentos, resp)
	}

	utils.RespondWithOK(c, respPagamentos)
}

// @Summary      Adiciona um pagamento ao documento de venda
// @Description  Adiciona um novo pagamento a um documento de venda existente.
// @Tags         Documentos de Venda
// @Accept       json
// @Produce      json
// @Param        documento_id   path      int                                 true  "ID do Documento de Venda"
// @Param        pagamento  body      dto.DocumentoVendaPagamentoRequest  true  "Dados do pagamento"
// @Success      201        {object}  dto.DocumentoVendaResponse
// @Failure      400        {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      404        {object}  utils.ErrorResponse "Documento de venda não encontrado"
// @Failure      422        {object}  utils.ErrorResponse "Erro de validação de negócio (ex: valor excede o total)"
// @Router       /documentos-venda/{documento_id}/pagamentos [post]
func (h *DocumentoVendaHandler) AddPagamento(c *gin.Context) {
	documentoVendaID, ok := utils.ParseIDParam(c, "documento_id")
	if !ok {
		return
	}

	var req dto.DocumentoVendaPagamentoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	req.DocumentoVendaID = documentoVendaID

	ddv, err := h.service.AddPagamento(&req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.DocumentoVendaResponse
	resp.FromModel(ddv)

	utils.RespondWithCreated(c, resp)
}

// @Summary      Edita um pagamento do documento de venda
// @Description  Atualiza os dados de um pagamento específico de um documento de venda.
// @Tags         Documentos de Venda
// @Accept       json
// @Produce      json
// @Param        documento_id   path      int                                 true  "ID do Documento de Venda"
// @Param        pagamento  path      int                                 true  "ID do Pagamento"
// @Param        pagamento  body      dto.DocumentoVendaPagamentoRequest  true  "Dados atualizados do pagamento"
// @Success      200        {object}  dto.DocumentoVendaResponse
// @Failure      400        {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      404        {object}  utils.ErrorResponse "Documento de venda ou pagamento não encontrado"
// @Failure      422        {object}  utils.ErrorResponse "Erro de validação de negócio"
// @Router       /documentos-venda/{documento_id}/pagamentos/{pagamento} [put]
func (h *DocumentoVendaHandler) UpdatePagamento(c *gin.Context) {
	documentoVendaID, ok := utils.ParseIDParam(c, "documento_id")
	if !ok {
		return
	}
	pagamento, ok := utils.ParseIDParam(c, "pagamento")
	if !ok {
		return
	}
	var req dto.DocumentoVendaPagamentoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}
	req.DocumentoVendaID = documentoVendaID

	ddv, err := h.service.EditPagamento(documentoVendaID, pagamento, &req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}
	resp := dto.DocumentoVendaResponse{}
	resp.FromModel(ddv)
	utils.RespondWithOK(c, resp)
}
