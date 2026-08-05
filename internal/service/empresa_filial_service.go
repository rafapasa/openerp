package service

import (
	"context"
	"fmt"

	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/utils"
)

// EmpresaService define a interface para o serviço de Empresa (para validação de dependência).
type EmpresaService interface {
	GetByID(ctx context.Context, id int) (*models.Empresa, error)
}

// EmpresaFilialService define os métodos públicos para o serviço de filial de empresa.
type EmpresaFilialService interface {
	Create(ctx context.Context, req *dto.EmpresaFilialRequest) (*dto.EmpresaFilialResponse, error)
	GetByID(ctx context.Context, id int) (*dto.EmpresaFilialResponse, error)
	Update(ctx context.Context, id int, req *dto.EmpresaFilialRequest) (*dto.EmpresaFilialResponse, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]dto.EmpresaFilialResponse, int64, error)
	FindByID(ctx context.Context, id int) (*models.EmpresaFilial, error) // Para uso interno por outros serviços
}

type empresaFilialService struct {
	repo          repository.EmpresaFilialRepository
	empresaService EmpresaService // Dependência para validar a empresa
}

// NewEmpresaFilialService cria uma nova instância de EmpresaFilialService.
func NewEmpresaFilialService(repo repository.EmpresaFilialRepository, empService EmpresaService) EmpresaFilialService {
	return &empresaFilialService{
		repo:          repo,
		empresaService: empService,
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO (PRIVADOS)
// ============================================================

func (s *empresaFilialService) validateEmpresaFilial(ctx context.Context, id int, req *dto.EmpresaFilialRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	// 1. Validar se a Empresa existe e está ativa
	empresa, err := s.empresaService.GetByID(ctx, req.EmpresaID)
	if err != nil {
		return err // Retorna NotFoundError ou InternalError do EmpresaService
	}
	if !empresa.IsActive() {
		return apperrors.NewValidationError(fmt.Sprintf("Empresa com ID %d está inativa.", req.EmpresaID))
	}

	// 2. Validar unicidade do número da filial por empresa
	exists, err := s.repo.ExistsByNumero(req.Numero, req.EmpresaID, id)
	if err != nil {
		return err
	}
	if exists {
		return apperrors.NewConflictError(fmt.Sprintf("Já existe uma filial com o número '%d' para esta empresa.", req.Numero))
	}

	return nil
}

// mapModelToResponse mapeia um modelo EmpresaFilial para um DTO de resposta.
func (s *empresaFilialService) mapModelToResponse(filial *models.EmpresaFilial) (*dto.EmpresaFilialResponse, error) {
	resp := &dto.EmpresaFilialResponse{}
	resp.FromModel(filial) // Usa o método FromModel do DTO
	return resp, nil
}

// ============================================================
// MÉTODOS PÚBLICOS (CRUD)
// ============================================================

// Create cria uma nova filial de empresa.
func (s *empresaFilialService) Create(ctx context.Context, req *dto.EmpresaFilialRequest) (*dto.EmpresaFilialResponse, error) {
	if err := s.validateEmpresaFilial(ctx, 0, req); err != nil {
		return nil, err
	}

	filial, err := req.ToModel()
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear DTO para modelo.", err)
	}

	if err := s.repo.Create(filial); err != nil {
		return nil, err
	}

	// Recarregar para popular relacionamentos (Empresa)
	createdFilial, err := s.repo.FindByID(filial.ID)
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao buscar filial de empresa criada.", err)
	}

	return s.mapModelToResponse(createdFilial)
}

// GetByID busca uma filial de empresa pelo ID.
func (s *empresaFilialService) GetByID(ctx context.Context, id int) (*dto.EmpresaFilialResponse, error) {
	filial, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if filial.IsDeleted() {
		return nil, apperrors.NewNotFoundError(fmt.Sprintf("Filial de empresa com ID %d não encontrada.", id))
	}
	return s.mapModelToResponse(filial)
}

// FindByID busca uma filial de empresa pelo ID (retorna o modelo).
func (s *empresaFilialService) FindByID(ctx context.Context, id int) (*models.EmpresaFilial, error) {
	return s.repo.FindByID(id)
}

// Update atualiza uma filial de empresa existente.
func (s *empresaFilialService) Update(ctx context.Context, id int, req *dto.EmpresaFilialRequest) (*dto.EmpresaFilialResponse, error) {
	existingFilial, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existingFilial.IsDeleted() {
		return nil, apperrors.NewNotFoundError(fmt.Sprintf("Filial de empresa com ID %d não encontrada para atualização.", id))
	}

	if err := s.validateEmpresaFilial(ctx, id, req); err != nil {
		return nil, err
	}

	if err := utils.MapToModel(req, existingFilial); err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear DTO para modelo existente.", err)
	}

	if err := s.repo.Update(existingFilial); err != nil {
		return nil, err
	}

	// Recarregar para popular relacionamentos (Empresa)
	updatedFilial, err := s.repo.FindByID(existingFilial.ID)
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao buscar filial de empresa atualizada.", err)
	}

	return s.mapModelToResponse(updatedFilial)
}

// Delete realiza a exclusão lógica de uma filial de empresa.
func (s *empresaFilialService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return apperrors.NewValidationError("ID da filial de empresa inválido.")
	}

	// Verificar dependências
	hasDependents, err := s.repo.HasDependentRecords(id)
	if err != nil {
		return apperrors.NewInternalError("Erro ao verificar dependências da filial de empresa.", err)
	}
	if hasDependents {
		return apperrors.NewConflictError(fmt.Sprintf("Não é possível excluir a filial de empresa com ID %d pois existem registros associados.", id))
	}

	existingFilial, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if existingFilial.IsDeleted() {
		return apperrors.NewNotFoundError(fmt.Sprintf("Filial de empresa com ID %d não encontrada para exclusão.", id))
	}

	if err := s.repo.Delete(id); err != nil {
		return apperrors.NewInternalError("Erro ao excluir filial de empresa.", err)
	}
	return nil
}

// List lista filiais de empresa com paginação e filtros.
func (s *empresaFilialService) List(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]dto.EmpresaFilialResponse, int64, error) {
	if limit < 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	filiais, total, err := s.repo.List(limit, offset, filters)
	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro ao listar filiais de empresa.", err)
	}

	var responses []dto.EmpresaFilialResponse
	for _, filial := range filiais {
		resp, err := s.mapModelToResponse(&filial)
		if err != nil {
			fmt.Printf("Erro ao mapear filial de empresa %d para resposta: %v\n", filial.ID, err)
			continue
		}
		responses = append(responses, *resp)
	}

	return responses, total, nil
}