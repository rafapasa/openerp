// internal/repository/documento_venda_pagamento_repository.go
package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
)

// ============================================================
// INTERFACE - Define o contrato
// ============================================================

// DocumentoVendaPagamentoRepository define o contrato para operações de banco
type DocumentoVendaPagamentoRepository interface {
	// CRUD Básico
	Create(pagamento *models.DocumentoVendaPagamento) error
	Update(ddvId, dvpItem int, pagamento *models.DocumentoVendaPagamento) error
	Delete(ddvId, dvpItem int) error
	FindByID(ddvId, dvpItem int) (*models.DocumentoVendaPagamento, error)

	// Listagem com Filtros
	ListByDocumentoVendaID(ddvId int) ([]models.DocumentoVendaPagamento, int64, error)

	// Consultas de Validação (APENAS CONSULTAS)
	GetNextItemNumber(ddvId int) (int, error)
	CountByDocumentoVendaID(ddvId int) (int64, error)
	ExistsByDocumentoVendaID(ddvId int) (bool, error)
	SumPagamentosByDocumentoVendaID(ddvId int) (float64, error)
}

// ============================================================
// IMPLEMENTAÇÃO CONCRETA (privada)
// ============================================================

type documentoVendaPagamentoRepository struct {
	db *gorm.DB
}

// NewDocumentoVendaPagamentoRepository cria uma nova instância do repositório
// ✅ CORRETO: Retorna a interface, não a struct concreta
func NewDocumentoVendaPagamentoRepository(db *gorm.DB) DocumentoVendaPagamentoRepository {
	return &documentoVendaPagamentoRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD (APENAS PERSISTÊNCIA)
// ============================================================

// Create cria um novo pagamento com item automático
func (r *documentoVendaPagamentoRepository) Create(pagamento *models.DocumentoVendaPagamento) error {
	return r.db.Create(pagamento).Error
}

// Update atualiza um pagamento existente
func (r *documentoVendaPagamentoRepository) Update(ddvId, dvpItem int, pagamento *models.DocumentoVendaPagamento) error {
	return r.db.Model(&models.DocumentoVendaPagamento{}).
		Where("ddv_id = ? AND dvp_item = ?", ddvId, dvpItem).
		Omit("created_at", "deleted_at").
		Updates(pagamento).Error
}

// Delete realiza soft delete de um pagamento
func (r *documentoVendaPagamentoRepository) Delete(ddvId, dvpItem int) error {
	return r.db.Model(&models.DocumentoVendaPagamento{}).
		Where("ddv_id = ? AND dvp_item = ?", ddvId, dvpItem).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// FindByID busca um pagamento por ID composto (ddv_id + dvp_item)
func (r *documentoVendaPagamentoRepository) FindByID(ddvId, dvpItem int) (*models.DocumentoVendaPagamento, error) {
	var pagamento models.DocumentoVendaPagamento
	err := r.db.
		Where("ddv_id = ? AND dvp_item = ? AND deleted_at IS NULL", ddvId, dvpItem).
		First(&pagamento).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("Pagamento %d do documento %d não encontrado.", dvpItem, ddvId))
		}
		return nil, err
	}
	return &pagamento, nil
}

// ============================================================
// MÉTODOS DE LISTAGEM
// ============================================================

// ListByDocumentoVendaID lista todos os pagamentos de um documento de venda com paginação
func (r *documentoVendaPagamentoRepository) ListByDocumentoVendaID(ddvId int) ([]models.DocumentoVendaPagamento, int64, error) {
	var pagamentos []models.DocumentoVendaPagamento
	var total int64

	err := r.db.Model(&models.DocumentoVendaPagamento{}).
		Preload("Portador").
		Preload("FormaPagamento").
		Preload("DocumentoVenda").
		Where("ddv_id = ? AND deleted_at IS NULL", ddvId).
		Order("dvp_item ASC").
		Find(&pagamentos).Error

	if err != nil {
		return nil, 0, err
	}

	total = int64(len(pagamentos))

	return pagamentos, total, nil
}

// ============================================================
// MÉTODOS DE CONSULTA PARA VALIDAÇÕES (APENAS CONSULTAS)
// ============================================================

// GetNextItemNumber retorna o próximo número de item para um documento
func (r *documentoVendaPagamentoRepository) GetNextItemNumber(ddvId int) (int, error) {
	var maxItem int
	err := r.db.Model(&models.DocumentoVendaPagamento{}).
		Where("ddv_id = ?", ddvId).
		Select("COALESCE(MAX(dvp_item), 0) + 1").
		Scan(&maxItem).Error
	if err != nil {
		return 0, err
	}
	return maxItem, nil
}

// CountByDocumentoVendaID retorna a quantidade de pagamentos de um documento
func (r *documentoVendaPagamentoRepository) CountByDocumentoVendaID(ddvId int) (int64, error) {
	var count int64
	err := r.db.Model(&models.DocumentoVendaPagamento{}).
		Where("ddv_id = ? AND deleted_at IS NULL", ddvId).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// ExistsByDocumentoVendaID verifica se um documento possui pagamentos
func (r *documentoVendaPagamentoRepository) ExistsByDocumentoVendaID(ddvId int) (bool, error) {
	var count int64
	err := r.db.Model(&models.DocumentoVendaPagamento{}).
		Where("ddv_id = ? AND deleted_at IS NULL", ddvId).
		Limit(1).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SumPagamentosByDocumentoVendaID retorna o valor total dos pagamentos de um documento
func (r *documentoVendaPagamentoRepository) SumPagamentosByDocumentoVendaID(ddvId int) (float64, error) {
	var total float64
	err := r.db.Model(&models.DocumentoVendaPagamento{}).
		Where("ddv_id = ? AND deleted_at IS NULL", ddvId).
		Select("COALESCE(SUM(dvp_valor), 0)").
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}
