package repository

import (
	"fmt"

	"github.com/openerp/backend/internal/models"
	"gorm.io/gorm"
)

// ============================================================
// REPOSITORY: DocumentoVendaRepository
// ============================================================

type DocumentoVendaRepository struct {
	db *gorm.DB
}

// NewDocumentoVendaRepository cria uma nova instância do repositório
func NewDocumentoVendaRepository(db *gorm.DB) *DocumentoVendaRepository {
	return &DocumentoVendaRepository{db: db}
}

// Create cria um novo documento de venda no banco de dados
func (r *DocumentoVendaRepository) Create(doc *models.DocumentoVenda) error {
	return r.db.Create(doc).Error
}

// FindByID busca um documento de venda pelo ID, incluindo seus relacionamentos
func (r *DocumentoVendaRepository) FindByID(id int) (*models.DocumentoVenda, error) {
	var doc models.DocumentoVenda
	err := r.db.
		Preload("Itens").
		Preload("Pagamentos").
		Preload("Entidade").
		Preload("CondicaoPagamento").
		First(&doc, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("documento de venda com ID %d não encontrado", id)
		}
		return nil, err
	}
	return &doc, nil
}

// FindByID busca um documento de venda pelo ID, sem incluir seus relacionamentos
func (r *DocumentoVendaRepository) GetByID(id int) (*models.DocumentoVenda, error) {
	var ddv models.DocumentoVenda
	err := r.db.First(&ddv, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("documento de venda com ID %d não encontrado", id)
		}
		return nil, err
	}
	return &ddv, nil
}

// Update atualiza um documento de venda existente no banco de dados
func (r *DocumentoVendaRepository) Update(doc *models.DocumentoVenda) error {
	return r.db.Save(doc).Error
}

// Delete realiza a exclusão lógica de um documento de venda
func (r *DocumentoVendaRepository) Delete(id int) error {
	return r.db.Model(&models.DocumentoVenda{}).Where("ddv_id = ?", id).Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

// List retorna uma lista paginada de documentos de venda com base nos filtros
func (r *DocumentoVendaRepository) List(limit, offset int, filters map[string]interface{}) ([]models.DocumentoVenda, int64, error) {
	var documentos []models.DocumentoVenda
	var total int64

	query := r.db.Model(&models.DocumentoVenda{})
	countQuery := r.db.Model(&models.DocumentoVenda{})

	// Aplicar filtros
	for key, value := range filters {
		switch key {
		case "entidade_id", "tipo_documento", "tipo_operacao", "situacao", "empresa_filial_id":
			if v, ok := value.(int); ok && v > 0 {
				query = query.Where(fmt.Sprintf("ddv_%s = ?", key), v)
				countQuery = countQuery.Where(fmt.Sprintf("ddv_%s = ?", key), v)
			}
		case "data_documento_ini":
			query = query.Where("ddv_datadocumento >= ?", value)
			countQuery = countQuery.Where("ddv_datadocumento >= ?", value)
		case "data_documento_fim":
			query = query.Where("ddv_datadocumento <= ?", value)
			countQuery = countQuery.Where("ddv_datadocumento <= ?", value)
		}
	}

	// Contar o total de registros
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginação e ordenação
	err := query.
		Preload("Entidade").
		Preload("CondicaoPagamento").
		Limit(limit).
		Offset(offset).
		Order("ddv_datadocumento DESC, ddv_id DESC").
		Find(&documentos).Error

	if err != nil {
		return nil, 0, err
	}

	return documentos, total, nil
}

// ExistsByNumero verifica se um documento com um número específico já existe para uma filial
func (r *DocumentoVendaRepository) ExistsByNumero(numero, empresaFilialID int, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.DocumentoVenda{}).
		Where("ddv_numero = ?", numero).
		Where("emf_id = ?", empresaFilialID)

	if excludeID > 0 {
		query = query.Where("ddv_id <> ?", excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
