package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

// ProdutoModeloHandler gerencia as requisições HTTP para modelos de produto.
type ProdutoModeloHandler struct {
	service service.ProdutoModeloService
}

// NewProdutoModeloHandler cria uma nova instância de ProdutoModeloHandler.
func NewProdutoModeloHandler(s service.ProdutoModeloService) *ProdutoModeloHandler {
	return &ProdutoModeloHandler{
		service: s,
	}
}

// @Summary      Cria um novo modelo de produto
// @Description  Cria um novo modelo de produto com base nos dados fornecidos.
// @Tags         Produto - Modelos
// @Accept       json
// @Produce      json
// @Param        modelo  body      dto.ProdutoModeloRequest  true  "Dados para criar o modelo de produto"
// @Success      201     {object}  dto.ProdutoModeloResponse
// @Failure      400     {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      500     {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /produtos/modelos [post]
func (h *ProdutoModeloHandler) Create(c *gin.Context) {
	var req dto.ProdutoModeloRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	modelo, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.ProdutoModeloResponse
	resp.FromModel(modelo)
	utils.RespondWithCreated(c, resp)
}

// @Summary      Busca um modelo de produto por ID
// @Description  Retorna os detalhes de um modelo de produto específico.
// @Tags         Produto - Modelos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Modelo de Produto"
// @Success      200  {object}  dto.ProdutoModeloResponse
// @Failure      404  {object}  utils.ErrorResponse "Modelo de produto não encontrado"
// @Router       /produtos/modelos/{id} [get]
func (h *ProdutoModeloHandler) GetByID(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	modelo, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.ProdutoModeloResponse
	resp.FromModel(modelo)
	utils.RespondWithOK(c, resp)
}

// @Summary      Atualiza um modelo de produto
// @Description  Atualiza os dados de um modelo de produto existente.
// @Tags         Produto - Modelos
// @Accept       json
// @Produce      json
// @Param        id      path      int                         true  "ID do Modelo de Produto"
// @Param        modelo  body      dto.ProdutoModeloRequest  true  "Dados para atualizar o modelo de produto"
// @Success      200     {object}  dto.ProdutoModeloResponse
// @Failure      400     {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      404     {object}  utils.ErrorResponse "Modelo de produto não encontrado"
// @Router       /produtos/modelos/{id} [put]
func (h *ProdutoModeloHandler) Update(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.ProdutoModeloRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.UpdatedBy = &userID

	modelo, err := h.service.Update(id, &req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.ProdutoModeloResponse
	resp.FromModel(modelo)
	utils.RespondWithOK(c, resp)
}

// @Summary      Exclui um modelo de produto
// @Description  Realiza a exclusão lógica de um modelo de produto.
// @Tags         Produto - Modelos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Modelo de Produto"
// @Success      204  "Nenhum conteúdo"
// @Failure      400  {object}  utils.ErrorResponse "Erro ao excluir"
// @Router       /produtos/modelos/{id} [delete]
func (h *ProdutoModeloHandler) Delete(c *gin.Context) {
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

// @Summary      Lista os modelos de produto
// @Description  Retorna uma lista paginada de modelos de produto, com suporte a filtros.
// @Tags         Produto - Modelos
// @Accept       json
// @Produce      json
// @Param        limit      query     int     false  "Número de registros por página"
// @Param        offset     query     int     false  "Offset para a paginação"
// @Param        descricao  query     string  false  "Filtrar por descrição"
// @Param        situacao   query     int     false  "Filtrar por situação (1=Ativo, 2=Inativo)"
// @Success      200        {object}  dto.ProdutoModeloListResponse
// @Router       /produtos/modelos [get]
func (h *ProdutoModeloHandler) List(c *gin.Context) {
	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)

	filters := utils.QueryParamsToFilters(c)

	modelos, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	items := make([]dto.ProdutoModeloResponse, len(modelos))
	for i, modelo := range modelos {
		var resp dto.ProdutoModeloResponse
		resp.FromModel(&modelo)
		items[i] = resp
	}

	totalPages := utils.CalculateTotalPages(int(total), limit)

	utils.RespondWithOK(c, dto.ProdutoModeloListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}