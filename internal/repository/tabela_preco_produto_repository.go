// internal/repository/tabela_preco_produto_repository.go
package repository

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// INTERFACE
// ============================================================

// TabelaPrecoProdutoRepository define o contrato para operações de banco
type TabelaPrecoProdutoRepository interface {
	// CRUD Básico
	Create(item *models.TabelaPrecoProduto) error
	Update(tabelaPrecoID, item int, produto *models.TabelaPrecoProduto) error
	Delete(tabelaPrecoID, item int) error
	FindByID(tabelaPrecoID, item int) (*models.TabelaPrecoProduto, error)
	GetByID(tabelaPrecoID, item int) (*models.TabelaPrecoProduto, error)

	// Buscas Específicas
	FindByTabelaPrecoID(tabelaPrecoID int) ([]models.TabelaPrecoProduto, error)
	FindByTabelaPrecoIDWithPagination(tabelaPrecoID int, limit, offset int) ([]models.TabelaPrecoProduto, int64, error)
	FindByProdutoID(produtoID int) ([]models.TabelaPrecoProduto, error)
	FindByTabelaPrecoAndProduto(tabelaPrecoID, produtoID int) (*models.TabelaPrecoProduto, error)
	FindActiveByTabelaPrecoID(tabelaPrecoID int) ([]models.TabelaPrecoProduto, error)
	FindBySituacao(situacao int) ([]models.TabelaPrecoProduto, error)

	// Listagem com Filtros
	List(tabelaPrecoID int, limit, offset int, filters map[string]interface{}) ([]models.TabelaPrecoProduto, int64, error)
	ListWithFullPreload(limit, offset int, filters map[string]interface{}) ([]models.TabelaPrecoProduto, int64, error)
	FindAll() ([]models.TabelaPrecoProduto, error)

	// Consultas de Validação (APENAS CONSULTAS)
	ExistsByTabelaPrecoAndProduto(tabelaPrecoID, produtoID int, excludeItem int) (bool, error)
	ExistsByID(tabelaPrecoID, item int) (bool, error)
	CountByTabelaPrecoID(tabelaPrecoID int) (int64, error)
	GetNextItemNumber(tabelaPrecoID int) (int, error)

	// Operações em Lote
	BulkUpdateStatus(tabelaPrecoID int, ids []int, situacao int) error
	BulkDelete(tabelaPrecoID int, ids []int) error
	CopyFromTabela(sourceTabelaID, targetTabelaID int) error

	// Consultas de Dependências
	HasDependentRecords(tabelaPrecoID, item int) (bool, error)
	CountDependentRecords(tabelaPrecoID, item int) (map[string]int64, error)
}

// ============================================================
// IMPLEMENTAÇÃO CONCRETA (privada)
// ============================================================

type tabelaPrecoProdutoRepository struct {
	db *gorm.DB
}

// NewTabelaPrecoProdutoRepository cria uma nova instância (retorna a interface)
func NewTabelaPrecoProdutoRepository(db *gorm.DB) TabelaPrecoProdutoRepository {
	return &tabelaPrecoProdutoRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD (APENAS PERSISTÊNCIA)
// ============================================================

// Create salva um novo item na tabela de preço
func (r *tabelaPrecoProdutoRepository) Create(item *models.TabelaPrecoProduto) error {
	err := r.db.Create(item).Error
	if err != nil {
		return apperrors.NewInternalError("Erro ao criar item na tabela de preço", err)
	}
	return nil
}

// Update atualiza um item da tabela de preço existente
func (r *tabelaPrecoProdutoRepository) Update(tabelaPrecoID, item int, produto *models.TabelaPrecoProduto) error {
	err := r.db.
		Omit("TabelaPreco", "Produto", "created_at", "deleted_at").
		Model(&models.TabelaPrecoProduto{}).
		Where("tbp_id = ? AND tbpp_item = ?", tabelaPrecoID, item).
		Updates(produto).Error
	if err != nil {
		return apperrors.NewInternalError("Erro ao atualizar item na tabela de preço", err)
	}
	return nil
}

// Delete realiza exclusão lógica
func (r *tabelaPrecoProdutoRepository) Delete(tabelaPrecoID, item int) error {
	err := r.db.
		Model(&models.TabelaPrecoProduto{}).
		Where("tbp_id = ? AND tbpp_item = ?", tabelaPrecoID, item).
		Update("deleted_at", gorm.Expr("NOW()")).Error
	if err != nil {
		return apperrors.NewInternalError("Erro ao excluir item da tabela de preço", err)
	}
	return nil
}

// FindByID busca um item da tabela de preço pelo ID composto com relacionamentos
func (r *tabelaPrecoProdutoRepository) FindByID(tabelaPrecoID, item int) (*models.TabelaPrecoProduto, error) {
	var produto models.TabelaPrecoProduto
	err := r.db.
		Preload("TabelaPreco").
		Preload("Produto").
		Where("tbp_id = ? AND tbpp_item = ? AND deleted_at IS NULL", tabelaPrecoID, item).
		First(&produto).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("item da tabela de preço com tabela %d e item %d não encontrado", tabelaPrecoID, item))
		}
		return nil, apperrors.NewInternalError("Erro buscando item da tabela de preço", err)
	}
	return &produto, nil
}

