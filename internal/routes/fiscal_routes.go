// internal/routes/fiscal_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/appwire"
)

func RegisterFiscalRoutes(router *gin.RouterGroup, container *appwire.Container) {
	// Operações Fiscais
	opf := router.Group("/operacoes-fiscais")
	{
		opf.GET("", container.OperacaoFiscalHandler.List)
		opf.POST("", container.OperacaoFiscalHandler.Create) // Corrected: This was missing
		opf.GET("/:id", container.OperacaoFiscalHandler.GetByID)
		opf.PUT("/:id", container.OperacaoFiscalHandler.Update)
		opf.DELETE("/:id", container.OperacaoFiscalHandler.Delete)
		opf.GET("/cfop/:cfop", container.OperacaoFiscalHandler.GetByCFOP)
		opf.GET("/filial/:filial_id", container.OperacaoFiscalHandler.GetByFilial)
	}

	// Processos
	processo := router.Group("/processos")
	{
		processo.GET("", container.ProcessoHandler.List)
		processo.POST("", container.ProcessoHandler.Create) // Corrected: This was missing
		processo.GET("/:id", container.ProcessoHandler.GetByID)
		processo.PUT("/:id", container.ProcessoHandler.Update)
		processo.DELETE("/:id", container.ProcessoHandler.Delete)
		processo.GET("/codigo/:codigo", container.ProcessoHandler.GetByCodigo)
	}
}
