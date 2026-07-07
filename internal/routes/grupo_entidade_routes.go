package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/handler"
)

// RegisterGrupoEntidadeRoutes registra as rotas para o CRUD de Grupo de Entidades
func RegisterGrupoEntidadeRoutes(router *gin.RouterGroup, h *handler.GrupoEntidadeHandler) {
	grupos := router.Group("/grupos-entidades")
	{
		grupos.POST("/", h.Create)
		grupos.GET("/", h.List)
		grupos.GET("/:id", h.GetByID)
		grupos.PUT("/:id", h.Update)
		grupos.DELETE("/:id", h.Delete)
	}
}
