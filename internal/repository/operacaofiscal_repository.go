// internal/repository/operacao_fiscal_repository.go
package repository

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// INTERFACE (Contrato)
// ============================================================

// OperacaoFiscalRepository define o contrato para operações de banco
type OperacaoFiscalRepository interface {
	// CRUD Básico
	Create(operacao *models.OperacaoFiscal) error
	Update(id int, operacao *models.OperacaoFiscal) error
	Delete(id int) error
	FindByID(id int) (*models.OperacaoFiscal, error)
	GetByID(id int) (*models.OperacaoFiscal, error)

	// Buscas Específicas
	FindByCFOP(cfop string) ([]models.OperacaoFiscal, error)
	FindByCFOPAndFilial(cfop string, filialID int) (*models.OperacaoFiscal, error)
	FindByEmpresaFilialID(filialID int) ([]models.OperacaoFiscal, error)
	FindActive() ([]models.OperacaoFiscal, error)
	FindActiveByFilial(filialID int) ([]models.OperacaoFiscal, error)
	FindByTipoOperacao(tipo string) ([]models.OperacaoFiscal, error)
	FindByTipoAndFilial(tipo string, filialID int) ([]models.OperacaoFiscal, error)
	FindByCST(cstID int, cstType string) ([]models.OperacaoFiscal, error)

	// Listagem com Filtros
	List(limit, offset int, filters map[string]interface{}) ([]models.OperacaoFiscal, int64, error)
	ListWithFullPreload(limit, offset int, filters map[string]interface{}) ([]models.OperacaoFiscal, int64, error)
	FindAll() ([]models.OperacaoFiscal, error)

	// Consultas de Validação (APENAS CONSULTAS, SEM VALIDAÇÕES)
	ExistsByCFOPAndFilial(cfop string, filialID int, excludeID int) (bool, error)
	ExistsByID(id int) (bool, error)
	HasActivePeriodConflict(filialID int, dataIni time.Time, dataFim *time.Time, excludeID int) (bool, error)

	// Métodos de Negócio
	GetCurrentByFilial(filialID int) ([]models.OperacaoFiscal, error)
	GetByPeriod(filialID int, data time.Time) ([]models.OperacaoFiscal, error)
	BulkUpdateStatus(ids []int, status string) error
	BulkDelete(ids []int) error
	CopyToFilial(sourceID int, targetFilialID int) error

	// Método para verificar dependências
	HasDependentRecords(id int) (bool, error)
}

// ============================================================
// IMPLEMENTAÇÃO CONCRETA (privada)
// ============================================================

// operacaoFiscalRepository é a implementação concreta
type operacaoFiscalRepository struct {
	db *gorm.DB
}

// NewOperacaoFiscalRepository cria uma nova instância (retorna a interface)
func NewOperacaoFiscalRepository(db *gorm.DB) OperacaoFiscalRepository {
	return &operacaoFiscalRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD (APENAS PERSISTÊNCIA)
// ============================================================

// Create salva uma nova operação fiscal
func (r *operacaoFiscalRepository) Create(operacao *models.OperacaoFiscal) error {
	return r.db.Create(operacao).Error
}

// Update atualiza uma operação fiscal existente
func (r *operacaoFiscalRepository) Update(id int, operacao *models.OperacaoFiscal) error {
	return r.db.
		Omit("EmpresaFilial", "CSTIPI", "CSTPISCOFINS", "CSTICMS", "created_at", "deleted_at").
		Model(&models.OperacaoFiscal{}).
		Where("opf_id = ?", id).
		Updates(operacao).Error
}

// Delete realiza exclusão lógica
func (r *operacaoFiscalRepository) Delete(id int) error {
	return r.db.
		Model(&models.OperacaoFiscal{}).
		Where("opf_id = ?", id).
		Update("deleted_at", time.Now()).Error
}

// FindByID busca uma operação fiscal pelo ID
func (r *operacaoFiscalRepository) FindByID(id int) (*models.OperacaoFiscal, error) {
	var operacao models.OperacaoFiscal
	err := r.db.
		Preload("EmpresaFilial").
		Preload("CSTIPI").
		Preload("CSTPISCOFINS").
		Preload("CSTICMS").
		Where("opf_id = ? AND deleted_at IS NULL", id).
		First(&operacao).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("operação fiscal com ID %d não encontrada", id))
		}
		return nil, err
	}
	return &operacao, nil
}

