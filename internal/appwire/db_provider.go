package appwire

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DBProvider é um tipo de função que fornece uma instância de *gorm.DB
// que está vinculada ao contexto da requisição atual.
type DBProvider func(c *gin.Context) (*gorm.DB, error)

// NewDBProvider cria uma função DBProvider.
// Esta função será injetada pelo Wire nos handlers.
func NewDBProvider() DBProvider {
	return func(c *gin.Context) (*gorm.DB, error) {
		dbWithContext, exists := c.Get("gorm_db_with_context")
		if !exists {
			return nil, errors.New("database context not found in request")
		}
		db, ok := dbWithContext.(*gorm.DB)
		if !ok {
			return nil, errors.New("invalid database context type")
		}
		return db, nil
	}
}
