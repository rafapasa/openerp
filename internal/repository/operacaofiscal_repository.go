package repository

import (
	"fmt"

	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"gorm.io/gorm"
)

type OperacaoFiscalRepository struct {
	db *gorm.DB
}

func NewOperacaoFiscalRepository(db *gorm.DB) *OperacaoFiscalRepository {
	return &OperacaoFiscalRepository{db: db}
}

func (r *OperacaoFiscalRepository) FindByID(id int) (*models.OperacaoFiscal, error) {
	var opf models.OperacaoFiscal
	err := r.db.
		Preload("CSTIPI").
		Preload("CSTICMS").
		Preload("CSTPISCOFINS").
		First(&opf, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("Operação fiscal com ID %d não encontrada", id))
		}
		return nil, err
	}
	return &opf, nil
}