// GetByID busca uma operação fiscal pelo ID (alias)
func (r *operacaoFiscalRepository) GetByID(id int) (*models.OperacaoFiscal, error) {
	return r.FindByID(id)
}

// ============================================================
// MÉTODOS DE BUSCA ESPECÍFICOS
// ============================================================

// FindByCFOP busca operações fiscais pelo CFOP
func (r *operacaoFiscalRepository) FindByCFOP(cfop string) ([]models.OperacaoFiscal, error) {
	var operacoes []models.OperacaoFiscal
	err := r.db.
		Preload("EmpresaFilial").
		Where("opf_cfop = ? AND deleted_at IS NULL", cfop).
		Order("opf_dataini DESC").
		Find(&operacoes).Error

	if err != nil {
		return nil, err
	}
	return operacoes, nil
}

// FindByCFOPAndFilial busca uma operação fiscal por CFOP e filial
func (r *operacaoFiscalRepository) FindByCFOPAndFilial(cfop string, filialID int) (*models.OperacaoFiscal, error) {
	var operacao models.OperacaoFiscal
	err := r.db.
		Preload("EmpresaFilial").
		Preload("CSTIPI").
		Preload("CSTPISCOFINS").
		Preload("CSTICMS").
		Where("opf_cfop = ? AND emf_id = ? AND deleted_at IS NULL", cfop, filialID).
		First(&operacao).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(
				fmt.Sprintf("operação fiscal com CFOP %s para filial %d não encontrada", cfop, filialID),
			)
		}
		return nil, err
	}
	return &operacao, nil
}

// FindByEmpresaFilialID busca operações fiscais por filial
func (r *operacaoFiscalRepository) FindByEmpresaFilialID(filialID int) ([]models.OperacaoFiscal, error) {
	var operacoes []models.OperacaoFiscal
	err := r.db.
		Preload("CSTIPI").
		Preload("CSTPISCOFINS").
		Preload("CSTICMS").
		Where("emf_id = ? AND deleted_at IS NULL", filialID).
		Order("opf_cfop ASC").
		Find(&operacoes).Error

	if err != nil {
		return nil, err
	}
	return operacoes, nil
}

// FindActive busca todas as operações fiscais ativas
func (r *operacaoFiscalRepository) FindActive() ([]models.OperacaoFiscal, error) {
	var operacoes []models.OperacaoFiscal
	now := time.Now()

	err := r.db.
		Preload("EmpresaFilial").
		Where("deleted_at IS NULL").
		Where("opf_dataini <= ?", now).
		Where("opf_datafim IS NULL OR opf_datafim >= ?", now).
		Order("opf_cfop ASC").
		Find(&operacoes).Error

	if err != nil {
		return nil, err
	}
	return operacoes, nil
}

// FindActiveByFilial busca operações fiscais ativas por filial
func (r *operacaoFiscalRepository) FindActiveByFilial(filialID int) ([]models.OperacaoFiscal, error) {
	var operacoes []models.OperacaoFiscal
	now := time.Now()

	err := r.db.
		Preload("CSTIPI").
		Preload("CSTPISCOFINS").
		Preload("CSTICMS").
		Where("emf_id = ? AND deleted_at IS NULL", filialID).
		Where("opf_dataini <= ?", now).
		Where("opf_datafim IS NULL OR opf_datafim >= ?", now).
		Order("opf_cfop ASC").
		Find(&operacoes).Error

	if err != nil {
		return nil, err
	}
	return operacoes, nil
}

// FindByTipoOperacao busca operações fiscais por tipo (Entrada/Saída)
func (r *operacaoFiscalRepository) FindByTipoOperacao(tipo string) ([]models.OperacaoFiscal, error) {
	var operacoes []models.OperacaoFiscal
	err := r.db.
		Preload("EmpresaFilial").
		Where("deleted_at IS NULL").
		Where("SUBSTRING(opf_cfop, 1, 1) = ?", r.getTipoPrefix(tipo)).
		Order("opf_cfop ASC").
		Find(&operacoes).Error

	if err != nil {
		return nil, err
	}
	return operacoes, nil
}