// GetByID busca um item da tabela de preço pelo ID composto (sem relacionamentos)
func (r *tabelaPrecoProdutoRepository) GetByID(tabelaPrecoID, item int) (*models.TabelaPrecoProduto, error) {
	var produto models.TabelaPrecoProduto
	err := r.db.
		Where("tbp_id = ? AND tbpp_item = ? AND deleted_at IS NULL", tabelaPrecoID, item).
		First(&produto).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("item da tabela de preço com tabela %d e item %d não encontrado", tabelaPrecoID, item))
		}
		return nil, apperrors.NewInternalError("Erro buscando item da tabela de preço", err)
	}
	return &produto, nil
}

// ============================================================
// MÉTODOS DE BUSCA ESPECÍFICOS
// ============================================================

// FindByTabelaPrecoID busca todos os itens de uma tabela de preço
func (r *tabelaPrecoProdutoRepository) FindByTabelaPrecoID(tabelaPrecoID int) ([]models.TabelaPrecoProduto, error) {
	var produtos []models.TabelaPrecoProduto
	err := r.db.
		Preload("Produto").
		Where("tbp_id = ? AND deleted_at IS NULL", tabelaPrecoID).
		Order("tbpp_item ASC").
		Find(&produtos).Error

	if err != nil {
		return nil, apperrors.NewInternalError("Erro buscando item da tabela de preço", err)
	}
	return produtos, nil
}

// FindByTabelaPrecoIDWithPagination busca itens de uma tabela de preço com paginação
func (r *tabelaPrecoProdutoRepository) FindByTabelaPrecoIDWithPagination(tabelaPrecoID int, limit, offset int) ([]models.TabelaPrecoProduto, int64, error) {
	var produtos []models.TabelaPrecoProduto
	var total int64

	query := r.db.Model(&models.TabelaPrecoProduto{}).
		Where("tbp_id = ? AND deleted_at IS NULL", tabelaPrecoID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("Erro buscando item da tabela de preço", err)
	}

	err := query.
		Preload("Produto").
		Limit(limit).
		Offset(offset).
		Order("tbpp_item ASC").
		Find(&produtos).Error

	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro buscando item da tabela de preço", err)
	}

	return produtos, total, nil
}

// FindByProdutoID busca itens de tabela de preço por produto
func (r *tabelaPrecoProdutoRepository) FindByProdutoID(produtoID int) ([]models.TabelaPrecoProduto, error) {
	var produtos []models.TabelaPrecoProduto
	err := r.db.
		Preload("TabelaPreco").
		Where("pro_id = ? AND deleted_at IS NULL", produtoID).
		Order("tbp_id ASC").
		Find(&produtos).Error

	if err != nil {
		return nil, apperrors.NewInternalError("Erro buscando item da tabela de preço", err)
	}
	return produtos, nil
}

// FindByTabelaPrecoAndProduto busca um item específico por tabela e produto
func (r *tabelaPrecoProdutoRepository) FindByTabelaPrecoAndProduto(tabelaPrecoID, produtoID int) (*models.TabelaPrecoProduto, error) {
	var produto models.TabelaPrecoProduto
	err := r.db.
		Preload("TabelaPreco").
		Preload("Produto").
		Where("tbp_id = ? AND pro_id = ? AND deleted_at IS NULL", tabelaPrecoID, produtoID).
		First(&produto).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("Produto %d não encontrado na tabela de preço %d", produtoID, tabelaPrecoID))
		}
		return nil, apperrors.NewInternalError("Erro buscando item da tabela de preço", err)
	}
	return &produto, nil
}

