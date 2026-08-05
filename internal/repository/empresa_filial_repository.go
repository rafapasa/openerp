package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// EmpresaFilialRepository define o contrato para operações de persistência de EmpresaFilial.
type EmpresaFilialRepository interface {
	Create(filial *models.EmpresaFilial) error
	FindByID(id int) (*models.EmpresaFilial, error)
	Update(filial *models.EmpresaFilial) error
	Delete(id int) error
	List(limit, offset int, filters map[string]interface{}) ([]models.EmpresaFilial, int64, error)
	ExistsByNumero(numero, empresaID, excludeID int) (bool, error)
	ExistsByID(id int) (bool, error)
	HasDependentRecords(id int) (bool, error)
}

type empresaFilialRepository struct {
	db *gorm.DB
}

// NewEmpresaFilialRepository cria uma nova instância de EmpresaFilialRepository.
func NewEmpresaFilialRepository(db *gorm.DB) EmpresaFilialRepository {
	return &empresaFilialRepository{db: db}
}

// Create insere uma nova filial de empresa no banco de dados.
func (r *empresaFilialRepository) Create(filial *models.EmpresaFilial) error {
	if err := r.db.Create(filial).Error; err != nil {
		return apperrors.NewInternalError("Erro ao criar filial de empresa.", err)
	}
	return nil
}

// FindByID busca uma filial de empresa pelo ID.
func (r *empresaFilialRepository) FindByID(id int) (*models.EmpresaFilial, error) {
	var filial models.EmpresaFilial
	err := r.db.
		Preload("Empresa"). // Carrega o relacionamento com Empresa
		Where("emf_id = ? AND deleted_at IS NULL", id).
		First(&filial).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("Filial de empresa com ID %d não encontrada.", id))
		}
		return nil, apperrors.NewInternalError(fmt.Sprintf("Erro ao buscar filial de empresa com ID %d.", id), err)
	}
	return &filial, nil
}

// Update atualiza uma filial de empresa existente no banco de dados.
func (r *empresaFilialRepository) Update(filial *models.EmpresaFilial) error {
	// Omitir campos que não devem ser atualizados diretamente ou que são gerenciados por GORM
	if err := r.db.Omit("Empresa", "created_at", "deleted_at").Save(filial).Error; err != nil {
		return apperrors.NewInternalError("Erro ao atualizar filial de empresa.", err)
	}
	return nil
}

// Delete realiza a exclusão lógica de uma filial de empresa.
func (r *empresaFilialRepository) Delete(id int) error {
	filial, err := r.FindByID(id)
	if err != nil {
		return err // Retorna NotFoundError ou InternalError do FindByID
	}
	if filial.IsDeleted() {
		return apperrors.NewConflictError(fmt.Sprintf("Filial de empresa com ID %d já está excluída.", id))
	}

	filial.SoftDelete()
	if err := r.db.Save(filial).Error; err != nil {
		return apperrors.NewInternalError("Erro ao excluir filial de empresa.", err)
	}
	return nil
}

// List lista filiais de empresa com paginação e filtros.
func (r *empresaFilialRepository) List(limit, offset int, filters map[string]interface{}) ([]models.EmpresaFilial, int64, error) {
	var filiais []models.EmpresaFilial
	var total int64

	query := r.db.Model(&models.EmpresaFilial{}).Where("deleted_at IS NULL")

	// Aplica filtros dinâmicos
	query = utils.ApplyFilters(query, models.EmpresaFilial{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("Erro ao contar filiais de empresa.", err)
	}

	if err := query.
		Preload("Empresa"). // Carrega o relacionamento com Empresa
		Limit(limit).
		Offset(offset).
		Find(&filiais).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("Erro ao listar filiais de empresa.", err)
	}

	return filiais, total, nil
}

// ExistsByNumero verifica se uma filial de empresa com o número e empresa especificados já existe.
func (r *empresaFilialRepository) ExistsByNumero(numero, empresaID, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.EmpresaFilial{}).
		Where("emf_numero = ? AND emp_id = ? AND deleted_at IS NULL", numero, empresaID)
	if excludeID > 0 {
		query = query.Where("emf_id <> ?", excludeID)
	}
	err := query.Count(&count).Error
	if err != nil {
		return false, apperrors.NewInternalError("Erro ao verificar número de filial de empresa existente.", err)
	}
	return count > 0, nil
}

// ExistsByID verifica se uma filial de empresa existe pelo ID.
func (r *empresaFilialRepository) ExistsByID(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.EmpresaFilial{}).
		Where("emf_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	if err != nil {
		return false, apperrors.NewInternalError("Erro ao verificar existência de filial de empresa por ID.", err)
	}
	return count > 0, nil
}

// HasDependentRecords verifica se a filial de empresa tem registros dependentes.
func (r *empresaFilialRepository) HasDependentRecords(id int) (bool, error) {
	var count int64
	// Verifica se existem grupos de usuário associados a esta filial
	err := r.db.Model(&models.GrupoUsuario{}).
		Where("emf_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	if err != nil {
		return false, apperrors.NewInternalError("Erro ao verificar dependências de GrupoUsuario para filial de empresa.", err)
	}
	if count > 0 {
		return true, nil
	}

	// Verifica se existem produtos associados a esta filial
	err = r.db.Model(&models.Produto{}).
		Where("emf_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	if err != nil {
		return false, apperrors.NewInternalError("Erro ao verificar dependências de Produto para filial de empresa.", err)
	}
	if count > 0 {
		return true, nil
	}

	// Adicione outras verificações de dependência conforme necessário (e.g., Entidades, Documentos de Venda, etc.)

	return false, nil
}