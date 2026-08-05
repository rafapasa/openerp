package service

import (
	"context"
	"strings"

	"github.com/openerp/backend/internal/apperrors"

	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"gorm.io/gorm"
)

// EntidadeLimiteCreditoService define os métodos públicos para o serviço de limite de crédito de entidade.
type EntidadeLimiteCreditoService interface {
	Create(ctx context.Context, req *dto.EntidadeLimiteCreditoRequest) (*models.EntidadeLimiteCredito, error)
	GetByID(ctx context.Context, id int) (*models.EntidadeLimiteCredito, error)
	Update(ctx context.Context, id int, req *dto.EntidadeLimiteCreditoRequest) (*models.EntidadeLimiteCredito, error)
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
func (s *entidadeLimiteCreditoService) Create(ctx context.Context, req *dto.EntidadeLimiteCreditoRequest) (*models.EntidadeLimiteCredito, error) {
	descricao := strings.TrimSpace(req.Descricao) // Context not used in ExistsByDescricao
	exists, err := s.limiteRepo.ExistsByDescricao(descricao, 0)
	if err != nil { //
		return nil, apperrors.NewInternalError("Erro ao verificar descrição.", err) //
	}
	if exists {
		return nil, apperrors.NewConflictError("Já existe um limite de crédito com esta descrição.") //
	}

	limite := &models.EntidadeLimiteCredito{
		Descricao: &descricao,
		Valor:     req.Valor,
		CreatedBy: req.CreatedBy,
		UpdatedBy: req.UpdatedBy,
	}

	if err := s.limiteRepo.Create(limite); err != nil {
		return nil, err
	}

	return limite, nil
}

// GetByID busca um limite de crédito por ID.
func (s *entidadeLimiteCreditoService) GetByID(ctx context.Context, id int) (*models.EntidadeLimiteCredito, error) { //
	return s.limiteRepo.FindByID(id) // Context not used in FindByID
}

// Update atualiza um limite de crédito.
func (s *entidadeLimiteCreditoService) Update(ctx context.Context, id int, req *dto.EntidadeLimiteCreditoRequest) (*models.EntidadeLimiteCredito, error) {
	limite, err := s.limiteRepo.FindByID(id) // Context not used in FindByID
	if err != nil {
		return nil, err //
	}

	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.limiteRepo.ExistsByDescricao(descricao, id)
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao verificar descrição.", err) //
	}
	if exists {
		return nil, apperrors.NewConflictError("Já existe um limite de crédito com esta descrição.") //
	}

	limite.Descricao = &descricao
	limite.Valor = req.Valor
	limite.UpdatedBy = req.UpdatedBy

	if err := s.limiteRepo.Update(id, limite); err != nil {
		return nil, apperrors.NewInternalError("Erro ao atualizar limite de crédito.", err) //
	}

	return limite, nil
}

// Delete exclui um limite de crédito.
func (s *entidadeLimiteCreditoService) Delete(id int) error {
	if _, err := s.limiteRepo.FindByID(id); err != nil { //
		return err //
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
