package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/handler"
)

// RegisterEmpresaFilialRoutes configura as rotas para o recurso EmpresaFilial.
func RegisterEmpresaFilialRoutes(router *gin.Engine, h *handler.EmpresaFilialHandler) {
	empresaFilialRoutes := router.Group("/empresas-filiais")
	{
		// Rotas públicas (ou com autenticação JWT, dependendo da configuração global)
		empresaFilialRoutes.POST("", h.CreateEmpresaFilial)
		empresaFilialRoutes.GET("/:id", h.GetEmpresaFilialByID)
		empresaFilialRoutes.PUT("/:id", h.UpdateEmpresaFilial)
		empresaFilialRoutes.DELETE("/:id", h.DeleteEmpresaFilial)
		empresaFilialRoutes.GET("", h.ListEmpresasFiliais)

		// Adicione aqui rotas que necessitem de middlewares específicos, como autorização
	}
}
