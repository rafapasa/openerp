// internal/modules/product/module.go
package product

import (
	"github.com/openerp/backend/internal/container"
	"github.com/openerp/backend/internal/handler"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/service"
)

const ModuleName = "product"

// ProductModule agrupa todas as dependências relacionadas a produto
type ProductModule struct {
	Handlers map[string]interface{}
}

func NewProductModule() *ProductModule {
	return &ProductModule{
		Handlers: make(map[string]interface{}),
	}
}

// Register registra todas as dependências do módulo de produto
func (m *ProductModule) Register(container *container.Container) error {
	db := container.GetDB()

	// 1. Repositórios
	produtoRepo := repository.NewProdutoRepository(db)
	grupoRepo := repository.NewProdutoGrupoRepository(db)
	subgrupoRepo := repository.NewProdutoSubgrupoRepository(db)
	marcaRepo := repository.NewProdutoMarcaRepository(db)
	modeloRepo := repository.NewProdutoModeloRepository(db)
	variacaoRepo := repository.NewProdutoVariacaoRepository(db)

	container.RegisterRepository("produto", produtoRepo)
	container.RegisterRepository("produto_grupo", grupoRepo)
	container.RegisterRepository("produto_subgrupo", subgrupoRepo)
	container.RegisterRepository("produto_marca", marcaRepo)
	container.RegisterRepository("produto_modelo", modeloRepo)
	container.RegisterRepository("produto_variacao", variacaoRepo)

	// 2. Services
	produtoService := service.NewProdutoService(produtoRepo)
	grupoService := service.NewProdutoGrupoService(grupoRepo)
	subgrupoService := service.NewProdutoSubgrupoService(subgrupoRepo)
	marcaService := service.NewProdutoMarcaService(marcaRepo)
	modeloService := service.NewProdutoModeloService(modeloRepo)
	// variacaoService := service.NewProdutoVariacaoService(variacaoRepo)

	container.RegisterService("produto", produtoService)
	container.RegisterService("produto_grupo", grupoService)
	container.RegisterService("produto_subgrupo", subgrupoService)
	container.RegisterService("produto_marca", marcaService)
	container.RegisterService("produto_modelo", modeloService)
	// container.RegisterService("produto_variacao", variacaoService)

	// 3. Handlers (usando os existentes)
	produtoHandler := handler.NewProdutoHandler(produtoService)
	grupoHandler := handler.NewProdutoGrupoHandler(grupoService)
	subgrupoHandler := handler.NewProdutoSubgrupoHandler(subgrupoService)
	marcaHandler := handler.NewProdutoMarcaHandler(marcaService)
	modeloHandler := handler.NewProdutoModeloHandler(modeloService)
	// variacaoHandler := handler.NewProdutoVariacaoHandler(variacaoService)

	container.RegisterHandler("produto", produtoHandler)
	container.RegisterHandler("produto_grupo", grupoHandler)
	container.RegisterHandler("produto_subgrupo", subgrupoHandler)
	container.RegisterHandler("produto_marca", marcaHandler)
	container.RegisterHandler("produto_modelo", modeloHandler)
	// container.RegisterHandler("produto_variacao", variacaoHandler)

	// Armazena os handlers para acesso rápido
	m.Handlers["produto"] = produtoHandler
	m.Handlers["produto_grupo"] = grupoHandler
	m.Handlers["produto_subgrupo"] = subgrupoHandler
	m.Handlers["produto_marca"] = marcaHandler
	m.Handlers["produto_modelo"] = modeloHandler
	// m.Handlers["produto_variacao"] = variacaoHandler

	return nil
}

// GetHandlers retorna todos os handlers do módulo
func (m *ProductModule) GetHandlers() map[string]interface{} {
	return m.Handlers
}
