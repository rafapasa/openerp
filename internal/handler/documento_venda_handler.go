package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

type DocumentoVendaHandler struct {
	service *service.DocumentoVendaService
}

func NewDocumentoVendaHandler(s *service.DocumentoVendaService) *DocumentoVendaHandler {
	return &DocumentoVendaHandler{
		service: s,
	}
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
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	doc, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	utils.RespondWithCreated(c, doc)
}

// @Summary      Busca um documento de venda por ID
// @Description  Retorna os detalhes de um documento de venda específico.
// @Tags         Documentos de Venda
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Documento de Venda"
// @Success      200  {object}  dto.DocumentoVendaResponse
// @Failure      404  {object}  utils.ErrorResponse "Documento de venda não encontrado"
// @Router       /documentos-venda/{id} [get]
func (h *DocumentoVendaHandler) GetByID(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	doc, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var response dto.DocumentoVendaResponse
	if err := utils.MapToModel(doc, &response); err != nil {
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
// @Param        id         path      int                       true  "ID do Documento de Venda"
// @Param        documento  body      dto.DocumentoVendaRequest  true  "Dados para atualizar o documento de venda"
// @Success      200        {object}  dto.DocumentoVendaResponse
// @Failure      400        {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      404        {object}  utils.ErrorResponse "Documento de venda não encontrado"
// @Router       /documentos-venda/{id} [put]
func (h *DocumentoVendaHandler) Update(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.DocumentoVendaRequest
	if utils.BindAndValidateOrRespond(c, &req) {
		return
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
		Documentos: respDocs,
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
// @Param        id   path      int  true  "ID do Documento de Venda"
// @Success      204  "Nenhum conteúdo"
// @Failure      400  {object}  utils.ErrorResponse "Erro ao excluir"
// @Failure      404  {object}  utils.ErrorResponse "Documento de venda não encontrado"
// @Router       /documentos-venda/{id} [delete]
func (h *DocumentoVendaHandler) Delete(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	utils.RespondWithNoContent(c)
}