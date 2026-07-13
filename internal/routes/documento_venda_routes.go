package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/handler"
)

func RegisterDocumentoVendaRoutes(router *gin.RouterGroup, h *handler.DocumentoVendaHandler) {
	docs := router.Group("/documentos-venda")
	{
		docs.POST("", h.Create)
		docs.GET("", h.List)
		docs.GET("/:id", h.GetByID)
		docs.PUT("/:id", h.Update)
		docs.DELETE("/:id", h.Delete)
	}
}
