package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

// EntidadeDocumentoHandler é o handler para documentos de entidade.
type EntidadeDocumentoHandler struct {
	service service.EntidadeDocumentoService
}

// NewEntidadeDocumentoHandler cria uma nova instância de EntidadeDocumentoHandler.
func NewEntidadeDocumentoHandler(s service.EntidadeDocumentoService) *EntidadeDocumentoHandler {
	return &EntidadeDocumentoHandler{
		service: s,
	}
}

// getEntidadeID extrai e valida o ID da entidade da URL.
func (h *EntidadeDocumentoHandler) getEntidadeID(c *gin.Context) (int, bool) {
	entidadeID, err := strconv.Atoi(c.Param("entidade_id"))
	if err != nil {
		utils.RespondWithValidationError(c, "ID da entidade inválido")
		return 0, false
	}
	return entidadeID, true
}

// getItem extrai e valida o item da URL.
func (h *EntidadeDocumentoHandler) getItem(c *gin.Context) (int, bool) {
	item, err := strconv.Atoi(c.Param("item"))
	if err != nil {
		utils.RespondWithValidationError(c, "ID do documento inválido")
		return 0, false
	}
	return item, true
}

// @Summary      Cria um novo documento
// @Description  Cadastra um novo documento para uma entidade
// @Tags         EntidadeDocumentos
// @Accept       json
// @Produce      json
// @Param        entidadeId  path      int                           true  "ID da entidade"
// @Param        request     body      dto.EntidadeDocumentoRequest  true  "Dados do documento"
// @Success      201         {object}  dto.EntidadeDocumentoResponse
// @Failure      400         {object}  map[string]string
// @Failure      404         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Security     BearerAuth
// @Router       /entidades/{entidadeId}/documentos [post]
func (h *EntidadeDocumentoHandler) Create(c *gin.Context) {
	entidadeID, ok := h.getEntidadeID(c)
	if !ok {
		return
	}

	var req dto.EntidadeDocumentoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	req.EntidadeID = entidadeID
	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID

	documento, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.EntidadeDocumentoResponse
	resp.FromModel(documento)
	utils.RespondWithCreated(c, resp)
}

// @Summary      Busca um documento específico
// @Description  Retorna os dados de um documento específico
// @Tags         EntidadeDocumentos
// @Produce      json
// @Param        entidadeId  path      int  true  "ID da entidade"
// @Param        item        path      int  true  "ID do documento"
// @Success      200         {object}  dto.EntidadeDocumentoResponse
// @Failure      400         {object}  map[string]string
// @Failure      404         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Security     BearerAuth
// @Router       /entidades/{entidadeId}/documentos/{item} [get]
func (h *EntidadeDocumentoHandler) GetByID(c *gin.Context) {
	entidadeID, ok := h.getEntidadeID(c)
	if !ok {
		return
	}
	item, ok := h.getItem(c)
	if !ok {
		return
	}

	documento, err := h.service.GetByID(entidadeID, item)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.EntidadeDocumentoResponse
	resp.FromModel(documento)
	utils.RespondWithOK(c, resp)
}

