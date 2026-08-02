//go:build wireinject
// +build wireinject

package appwire

import (
	"github.com/google/wire"
	"github.com/openerp/backend/internal/database"
	"gorm.io/gorm"
)

// InitializeContainer injeta todas as dependências e retorna o Container
func InitializeContainer(db *gorm.DB, redis *database.Redis) *Container {
	wire.Build(
		AllModules,
		NewContainer,
	)
	return nil
}
