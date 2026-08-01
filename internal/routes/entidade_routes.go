// internal/routes/entidade_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/appwire"
)

func RegisterEntidadeRoutes(router *gin.RouterGroup, container *appwire.Container) {
	// Grupo de rotas de entidade
	entidade := router.Group("/entidades")
	{
		entidade.GET("", container.EntidadeHandler.List)
		entidade.POST("", container.EntidadeHandler.Create)
		entidade.GET("/:id", container.EntidadeHandler.GetByID)
		entidade.PUT("/:id", container.EntidadeHandler.Update)
		entidade.DELETE("/:id", container.EntidadeHandler.Delete) // Corrected: This was missing
		// entidade.GET("/documento/:documento", container.EntidadeHandler.GetByDocumento)
	}

	// Endereços
	endereco := router.Group("/entidades/:entidade_id/enderecos")
	{
		endereco.GET("", container.EntidadeEnderecoHandler.List)
		endereco.POST("", container.EntidadeEnderecoHandler.Create)
		endereco.GET("/:item", container.EntidadeEnderecoHandler.GetByID)
		endereco.PUT("/:item", container.EntidadeEnderecoHandler.Update)
		endereco.DELETE("/:item", container.EntidadeEnderecoHandler.Delete)
	}

	// Contatos
	contato := router.Group("/entidades/:entidade_id/contatos")
	{
		contato.GET("", container.EntidadeContatoHandler.List)
		contato.POST("", container.EntidadeContatoHandler.Create)
		contato.GET("/:item", container.EntidadeContatoHandler.GetByID)
		contato.PUT("/:item", container.EntidadeContatoHandler.Update)
		contato.DELETE("/:item", container.EntidadeContatoHandler.Delete)
	}

	// Documentos
	documento := router.Group("/entidades/:entidade_id/documentos")
	{
		documento.GET("", container.EntidadeDocumentoHandler.List)
		documento.POST("", container.EntidadeDocumentoHandler.Create)
		documento.GET("/:item", container.EntidadeDocumentoHandler.GetByID)
		documento.PUT("/:item", container.EntidadeDocumentoHandler.Update) // Corrected: This was missing
		documento.DELETE("/:item", container.EntidadeDocumentoHandler.Delete)
		documento.GET("/:item/download", container.EntidadeDocumentoHandler.Download)
	}

	// Limites de Crédito
	limite := router.Group("/entidades/:entidade_id/limites-credito")
	{
		limite.GET("", container.EntidadeLimiteCreditoHandler.List)
		limite.POST("", container.EntidadeLimiteCreditoHandler.Create)
		limite.GET("/:id", container.EntidadeLimiteCreditoHandler.GetByID) // Corrected: This was missing
		limite.PUT("/:id", container.EntidadeLimiteCreditoHandler.Update)
		limite.DELETE("/:id", container.EntidadeLimiteCreditoHandler.Delete)
	}

	// Regimes Tributários
	regime := router.Group("/entidades/:entidade_id/regimes-tributarios")
	{
		regime.GET("", container.EntidadeRegimeTributarioHandler.List)
		regime.POST("", container.EntidadeRegimeTributarioHandler.Create) // Corrected: This was missing
		regime.GET("/:item", container.EntidadeRegimeTributarioHandler.GetByID)
		regime.PUT("/:item", container.EntidadeRegimeTributarioHandler.Update)
		regime.DELETE("/:item", container.EntidadeRegimeTributarioHandler.Delete)
	}

	// Grupos de Entidade
	grupo := router.Group("/grupos-entidade")
	{
		grupo.GET("", container.GrupoEntidadeHandler.List) // Corrected: This was missing
		grupo.POST("", container.GrupoEntidadeHandler.Create)
		grupo.GET("/:id", container.GrupoEntidadeHandler.GetByID)
		grupo.PUT("/:id", container.GrupoEntidadeHandler.Update)
		grupo.DELETE("/:id", container.GrupoEntidadeHandler.Delete)
	}
}
