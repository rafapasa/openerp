package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/handler"
)

// RegisterTabelaPrecoRoutes registra as rotas para o CRUD de tabelas de preço.
func RegisterTabelaPrecoRoutes(router *gin.RouterGroup, h *handler.TabelaPrecoHandler) {
	tabelas := router.Group("/tabelas-preco")
	{
		tabelas.POST("", h.Create)
		tabelas.GET("", h.List)
		tabelas.GET("/:id", h.GetByID)
		tabelas.PUT("/:id", h.Update)
		tabelas.DELETE("/:id", h.Delete)
	}
}