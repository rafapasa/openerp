// internal/repository/produto_repository.go
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

// ProdutoRepository define o contrato para operações de banco
type ProdutoRepository interface {
	// CRUD Básico
	Create(produto *models.Produto) error
	Update(id int, produto *models.Produto) error
	Delete(id int) error
	FindByID(id int) (*models.Produto, error)
	GetByID(id int) (*models.Produto, error)

	// Buscas Específicas
	FindByCodigo(codigo int) (*models.Produto, error)
	FindByCodigoAndFilial(codigo int, filialID int) (*models.Produto, error)
	FindByNome(nome string, limit int) ([]models.Produto, error)
	FindByResumo(resumo string, limit int) ([]models.Produto, error)
	FindByCodigoBarras(codigoBarras string) (*models.Produto, error)
	FindBySituacao(situacao int) ([]models.Produto, error)
	FindByGrupoID(grupoID int) ([]models.Produto, error)
	FindBySubgrupoID(subgrupoID int) ([]models.Produto, error)
	FindByMarcaID(marcaID int) ([]models.Produto, error)
	FindByModeloID(modeloID int) ([]models.Produto, error)
	FindByFilialID(filialID int) ([]models.Produto, error)
	FindActive() ([]models.Produto, error)
	FindActiveByFilial(filialID int) ([]models.Produto, error)

	// Listagem com Filtros
	List(limit, offset int, filters map[string]interface{}) ([]models.Produto, int64, error)
	ListWithFullPreload(limit, offset int, filters map[string]interface{}) ([]models.Produto, int64, error)
	FindAll() ([]models.Produto, error)

	// Consultas de Validação (APENAS CONSULTAS)
	ExistsByCodigo(codigo int, excludeID int) (bool, error)
	ExistsByCodigoAndFilial(codigo int, filialID int, excludeID int) (bool, error)
	ExistsByCodigoBarras(codigoBarras string, excludeID int) (bool, error)
	ExistsByNome(nome string, excludeID int) (bool, error)
	ExistsByID(id int) (bool, error)
	CountByFilialID(filialID int) (int64, error)
	GetMaxCodigoByFilial(filialID int) (int, error)

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

type produtoRepository struct {
	db *gorm.DB
}

// NewProdutoRepository cria uma nova instância (retorna a interface)
func NewProdutoRepository(db *gorm.DB) ProdutoRepository {
	return &produtoRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD (APENAS PERSISTÊNCIA)
// ============================================================

// Create salva um novo produto
func (r *produtoRepository) Create(produto *models.Produto) error {
	return r.db.Create(produto).Error
}

// Update atualiza um produto existente
func (r *produtoRepository) Update(id int, produto *models.Produto) error {
	return r.db.
		Omit("TipoProduto", "ProdutoGrupo", "ProdutoSubgrupo", "Marca", "Modelo", "Serie", "Especie", "ItensPedido", "created_at", "deleted_at").
		Model(&models.Produto{}).
		Where("pro_id = ?", id).
		Updates(produto).Error
}

// Delete realiza exclusão lógica
func (r *produtoRepository) Delete(id int) error {
	return r.db.
		Model(&models.Produto{}).
		Where("pro_id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// FindByID busca um produto pelo ID com relacionamentos
func (r *produtoRepository) FindByID(id int) (*models.Produto, error) {
	var produto models.Produto
	err := r.db.
		Preload("TipoProduto").
		Preload("ProdutoGrupo").
		Preload("ProdutoSubgrupo").
		Preload("Marca").
		Preload("Modelo").
		Preload("Serie").
		Preload("Especie").
		Preload("ItensPedido", "deleted_at IS NULL").
		Where("pro_id = ? AND deleted_at IS NULL", id).
		First(&produto).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("produto com ID %d não encontrado", id))
		}
		return nil, apperrors.NewInternalError("Erro ao buscar produto.", err)
	}
	return &produto, nil
}

// GetByID busca um produto pelo ID (sem relacionamentos)
func (r *produtoRepository) GetByID(id int) (*models.Produto, error) {
	var produto models.Produto
	err := r.db.
		Where("pro_id = ? AND deleted_at IS NULL", id).
		First(&produto).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("produto com ID %d não encontrado", id))
		}
		return nil, err
	}
	return &produto, nil
}

