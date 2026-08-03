// internal/middleware/gorm_context_middleware.go
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/appcontext"
	"gorm.io/gorm"
)

// GormContextMiddleware injeta o userID no contexto do GORM para cada requisição
func GormContextMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Obter userID do contexto do Gin
		userID := GetUserID(c)

		// 2. ✅ Criar contexto com o userID
		ctx := appcontext.WithUserID(c.Request.Context(), userID)

		// 3. ✅ Substituir o contexto da requisição
		c.Request = c.Request.WithContext(ctx)

		// 4. ✅ Salvar o DB com contexto no Gin para uso nos handlers
		c.Set("db_with_context", db.WithContext(ctx))

		c.Next()
	}
}

// GetDBWithContext obtém o DB com contexto da requisição atual
func GetDBWithContext(c *gin.Context) *gorm.DB {
	if db, exists := c.Get("db_with_context"); exists {
		if gormDB, ok := db.(*gorm.DB); ok {
			return gormDB
		}
	}
	return nil
}
