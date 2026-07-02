package repository

import (
	"errors"

	"github.com/openerp/backend/internal/models"
	"gorm.io/gorm"
)

type ProdutoRepository struct {
	db *gorm.DB
}

func NewProdutoRepository(db *gorm.DB) *ProdutoRepository {
	return &ProdutoRepository{db: db}
}

func (r *ProdutoRepository) Create(produto *models.Produto) error {
	return r.db.Create(produto).Error
}

func (r *ProdutoRepository) Update(produto *models.Produto) error {
	return r.db.Save(produto).Error
}

func (r *ProdutoRepository) Delete(id int) error {
	produto, err := r.FindByID(id)
	if err != nil {
		return err
	}
	if produto.IsDeleted() {
		return errors.New("produto já foi deletado")
	}
	produto.SoftDelete()
	return r.db.Save(produto).Error
}

func (r *ProdutoRepository) FindByID(id int) (*models.Produto, error) {
	var produto models.Produto
	err := r.db.
		Preload("produto_especie").
		Preload("produto_grupo").
		Preload("produto_unidade").
		Preload("produto_marca").
		Preload("produto_subgrupo"). // Corrected from preload to Preload
		Preload("produto_serie").
		First(&produto, id).Error
	if err != nil {
		return nil, err
	}
	return &produto, nil
}
