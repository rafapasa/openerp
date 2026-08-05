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

// GrupoUsuarioHandler gerencia as requisições HTTP para grupos de usuário.
type GrupoUsuarioHandler struct {
	service service.GrupoUsuarioService
}

// NewGrupoUsuarioHandler cria uma nova instância de GrupoUsuarioHandler.
func NewGrupoUsuarioHandler(s service.GrupoUsuarioService) *GrupoUsuarioHandler {
	return &GrupoUsuarioHandler{service: s}
}

// CreateGrupoUsuario godoc
// @Summary Cria um novo grupo de usuário
// @Description Cria um novo grupo de usuário com os dados fornecidos
// @Tags Grupo de Usuário
// @Accept json
// @Produce json
// @Param grupo body dto.GrupoUsuarioRequest true "Dados do Grupo de Usuário"
// @Success 201 {object} dto.GrupoUsuarioResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 409 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /grupos-usuario [post]
func (h *GrupoUsuarioHandler) CreateGrupoUsuario(c *gin.Context) {
	var req dto.GrupoUsuarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithErrorAny(c, apperrors.NewValidationError(err.Error()))
		return
	}

	grupo, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	c.JSON(http.StatusCreated, grupo)
}

// GetGrupoUsuarioByID godoc
// @Summary Busca um grupo de usuário pelo ID
// @Description Retorna um grupo de usuário específico pelo seu ID
// @Tags Grupo de Usuário
// @Produce json
// @Param id path int true "ID do Grupo de Usuário"
// @Success 200 {object} dto.GrupoUsuarioResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /grupos-usuario/{id} [get]
func (h *GrupoUsuarioHandler) GetGrupoUsuarioByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondWithErrorAny(c, apperrors.NewValidationError("ID inválido."))
		return
	}

	grupo, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	c.JSON(http.StatusOK, grupo)
}

// UpdateGrupoUsuario godoc
// @Summary Atualiza um grupo de usuário existente
// @Description Atualiza um grupo de usuário com os dados fornecidos pelo ID
// @Tags Grupo de Usuário
// @Accept json
// @Produce json
// @Param id path int true "ID do Grupo de Usuário"
// @Param grupo body dto.GrupoUsuarioRequest true "Dados do Grupo de Usuário para atualização"
// @Success 200 {object} dto.GrupoUsuarioResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 409 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /grupos-usuario/{id} [put]
func (h *GrupoUsuarioHandler) UpdateGrupoUsuario(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondWithErrorAny(c, apperrors.NewValidationError("ID inválido."))
		return
	}

	var req dto.GrupoUsuarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithErrorAny(c, apperrors.NewValidationError(err.Error()))
		return
	}

	grupo, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	c.JSON(http.StatusOK, grupo)
}

// DeleteGrupoUsuario godoc
// @Summary Exclui um grupo de usuário
// @Description Exclui logicamente um grupo de usuário pelo ID
// @Tags Grupo de Usuário
// @Produce json
// @Param id path int true "ID do Grupo de Usuário"
// @Success 204 "No Content"
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 409 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /grupos-usuario/{id} [delete]
func (h *GrupoUsuarioHandler) DeleteGrupoUsuario(c *gin.Context) {
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

// ListGrupoUsuarios godoc
// @Summary Lista todos os grupos de usuário
// @Description Retorna uma lista paginada de grupos de usuário
// @Tags Grupo de Usuário
// @Produce json
// @Param limit query int false "Limite de itens por página" default(10)
// @Param offset query int false "Número de itens a pular" default(0)
// @Success 200 {object} dto.GrupoUsuarioListResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /grupos-usuario [get]
func (h *GrupoUsuarioHandler) ListGrupoUsuarios(c *gin.Context) {
	limit := utils.GetQueryInt(c, "limit", 0)
	offset := utils.GetQueryInt(c, "offset", 0)
	filters := utils.QueryParamsToFilters(c)

	grupos, total, err := h.service.List(c.Request.Context(), limit, offset, filters)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.GrupoUsuarioListResponse{
		Items:      grupos,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: (int(total) + limit - 1) / limit,
	})
}
