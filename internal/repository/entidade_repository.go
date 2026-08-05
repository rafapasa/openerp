// internal/repository/entidade_repository.go
package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// INTERFACE (Contrato)
// ============================================================

// EntidadeRepository define o contrato para operações de banco
type EntidadeRepository interface {
	// CRUD Básico
	Create(entidade *models.Entidade) error
	Update(id int, entidade *models.Entidade) error
	Delete(id int) error
	GetByID(id int) (*models.Entidade, error)
	FindByID(id int) (*models.Entidade, error)

	// Buscas Específicas
	FindByDocumento(documento string) (*models.Entidade, error)
	FindByRazaoSocial(nome string, limit int) ([]models.Entidade, error)
	FindByNomeFantasia(nome string, limit int) ([]models.Entidade, error)
	FindByGrupoID(grupoID int) ([]models.Entidade, error)
	FindByFilialID(filialID int) ([]models.Entidade, error)

	// Listagem com Filtros
	List(limit, offset int, filters map[string]interface{}) ([]models.Entidade, int64, error)
	ListWithFullPreload(limit, offset int, filters map[string]interface{}) ([]models.Entidade, int64, error)

	// Consultas de Validação (APENAS CONSULTAS)
	ExistsByDocumento(documento string, excludeID int) (bool, error)
	ExistsByRazaoSocial(razaoSocial string, excludeID int) (bool, error)
	ExistsByNomeFantasia(nomeFantasia string, excludeID int) (bool, error)
	ExistsByID(id int) (bool, error)

	// Operações em Lote
	BulkUpdateStatus(ids []int, status string) error
	BulkDelete(ids []int) error

	// Consultas de Dependências
	HasDependentRecords(id int) (bool, error)
	CountDependentRecords(id int) (map[string]int64, error)
}

// ============================================================
// IMPLEMENTAÇÃO CONCRETA (privada)
// ============================================================

// entidadeRepository é a implementação concreta
type entidadeRepository struct {
	db *gorm.DB
}

// NewEntidadeRepository cria uma nova instância (retorna a interface)
func NewEntidadeRepository(db *gorm.DB) EntidadeRepository {
	return &entidadeRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD (APENAS PERSISTÊNCIA)
// ============================================================

// Create salva uma nova entidade
func (r *entidadeRepository) Create(entidade *models.Entidade) error {
	err := r.db.Create(entidade).Error
	if err != nil {
		return apperrors.NewInternalError("Erro ao criar entidade.", err)
	}
	return nil
}

// Update atualiza uma entidade existente
func (r *entidadeRepository) Update(id int, entidade *models.Entidade) error {
	err := r.db.
		Omit("EmpresaFilial", "GrupoEntidade", "TabelaPreco", "TabelaDesconto", "Horario", "created_at", "deleted_at").
		Model(&models.Entidade{}).
		Where("ent_id = ?", id).
		Updates(entidade).Error
	if err != nil {
		return apperrors.NewInternalError("Erro ao atualizar entidade.", err)
	}
	return nil
}

// Delete realiza exclusão lógica de uma entidade pelo ID
func (r *entidadeRepository) Delete(id int) error {
	err := r.db.
		Model(&models.Entidade{}).
		Where("ent_id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
	if err != nil {
		return apperrors.NewInternalError("Erro ao excluir entidade.", err)
	}
	return nil
}

// GetByID busca uma entidade pelo ID (sem relacionamentos, mais leve)
func (r *entidadeRepository) GetByID(id int) (*models.Entidade, error) {
	var entidade models.Entidade
	result := r.db.
		Where("ent_id = ? AND deleted_at IS NULL", id).
		First(&entidade)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("Entidade com ID %d não encontrada", id))
		}
		return nil, apperrors.NewInternalError("Erro buscando entidade", result.Error)
	}
	return &entidade, nil
}

// FindByID busca uma entidade pelo ID com todos os relacionamentos
func (r *entidadeRepository) FindByID(id int) (*models.Entidade, error) {
	var entidade models.Entidade
	err := r.db.
		Preload("GrupoEntidade").
		Preload("EmpresaFilial").
		Preload("TabelaPreco").
		Preload("TabelaDesconto").
		Preload("Horario").
		Preload("Enderecos", "deleted_at IS NULL").
		Preload("Contatos", "deleted_at IS NULL").
		Preload("Documentos", "deleted_at IS NULL").
		Preload("Regimes", "deleted_at IS NULL").
		Where("ent_id = ? AND deleted_at IS NULL", id).
		First(&entidade).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("entidade com ID %d não encontrada", id))
		}
		return nil, apperrors.NewInternalError("Erro buscando entidade.", err)
	}
	return &entidade, nil
}

// ============================================================
// MÉTODOS DE BUSCA ESPECÍFICOS
// ============================================================

