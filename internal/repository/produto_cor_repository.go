package repository

import (
	"fmt"

	"github.com/openerp/backend/internal/models"
	"gorm.io/gorm"
)

type ProdutoCorRepository struct {
	db *gorm.DB
}

func NewProdutoCorRepository(db *gorm.DB) *ProdutoCorRepository {
	return &ProdutoCorRepository{db: db}
}

func (r *ProdutoCorRepository) FindByID(id int) (*models.ProdutoCor, error) {
	var cor models.ProdutoCor
	err := r.db.First(&cor, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil, nil if not found
		}
		return nil, fmt.Errorf("erro ao buscar cor do produto: %w", err)
	}
	return &cor, nil
}