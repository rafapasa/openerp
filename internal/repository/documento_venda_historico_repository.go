package repository

import (
	"context"
	"errors"
	"time" // Adicionado para parsing de datas

	"gorm.io/gorm"

	apperrors "github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
)

// ============================================================
// INTERFACE
// ============================================================

type DocumentVendaHistoricoRepository interface {
	Create(ctx context.Context, historico *models.DocumentoVendaHistorico) error
	FindByID(ctx context.Context, documentVendaID int, item int) (*models.DocumentoVendaHistorico, error)
	ExistByID(ctx context.Context, documentVendaID int, item int) (bool, error)
	FindByDocumentVendaID(ctx context.Context, documentVendaID int) ([]models.DocumentoVendaHistorico, error)
	FindAll(ctx context.Context, filter *dto.DocumentoVendaHistoricoFilter) ([]models.DocumentoVendaHistorico, int64, error)
	Update(ctx context.Context, historico *models.DocumentoVendaHistorico) error
	Delete(ctx context.Context, documentVendaID int, item int) error
	GetLastByDocumentVendaID(ctx context.Context, documentVendaID int) (*models.DocumentoVendaHistorico, error)
}

// ============================================================
// TYPES
// ============================================================

type documentVendaHistoricoRepository struct {
	db *gorm.DB
}

// ============================================================
// CONSTRUCTOR
// ============================================================

func NewDocumentoVendaHistoricoRepository(db *gorm.DB) DocumentVendaHistoricoRepository {
	return &documentVendaHistoricoRepository{db: db}
}

// ============================================================
// MÉTODOS CRUD
// ============================================================

// Create cria um novo histórico de documento de venda
func (r *documentVendaHistoricoRepository) Create(ctx context.Context, historico *models.DocumentoVendaHistorico) error {
	err := r.db.
		WithContext(ctx).
		Create(historico).
		Error

	if err != nil {
		return apperrors.NewInternalError("Erro ao criar histórico do documento de venda.", err)
	}

	return nil
}

// Update atualiza um histórico existente
func (r *documentVendaHistoricoRepository) Update(ctx context.Context, historico *models.DocumentoVendaHistorico) error {
	err := r.db.
		WithContext(ctx).
		Save(historico).
		Error

	if err != nil {
		return apperrors.NewInternalError("Erro ao atualizar histórico do documento de venda.", err)
	}

	return nil
}

// Delete realiza exclusão lógica de um histórico
func (r *documentVendaHistoricoRepository) Delete(ctx context.Context, documentVendaID int, item int) error {
	// Buscar o histórico primeiro
	historico, err := r.FindByID(ctx, documentVendaID, item)
	if err != nil {
		return err
	}
	if historico == nil {
		return apperrors.NewNotFoundError("Histórico do documento de venda não encontrado")
	}

	// Realizar soft delete
	historico.SoftDelete()

	err = r.db.
		WithContext(ctx).
		Save(historico).
		Error

	if err != nil {
		return apperrors.NewInternalError("Erro ao deletar histórico do documento de venda.", err)
	}

	return nil
}

// ============================================================
// MÉTODOS DE BUSCA
// ============================================================

// FindByID busca um histórico pelo ID do documento e item
func (r *documentVendaHistoricoRepository) FindByID(ctx context.Context, documentVendaID int, item int) (*models.DocumentoVendaHistorico, error) {
	var historico models.DocumentoVendaHistorico

	err := r.db.
		WithContext(ctx).
		Where("doc_id = ? AND item = ? AND deleted_at IS NULL", documentVendaID, item).
		First(&historico).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, apperrors.NewInternalError("Erro ao buscar histórico do documento de venda.", err)
	}

	return &historico, nil
}

