package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/apperrors"
)

// ============================================================
// TYPES
// ============================================================

// ErrorResponse representa uma resposta de erro padronizada
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// ============================================================
// FUNÇÕES DE RESPOSTA
// ============================================================

// RespondWithError envia uma resposta de erro padronizada
func RespondWithError(c *gin.Context, status int, errorCode string, message string) {
	c.JSON(status, ErrorResponse{
		Error:   errorCode,
		Message: message,
	})
}

// RespondWithValidationError envia erro de validação (400)
func RespondWithValidationError(c *gin.Context, message string) {
	RespondWithError(c, http.StatusBadRequest, "validation_error", message)
}

// RespondWithNotFoundError envia erro de não encontrado (404)
func RespondWithNotFoundError(c *gin.Context, message string) {
	RespondWithError(c, http.StatusNotFound, "not_found", message)
}

// RespondWithInternalError envia erro interno (500)
func RespondWithInternalError(c *gin.Context, message string) {
	RespondWithError(c, http.StatusInternalServerError, "internal_error", message)
}

// RespondWithUnauthorizedError envia erro de não autorizado (401)
func RespondWithUnauthorizedError(c *gin.Context, message string) {
	RespondWithError(c, http.StatusUnauthorized, "unauthorized", message)
}

// RespondWithForbiddenError envia erro de proibido (403)
func RespondWithForbiddenError(c *gin.Context, message string) {
	RespondWithError(c, http.StatusForbidden, "forbidden", message)
}

// RespondWithConflictError envia erro de conflito (409)
func RespondWithConflictError(c *gin.Context, message string) {
	RespondWithError(c, http.StatusConflict, "conflict", message)
}

// RespondWithSuccess envia uma resposta de sucesso padronizada
func RespondWithSuccess(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

// RespondWithCreated envia resposta de criação (201)
func RespondWithCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

// RespondWithOK envia resposta de sucesso (200)
func RespondWithOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// RespondWithNoContent envia resposta sem conteúdo (204)
func RespondWithNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// ============================================================
// FUNÇÕES DE BIND COM VALIDAÇÃO
// ============================================================

// BindAndValidate faz o bind e validação do JSON, retornando erro se falhar
func BindAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return err
	}
	return nil
}

// BindAndValidateOrRespond faz o bind e validação, respondendo com erro se falhar
func BindAndValidateOrRespond(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		RespondWithValidationError(c, "Dados inválidos: "+err.Error())
		return false
	}
	return true
}

// ============================================================
// FUNÇÕES DE EXTRAÇÃO DE PARÂMETROS
// ============================================================

// ParseIDParam extrai e valida um parâmetro ID da URL
func ParseIDParam(c *gin.Context, param string) (int, bool) {
	id, err := parseIntParam(c.Param(param))
	if err != nil {
		RespondWithValidationError(c, "ID deve ser um número válido")
		return 0, false
	}
	return id, true
}

// parseIntParam converte string para int
func parseIntParam(value string) (int, error) {
	if value == "" {
		return 0, nil
	}
	return strconv.Atoi(value)
}

// ============================================================
// FUNÇÕES DE FILTROS
// ============================================================

// GetQueryInt obtém um parâmetro query como int
func GetQueryInt(c *gin.Context, key string, defaultValue int) int {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// GetQueryString obtém um parâmetro query como string
func GetQueryString(c *gin.Context, key string, defaultValue string) string {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func CalculateTotalPages(total, limit int) int {
	return (total + limit - 1) / limit
}

// RespondWithAppError envia uma resposta baseada no AppError
func RespondWithAppError(c *gin.Context, err *apperrors.AppError) {
	switch err.Code {
	case http.StatusBadRequest:
		RespondWithValidationError(c, err.Message)
	case http.StatusConflict:
		RespondWithConflictError(c, err.Message)
	case http.StatusNotFound:
		RespondWithNotFoundError(c, err.Message)
	case http.StatusInternalServerError:
		RespondWithInternalError(c, err.Message)
	case http.StatusUnauthorized:
		RespondWithUnauthorizedError(c, err.Message)
	case http.StatusForbidden:
		RespondWithForbiddenError(c, err.Message)
	default:
		// Caso não mapeado, responde com o código genérico
		RespondWithError(c, err.Code, "error", err.Message)
	}
}

// RespondWithErrorAny responde com qualquer erro, extraindo o código se for AppError
func RespondWithErrorAny(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		RespondWithAppError(c, appErr)
		return
	}
	// Fallback para erro interno
	RespondWithInternalError(c, err.Error())
}
