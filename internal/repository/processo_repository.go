// internal/repository/processo_repository.go
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
// INTERFACE
// ============================================================

// ProcessoRepository define o contrato para operações de banco
type ProcessoRepository interface {
	// CRUD Básico
	Create(processo *models.Processo) error
	Update(id int, processo *models.Processo) error
	Delete(id int) error
	FindByID(id int) (*models.Processo, error)
	GetByID(id int) (*models.Processo, error)

	// Buscas Específicas
	FindByCodigo(codigo int) (*models.Processo, error)
	FindByCodigoAndFilial(codigo int, filialID int) (*models.Processo, error)
	FindByDescricao(descricao string, limit int) ([]models.Processo, error)
	FindBySituacao(situacao int) ([]models.Processo, error)
	FindByTipoOperacao(tipoOperacao int) ([]models.Processo, error)
	FindByFilialID(filialID int) ([]models.Processo, error)
	FindActive() ([]models.Processo, error)
	FindActiveByFilial(filialID int) ([]models.Processo, error)

	// Listagem com Filtros
	List(limit, offset int, filters map[string]interface{}) ([]models.Processo, int64, error)
	ListWithFullPreload(limit, offset int, filters map[string]interface{}) ([]models.Processo, int64, error)
	FindAll() ([]models.Processo, error)

	// Consultas de Validação (APENAS CONSULTAS)
	ExistsByCodigo(codigo int, excludeID int) (bool, error)
	ExistsByCodigoAndFilial(codigo int, filialID int, excludeID int) (bool, error)
	ExistsByDescricao(descricao string, excludeID int) (bool, error)
	ExistsByID(id int) (bool, error)
	CountByFilialID(filialID int) (int64, error)
	GetMaxCodigoByFilial(filialID int) (int, error)

	// Operações em Lote
	BulkUpdateStatus(ids []int, situacao int) error
	BulkDelete(ids []int) error

	// Consultas de Dependências
	HasDependentRecords(id int) (bool, error)
	CountDependentRecords(id int) (map[string]int64, error)
}

// ============================================================
// IMPLEMENTAÇÃO CONCRETA (privada)
// ============================================================

type processoRepository struct {
	db *gorm.DB
}

// NewProcessoRepository cria uma nova instância (retorna a interface)
func NewProcessoRepository(db *gorm.DB) ProcessoRepository {
	return &processoRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD (APENAS PERSISTÊNCIA)
// ============================================================

// Create salva um novo processo
func (r *processoRepository) Create(processo *models.Processo) error {
	return r.db.Create(processo).Error
}

// Update atualiza um processo existente
func (r *processoRepository) Update(id int, processo *models.Processo) error {
	return r.db.
		Omit("EmpresaFilial", "RotinaContabil", "OperacaoFiscalForaEst", "OperacaoFiscalNoEst", "OperacaoFiscalForaEstST", "OperacaoFiscalNoEstST", "PlanoContasFinanceiro", "NaturezaOperacao", "Receita", "Despesa", "ProcessoNF", "DocumentoVendas", "NotaFiscais", "Operacoes", "created_at", "deleted_at").
		Model(&models.Processo{}).
		Where("prc_id = ?", id).
		Updates(processo).Error
}

// Delete realiza exclusão lógica
func (r *processoRepository) Delete(id int) error {
	return r.db.
		Model(&models.Processo{}).
		Where("prc_id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// FindByID busca um processo pelo ID com relacionamentos
func (r *processoRepository) FindByID(id int) (*models.Processo, error) {
	var processo models.Processo
	err := r.db.
		Preload("EmpresaFilial").
		Preload("RotinaContabil").
		Preload("OperacaoFiscalForaEst").
		Preload("OperacaoFiscalNoEst").
		Preload("OperacaoFiscalForaEstST").
		Preload("OperacaoFiscalNoEstST").
		Preload("PlanoContasFinanceiro").
		Preload("NaturezaOperacao").
		Preload("Receita").
		Preload("Despesa").
		Preload("ProcessoNF").
		Preload("DocumentoVendas", "deleted_at IS NULL").
		Preload("NotaFiscais", "deleted_at IS NULL").
		Preload("Operacoes", "deleted_at IS NULL").
		Where("prc_id = ? AND deleted_at IS NULL", id).
		First(&processo).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("processo com ID %d não encontrado", id))
		}
		return nil, err
	}
	return &processo, nil
}

// GetByID busca um processo pelo ID (sem relacionamentos)
func (r *processoRepository) GetByID(id int) (*models.Processo, error) {
	var processo models.Processo
	err := r.db.
		Where("prc_id = ? AND deleted_at IS NULL", id).
		First(&processo).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("processo com ID %d não encontrado", id))
		}
		return nil, err
	}
	return &processo, nil
}

// ============================================================
// MÉTODOS DE BUSCA ESPECÍFICOS
// ============================================================