// FindActiveByTabelaPrecoID busca itens ativos de uma tabela de preço
func (r *tabelaPrecoProdutoRepository) FindActiveByTabelaPrecoID(tabelaPrecoID int) ([]models.TabelaPrecoProduto, error) {
	var produtos []models.TabelaPrecoProduto
	err := r.db.
		Preload("Produto").
		Where("tbp_id = ? AND tbpp_situacao = 1 AND deleted_at IS NULL", tabelaPrecoID).
		Order("tbpp_item ASC").
		Find(&produtos).Error

	if err != nil {
		return nil, apperrors.NewInternalError("Erro buscando item da tabela de preço", err)
	}
	return produtos, nil
}

// FindBySituacao busca itens por situação
func (r *tabelaPrecoProdutoRepository) FindBySituacao(situacao int) ([]models.TabelaPrecoProduto, error) {
	var produtos []models.TabelaPrecoProduto
	err := r.db.
		Preload("TabelaPreco").
		Preload("Produto").
		Where("tbpp_situacao = ? AND deleted_at IS NULL", situacao).
		Order("tbp_id ASC, tbpp_item ASC").
		Find(&produtos).Error

	if err != nil {
		return nil, apperrors.NewInternalError("Erro buscando item da tabela de preço", err)
	}
	return produtos, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de itens da tabela de preço com paginação e filtros
func (r *tabelaPrecoProdutoRepository) List(tabelaPrecoID, limit, offset int, filters map[string]interface{}) ([]models.TabelaPrecoProduto, int64, error) {
	var produtos []models.TabelaPrecoProduto
	var total int64

	query := r.db.Model(&models.TabelaPrecoProduto{}).
		Where("deleted_at IS NULL").
		Where("tbp_id = ?", tabelaPrecoID)
	query = utils.ApplyFilters(query, models.TabelaPrecoProduto{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("Erro buscando item da tabela de preço", err)
	}

	err := query.
		Preload("TabelaPreco").
		Preload("Produto").
		Limit(limit).
		Offset(offset).
		Order("tbp_id ASC, tbpp_item ASC").
		Find(&produtos).Error

	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro buscando item da tabela de preço", err)
	}

	return produtos, total, nil
}

// ListWithFullPreload retorna uma lista com todos os relacionamentos
func (r *tabelaPrecoProdutoRepository) ListWithFullPreload(limit, offset int, filters map[string]interface{}) ([]models.TabelaPrecoProduto, int64, error) {
	var produtos []models.TabelaPrecoProduto
	var total int64

	query := r.db.Model(&models.TabelaPrecoProduto{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.TabelaPrecoProduto{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("Erro buscando item da tabela de preço", err)
	}

	err := query.
		Preload("TabelaPreco").
		Preload("Produto").
		Preload("Produto.TipoProduto").
		Preload("Produto.ProdutoGrupo").
		Preload("Produto.Marca").
		Preload("Produto.Modelo").
		Limit(limit).
		Offset(offset).
		Order("tbp_id ASC, tbpp_item ASC").
		Find(&produtos).Error

	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro buscando item da tabela de preço", err)
	}

	return produtos, total, nil
}

// FindAll busca todos os itens da tabela de preço
func (r *tabelaPrecoProdutoRepository) FindAll() ([]models.TabelaPrecoProduto, error) {
	var produtos []models.TabelaPrecoProduto
	err := r.db.
		Where("deleted_at IS NULL").
		Order("tbp_id ASC, tbpp_item ASC").
		Find(&produtos).Error
	if err != nil {
		return nil, apperrors.NewInternalError("Erro buscando item da tabela de preço", err)
	}
	return produtos, nil
}

// ============================================================
// MÉTODOS DE CONSULTA PARA VALIDAÇÕES (APENAS CONSULTAS)
// ============================================================

// ExistsByTabelaPrecoAndProduto verifica se já existe um item com o produto na tabela
func (r *tabelaPrecoProdutoRepository) ExistsByTabelaPrecoAndProduto(tabelaPrecoID, produtoID int, excludeItem int) (bool, error) {
	var count int64
	query := r.db.Model(&models.TabelaPrecoProduto{}).
		Where("tbp_id = ? AND pro_id = ? AND deleted_at IS NULL", tabelaPrecoID, produtoID)

	if excludeItem > 0 {
		query = query.Where("tbpp_item != ?", excludeItem)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, apperrors.NewInternalError("Erro buscando item da tabela de preço", err)
	}
	return count > 0, nil
}

// ExistsByID verifica se um item existe pelo ID composto
func (r *tabelaPrecoProdutoRepository) ExistsByID(tabelaPrecoID, item int) (bool, error) {
	var count int64
	err := r.db.Model(&models.TabelaPrecoProduto{}).
		Where("tbp_id = ? AND tbpp_item = ? AND deleted_at IS NULL", tabelaPrecoID, item).
		Count(&count).Error
	if err != nil {
		return false, apperrors.NewInternalError("Erro buscando item da tabela de preço", err)
	}
	return count > 0, nil
}

// CountByTabelaPrecoID retorna a quantidade de itens em uma tabela de preço
func (r *tabelaPrecoProdutoRepository) CountByTabelaPrecoID(tabelaPrecoID int) (int64, error) {
	var count int64
	err := r.db.Model(&models.TabelaPrecoProduto{}).
		Where("tbp_id = ? AND deleted_at IS NULL", tabelaPrecoID).
		Count(&count).Error
	if err != nil {
		return 0, apperrors.NewInternalError("Erro buscando item da tabela de preço", err)
	}
	return count, nil
}

// GetNextItemNumber retorna o próximo número de item para uma tabela de preço
func (r *tabelaPrecoProdutoRepository) GetNextItemNumber(tabelaPrecoID int) (int, error) {
	var maxItem int
	err := r.db.Model(&models.TabelaPrecoProduto{}).
		Where("tbp_id = ?", tabelaPrecoID).
		Select("COALESCE(MAX(tbpp_item), 0) + 1").
		Scan(&maxItem).Error
	if err != nil {
		return 0, apperrors.NewInternalError("Erro buscando item da tabela de preço", err)
	}
	return maxItem, nil
}

// ============================================================
// OPERAÇÕES EM LOTE
// ============================================================

// BulkUpdateStatus atualiza a situação de múltiplos itens
func (r *tabelaPrecoProdutoRepository) BulkUpdateStatus(tabelaPrecoID int, ids []int, situacao int) error {
	return r.db.Model(&models.TabelaPrecoProduto{}).
		Where("tbp_id = ? AND tbpp_item IN ? AND deleted_at IS NULL", tabelaPrecoID, ids).
		Update("tbpp_situacao", situacao).Error
}

// BulkDelete realiza exclusão lógica de múltiplos itens
func (r *tabelaPrecoProdutoRepository) BulkDelete(tabelaPrecoID int, ids []int) error {
	return r.db.Model(&models.TabelaPrecoProduto{}).
		Where("tbp_id = ? AND tbpp_item IN ? AND deleted_at IS NULL", tabelaPrecoID, ids).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// CopyFromTabela copia todos os itens de uma tabela de preço para outra
func (r *tabelaPrecoProdutoRepository) CopyFromTabela(sourceTabelaID, targetTabelaID int) error {
	// Busca itens da tabela fonte
	var itens []models.TabelaPrecoProduto
	err := r.db.
		Where("tbp_id = ? AND deleted_at IS NULL", sourceTabelaID).
		Find(&itens).Error
	if err != nil {
		return err
	}

	if len(itens) == 0 {
		return nil // Nada para copiar
	}

	// Remove itens existentes na tabela destino
	if err := r.db.
		Where("tbp_id = ?", targetTabelaID).
		Delete(&models.TabelaPrecoProduto{}).Error; err != nil {
		return err
	}

	// Prepara para cópia
	for i := range itens {
		itens[i].TabelaPrecoID = targetTabelaID
		itens[i].Item = 0 // Será gerado automaticamente
		itens[i].CreatedAt = time.Now()
		itens[i].UpdatedAt = time.Now()
		itens[i].DeletedAt = nil
	}

	// Insere em lote
	return r.db.Create(&itens).Error
}

// ============================================================
// MÉTODOS DE CONSULTA DE DEPENDÊNCIAS
// ============================================================

// HasDependentRecords verifica se o item tem registros dependentes
func (r *tabelaPrecoProdutoRepository) HasDependentRecords(tabelaPrecoID, item int) (bool, error) {
	counts, err := r.CountDependentRecords(tabelaPrecoID, item)
	if err != nil {
		return false, err
	}

	for _, count := range counts {
		if count > 0 {
			return true, nil
		}
	}
	return false, nil
}

// CountDependentRecords retorna a contagem de registros dependentes por tipo
func (r *tabelaPrecoProdutoRepository) CountDependentRecords(tabelaPrecoID, item int) (map[string]int64, error) {
	result := make(map[string]int64)

	// Verifica documentos que usam este preço
	var countDocumentos int64
	if err := r.db.Model(&models.DocumentoVendaItem{}).
		Where("tbp_id = ? AND tbpp_item = ? AND deleted_at IS NULL", tabelaPrecoID, item).
		Count(&countDocumentos).Error; err != nil {
		return nil, err
	}
	if countDocumentos > 0 {
		result["documentos"] = countDocumentos
	}

	return result, nil
}
