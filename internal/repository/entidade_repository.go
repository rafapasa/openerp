package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// TYPES
// ============================================================

// EntidadeRepository é o repositório para Entidade
type EntidadeRepository struct {
	db *gorm.DB
}

// ============================================================
// CONSTRUCTOR
// ============================================================

// NewEntidadeRepository cria uma nova instância
func NewEntidadeRepository(db *gorm.DB) *EntidadeRepository {
	return &EntidadeRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD
// ============================================================

// Create salva uma nova entidade
func (r *EntidadeRepository) Create(entidade *models.Entidade) error {
	return r.db.Create(entidade).Error
}

// Update atualiza uma entidade existente
func (r *EntidadeRepository) Update(entidade *models.Entidade) error {
	return r.db.Save(entidade).Error
}

// Delete realiza exclusão lógica de uma entidade pelo ID (CORRIGIDO)
func (r *EntidadeRepository) Delete(id int) error {
	// 1. Buscar a entidade
	entidade, err := r.FindByID(id)
	if err != nil {
		return err
	}

	// 2. Verificar se já foi deletada
	if entidade.IsDeleted() {
		return errors.New("entidade já foi deletada")
	}

	// 3. Realizar soft delete
	entidade.SoftDelete()

	// 4. Salvar
	return r.db.Save(entidade).Error
}

// ============================================================
// MÉTODOS DE BUSCA
// ============================================================

// FindByID busca uma entidade pelo ID com relacionamentos
func (r *EntidadeRepository) FindByID(id int) (*models.Entidade, error) {
	var entidade models.Entidade
	err := r.db.
		Preload("GrupoEntidade").
		Preload("EmpresaFilial").
		Preload("TabelaPreco").
		Preload("TabelaDesconto").
		Preload("Horario").
		Preload("Enderecos").
		Preload("Contatos").
		Preload("Documentos").
		Preload("Regimes").
		Where("ent_id = ? AND deleted_at IS NULL", id).
		First(&entidade).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("entidade não encontrada")
		}
		return nil, err
	}
	return &entidade, nil
}

// FindByDocumento busca uma entidade pelo CPF/CNPJ
func (r *EntidadeRepository) FindByDocumento(documento string) (*models.Entidade, error) {
	var entidade models.Entidade
	err := r.db.
		Where("ent_inscricaofederal = ? AND deleted_at IS NULL", documento).
		First(&entidade).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("entidade não encontrada")
		}
		return nil, err
	}
	return &entidade, nil
}

// FindByRazaoSocial busca entidades pela razão social (autocomplete)
func (r *EntidadeRepository) FindByRazaoSocial(nome string, limit int) ([]models.Entidade, error) {
	var entidades []models.Entidade
	err := r.db.
		Where("ent_razaosocial LIKE ? AND deleted_at IS NULL", "%"+nome+"%").
		Limit(limit).
		Find(&entidades).Error

	if err != nil {
		return nil, err
	}
	return entidades, nil
}

// FindByNomeFantasia busca entidades pelo nome fantasia (autocomplete)
func (r *EntidadeRepository) FindByNomeFantasia(nome string, limit int) ([]models.Entidade, error) {
	var entidades []models.Entidade
	err := r.db.
		Where("ent_nomefantasia LIKE ? AND deleted_at IS NULL", "%"+nome+"%").
		Limit(limit).
		Find(&entidades).Error

	if err != nil {
		return nil, err
	}
	return entidades, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de entidades com paginação e filtros
func (r *EntidadeRepository) List(limit, offset int, filters map[string]interface{}) ([]models.Entidade, int64, error) {
	var entidades []models.Entidade
	var total int64

	// Construir a query base
	query := r.db.Model(&models.Entidade{}).Where("deleted_at IS NULL")

	// Aplicar filtros dinamicamente
	query = utils.ApplyFilters(query, models.Entidade{}, filters)

	// Contar total de registros
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Buscar com paginação e relacionamentos
	err := query.
		Preload("GrupoEntidade").
		Preload("EmpresaFilial").
		Preload("TabelaPreco").
		Preload("TabelaDesconto").
		Preload("Horario").
		Preload("Enderecos").
		Preload("Contatos").
		Limit(limit).
		Offset(offset).
		Order("ent_id DESC").
		Find(&entidades).Error

	if err != nil {
		return nil, 0, err
	}

	return entidades, total, nil
}

// ============================================================
// MÉTODO ADICIONAL: Verificar duplicidade
// ============================================================

// ExistsByDocumento verifica se já existe uma entidade com o documento
func (r *EntidadeRepository) ExistsByDocumento(documento string, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.Entidade{}).
		Where("ent_inscricaofederal = ? AND deleted_at IS NULL", documento)

	// Se for para excluir um ID (atualização)
	if excludeID > 0 {
		query = query.Where("ent_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
