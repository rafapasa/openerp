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

// EntidadeEnderecoHandler é o handler para endereços de entidade
type EntidadeEnderecoHandler struct {
	service *service.EntidadeEnderecoService
}

// ============================================================
// CONSTRUCTOR
// ============================================================

// NewEntidadeEnderecoHandler cria uma nova instância
func NewEntidadeEnderecoHandler(service *service.EntidadeEnderecoService) *EntidadeEnderecoHandler {
	return &EntidadeEnderecoHandler{
		service: service,
	}
}

// ============================================================
// HELPERS PRIVADOS
// ============================================================

// getEntidadeID extrai e valida o ID da entidade da URL
func (h *EntidadeEnderecoHandler) getEntidadeID(c *gin.Context) (int, bool) {
	entidadeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondWithValidationError(c, "ID da entidade inválido")
		return 0, false
	}
	return entidadeID, true
}

// getItem extrai e valida o item da URL
func (h *EntidadeEnderecoHandler) getItem(c *gin.Context) (int, bool) {
	item, err := strconv.Atoi(c.Param("item"))
	if err != nil {
		utils.RespondWithValidationError(c, "ID do endereço inválido")
		return 0, false
	}
	return item, true
}

// ============================================================
// HANDLERS
// ============================================================

// Create cria um novo endereço para uma entidade
// @Summary Cria um novo endereço
// @Description Cadastra um novo endereço para uma entidade
// @Tags EntidadeEnderecos
// @Accept json
// @Produce json
// @Param entidadeId path int true "ID da entidade"
// @Param request body dto.EntidadeEnderecoRequest true "Dados do endereço"
// @Success 201 {object} dto.EntidadeEnderecoResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /entidades/{entidadeId}/enderecos [post]
func (h *EntidadeEnderecoHandler) Create(c *gin.Context) {
	// 1. Extrair entidadeID da URL
	entidadeID, ok := h.getEntidadeID(c)
	if !ok {
		return
	}

	// 2. Bind do JSON
	var req dto.EntidadeEnderecoRequest
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
	endereco, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	// 6. Mapear para response
	var resp dto.EntidadeEnderecoResponse
	resp.FromModel(endereco)

	// 7. Retornar 201 Created
	utils.RespondWithCreated(c, resp)
}

// List lista todos os endereços de uma entidade
// @Summary Lista endereços de uma entidade
// @Description Retorna todos os endereços de uma entidade
// @Tags EntidadeEnderecos
// @Produce json
// @Param entidadeId path int true "ID da entidade"
// @Param limit query int false "Limite de registros" default(10)
// @Param offset query int false "Offset para paginação" default(0)
// @Param tipo query int false "Filtrar por tipo"
// @Success 200 {object} dto.EntidadeEnderecoListResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /entidades/{entidadeId}/enderecos [get]
func (h *EntidadeEnderecoHandler) List(c *gin.Context) {
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
		filters["tipo"] = tipo
	}

	// 4. Chamar service
	enderecos, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithInternalError(c, err.Error())
		return
	}

	// 5. Mapear para response
	items := make([]dto.EntidadeEnderecoResponse, len(enderecos))
	for i, endereco := range enderecos {
		var resp dto.EntidadeEnderecoResponse
		resp.FromModel(&endereco)
		items[i] = resp
	}

	// 6. Calcular total de páginas
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	// 7. Retornar 200 OK
	utils.RespondWithOK(c, dto.EntidadeEnderecoListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// GetByID busca um endereço específico
// @Summary Busca um endereço específico
// @Description Retorna os dados de um endereço específico
// @Tags EntidadeEnderecos
// @Produce json
// @Param entidadeId path int true "ID da entidade"
// @Param item path int true "ID do endereço"
// @Success 200 {object} dto.EntidadeEnderecoResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /entidades/{entidadeId}/enderecos/{item} [get]
func (h *EntidadeEnderecoHandler) GetByID(c *gin.Context) {
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
	endereco, err := h.service.GetByID(entidadeID, item)
	if err != nil {
		utils.RespondWithNotFoundError(c, err.Error())
		return
	}

	// 3. Mapear para response
	var resp dto.EntidadeEnderecoResponse
	resp.FromModel(endereco)

	// 4. Retornar 200 OK
	utils.RespondWithOK(c, resp)
}

// Update atualiza um endereço
// @Summary Atualiza um endereço
// @Description Atualiza os dados de um endereço existente
// @Tags EntidadeEnderecos
// @Accept json
// @Produce json
// @Param entidadeId path int true "ID da entidade"
// @Param item path int true "ID do endereço"
// @Param request body dto.EntidadeEnderecoRequest true "Dados atualizados"
// @Success 200 {object} dto.EntidadeEnderecoResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /entidades/{entidadeId}/enderecos/{item} [put]
func (h *EntidadeEnderecoHandler) Update(c *gin.Context) {
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
	var req dto.EntidadeEnderecoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	// 3. Associar ao ID da entidade (segurança)
	req.EntidadeID = entidadeID

	// 4. Obter ID do usuário
	userID := middleware.GetUserID(c)
	req.UpdatedBy = &userID

	// 5. Chamar service
	endereco, err := h.service.Update(entidadeID, item, &req)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	// 6. Mapear para response
	var resp dto.EntidadeEnderecoResponse
	resp.FromModel(endereco)

	// 7. Retornar 200 OK
	utils.RespondWithOK(c, resp)
}

// Delete exclui um endereço
// @Summary Exclui um endereço
// @Description Realiza a exclusão lógica de um endereço
// @Tags EntidadeEnderecos
// @Produce json
// @Param entidadeId path int true "ID da entidade"
// @Param item path int true "ID do endereço"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /entidades/{entidadeId}/enderecos/{item} [delete]
func (h *EntidadeEnderecoHandler) Delete(c *gin.Context) {
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
