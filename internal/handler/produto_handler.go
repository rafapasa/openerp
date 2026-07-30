package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

type ProdutoHandler struct {
	service service.ProdutoService
}

func NewProdutoHandler(service service.ProdutoService) *ProdutoHandler {
	return &ProdutoHandler{
		service: service,
	}
}

// @Summary      Cria um novo produto
// @Description  Cria um novo produto com base nos dados fornecidos.
// @Tags         Produtos
// @Accept       json
// @Produce      json
// @Param        produto  body      dto.ProdutoRequest  true  "Dados para criar o produto"
// @Success      201      {object}  dto.ProdutoResponse
// @Failure      400      {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      500      {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /produtos [post]
func (h *ProdutoHandler) Create(c *gin.Context) {
	var req dto.ProdutoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}
	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID
	produto, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}
	var resp dto.ProdutoResponse
	resp.FromModel(produto)
	utils.RespondWithCreated(c, resp)
}

// @Summary      Busca um produto por ID
// @Description  Retorna os detalhes de um produto específico.
// @Tags         Produtos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Produto"
// @Success      200  {object}  dto.ProdutoResponse
// @Failure      404  {object}  utils.ErrorResponse "Produto não encontrado"
// @Router       /produtos/{id} [get]
func (h *ProdutoHandler) GetByID(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}
	produto, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondWithNotFoundError(c, err.Error())
		return
	}
	var resp dto.ProdutoResponse
	resp.FromModel(produto)
	utils.RespondWithOK(c, resp)
}

// @Summary      Atualiza um produto
// @Description  Atualiza os dados de um produto existente.
// @Tags         Produtos
// @Accept       json
// @Produce      json
// @Param        id       path      int                 true  "ID do Produto"
// @Param        produto  body      dto.ProdutoRequest  true  "Dados para atualizar o produto"
// @Success      200      {object}  dto.ProdutoResponse
// @Failure      400      {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      404      {object}  utils.ErrorResponse "Produto não encontrado"
// @Router       /produtos/{id} [put]
func (h *ProdutoHandler) Update(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}
	var req dto.ProdutoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}
	userID := middleware.GetUserID(c)
	req.UpdatedBy = &userID
	produto, err := h.service.Update(id, &req)
	if err != nil {
		utils.RespondWithNotFoundError(c, err.Error())
		return
	}
	var resp dto.ProdutoResponse
	resp.FromModel(produto)
	utils.RespondWithOK(c, resp)
}

// @Summary      Exclui um produto
// @Description  Realiza a exclusão lógica de um produto.
// @Tags         Produtos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Produto"
// @Success      204  "Nenhum conteúdo"
// @Failure      400  {object}  utils.ErrorResponse "Erro ao excluir"
// @Router       /produtos/{id} [delete]
func (h *ProdutoHandler) Delete(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}
	err := h.service.Delete(id)
	if err != nil {
		utils.RespondWithNotFoundError(c, err.Error())
		return
	}
	utils.RespondWithNoContent(c)
}

// @Summary      Lista os produtos
// @Description  Retorna uma lista paginada de produtos, com suporte a filtros.
// @Tags         Produtos
// @Accept       json
// @Produce      json
// @Param        limit      query     int  false  "Número de registros por página"
// @Param        offset     query     int  false  "Offset para a paginação"
// @Param        nome       query     string  false  "Filtrar por nome"
// @Param        codigo     query     int  false  "Filtrar por código"
// @Param        referencia query     string  false  "Filtrar por referência"
// @Param        situacao   query     int  false  "Filtrar por situação (1=Ativo, 2=Inativo)"
// @Success      200        {object}  dto.ProdutoListResponse
// @Router       /produtos [get]
func (h *ProdutoHandler) List(c *gin.Context) {
	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)

	filters := make(map[string]interface{})
	if nome := utils.GetQueryString(c, "nome", ""); nome != "" {
		filters["pro_nome"] = nome
	}
	if codigo := utils.GetQueryString(c, "codigo", ""); codigo != "" {
		if val, err := strconv.Atoi(codigo); err == nil {
			filters["pro_codigo"] = val
		}
	}
	if referencia := utils.GetQueryString(c, "referencia", ""); referencia != "" {
		filters["pro_referencia"] = referencia
	}
	if situacao := utils.GetQueryString(c, "situacao", ""); situacao != "" {
		if val, err := strconv.Atoi(situacao); err == nil {
			filters["pro_situacao"] = val
		}
	}

	produtos, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithInternalError(c, err.Error())
		return
	}

	items := make([]dto.ProdutoResponse, len(produtos))
	for i, produto := range produtos {
		var resp dto.ProdutoResponse
		resp.FromModel(&produto)
		items[i] = resp
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	utils.RespondWithOK(c, dto.ProdutoListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}
