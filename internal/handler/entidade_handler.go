package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

// EntidadeHandler gerencia as requisições HTTP para entidades.
type EntidadeHandler struct {
	service service.EntidadeService
}

// NewEntidadeHandler cria uma nova instância de EntidadeHandler.
func NewEntidadeHandler(s service.EntidadeService) *EntidadeHandler {
	return &EntidadeHandler{
		service: s,
	}
}

// @Summary      Cria uma nova entidade
// @Description  Cadastra uma nova entidade (cliente/fornecedor)
// @Tags         Entidades
// @Accept       json
// @Produce      json
// @Param        request  body      dto.EntidadeRequest  true  "Dados da entidade"
// @Success      201      {object}  dto.EntidadeResponse
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Security     BearerAuth
// @Router       /entidades [post]
func (h *EntidadeHandler) Create(c *gin.Context) {
	var req dto.EntidadeRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	entidade, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.EntidadeResponse
	resp.FromModel(entidade)
	utils.RespondWithCreated(c, resp)
}

// @Summary      Busca entidade por ID
// @Description  Retorna os dados de uma entidade específica
// @Tags         Entidades
// @Produce      json
// @Param        id   path      int  true  "ID da entidade"
// @Success      200  {object}  dto.EntidadeResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /entidades/{entidade_id} [get]
func (h *EntidadeHandler) GetByID(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "entidade_id")
	if !ok {
		return
	}

	entidade, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.EntidadeResponse
	resp.FromModel(entidade)
	utils.RespondWithOK(c, resp)
}

// @Summary      Busca entidade por documento
// @Description  Retorna os dados de uma entidade pelo CPF/CNPJ
// @Tags         Entidades
// @Produce      json
// @Param        documento  path      string  true  "CPF/CNPJ da entidade"
// @Success      200        {object}  dto.EntidadeResponse
// @Failure      404        {object}  map[string]string
// @Failure      500        {object}  map[string]string
// @Security     BearerAuth
// @Router       /entidades/documento/{documento} [get]
func (h *EntidadeHandler) GetByDocumento(c *gin.Context) {
	documento := c.Param("documento")
	if documento == "" {
		utils.RespondWithErrorAny(c, apperrors.NewValidationError("Documento não pode ser vazio."))
		return
	}

	entidade, err := h.service.GetByDocumento(documento)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.EntidadeResponse
	resp.FromModel(entidade)
	utils.RespondWithOK(c, resp)
}

// @Summary      Atualiza uma entidade
// @Description  Atualiza os dados de uma entidade existente
// @Tags         Entidades
// @Accept       json
// @Produce      json
// @Param        id       path      int                  true  "ID da entidade"
// @Param        request  body      dto.EntidadeRequest  true  "Dados atualizados"
// @Success      200      {object}  dto.EntidadeResponse
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Security     BearerAuth
// @Router       /entidades/{entidade_id} [put]
func (h *EntidadeHandler) Update(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "entidade_id")
	if !ok {
		return
	}

	var req dto.EntidadeRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.UpdatedBy = &userID

	entidade, err := h.service.Update(id, &req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.EntidadeResponse
	resp.FromModel(entidade)
	utils.RespondWithOK(c, resp)
}

// @Summary      Exclui uma entidade
// @Description  Realiza a exclusão lógica de uma entidade
// @Tags         Entidades
// @Produce      json
// @Param        id   path      int  true  "ID da entidade"
// @Success      204  "No Content"
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /entidades/{entidade_id} [delete]
func (h *EntidadeHandler) Delete(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "entidade_id")
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary      Lista entidades
// @Description  Retorna uma lista paginada de entidades
// @Tags         Entidades
// @Produce      json
// @Param        limit      query     int     false  "Limite de registros"  default(10)
// @Param        offset     query     int     false  "Offset para paginação"  default(0)
// @Param        nome       query     string  false  "Filtrar por nome (razao social)"
// @Param        documento  query     string  false  "Filtrar por documento"
// @Param        tipo_pessoa query    int     false  "Filtrar por tipo (1-Física, 2-Jurídica)"
// @Param        situacao   query     int     false  "Filtrar por situação (1-Ativo, 2-Inativo)"
// @Success      200        {object}  dto.EntidadeListResponse
// @Failure      500        {object}  map[string]string
// @Security     BearerAuth
// @Router       /entidades [get]
func (h *EntidadeHandler) List(c *gin.Context) {
	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)

	filters := utils.QueryParamsToFilters(c)

	entidades, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	items := make([]dto.EntidadeResponse, len(entidades))
	for i, entidade := range entidades {
		var resp dto.EntidadeResponse
		resp.FromModel(&entidade)
		items[i] = resp
	}

	totalPages := utils.CalculateTotalPages(int(total), limit)

	utils.RespondWithOK(c, dto.EntidadeListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}
