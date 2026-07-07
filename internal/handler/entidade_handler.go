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

// EntidadeHandler é o handler para entidades
type EntidadeHandler struct {
	service *service.EntidadeService
}

// ============================================================
// CONSTRUCTOR
// ============================================================

// NewEntidadeHandler cria uma nova instância do EntidadeHandler
func NewEntidadeHandler(service *service.EntidadeService) *EntidadeHandler {
	return &EntidadeHandler{
		service: service,
	}
}

// ============================================================
// HANDLERS
// ============================================================

// Create cria uma nova entidade
//
//	@Summary		Cria uma nova entidade
//	@Description	Cadastra uma nova entidade (cliente/fornecedor)
//	@Tags			Entidades
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.EntidadeRequest	true	"Dados da entidade"
//	@Success		201		{object}	dto.EntidadeResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/entidades [post]
func (h *EntidadeHandler) Create(c *gin.Context) {
	// TODO: Implementar
	var req dto.EntidadeRequest

	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	// 2. Obter ID do usuário logado (do middleware)
	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	entidade, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	var resp dto.EntidadeResponse
	resp.FromModel(entidade)
	utils.RespondWithCreated(c, resp)
}

// GetByID busca uma entidade por ID
//
//	@Summary		Busca entidade por ID
//	@Description	Retorna os dados de uma entidade específica
//	@Tags			Entidades
//	@Produce		json
//	@Param			id	path		int	true	"ID da entidade"
//	@Success		200	{object}	dto.EntidadeResponse
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/entidades/{id} [get]
func (h *EntidadeHandler) GetByID(c *gin.Context) {
	// TODO: Implementar
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}

	entidade, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondWithNotFoundError(c, err.Error())
		return
	}

	var resp dto.EntidadeResponse
	resp.FromModel(entidade)
	utils.RespondWithOK(c, resp)
}

// GetByDocumento busca uma entidade por CPF/CNPJ
//
//	@Summary		Busca entidade por documento
//	@Description	Retorna os dados de uma entidade pelo CPF/CNPJ
//	@Tags			Entidades
//	@Produce		json
//	@Param			documento	path		string	true	"CPF/CNPJ da entidade"
//	@Success		200			{object}	dto.EntidadeResponse
//	@Failure		404			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/entidades/documento/{documento} [get]
func (h *EntidadeHandler) GetByDocumento(c *gin.Context) {
	// TODO: Implementar
	documento := c.Param("documento")
	if documento == "" {
		utils.RespondWithValidationError(c, "documento é obrigatorio")
		return
	}
	entidade, err := h.service.GetByDocumento(documento)
	if err != nil {
		utils.RespondWithNotFoundError(c, err.Error())
		return
	}

	var resp dto.EntidadeResponse
	resp.FromModel(entidade)
	utils.RespondWithOK(c, resp)
}

// Update atualiza uma entidade
//
//	@Summary		Atualiza uma entidade
//	@Description	Atualiza os dados de uma entidade existente
//	@Tags			Entidades
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"ID da entidade"
//	@Param			request	body		dto.EntidadeRequest	true	"Dados atualizados"
//	@Success		200		{object}	dto.EntidadeResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/entidades/{id} [put]
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

	entidade, err := h.service.Update(id, &req)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	var resp dto.EntidadeResponse
	resp.FromModel(entidade)
	utils.RespondWithOK(c, resp)
}

// Delete exclui uma entidade
//
//	@Summary		Exclui uma entidade
//	@Description	Realiza a exclusão lógica de uma entidade
//	@Tags			Entidades
//	@Produce		json
//	@Param			id	path	int	true	"ID da entidade"
//	@Success		204	"No Content"
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/entidades/{id} [delete]
func (h *EntidadeHandler) Delete(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}
	if err := h.service.Delete(id); err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}
	utils.RespondWithNoContent(c)
}

// List lista entidades com paginação e filtros
//
//	@Summary		Lista entidades
//	@Description	Retorna uma lista paginada de entidades
//	@Tags			Entidades
//	@Produce		json
//	@Param			limit		query		int		false	"Limite de registros"	default(10)
//	@Param			offset		query		int		false	"Offset para paginação"	default(0)
//	@Param			nome		query		string	false	"Filtrar por nome (razao social)"
//	@Param			documento	query		string	false	"Filtrar por documento"
//	@Param			tipo_pessoa	query		int		false	"Filtrar por tipo (1-Física, 2-Jurídica)"
//	@Param			situacao	query		int		false	"Filtrar por situação (1-Ativo, 2-Inativo)"
//	@Success		200			{object}	dto.EntidadeListResponse
//	@Failure		500			{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/entidades [get]
func (h *EntidadeHandler) List(c *gin.Context) {
	// Extrair parâmetros de paginação
	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)

	// Construir filtros
	filters := make(map[string]interface{})

	if nome := utils.GetQueryString(c, "nome", ""); nome != "" {
		filters["nome"] = nome
	}
	if documento := utils.GetQueryString(c, "documento", ""); documento != "" {
		filters["documento"] = documento
	}
	if tipoPessoa := utils.GetQueryString(c, "tipo_pessoa", ""); tipoPessoa != "" {
		if val, err := strconv.Atoi(tipoPessoa); err == nil {
			filters["tipo_pessoa"] = val
		}
	}
	if situacao := utils.GetQueryString(c, "situacao", ""); situacao != "" {
		if val, err := strconv.Atoi(situacao); err == nil {
			filters["situacao"] = val
		}
	}

	// Chamar service
	entidades, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithInternalError(c, err.Error())
		return
	}

	// Mapear para response
	items := make([]dto.EntidadeResponse, len(entidades))
	for i, entidade := range entidades {
		var resp dto.EntidadeResponse
		resp.FromModel(&entidade)
		items[i] = resp
	}

	// Calcular total de páginas
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	// Retornar 200 OK
	utils.RespondWithOK(c, dto.EntidadeListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}
