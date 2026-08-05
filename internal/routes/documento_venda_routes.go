// internal/routes/venda_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/appwire"
)

func RegisterVendaRoutes(router *gin.RouterGroup, container *appwire.Container) {
	// Documentos de Venda
	documento := router.Group("/documentos-venda")

	// Itens do Documento
	item := router.Group("/documentos-venda/:documento_id/itens")
	{
		item.GET("", container.DocumentoVendaHandler.ListItem)
		item.POST("", container.DocumentoVendaHandler.AddItem) // Corrected: This was missing
		item.GET("/:item", container.DocumentoVendaHandler.GetByID)
		item.PUT("/:item", container.DocumentoVendaHandler.EditItem)
		item.DELETE("/:item", container.DocumentoVendaHandler.DeleteItem)
	}

	// Pagamentos
	pagamento := router.Group("/documentos-venda/:documento_id/pagamentos")
	{
		pagamento.GET("", container.DocumentoVendaHandler.List)
		pagamento.POST("", container.DocumentoVendaHandler.Create) // Corrected: This was missing
		pagamento.GET("/:item", container.DocumentoVendaHandler.GetByID)
		pagamento.PUT("/:item", container.DocumentoVendaHandler.Update)
		pagamento.DELETE("/:item", container.DocumentoVendaHandler.Delete)
	}

	// Condições de Pagamento
	condicao := router.Group("/condicoes-pagamento")
	{
		condicao.GET("", container.CondicaoPagamentoHandler.List)
		condicao.POST("", container.CondicaoPagamentoHandler.Create)
		condicao.GET("/:id", container.CondicaoPagamentoHandler.GetByID)
		condicao.PUT("/:id", container.CondicaoPagamentoHandler.Update)
		condicao.DELETE("/:id", container.CondicaoPagamentoHandler.Delete)
	}

	{
		documento.GET("", container.DocumentoVendaHandler.List)
		documento.POST("", container.DocumentoVendaHandler.Create) // Corrected: This was missing
		documento.GET("/:documento_id", container.DocumentoVendaHandler.GetByID)
		documento.PUT("/:documento_id", container.DocumentoVendaHandler.Update)
		documento.DELETE("/:documento_id", container.DocumentoVendaHandler.Delete)
		documento.POST("/:documento_id/emitir", container.DocumentoVendaHandler.Emitir)
		documento.POST("/:documento_id/cancelar", container.DocumentoVendaHandler.Cancelar)
	}
}
