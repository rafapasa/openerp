package service

import (
	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
)

// OperacaoFiscalService define os métodos públicos para o serviço de operação fiscal.
type OperacaoFiscalService interface {
	FindByID(id int) (*models.OperacaoFiscal, error)
}

type operacaoFiscalService struct {
	opfRepo repository.OperacaoFiscalRepository
}

func NewOperacaoFiscalService(opfRepo repository.OperacaoFiscalRepository) OperacaoFiscalService { // Removed db parameter as it's not used
	return &operacaoFiscalService{
		opfRepo: opfRepo,
	}
}

func (s *operacaoFiscalService) FindByID(id int) (*models.OperacaoFiscal, error) {
	if id <= 0 {
		return nil, apperrors.NewValidationError("ID da operação fiscal inválido.")
	}
	return s.opfRepo.FindByID(id)
}
