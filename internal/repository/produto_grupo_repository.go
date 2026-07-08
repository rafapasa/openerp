package repository

import (
	"fmt"

	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
	"gorm.io/gorm"
)

type ProdutoGrupoRepository struct {
	db *gorm.DB
}

func NewProdutoGrupoRepository(db *gorm.DB) *ProdutoGrupoRepository {
	return &ProdutoGrupoRepository{db: db}
}

// FindByID busca um grupo de produto pelo ID
func (r *ProdutoGrupoRepository) FindByID(id int) (*models.ProdutoGrupo, error) {
	var produtoGrupo models.ProdutoGrupo
	err := r.db.Where("id = ?", id).First(&produtoGrupo).Error
	if err != nil {
		return nil, err
	}
	return &produtoGrupo, nil
}

func (r *ProdutoGrupoRepository) List(limit, offset int, filters map[string]any) ([]models.ProdutoGrupo, int64, error) {
	var produtoGrupos []models.ProdutoGrupo
	var total int64
	query := r.db

	// Aplicar filtros
	query.Where("deleted_at IS NULL")
	utils.ApplyFilters(query, models.ProdutoGrupo{}, filters)

	err := query.Limit(limit).Offset(offset).
		Order("id desc").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Limit(limit).Offset(offset).
		Order("id desc").
		Find(&produtoGrupos).Error
	if err != nil {
		return nil, 0, err
	}

	return produtoGrupos, total, nil
}

func (r *ProdutoGrupoRepository) Create(produtoGrupo *models.ProdutoGrupo) error {
	return r.db.Create(produtoGrupo).Error
}

func (r *ProdutoGrupoRepository) Update(id int, produtoGrupo *models.ProdutoGrupo) error {
	return r.db.
		Omit("Produto").
		Where("prog_id = ?", id).
		Updates(produtoGrupo).Error
}

func (r *ProdutoGrupoRepository) Delete(id int) error {
	produtoGrupo, err := r.FindByID(id)
	if err != nil {
		return err
	}
	if produtoGrupo.IsDeleted() {
		return fmt.Errorf("Grupo de produto %s já está deletado", produtoGrupo.Descricao)
	}

	produtoGrupo.SoftDelete()
	return r.Update(id, produtoGrupo)
}
