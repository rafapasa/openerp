package service

import (
	"fmt"

	"github.com/openerp/backend/internal/dto"
	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/utils"
	"gorm.io/gorm"
)

type DocumentoVendaPagamentoService struct {
	db      *gorm.DB
	dvpRepo *repository.DocumentoVendaPagamentoRepository
}

func NewDocumentoVendaPagamentoService(db *gorm.DB) *DocumentoVendaPagamentoService {
	return &DocumentoVendaPagamentoService{
		db:      db,
		dvpRepo: repository.NewDocumentoVendaPagamentoRepository(db),
	}
}

func (s *DocumentoVendaPagamentoService) Create(req *dto.DocumentoVendaPagamentoRequest) error {
	pagamento := &models.DocumentoVendaPagamento{}
	if err := utils.MapToModel(req, pagamento); err != nil {
		return apperrors.NewInternalError("Erro ao mapear dados do pagamento.", err)
	}

	return s.dvpRepo.Create(pagamento)
}

func (s *DocumentoVendaPagamentoService) Update(ddvId, dvpItem int, req *dto.DocumentoVendaPagamentoRequest) error {
	pagamento, err := s.dvpRepo.FindByID(ddvId, dvpItem)
	if err != nil {
		return apperrors.NewNotFoundError(fmt.Sprintf("Pagamento %d do documento %d não encontrado.", dvpItem, ddvId))
	}

	if err := utils.MapToModel(req, pagamento); err != nil {
		return apperrors.NewInternalError("Erro ao mapear dados do pagamento.", err)
	}

	return s.dvpRepo.Update(ddvId, dvpItem, pagamento)
}

func (s *DocumentoVendaPagamentoService) Delete(ddvId, dvpItem int) error {
	return s.dvpRepo.Delete(ddvId, dvpItem)
}

func (s *DocumentoVendaPagamentoService) ListByDocumentoVendaID(limit, offset int, ddvId int) ([]models.DocumentoVendaPagamento, int64, error) {
	return s.dvpRepo.ListByDocumentoVendaID(limit, offset, ddvId)
}

func (s *DocumentoVendaPagamentoService) FindByID(ddvId, dvpItem int) (*models.DocumentoVendaPagamento, error) {
	return s.dvpRepo.FindByID(ddvId, dvpItem)
}