package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"gorm.io/gorm"
)

// ============================================================
// TYPES
// ============================================================

type GrupoEntidadeService struct {
	repo *repository.GrupoEntidadeRepository
}

// ============================================================
// CONSTRUCTOR
// ============================================================

func NewGrupoEntidadeService(db *gorm.DB) *GrupoEntidadeService {
	return &GrupoEntidadeService{
		repo: repository.NewGrupoEntidadeRepository(db),
	}
}

// ============================================================
// MÉTODOS PRINCIPAIS (CRUD)
// ============================================================

// Create cria um novo grupo de entidade
func (s *GrupoEntidadeService) Create(req *dto.GrupoEntidadeRequest) (*models.GrupoEntidade, error) {
	nome := strings.TrimSpace(req.Descricao)
	exists, err := s.repo.ExistsByName(nome, 0)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar nome: %w", err)
	}
	if exists {
		return nil, errors.New("já existe um grupo de entidade com este nome")
	}

	grupo := &models.GrupoEntidade{
		Descricao: nome,
		CreatedBy: req.CreatedBy,
		UpdatedBy: req.UpdatedBy,
	}

	if err := s.repo.Create(grupo); err != nil {
		return nil, fmt.Errorf("erro ao criar grupo de entidade: %w", err)
	}

	return grupo, nil
}

// GetByID busca um grupo de entidade por ID
func (s *GrupoEntidadeService) GetByID(id int) (*models.GrupoEntidade, error) {
	return s.repo.FindByID(id)
}

// Update atualiza um grupo de entidade
func (s *GrupoEntidadeService) Update(id int, req *dto.GrupoEntidadeRequest) (*models.GrupoEntidade, error) {
	grupo, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	nome := strings.TrimSpace(req.Descricao)
	exists, err := s.repo.ExistsByName(nome, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar nome: %w", err)
	}
	if exists {
		return nil, errors.New("já existe um grupo de entidade com este nome")
	}

	grupo.Descricao = nome
	grupo.UpdatedBy = req.UpdatedBy

	if err := s.repo.Update(id, grupo); err != nil {
		return nil, fmt.Errorf("erro ao atualizar grupo de entidade: %w", err)
	}

	return grupo, nil
}

// Delete exclui um grupo de entidade
func (s *GrupoEntidadeService) Delete(id int) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// TODO: Adicionar verificação se o grupo está em uso por alguma entidade

	return s.repo.Delete(id)
}

// List lista todos os grupos de entidade
func (s *GrupoEntidadeService) List(limit, offset int, filters map[string]interface{}) ([]models.GrupoEntidade, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(limit, offset, filters)
}
