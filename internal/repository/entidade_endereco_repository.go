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

type EntidadeEnderecoRepository struct {
	db *gorm.DB
}

// ============================================================
// CONSTRUCTOR
// ============================================================

func NewEntidadeEnderecoRepository(db *gorm.DB) *EntidadeEnderecoRepository {
	return &EntidadeEnderecoRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD
// ============================================================

// Create salva um novo endereço com sequencial manual
func (r *EntidadeEnderecoRepository) Create(endereco *models.EntidadeEndereco) error {
	// 1. Buscar o próximo número para esta entidade
	var maxItem int
	err := r.db.Model(&models.EntidadeEndereco{}).
		Where("ent_id = ?", endereco.EntidadeID).
		Select("COALESCE(MAX(ete_item), 0) + 1").
		Scan(&maxItem).Error
	if err != nil {
		return err
	}

	// 2. Atribuir o próximo número
	endereco.Item = maxItem

	// 3. Salvar
	return r.db.Create(endereco).Error
}

// Update atualiza um endereço existente
func (r *EntidadeEnderecoRepository) Update(endereco *models.EntidadeEndereco) error {
	return r.db.Save(endereco).Error
}

// Delete realiza exclusão lógica
func (r *EntidadeEnderecoRepository) Delete(entidadeID, item int) error {
	// 1. Buscar o endereço
	endereco, err := r.FindByID(entidadeID, item)
	if err != nil {
		return err
	}

	// 2. Verificar se já foi deletado
	if endereco.IsDeleted() {
		return errors.New("endereço já foi deletado")
	}

	// 3. Realizar soft delete
	endereco.SoftDelete()

	// 4. Salvar
	return r.db.Save(endereco).Error
}

// ============================================================
// MÉTODOS DE BUSCA
// ============================================================

// FindByID busca um endereço pelo ID composto (ent_id + ete_item)
func (r *EntidadeEnderecoRepository) FindByID(entidadeID, item int) (*models.EntidadeEndereco, error) {
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
func (r *EntidadeEnderecoRepository) FindByEntidadeID(entidadeID int) ([]models.EntidadeEndereco, error) {
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
func (r *EntidadeEnderecoRepository) FindByEntidadeIDAndTipo(entidadeID, tipo int) ([]models.EntidadeEndereco, error) {
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
// MÉTODOS DE VERIFICAÇÃO
// ============================================================

// ExistsByEntidadeTipo verifica se a entidade já tem um endereço do tipo especificado
func (r *EntidadeEnderecoRepository) ExistsByEntidadeTipo(entidadeID, tipo int, excludeItem int) (bool, error) {
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

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de endereços com paginação e filtros
func (r *EntidadeEnderecoRepository) List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeEndereco, int64, error) {
	var enderecos []models.EntidadeEndereco
	var total int64

	// Construir query base
	query := r.db.Model(&models.EntidadeEndereco{}).Where("deleted_at IS NULL")

	// Aplicar filtros
	query = utils.ApplyFilters(query, models.EntidadeEndereco{}, filters)

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Buscar com paginação e relacionamentos
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
