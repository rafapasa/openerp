//go:build wireinject
// +build wireinject

package appvalidation

import (
	"github.com/google/wire"
)

// ProviderSet é o conjunto de providers para o pacote de validação
var ProviderSet = wire.NewSet(
	NewPlayValidator,
	NewDocumentValidator,
	// Bind interfaces a implementações
	wire.Bind(new(Validator), new(*playValidator)),
	wire.Bind(new(DocumentValidator), new(*documentValidator)),
)

// ValidateProvider é um provider simplificado que retorna a interface
// Usado quando você só precisa do Validator
var ValidateProvider = wire.NewSet(
	NewPlayValidator,
	wire.Bind(new(Validator), new(*playValidator)),
)
