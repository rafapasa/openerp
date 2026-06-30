package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/models"
)

// ============================================================
// TYPES
// ============================================================

// EntidadeDocumentoRepository é o repositório para documentos de entidade
type EntidadeDocumentoRepository struct {
	db *gorm.DB
}

// ============================================================
// CONSTRUCTOR
// ============================================================

// NewEntidadeDocumentoRepository cria uma nova instância
func NewEntidadeDocumentoRepository(db *gorm.DB) *EntidadeDocumentoRepository {
	return &EntidadeDocumentoRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD
// ============================================================

// Create salva um novo documento com sequencial manual
func (r *EntidadeDocumentoRepository) Create(documento *models.EntidadeDocumento) error {
	// 1. Buscar o próximo número para esta entidade
	var maxItem int
	err := r.db.Model(&models.EntidadeDocumento{}).
		Where("ent_id = ?", documento.EntidadeID).
		Select("COALESCE(MAX(edoc_item), 0) + 1").
		Scan(&maxItem).Error
	if err != nil {
		return err
	}

	// 2. Atribuir o próximo número
	documento.Item = maxItem

	// 3. Salvar
	return r.db.Create(documento).Error
}

// Update atualiza um documento existente
func (r *EntidadeDocumentoRepository) Update(documento *models.EntidadeDocumento) error {
	return r.db.Save(documento).Error
}

// Delete realiza exclusão lógica de um documento
func (r *EntidadeDocumentoRepository) Delete(entidadeID, item int) error {
	// 1. Buscar o documento
	documento, err := r.FindByID(entidadeID, item)
	if err != nil {
		return err
	}

	// 2. Verificar se já foi deletado
	if documento.IsDeleted() {
		return errors.New("documento já foi deletado")
	}

	// 3. Realizar soft delete
	documento.SoftDelete()

	// 4. Salvar
	return r.db.Save(documento).Error
}

// ============================================================
// MÉTODOS DE BUSCA
// ============================================================

// FindByID busca um documento pelo ID composto (ent_id + edoc_item)
func (r *EntidadeDocumentoRepository) FindByID(entidadeID, item int) (*models.EntidadeDocumento, error) {
	var documento models.EntidadeDocumento
	err := r.db.
		Where("ent_id = ? AND edoc_item = ? AND deleted_at IS NULL", entidadeID, item).
		First(&documento).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("documento não encontrado")
		}
		return nil, err
	}
	return &documento, nil
}

// FindByEntidadeID busca todos os documentos de uma entidade
func (r *EntidadeDocumentoRepository) FindByEntidadeID(entidadeID int) ([]models.EntidadeDocumento, error) {
	var documentos []models.EntidadeDocumento
	err := r.db.
		Where("ent_id = ? AND deleted_at IS NULL", entidadeID).
		Order("edoc_item ASC").
		Find(&documentos).Error

	if err != nil {
		return nil, err
	}
	return documentos, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de documentos com paginação e filtros
func (r *EntidadeDocumentoRepository) List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeDocumento, int64, error) {
	var documentos []models.EntidadeDocumento
	var total int64

	// Construir query base
	query := r.db.Model(&models.EntidadeDocumento{}).Where("deleted_at IS NULL")

	// Aplicar filtros
	query = r.applyFilters(query, filters)

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Buscar com paginação
	err := query.
		Limit(limit).
		Offset(offset).
		Order("edoc_item DESC").
		Find(&documentos).Error

	if err != nil {
		return nil, 0, err
	}

	return documentos, total, nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// applyFilters aplica os filtros à query
func (r *EntidadeDocumentoRepository) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for key, value := range filters {
		if value == nil || value == "" {
			continue
		}

		switch key {
		case "entidade_id":
			query = query.Where("ent_id = ?", value)
		case "tipo":
			query = query.Where("edoc_tipo = ?", value)
		case "descricao":
			query = query.Where("edoc_descricao LIKE ?", "%"+value.(string)+"%")
		}
	}
	return query
}
