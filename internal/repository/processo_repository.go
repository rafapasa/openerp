package repository

import (
	"fmt"

	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"gorm.io/gorm"
)

type ProcessoRepository struct {
	db *gorm.DB
}

func NewProcessoRepository(db *gorm.DB) *ProcessoRepository {
	return &ProcessoRepository{db: db}
}

func (r *ProcessoRepository) FindByID(id int) (*models.Processo, error) {
	var processo models.Processo
	err := r.db.
		Preload("OperacaoFiscal").
		First(&processo, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("Processo com ID %d não encontrado", id))
		}
		return nil, err
	}
	return &processo, nil

}