// FindByCodigo busca um processo pelo código
func (r *processoRepository) FindByCodigo(codigo int) (*models.Processo, error) {
	var processo models.Processo
	err := r.db.
		Preload("EmpresaFilial").
		Where("prc_codigo = ? AND deleted_at IS NULL", codigo).
		First(&processo).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("processo com código %d não encontrado", codigo))
		}
		return nil, err
	}
	return &processo, nil
}

// FindByCodigoAndFilial busca um processo pelo código e filial
func (r *processoRepository) FindByCodigoAndFilial(codigo int, filialID int) (*models.Processo, error) {
	var processo models.Processo
	err := r.db.
		Preload("EmpresaFilial").
		Where("prc_codigo = ? AND emf_id = ? AND deleted_at IS NULL", codigo, filialID).
		First(&processo).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("processo com código %d e filial %d não encontrado", codigo, filialID))
		}
		return nil, err
	}
	return &processo, nil
}

// FindByDescricao busca processos pela descrição (autocomplete)
func (r *processoRepository) FindByDescricao(descricao string, limit int) ([]models.Processo, error) {
	var processos []models.Processo
	err := r.db.
		Preload("EmpresaFilial").
		Where("prc_descricao LIKE ? AND deleted_at IS NULL", "%"+descricao+"%").
		Order("prc_descricao ASC").
		Limit(limit).
		Find(&processos).Error

	if err != nil {
		return nil, err
	}
	return processos, nil
}

// FindBySituacao busca processos por situação
func (r *processoRepository) FindBySituacao(situacao int) ([]models.Processo, error) {
	var processos []models.Processo
	err := r.db.
		Preload("EmpresaFilial").
		Where("prc_situacao = ? AND deleted_at IS NULL", situacao).
		Order("prc_descricao ASC").
		Find(&processos).Error

	if err != nil {
		return nil, err
	}
	return processos, nil
}

// FindByTipoOperacao busca processos por tipo de operação
func (r *processoRepository) FindByTipoOperacao(tipoOperacao int) ([]models.Processo, error) {
	var processos []models.Processo
	err := r.db.
		Preload("EmpresaFilial").
		Where("prc_tipooperacao = ? AND deleted_at IS NULL", tipoOperacao).
		Order("prc_descricao ASC").
		Find(&processos).Error

	if err != nil {
		return nil, err
	}
	return processos, nil
}

// FindByFilialID busca processos por filial
func (r *processoRepository) FindByFilialID(filialID int) ([]models.Processo, error) {
	var processos []models.Processo
	err := r.db.
		Where("emf_id = ? AND deleted_at IS NULL", filialID).
		Order("prc_descricao ASC").
		Find(&processos).Error

	if err != nil {
		return nil, err
	}
	return processos, nil
}

// FindActive busca todos os processos ativos
func (r *processoRepository) FindActive() ([]models.Processo, error) {
	var processos []models.Processo
	err := r.db.
		Preload("EmpresaFilial").
		Where("prc_situacao = 1 AND deleted_at IS NULL").
		Order("prc_descricao ASC").
		Find(&processos).Error

	if err != nil {
		return nil, err
	}
	return processos, nil
}

