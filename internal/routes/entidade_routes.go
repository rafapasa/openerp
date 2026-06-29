package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/openerp/backend/internal/handler"
)

// RegisterEntidadeRoutes registra as rotas de entidade
func RegisterEntidadeRoutes(router *gin.RouterGroup, handler *handler.EntidadeHandler) {
	entidades := router.Group("/entidades")
	{
		entidades.POST("/", handler.Create)                            // POST /entidades
		entidades.GET("/", handler.List)                               // GET /entidades
		entidades.GET("/:id", handler.GetByID)                         // GET /entidades/:id
		entidades.PUT("/:id", handler.Update)                          // PUT /entidades/:id
		entidades.DELETE("/:id", handler.Delete)                       // DELETE /entidades/:id
		entidades.GET("/documento/:documento", handler.GetByDocumento) // GET /entidades/documento/:documento
	}
}
