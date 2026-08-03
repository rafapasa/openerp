package service

import (
	apperrors "github.com/openerp/backend/internal/erros"
	"gorm.io/gorm"

	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
)

// EntidadeRegimeTributarioService define os métodos públicos para o serviço de regimes tributários de entidade.
type EntidadeRegimeTributarioService interface {
	Create(req *dto.EntidadeRegimeTributarioRequest) (*models.EntidadeRegimeTributario, error)
	GetByID(entidadeID, item int) (*models.EntidadeRegimeTributario, error)
	Update(entidadeID, item int, req *dto.EntidadeRegimeTributarioRequest) (*models.EntidadeRegimeTributario, error)
	Delete(entidadeID, item int) error
	List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeRegimeTributario, int64, error)
}

// ============================================================
// TYPES
// ============================================================

// EntidadeRegimeTributarioService é o serviço para os regimes da entidade.
type entidadeRegimeTributarioService struct {
	regimeRepo      repository.EntidadeRegimeTributarioRepository
	entidadeService EntidadeService
}

// ============================================================
// CONSTRUCTOR
// ============================================================

// NewEntidadeRegimeTributarioService cria uma nova instância.
func NewEntidadeRegimeTributarioService(db *gorm.DB, entidadeService EntidadeService, regimeRepo repository.EntidadeRegimeTributarioRepository) EntidadeRegimeTributarioService {
	return &entidadeRegimeTributarioService{
		regimeRepo:      regimeRepo,
		entidadeService: entidadeService,
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO (PRIVADOS)
// ============================================================

// validateEntidadeExists verifica se a entidade existe. (Receiver changed)
func (s *entidadeRegimeTributarioService) validateEntidadeExists(entidadeID int) error {
	_, err := s.entidadeService.GetByID(entidadeID)
	if err != nil {
		return err
	}

	return nil
}

// ============================================================
// MÉTODOS PRINCIPAIS (CRUD)
// ============================================================

// Create cria um novo regime para uma entidade. (Receiver changed)
func (s *entidadeRegimeTributarioService) Create(req *dto.EntidadeRegimeTributarioRequest) (*models.EntidadeRegimeTributario, error) { //
	if err := s.validateEntidadeExists(req.EntidadeID); err != nil {
		return nil, err
	}

	regime, err := req.ToModel()
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao converter dados.", err) //
	}

	if err := s.regimeRepo.Create(regime); err != nil {
		return nil, apperrors.NewInternalError("Erro ao criar regime tributário.", err) //
	}

	return regime, nil
}

// GetByID busca um regime tributário específico.
func (s *entidadeRegimeTributarioService) GetByID(entidadeID, item int) (*models.EntidadeRegimeTributario, error) {
	regime, err := s.regimeRepo.FindByID(entidadeID, item)
	if err != nil {
		return nil, err
	}
	return regime, nil
}

// Update atualiza um regime tributário existente.
func (s *entidadeRegimeTributarioService) Update(entidadeID, item int, req *dto.EntidadeRegimeTributarioRequest) (*models.EntidadeRegimeTributario, error) {
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
		return nil, apperrors.NewInternalError("Erro ao converter dados da requisição.", err) //
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
		return nil, apperrors.NewInternalError("Erro ao atualizar regime tributário.", err) //
	}

	return regime, nil
}

// Delete exclui logicamente um regime tributário.
func (s *entidadeRegimeTributarioService) Delete(entidadeID, item int) error {
	if _, err := s.regimeRepo.FindByID(entidadeID, item); err != nil {
		return err
	}

	if err := s.regimeRepo.Delete(entidadeID, item); err != nil {
		return apperrors.NewInternalError("Erro ao excluir regime tributário.", err) //
	}

	return nil
}

// List lista regimes com paginação e filtros.
func (s *entidadeRegimeTributarioService) List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeRegimeTributario, int64, error) {
	if _, ok := filters["ent_id"]; !ok { //
		return nil, 0, apperrors.NewValidationError("O filtro 'ent_id' é obrigatório.") //
	}

	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return s.regimeRepo.List(limit, offset, filters)
}