// FindActiveByFilial busca processos ativos por filial
func (r *processoRepository) FindActiveByFilial(filialID int) ([]models.Processo, error) {
	var processos []models.Processo
	err := r.db.
		Preload("EmpresaFilial").
		Where("emf_id = ? AND prc_situacao = 1 AND deleted_at IS NULL", filialID).
		Order("prc_descricao ASC").
		Find(&processos).Error

	if err != nil {
		return nil, err
	}
	return processos, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de processos com paginação e filtros
func (r *processoRepository) List(limit, offset int, filters map[string]interface{}) ([]models.Processo, int64, error) {
	var processos []models.Processo
	var total int64

	query := r.db.Model(&models.Processo{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.Processo{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("EmpresaFilial").
		Preload("RotinaContabil").
		Limit(limit).
		Offset(offset).
		Order("prc_id DESC").
		Find(&processos).Error

	if err != nil {
		return nil, 0, err
	}

	return processos, total, nil
}

// ListWithFullPreload retorna uma lista com todos os relacionamentos
func (r *processoRepository) ListWithFullPreload(limit, offset int, filters map[string]interface{}) ([]models.Processo, int64, error) {
	var processos []models.Processo
	var total int64

	query := r.db.Model(&models.Processo{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.Processo{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("EmpresaFilial").
		Preload("RotinaContabil").
		Preload("OperacaoFiscalForaEst").
		Preload("OperacaoFiscalNoEst").
		Preload("OperacaoFiscalForaEstST").
		Preload("OperacaoFiscalNoEstST").
		Preload("PlanoContasFinanceiro").
		Preload("NaturezaOperacao").
		Preload("Receita").
		Preload("Despesa").
		Preload("ProcessoNF").
		Limit(limit).
		Offset(offset).
		Order("prc_id DESC").
		Find(&processos).Error

	if err != nil {
		return nil, 0, err
	}

	return processos, total, nil
}

// FindAll busca todos os processos
func (r *processoRepository) FindAll() ([]models.Processo, error) {
	var processos []models.Processo
	err := r.db.
		Where("deleted_at IS NULL").
		Order("prc_descricao ASC").
		Find(&processos).Error
	if err != nil {
		return nil, err
	}
	return processos, nil
}

// ============================================================
// MÉTODOS DE CONSULTA PARA VALIDAÇÕES (APENAS CONSULTAS)
// ============================================================

// ExistsByCodigo verifica se já existe um processo com o código
func (r *processoRepository) ExistsByCodigo(codigo int, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.Processo{}).
		Where("prc_codigo = ? AND deleted_at IS NULL", codigo)

	if excludeID > 0 {
		query = query.Where("prc_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByCodigoAndFilial verifica se já existe um processo com o código e filial
func (r *processoRepository) ExistsByCodigoAndFilial(codigo int, filialID int, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.Processo{}).
		Where("prc_codigo = ? AND emf_id = ? AND deleted_at IS NULL", codigo, filialID)

	if excludeID > 0 {
		query = query.Where("prc_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByDescricao verifica se já existe um processo com a descrição
func (r *processoRepository) ExistsByDescricao(descricao string, excludeID int) (bool, error) {
	if descricao == "" {
		return false, nil
	}

	var count int64
	query := r.db.Model(&models.Processo{}).
		Where("prc_descricao = ? AND deleted_at IS NULL", descricao)

	if excludeID > 0 {
		query = query.Where("prc_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByID verifica se um processo existe pelo ID
func (r *processoRepository) ExistsByID(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.Processo{}).
		Where("prc_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountByFilialID retorna a quantidade de processos de uma filial
func (r *processoRepository) CountByFilialID(filialID int) (int64, error) {
	var count int64
	err := r.db.Model(&models.Processo{}).
		Where("emf_id = ? AND deleted_at IS NULL", filialID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetMaxCodigoByFilial retorna o maior código de processo para uma filial
func (r *processoRepository) GetMaxCodigoByFilial(filialID int) (int, error) {
	var maxCodigo int
	err := r.db.Model(&models.Processo{}).
		Where("emf_id = ? AND deleted_at IS NULL", filialID).
		Select("COALESCE(MAX(prc_codigo), 0)").
		Scan(&maxCodigo).Error
	if err != nil {
		return 0, err
	}
	return maxCodigo, nil
}

// ============================================================
// OPERAÇÕES EM LOTE
// ============================================================

// BulkUpdateStatus atualiza a situação de múltiplos processos
func (r *processoRepository) BulkUpdateStatus(ids []int, situacao int) error {
	return r.db.Model(&models.Processo{}).
		Where("prc_id IN ? AND deleted_at IS NULL", ids).
		Update("prc_situacao", situacao).Error
}

// BulkDelete realiza exclusão lógica de múltiplos processos
func (r *processoRepository) BulkDelete(ids []int) error {
	return r.db.Model(&models.Processo{}).
		Where("prc_id IN ? AND deleted_at IS NULL", ids).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// ============================================================
// MÉTODOS DE CONSULTA DE DEPENDÊNCIAS
// ============================================================

// HasDependentRecords verifica se o processo tem registros dependentes
func (r *processoRepository) HasDependentRecords(id int) (bool, error) {
	counts, err := r.CountDependentRecords(id)
	if err != nil {
		return false, err
	}

	for _, count := range counts {
		if count > 0 {
			return true, nil
		}
	}
	return false, nil
}

// CountDependentRecords retorna a contagem de registros dependentes por tipo
func (r *processoRepository) CountDependentRecords(id int) (map[string]int64, error) {
	result := make(map[string]int64)

	// Verifica documentos de venda
	var countVendas int64
	if err := r.db.Model(&models.DocumentoVenda{}).
		Where("prc_id = ? AND deleted_at IS NULL", id).
		Count(&countVendas).Error; err != nil {
		return nil, err
	}
	if countVendas > 0 {
		result["documento_vendas"] = countVendas
	}

	// Verifica notas fiscais
	var countNotas int64
	if err := r.db.Model(&models.NotaFiscal{}).
		Where("prc_id = ? AND deleted_at IS NULL", id).
		Count(&countNotas).Error; err != nil {
		return nil, err
	}
	if countNotas > 0 {
		result["nota_fiscais"] = countNotas
	}

	// Verifica operações do processo
	var countOperacoes int64
	if err := r.db.Model(&models.ProcessoOperacao{}).
		Where("prc_id = ? AND deleted_at IS NULL", id).
		Count(&countOperacoes).Error; err != nil {
		return nil, err
	}
	if countOperacoes > 0 {
		result["operacoes"] = countOperacoes
	}

	return result, nil
}
