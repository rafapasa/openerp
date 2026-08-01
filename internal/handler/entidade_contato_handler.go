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

// EntidadeContatoHandler é o handler para contatos de entidade
type EntidadeContatoHandler struct {
	service service.EntidadeContatoService
}

// ============================================================
// CONSTRUCTOR
// ============================================================

// NewEntidadeContatoHandler cria uma nova instância
func NewEntidadeContatoHandler(service service.EntidadeContatoService) *EntidadeContatoHandler {
	return &EntidadeContatoHandler{
		service: service,
	}
}

// ============================================================
// HELPERS PRIVADOS
// ============================================================

// getEntidadeID extrai e valida o ID da entidade da URL
func (h *EntidadeContatoHandler) getEntidadeID(c *gin.Context) (int, bool) {
	entidadeID, err := strconv.Atoi(c.Param("entidade_id"))
	if err != nil {
		utils.RespondWithValidationError(c, "ID da entidade inválido")
		return 0, false
	}
	return entidadeID, true
}

// getItem extrai e valida o item da URL
func (h *EntidadeContatoHandler) getItem(c *gin.Context) (int, bool) {
	item, err := strconv.Atoi(c.Param("item"))
	if err != nil {
		utils.RespondWithValidationError(c, "ID do contato inválido")
		return 0, false
	}
	return item, true
}

// ============================================================
// HANDLERS
// ============================================================

// Create cria um novo contato para uma entidade
//
//	@Summary		Cria um novo contato
//	@Description	Cadastra um novo contato para uma entidade
//	@Tags			EntidadeContatos
//	@Accept			json
//	@Produce		json
//	@Param			entidadeId	path		int							true	"ID da entidade"
//	@Param			request		body		dto.EntidadeContatoRequest	true	"Dados do contato"
//	@Success		201			{object}	dto.EntidadeContatoResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/entidades/{entidadeId}/contatos [post]
func (h *EntidadeContatoHandler) Create(c *gin.Context) {
	// 1. Extrair entidadeID da URL
	entidadeID, ok := h.getEntidadeID(c)
	if !ok {
		return
	}

	// 2. Bind do JSON
	var req dto.EntidadeContatoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	// 3. Associar ao ID da entidade (segurança)
	req.EntidadeID = entidadeID

	// 4. Obter ID do usuário
	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	// 5. Chamar service
	contato, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	// 6. Mapear para response
	var resp dto.EntidadeContatoResponse
	resp.FromModel(contato)

	// 7. Retornar 201 Created
	utils.RespondWithCreated(c, resp)
}

// List lista todos os contatos de uma entidade
//
//	@Summary		Lista contatos de uma entidade
//	@Description	Retorna todos os contatos de uma entidade
//	@Tags			EntidadeContatos
//	@Produce		json
//	@Param			entidadeId	path		int	true	"ID da entidade"
//	@Param			limit		query		int	false	"Limite de registros"	default(10)
//	@Param			offset		query		int	false	"Offset para paginação"	default(0)
//	@Param			tipo		query		int	false	"Filtrar por tipo de contato"
//	@Success		200			{object}	dto.EntidadeContatoListResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/entidades/{entidadeId}/contatos [get]
func (h *EntidadeContatoHandler) List(c *gin.Context) {
	// 1. Extrair entidadeID da URL
	entidadeID, ok := h.getEntidadeID(c)
	if !ok {
		return
	}

	// 2. Extrair parâmetros
	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)

	// 3. Construir filtros (incluir entidadeID obrigatório)
	filters := make(map[string]interface{})
	filters["ent_id"] = entidadeID

	if tipo := utils.GetQueryInt(c, "tipo", 0); tipo > 0 {
		filters["frc_id"] = tipo
	}

	// 4. Chamar service
	contatos, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithInternalError(c, err.Error())
		return
	}

	// 5. Mapear para response
	items := make([]dto.EntidadeContatoResponse, len(contatos))
	for i, contato := range contatos {
		var resp dto.EntidadeContatoResponse
		resp.FromModel(&contato)
		items[i] = resp
	}

	// 6. Calcular total de páginas
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	// 7. Retornar 200 OK
	utils.RespondWithOK(c, dto.EntidadeContatoListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// GetByID busca um contato específico
//
//	@Summary		Busca um contato específico
//	@Description	Retorna os dados de um contato específico
//	@Tags			EntidadeContatos
//	@Produce		json
//	@Param			entidadeId	path		int	true	"ID da entidade"
//	@Param			item		path		int	true	"ID do contato"
//	@Success		200			{object}	dto.EntidadeContatoResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/entidades/{entidadeId}/contatos/{item} [get]
func (h *EntidadeContatoHandler) GetByID(c *gin.Context) {
	// 1. Extrair IDs
	entidadeID, ok := h.getEntidadeID(c)
	if !ok {
		return
	}

	item, ok := h.getItem(c)
	if !ok {
		return
	}

	// 2. Chamar service
	contato, err := h.service.GetByID(entidadeID, item)
	if err != nil {
		utils.RespondWithNotFoundError(c, err.Error())
		return
	}

	// 3. Mapear para response
	var resp dto.EntidadeContatoResponse
	resp.FromModel(contato)

	// 4. Retornar 200 OK
	utils.RespondWithOK(c, resp)
}

// Update atualiza um contato
//
//	@Summary		Atualiza um contato
//	@Description	Atualiza os dados de um contato existente
//	@Tags			EntidadeContatos
//	@Accept			json
//	@Produce		json
//	@Param			entidadeId	path		int							true	"ID da entidade"
//	@Param			item		path		int							true	"ID do contato"
//	@Param			request		body		dto.EntidadeContatoRequest	true	"Dados atualizados"
//	@Success		200			{object}	dto.EntidadeContatoResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/entidades/{entidadeId}/contatos/{item} [put]
func (h *EntidadeContatoHandler) Update(c *gin.Context) {
	// 1. Extrair IDs
	entidadeID, ok := h.getEntidadeID(c)
	if !ok {
		return
	}

	item, ok := h.getItem(c)
	if !ok {
		return
	}

	// 2. Bind do JSON
	var req dto.EntidadeContatoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	// 3. Associar ao ID da entidade (segurança)
	req.EntidadeID = entidadeID

	// 4. Obter ID do usuário
	userID := middleware.GetUserID(c)
	req.UpdatedBy = &userID

	// 5. Chamar service
	contato, err := h.service.Update(entidadeID, item, &req)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	// 6. Mapear para response
	var resp dto.EntidadeContatoResponse
	resp.FromModel(contato)

	// 7. Retornar 200 OK
	utils.RespondWithOK(c, resp)
}

// Delete exclui um contato
//
//	@Summary		Exclui um contato
//	@Description	Realiza a exclusão lógica de um contato
//	@Tags			EntidadeContatos
//	@Produce		json
//	@Param			entidadeId	path	int	true	"ID da entidade"
//	@Param			item		path	int	true	"ID do contato"
//	@Success		204			"No Content"
//	@Failure		400			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/entidades/{entidadeId}/contatos/{item} [delete]
func (h *EntidadeContatoHandler) Delete(c *gin.Context) {
	// 1. Extrair IDs
	entidadeID, ok := h.getEntidadeID(c)
	if !ok {
		return
	}

	item, ok := h.getItem(c)
	if !ok {
		return
	}

	// 2. Chamar service
	if err := h.service.Delete(entidadeID, item); err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	// 3. Retornar 204 No Content
	utils.RespondWithNoContent(c)
}
