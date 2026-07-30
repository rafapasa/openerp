// internal/repository/entidade_endereco_repository.go
package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// INTERFACE
// ============================================================

// EntidadeEnderecoRepository define as operações para endereços de entidade
type EntidadeEnderecoRepository interface {
	// CRUD Básico
	Create(endereco *models.EntidadeEndereco) error
	Update(endereco *models.EntidadeEndereco) error
	Delete(entidadeID, item int) error
	FindByID(entidadeID, item int) (*models.EntidadeEndereco, error)

	// Buscas Específicas
	FindByEntidadeID(entidadeID int) ([]models.EntidadeEndereco, error)
	FindByEntidadeIDAndTipo(entidadeID, tipo int) ([]models.EntidadeEndereco, error)

	// Listagem com Filtros
	List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeEndereco, int64, error)

	// Consultas de Validação (APENAS CONSULTAS)
	ExistsByEntidadeTipo(entidadeID, tipo int, excludeItem int) (bool, error)
	GetNextItemNumber(entidadeID int) (int, error)
	CountByEntidadeID(entidadeID int) (int64, error)
}

// ============================================================
// IMPLEMENTAÇÃO CONCRETA
// ============================================================

type entidadeEnderecoRepository struct {
	db *gorm.DB
}

// NewEntidadeEnderecoRepository cria uma nova instância do repositório
// ✅ Retorna a interface, não a struct concreta
func NewEntidadeEnderecoRepository(db *gorm.DB) EntidadeEnderecoRepository {
	return &entidadeEnderecoRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD (APENAS PERSISTÊNCIA)
// ============================================================

// Create salva um novo endereço com sequencial manual
func (r *entidadeEnderecoRepository) Create(endereco *models.EntidadeEndereco) error {
	return r.db.Create(endereco).Error
}

// Update atualiza um endereço existente
func (r *entidadeEnderecoRepository) Update(endereco *models.EntidadeEndereco) error {
	return r.db.
		Omit("municipio", "estado", "pais", "created_at", "deleted_at").
		Model(&models.EntidadeEndereco{}).
		Where("ent_id = ? AND ete_item = ?", endereco.EntidadeID, endereco.Item).
		Updates(endereco).Error
}

// Delete realiza exclusão lógica
func (r *entidadeEnderecoRepository) Delete(entidadeID, item int) error {
	return r.db.
		Model(&models.EntidadeEndereco{}).
		Where("ent_id = ? AND ete_item = ?", entidadeID, item).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// ============================================================
// MÉTODOS DE BUSCA
// ============================================================

// FindByID busca um endereço pelo ID composto (ent_id + ete_item)
func (r *entidadeEnderecoRepository) FindByID(entidadeID, item int) (*models.EntidadeEndereco, error) {
	var endereco models.EntidadeEndereco
	err := r.db.
		Preload("Pais").
		Preload("Estado").
		Preload("Municipio").
		Where("ent_id = ? AND ete_item = ? AND deleted_at IS NULL", entidadeID, item).
		First(&endereco).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("endereço não encontrado")
		}
		return nil, err
	}
	return &endereco, nil
}

// FindByEntidadeID busca todos os endereços de uma entidade
func (r *entidadeEnderecoRepository) FindByEntidadeID(entidadeID int) ([]models.EntidadeEndereco, error) {
	var enderecos []models.EntidadeEndereco
	err := r.db.
		Preload("Pais").
		Preload("Estado").
		Preload("Municipio").
		Where("ent_id = ? AND deleted_at IS NULL", entidadeID).
		Order("ete_item ASC").
		Find(&enderecos).Error

	if err != nil {
		return nil, err
	}
	return enderecos, nil
}

// FindByEntidadeIDAndTipo busca endereços de uma entidade por tipo
func (r *entidadeEnderecoRepository) FindByEntidadeIDAndTipo(entidadeID, tipo int) ([]models.EntidadeEndereco, error) {
	var enderecos []models.EntidadeEndereco
	err := r.db.
		Preload("Pais").
		Preload("Estado").
		Preload("Municipio").
		Where("ent_id = ? AND ete_tipo = ? AND deleted_at IS NULL", entidadeID, tipo).
		Order("ete_item ASC").
		Find(&enderecos).Error

	if err != nil {
		return nil, err
	}
	return enderecos, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de endereços com paginação e filtros
func (r *entidadeEnderecoRepository) List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeEndereco, int64, error) {
	var enderecos []models.EntidadeEndereco
	var total int64

	query := r.db.Model(&models.EntidadeEndereco{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.EntidadeEndereco{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("Pais").
		Preload("Estado").
		Preload("Municipio").
		Limit(limit).
		Offset(offset).
		Order("ete_item DESC").
		Find(&enderecos).Error

	if err != nil {
		return nil, 0, err
	}

	return enderecos, total, nil
}

// ============================================================
// MÉTODOS DE CONSULTA PARA VALIDAÇÕES (APENAS CONSULTAS)
// ============================================================

// ExistsByEntidadeTipo verifica se a entidade já tem um endereço do tipo especificado
func (r *entidadeEnderecoRepository) ExistsByEntidadeTipo(entidadeID, tipo int, excludeItem int) (bool, error) {
	var count int64
	query := r.db.Model(&models.EntidadeEndereco{}).
		Where("ent_id = ? AND ete_tipo = ? AND deleted_at IS NULL", entidadeID, tipo)

	if excludeItem > 0 {
		query = query.Where("ete_item != ?", excludeItem)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetNextItemNumber retorna o próximo número de item para uma entidade
func (r *entidadeEnderecoRepository) GetNextItemNumber(entidadeID int) (int, error) {
	var maxItem int
	err := r.db.Model(&models.EntidadeEndereco{}).
		Where("ent_id = ?", entidadeID).
		Select("COALESCE(MAX(ete_item), 0) + 1").
		Scan(&maxItem).Error
	if err != nil {
		return 0, err
	}
	return maxItem, nil
}

// CountByEntidadeID retorna a quantidade de endereços de uma entidade
func (r *entidadeEnderecoRepository) CountByEntidadeID(entidadeID int) (int64, error) {
	var count int64
	err := r.db.Model(&models.EntidadeEndereco{}).
		Where("ent_id = ? AND deleted_at IS NULL", entidadeID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
