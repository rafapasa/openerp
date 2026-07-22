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

type DocumentoVendaItemService struct {
	db             *gorm.DB
	docRepo        *repository.DocumentoVendaRepository
	itemRepo       *repository.DocumentoVendaItemRepository
	produtoService *ProdutoService
	configService  *ConfiguracaoService
}

func NewDocumentoVendaItemService(db *gorm.DB) *DocumentoVendaItemService {
	return &DocumentoVendaItemService{
		db:             db,
		docRepo:        repository.NewDocumentoVendaRepository(db),
		itemRepo:       repository.NewDocumentoVendaItemRepository(db),
		produtoService: NewProdutoService(db),
		configService:  NewConfiguracaoService(db),
	}
}

func (s *DocumentoVendaItemService) Create(req *dto.DocumentoVendaItemRequest) (*models.DocumentoVendaItem, error) {
	// if err := s.isDataValid(req); err != nil {
	// 	return nil, err
	// }

	if req.ProdutoID <= 0 {
		return nil, apperrors.NewValidationError("O 'produto_id' é obrigatório.")
	}

	produto, err := s.produtoService.FindById(req.ProdutoID)
	if err != nil {
		return nil, apperrors.NewNotFoundError(fmt.Sprintf("Produto com ID %d não encontrado.", req.ProdutoID))
	}
	if !produto.IsActive() {
		return nil, apperrors.NewValidationError(fmt.Sprintf("O produto '%s' não está ativo.", produto.GetNomeCompleto()))
	}

	config, err := s.configService.GetConfig(constants.CONFIG_VALOR_UNITARIO_VENDA_HABILITADO)
	configTbpId, err := s.configService.GetConfig(constants.CONFIG_TABELA_PRECO_PADRAO)
	if err != nil {
		return nil, err
	}

	item := &models.DocumentoVendaItem{}
	if err := utils.MapToModel(req, item); err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear dados do item.", err)
	}
	if config.(int) != constants.SIM {
		valor, err := s.produtoService.GetValorUnitario(configTbpId.(int), req.ProdutoID)
		if err != nil {
			return nil, err
		}
		item.ValorUnitario = valor
	}
	return nil, nil
}
