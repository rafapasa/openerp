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

// EntidadeDocumentoHandler é o handler para documentos de entidade
type EntidadeDocumentoHandler struct {
	service *service.EntidadeDocumentoService
}

// ============================================================
// CONSTRUCTOR
// ============================================================

// NewEntidadeDocumentoHandler cria uma nova instância
func NewEntidadeDocumentoHandler(service *service.EntidadeDocumentoService) *EntidadeDocumentoHandler {
	return &EntidadeDocumentoHandler{
		service: service,
	}
}

// ============================================================
// HELPERS PRIVADOS
// ============================================================

// getEntidadeID extrai e valida o ID da entidade da URL
func (h *EntidadeDocumentoHandler) getEntidadeID(c *gin.Context) (int, bool) {
	entidadeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondWithValidationError(c, "ID da entidade inválido")
		return 0, false
	}
	return entidadeID, true
}

// getItem extrai e valida o item da URL
func (h *EntidadeDocumentoHandler) getItem(c *gin.Context) (int, bool) {
	item, err := strconv.Atoi(c.Param("item"))
	if err != nil {
		utils.RespondWithValidationError(c, "ID do documento inválido")
		return 0, false
	}
	return item, true
}

// ============================================================
// HANDLERS
// ============================================================

// Create cria um novo documento para uma entidade
//
//	@Summary		Cria um novo documento
//	@Description	Cadastra um novo documento para uma entidade
//	@Tags			EntidadeDocumentos
//	@Accept			json
//	@Produce		json
//	@Param			entidadeId	path		int								true	"ID da entidade"
//	@Param			request		body		dto.EntidadeDocumentoRequest	true	"Dados do documento"
//	@Success		201			{object}	dto.EntidadeDocumentoResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/entidades/{entidadeId}/documentos [post]
func (h *EntidadeDocumentoHandler) Create(c *gin.Context) {
	// 1. Extrair entidadeID da URL
	entidadeID, ok := h.getEntidadeID(c)
	if !ok {
		return
	}

	// 2. Bind do JSON
	var req dto.EntidadeDocumentoRequest
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
	documento, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	// 6. Mapear para response
	var resp dto.EntidadeDocumentoResponse
	resp.FromModel(documento)

	// 7. Retornar 201 Created
	utils.RespondWithCreated(c, resp)
}

// List lista todos os documentos de uma entidade
//
//	@Summary		Lista documentos de uma entidade
//	@Description	Retorna todos os documentos de uma entidade
//	@Tags			EntidadeDocumentos
//	@Produce		json
//	@Param			entidadeId	path		int		true	"ID da entidade"
//	@Param			limit		query		int		false	"Limite de registros"	default(10)
//	@Param			offset		query		int		false	"Offset para paginação"	default(0)
//	@Param			tipo		query		string	false	"Filtrar por tipo de documento"
//	@Success		200			{object}	dto.EntidadeDocumentoListResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/entidades/{entidadeId}/documentos [get]
func (h *EntidadeDocumentoHandler) List(c *gin.Context) {
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
	filters["entidade_id"] = entidadeID

	if tipo := utils.GetQueryString(c, "tipo", ""); tipo != "" {
		filters["tipo"] = tipo
	}

	// 4. Chamar service
	documentos, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithInternalError(c, err.Error())
		return
	}

	// 5. Mapear para response
	items := make([]dto.EntidadeDocumentoResponse, len(documentos))
	for i, documento := range documentos {
		var resp dto.EntidadeDocumentoResponse
		resp.FromModel(&documento)
		items[i] = resp
	}

	// 6. Calcular total de páginas
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	// 7. Retornar 200 OK
	utils.RespondWithOK(c, dto.EntidadeDocumentoListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// GetByID busca um documento específico
//
//	@Summary		Busca um documento específico
//	@Description	Retorna os dados de um documento específico
//	@Tags			EntidadeDocumentos
//	@Produce		json
//	@Param			entidadeId	path		int	true	"ID da entidade"
//	@Param			item		path		int	true	"ID do documento"
//	@Success		200			{object}	dto.EntidadeDocumentoResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/entidades/{entidadeId}/documentos/{item} [get]
func (h *EntidadeDocumentoHandler) GetByID(c *gin.Context) {
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
	documento, err := h.service.GetByID(entidadeID, item)
	if err != nil {
		utils.RespondWithNotFoundError(c, err.Error())
		return
	}

	// 3. Mapear para response (com arquivo)
	var resp dto.EntidadeDocumentoResponse
	resp.FromModelWithFile(documento)

	// 4. Retornar 200 OK
	utils.RespondWithOK(c, resp)
}

// Update atualiza um documento
//
//	@Summary		Atualiza um documento
//	@Description	Atualiza os dados de um documento existente
//	@Tags			EntidadeDocumentos
//	@Accept			json
//	@Produce		json
//	@Param			entidadeId	path		int								true	"ID da entidade"
//	@Param			item		path		int								true	"ID do documento"
//	@Param			request		body		dto.EntidadeDocumentoRequest	true	"Dados atualizados"
//	@Success		200			{object}	dto.EntidadeDocumentoResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/entidades/{entidadeId}/documentos/{item} [put]
func (h *EntidadeDocumentoHandler) Update(c *gin.Context) {
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
	var req dto.EntidadeDocumentoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	// 3. Associar ao ID da entidade (segurança)
	req.EntidadeID = entidadeID

	// 4. Obter ID do usuário
	userID := middleware.GetUserID(c)
	req.UpdatedBy = &userID

	// 5. Chamar service
	documento, err := h.service.Update(entidadeID, item, &req)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}

	// 6. Mapear para response (sem arquivo, para não sobrecarregar)
	var resp dto.EntidadeDocumentoResponse
	resp.FromModel(documento)

	// 7. Retornar 200 OK
	utils.RespondWithOK(c, resp)
}

