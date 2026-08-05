// internal/repository/produto_variacao_repository.go
package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// INTERFACE
// ============================================================

// ProdutoVariacaoRepository define o contrato para operações de banco
type ProdutoVariacaoRepository interface {
	// CRUD Básico
	Create(variacao *models.ProdutoVariacao) error
	Update(id int, variacao *models.ProdutoVariacao) error
	Delete(id int) error
	FindByID(id int) (*models.ProdutoVariacao, error)
	GetByID(id int) (*models.ProdutoVariacao, error)

	// Buscas Específicas
	FindBySKU(sku string) (*models.ProdutoVariacao, error)
	FindBySKUAndFilial(sku string, filialID int) (*models.ProdutoVariacao, error)
	FindByProdutoID(produtoID int) ([]models.ProdutoVariacao, error)
	FindByProdutoIDAndFilial(produtoID int, filialID int) ([]models.ProdutoVariacao, error)
	FindByCorID(corID int) ([]models.ProdutoVariacao, error)
	FindByTamanhoID(tamanhoID int) ([]models.ProdutoVariacao, error)
	FindByProdutoCorTamanho(produtoID int, corID, tamanhoID *int) (*models.ProdutoVariacao, error)
	FindByFilialID(filialID int) ([]models.ProdutoVariacao, error)

	// Listagem com Filtros
	List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoVariacao, int64, error)
	ListWithFullPreload(limit, offset int, filters map[string]interface{}) ([]models.ProdutoVariacao, int64, error)
	FindAll() ([]models.ProdutoVariacao, error)

	// Consultas de Validação (APENAS CONSULTAS)
	ExistsBySKU(sku string, excludeID int) (bool, error)
	ExistsBySKUAndFilial(sku string, filialID int, excludeID int) (bool, error)
	ExistsByProdutoCorTamanho(produtoID int, corID, tamanhoID *int, excludeID int) (bool, error)
	ExistsByID(id int) (bool, error)
	CountByProdutoID(produtoID int) (int64, error)
	CountByFilialID(filialID int) (int64, error)

	// Operações em Lote
	BulkUpdateEstoque(ids []int, quantidade float64) error
	BulkDelete(ids []int) error
}

// ============================================================
// IMPLEMENTAÇÃO CONCRETA (privada)
// ============================================================

type produtoVariacaoRepository struct {
	db *gorm.DB
}

// NewProdutoVariacaoRepository cria uma nova instância (retorna a interface)
func NewProdutoVariacaoRepository(db *gorm.DB) ProdutoVariacaoRepository {
	return &produtoVariacaoRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD (APENAS PERSISTÊNCIA)
// ============================================================

// Create salva uma nova variação de produto
func (r *produtoVariacaoRepository) Create(variacao *models.ProdutoVariacao) error {
	return r.db.Create(variacao).Error
}

// Update atualiza uma variação de produto existente
func (r *produtoVariacaoRepository) Update(id int, variacao *models.ProdutoVariacao) error {
	return r.db.
		Omit("Produto", "EmpresaFilial", "Cor", "Tamanho", "created_at", "deleted_at").
		Model(&models.ProdutoVariacao{}).
		Where("provar_id = ?", id).
		Updates(variacao).Error
}

// Delete realiza exclusão lógica
func (r *produtoVariacaoRepository) Delete(id int) error {
	return r.db.
		Model(&models.ProdutoVariacao{}).
		Where("provar_id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// FindByID busca uma variação pelo ID com relacionamentos
func (r *produtoVariacaoRepository) FindByID(id int) (*models.ProdutoVariacao, error) {
	var variacao models.ProdutoVariacao
	err := r.db.
		Preload("Produto").
		Preload("EmpresaFilial").
		Preload("Cor").
		Preload("Tamanho").
		Where("provar_id = ? AND deleted_at IS NULL", id).
		First(&variacao).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("variação de produto com ID %d não encontrada", id))
		}
		return nil, err
	}
	return &variacao, nil
}

// GetByID busca uma variação pelo ID (sem relacionamentos)
func (r *produtoVariacaoRepository) GetByID(id int) (*models.ProdutoVariacao, error) {
	var variacao models.ProdutoVariacao
	err := r.db.
		Where("provar_id = ? AND deleted_at IS NULL", id).
		First(&variacao).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("variação de produto com ID %d não encontrada", id))
		}
		return nil, err
	}
	return &variacao, nil
}

// ============================================================
// MÉTODOS DE BUSCA ESPECÍFICOS
// ============================================================

// FindBySKU busca uma variação pelo SKU
func (r *produtoVariacaoRepository) FindBySKU(sku string) (*models.ProdutoVariacao, error) {
	var variacao models.ProdutoVariacao
	err := r.db.
		Preload("Produto").
		Preload("Cor").
		Preload("Tamanho").
		Where("provar_sku = ? AND deleted_at IS NULL", sku).
		First(&variacao).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("variação com SKU %s não encontrada", sku))
		}
		return nil, err
	}
	return &variacao, nil
}

