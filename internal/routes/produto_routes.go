package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/handler"
)

func RegisterProdutoRoutes(router *gin.RouterGroup, produtoHandler *handler.ProdutoHandler) {
	produtos := router.Group("/produtos")
	{
		produtos.POST("/", produtoHandler.Create)
		produtos.GET("/:id", produtoHandler.GetByID)
		produtos.PUT("/:id", produtoHandler.Update)
		produtos.DELETE("/:id", produtoHandler.Delete)
	}
}
