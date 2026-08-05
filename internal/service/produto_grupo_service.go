package service

import (
	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/utils"
)

// ProdutoGrupoService define os métodos públicos para o serviço de grupo de produto.
type ProdutoGrupoService interface {
	List(limit, offset int, filters map[string]any) ([]models.ProdutoGrupo, int64, error)
	GetByID(id int) (*models.ProdutoGrupo, error)
	Create(req *dto.ProdutoGrupoRequest) (*models.ProdutoGrupo, error)
	Update(id int, req *dto.ProdutoGrupoRequest) (*models.ProdutoGrupo, error)
	Delete(id int) error
}

type produtoGrupoService struct {
	progRepo repository.ProdutoGrupoRepository
}

func NewProdutoGrupoService(progRepo repository.ProdutoGrupoRepository) ProdutoGrupoService {
	return &produtoGrupoService{
		progRepo: progRepo,
	}
}

func (s *produtoGrupoService) List(limit, offset int, filters map[string]any) ([]models.ProdutoGrupo, int64, error) {
	return s.progRepo.List(limit, offset, filters)
}

func (s *produtoGrupoService) GetByID(id int) (*models.ProdutoGrupo, error) {
	return s.progRepo.FindByID(id)
}

func (s *produtoGrupoService) Create(req *dto.ProdutoGrupoRequest) (*models.ProdutoGrupo, error) {
	produtoGrupo, err := req.ToModel()
	if err != nil { //
		return nil, apperrors.NewInternalError("Erro ao converter dados do grupo de produto.", err) //
	}

	// Validar o modelo antes de persistir
	if err := s.validateProdutoGrupo(req); err != nil {
		return nil, err
	}

	result := s.progRepo.Create(produtoGrupo)
	if result != nil {
		return nil, apperrors.NewInternalError("Erro ao criar grupo de produto.", result) //
	}

	return produtoGrupo, nil
}

func (s *produtoGrupoService) Update(id int, req *dto.ProdutoGrupoRequest) (*models.ProdutoGrupo, error) {
	produtoGrupo, err := req.ToModel()
	if err != nil { //
		return nil, apperrors.NewInternalError("Erro ao converter dados do grupo de produto.", err) //
	}

	if err := s.validateProdutoGrupo(req); err != nil {
		return nil, err
	}
	result := s.progRepo.Update(id, produtoGrupo)
	if result != nil {
		return nil, apperrors.NewInternalError("Erro ao atualizar grupo de produto.", result) //
	}
	return produtoGrupo, nil
}

func (s *produtoGrupoService) Delete(id int) error {
	// TODO: Implementar a checagem de dependencias
	// FIXME: Descrição do bug a ser corrigido
	// OPTIMIZE: Descrição da otimização
	// SECURITY: Descrição da melhoria de segurança
	return s.progRepo.Delete(id)
}

func (s *produtoGrupoService) validateProdutoGrupo(req *dto.ProdutoGrupoRequest) error {
	if err := utils.ValidateMandatoryFields(req); err != nil {
		return err
	}

	return nil
}
