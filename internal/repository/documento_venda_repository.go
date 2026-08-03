// internal/repository/documento_venda_repository.go
package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/models"
)

// ============================================================
// INTERFACE
// ============================================================

// DocumentoVendaRepository define o contrato para operações de banco
type DocumentoVendaRepository interface {
	// CRUD Básico
	Create(doc *models.DocumentoVenda) error
	FindByID(id int) (*models.DocumentoVenda, error)
	GetByID(id int) (*models.DocumentoVenda, error)
	Update(doc *models.DocumentoVenda) error
	Delete(id int) error

	// Listagem com Filtros
	List(limit, offset int, filters map[string]interface{}) ([]models.DocumentoVenda, int64, error)

	// Consultas de Validação (APENAS CONSULTAS)
	ExistsByNumero(numero, empresaFilialID int, excludeID int) (bool, error)
	ExistsByID(id int) (bool, error)
	CountByEntidadeID(entidadeID int) (int64, error)
	GetMaxNumero(empresaFilialID int) (int, error)
}

// ============================================================
// IMPLEMENTAÇÃO CONCRETA (privada)
// ============================================================

type documentoVendaRepository struct {
	db *gorm.DB
}

// NewDocumentoVendaRepository cria uma nova instância do repositório
// ✅ Retorna a interface, não a struct concreta
func NewDocumentoVendaRepository(db *gorm.DB) DocumentoVendaRepository {
	return &documentoVendaRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD (APENAS PERSISTÊNCIA)
// ============================================================

// Create cria um novo documento de venda no banco de dados
func (r *documentoVendaRepository) Create(doc *models.DocumentoVenda) error {
	return r.db.Create(doc).Error
}

// FindByID busca um documento de venda pelo ID, incluindo seus relacionamentos
func (r *documentoVendaRepository) FindByID(id int) (*models.DocumentoVenda, error) {
	var doc models.DocumentoVenda
	err := r.db.
		Preload("Itens").
		Preload("Pagamentos").
		Preload("Entidade").
		Preload("CondicaoPagamento").
		Where("ddv_id = ? AND deleted_at IS NULL", id).
		First(&doc).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("documento de venda com ID %d não encontrado", id)
		}
		return nil, err
	}
	return &doc, nil
}

// GetByID busca um documento de venda pelo ID, sem incluir seus relacionamentos
func (r *documentoVendaRepository) GetByID(id int) (*models.DocumentoVenda, error) {
	var doc models.DocumentoVenda
	err := r.db.
		Where("ddv_id = ? AND deleted_at IS NULL", id).
		First(&doc).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("documento de venda com ID %d não encontrado", id)
		}
		return nil, err
	}
	return &doc, nil
}

// Update atualiza um documento de venda existente no banco de dados
func (r *documentoVendaRepository) Update(doc *models.DocumentoVenda) error {
	return r.db.
		Omit("created_at", "deleted_at").
		Model(&models.DocumentoVenda{}).
		Where("ddv_id = ?", doc.ID).
		Updates(doc).Error
}

// Delete realiza a exclusão lógica de um documento de venda
func (r *documentoVendaRepository) Delete(id int) error {
	return r.db.
		Model(&models.DocumentoVenda{}).
		Where("ddv_id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista paginada de documentos de venda com base nos filtros
func (r *documentoVendaRepository) List(limit, offset int, filters map[string]interface{}) ([]models.DocumentoVenda, int64, error) {
	var documentos []models.DocumentoVenda
	var total int64

	query := r.db.Model(&models.DocumentoVenda{}).Where("deleted_at IS NULL")
	countQuery := r.db.Model(&models.DocumentoVenda{}).Where("deleted_at IS NULL")

	// Aplicar filtros
	for key, value := range filters {
		if value == nil {
			continue
		}

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
		case "numero":
			if v, ok := value.(int); ok && v > 0 {
				query = query.Where("ddv_numero = ?", v)
				countQuery = countQuery.Where("ddv_numero = ?", v)
			}
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

// ============================================================
// MÉTODOS DE CONSULTA PARA VALIDAÇÕES (APENAS CONSULTAS)
// ============================================================

// ExistsByNumero verifica se um documento com um número específico já existe para uma filial
func (r *documentoVendaRepository) ExistsByNumero(numero, empresaFilialID int, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.DocumentoVenda{}).
		Where("ddv_numero = ?", numero).
		Where("emf_id = ?", empresaFilialID).
		Where("deleted_at IS NULL")

	if excludeID > 0 {
		query = query.Where("ddv_id != ?", excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ExistsByID verifica se um documento existe pelo ID
func (r *documentoVendaRepository) ExistsByID(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.DocumentoVenda{}).
		Where("ddv_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountByEntidadeID retorna a quantidade de documentos de uma entidade
func (r *documentoVendaRepository) CountByEntidadeID(entidadeID int) (int64, error) {
	var count int64
	err := r.db.Model(&models.DocumentoVenda{}).
		Where("ent_id = ? AND deleted_at IS NULL", entidadeID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetMaxNumero retorna o maior número de documento para uma filial
func (r *documentoVendaRepository) GetMaxNumero(empresaFilialID int) (int, error) {
	var maxNumero int
	err := r.db.Model(&models.DocumentoVenda{}).
		Where("emf_id = ? AND deleted_at IS NULL", empresaFilialID).
		Select("COALESCE(MAX(ddv_numero), 0)").
		Scan(&maxNumero).Error
	if err != nil {
		return 0, err
	}
	return maxNumero, nil
}
