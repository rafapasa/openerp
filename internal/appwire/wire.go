//go:build wireinject
// +build wireinject

package appwire

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

// InitializeContainer injeta todas as dependências e retorna o Container
func InitializeContainer(db *gorm.DB) *Container {
	wire.Build(
		AllModules,
		NewContainer,
	)
	return nil
}
