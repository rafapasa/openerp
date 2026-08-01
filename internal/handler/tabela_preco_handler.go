package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

// TabelaPrecoHandler gerencia as requisições HTTP para tabelas de preço.
type TabelaPrecoHandler struct {
	service service.TabelaPrecoService
}

// NewTabelaPrecoHandler cria uma nova instância de TabelaPrecoHandler.
func NewTabelaPrecoHandler(s service.TabelaPrecoService) *TabelaPrecoHandler {
	return &TabelaPrecoHandler{
		service: s,
	}
}

// @Summary      Cria uma nova tabela de preço
// @Description  Cria uma nova tabela de preço com base nos dados fornecidos.
// @Tags         Produto - Tabela de Preços
// @Accept       json
// @Produce      json
// @Param        tabela  body      dto.TabelaPrecoRequest  true  "Dados para criar a tabela de preço"
// @Success      201     {object}  dto.TabelaPrecoResponse
// @Failure      400     {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      500     {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /tabelas-preco [post]
func (h *TabelaPrecoHandler) Create(c *gin.Context) {
	var req dto.TabelaPrecoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	tabela, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.TabelaPrecoResponse
	resp.FromModel(tabela)
	utils.RespondWithCreated(c, resp)
}

// @Summary      Busca uma tabela de preço por ID
// @Description  Retorna os detalhes de uma tabela de preço específica.
// @Tags         Produto - Tabela de Preços
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID da Tabela de Preço"
// @Success      200  {object}  dto.TabelaPrecoResponse
// @Failure      404  {object}  utils.ErrorResponse "Tabela de preço não encontrada"
// @Router       /tabelas-preco/{id} [get]
func (h *TabelaPrecoHandler) GetByID(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	tabela, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.TabelaPrecoResponse
	resp.FromModel(tabela)
	utils.RespondWithOK(c, resp)
}

// @Summary      Atualiza uma tabela de preço
// @Description  Atualiza os dados de uma tabela de preço existente.
// @Tags         Produto - Tabela de Preços
// @Accept       json
// @Produce      json
// @Param        id      path      int                       true  "ID da Tabela de Preço"
// @Param        tabela  body      dto.TabelaPrecoRequest  true  "Dados para atualizar a tabela de preço"
// @Success      200     {object}  dto.TabelaPrecoResponse
// @Failure      400     {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      404     {object}  utils.ErrorResponse "Tabela de preço não encontrada"
// @Router       /tabelas-preco/{id} [put]
func (h *TabelaPrecoHandler) Update(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.TabelaPrecoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.UpdatedBy = &userID

	tabela, err := h.service.Update(id, &req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.TabelaPrecoResponse
	resp.FromModel(tabela)
	utils.RespondWithOK(c, resp)
}

// @Summary      Exclui uma tabela de preço
// @Description  Realiza a exclusão lógica de uma tabela de preço.
// @Tags         Produto - Tabela de Preços
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID da Tabela de Preço"
// @Success      204  "Nenhum conteúdo"
// @Failure      400  {object}  utils.ErrorResponse "Erro ao excluir"
// @Router       /tabelas-preco/{id} [delete]
func (h *TabelaPrecoHandler) Delete(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary      Lista as tabelas de preço
// @Description  Retorna uma lista paginada de tabelas de preço, com suporte a filtros.
// @Tags         Produto - Tabela de Preços
// @Accept       json
// @Produce      json
// @Param        limit      query     int     false  "Número de registros por página"
// @Param        offset     query     int     false  "Offset para a paginação"
// @Param        descricao  query     string  false  "Filtrar por descrição"
// @Param        situacao   query     int     false  "Filtrar por situação (1=Ativo, 2=Inativo)"
// @Success      200        {object}  dto.TabelaPrecoListResponse
// @Router       /tabelas-preco [get]
func (h *TabelaPrecoHandler) List(c *gin.Context) {
	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)

	filters := utils.QueryParamsToFilters(c)

	tabelas, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	items := make([]dto.TabelaPrecoResponse, len(tabelas))
	for i, tabela := range tabelas {
		var resp dto.TabelaPrecoResponse
		resp.FromModel(&tabela)
		items[i] = resp
	}

	totalPages := utils.CalculateTotalPages(int(total), limit)

	utils.RespondWithOK(c, dto.TabelaPrecoListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}