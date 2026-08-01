package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

// ProdutoTamanhoHandler gerencia as requisições HTTP para tamanhos de produto.
type ProdutoTamanhoHandler struct {
	service service.ProdutoTamanhoService
}

// NewProdutoTamanhoHandler cria uma nova instância de ProdutoTamanhoHandler.
func NewProdutoTamanhoHandler(s service.ProdutoTamanhoService) *ProdutoTamanhoHandler {
	return &ProdutoTamanhoHandler{
		service: s,
	}
}

// @Summary      Cria um novo tamanho de produto
// @Description  Cria um novo tamanho de produto com base nos dados fornecidos.
// @Tags         Produto - Tamanhos
// @Accept       json
// @Produce      json
// @Param        tamanho  body      dto.ProdutoTamanhoRequest  true  "Dados para criar o tamanho de produto"
// @Success      201      {object}  dto.ProdutoTamanhoResponse
// @Failure      400      {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      500      {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /produtos/tamanhos [post]
func (h *ProdutoTamanhoHandler) Create(c *gin.Context) {
	var req dto.ProdutoTamanhoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	tamanho, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	utils.RespondWithCreated(c, tamanho)
}

// @Summary      Busca um tamanho de produto por ID
// @Description  Retorna os detalhes de um tamanho de produto específico.
// @Tags         Produto - Tamanhos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Tamanho de Produto"
// @Success      200  {object}  dto.ProdutoTamanhoResponse
// @Failure      404  {object}  utils.ErrorResponse "Tamanho de produto não encontrado"
// @Router       /produtos/tamanhos/{id} [get]
func (h *ProdutoTamanhoHandler) GetByID(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	tamanho, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	utils.RespondWithOK(c, tamanho)
}

// @Summary      Atualiza um tamanho de produto
// @Description  Atualiza os dados de um tamanho de produto existente.
// @Tags         Produto - Tamanhos
// @Accept       json
// @Produce      json
// @Param        id       path      int                          true  "ID do Tamanho de Produto"
// @Param        tamanho  body      dto.ProdutoTamanhoRequest  true  "Dados para atualizar o tamanho de produto"
// @Success      200      {object}  dto.ProdutoTamanhoResponse
// @Failure      400      {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      404      {object}  utils.ErrorResponse "Tamanho de produto não encontrado"
// @Router       /produtos/tamanhos/{id} [put]
func (h *ProdutoTamanhoHandler) Update(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.ProdutoTamanhoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.UpdatedBy = &userID

	tamanho, err := h.service.Update(id, &req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	utils.RespondWithOK(c, tamanho)
}

// @Summary      Exclui um tamanho de produto
// @Description  Realiza a exclusão lógica de um tamanho de produto.
// @Tags         Produto - Tamanhos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Tamanho de Produto"
// @Success      204  "Nenhum conteúdo"
// @Failure      400  {object}  utils.ErrorResponse "Erro ao excluir"
// @Router       /produtos/tamanhos/{id} [delete]
func (h *ProdutoTamanhoHandler) Delete(c *gin.Context) {
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

// @Summary      Lista os tamanhos de produto
// @Description  Retorna uma lista paginada de tamanhos de produto, com suporte a filtros.
// @Tags         Produto - Tamanhos
// @Accept       json
// @Produce      json
// @Param        limit      query     int     false  "Número de registros por página"
// @Param        offset     query     int     false  "Offset para a paginação"
// @Param        descricao  query     string  false  "Filtrar por descrição"
// @Param        situacao   query     int     false  "Filtrar por situação (1=Ativo, 2=Inativo)"
// @Success      200        {object}  dto.ProdutoTamanhoListResponse
// @Router       /produtos/tamanhos [get]
func (h *ProdutoTamanhoHandler) List(c *gin.Context) {
	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)

	filters := utils.QueryParamsToFilters(c)

	tamanhos, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	totalPages := utils.CalculateTotalPages(int(total), limit)

	utils.RespondWithOK(c, dto.ProdutoTamanhoListResponse{
		Items:      tamanhos,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}
