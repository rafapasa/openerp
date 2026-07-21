package repository

import (
	"errors"

	"fmt"

	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
	"gorm.io/gorm"
)

// ============================================================
// TYPES
// ============================================================

// CondicaoPagamentoRepository é o repositório para CondicaoPagamento
type CondicaoPagamentoRepository struct {
	db *gorm.DB
}

// ============================================================
// CONSTRUCTOR
// ============================================================

// NewCondicaoPagamentoRepository cria uma nova instância
func NewCondicaoPagamentoRepository(db *gorm.DB) *CondicaoPagamentoRepository {
	return &CondicaoPagamentoRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD
// ============================================================

// Create salva uma nova CondicaoPagamento
func (r *CondicaoPagamentoRepository) Create(condicaoPagamento *models.CondicaoPagamento) error {
	return r.db.Create(condicaoPagamento).Error
}

// Update atualiza uma CondicaoPagamento existente
func (r *CondicaoPagamentoRepository) Update(id int, condicaoPagamento *models.CondicaoPagamento) error {
	return r.db.
		Omit("Tipo_Documento", "Forma_Pagamento", "Portador").
		Model(&models.CondicaoPagamento{}).
		Where("cdpgt_id = ?", id).
		Updates(condicaoPagamento).Error
}

// Delete realiza exclusão lógica de uma CondicaoPagamento pelo ID
func (r *CondicaoPagamentoRepository) Delete(id int) error {
	condicaoPagamento, err := r.FindByID(id)
	if err != nil {
		return err
	}
	if condicaoPagamento.IsDeleted() {
		return errors.New("CondicaoPagamento já foi deletada")
	}
	condicaoPagamento.SoftDelete()
	return r.Update(id, condicaoPagamento)
}

func (r *CondicaoPagamentoRepository) GetByID(id int) (*models.CondicaoPagamento, error) {
	var condicaoPagamento models.CondicaoPagamento
	result := r.db.Where("cdpgt_id = ?", id).First(&condicaoPagamento)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("CondicaoPagamento com ID %d não encontrada", id))
		}
		return nil, apperrors.NewInternalError("Erro buscando CondicaoPagamento", result.Error)
	}
	return &condicaoPagamento, nil
}

// ============================================================
// MÉTODOS DE BUSCA
// ============================================================

// FindByID busca uma CondicaoPagamento pelo ID com relacionamentos
func (r *CondicaoPagamentoRepository) FindByID(id int) (*models.CondicaoPagamento, error) {
	var condicaoPagamento models.CondicaoPagamento
	err := r.db.
		Preload("TipoDocumento").
		Preload("FormaPagamento").
		Preload("Portador").
		Where("cdpgt_id = ? AND deleted_at IS NULL", id).
		First(&condicaoPagamento).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("CondicaoPagamento com ID %d não encontrada", id))
		}
		return nil, err
	}
	return &condicaoPagamento, nil
}

// FindByTipoDocumento busca uma CondicaoPagamento pelo TipoDocumento
func (r *CondicaoPagamentoRepository) FindByTipoDocumento(id string) (*models.CondicaoPagamento, error) {
	var condicaoPagamento models.CondicaoPagamento
	err := r.db.
		Where("tdoc_id = ? AND deleted_at IS NULL", id).
		First(&condicaoPagamento).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("CondicaoPagamento não encontrada")
		}
		return nil, err
	}
	return &condicaoPagamento, nil
}

// FindByNome busca CondicaoPagamento pela Nome (autocomplete)
func (r *CondicaoPagamentoRepository) FindByNome(nome string, limit int) ([]models.CondicaoPagamento, error) {
	var Condicoes []models.CondicaoPagamento
	err := r.db.
		Where("cdpgt_descricao LIKE ? AND deleted_at IS NULL", "%"+nome+"%").
		Limit(limit).
		Find(&Condicoes).Error
	if err != nil {
		return nil, err
	}
	return Condicoes, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de CondicaoPagamento com paginação e filtros
func (r *CondicaoPagamentoRepository) List(limit, offset int, filters map[string]interface{}) ([]models.CondicaoPagamento, int64, error) {
	var Condicoes []models.CondicaoPagamento
	var total int64

	query := r.db.Model(&models.CondicaoPagamento{}).Where("deleted_at IS NULL")

	query = utils.ApplyFilters(query, models.CondicaoPagamento{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("TipoDocumento").
		Preload("FormaPagamento").
		Preload("Portador").
		Limit(limit).
		Offset(offset).
		Order("codpgt_id DESC").
		Find(&Condicoes).Error

	if err != nil {
		return nil, 0, err
	}

	return Condicoes, total, nil
}

// ============================================================
// MÉTODO ADICIONAL: Verificar duplicidade
// ============================================================

// ExistsByTipoDocumento verifica se já existe uma CondicaoPagamento com o TipoDocumento
func (r *CondicaoPagamentoRepository) ExistsByTipoDocumento(id string, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.CondicaoPagamento{}).Where("tdoc_id = ? AND deleted_at IS NULL", utils.LimparDocumento(id))
	if excludeID > 0 {
		query = query.Where("codpgt_id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