// @Summary      Atualiza um documento
// @Description  Atualiza os dados de um documento existente
// @Tags         EntidadeDocumentos
// @Accept       json
// @Produce      json
// @Param        entidadeId  path      int                           true  "ID da entidade"
// @Param        item        path      int                           true  "ID do documento"
// @Param        request     body      dto.EntidadeDocumentoRequest  true  "Dados atualizados"
// @Success      200         {object}  dto.EntidadeDocumentoResponse
// @Failure      400         {object}  map[string]string
// @Failure      404         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Security     BearerAuth
// @Router       /entidades/{entidadeId}/documentos/{item} [put]
func (h *EntidadeDocumentoHandler) Update(c *gin.Context) {
	entidadeID, ok := h.getEntidadeID(c)
	if !ok {
		return
	}
	item, ok := h.getItem(c)
	if !ok {
		return
	}

	var req dto.EntidadeDocumentoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	req.EntidadeID = entidadeID
	userID := middleware.GetUserID(c)
	req.UpdatedBy = &userID

	documento, err := h.service.Update(entidadeID, item, &req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	var resp dto.EntidadeDocumentoResponse
	resp.FromModel(documento)
	utils.RespondWithOK(c, resp)
}

// @Summary      Exclui um documento
// @Description  Realiza a exclusão lógica de um documento
// @Tags         EntidadeDocumentos
// @Produce      json
// @Param        entidadeId  path      int  true  "ID da entidade"
// @Param        item        path      int  true  "ID do documento"
// @Success      204         "No Content"
// @Failure      400         {object}  map[string]string
// @Failure      404         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Security     BearerAuth
// @Router       /entidades/{entidadeId}/documentos/{item} [delete]
func (h *EntidadeDocumentoHandler) Delete(c *gin.Context) {
	entidadeID, ok := h.getEntidadeID(c)
	if !ok {
		return
	}
	item, ok := h.getItem(c)
	if !ok {
		return
	}

	if err := h.service.Delete(entidadeID, item); err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary      Lista documentos de uma entidade
// @Description  Retorna todos os documentos de uma entidade
// @Tags         EntidadeDocumentos
// @Produce      json
// @Param        entidadeId  path      int     true  "ID da entidade"
// @Param        limit       query     int     false  "Limite de registros"  default(10)
// @Param        offset      query     int     false  "Offset para paginação"  default(0)
// @Param        tipo        query     string  false  "Filtrar por tipo de documento"
// @Success      200         {object}  dto.EntidadeDocumentoListResponse
// @Failure      400         {object}  map[string]string
// @Failure      404         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Security     BearerAuth
// @Router       /entidades/{entidadeId}/documentos [get]
func (h *EntidadeDocumentoHandler) List(c *gin.Context) {
	entidadeID, ok := h.getEntidadeID(c)
	if !ok {
		return
	}

	limit := utils.GetQueryInt(c, "limit", 10)
	offset := utils.GetQueryInt(c, "offset", 0)

	filters := make(map[string]interface{})
	filters["ent_id"] = entidadeID
	if tipo := utils.GetQueryString(c, "tipo", ""); tipo != "" {
		filters["edoc_tipo"] = tipo
	}

	documentos, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	items := make([]dto.EntidadeDocumentoResponse, len(documentos))
	for i, documento := range documentos {
		var resp dto.EntidadeDocumentoResponse
		resp.FromModel(&documento)
		items[i] = resp
	}

	totalPages := utils.CalculateTotalPages(int(total), limit)

	utils.RespondWithOK(c, dto.EntidadeDocumentoListResponse{
		Items:      items,
		Total:      total,
		Page:       offset/limit + 1,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// @Summary      Download do arquivo
// @Description  Download do arquivo de um documento
// @Tags         EntidadeDocumentos
// @Produce      application/octet-stream
// @Param        entidadeId  path      int  true  "ID da entidade"
// @Param        item        path      int  true  "ID do documento"
// @Success      200         {file}    string "Arquivo do documento"
// @Failure      400         {object}  map[string]string
// @Failure      404         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Security     BearerAuth
// @Router       /entidades/{entidadeId}/documentos/{item}/download [get]
func (h *EntidadeDocumentoHandler) Download(c *gin.Context) {
	entidadeID, ok := h.getEntidadeID(c)
	if !ok {
		return
	}
	item, ok := h.getItem(c)
	if !ok {
		return
	}

	documento, err := h.service.GetByID(entidadeID, item)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	if len(documento.Arquivo) == 0 {
		utils.RespondWithErrorAny(c, apperrors.NewNotFoundError("Arquivo não encontrado para download."))
		return
	}

	// Define o nome do arquivo para download
	filename := fmt.Sprintf("documento_entidade_%d_item_%d", entidadeID, item)
	if documento.Descricao != nil && *documento.Descricao != "" {
		filename = *documento.Descricao
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	// c.Header("Content-Type", dto.GetContentType(documento.Arquivo))
	// c.Data(http.StatusOK, dto.GetContentType(documento.Arquivo), documento.Arquivo)
}
