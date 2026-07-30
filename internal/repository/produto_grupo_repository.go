// internal/repository/produto_grupo_repository.go
package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// INTERFACE
// ============================================================

// ProdutoGrupoRepository define o contrato para operações de banco
type ProdutoGrupoRepository interface {
	// CRUD Básico
	Create(grupo *models.ProdutoGrupo) error
	Update(id int, grupo *models.ProdutoGrupo) error
	Delete(id int) error
	FindByID(id int) (*models.ProdutoGrupo, error)
	GetByID(id int) (*models.ProdutoGrupo, error)

	// Buscas Específicas
	FindByDescricao(descricao string, limit int) ([]models.ProdutoGrupo, error)
	FindBySituacao(situacao int) ([]models.ProdutoGrupo, error)
	FindActive() ([]models.ProdutoGrupo, error)
	FindVisivelNoCaixa() ([]models.ProdutoGrupo, error)

	// Listagem com Filtros
	List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoGrupo, int64, error)
	ListWithProdutos(limit, offset int, filters map[string]interface{}) ([]models.ProdutoGrupo, int64, error)
	FindAll() ([]models.ProdutoGrupo, error)

	// Consultas de Validação (APENAS CONSULTAS)
	ExistsByDescricao(descricao string, excludeID int) (bool, error)
	ExistsByID(id int) (bool, error)
	Count() (int64, error)
	CountProdutosByGrupo(grupoID int) (int64, error)

	// Operações em Lote
	BulkUpdateStatus(ids []int, situacao int) error
	BulkDelete(ids []int) error

	// Consultas de Dependências
	HasDependentRecords(id int) (bool, error)
	CountDependentRecords(id int) (map[string]int64, error)
}

// ============================================================
// IMPLEMENTAÇÃO CONCRETA (privada)
// ============================================================

type produtoGrupoRepository struct {
	db *gorm.DB
}

// NewProdutoGrupoRepository cria uma nova instância (retorna a interface)
func NewProdutoGrupoRepository(db *gorm.DB) ProdutoGrupoRepository {
	return &produtoGrupoRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD (APENAS PERSISTÊNCIA)
// ============================================================

// Create salva um novo grupo de produto
func (r *produtoGrupoRepository) Create(grupo *models.ProdutoGrupo) error {
	return r.db.Create(grupo).Error
}

// Update atualiza um grupo de produto existente
func (r *produtoGrupoRepository) Update(id int, grupo *models.ProdutoGrupo) error {
	return r.db.
		Omit("Produtos", "created_at", "deleted_at").
		Model(&models.ProdutoGrupo{}).
		Where("prog_id = ?", id).
		Updates(grupo).Error
}

// Delete realiza exclusão lógica
func (r *produtoGrupoRepository) Delete(id int) error {
	return r.db.
		Model(&models.ProdutoGrupo{}).
		Where("prog_id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// FindByID busca um grupo de produto pelo ID com relacionamentos
func (r *produtoGrupoRepository) FindByID(id int) (*models.ProdutoGrupo, error) {
	var grupo models.ProdutoGrupo
	err := r.db.
		Preload("Produtos", "deleted_at IS NULL").
		Where("prog_id = ? AND deleted_at IS NULL", id).
		First(&grupo).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("grupo de produto com ID %d não encontrado", id))
		}
		return nil, err
	}
	return &grupo, nil
}

// GetByID busca um grupo de produto pelo ID (sem relacionamentos)
func (r *produtoGrupoRepository) GetByID(id int) (*models.ProdutoGrupo, error) {
	var grupo models.ProdutoGrupo
	err := r.db.
		Where("prog_id = ? AND deleted_at IS NULL", id).
		First(&grupo).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("grupo de produto com ID %d não encontrado", id))
		}
		return nil, err
	}
	return &grupo, nil
}

// ============================================================
// MÉTODOS DE BUSCA ESPECÍFICOS
// ============================================================

// FindByDescricao busca grupos de produto pela descrição (autocomplete)
func (r *produtoGrupoRepository) FindByDescricao(descricao string, limit int) ([]models.ProdutoGrupo, error) {
	var grupos []models.ProdutoGrupo
	err := r.db.
		Where("prog_descricao LIKE ? AND deleted_at IS NULL", "%"+descricao+"%").
		Order("prog_descricao ASC").
		Limit(limit).
		Find(&grupos).Error

	if err != nil {
		return nil, err
	}
	return grupos, nil
}

// FindBySituacao busca grupos de produto por situação
func (r *produtoGrupoRepository) FindBySituacao(situacao int) ([]models.ProdutoGrupo, error) {
	var grupos []models.ProdutoGrupo
	err := r.db.
		Where("prog_situacao = ? AND deleted_at IS NULL", situacao).
		Order("prog_descricao ASC").
		Find(&grupos).Error

	if err != nil {
		return nil, err
	}
	return grupos, nil
}

