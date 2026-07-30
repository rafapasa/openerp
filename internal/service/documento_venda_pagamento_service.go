package service

import (
	"github.com/openerp/backend/internal/dto"
	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/utils"
	"gorm.io/gorm"
)

// DocumentoVendaPagamentoService define o contrato para as operações de serviço de pagamento de documento de venda.
type DocumentoVendaPagamentoService interface {
	Create(req *dto.DocumentoVendaPagamentoRequest) error
	Update(ddvId, dvpItem int, req *dto.DocumentoVendaPagamentoRequest) error
	Delete(ddvId, dvpItem int) error
	ListByDocumentoVendaID(limit, offset int, ddvId int) ([]models.DocumentoVendaPagamento, int64, error)
	FindByID(ddvId, dvpItem int) (*models.DocumentoVendaPagamento, error)
}

type documentoVendaPagamentoService struct {
	db      *gorm.DB
	dvpRepo repository.DocumentoVendaPagamentoRepository
}

func NewDocumentoVendaPagamentoService(db *gorm.DB, dvpRepo repository.DocumentoVendaPagamentoRepository) DocumentoVendaPagamentoService {
	return &documentoVendaPagamentoService{
		db:      db,
		dvpRepo: dvpRepo,
	}
}

func (s *documentoVendaPagamentoService) Create(req *dto.DocumentoVendaPagamentoRequest) error { // Renamed from NewdocumentoVendaPagamentoService
	pagamento := &models.DocumentoVendaPagamento{}
	if err := utils.MapToModel(req, pagamento); err != nil {
		return apperrors.NewInternalError("Erro ao mapear dados do pagamento.", err)
	}

	return s.dvpRepo.Create(pagamento)
}

func (s *documentoVendaPagamentoService) Update(ddvId, dvpItem int, req *dto.DocumentoVendaPagamentoRequest) error {
	pagamento, err := s.dvpRepo.FindByID(ddvId, dvpItem)
	if err != nil {
		return err //
	}

	if err := utils.MapToModel(req, pagamento); err != nil {
		return apperrors.NewInternalError("Erro ao mapear dados do pagamento.", err)
	}

	return s.dvpRepo.Update(ddvId, dvpItem, pagamento)
}

func (s *documentoVendaPagamentoService) Delete(ddvId, dvpItem int) error {
	return s.dvpRepo.Delete(ddvId, dvpItem)
}

func (s *documentoVendaPagamentoService) ListByDocumentoVendaID(limit, offset int, ddvId int) ([]models.DocumentoVendaPagamento, int64, error) {
	return s.dvpRepo.ListByDocumentoVendaID(limit, offset, ddvId)
}

func (s *documentoVendaPagamentoService) FindByID(ddvId, dvpItem int) (*models.DocumentoVendaPagamento, error) {
	return s.dvpRepo.FindByID(ddvId, dvpItem)
}
