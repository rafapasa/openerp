package service

import (
	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
)

// ProcessoService define os métodos públicos para o serviço de processo.
type ProcessoService interface {
	GetOperacaoFiscal(prcId int, opInterna, opSubTrib bool) (*models.OperacaoFiscal, error)
	FindByID(id int) (*models.Processo, error)
}

type processoService struct {
	prcRepo    repository.ProcessoRepository
	opfService OperacaoFiscalService
}

func NewProcessoService(prcRepo repository.ProcessoRepository, opfService OperacaoFiscalService) ProcessoService {
	return &processoService{
		prcRepo:    prcRepo,
		opfService: opfService,
	}
}

func (s *processoService) GetOperacaoFiscal(prcId int, opInterna, opSubTrib bool) (*models.OperacaoFiscal, error) {
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
	if !opInterna && !opSubTrib { // Corrected: Changed condition
		return s.opfService.FindByID(*prc.OperacaoFiscalForaEstID)
	}
	return nil, apperrors.NewInternalError("Operação fiscal não encontrada.", nil)
}

func (s *processoService) FindByID(id int) (*models.Processo, error) { // Added this method to satisfy the interface
	if id <= 0 {
		return nil, apperrors.NewValidationError("ID do processo inválido.")
	}
	return s.prcRepo.FindByID(id)

}
