package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

// EmpresaFilialHandler gerencia as requisições HTTP para filiais de empresa.
type EmpresaFilialHandler struct {
	service service.EmpresaFilialService
}

// NewEmpresaFilialHandler cria uma nova instância de EmpresaFilialHandler.
func NewEmpresaFilialHandler(s service.EmpresaFilialService) *EmpresaFilialHandler {
	return &EmpresaFilialHandler{service: s}
}

// CreateEmpresaFilial godoc
// @Summary Cria uma nova filial de empresa
// @Description Cria uma nova filial de empresa com os dados fornecidos
// @Tags Empresa - Filiais
// @Accept json
// @Produce json
// @Param filial body dto.EmpresaFilialRequest true "Dados da Filial de Empresa"
// @Success 201 {object} dto.EmpresaFilialResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 409 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /empresas-filiais [post]
func (h *EmpresaFilialHandler) CreateEmpresaFilial(c *gin.Context) {
	var req dto.EmpresaFilialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithErrorAny(c, apperrors.NewValidationError(err.Error()))
		return
	}

	filial, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	c.JSON(http.StatusCreated, filial)
}

// GetEmpresaFilialByID godoc
// @Summary Busca uma filial de empresa pelo ID
// @Description Retorna uma filial de empresa específica pelo seu ID
// @Tags Empresa - Filiais
// @Produce json
// @Param id path int true "ID da Filial de Empresa"
// @Success 200 {object} dto.EmpresaFilialResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /empresas-filiais/{id} [get]
func (h *EmpresaFilialHandler) GetEmpresaFilialByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondWithErrorAny(c, apperrors.NewValidationError("ID inválido."))
		return
	}

	filial, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	c.JSON(http.StatusOK, filial)
}

// UpdateEmpresaFilial godoc
// @Summary Atualiza uma filial de empresa existente
// @Description Atualiza uma filial de empresa com os dados fornecidos pelo ID
// @Tags Empresa - Filiais
// @Accept json
// @Produce json
// @Param id path int true "ID da Filial de Empresa"
// @Param filial body dto.EmpresaFilialRequest true "Dados da Filial de Empresa para atualização"
// @Success 200 {object} dto.EmpresaFilialResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 409 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /empresas-filiais/{id} [put]
func (h *EmpresaFilialHandler) UpdateEmpresaFilial(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondWithErrorAny(c, apperrors.NewValidationError("ID inválido."))
		return
	}

	var req dto.EmpresaFilialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithErrorAny(c, apperrors.NewValidationError(err.Error()))
		return
	}

	filial, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	c.JSON(http.StatusOK, filial)
}

// DeleteEmpresaFilial godoc
// @Summary Exclui uma filial de empresa
// @Description Exclui logicamente uma filial de empresa pelo ID
// @Tags Empresa - Filiais
// @Produce json
// @Param id path int true "ID da Filial de Empresa"
// @Success 204 "No Content"
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 409 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /empresas-filiais/{id} [delete]
func (h *EmpresaFilialHandler) DeleteEmpresaFilial(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondWithErrorAny(c, apperrors.NewValidationError("ID inválido."))
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListEmpresasFiliais godoc
// @Summary Lista todas as filiais de empresa
// @Description Retorna uma lista paginada de filiais de empresa
// @Tags Empresa - Filiais
// @Produce json
// @Param limit query int false "Limite de itens por página" default(10)
// @Param offset query int false "Número de itens a pular" default(0)
// @Success 200 {object} dto.EmpresaFilialListResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /empresas-filiais [get]
func (h *EmpresaFilialHandler) ListEmpresasFiliais(c *gin.Context) {
	limit := utils.GetQueryInt(c, "limit", 0)
	offset := utils.GetQueryInt(c, "offset", 0)
	filters := utils.QueryParamsToFilters(c)

	filiais, total, err := h.service.List(c.Request.Context(), limit, offset, filters)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.EmpresaFilialListResponse{
		Items:      filiais,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: (int(total) + limit - 1) / limit,
	})
}
