package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

type ProdutoGrupoHandler struct {
	produtoGrupoService *service.ProdutoGrupoService
}

func NewProdutoGrupoHandler(produtoGrupoService *service.ProdutoGrupoService) *ProdutoGrupoHandler {
	return &ProdutoGrupoHandler{
		produtoGrupoService: produtoGrupoService,
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
// @Router       /produto-grupos [post]
func (h *ProdutoGrupoHandler) Create(c *gin.Context) {
	var req dto.ProdutoGrupoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	grupo, err := h.produtoGrupoService.Create(&req)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	utils.RespondWithCreated(c, grupo)
}

// @Summary      Busca um grupo de produto por ID
// @Description  Retorna os detalhes de um grupo de produto específico.
// @Tags         Produto - Grupos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Grupo de Produto"
// @Success      200  {object}  dto.ProdutoGrupoResponse
// @Failure      404  {object}  utils.ErrorResponse "Grupo de produto não encontrado"
// @Router       /produto-grupos/{id} [get]
func (h *ProdutoGrupoHandler) GetByID(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	grupo, err := h.produtoGrupoService.GetByID(id)
	if err != nil {
		utils.RespondWithNotFoundError(c, err.Error())
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
// @Param        id     path      int                      true  "ID do Grupo de Produto"
// @Param        grupo  body      dto.ProdutoGrupoRequest  true  "Dados para atualizar o grupo de produto"
// @Success      200    {object}  dto.ProdutoGrupoResponse
// @Failure      400    {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      404    {object}  utils.ErrorResponse "Grupo de produto não encontrado"
// @Router       /produto-grupos/{id} [put]
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

	grupo, err := h.produtoGrupoService.Update(id, &req)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	utils.RespondWithOK(c, grupo)
}

// @Summary      Exclui um grupo de produto
// @Description  Realiza a exclusão lógica de um grupo de produto.
// @Tags         Produto - Grupos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Grupo de Produto"
// @Success      204  "Nenhum conteúdo"
// @Failure      400  {object}  utils.ErrorResponse "Erro ao excluir"
// @Router       /produto-grupos/{id} [delete]
func (h *ProdutoGrupoHandler) Delete(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	err := h.produtoGrupoService.Delete(id)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	utils.RespondWithNoContent(c)
}

// @Summary      Lista os grupos de produto
// @Description  Retorna uma lista paginada de grupos de produto, com suporte a filtros.
// @Tags         Produto - Grupos
// @Accept       json
// @Produce      json
// @Param        limit                query     int  false  "Número de registros por página"
// @Param        offset               query     int  false  "Offset para a paginação"
// @Param        descricao            query     string  false  "Filtrar por descrição"
// @Param        situacao             query     int  false  "Filtrar por situação (1=Ativo, 2=Inativo)"
// @Param        visivel_frente_caixa query     int  false  "Filtrar por visibilidade no PDV (1=Sim, 2=Não)"
// @Success      200                  {object}  dto.ProdutoGrupoListResponse
// @Router       /produto-grupos [get]
func (h *ProdutoGrupoHandler) List(c *gin.Context) {
	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)

	filters := make(map[string]any)
	if descricao := utils.GetQueryString(c, "descricao", ""); descricao != "" {
		filters["prog_descricao"] = descricao
	}
	if situacao := utils.GetQueryInt(c, "situacao", 0); situacao != 0 {
		filters["prog_situacao"] = situacao
	}
	if visivelfrentecaixa := utils.GetQueryInt(c, "visivel_frente_caixa", 0); visivelfrentecaixa != 0 {
		filters["prog_visivelfrentecaixa"] = visivelfrentecaixa
	}
	if agenda := utils.GetQueryInt(c, "agenda", 0); agenda != 0 {
		filters["prog_agenda"] = agenda
	}

	produtoGrupos, total, err := h.produtoGrupoService.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithInternalError(c, err.Error())
		return
	}

	items := make([]dto.ProdutoGrupoResponse, len(produtoGrupos))
	for i, grupo := range produtoGrupos {
		var resp dto.ProdutoGrupoResponse
		resp.FromModel(&grupo)
		items[i] = resp
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	utils.RespondWithOK(c, dto.ProdutoGrupoListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}
