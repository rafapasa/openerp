package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/models"
)

// ============================================================
// TYPES
// ============================================================

// EntidadeContatoRepository é o repositório para contatos de entidade
type EntidadeContatoRepository struct {
	db *gorm.DB
}

// ============================================================
// CONSTRUCTOR
// ============================================================

// NewEntidadeContatoRepository cria uma nova instância
func NewEntidadeContatoRepository(db *gorm.DB) *EntidadeContatoRepository {
	return &EntidadeContatoRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD
// ============================================================

// Create salva um novo contato com sequencial manual
func (r *EntidadeContatoRepository) Create(contato *models.EntidadeContato) error {
	// 1. Buscar o próximo número para esta entidade
	var maxItem int
	err := r.db.Model(&models.EntidadeContato{}).
		Where("ent_id = ?", contato.EntidadeID).
		Select("COALESCE(MAX(efc_item), 0) + 1").
		Scan(&maxItem).Error
	if err != nil {
		return err
	}

	// 2. Atribuir o próximo número
	contato.Item = maxItem

	// 3. Salvar
	return r.db.Create(contato).Error
}

// Update atualiza um contato existente
func (r *EntidadeContatoRepository) Update(contato *models.EntidadeContato) error {
	return r.db.
		Omit("FormaContato").
		Model(&models.EntidadeContato{}).
		Where("ent_id = ? AND efc_item = ?", contato.EntidadeID, contato.Item).
		Updates(contato).Error
}

// Delete realiza exclusão lógica de um contato
func (r *EntidadeContatoRepository) Delete(entidadeID, item int) error {
	// 1. Buscar o contato
	contato, err := r.FindByID(entidadeID, item)
	if err != nil {
		return err
	}

	// 2. Verificar se já foi deletado
	if contato.IsDeleted() {
		return errors.New("contato já foi deletado")
	}

	// 3. Realizar soft delete
	contato.SoftDelete()

	// 4. Salvar
	return r.Update(contato)
}

// ============================================================
// MÉTODOS DE BUSCA
// ============================================================

// FindByID busca um contato pelo ID composto (ent_id + efc_item)
func (r *EntidadeContatoRepository) FindByID(entidadeID, item int) (*models.EntidadeContato, error) {
	var contato models.EntidadeContato
	err := r.db.
		Preload("FormaContato").
		Where("ent_id = ? AND efc_item = ? AND deleted_at IS NULL", entidadeID, item).
		First(&contato).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("contato não encontrado")
		}
		return nil, err
	}
	return &contato, nil
}

// FindByEntidadeID busca todos os contatos de uma entidade
func (r *EntidadeContatoRepository) FindByEntidadeID(entidadeID int) ([]models.EntidadeContato, error) {
	var contatos []models.EntidadeContato
	err := r.db.
		Preload("FormaContato").
		Where("ent_id = ? AND deleted_at IS NULL", entidadeID).
		Order("efc_item ASC").
		Find(&contatos).Error

	if err != nil {
		return nil, err
	}
	return contatos, nil
}

// FindByEntidadeIDAndTipo busca contatos de uma entidade por tipo
func (r *EntidadeContatoRepository) FindByEntidadeIDAndTipo(entidadeID, formaContatoID int) ([]models.EntidadeContato, error) {
	var contatos []models.EntidadeContato
	err := r.db.
		Preload("FormaContato").
		Where("ent_id = ? AND frc_id = ? AND deleted_at IS NULL", entidadeID, formaContatoID).
		Order("efc_item ASC").
		Find(&contatos).Error

	if err != nil {
		return nil, err
	}
	return contatos, nil
}

// ============================================================
// MÉTODOS DE VERIFICAÇÃO
// ============================================================

// ExistsByEntidadeTipo verifica se a entidade já tem um contato do tipo especificado
func (r *EntidadeContatoRepository) ExistsByEntidadeTipo(entidadeID, formaContatoID int, excludeItem int) (bool, error) {
	var count int64
	query := r.db.Model(&models.EntidadeContato{}).
		Where("ent_id = ? AND frc_id = ? AND deleted_at IS NULL", entidadeID, formaContatoID)

	if excludeItem > 0 {
		query = query.Where("efc_item != ?", excludeItem)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de contatos com paginação e filtros
func (r *EntidadeContatoRepository) List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeContato, int64, error) {
	var contatos []models.EntidadeContato
	var total int64

	// Construir query base
	query := r.db.Model(&models.EntidadeContato{}).Where("deleted_at IS NULL")

	// Aplicar filtros
	query = r.applyFilters(query, filters)

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Buscar com paginação e relacionamentos
	err := query.
		Preload("FormaContato").
		Limit(limit).
		Offset(offset).
		Order("efc_item DESC").
		Find(&contatos).Error

	if err != nil {
		return nil, 0, err
	}

	return contatos, total, nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// applyFilters aplica os filtros à query
func (r *EntidadeContatoRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for key, value := range filters {
		if value == nil || value == "" {
			continue
		}

		switch key {
		case "entidade_id":
			query = query.Where("ent_id = ?", value)
		case "forma_contato_id":
			query = query.Where("frc_id = ?", value)
		case "informacao":
			query = query.Where("efc_informacao LIKE ?", "%"+value.(string)+"%")
		}
	}
	return query
}
