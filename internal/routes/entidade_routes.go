package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/openerp/backend/internal/handler"
)

// RegisterEntidadeRoutes registra as rotas de entidade e seus sub-recursos
func RegisterEntidadeRoutes(
	router *gin.RouterGroup,
	entidadeHandler *handler.EntidadeHandler,
	enderecoHandler *handler.EntidadeEnderecoHandler,
	contatoHandler *handler.EntidadeContatoHandler,
	documentoHandler *handler.EntidadeDocumentoHandler,
	regimeHandler *handler.EntidadeRegimeTributarioHandler,
	limiteCreditoHandler *handler.EntidadeLimiteCreditoHandler,
) {
	entidades := router.Group("/entidades")
	{
		// ============================================================
		// ROTAS MAIS ESPECÍFICAS (primeiro)
		// ============================================================
		entidades.GET("/documento/:documento", entidadeHandler.GetByDocumento)

		// ============================================================
		// SUB-RECURSOS: ENDEREÇOS
		// ============================================================
		enderecos := entidades.Group("/:id/enderecos")
		{
			enderecos.POST("/", enderecoHandler.Create)
			enderecos.GET("/", enderecoHandler.List)
			enderecos.GET("/:item", enderecoHandler.GetByID)
			enderecos.PUT("/:item", enderecoHandler.Update)
			enderecos.DELETE("/:item", enderecoHandler.Delete)
		}

		// ============================================================
		// SUB-RECURSOS: CONTATOS
		// ============================================================
		contatos := entidades.Group("/:id/contatos")
		{
			contatos.POST("/", contatoHandler.Create)
			contatos.GET("/", contatoHandler.List)
			contatos.GET("/:item", contatoHandler.GetByID)
			contatos.PUT("/:item", contatoHandler.Update)
			contatos.DELETE("/:item", contatoHandler.Delete)
		}

		// ============================================================
		// SUB-RECURSOS: DOCUMENTOS
		// ============================================================
		documentos := entidades.Group("/:id/documentos")
		{
			documentos.POST("/", documentoHandler.Create)
			documentos.GET("/", documentoHandler.List)
			documentos.GET("/:item", documentoHandler.GetByID)
			documentos.PUT("/:item", documentoHandler.Update)
			documentos.DELETE("/:item", documentoHandler.Delete)
			documentos.GET("/:item/download", documentoHandler.Download) // Download do arquivo
		}

		// ============================================================
		// SUB-RECURSOS: REGIMES TRIBUTÁRIOS
		// ============================================================
		RegisterEntidadeRegimeTributarioRoutes(entidades, regimeHandler)
		RegisterEntidadeLimiteCreditoRoutes(entidades, limiteCreditoHandler) // Limites de crédito

		// ============================================================
		// ROTAS PRINCIPAIS (mais genéricas, por último)
		// ============================================================
		entidades.POST("/", entidadeHandler.Create)
		entidades.GET("/", entidadeHandler.List)
		entidades.GET("/:id", entidadeHandler.GetByID)
		entidades.PUT("/:id", entidadeHandler.Update)
		entidades.DELETE("/:id", entidadeHandler.Delete)
	}
}
