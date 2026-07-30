package service

import (
	"fmt"

	"github.com/openerp/backend/internal/constants"
	"github.com/openerp/backend/internal/dto"
	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/utils"
	"gorm.io/gorm"
)

// DocumentoVendaItemServiceFactory é uma fábrica para criar instâncias do service
type DocumentoVendaItemServiceFactory interface {
	// CreateWithTx cria uma instância com um DB transacional
	CreateWithTx(tx *gorm.DB) DocumentoVendaItemService
}

// DocumentoVendaItemServiceInterface define os métodos públicos para o serviço de itens de documento de venda.
type DocumentoVendaItemService interface {
	Create(req *dto.DocumentoVendaItemRequest) error
	Update(ddvId, dviItem int, req *dto.DocumentoVendaItemRequest) error
	GetByID(ddvId, dviItem int) (*models.DocumentoVendaItem, error)
	List(limit, offset, ddvId int, filters map[string]interface{}) ([]models.DocumentoVendaItem, int64, error)
	Delete(ddvId, dviItem int) error
}

// documentoVendaItemServiceFactory implementa a fábrica
type documentoVendaItemServiceFactory struct {
	// Dependências fixas (não mudam com a transação)
	proService    ProdutoService
	prcService    ProcessoService
	configService ConfiguracaoService
}

// NewDocumentoVendaItemServiceFactory cria a fábrica com as dependências injetadas
func NewDocumentoVendaItemServiceFactory(
	proService ProdutoService,
	prcService ProcessoService,
	configService ConfiguracaoService,
) DocumentoVendaItemServiceFactory {
	return &documentoVendaItemServiceFactory{
		proService:    proService,
		prcService:    prcService,
		configService: configService,
	}
}

// CreateWithTx cria uma nova instância com DB transacional
func (f *documentoVendaItemServiceFactory) CreateWithTx(tx *gorm.DB) DocumentoVendaItemService {
	// Cria o repositório com o DB transacional
	dviRepo := repository.NewDocumentoVendaItemRepository(tx)

	// Cria o service com as dependências fixas + repositório transacional
	return &documentoVendaItemService{
		dviRepo:       dviRepo,
		proService:    f.proService,
		prcService:    f.prcService,
		configService: f.configService,
	}
}

// documentoVendaItemService é a implementação concreta de DocumentoVendaItemServiceInterface.
type documentoVendaItemService struct {
	ddvService    DocumentoVendaService
	dviRepo       repository.DocumentoVendaItemRepository
	proService    ProdutoService
	prcService    ProcessoService
	configService ConfiguracaoService
}

func NewDocumentoVendaItemService(ddvService DocumentoVendaService,
	dviRepo repository.DocumentoVendaItemRepository,
	proService ProdutoService,
	prcService ProcessoService,
	configService ConfiguracaoService) DocumentoVendaItemService {
	return &documentoVendaItemService{
		ddvService:    ddvService,
		dviRepo:       dviRepo,
		proService:    proService,
		prcService:    prcService,
		configService: configService}
}

func (s *documentoVendaItemService) Create(req *dto.DocumentoVendaItemRequest) error {
	if err := s.isDataValid(req); err != nil {
		return err
	}

	configEditValUnit, err := s.configService.GetConfig(constants.CONFIG_VALOR_UNITARIO_VENDA_HABILITADO)
	configTbpId, err := s.configService.GetConfig(constants.CONFIG_TABELA_PRECO_PADRAO)
	if err != nil {
		return err
	}

	dvi := &models.DocumentoVendaItem{}
	if err := utils.MapToModel(req, dvi); err != nil {
		return apperrors.NewInternalError("Erro ao mapear dados do item.", err)
	}

	if configEditValUnit.(int) != constants.SIM {
		valor, err := s.proService.GetValorUnitario(configTbpId.(int), req.ProdutoID)
		if err != nil {
			return err
		}
		dvi.ValorUnitario = valor
	} else {
		dvi.ValorUnitario = req.ValorUnitario
	}

	s.recalcularValoresItem(dvi)

	ddv, err := s.ddvService.GetByID(req.DocumentoVendaID)
	if err != nil {
		return err
	}

	opInterna := true
	opSubTrib := true

	prcId := ddv.ProcessoID
	opf, err := s.prcService.GetOperacaoFiscal(prcId, opInterna, opSubTrib)
	if err != nil {
		return apperrors.NewInternalError("Erro ao buscar operação fiscal.", err)
	}

	dvi.OperacaoFiscalID = opf.ID
	dvi.CSTICMSID = *opf.CSTICMSID
	dvi.CSTIPIID = *opf.CSTIPIID
	dvi.CSTPISCOFINSID = *opf.CSTPISCOFINSID

	return s.dviRepo.Create(dvi)
}

func (s *documentoVendaItemService) Update(ddvId, dviItem int, req *dto.DocumentoVendaItemRequest) error {
	if err := s.isDataValid(req); err != nil {
		return err
	}

	item, err := s.dviRepo.FindByID(ddvId, dviItem)
	if err != nil {
		return apperrors.NewNotFoundError(fmt.Sprintf("Item %d do documento %d não encontrado.", dviItem, ddvId))
	}

	if err := utils.MapToModel(req, item); err != nil {
		return apperrors.NewInternalError("Erro ao mapear dados do item.", err)
	}

	s.recalcularValoresItem(item)

	return s.dviRepo.Update(ddvId, dviItem, item)
}

func (s *documentoVendaItemService) Delete(ddvId, dviItem int) error {
	return s.dviRepo.Delete(ddvId, dviItem)
}

// GetByID busca um item de documento de venda específico.
func (s *documentoVendaItemService) GetByID(ddvId, dviItem int) (*models.DocumentoVendaItem, error) {
	item, err := s.dviRepo.FindByID(ddvId, dviItem)
	if err != nil {
		return nil, apperrors.NewNotFoundError(fmt.Sprintf("Item %d do documento %d não encontrado.", dviItem, ddvId))
	}
	return item, nil
}

// List lista itens de documento de venda com paginação e filtros.
func (s *documentoVendaItemService) List(limit, offset, ddvId int, filters map[string]interface{}) ([]models.DocumentoVendaItem, int64, error) {
	return s.dviRepo.ListByDocumentoVendaID(limit, offset, ddvId)
}

// recalcularValoresItem calcula os totais do item.
func (s *documentoVendaItemService) recalcularValoresItem(item *models.DocumentoVendaItem) {
	item.TotalProdutos = item.Quantidade * item.ValorUnitario
	item.TotalItem = item.TotalProdutos
	if item.ValorDesconto != nil {
		item.TotalItem -= *item.ValorDesconto
	}
}

func (s *documentoVendaItemService) isDataValid(req *dto.DocumentoVendaItemRequest) error {
	if req.ProdutoID <= 0 {
		return apperrors.NewValidationError("O 'produto_id' é obrigatório.")
	}
	if req.DocumentoVendaID <= 0 {
		return apperrors.NewValidationError("O 'documento_venda_id' é obrigatório.")
	}
	if req.Quantidade <= 0 {
		return apperrors.NewValidationError("A 'quantidade' deve ser maior que zero.")
	}

	produto, err := s.proService.FindById(req.ProdutoID)
	if err != nil { //
		return apperrors.NewNotFoundError(fmt.Sprintf("Produto com ID %d não encontrado.", req.ProdutoID)) //
	}
	if !produto.IsActive() {
		return apperrors.NewValidationError(fmt.Sprintf("O produto '%s' não está ativo.", produto.GetNomeCompleto())) //
	}
	return nil
}
