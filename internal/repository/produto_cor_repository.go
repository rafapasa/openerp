package repository

import (
	"errors"
	"fmt"

	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
	"gorm.io/gorm"
)

// ProdutoCorRepository define o contrato para operações de persistência de ProdutoCor.
type ProdutoCorRepository interface {
	Create(cor *models.ProdutoCor) error
	FindByID(id int) (*models.ProdutoCor, error)
	Update(cor *models.ProdutoCor) error
	Delete(id int) error
	List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoCor, int64, error)
	ExistsBySigla(sigla string, empresaFilialID, excludeID int) (bool, error)
	ExistsByNome(nome string, empresaFilialID, excludeID int) (bool, error)
}

type produtoCorRepository struct {
	db *gorm.DB
}

// NewProdutoCorRepository cria uma nova instância de ProdutoCorRepository.
func NewProdutoCorRepository(db *gorm.DB) ProdutoCorRepository {
	return &produtoCorRepository{db: db}
}

// Create insere uma nova cor de produto no banco de dados.
func (r *produtoCorRepository) Create(cor *models.ProdutoCor) error {
	if err := r.db.Create(cor).Error; err != nil {
		return apperrors.NewInternalError("Erro ao criar cor de produto.", err)
	}
	return nil
}

// FindByID busca uma cor de produto pelo ID.
func (r *produtoCorRepository) FindByID(id int) (*models.ProdutoCor, error) {
	var cor models.ProdutoCor
	if err := r.db.Where("cor_id = ?", id).First(&cor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("Cor de produto com ID %d não encontrada.", id))
		}
		return nil, apperrors.NewInternalError(fmt.Sprintf("Erro ao buscar cor de produto com ID %d.", id), err)
	}
	return &cor, nil
}

// Update atualiza uma cor de produto existente no banco de dados.
func (r *produtoCorRepository) Update(cor *models.ProdutoCor) error {
	if err := r.db.Save(cor).Error; err != nil {
		return apperrors.NewInternalError("Erro ao atualizar cor de produto.", err)
	}
	return nil
}

// Delete realiza a exclusão lógica de uma cor de produto.
func (r *produtoCorRepository) Delete(id int) error {
	cor, err := r.FindByID(id)
	if err != nil {
		return err // Retorna NotFoundError ou InternalError do FindByID
	}
	if cor.IsDeleted() {
		return apperrors.NewConflictError(fmt.Sprintf("Cor de produto com ID %d já está excluída.", id))
	}

	cor.SoftDelete()
	if err := r.db.Save(cor).Error; err != nil {
		return apperrors.NewInternalError("Erro ao excluir cor de produto.", err)
	}
	return nil
}

// List lista cores de produto com paginação e filtros.
func (r *produtoCorRepository) List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoCor, int64, error) {
	var cores []models.ProdutoCor
	var total int64

	query := r.db.Model(&models.ProdutoCor{}).Where("deleted_at IS NULL")

	// Aplica filtros dinâmicos
	query = utils.ApplyFilters(query, models.ProdutoCor{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("Erro ao contar cores de produto.", err)
	}

	if err := query.Limit(limit).Offset(offset).Find(&cores).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("Erro ao listar cores de produto.", err)
	}

	return cores, total, nil
}

// ExistsBySigla verifica se uma cor com a sigla e filial especificadas já existe.
func (r *produtoCorRepository) ExistsBySigla(sigla string, empresaFilialID, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.ProdutoCor{}).Where("cor_sigla = ? AND emf_id = ? AND deleted_at IS NULL", sigla, empresaFilialID)
	if excludeID > 0 {
		query = query.Where("cor_id <> ?", excludeID)
	}
	err := query.Count(&count).Error
	if err != nil {
		return false, apperrors.NewInternalError("Erro ao verificar sigla de cor existente.", err)
	}
	return count > 0, nil
}

// ExistsByNome verifica se uma cor com o nome e filial especificadas já existe.
func (r *produtoCorRepository) ExistsByNome(nome string, empresaFilialID, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.ProdutoCor{}).Where("cor_nome = ? AND emf_id = ? AND deleted_at IS NULL", nome, empresaFilialID)
	if excludeID > 0 {
		query = query.Where("cor_id <> ?", excludeID)
	}
	err := query.Count(&count).Error
	if err != nil {
		return false, apperrors.NewInternalError("Erro ao verificar nome de cor existente.", err)
	}
	return count > 0, nil
}
