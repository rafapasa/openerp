// internal/handler/entidade_handler.go (ENXUTO!)
package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

type EntidadeHandler struct {
	service service.EntidadeService
}

func NewEntidadeHandler(service service.EntidadeService) *EntidadeHandler {
	return &EntidadeHandler{service: service}
}

// Create cria uma nova entidade
func (h *EntidadeHandler) Create(c *gin.Context) {
	var req dto.EntidadeRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	// ✅ Service recebe DTO e faz a conversão
	entidade, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.EntidadeResponse
	resp.FromModel(entidade)
	utils.RespondWithCreated(c, resp)
}

// Update atualiza uma entidade
func (h *EntidadeHandler) Update(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	var req dto.EntidadeRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.UpdatedBy = &userID

	// ✅ Service recebe DTO e faz a conversão
	entidade, err := h.service.Update(id, &req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.EntidadeResponse
	resp.FromModel(entidade)
	utils.RespondWithOK(c, resp)
}

// GetByID busca uma entidade por ID
func (h *EntidadeHandler) GetByID(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
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

// Delete exclui uma entidade
func (h *EntidadeHandler) Delete(c *gin.Context) {
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

// List lista entidades com paginação
func (h *EntidadeHandler) List(c *gin.Context) {
	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)
	filters := h.buildFilters(c)

	entidades, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithInternalError(c, err.Error())
		return
	}

	items := make([]dto.EntidadeResponse, len(entidades))
	for i, entidade := range entidades {
		items[i].FromModel(&entidade)
	}

	utils.RespondWithOK(c, dto.EntidadeListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	})
}

// buildFilters constrói filtros da query string
func (h *EntidadeHandler) buildFilters(c *gin.Context) map[string]interface{} {
	filters := make(map[string]interface{})

	if nome := c.Query("nome"); nome != "" {
		filters["razao_social"] = nome
	}
	if documento := c.Query("documento"); documento != "" {
		filters["documento"] = documento
	}
	if tipoPessoa := c.Query("tipo_pessoa"); tipoPessoa != "" {
		if val, err := strconv.Atoi(tipoPessoa); err == nil {
			filters["tipo_pessoa"] = val
		}
	}
	if situacao := c.Query("situacao"); situacao != "" {
		if val, err := strconv.Atoi(situacao); err == nil {
			filters["situacao"] = val
		}
	}

	return filters
}