// FindByDocumento busca uma entidade pelo CPF/CNPJ
func (r *entidadeRepository) FindByDocumento(documento string) (*models.Entidade, error) {
	documentoLimpo := utils.LimparDocumento(documento)

	var entidade models.Entidade
	err := r.db.
		Preload("GrupoEntidade").
		Preload("EmpresaFilial").
		Preload("Enderecos", "deleted_at IS NULL").
		Where("ent_inscricaofederal = ? AND deleted_at IS NULL", documentoLimpo).
		First(&entidade).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("entidade com documento %s não encontrada", documento))
		}
		return nil, apperrors.NewInternalError("Erro buscando entidade.", err)
	}
	return &entidade, nil
}

// FindByRazaoSocial busca entidades pela razão social (autocomplete)
func (r *entidadeRepository) FindByRazaoSocial(nome string, limit int) ([]models.Entidade, error) {
	var entidades []models.Entidade
	err := r.db.
		Where("ent_razaosocial LIKE ? AND deleted_at IS NULL", "%"+nome+"%").
		Order("ent_razaosocial ASC").
		Limit(limit).
		Find(&entidades).Error

	if err != nil {
		return nil, apperrors.NewInternalError("Erro buscando entidade.", err)
	}
	return entidades, nil
}

// FindByNomeFantasia busca entidades pelo nome fantasia (autocomplete)
func (r *entidadeRepository) FindByNomeFantasia(nome string, limit int) ([]models.Entidade, error) {
	var entidades []models.Entidade
	err := r.db.
		Where("ent_nomefantasia LIKE ? AND deleted_at IS NULL", "%"+nome+"%").
		Order("ent_nomefantasia ASC").
		Limit(limit).
		Find(&entidades).Error

	if err != nil {
		return nil, apperrors.NewInternalError("Erro buscando entidade.", err)
	}
	return entidades, nil
}

// FindByGrupoID busca entidades por grupo
func (r *entidadeRepository) FindByGrupoID(grupoID int) ([]models.Entidade, error) {
	var entidades []models.Entidade
	err := r.db.
		Where("gpe_id = ? AND deleted_at IS NULL", grupoID).
		Order("ent_razaosocial ASC").
		Find(&entidades).Error
	if err != nil {
		return nil, apperrors.NewInternalError("Erro buscando entidade.", err)
	}
	return entidades, nil
}

