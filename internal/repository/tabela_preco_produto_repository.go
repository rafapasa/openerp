package repository

import (
	"errors"
	"fmt"

	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
	"gorm.io/gorm"
)

type TabelaPrecoProdutoRepository struct {
	db *gorm.DB
}

func NewTabelaPrecoProdutoRepository(db *gorm.DB) *TabelaPrecoProdutoRepository {
	return &TabelaPrecoProdutoRepository{db: db}
}

func (r *TabelaPrecoProdutoRepository) Create(createItem *models.TabelaPrecoProduto) error {
	var maxItem int
	err := r.db.Model(&models.TabelaPrecoProduto{}).
		Where("tbp_id = ?", createItem.TabelaPrecoID).
		Select("COALESCE(MAX(tbpp_item), 0) + 1").
		Scan(&maxItem).Error
	if err != nil {
		return err
	}
	createItem.Item = maxItem
	return r.db.Create(createItem).Error
}

func (r *TabelaPrecoProdutoRepository) Update(id, item int, updateItem *models.TabelaPrecoProduto) error {
	return r.db.
		Omit("TabelaPreco, Produto").
		Model(&models.TabelaPrecoProduto{}).
		Where("tbp_id = ? AND tbpp_item = ?", id, item).
		Updates(updateItem).
		Error
}

func (r *TabelaPrecoProdutoRepository) Delete(id, item int) error {
	deleteItem, err := r.FindByID(id, item)
	if err != nil {
		return err
	}
	if deleteItem.IsDeleted() {
		return fmt.Errorf("item %d já está deletado", id)
	}
	deleteItem.SoftDelete()
	return r.Update(id, item, deleteItem)
}

func (r *TabelaPrecoProdutoRepository) FindByID(id, item int) (*models.TabelaPrecoProduto, error) {
	var tabItem models.TabelaPrecoProduto
	err := r.db.Preload("Produto").
		Where("tbpp_item = ? AND tbp_id = ? AND deleted_at IS NULL", item, id).
		First(&tabItem).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item da tabela de preço não encontrado")
		}
		return nil, err
	}
	return &tabItem, nil
}

func (r *TabelaPrecoProdutoRepository) List(tabelaPrecoID, limit, offset int, filters map[string]interface{}) ([]models.TabelaPrecoProduto, int64, error) {
	var items []models.TabelaPrecoProduto
	var total int64

	query := r.db.
		Model(&models.TabelaPrecoProduto{}).
		Joins("LEFT JOIN produto ON produto.pro_id = tabela_preco_produto.pro_id").
		Where("tbp_id = ? AND deleted_at IS NULL", tabelaPrecoID)

	// Aplica filtros
	if nome, ok := filters["produto_nome"]; ok {
		query = query.Where("produto.pro_nome LIKE ?", "%"+nome.(string)+"%")
		delete(filters, "produto_nome")
	}
	query = utils.ApplyFilters(query, models.TabelaPrecoProduto{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("Produto").
		Limit(limit).
		Offset(offset).
		Order("tbpp_item DESC").
		Find(&items).
		Error
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *TabelaPrecoProdutoRepository) ExistsByProduto(tabelaPrecoID, produtoID, excludeItemID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.TabelaPrecoProduto{}).
		Where("tbp_id = ? AND pro_id = ? AND deleted_at IS NULL", tabelaPrecoID, produtoID)

	if excludeItemID > 0 {
		query = query.Where("tbpp_item != ?", excludeItemID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
