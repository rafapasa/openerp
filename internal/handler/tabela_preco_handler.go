package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

type TabelaPrecoHandler struct {
	service *service.TabelaPrecoService
}

func NewTabelaPrecoHandler(service *service.TabelaPrecoService) *TabelaPrecoHandler {
	return &TabelaPrecoHandler{
		service: service,
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
		utils.RespondWithValidationError(c, err.Error())
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
		utils.RespondWithNotFoundError(c, err.Error())
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
// @Param        id      path      int                     true  "ID da Tabela de Preço"
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
		utils.RespondWithValidationError(c, err.Error())
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
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	utils.RespondWithNoContent(c)
}

// @Summary      Lista as tabelas de preço
// @Description  Retorna uma lista paginada de tabelas de preço, com suporte a filtros.
// @Tags         Produto - Tabela de Preços
// @Accept       json
// @Produce      json
// @Param        limit     query     int  false  "Número de registros por página"
// @Param        offset    query     int  false  "Offset para a paginação"
// @Param        descricao query     string  false  "Filtrar por descrição"
// @Param        situacao  query     int  false  "Filtrar por situação (1=Ativo, 2=Inativo)"
// @Success      200       {object}  dto.TabelaPrecoListResponse
// @Router       /tabelas-preco [get]
func (h *TabelaPrecoHandler) List(c *gin.Context) {
	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)

	filters := make(map[string]interface{})
	if descricao := utils.GetQueryString(c, "descricao", ""); descricao != "" {
		filters["tbp_descricao"] = descricao
	}
	if situacaoStr := utils.GetQueryString(c, "situacao", ""); situacaoStr != "" {
		if situacao, err := strconv.Atoi(situacaoStr); err == nil {
			filters["tbp_situacao"] = situacao
		}
	}

	tabelas, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithInternalError(c, err.Error())
		return
	}

	items := make([]dto.TabelaPrecoResponse, len(tabelas))
	for i, tabela := range tabelas {
		items[i].FromModel(&tabela)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	utils.RespondWithOK(c, dto.TabelaPrecoListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}
