package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// GrupoUsuarioRepository define o contrato para operações de persistência de GrupoUsuario.
type GrupoUsuarioRepository interface {
	Create(grupo *models.GrupoUsuario) error
	FindByID(id int) (*models.GrupoUsuario, error)
	Update(grupo *models.GrupoUsuario) error
	Delete(id int) error
	List(limit, offset int, filters map[string]interface{}) ([]models.GrupoUsuario, int64, error)
	ExistsByDescricao(descricao string, empresaFilialID, excludeID int) (bool, error)
	ExistsByID(id int) (bool, error)
	HasDependentRecords(id int) (bool, error)
}

type grupoUsuarioRepository struct {
	db *gorm.DB
}

// NewGrupoUsuarioRepository cria uma nova instância de GrupoUsuarioRepository.
func NewGrupoUsuarioRepository(db *gorm.DB) GrupoUsuarioRepository {
	return &grupoUsuarioRepository{db: db}
}

// Create insere um novo grupo de usuário no banco de dados.
func (r *grupoUsuarioRepository) Create(grupo *models.GrupoUsuario) error {
	if err := r.db.Create(grupo).Error; err != nil {
		return apperrors.NewInternalError("Erro ao criar grupo de usuário.", err)
	}
	return nil
}

// FindByID busca um grupo de usuário pelo ID.
func (r *grupoUsuarioRepository) FindByID(id int) (*models.GrupoUsuario, error) {
	var grupo models.GrupoUsuario
	err := r.db.
		Preload("EmpresaFilial"). // Carrega o relacionamento com EmpresaFilial
		Where("gpu_id = ? AND deleted_at IS NULL", id).
		First(&grupo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("Grupo de usuário com ID %d não encontrado.", id))
		}
		return nil, apperrors.NewInternalError(fmt.Sprintf("Erro ao buscar grupo de usuário com ID %d.", id), err)
	}
	return &grupo, nil
}

// Update atualiza um grupo de usuário existente no banco de dados.
func (r *grupoUsuarioRepository) Update(grupo *models.GrupoUsuario) error {
	// Omitir campos que não devem ser atualizados diretamente ou que são gerenciados por GORM
	if err := r.db.Omit("EmpresaFilial", "Usuarios", "created_at", "deleted_at").Save(grupo).Error; err != nil {
		return apperrors.NewInternalError("Erro ao atualizar grupo de usuário.", err)
	}
	return nil
}

// Delete realiza a exclusão lógica de um grupo de usuário.
func (r *grupoUsuarioRepository) Delete(id int) error {
	grupo, err := r.FindByID(id)
	if err != nil {
		return err // Retorna NotFoundError ou InternalError do FindByID
	}
	if grupo.IsDeleted() {
		return apperrors.NewConflictError(fmt.Sprintf("Grupo de usuário com ID %d já está excluído.", id))
	}

	grupo.SoftDelete()
	if err := r.db.Save(grupo).Error; err != nil {
		return apperrors.NewInternalError("Erro ao excluir grupo de usuário.", err)
	}
	return nil
}

// List lista grupos de usuário com paginação e filtros.
func (r *grupoUsuarioRepository) List(limit, offset int, filters map[string]interface{}) ([]models.GrupoUsuario, int64, error) {
	var grupos []models.GrupoUsuario
	var total int64

	query := r.db.Model(&models.GrupoUsuario{}).Where("deleted_at IS NULL")

	// Aplica filtros dinâmicos
	query = utils.ApplyFilters(query, models.GrupoUsuario{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("Erro ao contar grupos de usuário.", err)
	}

	if err := query.
		Preload("EmpresaFilial"). // Carrega o relacionamento com EmpresaFilial
		Limit(limit).
		Offset(offset).
		Find(&grupos).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("Erro ao listar grupos de usuário.", err)
	}

	return grupos, total, nil
}

// ExistsByDescricao verifica se um grupo de usuário com a descrição e filial especificadas já existe.
func (r *grupoUsuarioRepository) ExistsByDescricao(descricao string, empresaFilialID, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.GrupoUsuario{}).
		Where("gpu_descricao = ? AND emf_id = ? AND deleted_at IS NULL", descricao, empresaFilialID)
	if excludeID > 0 {
		query = query.Where("gpu_id <> ?", excludeID)
	}
	err := query.Count(&count).Error
	if err != nil {
		return false, apperrors.NewInternalError("Erro ao verificar descrição de grupo de usuário existente.", err)
	}
	return count > 0, nil
}

// ExistsByID verifica se um grupo de usuário existe pelo ID.
func (r *grupoUsuarioRepository) ExistsByID(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.GrupoUsuario{}).
		Where("gpu_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	if err != nil {
		return false, apperrors.NewInternalError("Erro ao verificar existência de grupo de usuário por ID.", err)
	}
	return count > 0, nil
}

// HasDependentRecords verifica se o grupo de usuário tem registros dependentes (usuários).
func (r *grupoUsuarioRepository) HasDependentRecords(id int) (bool, error) {
	var count int64
	// Verifica se existem usuários associados a este grupo
	err := r.db.Model(&models.Usuario{}).
		Where("gpu_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	if err != nil {
		return false, apperrors.NewInternalError("Erro ao verificar registros dependentes para grupo de usuário.", err)
	}
	return count > 0, nil
}