// FindBySKUAndFilial busca uma variação pelo SKU e filial
func (r *produtoVariacaoRepository) FindBySKUAndFilial(sku string, filialID int) (*models.ProdutoVariacao, error) {
	var variacao models.ProdutoVariacao
	err := r.db.
		Preload("Produto").
		Preload("Cor").
		Preload("Tamanho").
		Where("provar_sku = ? AND emf_id = ? AND deleted_at IS NULL", sku, filialID).
		First(&variacao).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("variação com SKU %s e filial %d não encontrada", sku, filialID))
		}
		return nil, err
	}
	return &variacao, nil
}

// FindByProdutoID busca variações de um produto
func (r *produtoVariacaoRepository) FindByProdutoID(produtoID int) ([]models.ProdutoVariacao, error) {
	var variacoes []models.ProdutoVariacao
	err := r.db.
		Preload("Cor").
		Preload("Tamanho").
		Where("pro_id = ? AND deleted_at IS NULL", produtoID).
		Order("provar_sku ASC").
		Find(&variacoes).Error

	if err != nil {
		return nil, err
	}
	return variacoes, nil
}

// FindByProdutoIDAndFilial busca variações de um produto por filial
func (r *produtoVariacaoRepository) FindByProdutoIDAndFilial(produtoID int, filialID int) ([]models.ProdutoVariacao, error) {
	var variacoes []models.ProdutoVariacao
	err := r.db.
		Preload("Cor").
		Preload("Tamanho").
		Where("pro_id = ? AND emf_id = ? AND deleted_at IS NULL", produtoID, filialID).
		Order("provar_sku ASC").
		Find(&variacoes).Error

	if err != nil {
		return nil, err
	}
	return variacoes, nil
}

// FindByCorID busca variações por cor
func (r *produtoVariacaoRepository) FindByCorID(corID int) ([]models.ProdutoVariacao, error) {
	var variacoes []models.ProdutoVariacao
	err := r.db.
		Preload("Produto").
		Preload("Tamanho").
		Where("cor_id = ? AND deleted_at IS NULL", corID).
		Order("provar_sku ASC").
		Find(&variacoes).Error

	if err != nil {
		return nil, err
	}
	return variacoes, nil
}

// FindByTamanhoID busca variações por tamanho
func (r *produtoVariacaoRepository) FindByTamanhoID(tamanhoID int) ([]models.ProdutoVariacao, error) {
	var variacoes []models.ProdutoVariacao
	err := r.db.
		Preload("Produto").
		Preload("Cor").
		Where("ptam_id = ? AND deleted_at IS NULL", tamanhoID).
		Order("provar_sku ASC").
		Find(&variacoes).Error

	if err != nil {
		return nil, err
	}
	return variacoes, nil
}

// FindByProdutoCorTamanho busca uma variação específica por produto, cor e tamanho
func (r *produtoVariacaoRepository) FindByProdutoCorTamanho(produtoID int, corID, tamanhoID *int) (*models.ProdutoVariacao, error) {
	var variacao models.ProdutoVariacao
	query := r.db.
		Where("pro_id = ? AND deleted_at IS NULL", produtoID)

	if corID != nil {
		query = query.Where("cor_id = ?", *corID)
	} else {
		query = query.Where("cor_id IS NULL")
	}

	if tamanhoID != nil {
		query = query.Where("ptam_id = ?", *tamanhoID)
	} else {
		query = query.Where("ptam_id IS NULL")
	}

	err := query.
		Preload("Cor").
		Preload("Tamanho").
		First(&variacao).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Não encontrado, não é erro
		}
		return nil, err
	}
	return &variacao, nil
}

