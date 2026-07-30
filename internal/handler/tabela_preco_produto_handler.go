package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

type TabelaPrecoProdutoHandler struct {
	service service.TabelaPrecoProdutoService
}

func NewTabelaPrecoProdutoHandler(service service.TabelaPrecoProdutoService) *TabelaPrecoProdutoHandler {
	return &TabelaPrecoProdutoHandler{
		service: service,
	}
}

func (h *TabelaPrecoProdutoHandler) getIDs(c *gin.Context) (int, int, bool) {
	tabelaID := utils.GetQueryInt(c, "id", 0)
	if tabelaID == 0 {
		utils.RespondWithValidationError(c, "ID da tabela de preço inválido")
		return 0, 0, false
	}
	item := utils.GetQueryInt(c, "item", 0)
	if item == 0 {
		utils.RespondWithValidationError(c, "item inválido")
		return 0, 0, false
	}
	return tabelaID, item, true
}

// @Summary      Adiciona um produto a uma tabela de preço
// @Description  Adiciona um novo produto com seus valores a uma tabela de preço existente.
// @Tags         Produto - Tabela de Preços
// @Accept       json
// @Produce      json
// @Param        tabela_id  path      int                                  true  "ID da Tabela de Preço"
// @Param        produto    body      dto.TabelaPrecoProdutoRequest        true  "Dados do produto na tabela"
// @Success      201        {object}  dto.TabelaPrecoProdutoResponse
// @Failure      400        {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Router       /tabelas-preco/{tabela_id}/produtos [post]
func (h *TabelaPrecoProdutoHandler) Create(c *gin.Context) {
	tabelaID, ok := utils.ParseIDParam(c, "tabela_id")
	if !ok {
		return
	}

	var req dto.TabelaPrecoProdutoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	req.TabelaPrecoID = tabelaID
	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	item, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	var resp dto.TabelaPrecoProdutoResponse
	resp.FromModel(item)
	utils.RespondWithCreated(c, resp)
}

// @Summary      Busca um produto em uma tabela de preço
// @Description  Retorna os detalhes de um produto específico dentro de uma tabela de preço.
// @Tags         Produto - Tabela de Preços
// @Accept       json
// @Produce      json
// @Param        tabela_id  path      int  true  "ID da Tabela de Preço"
// @Param        item_id    path      int  true  "ID do Item (Produto na Tabela)"
// @Success      200        {object}  dto.TabelaPrecoProdutoResponse
// @Failure      404        {object}  utils.ErrorResponse "Item não encontrado"
// @Router       /tabelas-preco/{tabela_id}/produtos/{item_id} [get]
func (h *TabelaPrecoProdutoHandler) GetByID(c *gin.Context) {
	tabelaID, item, ok := h.getIDs(c)
	if !ok {
		return
	}

	tabelaItem, err := h.service.GetByID(tabelaID, item)
	if err != nil {
		utils.RespondWithNotFoundError(c, err.Error())
		return
	}

	var resp dto.TabelaPrecoProdutoResponse
	resp.FromModel(tabelaItem)
	utils.RespondWithOK(c, resp)
}

// @Summary      Atualiza um produto em uma tabela de preço
// @Description  Atualiza os dados de um produto em uma tabela de preço.
// @Tags         Produto - Tabela de Preços
// @Accept       json
// @Produce      json
// @Param        tabela_id  path      int                                  true  "ID da Tabela de Preço"
// @Param        item_id    path      int                                  true  "ID do Item (Produto na Tabela)"
// @Param        produto    body      dto.TabelaPrecoProdutoRequest        true  "Dados para atualizar"
// @Success      200        {object}  dto.TabelaPrecoProdutoResponse
// @Failure      400        {object}  utils.ErrorResponse "Erro de validação"
// @Failure      404        {object}  utils.ErrorResponse "Item não encontrado"
// @Router       /tabelas-preco/{tabela_id}/produtos/{item_id} [put]
func (h *TabelaPrecoProdutoHandler) Update(c *gin.Context) {
	tabelaID, item, ok := h.getIDs(c)
	if !ok {
		return
	}

	var req dto.TabelaPrecoProdutoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	req.TabelaPrecoID = tabelaID
	userID := middleware.GetUserID(c)
	req.UpdatedBy = &userID

	updateItem, err := h.service.Update(tabelaID, item, &req)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	var resp dto.TabelaPrecoProdutoResponse
	resp.FromModel(updateItem)
	utils.RespondWithOK(c, resp)
}

// @Summary      Remove um produto de uma tabela de preço
// @Description  Realiza a exclusão lógica de um produto de uma tabela de preço.
// @Tags         Produto - Tabela de Preços
// @Accept       json
// @Produce      json
// @Param        tabela_id  path      int  true  "ID da Tabela de Preço"
// @Param        item_id    path      int  true  "ID do Item (Produto na Tabela)"
// @Success      204        "Nenhum conteúdo"
// @Failure      400        {object}  utils.ErrorResponse "Erro ao excluir"
// @Router       /tabelas-preco/{tabela_id}/produtos/{item_id} [delete]
func (h *TabelaPrecoProdutoHandler) Delete(c *gin.Context) {
	tabelaID, item, ok := h.getIDs(c)
	if !ok {
		return
	}

	if err := h.service.Delete(tabelaID, item); err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	utils.RespondWithNoContent(c)
}

// @Summary      Lista os produtos de uma tabela de preço
// @Description  Retorna uma lista paginada de produtos de uma tabela de preço.
// @Tags         Produto - Tabela de Preços
// @Accept       json
// @Produce      json
// @Param        tabela_id     path      int  true  "ID da Tabela de Preço"
// @Param        limit         query     int  false "Número de registros por página"
// @Param        offset        query     int  false "Offset para a paginação"
// @Param        produto_nome  query     string false "Filtrar por nome do produto"
// @Success      200           {object}  dto.TabelaPrecoProdutoListResponse
// @Router       /tabelas-preco/{tabela_id}/produtos [get]
func (h *TabelaPrecoProdutoHandler) List(c *gin.Context) {
	tabelaID, ok := utils.ParseIDParam(c, "tabela_id")
	if !ok {
		return
	}

	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)
	filters := make(map[string]interface{})
	if nome := utils.GetQueryString(c, "produto_nome", ""); nome != "" {
		filters["produto_nome"] = nome
	}

	items, total, err := h.service.List(tabelaID, limit, offset, filters)
	if err != nil {
		utils.RespondWithInternalError(c, err.Error())
		return
	}

	respItems := make([]dto.TabelaPrecoProdutoResponse, len(items))
	for i, item := range items {
		respItems[i].FromModel(&item)
	}

	utils.RespondWithOK(c, dto.TabelaPrecoProdutoListResponse{
		Items:      respItems,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: utils.CalculateTotalPages(int(total), limit),
	})
}
