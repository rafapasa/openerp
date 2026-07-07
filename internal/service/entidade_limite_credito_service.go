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

type EntidadeLimiteCreditoService struct {
	repo *repository.EntidadeLimiteCreditoRepository
}

// ============================================================
// CONSTRUCTOR
// ============================================================

func NewEntidadeLimiteCreditoService(db *gorm.DB) *EntidadeLimiteCreditoService {
	return &EntidadeLimiteCreditoService{
		repo: repository.NewEntidadeLimiteCreditoRepository(db),
	}
}

// ============================================================
// MÉTODOS PRINCIPAIS (CRUD)
// ============================================================

// Create cria um novo limite de crédito.
func (s *EntidadeLimiteCreditoService) Create(req *dto.EntidadeLimiteCreditoRequest) (*models.EntidadeLimiteCredito, error) {
	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.repo.ExistsByDescricao(descricao, 0)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar descrição: %w", err)
	}
	if exists {
		return nil, errors.New("já existe um limite de crédito com esta descrição")
	}

	limite := &models.EntidadeLimiteCredito{
		Descricao: &descricao,
		Valor:     req.Valor,
		CreatedBy: req.CreatedBy,
		UpdatedBy: req.UpdatedBy,
	}

	if err := s.repo.Create(limite); err != nil {
		return nil, fmt.Errorf("erro ao criar limite de crédito: %w", err)
	}

	return limite, nil
}

// GetByID busca um limite de crédito por ID.
func (s *EntidadeLimiteCreditoService) GetByID(id int) (*models.EntidadeLimiteCredito, error) {
	return s.repo.FindByID(id)
}

// Update atualiza um limite de crédito.
func (s *EntidadeLimiteCreditoService) Update(id int, req *dto.EntidadeLimiteCreditoRequest) (*models.EntidadeLimiteCredito, error) {
	limite, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.repo.ExistsByDescricao(descricao, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar descrição: %w", err)
	}
	if exists {
		return nil, errors.New("já existe um limite de crédito com esta descrição")
	}

	limite.Descricao = &descricao
	limite.Valor = req.Valor
	limite.UpdatedBy = req.UpdatedBy

	if err := s.repo.Update(id, limite); err != nil {
		return nil, fmt.Errorf("erro ao atualizar limite de crédito: %w", err)
	}

	return limite, nil
}

// Delete exclui um limite de crédito.
func (s *EntidadeLimiteCreditoService) Delete(id int) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// TODO: Adicionar verificação se o limite de crédito está em uso por alguma entidade.

	return s.repo.Delete(id)
}

// List lista todos os limites de crédito.
func (s *EntidadeLimiteCreditoService) List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeLimiteCredito, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(limit, offset, filters)
}
