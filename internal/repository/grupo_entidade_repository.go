package repository

import (
	"errors"

	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
	"gorm.io/gorm"
)

// ============================================================
// TYPES
// ============================================================

// GrupoEntidadeRepository é o repositório para GrupoEntidade
type GrupoEntidadeRepository struct {
	db *gorm.DB
}

// ============================================================
// CONSTRUCTOR
// ============================================================

// NewGrupoEntidadeRepository cria uma nova instância
func NewGrupoEntidadeRepository(db *gorm.DB) *GrupoEntidadeRepository {
	return &GrupoEntidadeRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD
// ============================================================

// Create salva um novo grupo de entidade
func (r *GrupoEntidadeRepository) Create(grupo *models.GrupoEntidade) error {
	return r.db.Create(grupo).Error
}

// Update atualiza um grupo de entidade existente
func (r *GrupoEntidadeRepository) Update(id int, grupo *models.GrupoEntidade) error {
	return r.db.
		Omit("EmpresaFilial").
		Model(grupo).
		Where("gre_id = ?", id).
		Updates(grupo).Error
}

// Delete realiza exclusão lógica
func (r *GrupoEntidadeRepository) Delete(id int) error {
	grupo, err := r.FindByID(id)
	if err != nil {
		return err
	}
	if grupo.IsDeleted() {
		return errors.New("grupo de entidade já foi deletado")
	}
	grupo.SoftDelete()
	return r.Update(id, grupo)
}

// ============================================================
// MÉTODOS DE BUSCA
// ============================================================

// FindByID busca um grupo de entidade pelo ID
func (r *GrupoEntidadeRepository) FindByID(id int) (*models.GrupoEntidade, error) {
	var grupo models.GrupoEntidade
	err := r.db.Where("gre_id = ? AND deleted_at IS NULL", id).First(&grupo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("grupo de entidade não encontrado")
		}
		return nil, err
	}
	return &grupo, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de grupos com paginação e filtros
func (r *GrupoEntidadeRepository) List(limit, offset int, filters map[string]interface{}) ([]models.GrupoEntidade, int64, error) {
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
		Order("gre_nome ASC").
		Find(&grupos).Error

	if err != nil {
		return nil, 0, err
	}

	return grupos, total, nil
}

// ============================================================
// MÉTODO ADICIONAL: Verificar duplicidade
// ============================================================

// ExistsByName verifica se já existe um grupo com o nome
func (r *GrupoEntidadeRepository) ExistsByName(nome string, excludeID int) (bool, error) {
	var count int64

	query := r.db.
		Model(&models.GrupoEntidade{}).
		Where("gre_nome = ? AND deleted_at IS NULL", nome)

	if excludeID > 0 {
		query = query.Where("gre_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
