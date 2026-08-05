package appwire

import (
	"github.com/openerp/backend/internal/handler"
)

// Container agrupa todos os handlers para fácil registro de rotas
type Container struct {
	AuthHandler                     *handler.AuthHandler
	EntidadeHandler                 *handler.EntidadeHandler
	EntidadeEnderecoHandler         *handler.EntidadeEnderecoHandler
	EntidadeContatoHandler          *handler.EntidadeContatoHandler
	EntidadeDocumentoHandler        *handler.EntidadeDocumentoHandler
	EntidadeLimiteCreditoHandler    *handler.EntidadeLimiteCreditoHandler
	EntidadeRegimeTributarioHandler *handler.EntidadeRegimeTributarioHandler
	GrupoEntidadeHandler            *handler.GrupoEntidadeHandler

	// Produto Handlers
	ProdutoHandler         *handler.ProdutoHandler
	ProdutoGrupoHandler    *handler.ProdutoGrupoHandler
	ProdutoSubgrupoHandler *handler.ProdutoSubgrupoHandler
	ProdutoMarcaHandler    *handler.ProdutoMarcaHandler
	ProdutoModeloHandler   *handler.ProdutoModeloHandler
	ProdutoCorHandler      *handler.ProdutoCorHandler
	ProdutoTamanhoHandler  *handler.ProdutoTamanhoHandler
	ProdutoVariacaoHandler *handler.ProdutoVariacaoHandler

	// Venda Handlers
	DocumentoVendaHandler    *handler.DocumentoVendaHandler
	CondicaoPagamentoHandler *handler.CondicaoPagamentoHandler

	// Fiscal Handlers
	OperacaoFiscalHandler *handler.OperacaoFiscalHandler
	ProcessoHandler       *handler.ProcessoHandler

	// Tabela Preço Handlers
	TabelaPrecoHandler        *handler.TabelaPrecoHandler
	TabelaPrecoProdutoHandler *handler.TabelaPrecoProdutoHandler
}

// NewContainer cria um novo container com todos os handlers
func NewContainer(
	authHandler *handler.AuthHandler,
	// Entidade Handlers
	entidadeHandler *handler.EntidadeHandler,
	entidadeEnderecoHandler *handler.EntidadeEnderecoHandler,
	entidadeContatoHandler *handler.EntidadeContatoHandler,
	entidadeDocumentoHandler *handler.EntidadeDocumentoHandler,
	entidadeLimiteCreditoHandler *handler.EntidadeLimiteCreditoHandler,
	entidadeRegimeTributarioHandler *handler.EntidadeRegimeTributarioHandler,
	grupoEntidadeHandler *handler.GrupoEntidadeHandler,
	// Produto Handlers
	produtoHandler *handler.ProdutoHandler,
	produtoGrupoHandler *handler.ProdutoGrupoHandler,
	produtoSubgrupoHandler *handler.ProdutoSubgrupoHandler,
	produtoMarcaHandler *handler.ProdutoMarcaHandler,
	produtoModeloHandler *handler.ProdutoModeloHandler,
	produtoCorHandler *handler.ProdutoCorHandler,
	produtoTamanhoHandler *handler.ProdutoTamanhoHandler,
	produtoVariacaoHandler *handler.ProdutoVariacaoHandler,
	// Venda Handlers
	documentoVendaHandler *handler.DocumentoVendaHandler,
	// Fiscal Handlers
	operacaoFiscalHandler *handler.OperacaoFiscalHandler,
	processoHandler *handler.ProcessoHandler,
	// Tabela Preço Handlers
	tabelaPrecoHandler *handler.TabelaPrecoHandler,
	tabelaPrecoProdutoHandler *handler.TabelaPrecoProdutoHandler,
) *Container {
	return &Container{
		AuthHandler: authHandler,
		// Entidade
		EntidadeHandler:                 entidadeHandler,
		EntidadeEnderecoHandler:         entidadeEnderecoHandler,
		EntidadeContatoHandler:          entidadeContatoHandler,
		EntidadeDocumentoHandler:        entidadeDocumentoHandler,
		EntidadeLimiteCreditoHandler:    entidadeLimiteCreditoHandler,
		EntidadeRegimeTributarioHandler: entidadeRegimeTributarioHandler,
		GrupoEntidadeHandler:            grupoEntidadeHandler,
		// Produto
		ProdutoHandler:         produtoHandler,
		ProdutoGrupoHandler:    produtoGrupoHandler,
		ProdutoSubgrupoHandler: produtoSubgrupoHandler,
		ProdutoMarcaHandler:    produtoMarcaHandler,
		ProdutoModeloHandler:   produtoModeloHandler,
		ProdutoCorHandler:      produtoCorHandler,
		ProdutoTamanhoHandler:  produtoTamanhoHandler,
		ProdutoVariacaoHandler: produtoVariacaoHandler,
		// Venda
		DocumentoVendaHandler: documentoVendaHandler,
		// Fiscal
		OperacaoFiscalHandler: operacaoFiscalHandler,
		ProcessoHandler:       processoHandler,
		// Tabela Preço
		TabelaPrecoHandler:        tabelaPrecoHandler,
		TabelaPrecoProdutoHandler: tabelaPrecoProdutoHandler,
	}
}
