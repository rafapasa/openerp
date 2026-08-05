package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/handler"
)

// RegisterGrupoUsuarioRoutes configura as rotas para o recurso GrupoUsuario.
func RegisterGrupoUsuarioRoutes(router *gin.Engine, h *handler.GrupoUsuarioHandler) {
	grupoUsuarioRoutes := router.Group("/grupos-usuario")
	{
		// Rotas públicas (ou com autenticação JWT, dependendo da configuração global)
		grupoUsuarioRoutes.POST("", h.CreateGrupoUsuario)
		grupoUsuarioRoutes.GET("/:id", h.GetGrupoUsuarioByID)
		grupoUsuarioRoutes.PUT("/:id", h.UpdateGrupoUsuario)
		grupoUsuarioRoutes.DELETE("/:id", h.DeleteGrupoUsuario)
		grupoUsuarioRoutes.GET("", h.ListGrupoUsuarios)

		// Adicione aqui rotas que necessitem de middlewares específicos, como autorização
	}
}