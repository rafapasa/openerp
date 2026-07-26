package service

import (
	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"gorm.io/gorm"
)

type ProcessoService struct {
	db         *gorm.DB
	prcRepo    *repository.ProcessoRepository
	opfService *OperacaoFiscalService
}

func (s *ProcessoService) GetOperacaoFiscal(prcId int, opInterna, opSubTrib bool) (*models.OperacaoFiscal, error) {
	prc, err := s.prcRepo.FindByID(prcId)
	if err != nil {
		return nil, err
	}
	if opInterna && opSubTrib {
		return s.opfService.FindByID(*prc.OperacaoFiscalNoEstSTID)
	}
	if opInterna && !opSubTrib {
		return s.opfService.FindByID(*prc.OperacaoFiscalNoEstID)
	}
	if !opInterna && opSubTrib {
		return s.opfService.FindByID(*prc.OperacaoFiscalForaEstSTID)
	}
	if !opInterna && opSubTrib {
		return s.opfService.FindByID(*prc.OperacaoFiscalForaEstID)
	}
	return nil, apperrors.NewInternalError("Operação fiscal não encontrada.", nil)
}

func NewProcessoService(db *gorm.DB) *ProcessoService {
	return &ProcessoService{
		db:      db,
		prcRepo: repository.NewProcessoRepository(db),
	}
}

func (s *ProcessoService) FindByID(id int) (*models.Processo, error) {
	if id <= 0 {
		return nil, apperrors.NewValidationError("ID do processo inválido.")
	}
	return s.prcRepo.FindByID(id)

}
