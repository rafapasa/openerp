// internal/repository/grupo_entidade_repository.go
package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// INTERFACE (Contrato)
// ============================================================

// GrupoEntidadeRepository define o contrato para operações de banco
type GrupoEntidadeRepository interface {
	// CRUD Básico
	Create(grupo *models.GrupoEntidade) error
	Update(id int, grupo *models.GrupoEntidade) error
	Delete(id int) error
	FindByID(id int) (*models.GrupoEntidade, error)
	GetByID(id int) (*models.GrupoEntidade, error)

	// Buscas Específicas
	FindByDescricao(descricao string) (*models.GrupoEntidade, error)

	// Listagem com Filtros
	List(limit, offset int, filters map[string]interface{}) ([]models.GrupoEntidade, int64, error)
	ListWithEntidades(limit, offset int, filters map[string]interface{}) ([]models.GrupoEntidade, int64, error)
	FindAll() ([]models.GrupoEntidade, error)

	// Consultas de Validação (APENAS CONSULTAS)
	ExistsByDescricao(descricao string, excludeID int) (bool, error)
	ExistsByID(id int) (bool, error)
	CountEntidadesByGrupo(grupoID int) (int64, error)
	GetDefault() (*models.GrupoEntidade, error)
	Count() (int64, error)
}

// ============================================================
// IMPLEMENTAÇÃO CONCRETA (privada)
// ============================================================

// grupoEntidadeRepository é a implementação concreta
type grupoEntidadeRepository struct {
	db *gorm.DB
}

// NewGrupoEntidadeRepository cria uma nova instância (retorna a interface)
func NewGrupoEntidadeRepository(db *gorm.DB) GrupoEntidadeRepository {
	return &grupoEntidadeRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD (APENAS PERSISTÊNCIA)
// ============================================================

// Create salva um novo grupo de entidade
func (r *grupoEntidadeRepository) Create(grupo *models.GrupoEntidade) error {
	return r.db.Create(grupo).Error
}

// Update atualiza um grupo de entidade existente
func (r *grupoEntidadeRepository) Update(id int, grupo *models.GrupoEntidade) error {
	return r.db.
		Omit("EmpresaFilial", "created_at", "deleted_at").
		Model(&models.GrupoEntidade{}).
		Where("gpe_id = ?", id).
		Updates(grupo).Error
}

// Delete realiza exclusão lógica
func (r *grupoEntidadeRepository) Delete(id int) error {
	return r.db.
		Model(&models.GrupoEntidade{}).
		Where("gpe_id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// FindByID busca um grupo de entidade pelo ID
func (r *grupoEntidadeRepository) FindByID(id int) (*models.GrupoEntidade, error) {
	var grupo models.GrupoEntidade
	err := r.db.
		Where("gpe_id = ? AND deleted_at IS NULL", id).
		First(&grupo).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("grupo de entidade com ID %d não encontrado", id))
		}
		return nil, apperrors.NewInternalError("Erro buscando grupo de entidades: ", err)
	}
	return &grupo, nil
}

// GetByID busca um grupo pelo ID (alias para FindByID)
func (r *grupoEntidadeRepository) GetByID(id int) (*models.GrupoEntidade, error) {
	return r.FindByID(id)
}

// ============================================================
// MÉTODOS DE BUSCA ESPECÍFICOS
// ============================================================

// FindByDescricao busca um grupo pela descrição
func (r *grupoEntidadeRepository) FindByDescricao(descricao string) (*models.GrupoEntidade, error) {
	var grupo models.GrupoEntidade
	err := r.db.
		Where("gpe_descricao = ? AND deleted_at IS NULL", descricao).
		First(&grupo).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("grupo de entidade com descrição %s não encontrado", descricao))
		}
		return nil, err
	}
	return &grupo, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de grupos com paginação e filtros
func (r *grupoEntidadeRepository) List(limit, offset int, filters map[string]interface{}) ([]models.GrupoEntidade, int64, error) {
	var grupos []models.GrupoEntidade
	var total int64

	query := r.db.Model(&models.GrupoEntidade{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.GrupoEntidade{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Order("gpe_descricao ASC").
		Find(&grupos).Error

	if err != nil {
		return nil, 0, err
	}

	return grupos, total, nil
}

// ListWithEntidades retorna grupos com contagem de entidades
func (r *grupoEntidadeRepository) ListWithEntidades(limit, offset int, filters map[string]interface{}) ([]models.GrupoEntidade, int64, error) {
	var grupos []models.GrupoEntidade
	var total int64

	query := r.db.Model(&models.GrupoEntidade{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.GrupoEntidade{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Busca com contagem de entidades associadas
	err := query.
		Select(`
			gpe.*,
			(
				SELECT COUNT(*) 
				FROM entidades 
				WHERE ent_grupo_id = gpe.gpe_id 
				AND deleted_at IS NULL
			) as total_entidades
		`).
		Limit(limit).
		Offset(offset).
		Order("gpe_descricao ASC").
		Find(&grupos).Error

	if err != nil {
		return nil, 0, err
	}

	return grupos, total, nil
}

// FindAll busca todos os grupos (sem paginação)
func (r *grupoEntidadeRepository) FindAll() ([]models.GrupoEntidade, error) {
	var grupos []models.GrupoEntidade
	err := r.db.
		Where("deleted_at IS NULL").
		Order("gpe_descricao ASC").
		Find(&grupos).Error
	if err != nil {
		return nil, err
	}
	return grupos, nil
}

// ============================================================
// MÉTODOS DE CONSULTA PARA VALIDAÇÕES (APENAS CONSULTAS)
// ============================================================

// ExistsByDescricao verifica se já existe um grupo com a descrição
func (r *grupoEntidadeRepository) ExistsByDescricao(descricao string, excludeID int) (bool, error) {
	var count int64
	query := r.db.
		Model(&models.GrupoEntidade{}).
		Where("gpe_descricao = ? AND deleted_at IS NULL", descricao)

	if excludeID > 0 {
		query = query.Where("gpe_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, apperrors.NewInternalError("Erro buscando Grupo de entidade. ", err)
	}

	return count > 0, nil
}

// ExistsByID verifica se um grupo existe pelo ID
func (r *grupoEntidadeRepository) ExistsByID(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.GrupoEntidade{}).
		Where("gpe_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CountEntidadesByGrupo conta quantas entidades estão associadas a um grupo
func (r *grupoEntidadeRepository) CountEntidadesByGrupo(grupoID int) (int64, error) {
	var count int64
	err := r.db.Model(&models.Entidade{}).
		Where("ent_grupo_id = ? AND deleted_at IS NULL", grupoID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetDefault busca o grupo padrão
func (r *grupoEntidadeRepository) GetDefault() (*models.GrupoEntidade, error) {
	var grupo models.GrupoEntidade
	err := r.db.
		Where("deleted_at IS NULL AND gpe_padrao = true").
		First(&grupo).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Nenhum grupo padrão encontrado
		}
		return nil, err
	}
	return &grupo, nil
}

// Count retorna o total de grupos
func (r *grupoEntidadeRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.GrupoEntidade{}).
		Where("deleted_at IS NULL").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
