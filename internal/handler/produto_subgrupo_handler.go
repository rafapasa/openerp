package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// TYPES
// ============================================================

type ProdutoSubgrupoHandler struct {
	service service.ProdutoSubgrupoService
}

// ============================================================
// CONSTRUCTOR
// ============================================================

func NewProdutoSubgrupoHandler(service service.ProdutoSubgrupoService) *ProdutoSubgrupoHandler {
	return &ProdutoSubgrupoHandler{
		service: service,
	}
}

// ============================================================
// HANDLERS
// ============================================================

// @Summary      Cria um novo subgrupo de produto
// @Description  Cria um novo subgrupo de produto com base nos dados fornecidos.
// @Tags         Produto - Subgrupos
// @Accept       json
// @Produce      json
// @Param        subgrupo  body      dto.ProdutoSubgrupoRequest  true  "Dados para criar o subgrupo de produto"
// @Success      201       {object}  dto.ProdutoSubgrupoResponse
// @Failure      400       {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      500       {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /produto-subgrupos [post]
// Create cria um novo subgrupo de produto.
func (h *ProdutoSubgrupoHandler) Create(c *gin.Context) {
	var req dto.ProdutoSubgrupoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	subgrupo, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
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
// @Router       /produto-subgrupos/{id} [get]
// GetByID busca um subgrupo de produto por ID.
func (h *ProdutoSubgrupoHandler) GetByID(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	subgrupo, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondWithNotFoundError(c, err.Error())
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
// @Param        id        path      int                         true  "ID do Subgrupo de Produto"
// @Param        subgrupo  body      dto.ProdutoSubgrupoRequest  true  "Dados para atualizar o subgrupo de produto"
// @Success      200       {object}  dto.ProdutoSubgrupoResponse
// @Failure      400       {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      404       {object}  utils.ErrorResponse "Subgrupo de produto não encontrado"
// @Router       /produto-subgrupos/{id} [put]
// Update atualiza um subgrupo de produto.
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

	subgrupo, err := h.service.Update(id, &req)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
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
// @Router       /produto-subgrupos/{id} [delete]
// Delete exclui um subgrupo de produto.
func (h *ProdutoSubgrupoHandler) Delete(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	utils.RespondWithNoContent(c)
}

// @Summary      Lista os subgrupos de produto
// @Description  Retorna uma lista paginada de subgrupos de produto, com suporte a filtros.
// @Tags         Produto - Subgrupos
// @Accept       json
// @Produce      json
// @Param        limit     query     int  false  "Número de registros por página"
// @Param        offset    query     int  false  "Offset para a paginação"
// @Param        descricao query     string  false  "Filtrar por descrição"
// @Param        situacao  query     int  false  "Filtrar por situação (1=Ativo, 2=Inativo)"
// @Success      200       {object}  dto.ProdutoSubgrupoListResponse
// @Router       /produto-subgrupos [get]
// List lista todos os subgrupos de produto.
func (h *ProdutoSubgrupoHandler) List(c *gin.Context) {
	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)
	descricao := utils.GetQueryString(c, "descricao", "")
	filters := make(map[string]interface{})
	if descricao != "" {
		filters["prosg_descricao"] = descricao
	}
	situacao := utils.GetQueryInt(c, "situacao", 0)
	if situacao != 0 {
		filters["prosg_situacao"] = situacao
	}

	subgrupos, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithInternalError(c, err.Error())
		return
	}

	items := make([]dto.ProdutoSubgrupoResponse, len(subgrupos))
	for i, subgrupo := range subgrupos {
		items[i].FromModel(&subgrupo)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	utils.RespondWithOK(c, dto.ProdutoSubgrupoListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}
