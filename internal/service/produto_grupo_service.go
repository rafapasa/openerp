package service

import (
	"fmt"

	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/utils"
	"gorm.io/gorm"
)

type ProdutoGrupoService struct {
	produtoGrupoRepository *repository.ProdutoGrupoRepository
}

func NewProdutoGrupoService(db *gorm.DB) *ProdutoGrupoService {
	return &ProdutoGrupoService{
		produtoGrupoRepository: repository.NewProdutoGrupoRepository(db),
	}
}

func (s *ProdutoGrupoService) List(limit, offset int, filters map[string]any) ([]models.ProdutoGrupo, int64, error) {
	return s.produtoGrupoRepository.List(limit, offset, filters)
}

func (s *ProdutoGrupoService) GetByID(id int) (*models.ProdutoGrupo, error) {
	return s.produtoGrupoRepository.FindByID(id)
}

func (s *ProdutoGrupoService) Create(req *dto.ProdutoGrupoRequest) (*models.ProdutoGrupo, error) {
	produtoGrupo, err := req.ToModel()
	if err != nil {
		return nil, fmt.Errorf("erro ao converter dados do grupo de produto: %w", err)
	}

	// Validar o modelo antes de persistir
	if err := s.validateProdutoGrupo(req); err != nil {
		return nil, err
	}

	result := s.produtoGrupoRepository.Create(produtoGrupo)
	if result != nil {
		return nil, fmt.Errorf("erro ao criar grupo de produto: %w", result)
	}

	return produtoGrupo, nil
}

func (s *ProdutoGrupoService) Update(id int, req *dto.ProdutoGrupoRequest) (*models.ProdutoGrupo, error) {
	produtoGrupo, err := req.ToModel()
	if err != nil {
		return nil, fmt.Errorf("erro ao converter dados do grupo de produto: %w", err)
	}

	if err := s.validateProdutoGrupo(req); err != nil {
		return nil, err
	}
	result := s.produtoGrupoRepository.Update(id, produtoGrupo)
	if result != nil {
		return nil, fmt.Errorf("erro ao atualizar grupo de produto: %w", result)
	}
	return produtoGrupo, nil
}

func (s *ProdutoGrupoService) Delete(id int) error {
	// TODO: Implementar a checagem de dependencias
	// FIXME: Descrição do bug a ser corrigido
	// OPTIMIZE: Descrição da otimização
	// SECURITY: Descrição da melhoria de segurança
	return s.produtoGrupoRepository.Delete(id)
}

func (s *ProdutoGrupoService) validateProdutoGrupo(req *dto.ProdutoGrupoRequest) error {
	if err := utils.ValidateMandatoryFields(req); err != nil {
		return err
	}

	return nil
}
