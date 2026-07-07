package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// TYPES
// ============================================================

type EntidadeLimiteCreditoHandler struct {
	service *service.EntidadeLimiteCreditoService
}

// ============================================================
// CONSTRUCTOR
// ============================================================

func NewEntidadeLimiteCreditoHandler(service *service.EntidadeLimiteCreditoService) *EntidadeLimiteCreditoHandler {
	return &EntidadeLimiteCreditoHandler{
		service: service,
	}
}

// ============================================================
// HANDLERS
// ============================================================

// Create cria um novo limite de crédito.
//
//	@Summary		Cria um novo limite de crédito
//	@Description	Cadastra um novo limite de crédito para ser usado nas entidades
//	@Tags			EntidadeLimiteCredito
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.EntidadeLimiteCreditoRequest	true	"Dados do limite de crédito"
//	@Success		201		{object}	dto.EntidadeLimiteCreditoResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/limites-credito [post]
func (h *EntidadeLimiteCreditoHandler) Create(c *gin.Context) {
	var req dto.EntidadeLimiteCreditoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	limite, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	var resp dto.EntidadeLimiteCreditoResponse
	resp.FromModel(limite)
	utils.RespondWithCreated(c, resp)
}

// List lista todos os limites de crédito.
//
//	@Summary		Lista limites de crédito
//	@Description	Retorna uma lista paginada de limites de crédito
//	@Tags			EntidadeLimiteCredito
//	@Produce		json
//	@Param			limit		query		int		false	"Limite de registros"	default(10)
//	@Param			offset		query		int		false	"Offset para paginação"	default(0)
//	@Param			descricao	query		string	false	"Filtrar por descrição"
//	@Success		200			{object}	dto.EntidadeLimiteCreditoListResponse
//	@Failure		500			{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/limites-credito [get]
func (h *EntidadeLimiteCreditoHandler) List(c *gin.Context) {
	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)

	filters := make(map[string]interface{})
	if descricao := utils.GetQueryString(c, "descricao", ""); descricao != "" {
		filters["descricao"] = descricao
	}

	limites, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithInternalError(c, err.Error())
		return
	}

	items := make([]dto.EntidadeLimiteCreditoResponse, len(limites))
	for i, limite := range limites {
		var resp dto.EntidadeLimiteCreditoResponse
		resp.FromModel(&limite)
		items[i] = resp
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	utils.RespondWithOK(c, dto.EntidadeLimiteCreditoListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// GetByID busca um limite de crédito por ID.
//
//	@Summary		Busca limite de crédito por ID
//	@Description	Retorna os dados de um limite de crédito específico
//	@Tags			EntidadeLimiteCredito
//	@Produce		json
//	@Param			id	path		int	true	"ID do Limite de Crédito"
//	@Success		200	{object}	dto.EntidadeLimiteCreditoResponse
//	@Failure		404	{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/limites-credito/{id} [get]
func (h *EntidadeLimiteCreditoHandler) GetByID(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	limite, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondWithNotFoundError(c, err.Error())
		return
	}

	var resp dto.EntidadeLimiteCreditoResponse
	resp.FromModel(limite)
	utils.RespondWithOK(c, resp)
}

// Update atualiza um limite de crédito.
//
//	@Summary		Atualiza um limite de crédito
//	@Description	Atualiza a descrição e o valor de um limite de crédito
//	@Tags			EntidadeLimiteCredito
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"ID do Limite de Crédito"
//	@Param			request	body		dto.EntidadeLimiteCreditoRequest	true	"Dados para atualização"
//	@Success		200		{object}	dto.EntidadeLimiteCreditoResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/limites-credito/{id} [put]
func (h *EntidadeLimiteCreditoHandler) Update(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.EntidadeLimiteCreditoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.UpdatedBy = &userID

	limite, err := h.service.Update(id, &req)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	var resp dto.EntidadeLimiteCreditoResponse
	resp.FromModel(limite)
	utils.RespondWithOK(c, resp)
}

// Delete exclui um limite de crédito.
//
//	@Summary		Exclui um limite de crédito
//	@Description	Realiza a exclusão lógica de um limite de crédito
//	@Tags			EntidadeLimiteCredito
//	@Param			id	path	int	true	"ID do Limite de Crédito"
//	@Success		204	"No Content"
//	@Failure		404	{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/limites-credito/{id} [delete]
func (h *EntidadeLimiteCreditoHandler) Delete(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.RespondWithNotFoundError(c, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