// FindByFilialID busca entidades por filial
func (r *entidadeRepository) FindByFilialID(filialID int) ([]models.Entidade, error) {
	var entidades []models.Entidade
	err := r.db.
		Where("enf_id = ? AND deleted_at IS NULL", filialID).
		Order("ent_razaosocial ASC").
		Find(&entidades).Error
	if err != nil {
		return nil, apperrors.NewInternalError("Erro buscando entidade.", err)
	}
	return entidades, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de entidades com paginação e filtros (básico)
func (r *entidadeRepository) List(limit, offset int, filters map[string]interface{}) ([]models.Entidade, int64, error) {
	var entidades []models.Entidade
	var total int64

	query := r.db.Model(&models.Entidade{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.Entidade{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("Erro buscando entidade.", err)
	}

	err := query.
		Preload("GrupoEntidade").
		Preload("EmpresaFilial").
		Preload("TabelaPreco").
		Preload("Horario").
		Limit(limit).
		Offset(offset).
		Order("ent_id DESC").
		Find(&entidades).Error

	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro buscando entidade.", err)
	}

	return entidades, total, nil
}

// ListWithFullPreload retorna uma lista com todos os relacionamentos
func (r *entidadeRepository) ListWithFullPreload(limit, offset int, filters map[string]interface{}) ([]models.Entidade, int64, error) {
	var entidades []models.Entidade
	var total int64

	query := r.db.Model(&models.Entidade{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.Entidade{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("Erro buscando entidade.", err)
	}

	err := query.
		Preload("GrupoEntidade").
		Preload("EmpresaFilial").
		Preload("TabelaPreco").
		Preload("TabelaDesconto").
		Preload("Horario").
		Preload("Enderecos", "deleted_at IS NULL").
		Preload("Contatos", "deleted_at IS NULL").
		Preload("Documentos", "deleted_at IS NULL").
		Preload("Regimes", "deleted_at IS NULL").
		Limit(limit).
		Offset(offset).
		Order("ent_id DESC").
		Find(&entidades).Error

	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro buscando entidade.", err)
	}

	return entidades, total, nil
}

// ============================================================
// MÉTODOS DE CONSULTA PARA VALIDAÇÕES (APENAS CONSULTAS)
// ============================================================

// ExistsByDocumento verifica se já existe uma entidade com o documento
func (r *entidadeRepository) ExistsByDocumento(documento string, excludeID int) (bool, error) {
	var count int64
	documentoLimpo := utils.LimparDocumento(documento)

	query := r.db.
		Model(&models.Entidade{}).
		Where("ent_inscricaofederal = ? AND deleted_at IS NULL", documentoLimpo)

	if excludeID > 0 {
		query = query.Where("ent_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, apperrors.NewInternalError("Erro buscando entidade.", err)
	}

	return count > 0, nil
}

// ExistsByRazaoSocial verifica se já existe uma entidade com a razão social
func (r *entidadeRepository) ExistsByRazaoSocial(razaoSocial string, excludeID int) (bool, error) {
	var count int64
	query := r.db.
		Model(&models.Entidade{}).
		Where("ent_razaosocial = ? AND deleted_at IS NULL", razaoSocial)

	if excludeID > 0 {
		query = query.Where("ent_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, apperrors.NewInternalError("Erro buscando entidade.", err)
	}

	return count > 0, nil
}

// ExistsByNomeFantasia verifica se já existe uma entidade com o nome fantasia
func (r *entidadeRepository) ExistsByNomeFantasia(nomeFantasia string, excludeID int) (bool, error) {
	if nomeFantasia == "" {
		return false, nil // Nome fantasia pode ser vazio
	}

	var count int64
	query := r.db.
		Model(&models.Entidade{}).
		Where("ent_nomefantasia = ? AND deleted_at IS NULL", nomeFantasia)

	if excludeID > 0 {
		query = query.Where("ent_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, apperrors.NewInternalError("Erro buscando entidade.", err)
	}

	return count > 0, nil
}

// ExistsByID verifica se uma entidade existe pelo ID
func (r *entidadeRepository) ExistsByID(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.Entidade{}).
		Where("ent_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	if err != nil {
		return false, apperrors.NewInternalError("Erro buscando entidade.", err)
	}

	return count > 0, nil
}

// ============================================================
// OPERAÇÕES EM LOTE
// ============================================================

// BulkUpdateStatus atualiza o status de múltiplas entidades
func (r *entidadeRepository) BulkUpdateStatus(ids []int, status string) error {
	err := r.db.Model(&models.Entidade{}).
		Where("ent_id IN ? AND deleted_at IS NULL", ids).
		Update("ent_status", status).Error
	if err != nil {
		return apperrors.NewInternalError("Erro atualizando status das entidades.", err)
	}
	return nil
}

// BulkDelete realiza exclusão lógica de múltiplas entidades
func (r *entidadeRepository) BulkDelete(ids []int) error {
	err := r.db.Model(&models.Entidade{}).
		Where("ent_id IN ? AND deleted_at IS NULL", ids).
		Update("deleted_at", gorm.Expr("NOW()")).Error
	if err != nil {
		return apperrors.NewInternalError("Erro excluindo entidades.", err)
	}
	return nil
}

// ============================================================
// MÉTODOS DE CONSULTA DE DEPENDÊNCIAS
// ============================================================

// HasDependentRecords verifica se a entidade tem registros dependentes
func (r *entidadeRepository) HasDependentRecords(id int) (bool, error) {
	counts, err := r.CountDependentRecords(id)
	if err != nil {
		return false, apperrors.NewConflictError("Erro entidade com .")
	}

	for _, count := range counts {
		if count > 0 {
			return true, nil
		}
	}
	return false, nil
}

// CountDependentRecords retorna a contagem de registros dependentes por tipo
func (r *entidadeRepository) CountDependentRecords(id int) (map[string]int64, error) {
	result := make(map[string]int64)

	// Verifica documentos de venda
	var countVendas int64
	if err := r.db.Model(&models.DocumentoVenda{}).
		Where("ent_id = ? AND deleted_at IS NULL", id).
		Count(&countVendas).Error; err != nil {
		return nil, apperrors.NewInternalError("Erro buscando documentos de venda.", err)
	}
	if countVendas > 0 {
		result["documentos_venda"] = countVendas
	}

	// Verifica documentos de entidade
	var countDocumentos int64
	if err := r.db.Model(&models.EntidadeDocumento{}).
		Where("ent_id = ? AND deleted_at IS NULL", id).
		Count(&countDocumentos).Error; err != nil {
		return nil, apperrors.NewInternalError("Erro buscando docuemntos da entidade.", err)
	}
	if countDocumentos > 0 {
		result["documentos_entidade"] = countDocumentos
	}

	// Verifica títulos
	var countTitulos int64
	if err := r.db.Model(&models.Titulo{}).
		Where("ent_id = ? AND deleted_at IS NULL", id).
		Count(&countTitulos).Error; err != nil {
		return nil, apperrors.NewInternalError("Erro buscando titulos.", err)
	}
	if countTitulos > 0 {
		result["titulos"] = countTitulos
	}

	return result, nil
}
