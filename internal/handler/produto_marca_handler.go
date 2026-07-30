package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

type ProdutoMarcaHandler struct {
	service service.ProdutoMarcaService
}

func NewProdutoMarcaHandler(service service.ProdutoMarcaService) *ProdutoMarcaHandler {
	return &ProdutoMarcaHandler{
		service: service,
	}
}

// @Summary      Cria uma nova marca de produto
// @Description  Cria uma nova marca de produto com base nos dados fornecidos.
// @Tags         Produto - Marcas
// @Accept       json
// @Produce      json
// @Param        marca  body      dto.ProdutoMarcaRequest  true  "Dados para criar a marca de produto"
// @Success      201    {object}  dto.ProdutoMarcaResponse
// @Failure      400    {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      500    {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /produto-marcas [post]
func (h *ProdutoMarcaHandler) Create(c *gin.Context) {
	var req dto.ProdutoMarcaRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	marca, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	var resp dto.ProdutoMarcaResponse
	resp.FromModel(marca)
	utils.RespondWithCreated(c, resp)
}

// @Summary      Busca uma marca de produto por ID
// @Description  Retorna os detalhes de uma marca de produto específica.
// @Tags         Produto - Marcas
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID da Marca de Produto"
// @Success      200  {object}  dto.ProdutoMarcaResponse
// @Failure      404  {object}  utils.ErrorResponse "Marca de produto não encontrada"
// @Router       /produto-marcas/{id} [get]
func (h *ProdutoMarcaHandler) GetByID(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	marca, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondWithNotFoundError(c, err.Error())
		return
	}

	var resp dto.ProdutoMarcaResponse
	resp.FromModel(marca)
	utils.RespondWithOK(c, resp)
}

// @Summary      Atualiza uma marca de produto
// @Description  Atualiza os dados de uma marca de produto existente.
// @Tags         Produto - Marcas
// @Accept       json
// @Produce      json
// @Param        id     path      int                      true  "ID da Marca de Produto"
// @Param        marca  body      dto.ProdutoMarcaRequest  true  "Dados para atualizar a marca de produto"
// @Success      200    {object}  dto.ProdutoMarcaResponse
// @Failure      400    {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      404    {object}  utils.ErrorResponse "Marca de produto não encontrada"
// @Router       /produto-marcas/{id} [put]
func (h *ProdutoMarcaHandler) Update(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.ProdutoMarcaRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.UpdatedBy = &userID

	marca, err := h.service.Update(id, &req)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	var resp dto.ProdutoMarcaResponse
	resp.FromModel(marca)
	utils.RespondWithOK(c, resp)
}

// @Summary      Exclui uma marca de produto
// @Description  Realiza a exclusão lógica de uma marca de produto.
// @Tags         Produto - Marcas
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID da Marca de Produto"
// @Success      204  "Nenhum conteúdo"
// @Failure      400  {object}  utils.ErrorResponse "Erro ao excluir"
// @Router       /produto-marcas/{id} [delete]
func (h *ProdutoMarcaHandler) Delete(c *gin.Context) {
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

// @Summary      Lista as marcas de produto
// @Description  Retorna uma lista paginada de marcas de produto, com suporte a filtros.
// @Tags         Produto - Marcas
// @Accept       json
// @Produce      json
// @Param        limit     query     int  false  "Número de registros por página"
// @Param        offset    query     int  false  "Offset para a paginação"
// @Param        descricao query     string  false  "Filtrar por descrição"
// @Param        situacao  query     int  false  "Filtrar por situação (1=Ativo, 2=Inativo)"
// @Success      200       {object}  dto.ProdutoMarcaListResponse
// @Router       /produto-marcas [get]
func (h *ProdutoMarcaHandler) List(c *gin.Context) {
	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)
	filters := make(map[string]interface{})
	descricao := utils.GetQueryString(c, "descricao", "")
	if descricao != "" {
		filters["prom_descricao"] = descricao
	}
	situacao := utils.GetQueryString(c, "situacao", "")
	if situacao != "" {
		filters["prom_situacao"] = situacao
	}

	marcas, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithInternalError(c, err.Error())
		return
	}

	items := make([]dto.ProdutoMarcaResponse, len(marcas))
	for i, marca := range marcas {
		items[i].FromModel(&marca)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	utils.RespondWithOK(c, dto.ProdutoMarcaListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}
