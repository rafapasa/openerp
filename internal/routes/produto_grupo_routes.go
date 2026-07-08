package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/handler"
)

// RegisterProdutoGrupoRoutes registra as rotas para o CRUD de grupo de produtos.
func RegisterProdutoGrupoRoutes(router *gin.RouterGroup, h *handler.ProdutoGrupoHandler) {
	grupos := router.Group("/produto-grupos")
	{
		grupos.POST("", h.Create)
		grupos.GET("", h.List)
		grupos.GET("/:id", h.GetByID)
		grupos.PUT("/:id", h.Update)
		grupos.DELETE("/:id", h.Delete)
	}
}
