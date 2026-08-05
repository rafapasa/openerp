package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

// OperacaoFiscalHandler gerencia as requisições HTTP para operações fiscais.
type OperacaoFiscalHandler struct {
	service service.OperacaoFiscalService
}

// NewOperacaoFiscalHandler cria uma nova instância de OperacaoFiscalHandler.
func NewOperacaoFiscalHandler(s service.OperacaoFiscalService) *OperacaoFiscalHandler {
	return &OperacaoFiscalHandler{
		service: s,
	}
}

// @Summary      Cria uma nova operação fiscal
// @Description  Cria uma nova operação fiscal com base nos dados fornecidos.
// @Tags         Operações Fiscais
// @Accept       json
// @Produce      json
// @Param        operacao  body      dto.OperacaoFiscalRequest  true  "Dados para criar a operação fiscal"
// @Success      201       {object}  dto.OperacaoFiscalResponse
// @Failure      400       {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      500       {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /operacoes-fiscais [post]
func (h *OperacaoFiscalHandler) Create(c *gin.Context) {
	var req dto.OperacaoFiscalRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	operacao, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.OperacaoFiscalResponse
	resp.FromModel(operacao)
	utils.RespondWithCreated(c, resp)
}

// @Summary      Busca uma operação fiscal por ID
// @Description  Retorna os detalhes de uma operação fiscal específica.
// @Tags         Operações Fiscais
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID da Operação Fiscal"
// @Success      200  {object}  dto.OperacaoFiscalResponse
// @Failure      404  {object}  utils.ErrorResponse "Operação fiscal não encontrada"
// @Router       /operacoes-fiscais/{id} [get]
func (h *OperacaoFiscalHandler) GetByID(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	operacao, err := h.service.FindByID(id)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.OperacaoFiscalResponse
	resp.FromModel(operacao)
	utils.RespondWithOK(c, resp)
}

// @Summary      Atualiza uma operação fiscal
// @Description  Atualiza os dados de uma operação fiscal existente.
// @Tags         Operações Fiscais
// @Accept       json
// @Produce      json
// @Param        id        path      int                        true  "ID da Operação Fiscal"
// @Param        operacao  body      dto.OperacaoFiscalRequest  true  "Dados para atualizar a operação fiscal"
// @Success      200       {object}  dto.OperacaoFiscalResponse
// @Failure      400       {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      404       {object}  utils.ErrorResponse "Operação fiscal não encontrada"
// @Router       /operacoes-fiscais/{id} [put]
func (h *OperacaoFiscalHandler) Update(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.OperacaoFiscalRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.UpdatedBy = &userID

	operacao, err := h.service.Update(id, &req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.OperacaoFiscalResponse
	resp.FromModel(operacao)
	utils.RespondWithOK(c, resp)
}

// @Summary      Exclui uma operação fiscal
// @Description  Realiza a exclusão lógica de uma operação fiscal.
// @Tags         Operações Fiscais
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID da Operação Fiscal"
// @Success      204  "Nenhum conteúdo"
// @Failure      400  {object}  utils.ErrorResponse "Erro ao excluir"
// @Router       /operacoes-fiscais/{id} [delete]
func (h *OperacaoFiscalHandler) Delete(c *gin.Context) {
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

// @Summary      Lista as operações fiscais
// @Description  Retorna uma lista paginada de operações fiscais, com suporte a filtros.
// @Tags         Operações Fiscais
// @Accept       json
// @Produce      json
// @Param        limit      query     int     false  "Número de registros por página"
// @Param        offset     query     int     false  "Offset para a paginação"
// @Param        cfop       query     string  false  "Filtrar por CFOP"
// @Param        descricao  query     string  false  "Filtrar por descrição"
// @Param        empresa_filial_id query int  false  "Filtrar por ID da Empresa Filial"
// @Success      200        {object}  dto.OperacaoFiscalListResponse
// @Router       /operacoes-fiscais [get]
func (h *OperacaoFiscalHandler) List(c *gin.Context) {
	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)

	filters := utils.QueryParamsToFilters(c)

	operacoes, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	items := make([]dto.OperacaoFiscalResponse, len(operacoes))
	for i, operacao := range operacoes {
		var resp dto.OperacaoFiscalResponse
		resp.FromModel(&operacao)
		items[i] = resp
	}

	totalPages := utils.CalculateTotalPages(int(total), limit)

	utils.RespondWithOK(c, dto.OperacaoFiscalListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// @Summary      Busca operações fiscais por CFOP
// @Description  Retorna operações fiscais filtradas por CFOP.
// @Tags         Operações Fiscais
// @Accept       json
// @Produce      json
// @Param        cfop  path      string  true  "CFOP da Operação Fiscal"
// @Success      200   {array}   dto.OperacaoFiscalResponse
// @Failure      400   {object}  utils.ErrorResponse "CFOP inválido"
// @Failure      404   {object}  utils.ErrorResponse "Operação fiscal não encontrada"
// @Router       /operacoes-fiscais/cfop/{cfop} [get]
func (h *OperacaoFiscalHandler) GetByCFOP(c *gin.Context) {
	cfop := c.Param("cfop")
	if cfop == "" {
		utils.RespondWithErrorAny(c, apperrors.NewValidationError("CFOP não pode ser vazio."))
		return
	}

	operacoes, err := h.service.FindByCFOP(cfop)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	items := make([]dto.OperacaoFiscalResponse, len(operacoes))
	for i, op := range operacoes {
		var resp dto.OperacaoFiscalResponse
		resp.FromModel(&op)
		items[i] = resp
	}
	utils.RespondWithOK(c, items)
}

// @Summary      Busca operações fiscais por filial
// @Description  Retorna operações fiscais filtradas por ID da filial.
// @Tags         Operações Fiscais
// @Accept       json
// @Produce      json
// @Param        filial_id  path      int  true  "ID da Filial"
// @Success      200        {array}   dto.OperacaoFiscalResponse
// @Failure      400        {object}  utils.ErrorResponse "ID da filial inválido"
// @Failure      404        {object}  utils.ErrorResponse "Operação fiscal não encontrada"
// @Router       /operacoes-fiscais/filial/{filial_id} [get]
func (h *OperacaoFiscalHandler) GetByFilial(c *gin.Context) {
	filialID, ok := utils.ParseIDParam(c, "filial_id")
	if !ok {
		return
	}

	operacoes, err := h.service.FindByEmpresaFilialID(filialID)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	items := make([]dto.OperacaoFiscalResponse, len(operacoes))
	for i, op := range operacoes {
		var resp dto.OperacaoFiscalResponse
		resp.FromModel(&op)
		items[i] = resp
	}
	utils.RespondWithOK(c, items)
}
