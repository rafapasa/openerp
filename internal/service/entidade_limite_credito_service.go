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

// EntidadeLimiteCreditoService define os métodos públicos para o serviço de limite de crédito de entidade.
type EntidadeLimiteCreditoService interface {
	Create(req *dto.EntidadeLimiteCreditoRequest) (*models.EntidadeLimiteCredito, error)
	GetByID(id int) (*models.EntidadeLimiteCredito, error)
	Update(id int, req *dto.EntidadeLimiteCreditoRequest) (*models.EntidadeLimiteCredito, error)
	Delete(id int) error
	List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeLimiteCredito, int64, error)
}

// ============================================================
// TYPES
// ============================================================

type entidadeLimiteCreditoService struct {
	limiteRepo repository.EntidadeLimiteCreditoRepository // Mantém o ponteiro para o repositório concreto
}

// ============================================================
// CONSTRUCTOR
// ============================================================

func NewEntidadeLimiteCreditoService(db *gorm.DB, elcRepo repository.EntidadeLimiteCreditoRepository) EntidadeLimiteCreditoService {
	return &entidadeLimiteCreditoService{ // Retorna a implementação concreta que satisfaz a interface
		limiteRepo: elcRepo,
	}
}

// ============================================================
// MÉTODOS PRINCIPAIS (CRUD)
// ============================================================

// Create cria um novo limite de crédito.
func (s *entidadeLimiteCreditoService) Create(req *dto.EntidadeLimiteCreditoRequest) (*models.EntidadeLimiteCredito, error) {
	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.limiteRepo.ExistsByDescricao(descricao, 0)
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

	if err := s.limiteRepo.Create(limite); err != nil {
		return nil, fmt.Errorf("erro ao criar limite de crédito: %w", err)
	}

	return limite, nil
}

// GetByID busca um limite de crédito por ID.
func (s *entidadeLimiteCreditoService) GetByID(id int) (*models.EntidadeLimiteCredito, error) {
	return s.limiteRepo.FindByID(id)
}

// Update atualiza um limite de crédito.
func (s *entidadeLimiteCreditoService) Update(id int, req *dto.EntidadeLimiteCreditoRequest) (*models.EntidadeLimiteCredito, error) {
	limite, err := s.limiteRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.limiteRepo.ExistsByDescricao(descricao, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar descrição: %w", err)
	}
	if exists {
		return nil, errors.New("já existe um limite de crédito com esta descrição")
	}

	limite.Descricao = &descricao
	limite.Valor = req.Valor
	limite.UpdatedBy = req.UpdatedBy

	if err := s.limiteRepo.Update(id, limite); err != nil {
		return nil, fmt.Errorf("erro ao atualizar limite de crédito: %w", err)
	}

	return limite, nil
}

// Delete exclui um limite de crédito.
func (s *entidadeLimiteCreditoService) Delete(id int) error {
	_, err := s.limiteRepo.FindByID(id)
	if err != nil {
		return err
	}

	// TODO: Adicionar verificação se o limite de crédito está em uso por alguma entidade.

	return s.limiteRepo.Delete(id)
}

// List lista todos os limites de crédito.
func (s *entidadeLimiteCreditoService) List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeLimiteCredito, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.limiteRepo.List(limit, offset, filters)
}
