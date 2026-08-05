// service/documento_venda_historico_service.go
package service

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// INTERFACE DO SERVICE
// ============================================================

type DocumentoVendaHistoricoService interface {
	Create(ctx context.Context, req *dto.DocumentoVendaHistoricoRequest) (*models.DocumentoVendaHistorico, error)
	GetByID(ctx context.Context, documentoVendaID int, item int) (*models.DocumentoVendaHistorico, error)
	GetByDocumentoVendaID(ctx context.Context, documentoVendaID int) ([]models.DocumentoVendaHistorico, error)
	GetLastByDocumentoVendaID(ctx context.Context, documentoVendaID int) (*models.DocumentoVendaHistorico, error)
	List(ctx context.Context, filter *dto.DocumentoVendaHistoricoFilter) ([]models.DocumentoVendaHistorico, int64, error)
	Update(ctx context.Context, documentoVendaID int, item int, req *dto.DocumentoVendaHistoricoRequest) (*models.DocumentoVendaHistorico, error)
	Delete(ctx context.Context, documentoVendaID int, item int) error
	SoftDelete(ctx context.Context, documentoVendaID int, item int) error
}

// ============================================================
// IMPLEMENTAÇÃO
// ============================================================

type documentoVendaHistoricoService struct {
	repo repository.DocumentVendaHistoricoRepository
	// Dependências para validações de negócio
	documentoVendaService DocumentoVendaService
	usuarioService        UsuarioService
	// fluxoSetorService     FluxoSetorService
}

// ============================================================
// CONSTANTES DE VALIDAÇÃO
// ============================================================

const (
	maxLengthHistoricoObservacao = 500
	minLengthHistoricoObservacao = 3
)

// ============================================================
// CONSTRUTOR
// ============================================================

