package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

// ProdutoGrupoHandler gerencia as requisições HTTP para grupos de produto.
type ProdutoGrupoHandler struct {
	service service.ProdutoGrupoService
}

// NewProdutoGrupoHandler cria uma nova instância de ProdutoGrupoHandler.
func NewProdutoGrupoHandler(s service.ProdutoGrupoService) *ProdutoGrupoHandler {
	return &ProdutoGrupoHandler{
		service: s,
	}
}

// @Summary      Cria um novo grupo de produto
// @Description  Cria um novo grupo de produto com base nos dados fornecidos.
// @Tags         Produto - Grupos
// @Accept       json
// @Produce      json
// @Param        grupo  body      dto.ProdutoGrupoRequest  true  "Dados para criar o grupo de produto"
// @Success      201    {object}  dto.ProdutoGrupoResponse
// @Failure      400    {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      500    {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /produtos/grupos [post]
func (h *ProdutoGrupoHandler) Create(c *gin.Context) {
	var req dto.ProdutoGrupoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	grupo, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.ProdutoGrupoResponse
	resp.FromModel(grupo)
	utils.RespondWithCreated(c, resp)
}

// @Summary      Busca um grupo de produto por ID
// @Description  Retorna os detalhes de um grupo de produto específico.
// @Tags         Produto - Grupos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Grupo de Produto"
// @Success      200  {object}  dto.ProdutoGrupoResponse
// @Failure      404  {object}  utils.ErrorResponse "Grupo de produto não encontrado"
// @Router       /produtos/grupos/{id} [get]
func (h *ProdutoGrupoHandler) GetByID(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	grupo, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.ProdutoGrupoResponse
	resp.FromModel(grupo)
	utils.RespondWithOK(c, resp)
}

// @Summary      Atualiza um grupo de produto
// @Description  Atualiza os dados de um grupo de produto existente.
// @Tags         Produto - Grupos
// @Accept       json
// @Produce      json
// @Param        id     path      int                        true  "ID do Grupo de Produto"
// @Param        grupo  body      dto.ProdutoGrupoRequest  true  "Dados para atualizar o grupo de produto"
// @Success      200    {object}  dto.ProdutoGrupoResponse
// @Failure      400    {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      404    {object}  utils.ErrorResponse "Grupo de produto não encontrado"
// @Router       /produtos/grupos/{id} [put]
func (h *ProdutoGrupoHandler) Update(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.ProdutoGrupoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.UpdatedBy = &userID

	grupo, err := h.service.Update(id, &req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.ProdutoGrupoResponse
	resp.FromModel(grupo)
	utils.RespondWithOK(c, resp)
}

// @Summary      Exclui um grupo de produto
// @Description  Realiza a exclusão lógica de um grupo de produto.
// @Tags         Produto - Grupos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Grupo de Produto"
// @Success      204  "Nenhum conteúdo"
// @Failure      400  {object}  utils.ErrorResponse "Erro ao excluir"
// @Router       /produtos/grupos/{id} [delete]
func (h *ProdutoGrupoHandler) Delete(c *gin.Context) {
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

// @Summary      Lista os grupos de produto
// @Description  Retorna uma lista paginada de grupos de produto, com suporte a filtros.
// @Tags         Produto - Grupos
// @Accept       json
// @Produce      json
// @Param        limit      query     int     false  "Número de registros por página"
// @Param        offset     query     int     false  "Offset para a paginação"
// @Param        descricao  query     string  false  "Filtrar por descrição"
// @Param        situacao   query     int     false  "Filtrar por situação (1=Ativo, 2=Inativo)"
// @Success      200        {object}  dto.ProdutoGrupoListResponse
// @Router       /produtos/grupos [get]
func (h *ProdutoGrupoHandler) List(c *gin.Context) {
	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)

	filters := utils.QueryParamsToFilters(c)

	grupos, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	items := make([]dto.ProdutoGrupoResponse, len(grupos))
	for i, grupo := range grupos {
		var resp dto.ProdutoGrupoResponse
		resp.FromModel(&grupo)
		items[i] = resp
	}

	totalPages := utils.CalculateTotalPages(int(total), limit)

	utils.RespondWithOK(c, dto.ProdutoGrupoListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}