// FindByFilialID busca variações por filial
func (r *produtoVariacaoRepository) FindByFilialID(filialID int) ([]models.ProdutoVariacao, error) {
	var variacoes []models.ProdutoVariacao
	err := r.db.
		Preload("Produto").
		Preload("Cor").
		Preload("Tamanho").
		Where("emf_id = ? AND deleted_at IS NULL", filialID).
		Order("provar_sku ASC").
		Find(&variacoes).Error

	if err != nil {
		return nil, err
	}
	return variacoes, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de variações com paginação e filtros
func (r *produtoVariacaoRepository) List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoVariacao, int64, error) {
	var variacoes []models.ProdutoVariacao
	var total int64

	query := r.db.Model(&models.ProdutoVariacao{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.ProdutoVariacao{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("Produto").
		Preload("Cor").
		Preload("Tamanho").
		Limit(limit).
		Offset(offset).
		Order("provar_sku ASC").
		Find(&variacoes).Error

	if err != nil {
		return nil, 0, err
	}

	return variacoes, total, nil
}

// ListWithFullPreload retorna uma lista com todos os relacionamentos
func (r *produtoVariacaoRepository) ListWithFullPreload(limit, offset int, filters map[string]interface{}) ([]models.ProdutoVariacao, int64, error) {
	var variacoes []models.ProdutoVariacao
	var total int64

	query := r.db.Model(&models.ProdutoVariacao{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.ProdutoVariacao{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("Produto").
		Preload("EmpresaFilial").
		Preload("Cor").
		Preload("Tamanho").
		Limit(limit).
		Offset(offset).
		Order("provar_sku ASC").
		Find(&variacoes).Error

	if err != nil {
		return nil, 0, err
	}

	return variacoes, total, nil
}

// FindAll busca todas as variações
func (r *produtoVariacaoRepository) FindAll() ([]models.ProdutoVariacao, error) {
	var variacoes []models.ProdutoVariacao
	err := r.db.
		Where("deleted_at IS NULL").
		Order("provar_sku ASC").
		Find(&variacoes).Error
	if err != nil {
		return nil, err
	}
	return variacoes, nil
}

// ============================================================
// MÉTODOS DE CONSULTA PARA VALIDAÇÕES (APENAS CONSULTAS)
// ============================================================

// ExistsBySKU verifica se já existe uma variação com o SKU
func (r *produtoVariacaoRepository) ExistsBySKU(sku string, excludeID int) (bool, error) {
	if sku == "" {
		return false, nil
	}

	var count int64
	query := r.db.Model(&models.ProdutoVariacao{}).
		Where("provar_sku = ? AND deleted_at IS NULL", sku)

	if excludeID > 0 {
		query = query.Where("provar_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsBySKUAndFilial verifica se já existe uma variação com o SKU e filial
func (r *produtoVariacaoRepository) ExistsBySKUAndFilial(sku string, filialID int, excludeID int) (bool, error) {
	if sku == "" {
		return false, nil
	}
	if filialID <= 0 {
		return false, errors.New("ID da filial inválido")
	}

	var count int64
	query := r.db.Model(&models.ProdutoVariacao{}).
		Where("provar_sku = ? AND emf_id = ? AND deleted_at IS NULL", sku, filialID)

	if excludeID > 0 {
		query = query.Where("provar_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByProdutoCorTamanho verifica se já existe uma variação com a combinação de produto, cor e tamanho
func (r *produtoVariacaoRepository) ExistsByProdutoCorTamanho(produtoID int, corID, tamanhoID *int, excludeID int) (bool, error) {
	if produtoID <= 0 {
		return false, errors.New("ID do produto inválido")
	}

	var count int64
	query := r.db.Model(&models.ProdutoVariacao{}).
		Where("pro_id = ? AND deleted_at IS NULL", produtoID)

	if corID != nil {
		query = query.Where("cor_id = ?", *corID)
	} else {
		query = query.Where("cor_id IS NULL")
	}

	if tamanhoID != nil {
		query = query.Where("ptam_id = ?", *tamanhoID)
	} else {
		query = query.Where("ptam_id IS NULL")
	}

	if excludeID > 0 {
		query = query.Where("provar_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByID verifica se uma variação existe pelo ID
func (r *produtoVariacaoRepository) ExistsByID(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.ProdutoVariacao{}).
		Where("provar_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountByProdutoID retorna a quantidade de variações de um produto
func (r *produtoVariacaoRepository) CountByProdutoID(produtoID int) (int64, error) {
	var count int64
	err := r.db.Model(&models.ProdutoVariacao{}).
		Where("pro_id = ? AND deleted_at IS NULL", produtoID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CountByFilialID retorna a quantidade de variações de uma filial
func (r *produtoVariacaoRepository) CountByFilialID(filialID int) (int64, error) {
	var count int64
	err := r.db.Model(&models.ProdutoVariacao{}).
		Where("emf_id = ? AND deleted_at IS NULL", filialID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// ============================================================
// OPERAÇÕES EM LOTE
// ============================================================

// BulkUpdateEstoque atualiza o estoque de múltiplas variações
func (r *produtoVariacaoRepository) BulkUpdateEstoque(ids []int, quantidade float64) error {
	return r.db.Model(&models.ProdutoVariacao{}).
		Where("provar_id IN ? AND deleted_at IS NULL", ids).
		Update("provar_estoque_atual", gorm.Expr("provar_estoque_atual + ?", quantidade)).Error
}

// BulkDelete realiza exclusão lógica de múltiplas variações
func (r *produtoVariacaoRepository) BulkDelete(ids []int) error {
	return r.db.Model(&models.ProdutoVariacao{}).
		Where("provar_id IN ? AND deleted_at IS NULL", ids).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}
