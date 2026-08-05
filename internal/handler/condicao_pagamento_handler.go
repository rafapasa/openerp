package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

type CondicaoPagamentoHandler struct {
	service service.CondicaoPagamentoService
}

func NewCondicaoPagamentoHandler(service service.CondicaoPagamentoService) *CondicaoPagamentoHandler {
	return &CondicaoPagamentoHandler{
		service: service,
	}
}

// @Summary      Cria uma nova condição de pagamento
// @Description  Cria uma nova condição de pagamento com base nos dados fornecidos.
// @Tags         Condições de Pagamento
// @Accept       json
// @Produce      json
// @Param        condicao  body      dto.CondicaoPagamentoRequest  true  "Dados para criar a condição de pagamento"
// @Success      201       {object}  dto.CondicaoPagamentoResponse
// @Failure      400       {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      500       {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /condicoes-pagamento [post]
func (h *CondicaoPagamentoHandler) Create(c *gin.Context) {
	var req dto.CondicaoPagamentoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	condicao, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.CondicaoPagamentoResponse
	if err := resp.FromModel(condicao); err != nil{
		utils.RespondWithErrorAny(c, err)
		return
	
	}
	utils.RespondWithCreated(c, resp)
}

// @Summary      Busca uma condição de pagamento por ID
// @Description  Retorna os detalhes de uma condição de pagamento específica.
// @Tags         Condições de Pagamento
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID da Condição de Pagamento"
// @Success      200  {object}  dto.CondicaoPagamentoResponse
// @Failure      404  {object}  utils.ErrorResponse "Condição de pagamento não encontrada"
// @Router       /condicoes-pagamento/{id} [get]
func (h *CondicaoPagamentoHandler) GetByID(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	condicao, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.CondicaoPagamentoResponse
	resp.FromModel(condicao)
	utils.RespondWithOK(c, resp)
}

// @Summary      Atualiza uma condição de pagamento
// @Description  Atualiza os dados de uma condição de pagamento existente.
// @Tags         Condições de Pagamento
// @Accept       json
// @Produce      json
// @Param        id        path      int                           true  "ID da Condição de Pagamento"
// @Param        condicao  body      dto.CondicaoPagamentoRequest  true  "Dados para atualizar a condição de pagamento"
// @Success      200       {object}  dto.CondicaoPagamentoResponse
// @Failure      400       {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      404       {object}  utils.ErrorResponse "Condição de pagamento não encontrada"
// @Router       /condicoes-pagamento/{id} [put]
func (h *CondicaoPagamentoHandler) Update(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.CondicaoPagamentoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.UpdatedBy = &userID

	condicao, err := h.service.Update(id, &req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.CondicaoPagamentoResponse
	resp.FromModel(condicao)
	utils.RespondWithOK(c, resp)
}

// @Summary      Exclui uma condição de pagamento
// @Description  Realiza a exclusão lógica de uma condição de pagamento.
// @Tags         Condições de Pagamento
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID da Condição de Pagamento"
// @Success      204  "Nenhum conteúdo"
// @Failure      400  {object}  utils.ErrorResponse "Erro ao excluir"
// @Router       /condicoes-pagamento/{id} [delete]
func (h *CondicaoPagamentoHandler) Delete(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	utils.RespondWithNoContent(c)
}

// @Summary      Lista as condições de pagamento
// @Description  Retorna uma lista paginada de condições de pagamento, com suporte a filtros.
// @Tags         Condições de Pagamento
// @Accept       json
// @Produce      json
// @Param        limit      query     int  false  "Número de registros por página"
// @Param        offset     query     int  false  "Offset para a paginação"
// @Param        descricao  query     string  false  "Filtrar por descrição"
// @Param        situacao   query     int  false  "Filtrar por situação (1=Ativo, 2=Inativo)"
// @Param        tipo_condicao query   int  false  "Filtrar por tipo de condição (0=À Vista, 1=À Prazo, 9=Sem Pagamento)"
// @Success      200        {object}  dto.CondicaoPagamentoListResponse
// @Failure      500        {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /condicoes-pagamento [get]
func (h *CondicaoPagamentoHandler) List(c *gin.Context) {
	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)

	filters := make(map[string]interface{})
	if descricao := utils.GetQueryString(c, "descricao", ""); descricao != "" {
		filters["cdpgt_descricao"] = descricao
	}
	if situacao := utils.GetQueryInt(c, "situacao", 0); situacao != 0 {
		filters["cdpgt_situacao"] = situacao
	}
	if tipoCondicao := utils.GetQueryInt(c, "tipo_condicao", -1); tipoCondicao != -1 { // -1 para permitir 0
		filters["cdpgt_tipocondicao"] = tipoCondicao
	}

	condicoes, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	items := make([]dto.CondicaoPagamentoResponse, len(condicoes))
	for i, condicao := range condicoes {
		var resp dto.CondicaoPagamentoResponse
		resp.FromModel(&condicao)
		items[i] = resp
	}

	totalPages := utils.CalculateTotalPages(int(total), limit)

	utils.RespondWithOK(c, dto.CondicaoPagamentoListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}
