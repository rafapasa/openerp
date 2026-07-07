package service

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
)

// ============================================================
// TYPES
// ============================================================

// EntidadeRegimeTributarioService é o serviço para os regimes da entidade.
type EntidadeRegimeTributarioService struct {
	regimeRepo   *repository.EntidadeRegimeTributarioRepository
	entidadeRepo *repository.EntidadeRepository
}

// ============================================================
// CONSTRUCTOR
// ============================================================

// NewEntidadeRegimeTributarioService cria uma nova instância.
func NewEntidadeRegimeTributarioService(db *gorm.DB) *EntidadeRegimeTributarioService {
	return &EntidadeRegimeTributarioService{
		regimeRepo:   repository.NewEntidadeRegimeTributarioRepository(db),
		entidadeRepo: repository.NewEntidadeRepository(db),
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO (PRIVADOS)
// ============================================================

// validateEntidadeExists verifica se a entidade existe.
func (s *EntidadeRegimeTributarioService) validateEntidadeExists(entidadeID int) error {
	_, err := s.entidadeRepo.FindByID(entidadeID)
	if err != nil {
		return fmt.Errorf("entidade com ID %d não encontrada", entidadeID)
	}
	return nil
}

// ============================================================
// MÉTODOS PRINCIPAIS (CRUD)
// ============================================================

// Create cria um novo regime para uma entidade.
func (s *EntidadeRegimeTributarioService) Create(req *dto.EntidadeRegimeTributarioRequest) (*models.EntidadeRegimeTributario, error) {
	if err := s.validateEntidadeExists(req.EntidadeID); err != nil {
		return nil, err
	}

	regime, err := req.ToModel()
	if err != nil {
		return nil, fmt.Errorf("erro ao converter dados: %w", err)
	}

	if err := s.regimeRepo.Create(regime); err != nil {
		return nil, fmt.Errorf("erro ao criar regime tributário: %w", err)
	}

	return regime, nil
}

// GetByID busca um regime tributário específico.
func (s *EntidadeRegimeTributarioService) GetByID(entidadeID, item int) (*models.EntidadeRegimeTributario, error) {
	regime, err := s.regimeRepo.FindByID(entidadeID, item)
	if err != nil {
		return nil, err
	}
	return regime, nil
}

// Update atualiza um regime tributário existente.
func (s *EntidadeRegimeTributarioService) Update(entidadeID, item int, req *dto.EntidadeRegimeTributarioRequest) (*models.EntidadeRegimeTributario, error) {
	if err := s.validateEntidadeExists(entidadeID); err != nil {
		return nil, err
	}

	regime, err := s.regimeRepo.FindByID(entidadeID, item)
	if err != nil {
		return nil, err
	}

	// Converte o DTO para um novo modelo para pegar os dados atualizados
	updatedModel, err := req.ToModel()
	if err != nil {
		return nil, fmt.Errorf("erro ao converter dados da requisição: %w", err)
	}

	// Atualiza os campos do modelo existente
	regime.Regime = updatedModel.Regime
	regime.Apuracao = updatedModel.Apuracao
	regime.Data = updatedModel.Data
	regime.RegimeEspecial = updatedModel.RegimeEspecial
	regime.SituacaoTribISS = updatedModel.SituacaoTribISS
	regime.RegimeMunicipal = updatedModel.RegimeMunicipal
	regime.UpdatedBy = req.UpdatedBy

	if err := s.regimeRepo.Update(regime); err != nil {
		return nil, fmt.Errorf("erro ao atualizar regime tributário: %w", err)
	}

	return regime, nil
}

// Delete exclui logicamente um regime tributário.
func (s *EntidadeRegimeTributarioService) Delete(entidadeID, item int) error {
	_, err := s.regimeRepo.FindByID(entidadeID, item)
	if err != nil {
		return err
	}

	if err := s.regimeRepo.Delete(entidadeID, item); err != nil {
		return fmt.Errorf("erro ao excluir regime tributário: %w", err)
	}

	return nil
}

// List lista regimes com paginação e filtros.
func (s *EntidadeRegimeTributarioService) List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeRegimeTributario, int64, error) {
	if _, ok := filters["ent_id"]; !ok {
		return nil, 0, errors.New("o filtro 'ent_id' é obrigatório")
	}

	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return s.regimeRepo.List(limit, offset, filters)
}