// ============================================================
// MÉTODOS DE BUSCA ESPECÍFICOS
// ============================================================

// FindByCodigo busca um produto pelo código
func (r *produtoRepository) FindByCodigo(codigo int) (*models.Produto, error) {
	var produto models.Produto
	err := r.db.
		Preload("TipoProduto").
		Preload("ProdutoGrupo").
		Where("pro_codigo = ? AND deleted_at IS NULL", codigo).
		First(&produto).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("produto com código %d não encontrado", codigo))
		}
		return nil, apperrors.NewInternalError("Erro ao buscar produto.", err)
	}
	return &produto, nil
}

// FindByCodigoAndFilial busca um produto pelo código e filial
func (r *produtoRepository) FindByCodigoAndFilial(codigo int, filialID int) (*models.Produto, error) {
	var produto models.Produto
	err := r.db.
		Preload("TipoProduto").
		Preload("ProdutoGrupo").
		Where("pro_codigo = ? AND emf_id = ? AND deleted_at IS NULL", codigo, filialID).
		First(&produto).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("produto com código %d e filial %d não encontrado", codigo, filialID))
		}
		return nil, err
	}
	return &produto, nil
}

// FindByNome busca produtos pelo nome (autocomplete)
func (r *produtoRepository) FindByNome(nome string, limit int) ([]models.Produto, error) {
	var produtos []models.Produto
	err := r.db.
		Preload("TipoProduto").
		Where("pro_nome LIKE ? AND deleted_at IS NULL", "%"+nome+"%").
		Order("pro_nome ASC").
		Limit(limit).
		Find(&produtos).Error

	if err != nil {
		return nil, err
	}
	return produtos, nil
}

// FindByResumo busca produtos pelo resumo (autocomplete)
func (r *produtoRepository) FindByResumo(resumo string, limit int) ([]models.Produto, error) {
	var produtos []models.Produto
	err := r.db.
		Where("pro_resumo LIKE ? AND deleted_at IS NULL", "%"+resumo+"%").
		Order("pro_resumo ASC").
		Limit(limit).
		Find(&produtos).Error

	if err != nil {
		return nil, err
	}
	return produtos, nil
}

// FindByCodigoBarras busca um produto pelo código de barras
func (r *produtoRepository) FindByCodigoBarras(codigoBarras string) (*models.Produto, error) {
	var produto models.Produto
	err := r.db.
		Preload("TipoProduto").
		Where("pro_codigobarra = ? AND deleted_at IS NULL", codigoBarras).
		First(&produto).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("produto com código de barras %s não encontrado", codigoBarras))
		}
		return nil, apperrors.NewInternalError("Erro ao buscar produto.", err)
	}
	return &produto, nil
}

// FindBySituacao busca produtos por situação
func (r *produtoRepository) FindBySituacao(situacao int) ([]models.Produto, error) {
	var produtos []models.Produto
	err := r.db.
		Preload("TipoProduto").
		Where("pro_situacao = ? AND deleted_at IS NULL", situacao).
		Order("pro_nome ASC").
		Find(&produtos).Error

	if err != nil {
		return nil, err
	}
	return produtos, nil
}

// FindByGrupoID busca produtos por grupo
func (r *produtoRepository) FindByGrupoID(grupoID int) ([]models.Produto, error) {
	var produtos []models.Produto
	err := r.db.
		Preload("TipoProduto").
		Where("prog_id = ? AND deleted_at IS NULL", grupoID).
		Order("pro_nome ASC").
		Find(&produtos).Error

	if err != nil {
		return nil, err
	}
	return produtos, nil
}

// FindBySubgrupoID busca produtos por subgrupo
func (r *produtoRepository) FindBySubgrupoID(subgrupoID int) ([]models.Produto, error) {
	var produtos []models.Produto
	err := r.db.
		Preload("TipoProduto").
		Where("prosg_id = ? AND deleted_at IS NULL", subgrupoID).
		Order("pro_nome ASC").
		Find(&produtos).Error

	if err != nil {
		return nil, err
	}
	return produtos, nil
}