func NewDocumentoVendaHistoricoService(
	repo repository.DocumentVendaHistoricoRepository,
	documentoVendaService DocumentoVendaService,
	usuarioService UsuarioService,
	// fluxoSetorService FluxoSetorService,
) DocumentoVendaHistoricoService {
	return &documentoVendaHistoricoService{
		repo:                  repo,
		documentoVendaService: documentoVendaService,
		usuarioService:        usuarioService,
		// fluxoSetorService:     fluxoSetorService,
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO (PRIVADOS)
// ============================================================

// validateRequest realiza as validações básicas do DTO
func (s *documentoVendaHistoricoService) validateRequest(req *dto.DocumentoVendaHistoricoRequest) error {
	// 1. Validar campos obrigatórios
	if req.DocumentoVendaID <= 0 {
		return apperrors.NewValidationError("ID do documento de venda é obrigatório")
	}

	if req.UsuarioID <= 0 {
		return apperrors.NewValidationError("ID do usuário é obrigatório")
	}

	if req.FluxoID <= 0 {
		return apperrors.NewValidationError("ID do fluxo/setor é obrigatório")
	}

	// 2. Validar tamanhos
	if len(req.Descricao) > maxLengthHistoricoObservacao {
		return apperrors.NewValidationError(
			fmt.Sprintf("observação deve ter no máximo %d caracteres", maxLengthHistoricoObservacao),
		)
	}

	return nil
}

// validateBusinessRules valida as regras de negócio
func (s *documentoVendaHistoricoService) validateBusinessRules(ctx context.Context, req *dto.DocumentoVendaHistoricoRequest, excludeItem int) error {
	// 1. Validar se o documento de venda existe
	doc, err := s.documentoVendaService.GetByID(req.DocumentoVendaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewNotFoundError(fmt.Sprintf("documento de venda com ID %d não encontrado", req.DocumentoVendaID))
		}
		return apperrors.NewInternalError("erro ao verificar documento de venda", err)
	}

	// 2. Validar se o documento está ativo (não deletado)
	if doc.IsDeleted() {
		return apperrors.NewValidationError("não é possível adicionar histórico a um documento deletado")
	}

	// 4. Validar se o fluxo/setor existe
	// fluxo, err := s.fluxoSetorService.GetByID(ctx, req.FluxoID)
	// if err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		return apperrors.NewNotFoundError(fmt.Sprintf("fluxo/setor com ID %d não encontrado", req.FluxoID))
	// 	}
	// 	return apperrors.NewInternalError("erro ao verificar fluxo/setor", err)
	// }

	// if !fluxo.IsActive() {
	// 	return apperrors.NewValidationError("fluxo/setor inativo não pode ser usado")
	// }

	// 5. Validar se o item já existe (para atualização)
	if excludeItem > 0 {
		existing, err := s.repo.FindByID(ctx, req.DocumentoVendaID, excludeItem)
		if err != nil {
			return apperrors.NewInternalError("erro ao verificar histórico existente", err)
		}
		if existing == nil {
			return apperrors.NewNotFoundError(fmt.Sprintf("histórico %d-%d não encontrado", req.DocumentoVendaID, excludeItem))
		}
	}

	// 6. Validar se não está tentando criar um registro duplicado
	if excludeItem == 0 {
		existing, err := s.repo.FindByID(ctx, req.DocumentoVendaID, req.Item)
		if err != nil {
			return apperrors.NewInternalError("erro ao verificar duplicidade", err)
		}
		if existing != nil {
			return apperrors.NewValidationError(
				fmt.Sprintf("já existe histórico com item %d para este documento", req.Item),
			)
		}
	}

	return nil
}

// validateCreateValid valida dados para criação
func (s *documentoVendaHistoricoService) validateCreateValid(ctx context.Context, req *dto.DocumentoVendaHistoricoRequest) error {
	// 1. Validações básicas do DTO
	if err := s.validateRequest(req); err != nil {
		return err
	}

	// 2. Regras de negócio (excludeItem = 0 para criação)
	if err := s.validateBusinessRules(ctx, req, 0); err != nil {
		return err
	}

	return nil
}

// validateUpdateValid valida dados para atualização
func (s *documentoVendaHistoricoService) validateUpdateValid(ctx context.Context, documentoVendaID int, item int, req *dto.DocumentoVendaHistoricoRequest) error {
	// 1. Validar se o ID é válido
	if documentoVendaID <= 0 {
		return apperrors.NewValidationError("ID do documento de venda inválido")
	}
	if item <= 0 {
		return apperrors.NewValidationError("item inválido")
	}

	// 2. Validações básicas do DTO
	if err := s.validateRequest(req); err != nil {
		return err
	}

	// 3. Regras de negócio (excludeItem = item para atualização)
	if err := s.validateBusinessRules(ctx, req, item); err != nil {
		return err
	}

	return nil
}

// ============================================================
// MÉTODOS PÚBLICOS (CRUD)
// ============================================================

// Create cria um novo histórico
func (s *documentoVendaHistoricoService) Create(ctx context.Context, req *dto.DocumentoVendaHistoricoRequest) (*models.DocumentoVendaHistorico, error) {
	// Validar
	if err := s.validateCreateValid(ctx, req); err != nil {
		return nil, err
	}

	// Criar modelo
	historico := &models.DocumentoVendaHistorico{
		DocumentoVendaID: req.DocumentoVendaID,
		Item:             req.Item,
		UsuarioID:        req.UsuarioID,
		FluxoID:          req.FluxoID,
		Descricao:        req.Descricao,
	}

	// Se a data foi fornecida, usa ela, senão usa agora
	d, err := utils.ParseDateTime(req.DataHistorico)
	if err != nil {
		return nil, apperrors.NewValidationError("data do histórico inválida")
	}
	historico.DataHistorico = d

	// Salvar
	if err := s.repo.Create(ctx, historico); err != nil {
		return nil, err
	}

	return historico, nil
}

// GetByID busca um histórico pelo ID composto (documentoVendaID + item)
func (s *documentoVendaHistoricoService) GetByID(ctx context.Context, documentoVendaID int, item int) (*models.DocumentoVendaHistorico, error) {
	if documentoVendaID <= 0 {
		return nil, apperrors.NewValidationError("ID do documento de venda inválido")
	}
	if item <= 0 {
		return nil, apperrors.NewValidationError("item inválido")
	}

	historico, err := s.repo.FindByID(ctx, documentoVendaID, item)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("histórico %d-%d não encontrado", documentoVendaID, item))
		}
		return nil, apperrors.NewInternalError("erro ao buscar histórico", err)
	}

	if historico == nil {
		return nil, apperrors.NewNotFoundError(fmt.Sprintf("histórico %d-%d não encontrado", documentoVendaID, item))
	}

	return historico, nil
}

// GetByDocumentoVendaID busca todos os históricos de um documento
func (s *documentoVendaHistoricoService) GetByDocumentoVendaID(ctx context.Context, documentoVendaID int) ([]models.DocumentoVendaHistorico, error) {
	if documentoVendaID <= 0 {
		return nil, apperrors.NewValidationError("ID do documento de venda inválido")
	}

	// Validar se o documento existe
	_, err := s.documentoVendaService.GetByID(documentoVendaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("documento de venda com ID %d não encontrado", documentoVendaID))
		}
		return nil, err
	}

	historicos, err := s.repo.FindByDocumentVendaID(ctx, documentoVendaID)
	if err != nil {
		return nil, apperrors.NewInternalError("erro ao buscar históricos do documento", err)
	}

	return historicos, nil
}