// FindActive busca todos os grupos de produto ativos
func (r *produtoGrupoRepository) FindActive() ([]models.ProdutoGrupo, error) {
	var grupos []models.ProdutoGrupo
	err := r.db.
		Where("prog_situacao = 1 AND deleted_at IS NULL").
		Order("prog_descricao ASC").
		Find(&grupos).Error

	if err != nil {
		return nil, err
	}
	return grupos, nil
}

// FindVisivelNoCaixa busca grupos de produto visíveis no frente de caixa
func (r *produtoGrupoRepository) FindVisivelNoCaixa() ([]models.ProdutoGrupo, error) {
	var grupos []models.ProdutoGrupo
	err := r.db.
		Where("prog_visivelfrentecaixa = 1 AND deleted_at IS NULL").
		Order("prog_descricao ASC").
		Find(&grupos).Error

	if err != nil {
		return nil, err
	}
	return grupos, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de grupos de produto com paginação e filtros
func (r *produtoGrupoRepository) List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoGrupo, int64, error) {
	var grupos []models.ProdutoGrupo
	var total int64

	query := r.db.Model(&models.ProdutoGrupo{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.ProdutoGrupo{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Order("prog_descricao ASC").
		Find(&grupos).Error

	if err != nil {
		return nil, 0, err
	}

	return grupos, total, nil
}

// ListWithProdutos retorna grupos de produto com contagem de produtos
func (r *produtoGrupoRepository) ListWithProdutos(limit, offset int, filters map[string]interface{}) ([]models.ProdutoGrupo, int64, error) {
	var grupos []models.ProdutoGrupo
	var total int64

	query := r.db.Model(&models.ProdutoGrupo{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.ProdutoGrupo{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Select(`
			prog.*,
			(
				SELECT COUNT(*) 
				FROM produto 
				WHERE prog_id = produto.prog_id 
				AND deleted_at IS NULL
			) as total_produtos
		`).
		Limit(limit).
		Offset(offset).
		Order("prog_descricao ASC").
		Find(&grupos).Error

	if err != nil {
		return nil, 0, err
	}

	return grupos, total, nil
}

// FindAll busca todos os grupos de produto
func (r *produtoGrupoRepository) FindAll() ([]models.ProdutoGrupo, error) {
	var grupos []models.ProdutoGrupo
	err := r.db.
		Where("deleted_at IS NULL").
		Order("prog_descricao ASC").
		Find(&grupos).Error
	if err != nil {
		return nil, err
	}
	return grupos, nil
}

// ============================================================
// MÉTODOS DE CONSULTA PARA VALIDAÇÕES (APENAS CONSULTAS)
// ============================================================

// ExistsByDescricao verifica se já existe um grupo com a descrição
func (r *produtoGrupoRepository) ExistsByDescricao(descricao string, excludeID int) (bool, error) {
	if descricao == "" {
		return false, nil
	}

	var count int64
	query := r.db.Model(&models.ProdutoGrupo{}).
		Where("prog_descricao = ? AND deleted_at IS NULL", descricao)

	if excludeID > 0 {
		query = query.Where("prog_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByID verifica se um grupo de produto existe pelo ID
func (r *produtoGrupoRepository) ExistsByID(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.ProdutoGrupo{}).
		Where("prog_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Count retorna o total de grupos de produto
func (r *produtoGrupoRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.ProdutoGrupo{}).
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CountProdutosByGrupo retorna a quantidade de produtos em um grupo
func (r *produtoGrupoRepository) CountProdutosByGrupo(grupoID int) (int64, error) {
	var count int64
	err := r.db.Model(&models.Produto{}).
		Where("prog_id = ? AND deleted_at IS NULL", grupoID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// ============================================================
// OPERAÇÕES EM LOTE
// ============================================================

// BulkUpdateStatus atualiza a situação de múltiplos grupos
func (r *produtoGrupoRepository) BulkUpdateStatus(ids []int, situacao int) error {
	return r.db.Model(&models.ProdutoGrupo{}).
		Where("prog_id IN ? AND deleted_at IS NULL", ids).
		Update("prog_situacao", situacao).Error
}

// BulkDelete realiza exclusão lógica de múltiplos grupos
func (r *produtoGrupoRepository) BulkDelete(ids []int) error {
	return r.db.Model(&models.ProdutoGrupo{}).
		Where("prog_id IN ? AND deleted_at IS NULL", ids).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// ============================================================
// MÉTODOS DE CONSULTA DE DEPENDÊNCIAS
// ============================================================

// HasDependentRecords verifica se o grupo tem registros dependentes
func (r *produtoGrupoRepository) HasDependentRecords(id int) (bool, error) {
	counts, err := r.CountDependentRecords(id)
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
func (r *produtoGrupoRepository) CountDependentRecords(id int) (map[string]int64, error) {
	result := make(map[string]int64)

	// Verifica produtos associados
	var countProdutos int64
	if err := r.db.Model(&models.Produto{}).
		Where("prog_id = ? AND deleted_at IS NULL", id).
		Count(&countProdutos).Error; err != nil {
		return nil, err
	}
	if countProdutos > 0 {
		result["produtos"] = countProdutos
	}

	return result, nil
}