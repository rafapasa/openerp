package repository

import (
	"errors"
	"fmt"

	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
	"gorm.io/gorm"
)

// ProdutoTamanhoRepository define o contrato para operações de persistência de ProdutoTamanho.
type ProdutoTamanhoRepository interface {
	Create(tamanho *models.ProdutoTamanho) error
	FindByID(id int) (*models.ProdutoTamanho, error)
	Update(tamanho *models.ProdutoTamanho) error
	Delete(id int) error
	List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoTamanho, int64, error)
	ExistsBySigla(sigla string, empresaFilialID, excludeID int) (bool, error)
	ExistsByNome(nome string, empresaFilialID, excludeID int) (bool, error)
}

type produtoTamanhoRepository struct {
	db *gorm.DB
}

// NewProdutoTamanhoRepository cria uma nova instância de ProdutoTamanhoRepository.
func NewProdutoTamanhoRepository(db *gorm.DB) ProdutoTamanhoRepository {
	return &produtoTamanhoRepository{db: db}
}

// Create insere um novo tamanho de produto no banco de dados.
func (r *produtoTamanhoRepository) Create(tamanho *models.ProdutoTamanho) error {
	if err := r.db.Create(tamanho).Error; err != nil {
		return apperrors.NewInternalError("Erro ao criar tamanho de produto.", err)
	}
	return nil
}

// FindByID busca um tamanho de produto pelo ID.
func (r *produtoTamanhoRepository) FindByID(id int) (*models.ProdutoTamanho, error) {
	var tamanho models.ProdutoTamanho
	if err := r.db.Where("ptam_id = ?", id).First(&tamanho).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("Tamanho de produto com ID %d não encontrado.", id))
		}
		return nil, apperrors.NewInternalError(fmt.Sprintf("Erro ao buscar tamanho de produto com ID %d.", id), err)
	}
	return &tamanho, nil
}

// Update atualiza um tamanho de produto existente no banco de dados.
func (r *produtoTamanhoRepository) Update(tamanho *models.ProdutoTamanho) error {
	if err := r.db.Save(tamanho).Error; err != nil {
		return apperrors.NewInternalError("Erro ao atualizar tamanho de produto.", err)
	}
	return nil
}

// Delete realiza a exclusão lógica de um tamanho de produto.
func (r *produtoTamanhoRepository) Delete(id int) error {
	tamanho, err := r.FindByID(id)
	if err != nil {
		return err // Retorna NotFoundError ou InternalError do FindByID
	}
	if tamanho.IsDeleted() {
		return apperrors.NewConflictError(fmt.Sprintf("Tamanho de produto com ID %d já está excluído.", id))
	}

	tamanho.SoftDelete()
	if err := r.db.Save(tamanho).Error; err != nil {
		return apperrors.NewInternalError("Erro ao excluir tamanho de produto.", err)
	}
	return nil
}

// List lista tamanhos de produto com paginação e filtros.
func (r *produtoTamanhoRepository) List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoTamanho, int64, error) {
	var tamanhos []models.ProdutoTamanho
	var total int64

	query := r.db.Model(&models.ProdutoTamanho{}).Where("deleted_at IS NULL")

	// Aplica filtros dinâmicos
	query = utils.ApplyFilters(query, models.ProdutoTamanho{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("Erro ao contar tamanhos de produto.", err)
	}

	if err := query.Limit(limit).Offset(offset).Find(&tamanhos).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("Erro ao listar tamanhos de produto.", err)
	}

	return tamanhos, total, nil
}

// ExistsBySigla verifica se um tamanho com a sigla e filial especificadas já existe.
func (r *produtoTamanhoRepository) ExistsBySigla(sigla string, empresaFilialID, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.ProdutoTamanho{}).Where("ptam_sigla = ? AND emf_id = ? AND deleted_at IS NULL", sigla, empresaFilialID)
	if excludeID > 0 {
		query = query.Where("ptam_id <> ?", excludeID)
	}
	err := query.Count(&count).Error
	if err != nil {
		return false, apperrors.NewInternalError("Erro ao verificar sigla de tamanho existente.", err)
	}
	return count > 0, nil
}

// ExistsByNome verifica se um tamanho com o nome e filial especificadas já existe.
func (r *produtoTamanhoRepository) ExistsByNome(nome string, empresaFilialID, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.ProdutoTamanho{}).Where("ptam_nome = ? AND emf_id = ? AND deleted_at IS NULL", nome, empresaFilialID)
	if excludeID > 0 {
		query = query.Where("ptam_id <> ?", excludeID)
	}
	err := query.Count(&count).Error
	if err != nil {
		return false, apperrors.NewInternalError("Erro ao verificar nome de tamanho existente.", err)
	}
	return count > 0, nil
}
