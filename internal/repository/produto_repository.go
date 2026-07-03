package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// TYPES
// ============================================================

// ProdutoRepository é o repositório para Produto
type ProdutoRepository struct {
	db *gorm.DB
}

// ============================================================
// CONSTRUCTOR
// ============================================================

// NewProdutoRepository cria uma nova instância
func NewProdutoRepository(db *gorm.DB) *ProdutoRepository {
	return &ProdutoRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD
// ============================================================

// Create salva um novo produto
func (r *ProdutoRepository) Create(produto *models.Produto) error {
	return r.db.Create(produto).Error
}

// Update atualiza um produto existente
func (r *ProdutoRepository) Update(produto *models.Produto) error {
	return r.db.Save(produto).Error
}

// Delete realiza exclusão lógica de um produto pelo ID
func (r *ProdutoRepository) Delete(id int) error {
	// 1. Buscar o produto
	produto, err := r.FindByID(id)
	if err != nil {
		return err
	}

	// 2. Verificar se já foi deletado
	if produto.IsDeleted() {
		return errors.New("produto já foi deletado")
	}

	// 3. Realizar soft delete
	produto.SoftDelete()

	// 4. Salvar
	return r.db.Save(produto).Error
}

// ============================================================
// MÉTODOS DE BUSCA
// ============================================================

// FindByID busca um produto pelo ID com relacionamentos
func (r *ProdutoRepository) FindByID(id int) (*models.Produto, error) {
	var produto models.Produto
	err := r.db.
		Preload("ProdutoEspecie").
		Preload("ProdutoGrupo").
		Preload("ProdutoMarca").
		Preload("ProdutoSubgrupo").
		Preload("ProdutoSerie").
		Preload("ProdutoTipoProduto").
		Where("pro_id = ? AND deleted_at IS NULL", id).
		First(&produto).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("produto não encontrado")
		}
		return nil, err
	}
	return &produto, nil
}

// FindByNome busca produtos pelo nome (autocomplete)
func (r *ProdutoRepository) FindByNome(nome string, limit int) ([]models.Produto, error) {
	var produtos []models.Produto
	err := r.db.
		Preload("ProdutoEspecie").
		Preload("ProdutoGrupo").
		Preload("ProdutoMarca").
		Preload("ProdutoSubgrupo").
		Preload("ProdutoSerie").
		Preload("ProdutoTipoProduto").
		Where("nome LIKE ? AND deleted_at IS NULL", "%"+nome+"%").
		Limit(limit).
		Order("pro_id DESC").
		Find(&produtos).Error

	if err != nil {
		return nil, err
	}
	return produtos, nil
}

// FindByCodigo busca um produto pelo código
func (r *ProdutoRepository) FindByCodigo(codigo string) (*models.Produto, error) {
	var produto models.Produto
	err := r.db.
		Preload("ProdutoEspecie").
		Preload("ProdutoGrupo").
		Preload("ProdutoMarca").
		Preload("ProdutoSubgrupo").
		Preload("ProdutoSerie").
		Preload("ProdutoTipoProduto").
		Where("pro_codigo = ? AND deleted_at IS NULL", codigo).
		First(&produto).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("produto não encontrado")
		}
		return nil, err
	}
	return &produto, nil
}

// FindByReferencia busca um produto pela referência
func (r *ProdutoRepository) FindByReferencia(referencia string) (*models.Produto, error) {
	var produto models.Produto
	err := r.db.
		Preload("ProdutoEspecie").
		Preload("ProdutoGrupo").
		Preload("ProdutoMarca").
		Preload("ProdutoSubgrupo").
		Preload("ProdutoSerie").
		Preload("ProdutoTipoProduto").
		Where("pro_referencia = ? AND deleted_at IS NULL", referencia).
		First(&produto).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("produto não encontrado")
		}
		return nil, err
	}
	return &produto, nil
}

// FindByGrupo busca produtos por grupo
func (r *ProdutoRepository) FindByGrupo(grupoID int, limit int) ([]models.Produto, error) {
	var produtos []models.Produto
	err := r.db.
		Preload("ProdutoEspecie").
		Preload("ProdutoGrupo").
		Preload("ProdutoMarca").
		Preload("ProdutoSubgrupo").
		Preload("ProdutoSerie").
		Preload("ProdutoTipoProduto").
		Where("produto_grupo_id = ? AND deleted_at IS NULL", grupoID).
		Limit(limit).
		Order("pro_id DESC").
		Find(&produtos).Error

	if err != nil {
		return nil, err
	}
	return produtos, nil
}

// FindByMarca busca produtos por marca
func (r *ProdutoRepository) FindByMarca(marcaID int, limit int) ([]models.Produto, error) {
	var produtos []models.Produto
	err := r.db.
		Preload("ProdutoEspecie").
		Preload("ProdutoGrupo").
		Preload("ProdutoMarca").
		Preload("ProdutoSubgrupo").
		Preload("ProdutoSerie").
		Preload("ProdutoTipoProduto").
		Where("produto_marca_id = ? AND deleted_at IS NULL", marcaID).
		Limit(limit).
		Order("pro_id DESC").
		Find(&produtos).Error

	if err != nil {
		return nil, err
	}
	return produtos, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de produtos com paginação e filtros
func (r *ProdutoRepository) List(limit, offset int, filters map[string]interface{}) ([]models.Produto, int64, error) {
	var produtos []models.Produto
	var total int64

	// Construir a query base
	query := r.db.Model(&models.Produto{}).Where("deleted_at IS NULL")

	// Aplicar filtros dinamicamente
	query = utils.ApplyFilters(query, models.Produto{}, filters)

	// Contar total de registros
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Buscar com paginação e relacionamentos
	err := query.
		Preload("ProdutoEspecie").
		Preload("ProdutoGrupo").
		Preload("ProdutoMarca").
		Preload("ProdutoSubgrupo").
		Preload("ProdutoSerie").
		Preload("ProdutoTipoProduto").
		Limit(limit).
		Offset(offset).
		Order("pro_id DESC").
		Find(&produtos).Error

	if err != nil {
		return nil, 0, err
	}

	return produtos, total, nil
}

// ============================================================
// MÉTODOS ADICIONAIS
// ============================================================

// ExistsByCodigo verifica se já existe um produto com o código
func (r *ProdutoRepository) ExistsByCodigo(codigo string, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.Produto{}).
		Where("pro_codigo = ? AND deleted_at IS NULL", codigo)

	// Se for para excluir um ID (atualização)
	if excludeID > 0 {
		query = query.Where("pro_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// ExistsByReferencia verifica se já existe um produto com a referência
func (r *ProdutoRepository) ExistsByReferencia(referencia string, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.Produto{}).
		Where("pro_referencia = ? AND deleted_at IS NULL", referencia)

	// Se for para excluir um ID (atualização)
	if excludeID > 0 {
		query = query.Where("pro_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// CountByGrupo conta quantos produtos existem em um grupo
func (r *ProdutoRepository) CountByGrupo(grupoID int) (int64, error) {
	var count int64
	err := r.db.Model(&models.Produto{}).
		Where("produto_grupo_id = ? AND deleted_at IS NULL", grupoID).
		Count(&count).Error

	return count, err
}

// CountByMarca conta quantos produtos existem em uma marca
func (r *ProdutoRepository) CountByMarca(marcaID int) (int64, error) {
	var count int64
	err := r.db.Model(&models.Produto{}).
		Where("produto_marca_id = ? AND deleted_at IS NULL", marcaID).
		Count(&count).Error

	return count, err
}