// FindByMarcaID busca produtos por marca
func (r *produtoRepository) FindByMarcaID(marcaID int) ([]models.Produto, error) {
	var produtos []models.Produto
	err := r.db.
		Preload("TipoProduto").
		Where("promar_id = ? AND deleted_at IS NULL", marcaID).
		Order("pro_nome ASC").
		Find(&produtos).Error

	if err != nil {
		return nil, err
	}
	return produtos, nil
}

// FindByModeloID busca produtos por modelo
func (r *produtoRepository) FindByModeloID(modeloID int) ([]models.Produto, error) {
	var produtos []models.Produto
	err := r.db.
		Preload("TipoProduto").
		Where("prom_id = ? AND deleted_at IS NULL", modeloID).
		Order("pro_nome ASC").
		Find(&produtos).Error

	if err != nil {
		return nil, err
	}
	return produtos, nil
}

// FindByFilialID busca produtos por filial
func (r *produtoRepository) FindByFilialID(filialID int) ([]models.Produto, error) {
	var produtos []models.Produto
	err := r.db.
		Preload("TipoProduto").
		Where("emf_id = ? AND deleted_at IS NULL", filialID).
		Order("pro_nome ASC").
		Find(&produtos).Error

	if err != nil {
		return nil, err
	}
	return produtos, nil
}

// FindActive busca todos os produtos ativos
func (r *produtoRepository) FindActive() ([]models.Produto, error) {
	var produtos []models.Produto
	err := r.db.
		Preload("TipoProduto").
		Preload("ProdutoGrupo").
		Where("pro_situacao = 1 AND deleted_at IS NULL").
		Order("pro_nome ASC").
		Find(&produtos).Error

	if err != nil {
		return nil, err
	}
	return produtos, nil
}

