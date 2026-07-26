package service

import (
	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"gorm.io/gorm"
)

type OperacaoFiscalService struct {
	db      *gorm.DB
	opfRepo *repository.OperacaoFiscalRepository
}

func NewOperacaoFiscalService(db *gorm.DB) *OperacaoFiscalService {
	return &OperacaoFiscalService{
		db: db,
	}
}

func (s *OperacaoFiscalService) FindByID(id int) (*models.OperacaoFiscal, error) {
	if id <= 0 {
		return nil, apperrors.NewValidationError("ID da operação fiscal inválido.")
	}
	return s.opfRepo.FindByID(id)
}