// Delete exclui um documento
//
//	@Summary		Exclui um documento
//	@Description	Realiza a exclusão lógica de um documento
//	@Tags			EntidadeDocumentos
//	@Produce		json
//	@Param			entidadeId	path	int	true	"ID da entidade"
//	@Param			item		path	int	true	"ID do documento"
//	@Success		204			"No Content"
//	@Failure		400			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/entidades/{entidadeId}/documentos/{item} [delete]
func (h *EntidadeDocumentoHandler) Delete(c *gin.Context) {
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

// Download arquivo
//
//	@Summary		Download do arquivo
//	@Description	Download do arquivo de um documento
//	@Tags			EntidadeDocumentos
//	@Produce		application/octet-stream
//	@Param			entidadeId	path		int	true	"ID da entidade"
//	@Param			item		path		int	true	"ID do documento"
//	@Success		200			{file}		file
//	@Failure		400			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/entidades/{entidadeId}/documentos/{item}/download [get]
func (h *EntidadeDocumentoHandler) Download(c *gin.Context) {
	// 1. Extrair IDs
	entidadeID, ok := h.getEntidadeID(c)
	if !ok {
		return
	}

	item, ok := h.getItem(c)
	if !ok {
		return
	}

	// 2. Buscar documento
	documento, err := h.service.GetByID(entidadeID, item)
	if err != nil {
		utils.RespondWithNotFoundError(c, err.Error())
		return
	}

	// 3. Verificar se tem arquivo
	if len(documento.Arquivo) == 0 {
		utils.RespondWithNotFoundError(c, "arquivo não encontrado")
		return
	}

	// 4. Determinar nome do arquivo
	filename := "documento_" + strconv.Itoa(entidadeID) + "_" + strconv.Itoa(item)
	if documento.Tipo != nil && *documento.Tipo != "" {
		filename = *documento.Tipo + "_" + strconv.Itoa(entidadeID) + "_" + strconv.Itoa(item)
	}

	// 5. Determinar Content-Type
	contentType := getContentTypeFromBytes(documento.Arquivo)

	// 6. Retornar arquivo
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(200, contentType, documento.Arquivo)
}

// getContentTypeFromBytes detecta o Content-Type
func getContentTypeFromBytes(data []byte) string {
	if len(data) == 0 {
		return "application/octet-stream"
	}

	// PDF
	if len(data) >= 4 && data[0] == 0x25 && data[1] == 0x50 && data[2] == 0x44 && data[3] == 0x46 {
		return "application/pdf"
	}

	// PNG
	if len(data) >= 8 && data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 &&
		data[4] == 0x0D && data[5] == 0x0A && data[6] == 0x1A && data[7] == 0x0A {
		return "image/png"
	}

	// JPEG
	if len(data) >= 2 && data[0] == 0xFF && data[1] == 0xD8 {
		return "image/jpeg"
	}

	// GIF
	if len(data) >= 6 && data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x38 &&
		(data[4] == 0x37 || data[4] == 0x39) && data[5] == 0x61 {
		return "image/gif"
	}

	// DOCX (ZIP)
	if len(data) >= 4 && data[0] == 0x50 && data[1] == 0x4B && data[2] == 0x03 && data[3] == 0x04 {
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	}

	return "application/octet-stream"
}