// FindByDocumentVendaID busca todos os históricos de um documento
func (r *documentVendaHistoricoRepository) FindByDocumentVendaID(ctx context.Context, documentVendaID int) ([]models.DocumentoVendaHistorico, error) {
	var historicos []models.DocumentoVendaHistorico

	err := r.db.
		WithContext(ctx).
		Where("doc_id = ? AND deleted_at IS NULL", documentVendaID).
		Order("item ASC").
		Find(&historicos).
		Error

	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao buscar históricos do documento de venda.", err)
	}

	return historicos, nil
}

// FindAll busca históricos com filtros e paginação
func (r *documentVendaHistoricoRepository) FindAll(ctx context.Context, filter *dto.DocumentoVendaHistoricoFilter) ([]models.DocumentoVendaHistorico, int64, error) {
	var historicos []models.DocumentoVendaHistorico
	var total int64

	// Construir query base
	query := r.db.WithContext(ctx).Model(&models.DocumentoVendaHistorico{}).Where("deleted_at IS NULL")

	// Aplicar filtros se existirem
	if filter != nil {
		// Filtro por DocumentVendaID
		if *filter.DocumentoVendaID > 0 {
			query = query.Where("doc_id = ?", filter.DocumentoVendaID)
		}

		// Filtro por data inicial
		if *filter.DataInicio != "" {
			query = query.Where("created_at >= ?", filter.DataInicio)
		}

		// Filtro por data inicial (correção similar à DataFim, embora não esteja no bloco selecionado, é um erro idêntico)
		if filter.DataInicio != nil && *filter.DataInicio != "" {
			parsedDate, err := time.Parse("2006-01-02", *filter.DataInicio)
			if err != nil {
				return nil, 0, apperrors.NewValidationError("Formato de data inicial inválido. Use YYYY-MM-DD.")
			}
			query = query.Where("created_at >= ?", parsedDate)
		}

		// Filtro por data final
		if filter.DataFim != nil && *filter.DataFim != "" { // Correção: Verifica se o ponteiro não é nulo e a string não está vazia
			parsedDate, err := time.Parse("2006-01-02", *filter.DataFim) // Converte a string para time.Time
			if err != nil {
				return nil, 0, apperrors.NewValidationError("Formato de data final inválido. Use YYYY-MM-DD.")
			}
			// Para incluir o dia inteiro na data final, ajustamos para o final do dia
			endOfDay := parsedDate.Add(24 * time.Hour).Add(-time.Second)
			query = query.Where("created_at <= ?", endOfDay) // Usa a data convertida e ajustada
		}
		// Filtro por usuário
		if filter.UsuarioID != nil && *filter.UsuarioID > 0 {
			query = query.Where("usuario_id = ?", filter.UsuarioID)
		}
	}

	// Contar total de registros
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperrors.NewInternalError("Erro ao contar históricos do documento de venda.", err)
	}

	// Buscar com ordenação
	err := query.
		Order("created_at DESC, item ASC").
		Find(&historicos).
		Error

	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro ao buscar históricos do documento de venda.", err)
	}

	return historicos, total, nil
}

// GetLastByDocumentVendaID busca o último histórico de um documento
func (r *documentVendaHistoricoRepository) GetLastByDocumentVendaID(ctx context.Context, documentVendaID int) (*models.DocumentoVendaHistorico, error) {
	var historico models.DocumentoVendaHistorico

	err := r.db.
		WithContext(ctx).
		Where("doc_id = ? AND deleted_at IS NULL", documentVendaID).
		Order("created_at DESC, item DESC").
		First(&historico).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, apperrors.NewInternalError("Erro ao buscar último histórico do documento de venda.", err)
	}

	return &historico, nil
}

func (r *documentVendaHistoricoRepository) ExistByID(ctx context.Context, documentVendaID int, item int) (bool, error) {
	var count int64
	err := r.db.
		WithContext(ctx).
		Where("doc_id = ? AND item = ? AND deleted_at IS NULL", documentVendaID, item).
		Count(&count).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, apperrors.NewNotFoundError("Histórico do documento de venda não encontrado")
		}
		return false, apperrors.NewInternalError("Erro ao verificar existência do histórico.", err)
	}

	return count > 0, nil
}
