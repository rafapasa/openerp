package wire

import (
	"github.com/google/wire"

	"github.com/openerp/backend/internal/handler"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/service"
)

// ============================================================
// REPOSITORY PROVIDERS
// ============================================================

// Repositories de Entidade
var EntidadeRepositories = wire.NewSet(
	repository.NewEntidadeRepository,
	repository.NewEntidadeEnderecoRepository,
	repository.NewEntidadeContatoRepository,
	repository.NewEntidadeDocumentoRepository,
	repository.NewEntidadeLimiteCreditoRepository,
	repository.NewEntidadeRegimeTributarioRepository,
	repository.NewGrupoEntidadeRepository,
)

// Repositories de Produto
var ProdutoRepositories = wire.NewSet(
	repository.NewProdutoRepository,
	repository.NewProdutoGrupoRepository,
	repository.NewProdutoSubgrupoRepository,
	repository.NewProdutoMarcaRepository,
	repository.NewProdutoModeloRepository,
	repository.NewProdutoCorRepository,
	repository.NewProdutoTamanhoRepository,
	repository.NewProdutoVariacaoRepository,
)

// Repositories de Venda
var VendaRepositories = wire.NewSet(
	repository.NewDocumentoVendaRepository,
	repository.NewDocumentoVendaItemRepository,
	repository.NewDocumentoVendaPagamentoRepository,
	repository.NewCondicaoPagamentoRepository,
)

// Repositories de Fiscal
var FiscalRepositories = wire.NewSet(
	repository.NewOperacaoFiscalRepository,
	repository.NewProcessoRepository,
)

// Repositories de Tabela Preço
var TabelaPrecoRepositories = wire.NewSet(
	repository.NewTabelaPrecoRepository,
	repository.NewTabelaPrecoProdutoRepository,
)

// Todos os Repositories
var AllRepositories = wire.NewSet(
	EntidadeRepositories,
	ProdutoRepositories,
	VendaRepositories,
	FiscalRepositories,
	TabelaPrecoRepositories,
)

// ============================================================
// SERVICE PROVIDERS
// ============================================================

// Services de Entidade
var EntidadeServices = wire.NewSet(
	service.NewEntidadeService,
	service.NewEntidadeEnderecoService,
	service.NewEntidadeContatoService,
	service.NewEntidadeDocumentoService,
	service.NewEntidadeLimiteCreditoService,
	service.NewEntidadeRegimeTributarioService,
	service.NewGrupoEntidadeService,
)

// Services de Produto
var ProdutoServices = wire.NewSet(
	service.NewProdutoService,
	service.NewProdutoGrupoService,
	service.NewProdutoSubgrupoService,
	service.NewProdutoMarcaService,
	service.NewProdutoModeloService,
	service.NewProdutoCorService,
	service.NewProdutoTamanhoService,
	service.NewProdutoVariacaoService,
)

// Services de Venda
var VendaServices = wire.NewSet(
	service.NewDocumentoVendaService,
	service.NewDocumentoVendaItemService,
	service.NewDocumentoVendaPagamentoService,
	// service.NewCondicaoPagamentoService,
)

// Services de Fiscal
var FiscalServices = wire.NewSet(
	service.NewOperacaoFiscalService,
	service.NewProcessoService,
)

// Services de Tabela Preço
var TabelaPrecoServices = wire.NewSet(
	service.NewTabelaPrecoService,
	service.NewTabelaPrecoProdutoService,
)

// Todos os Services
var AllServices = wire.NewSet(
	EntidadeServices,
	ProdutoServices,
	VendaServices,
	FiscalServices,
	TabelaPrecoServices,
)

// ============================================================
// HANDLER PROVIDERS
// ============================================================

// Handlers de Entidade
var EntidadeHandlers = wire.NewSet(
	handler.NewEntidadeHandler,
	handler.NewEntidadeEnderecoHandler,
	handler.NewEntidadeContatoHandler,
	handler.NewEntidadeDocumentoHandler,
	handler.NewEntidadeLimiteCreditoHandler,
	handler.NewEntidadeRegimeTributarioHandler,
	handler.NewGrupoEntidadeHandler,
)

// Handlers de Produto
var ProdutoHandlers = wire.NewSet(
	handler.NewProdutoHandler,
	handler.NewProdutoGrupoHandler,
	handler.NewProdutoSubgrupoHandler,
	handler.NewProdutoMarcaHandler,
	handler.NewProdutoModeloHandler,
	// handler.NewProdutoCorHandler,
	// handler.NewProdutoTamanhoHandler,
	handler.NewProdutoVariacaoHandler,
)

// Handlers de Venda
var VendaHandlers = wire.NewSet(
	handler.NewDocumentoVendaHandler,
	handler.NewDocumentoVendaItemHandler,
	handler.NewDocumentoVendaPagamentoHandler,
	handler.NewCondicaoPagamentoHandler,
)

// Handlers de Fiscal
var FiscalHandlers = wire.NewSet(
	handler.NewOperacaoFiscalHandler,
	handler.NewProcessoHandler,
)

// Handlers de Tabela Preço
var TabelaPrecoHandlers = wire.NewSet(
	handler.NewTabelaPrecoHandler,
	handler.NewTabelaPrecoProdutoHandler,
)

// Todos os Handlers
var AllHandlers = wire.NewSet(
	EntidadeHandlers,
	ProdutoHandlers,
	VendaHandlers,
	FiscalHandlers,
	TabelaPrecoHandlers,
)

// ============================================================
// MÓDULOS COMPLETOS
// ============================================================

var EntidadeModule = wire.NewSet(
	EntidadeRepositories,
	EntidadeServices,
	EntidadeHandlers,
)

var ProdutoModule = wire.NewSet(
	ProdutoRepositories,
	ProdutoServices,
	ProdutoHandlers,
)

var VendaModule = wire.NewSet(
	VendaRepositories,
	VendaServices,
	VendaHandlers,
)

var FiscalModule = wire.NewSet(
	FiscalRepositories,
	FiscalServices,
	FiscalHandlers,
)

var TabelaPrecoModule = wire.NewSet(
	TabelaPrecoRepositories,
	TabelaPrecoServices,
	TabelaPrecoHandlers,
)

var AllModules = wire.NewSet(
	EntidadeModule,
	ProdutoModule,
	VendaModule,
	FiscalModule,
	TabelaPrecoModule,
)
