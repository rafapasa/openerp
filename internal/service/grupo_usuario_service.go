package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/utils"
)

// GrupoUsuarioService define os métodos públicos para o serviço de grupo de usuário.
type GrupoUsuarioService interface {
	Create(ctx context.Context, req *dto.GrupoUsuarioRequest) (*dto.GrupoUsuarioResponse, error)
	GetByID(ctx context.Context, id int) (*dto.GrupoUsuarioResponse, error)
	Update(ctx context.Context, id int, req *dto.GrupoUsuarioRequest) (*dto.GrupoUsuarioResponse, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]dto.GrupoUsuarioResponse, int64, error)
	FindByID(ctx context.Context, id int) (*models.GrupoUsuario, error) // Para uso interno por outros serviços
}

type grupoUsuarioService struct {
	repo repository.GrupoUsuarioRepository
}

// NewGrupoUsuarioService cria uma nova instância de GrupoUsuarioService.
func NewGrupoUsuarioService(repo repository.GrupoUsuarioRepository) GrupoUsuarioService {
	return &grupoUsuarioService{
		repo: repo,
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO (PRIVADOS)
// ============================================================

func (s *grupoUsuarioService) validateGrupoUsuario(id int, req *dto.GrupoUsuarioRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	// Validar unicidade da descrição por filial
	exists, err := s.repo.ExistsByDescricao(strings.TrimSpace(req.Descricao), req.EmpresaFilialID, id)
	if err != nil {
		return err
	}
	if exists {
		return apperrors.NewConflictError(fmt.Sprintf("Já existe um grupo de usuário com a descrição '%s' para esta filial.", req.Descricao))
	}

	return nil
}

// mapModelToResponse mapeia um modelo GrupoUsuario para um DTO de resposta.
func (s *grupoUsuarioService) mapModelToResponse(grupo *models.GrupoUsuario) (*dto.GrupoUsuarioResponse, error) {
	resp := &dto.GrupoUsuarioResponse{}
	resp.FromModel(grupo) // Usa o método FromModel do DTO
	return resp, nil
}

// ============================================================
// MÉTODOS PÚBLICOS (CRUD)
// ============================================================

// Create cria um novo grupo de usuário.
func (s *grupoUsuarioService) Create(ctx context.Context, req *dto.GrupoUsuarioRequest) (*dto.GrupoUsuarioResponse, error) {
	if err := s.validateGrupoUsuario(0, req); err != nil {
		return nil, err
	}

	grupo, err := req.ToModel()
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear DTO para modelo.", err)
	}

	if err := s.repo.Create(grupo); err != nil {
		return nil, err
	}

	return s.mapModelToResponse(grupo)
}

// GetByID busca um grupo de usuário pelo ID.
func (s *grupoUsuarioService) GetByID(ctx context.Context, id int) (*dto.GrupoUsuarioResponse, error) {
	grupo, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if grupo.IsDeleted() {
		return nil, apperrors.NewNotFoundError(fmt.Sprintf("Grupo de usuário com ID %d não encontrado.", id))
	}
	return s.mapModelToResponse(grupo)
}

// FindByID busca um grupo de usuário pelo ID (retorna o modelo).
func (s *grupoUsuarioService) FindByID(ctx context.Context, id int) (*models.GrupoUsuario, error) {
	return s.repo.FindByID(id)
}

// Update atualiza um grupo de usuário existente.
func (s *grupoUsuarioService) Update(ctx context.Context, id int, req *dto.GrupoUsuarioRequest) (*dto.GrupoUsuarioResponse, error) {
	existingGrupo, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existingGrupo.IsDeleted() {
		return nil, apperrors.NewNotFoundError(fmt.Sprintf("Grupo de usuário com ID %d não encontrado para atualização.", id))
	}

	if err := s.validateGrupoUsuario(id, req); err != nil {
		return nil, err
	}

	if err := utils.MapToModel(req, existingGrupo); err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear DTO para modelo existente.", err)
	}

	if err := s.repo.Update(existingGrupo); err != nil {
		return nil, err
	}

	return s.mapModelToResponse(existingGrupo)
}

// Delete realiza a exclusão lógica de um grupo de usuário.
func (s *grupoUsuarioService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return apperrors.NewValidationError("ID do grupo de usuário inválido.")
	}

	// Verificar dependências
	hasDependents, err := s.repo.HasDependentRecords(id)
	if err != nil {
		return apperrors.NewInternalError("Erro ao verificar dependências do grupo de usuário.", err)
	}
	if hasDependents {
		return apperrors.NewConflictError(fmt.Sprintf("Não é possível excluir o grupo de usuário com ID %d pois existem usuários associados.", id))
	}

	existingGrupo, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if existingGrupo.IsDeleted() {
		return apperrors.NewNotFoundError(fmt.Sprintf("Grupo de usuário com ID %d não encontrado para exclusão.", id))
	}

	if err := s.repo.Delete(id); err != nil {
		return apperrors.NewInternalError("Erro ao excluir grupo de usuário.", err)
	}
	return nil
}

// List lista grupos de usuário com paginação e filtros.
func (s *grupoUsuarioService) List(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]dto.GrupoUsuarioResponse, int64, error) {
	if limit < 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	grupos, total, err := s.repo.List(limit, offset, filters)
	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro ao listar grupos de usuário.", err)
	}

	var responses []dto.GrupoUsuarioResponse
	for _, grupo := range grupos {
		resp, err := s.mapModelToResponse(&grupo)
		if err != nil {
			fmt.Printf("Erro ao mapear grupo de usuário %d para resposta: %v\n", grupo.ID, err)
			continue
		}
		responses = append(responses, *resp)
	}

	return responses, total, nil
}
