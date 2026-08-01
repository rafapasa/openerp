package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

// ProcessoHandler gerencia as requisições HTTP para processos.
type ProcessoHandler struct {
	service service.ProcessoService
}

// NewProcessoHandler cria uma nova instância de ProcessoHandler.
func NewProcessoHandler(s service.ProcessoService) *ProcessoHandler {
	return &ProcessoHandler{
		service: s,
	}
}

// @Summary      Cria um novo processo
// @Description  Cria um novo processo com base nos dados fornecidos.
// @Tags         Processos
// @Accept       json
// @Produce      json
// @Param        processo  body      dto.ProcessoRequest  true  "Dados para criar o processo"
// @Success      201       {object}  dto.ProcessoResponse
// @Failure      400       {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      500       {object}  utils.ErrorResponse "Erro interno do servidor"
// @Router       /processos [post]
func (h *ProcessoHandler) Create(c *gin.Context) {
	var req dto.ProcessoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	processo, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.ProcessoResponse
	resp.FromModel(processo)
	utils.RespondWithCreated(c, resp)
}

// @Summary      Busca um processo por ID
// @Description  Retorna os detalhes de um processo específico.
// @Tags         Processos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Processo"
// @Success      200  {object}  dto.ProcessoResponse
// @Failure      404  {object}  utils.ErrorResponse "Processo não encontrado"
// @Router       /processos/{id} [get]
func (h *ProcessoHandler) GetByID(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	processo, err := h.service.FindByID(id)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.ProcessoResponse
	resp.FromModel(processo)
	utils.RespondWithOK(c, resp)
}

// @Summary      Atualiza um processo
// @Description  Atualiza os dados de um processo existente.
// @Tags         Processos
// @Accept       json
// @Produce      json
// @Param        id        path      int                      true  "ID do Processo"
// @Param        processo  body      dto.ProcessoRequest  true  "Dados para atualizar o processo"
// @Success      200       {object}  dto.ProcessoResponse
// @Failure      400       {object}  utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure      404       {object}  utils.ErrorResponse "Processo não encontrado"
// @Router       /processos/{id} [put]
func (h *ProcessoHandler) Update(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.ProcessoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.UpdatedBy = &userID

	processo, err := h.service.Update(id, &req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.ProcessoResponse
	resp.FromModel(processo)
	utils.RespondWithOK(c, resp)
}

// @Summary      Exclui um processo
// @Description  Realiza a exclusão lógica de um processo.
// @Tags         Processos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do Processo"
// @Success      204  "Nenhum conteúdo"
// @Failure      400  {object}  utils.ErrorResponse "Erro ao excluir"
// @Router       /processos/{id} [delete]
func (h *ProcessoHandler) Delete(c *gin.Context) {
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

// @Summary      Lista os processos
// @Description  Retorna uma lista paginada de processos, com suporte a filtros.
// @Tags         Processos
// @Accept       json
// @Produce      json
// @Param        limit      query     int     false  "Número de registros por página"
// @Param        offset     query     int     false  "Offset para a paginação"
// @Param        codigo     query     int     false  "Filtrar por código"
// @Param        descricao  query     string  false  "Filtrar por descrição"
// @Param        tipo_operacao query  int     false  "Filtrar por tipo de operação (0=Entrada, 1=Saída)"
// @Param        situacao   query     int     false  "Filtrar por situação (1=Ativo, 2=Inativo)"
// @Success      200        {object}  dto.ProcessoListResponse
// @Router       /processos [get]
func (h *ProcessoHandler) List(c *gin.Context) {
	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)

	filters := utils.QueryParamsToFilters(c)

	processos, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	items := make([]dto.ProcessoResponse, len(processos))
	for i, processo := range processos {
		var resp dto.ProcessoResponse
		resp.FromModel(&processo)
		items[i] = resp
	}

	totalPages := utils.CalculateTotalPages(int(total), limit)

	utils.RespondWithOK(c, dto.ProcessoListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// @Summary      Busca um processo por código
// @Description  Retorna os detalhes de um processo específico pelo código.
// @Tags         Processos
// @Accept       json
// @Produce      json
// @Param        codigo  path      int  true  "Código do Processo"
// @Success      200     {object}  dto.ProcessoResponse
// @Failure      400     {object}  utils.ErrorResponse "Código inválido"
// @Failure      404     {object}  utils.ErrorResponse "Processo não encontrado"
// @Router       /processos/codigo/{codigo} [get]
func (h *ProcessoHandler) GetByCodigo(c *gin.Context) {
	codigo, ok := utils.ParseIDParam(c, "codigo")
	if !ok {
		return
	}

	processo, err := h.service.FindByCodigo(codigo)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.ProcessoResponse
	resp.FromModel(processo)
	utils.RespondWithOK(c, resp)
}
