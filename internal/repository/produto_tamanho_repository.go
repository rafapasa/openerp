package repository

import (
	"fmt"

	"github.com/openerp/backend/internal/models"
	"gorm.io/gorm"
)

type ProdutoTamanhoRepository struct {
	db *gorm.DB
}

func NewProdutoTamanhoRepository(db *gorm.DB) *ProdutoTamanhoRepository {
	return &ProdutoTamanhoRepository{db: db}
}

func (r *ProdutoTamanhoRepository) FindByID(id int) (*models.ProdutoTamanho, error) {
	var tamanho models.ProdutoTamanho
	err := r.db.First(&tamanho, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil, nil if not found
		}
		return nil, fmt.Errorf("erro ao buscar tamanho do produto: %w", err)
	}
	return &tamanho, nil
}