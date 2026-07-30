package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
)

// GrupoEntidadeService define os métodos públicos para o serviço de grupo de entidade.
type GrupoEntidadeService interface {
	Create(req *dto.GrupoEntidadeRequest) (*models.GrupoEntidade, error)
	GetByID(id int) (*models.GrupoEntidade, error)
	Update(id int, req *dto.GrupoEntidadeRequest) (*models.GrupoEntidade, error)
	Delete(id int) error
	List(limit, offset int, filters map[string]interface{}) ([]models.GrupoEntidade, int64, error)
}

// ============================================================
// TYPES
// ============================================================

type grupoEntidadeService struct {
	gpeRepo repository.GrupoEntidadeRepository
}

// ============================================================
// CONSTRUCTOR
// ============================================================

func NewGrupoEntidadeService(gpeRepo repository.GrupoEntidadeRepository) GrupoEntidadeService {
	return &grupoEntidadeService{
		gpeRepo: gpeRepo,
	}
}

// ============================================================
// MÉTODOS PRINCIPAIS (CRUD)
// ============================================================

// Create cria um novo grupo de entidade
func (s *grupoEntidadeService) Create(req *dto.GrupoEntidadeRequest) (*models.GrupoEntidade, error) {
	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.gpeRepo.ExistsByDescricao(descricao, 0)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar nome: %w", err)
	}
	if exists {
		return nil, errors.New("já existe um grupo de entidade com este nome")
	}

	grupo := &models.GrupoEntidade{
		Descricao: descricao,
		CreatedBy: req.CreatedBy,
		UpdatedBy: req.UpdatedBy,
	}

	if err := s.gpeRepo.Create(grupo); err != nil {
		return nil, fmt.Errorf("erro ao criar grupo de entidade: %w", err)
	}

	return grupo, nil
}

// GetByID busca um grupo de entidade por ID
func (s *grupoEntidadeService) GetByID(id int) (*models.GrupoEntidade, error) {
	return s.gpeRepo.FindByID(id)
}

// Update atualiza um grupo de entidade
func (s *grupoEntidadeService) Update(id int, req *dto.GrupoEntidadeRequest) (*models.GrupoEntidade, error) {
	grupo, err := s.gpeRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.gpeRepo.ExistsByDescricao(descricao, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar nome: %w", err)
	}
	if exists {
		return nil, errors.New("já existe um grupo de entidade com este nome")
	}

	grupo.Descricao = descricao
	grupo.UpdatedBy = req.UpdatedBy

	if err := s.gpeRepo.Update(id, grupo); err != nil {
		return nil, fmt.Errorf("erro ao atualizar grupo de entidade: %w", err)
	}

	return grupo, nil
}

// Delete exclui um grupo de entidade
func (s *grupoEntidadeService) Delete(id int) error {
	_, err := s.gpeRepo.FindByID(id)
	if err != nil {
		return err
	}

	// TODO: Adicionar verificação se o grupo está em uso por alguma entidade

	return s.gpeRepo.Delete(id)
}

// List lista todos os grupos de entidade
func (s *grupoEntidadeService) List(limit, offset int, filters map[string]interface{}) ([]models.GrupoEntidade, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.gpeRepo.List(limit, offset, filters)
}