// FindByTipoAndFilial busca operações fiscais por tipo e filial
func (r *operacaoFiscalRepository) FindByTipoAndFilial(tipo string, filialID int) ([]models.OperacaoFiscal, error) {
	var operacoes []models.OperacaoFiscal
	err := r.db.
		Preload("CSTIPI").
		Preload("CSTPISCOFINS").
		Preload("CSTICMS").
		Where("emf_id = ? AND deleted_at IS NULL", filialID).
		Where("SUBSTRING(opf_cfop, 1, 1) = ?", r.getTipoPrefix(tipo)).
		Order("opf_cfop ASC").
		Find(&operacoes).Error

	if err != nil {
		return nil, err
	}
	return operacoes, nil
}

// FindByCST busca operações fiscais por CST
func (r *operacaoFiscalRepository) FindByCST(cstID int, cstType string) ([]models.OperacaoFiscal, error) {
	var operacoes []models.OperacaoFiscal
	query := r.db.Where("deleted_at IS NULL")

	switch cstType {
	case "ICMS":
		query = query.Where("csticms_id = ?", cstID)
	case "IPI":
		query = query.Where("cstipi_id = ?", cstID)
	case "PISCOFINS":
		query = query.Where("cstpiscofins_id = ?", cstID)
	default:
		return nil, errors.New("tipo de CST inválido")
	}

	err := query.
		Preload("EmpresaFilial").
		Order("opf_cfop ASC").
		Find(&operacoes).Error

	if err != nil {
		return nil, err
	}
	return operacoes, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// List retorna uma lista de operações fiscais com paginação e filtros
func (r *operacaoFiscalRepository) List(limit, offset int, filters map[string]interface{}) ([]models.OperacaoFiscal, int64, error) {
	var operacoes []models.OperacaoFiscal
	var total int64

	query := r.db.Model(&models.OperacaoFiscal{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.OperacaoFiscal{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("EmpresaFilial").
		Preload("CSTIPI").
		Preload("CSTPISCOFINS").
		Preload("CSTICMS").
		Limit(limit).
		Offset(offset).
		Order("opf_cfop ASC").
		Find(&operacoes).Error

	if err != nil {
		return nil, 0, err
	}

	return operacoes, total, nil
}

// ListWithFullPreload lista com todos os relacionamentos
func (r *operacaoFiscalRepository) ListWithFullPreload(limit, offset int, filters map[string]interface{}) ([]models.OperacaoFiscal, int64, error) {
	var operacoes []models.OperacaoFiscal
	var total int64

	query := r.db.Model(&models.OperacaoFiscal{}).Where("deleted_at IS NULL")
	query = utils.ApplyFilters(query, models.OperacaoFiscal{}, filters)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("EmpresaFilial").
		Preload("CSTIPI").
		Preload("CSTPISCOFINS").
		Preload("CSTICMS").
		Limit(limit).
		Offset(offset).
		Order("opf_cfop ASC").
		Find(&operacoes).Error

	if err != nil {
		return nil, 0, err
	}

	return operacoes, total, nil
}

// FindAll busca todas as operações fiscais
func (r *operacaoFiscalRepository) FindAll() ([]models.OperacaoFiscal, error) {
	var operacoes []models.OperacaoFiscal
	err := r.db.
		Where("deleted_at IS NULL").
		Order("opf_cfop ASC").
		Find(&operacoes).Error
	if err != nil {
		return nil, err
	}
	return operacoes, nil
}

// ============================================================
// MÉTODOS DE CONSULTA PARA VALIDAÇÕES (APENAS CONSULTAS)
// ============================================================

// ExistsByCFOPAndFilial verifica se existe operação com CFOP para filial
func (r *operacaoFiscalRepository) ExistsByCFOPAndFilial(cfop string, filialID int, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.OperacaoFiscal{}).
		Where("opf_cfop = ? AND emf_id = ? AND deleted_at IS NULL", cfop, filialID)

	if excludeID > 0 {
		query = query.Where("opf_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// ExistsByID verifica se uma operação existe
func (r *operacaoFiscalRepository) ExistsByID(id int) (bool, error) {
	var count int64
	err := r.db.Model(&models.OperacaoFiscal{}).
		Where("opf_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// HasActivePeriodConflict verifica conflito de período ativo
func (r *operacaoFiscalRepository) HasActivePeriodConflict(filialID int, dataIni time.Time, dataFim *time.Time, excludeID int) (bool, error) {
	var count int64
	query := r.db.Model(&models.OperacaoFiscal{}).
		Where("emf_id = ? AND deleted_at IS NULL", filialID).
		Where("? BETWEEN opf_dataini AND COALESCE(opf_datafim, '9999-12-31')", dataIni)

	if dataFim != nil {
		query = query.Or("? BETWEEN opf_dataini AND COALESCE(opf_datafim, '9999-12-31')", *dataFim)
	}

	if excludeID > 0 {
		query = query.Where("opf_id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// HasDependentRecords verifica se existem registros dependentes
func (r *operacaoFiscalRepository) HasDependentRecords(id int) (bool, error) {
	var count int64

	// Verifica se existem itens de documento usando esta operação fiscal
	err := r.db.Model(&models.DocumentoVendaItem{}).
		Where("opf_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}

	return false, nil
}

// ============================================================
// MÉTODOS DE NEGÓCIO (APENAS OPERAÇÕES DE BANCO)
// ============================================================

// GetCurrentByFilial busca operações fiscais atuais por filial
func (r *operacaoFiscalRepository) GetCurrentByFilial(filialID int) ([]models.OperacaoFiscal, error) {
	return r.FindActiveByFilial(filialID)
}

// GetByPeriod busca operações fiscais em um período específico
func (r *operacaoFiscalRepository) GetByPeriod(filialID int, data time.Time) ([]models.OperacaoFiscal, error) {
	var operacoes []models.OperacaoFiscal
	err := r.db.
		Preload("CSTIPI").
		Preload("CSTPISCOFINS").
		Preload("CSTICMS").
		Where("emf_id = ? AND deleted_at IS NULL", filialID).
		Where("opf_dataini <= ?", data).
		Where("opf_datafim IS NULL OR opf_datafim >= ?", data).
		Order("opf_cfop ASC").
		Find(&operacoes).Error

	if err != nil {
		return nil, err
	}
	return operacoes, nil
}

// BulkUpdateStatus atualiza status de múltiplas operações
func (r *operacaoFiscalRepository) BulkUpdateStatus(ids []int, status string) error {
	now := time.Now()
	var err error

	if status == "ativo" {
		// Ativa: remove data_fim
		err = r.db.Model(&models.OperacaoFiscal{}).
			Where("opf_id IN ? AND deleted_at IS NULL", ids).
			Update("opf_datafim", nil).Error
	} else {
		// Inativa: define data_fim para hoje
		err = r.db.Model(&models.OperacaoFiscal{}).
			Where("opf_id IN ? AND deleted_at IS NULL", ids).
			Update("opf_datafim", now).Error
	}

	return err
}

// BulkDelete realiza exclusão lógica de múltiplas operações
func (r *operacaoFiscalRepository) BulkDelete(ids []int) error {
	return r.db.Model(&models.OperacaoFiscal{}).
		Where("opf_id IN ? AND deleted_at IS NULL", ids).
		Update("deleted_at", time.Now()).Error
}

// CopyToFilial copia operações fiscais de uma filial para outra
func (r *operacaoFiscalRepository) CopyToFilial(sourceID int, targetFilialID int) error {
	// Busca operações da filial fonte
	var operacoes []models.OperacaoFiscal
	err := r.db.
		Where("emf_id = ? AND deleted_at IS NULL", sourceID).
		Find(&operacoes).Error
	if err != nil {
		return err
	}

	if len(operacoes) == 0 {
		return nil // Nada para copiar
	}

	// Prepara para cópia
	for i := range operacoes {
		// Limpa IDs e ajusta filial
		operacoes[i].ID = 0
		operacoes[i].EmpresaFilialID = targetFilialID
		operacoes[i].CreatedAt = time.Now()
		operacoes[i].UpdatedAt = time.Now()
		operacoes[i].DeletedAt = nil
	}

	// Insere em lote
	return r.db.Create(&operacoes).Error
}

// ============================================================
// MÉTODOS PRIVADOS (AUXILIARES DE CONSULTA)
// ============================================================

// getTipoPrefix retorna o prefixo do CFOP baseado no tipo
func (r *operacaoFiscalRepository) getTipoPrefix(tipo string) string {
	switch tipo {
	case "Entrada":
		return "5"
	case "Saída":
		return "6"
	default:
		return ""
	}
}