// GetLastByDocumentoVendaID busca o último histórico de um documento
func (s *documentoVendaHistoricoService) GetLastByDocumentoVendaID(ctx context.Context, documentoVendaID int) (*models.DocumentoVendaHistorico, error) {
	if documentoVendaID <= 0 {
		return nil, apperrors.NewValidationError("ID do documento de venda inválido")
	}

	// Validar se o documento existe
	_, err := s.documentoVendaService.GetByID(documentoVendaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("documento de venda com ID %d não encontrado", documentoVendaID))
		}
		return nil, err
	}

	historico, err := s.repo.GetLastByDocumentVendaID(ctx, documentoVendaID)
	if err != nil {
		return nil, apperrors.NewInternalError("erro ao buscar último histórico", err)
	}

	if historico == nil {
		return nil, apperrors.NewNotFoundError(fmt.Sprintf("histórico %d não encontrado", documentoVendaID))
	}

	return historico, nil
}

// List lista históricos com paginação e filtros
func (s *documentoVendaHistoricoService) List(ctx context.Context, filter *dto.DocumentoVendaHistoricoFilter) ([]models.DocumentoVendaHistorico, int64, error) {
	if filter == nil {
		filter = &dto.DocumentoVendaHistoricoFilter{}
	}

	// Validar filtros
	if filter.DocumentoVendaID != nil && *filter.DocumentoVendaID <= 0 {
		return nil, 0, apperrors.NewValidationError("ID do documento de venda inválido")
	}
	if filter.UsuarioID != nil && *filter.UsuarioID <= 0 {
		return nil, 0, apperrors.NewValidationError("ID do usuário inválido")
	}
	if filter.FluxoID != nil && *filter.FluxoID <= 0 {
		return nil, 0, apperrors.NewValidationError("ID do fluxo/setor inválido")
	}

	historicos, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, 0, apperrors.NewInternalError("erro ao listar históricos", err)
	}

	return historicos, total, nil
}

// Update atualiza um histórico existente
func (s *documentoVendaHistoricoService) Update(ctx context.Context, documentoVendaID int, item int, req *dto.DocumentoVendaHistoricoRequest) (*models.DocumentoVendaHistorico, error) {
	// Validar
	if err := s.validateUpdateValid(ctx, documentoVendaID, item, req); err != nil {
		return nil, err
	}

	// Buscar histórico existente
	existing, err := s.repo.ExistByID(ctx, documentoVendaID, item)
	if err != nil {
		return nil, err
	}
	if !existing {
		return nil, apperrors.NewNotFoundError(fmt.Sprintf("histórico %d-%d não encontrado", documentoVendaID, item))
	}

	ddvh, err := req.ToModel()
	if err != nil {
		return nil, apperrors.NewBadRequestError(err.Error())
	}

	// Salvar
	if err := s.repo.Update(ctx, ddvh); err != nil {
		return nil, err
	}

	return ddvh, nil
}

// Delete realiza a exclusão física de um histórico
func (s *documentoVendaHistoricoService) Delete(ctx context.Context, documentoVendaID int, item int) error {
	if documentoVendaID <= 0 {
		return apperrors.NewValidationError("ID do documento de venda inválido")
	}
	if item <= 0 {
		return apperrors.NewValidationError("item inválido")
	}

	// Verificar se o histórico existe
	historico, err := s.repo.FindByID(ctx, documentoVendaID, item)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewNotFoundError(fmt.Sprintf("histórico %d-%d não encontrado", documentoVendaID, item))
		}
		return apperrors.NewInternalError("erro ao buscar histórico para exclusão", err)
	}

	if historico == nil {
		return apperrors.NewNotFoundError(fmt.Sprintf("histórico %d-%d não encontrado", documentoVendaID, item))
	}

	// Realizar exclusão
	if err := s.repo.Delete(ctx, documentoVendaID, item); err != nil {
		return err
	}

	return nil
}

// SoftDelete realiza a exclusão lógica de um histórico
func (s *documentoVendaHistoricoService) SoftDelete(ctx context.Context, documentoVendaID int, item int) error {
	if documentoVendaID <= 0 {
		return apperrors.NewValidationError("ID do documento de venda inválido")
	}
	if item <= 0 {
		return apperrors.NewValidationError("item inválido")
	}

	// Verificar se o histórico existe
	historico, err := s.repo.FindByID(ctx, documentoVendaID, item)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewNotFoundError(fmt.Sprintf("histórico %d-%d não encontrado", documentoVendaID, item))
		}
		return apperrors.NewInternalError("erro ao buscar histórico para exclusão lógica", err)
	}

	if historico == nil {
		return apperrors.NewNotFoundError(fmt.Sprintf("histórico %d-%d não encontrado", documentoVendaID, item))
	}

	// Realizar exclusão lógica
	if err := s.repo.Delete(ctx, documentoVendaID, item); err != nil {
		return err
	}

	return nil
}
