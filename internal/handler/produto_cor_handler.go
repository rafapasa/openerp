package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

// ProdutoCorHandler gerencia as requisições HTTP para cores de produto.
type ProdutoCorHandler struct {
	service service.ProdutoCorService
}

// NewProdutoCorHandler cria uma nova instância de ProdutoCorHandler.
func NewProdutoCorHandler(s service.ProdutoCorService) *ProdutoCorHandler {
	return &ProdutoCorHandler{
		service: s,
	}
}

// @Summary      Cria uma nova cor de produto
// @Description  Cria uma nova cor de produto com base nos dados fornecidos.
// @Tags         Produto - Cores
// @Accept       json
// @Produce      json
// @Param        cor  body      dto.ProdutoCorRequest  true  "Dados para criar a cor de produto"
// @Success      201  {object}  dto.ProdutoCorResponse
// @Failure      400  {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      500  {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /produtos/cores [post]
func (h *ProdutoCorHandler) Create(c *gin.Context) {
	var req dto.ProdutoCorRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	cor, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	utils.RespondWithCreated(c, cor)
}

// @Summary      Busca uma cor de produto por ID
// @Description  Retorna os detalhes de uma cor de produto específica.
// @Tags         Produto - Cores
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID da Cor de Produto"
// @Success      200  {object}  dto.ProdutoCorResponse
// @Failure      404  {object}  utils.ErrorResponse "Cor de produto não encontrada"
// @Router       /produtos/cores/{id} [get]
func (h *ProdutoCorHandler) GetByID(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	cor, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	utils.RespondWithOK(c, cor)
}

// @Summary      Atualiza uma cor de produto
// @Description  Atualiza os dados de uma cor de produto existente.
// @Tags         Produto - Cores
// @Accept       json
// @Produce      json
// @Param        id    path      int                      true  "ID da Cor de Produto"
// @Param        cor   body      dto.ProdutoCorRequest  true  "Dados para atualizar a cor de produto"
// @Success      200   {object}  dto.ProdutoCorResponse
// @Failure      400   {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      404   {object}  utils.ErrorResponse "Cor de produto não encontrada"
// @Router       /produtos/cores/{id} [put]
func (h *ProdutoCorHandler) Update(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.ProdutoCorRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.UpdatedBy = &userID

	cor, err := h.service.Update(id, &req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	utils.RespondWithOK(c, cor)
}

// @Summary      Exclui uma cor de produto
// @Description  Realiza a exclusão lógica de uma cor de produto.
// @Tags         Produto - Cores
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID da Cor de Produto"
// @Success      204  "Nenhum conteúdo"
// @Failure      400  {object}  utils.ErrorResponse "Erro ao excluir"
// @Router       /produtos/cores/{id} [delete]
func (h *ProdutoCorHandler) Delete(c *gin.Context) {
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

// @Summary      Lista as cores de produto
// @Description  Retorna uma lista paginada de cores de produto, com suporte a filtros.
// @Tags         Produto - Cores
// @Accept       json
// @Produce      json
// @Param        limit      query     int     false  "Número de registros por página"
// @Param        offset     query     int     false  "Offset para a paginação"
// @Param        descricao  query     string  false  "Filtrar por descrição"
// @Param        situacao   query     int     false  "Filtrar por situação (1=Ativo, 2=Inativo)"
// @Success      200        {object}  dto.ProdutoCorListResponse
// @Router       /produtos/cores [get]
func (h *ProdutoCorHandler) List(c *gin.Context) {
	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)

	filters := utils.QueryParamsToFilters(c)

	cores, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	items := make([]dto.ProdutoCorResponse, len(cores))
	for i, cor := range cores {
		items[i] = cor
	}

	totalPages := utils.CalculateTotalPages(int(total), limit)

	utils.RespondWithOK(c, dto.ProdutoCorListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}