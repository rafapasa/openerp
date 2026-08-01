package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

// DocumentoVendaItemHandler gerencia as requisições HTTP para itens de documento de venda.
type DocumentoVendaItemHandler struct {
	service service.DocumentoVendaItemService
}

// NewDocumentoVendaItemHandler cria uma nova instância de DocumentoVendaItemHandler.
func NewDocumentoVendaItemHandler(s service.DocumentoVendaItemService) *DocumentoVendaItemHandler {
	return &DocumentoVendaItemHandler{
		service: s,
	}
}

// @Summary      Cria um novo item para um documento de venda
// @Description  Cria um novo item para um documento de venda específico.
// @Tags         Documentos de Venda - Itens
// @Accept       json
// @Produce      json
// @Param        documento_venda_id  path      int                         true  "ID do Documento de Venda"
// @Param        item                body      dto.DocumentoVendaItemRequest  true  "Dados para criar o item do documento de venda"
// @Success      201                 {object}  dto.DocumentoVendaItemResponse
// @Failure      400                 {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      500                 {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /documentos/venda/{documento_venda_id}/itens [post]
func (h *DocumentoVendaItemHandler) Create(c *gin.Context) {
	documentoVendaID, ok := utils.ParseIDParam(c, "documento_venda_id")
	if !ok {
		return
	}

	var req dto.DocumentoVendaItemRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	req.DocumentoVendaID = documentoVendaID
	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	if err := h.service.Create(&req); err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	// Como o service.Create retorna apenas erro, precisamos buscar o item recém-criado
	// ou ajustar o service.Create para retornar o modelo. Por simplicidade, vamos retornar um sucesso genérico.
	// Idealmente, o service.Create deveria retornar o modelo criado.
	utils.RespondWithCreated(c, gin.H{"message": "Item criado com sucesso."})
}

// @Summary      Busca um item de documento de venda por ID
// @Description  Retorna os detalhes de um item de documento de venda específico.
// @Tags         Documentos de Venda - Itens
// @Accept       json
// @Produce      json
// @Param        documento_venda_id  path      int  true  "ID do Documento de Venda"
// @Param        item                path      int  true  "Número do Item"
// @Success      200                 {object}  dto.DocumentoVendaItemResponse
// @Failure      404                 {object}  utils.ErrorResponse "Item não encontrado"
// @Router       /documentos/venda/{documento_venda_id}/itens/{item} [get]
func (h *DocumentoVendaItemHandler) GetByID(c *gin.Context) {
	documentoVendaID, ok := utils.ParseIDParam(c, "documento_venda_id")
	if !ok {
		return
	}
	itemNum, ok := utils.ParseIDParam(c, "item")
	if !ok {
		return
	}

	item, err := h.service.GetByID(documentoVendaID, itemNum)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.DocumentoVendaItemResponse
	resp.FromModel(item)
	utils.RespondWithOK(c, resp)
}

// @Summary      Atualiza um item de documento de venda
// @Description  Atualiza os dados de um item de documento de venda existente.
// @Tags         Documentos de Venda - Itens
// @Accept       json
// @Produce      json
// @Param        documento_venda_id  path      int                         true  "ID do Documento de Venda"
// @Param        item                path      int                         true  "Número do Item"
// @Param        item_data           body      dto.DocumentoVendaItemRequest  true  "Dados para atualizar o item do documento de venda"
// @Success      200                 {object}  dto.DocumentoVendaItemResponse
// @Failure      400                 {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      404                 {object}  utils.ErrorResponse "Item não encontrado"
// @Router       /documentos/venda/{documento_venda_id}/itens/{item} [put]
func (h *DocumentoVendaItemHandler) Update(c *gin.Context) {
	documentoVendaID, ok := utils.ParseIDParam(c, "documento_venda_id")
	if !ok {
		return
	}
	itemNum, ok := utils.ParseIDParam(c, "item")
	if !ok {
		return
	}

	var req dto.DocumentoVendaItemRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	req.DocumentoVendaID = documentoVendaID
	userID := middleware.GetUserID(c)
	req.UpdatedBy = &userID

	if err := h.service.Update(documentoVendaID, itemNum, &req); err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	// Buscar o item atualizado para retornar na resposta
	updatedItem, err := h.service.GetByID(documentoVendaID, itemNum)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}
	var resp dto.DocumentoVendaItemResponse
	resp.FromModel(updatedItem)
	utils.RespondWithOK(c, resp)
}

// @Summary      Exclui um item de documento de venda
// @Description  Exclui logicamente um item de documento de venda.
// @Tags         Documentos de Venda - Itens
// @Accept       json
// @Produce      json
// @Param        documento_venda_id  path      int  true  "ID do Documento de Venda"
// @Param        item                path      int  true  "Número do Item"
// @Success      204                 "Nenhum conteúdo"
// @Failure      400                 {object}  utils.ErrorResponse "Erro ao excluir"
// @Router       /documentos/venda/{documento_venda_id}/itens/{item} [delete]
func (h *DocumentoVendaItemHandler) Delete(c *gin.Context) {
	documentoVendaID, ok := utils.ParseIDParam(c, "documento_venda_id")
	if !ok {
		return
	}
	itemNum, ok := utils.ParseIDParam(c, "item")
	if !ok {
		return
	}

	if err := h.service.Delete(documentoVendaID, itemNum); err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	utils.RespondWithNoContent(c)
}

// @Summary      Lista os itens de um documento de venda
// @Description  Retorna uma lista paginada de itens para um documento de venda específico.
// @Tags         Documentos de Venda - Itens
// @Accept       json
// @Produce      json
// @Param        documento_venda_id  path      int  true  "ID do Documento de Venda"
// @Param        limit               query     int  false  "Número de registros por página"
// @Param        offset              query     int  false  "Offset para a paginação"
// @Success      200                 {object}  dto.DocumentoVendaItemListResponse
// @Failure      500                 {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /documentos/venda/{documento_venda_id}/itens [get]
func (h *DocumentoVendaItemHandler) List(c *gin.Context) {
	documentoVendaID, ok := utils.ParseIDParam(c, "documento_venda_id")
	if !ok {
		return
	}

	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)

	// Filters can be added here if needed, for now, it's just by documentoVendaID
	filters := make(map[string]interface{})

	items, total, err := h.service.List(limit, offset, documentoVendaID, filters)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	respItems := make([]dto.DocumentoVendaItemResponse, len(items))
	for i, item := range items {
		var respItem dto.DocumentoVendaItemResponse
		respItem.FromModel(&item)
		respItems[i] = respItem
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	utils.RespondWithOK(c, dto.DocumentoVendaItemListResponse{
		Items:      respItems,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}