// FindActiveByFilial busca produtos ativos por filial
func (r *produtoRepository) FindActiveByFilial(filialID int) ([]models.Produto, error) {
	var produtos []models.Produto
	err := r.db.
		Preload("TipoProduto").
		Preload("ProdutoGrupo").
		Where("emf_id = ? AND pro_situacao = 1 AND deleted_at IS NULL", filialID).
		Order("pro_nome ASC").
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
func (r *produtoRepository) List(limit, offset int, filters map[string]interface{}) ([]models.Produto, int64, error) {
	var produtos []models.Produto
	var total int64

	query := r.db.Model(&models.Produto{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.Produto{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("TipoProduto").
		Preload("ProdutoGrupo").
		Preload("Marca").
		Preload("Modelo").
		Limit(limit).
		Offset(offset).
		Order("pro_nome ASC").
		Find(&produtos).Error

	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro ao listar produtos.", err)
	}

	return produtos, total, nil
}

// ListWithFullPreload retorna uma lista com todos os relacionamentos
func (r *produtoRepository) ListWithFullPreload(limit, offset int, filters map[string]interface{}) ([]models.Produto, int64, error) {
	var produtos []models.Produto
	var total int64

	query := r.db.Model(&models.Produto{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.Produto{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("TipoProduto").
		Preload("ProdutoGrupo").
		Preload("ProdutoSubgrupo").
		Preload("Marca").
		Preload("Modelo").
		Preload("Serie").
		Preload("Especie").
		Limit(limit).
		Offset(offset).
		Order("pro_nome ASC").
		Find(&produtos).Error

	if err != nil {
		return nil, 0, err
	}

	return produtos, total, nil
}

// FindAll busca todos os produtos
func (r *produtoRepository) FindAll() ([]models.Produto, error) {
	var produtos []models.Produto
	err := r.db.
		Where("deleted_at IS NULL").
		Order("pro_nome ASC").
		Find(&produtos).Error
	if err != nil {
		return nil, err
	}
	return produtos, nil
}

// ============================================================
// MÉTODOS DE CONSULTA PARA VALIDAÇÕES (APENAS CONSULTAS)
// ============================================================

// ExistsByCodigo verifica se já existe um produto com o código
func (r *produtoRepository) ExistsByCodigo(codigo int, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.Produto{}).
		Where("pro_codigo = ? AND deleted_at IS NULL", codigo)

	if excludeID > 0 {
		query = query.Where("pro_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByCodigoAndFilial verifica se já existe um produto com o código e filial
func (r *produtoRepository) ExistsByCodigoAndFilial(codigo int, filialID int, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.Produto{}).
		Where("pro_codigo = ? AND emf_id = ? AND deleted_at IS NULL", codigo, filialID)

	if excludeID > 0 {
		query = query.Where("pro_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByCodigoBarras verifica se já existe um produto com o código de barras
func (r *produtoRepository) ExistsByCodigoBarras(codigoBarras string, excludeID int) (bool, error) {
	if codigoBarras == "" {
		return false, nil
	}

	var count int64
	query := r.db.Model(&models.Produto{}).
		Where("pro_codigobarra = ? AND deleted_at IS NULL", codigoBarras)

	if excludeID > 0 {
		query = query.Where("pro_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByNome verifica se já existe um produto com o nome
func (r *produtoRepository) ExistsByNome(nome string, excludeID int) (bool, error) {
	if nome == "" {
		return false, nil
	}

	var count int64
	query := r.db.Model(&models.Produto{}).
		Where("pro_nome = ? AND deleted_at IS NULL", nome)

	if excludeID > 0 {
		query = query.Where("pro_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByID verifica se um produto existe pelo ID
func (r *produtoRepository) ExistsByID(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.Produto{}).
		Where("pro_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	if err != nil {
		return false, apperrors.NewInternalError("Erro ao buscar produto.", err)
	}
	return count > 0, nil
}

// CountByFilialID retorna a quantidade de produtos de uma filial
func (r *produtoRepository) CountByFilialID(filialID int) (int64, error) {
	var count int64
	err := r.db.Model(&models.Produto{}).
		Where("emf_id = ? AND deleted_at IS NULL", filialID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetMaxCodigoByFilial retorna o maior código de produto para uma filial
func (r *produtoRepository) GetMaxCodigoByFilial(filialID int) (int, error) {
	var maxCodigo int
	err := r.db.Model(&models.Produto{}).
		Where("emf_id = ? AND deleted_at IS NULL", filialID).
		Select("COALESCE(MAX(pro_codigo), 0)").
		Scan(&maxCodigo).Error
	if err != nil {
		return 0, err
	}
	return maxCodigo, nil
}

// ============================================================
// OPERAÇÕES EM LOTE
// ============================================================

// BulkUpdateStatus atualiza a situação de múltiplos produtos
func (r *produtoRepository) BulkUpdateStatus(ids []int, situacao int) error {
	return r.db.Model(&models.Produto{}).
		Where("pro_id IN ? AND deleted_at IS NULL", ids).
		Update("pro_situacao", situacao).Error
}

// BulkDelete realiza exclusão lógica de múltiplos produtos
func (r *produtoRepository) BulkDelete(ids []int) error {
	return r.db.Model(&models.Produto{}).
		Where("pro_id IN ? AND deleted_at IS NULL", ids).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// ============================================================
// MÉTODOS DE CONSULTA DE DEPENDÊNCIAS
// ============================================================

// HasDependentRecords verifica se o produto tem registros dependentes
func (r *produtoRepository) HasDependentRecords(id int) (bool, error) {
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
func (r *produtoRepository) CountDependentRecords(id int) (map[string]int64, error) {
	result := make(map[string]int64)

	// Verifica itens de pedido
	var countItensPedido int64
	if err := r.db.Model(&models.DocumentoVendaItem{}).
		Where("pro_id = ? AND deleted_at IS NULL", id).
		Count(&countItensPedido).Error; err != nil {
		return nil, err
	}
	if countItensPedido > 0 {
		result["itens_pedido"] = countItensPedido
	}

	// Verifica tabelas de preço
	var countTabelaPreco int64
	if err := r.db.Model(&models.TabelaPrecoProduto{}).
		Where("pro_id = ? AND deleted_at IS NULL", id).
		Count(&countTabelaPreco).Error; err != nil {
		return nil, err
	}
	if countTabelaPreco > 0 {
		result["tabela_preco"] = countTabelaPreco
	}

	return result, nil
}
