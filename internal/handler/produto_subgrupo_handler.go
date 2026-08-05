package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

// ProdutoSubgrupoHandler gerencia as requisições HTTP para subgrupos de produto.
type ProdutoSubgrupoHandler struct {
	service service.ProdutoSubgrupoService
}

// NewProdutoSubgrupoHandler cria uma nova instância de ProdutoSubgrupoHandler.
func NewProdutoSubgrupoHandler(s service.ProdutoSubgrupoService) *ProdutoSubgrupoHandler {
	return &ProdutoSubgrupoHandler{
		service: s,
	}
}

// @Summary      Cria um novo subgrupo de produto
// @Description  Cria um novo subgrupo de produto com base nos dados fornecidos.
// @Tags         Produto - Subgrupos
// @Accept       json
// @Produce      json
// @Param        subgrupo  body      dto.ProdutoSubgrupoRequest  true  "Dados para criar o subgrupo de produto"
// @Success      201       {object}  dto.ProdutoSubgrupoResponse
// @Failure      400       {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      500       {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /produtos/subgrupos [post]
func (h *ProdutoSubgrupoHandler) Create(c *gin.Context) {
	var req dto.ProdutoSubgrupoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	subgrupo, err := h.service.Create(c, &req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.ProdutoSubgrupoResponse
	resp.FromModel(subgrupo)
	utils.RespondWithCreated(c, resp)
}

// @Summary      Busca um subgrupo de produto por ID
// @Description  Retorna os detalhes de um subgrupo de produto específico.
// @Tags         Produto - Subgrupos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Subgrupo de Produto"
// @Success      200  {object}  dto.ProdutoSubgrupoResponse
// @Failure      404  {object}  utils.ErrorResponse "Subgrupo de produto não encontrado"
// @Router       /produtos/subgrupos/{id} [get]
func (h *ProdutoSubgrupoHandler) GetByID(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	subgrupo, err := h.service.GetByID(c, id)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.ProdutoSubgrupoResponse
	resp.FromModel(subgrupo)
	utils.RespondWithOK(c, resp)
}

// @Summary      Atualiza um subgrupo de produto
// @Description  Atualiza os dados de um subgrupo de produto existente.
// @Tags         Produto - Subgrupos
// @Accept       json
// @Produce      json
// @Param        id        path      int                           true  "ID do Subgrupo de Produto"
// @Param        subgrupo  body      dto.ProdutoSubgrupoRequest  true  "Dados para atualizar o subgrupo de produto"
// @Success      200       {object}  dto.ProdutoSubgrupoResponse
// @Failure      400       {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      404       {object}  utils.ErrorResponse "Subgrupo de produto não encontrado"
// @Router       /produtos/subgrupos/{id} [put]
func (h *ProdutoSubgrupoHandler) Update(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.ProdutoSubgrupoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.UpdatedBy = &userID

	subgrupo, err := h.service.Update(c, id, &req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.ProdutoSubgrupoResponse
	resp.FromModel(subgrupo)
	utils.RespondWithOK(c, resp)
}

// @Summary      Exclui um subgrupo de produto
// @Description  Realiza a exclusão lógica de um subgrupo de produto.
// @Tags         Produto - Subgrupos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Subgrupo de Produto"
// @Success      204  "Nenhum conteúdo"
// @Failure      400  {object}  utils.ErrorResponse "Erro ao excluir"
// @Router       /produtos/subgrupos/{id} [delete]
func (h *ProdutoSubgrupoHandler) Delete(c *gin.Context) {
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

// @Summary      Lista os subgrupos de produto
// @Description  Retorna uma lista paginada de subgrupos de produto, com suporte a filtros.
// @Tags         Produto - Subgrupos
// @Accept       json
// @Produce      json
// @Param        limit      query     int     false  "Número de registros por página"
// @Param        offset     query     int     false  "Offset para a paginação"
// @Param        descricao  query     string  false  "Filtrar por descrição"
// @Param        situacao   query     int     false  "Filtrar por situação (1=Ativo, 2=Inativo)"
// @Success      200        {object}  dto.ProdutoSubgrupoListResponse
// @Router       /produtos/subgrupos [get]
func (h *ProdutoSubgrupoHandler) List(c *gin.Context) {
	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)

	filters := utils.QueryParamsToFilters(c)

	subgrupos, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	items := make([]dto.ProdutoSubgrupoResponse, len(subgrupos))
	for i, subgrupo := range subgrupos {
		var resp dto.ProdutoSubgrupoResponse
		resp.FromModel(&subgrupo)
		items[i] = resp
	}

	totalPages := utils.CalculateTotalPages(int(total), limit)

	utils.RespondWithOK(c, dto.ProdutoSubgrupoListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}
