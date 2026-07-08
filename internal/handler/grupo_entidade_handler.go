package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// TYPES
// ============================================================

type GrupoEntidadeHandler struct {
	service *service.GrupoEntidadeService
}

// ============================================================
// CONSTRUCTOR
// ============================================================

func NewGrupoEntidadeHandler(service *service.GrupoEntidadeService) *GrupoEntidadeHandler {
	return &GrupoEntidadeHandler{
		service: service,
	}
}

// ============================================================
// HANDLERS
// ============================================================

// Create cria um novo grupo de entidade
//	@Summary		Cria um novo grupo de entidade
//	@Description	Cadastra um novo grupo para ser usado nas entidades
//	@Tags			GrupoEntidades
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.GrupoEntidadeRequest	true	"Dados do grupo"
//	@Success		201		{object}	dto.GrupoEntidadeResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/grupos-entidades [post]
func (h *GrupoEntidadeHandler) Create(c *gin.Context) {
	var req dto.GrupoEntidadeRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	grupo, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	var resp dto.GrupoEntidadeResponse
	resp.FromModel(grupo)
	utils.RespondWithCreated(c, resp)
}

// List lista todos os grupos de entidade
//	@Summary		Lista grupos de entidade
//	@Description	Retorna uma lista paginada de grupos de entidade
//	@Tags			GrupoEntidades
//	@Produce		json
//	@Param			limit	query		int		false	"Limite de registros"	default(10)
//	@Param			offset	query		int		false	"Offset para paginação"	default(0)
//	@Param			nome	query		string	false	"Filtrar por nome"
//	@Success		200		{object}	dto.GrupoEntidadeListResponse
//	@Failure		500		{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/grupos-entidades [get]
func (h *GrupoEntidadeHandler) List(c *gin.Context) {
	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)

	filters := make(map[string]interface{})
	if nome := utils.GetQueryString(c, "nome", ""); nome != "" {
		filters["nome"] = nome
	}

	grupos, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithInternalError(c, err.Error())
		return
	}

	items := make([]dto.GrupoEntidadeResponse, len(grupos))
	for i, grupo := range grupos {
		var resp dto.GrupoEntidadeResponse
		resp.FromModel(&grupo)
		items[i] = resp
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	utils.RespondWithOK(c, dto.GrupoEntidadeListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// GetByID busca um grupo de entidade por ID
//	@Summary		Busca grupo por ID
//	@Description	Retorna os dados de um grupo de entidade específico
//	@Tags			GrupoEntidades
//	@Produce		json
//	@Param			id	path		int	true	"ID do Grupo"
//	@Success		200	{object}	dto.GrupoEntidadeResponse
//	@Failure		404	{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/grupos-entidades/{id} [get]
func (h *GrupoEntidadeHandler) GetByID(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	grupo, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondWithNotFoundError(c, err.Error())
		return
	}

	var resp dto.GrupoEntidadeResponse
	resp.FromModel(grupo)
	utils.RespondWithOK(c, resp)
}

// Update atualiza um grupo de entidade
//	@Summary		Atualiza um grupo de entidade
//	@Description	Atualiza o nome de um grupo de entidade
//	@Tags			GrupoEntidades
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"ID do Grupo"
//	@Param			request	body		dto.GrupoEntidadeRequest	true	"Dados para atualização"
//	@Success		200		{object}	dto.GrupoEntidadeResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/grupos-entidades/{id} [put]
func (h *GrupoEntidadeHandler) Update(c *gin.Context) {
	// Implementação similar ao Create e GetByID
}

// Delete exclui um grupo de entidade
//	@Summary		Exclui um grupo de entidade
//	@Description	Realiza a exclusão lógica de um grupo de entidade
//	@Tags			GrupoEntidades
//	@Param			id	path	int	true	"ID do Grupo"
//	@Success		204	"No Content"
//	@Failure		404	{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/grupos-entidades/{id} [delete]
func (h *GrupoEntidadeHandler) Delete(c *gin.Context) {
	// Implementação similar ao GetByID
}
