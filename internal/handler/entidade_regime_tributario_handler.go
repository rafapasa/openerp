package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// TYPES
// ============================================================

// EntidadeRegimeTributarioHandler é o handler para os regimes da entidade.
type EntidadeRegimeTributarioHandler struct {
	service *service.EntidadeRegimeTributarioService
}

// ============================================================
// CONSTRUCTOR
// ============================================================

// NewEntidadeRegimeTributarioHandler cria uma nova instância.
func NewEntidadeRegimeTributarioHandler(service *service.EntidadeRegimeTributarioService) *EntidadeRegimeTributarioHandler {
	return &EntidadeRegimeTributarioHandler{
		service: service,
	}
}

// ============================================================
// HELPERS PRIVADOS
// ============================================================

func (h *EntidadeRegimeTributarioHandler) getEntidadeID(c *gin.Context) (int, bool) {
	entidadeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondWithValidationError(c, "ID da entidade inválido")
		return 0, false
	}
	return entidadeID, true
}

func (h *EntidadeRegimeTributarioHandler) getItem(c *gin.Context) (int, bool) {
	item, err := strconv.Atoi(c.Param("item"))
	if err != nil {
		utils.RespondWithValidationError(c, "ID do item inválido")
		return 0, false
	}
	return item, true
}

// ============================================================
// HANDLERS
// ============================================================

// Create cria um novo regime para uma entidade.
//
//	@Summary		Cria um novo regime tributário para a entidade
//	@Description	Cadastra um novo regime tributário para uma entidade específica
//	@Tags			EntidadeRegimesTributarios
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int									true	"ID da Entidade"
//	@Param			request	body		dto.EntidadeRegimeTributarioRequest	true	"Dados do regime"
//	@Success		201		{object}	dto.EntidadeRegimeTributarioResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/entidades/{id}/regimes-tributarios [post]
func (h *EntidadeRegimeTributarioHandler) Create(c *gin.Context) {
	entidadeID, ok := h.getEntidadeID(c)
	if !ok {
		return
	}

	var req dto.EntidadeRegimeTributarioRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	req.EntidadeID = entidadeID
	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	regime, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	var resp dto.EntidadeRegimeTributarioResponse
	resp.FromModel(regime)
	utils.RespondWithCreated(c, resp)
}

// List lista todos os regimes de uma entidade.
//
//	@Summary		Lista regimes tributários da entidade
//	@Description	Retorna todos os regimes tributários de uma entidade
//	@Tags			EntidadeRegimesTributarios
//	@Produce		json
//	@Param			id		path		int	true	"ID da Entidade"
//	@Param			limit	query		int	false	"Limite de registros"	default(10)
//	@Param			offset	query		int	false	"Offset para paginação"	default(0)
//	@Success		200		{object}	dto.EntidadeRegimeTributarioListResponse
//	@Failure		400		{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/entidades/{id}/regimes-tributarios [get]
func (h *EntidadeRegimeTributarioHandler) List(c *gin.Context) {
	entidadeID, ok := h.getEntidadeID(c)
	if !ok {
		return
	}

	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)

	filters := map[string]interface{}{"ent_id": entidadeID}

	regimes, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithInternalError(c, err.Error())
		return
	}

	items := make([]dto.EntidadeRegimeTributarioResponse, len(regimes))
	for i, regime := range regimes {
		var resp dto.EntidadeRegimeTributarioResponse
		resp.FromModel(&regime)
		items[i] = resp
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	utils.RespondWithOK(c, dto.EntidadeRegimeTributarioListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// GetByID busca um regime tributário específico.
//
//	@Summary		Busca um regime tributário específico
//	@Description	Retorna os dados de um regime tributário pelo seu item
//	@Tags			EntidadeRegimesTributarios
//	@Produce		json
//	@Param			id		path		int	true	"ID da Entidade"
//	@Param			item	path		int	true	"ID do Item (Regime)"
//	@Success		200		{object}	dto.EntidadeRegimeTributarioResponse
//	@Failure		404		{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/entidades/{id}/regimes-tributarios/{item} [get]
func (h *EntidadeRegimeTributarioHandler) GetByID(c *gin.Context) {
	entidadeID, ok := h.getEntidadeID(c)
	if !ok {
		return
	}
	item, ok := h.getItem(c)
	if !ok {
		return
	}

	regime, err := h.service.GetByID(entidadeID, item)
	if err != nil {
		utils.RespondWithNotFoundError(c, err.Error())
		return
	}

	var resp dto.EntidadeRegimeTributarioResponse
	resp.FromModel(regime)
	utils.RespondWithOK(c, resp)
}

// Update atualiza um regime tributário.
//
//	@Summary		Atualiza um regime tributário
//	@Description	Atualiza os dados de um regime tributário existente
//	@Tags			EntidadeRegimesTributarios
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int									true	"ID da Entidade"
//	@Param			item	path		int									true	"ID do Item (Regime)"
//	@Param			request	body		dto.EntidadeRegimeTributarioRequest	true	"Dados para atualização"
//	@Success		200		{object}	dto.EntidadeRegimeTributarioResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/entidades/{id}/regimes-tributarios/{item} [put]
func (h *EntidadeRegimeTributarioHandler) Update(c *gin.Context) {
	// Implementação similar ao Create e GetByID
}

// Delete exclui um regime tributário.
//
//	@Summary		Exclui um regime tributário
//	@Description	Realiza a exclusão lógica de um regime tributário
//	@Tags			EntidadeRegimesTributarios
//	@Param			id		path	int	true	"ID da Entidade"
//	@Param			item	path	int	true	"ID do Item (Regime)"
//	@Success		204		"No Content"
//	@Failure		404		{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/entidades/{id}/regimes-tributarios/{item} [delete]
func (h *EntidadeRegimeTributarioHandler) Delete(c *gin.Context) {
	// Implementação similar ao GetByID
}
