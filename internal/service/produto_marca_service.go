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

// ProdutoMarcaService define os métodos públicos para o serviço de marca de produto.
type ProdutoMarcaService interface {
	Create(req *dto.ProdutoMarcaRequest) (*models.ProdutoMarca, error)
	GetByID(id int) (*models.ProdutoMarca, error)
	Update(id int, req *dto.ProdutoMarcaRequest) (*models.ProdutoMarca, error)
	Delete(id int) error
	List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoMarca, int64, error)
}

type produtoMarcaService struct {
	repo *repository.ProdutoMarcaRepository
}

func NewProdutoMarcaService(db *gorm.DB) ProdutoMarcaService {
	return &produtoMarcaService{
		repo: repository.NewProdutoMarcaRepository(db),
	}
}

func (s *produtoMarcaService) Create(req *dto.ProdutoMarcaRequest) (*models.ProdutoMarca, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.repo.ExistsByDescricao(descricao, 0)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar descrição: %w", err)
	}
	if exists {
		return nil, errors.New("já existe uma marca de produto com esta descrição")
	}

	marca, err := req.ToModel()
	if err != nil {
		return nil, fmt.Errorf("erro ao converter dados: %w", err)
	}

	if err := s.repo.Create(marca); err != nil {
		return nil, fmt.Errorf("erro ao criar marca de produto: %w", err)
	}

	return marca, nil
}

func (s *produtoMarcaService) GetByID(id int) (*models.ProdutoMarca, error) {
	return s.repo.FindByID(id)
}

func (s *produtoMarcaService) Update(id int, req *dto.ProdutoMarcaRequest) (*models.ProdutoMarca, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	marca, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.repo.ExistsByDescricao(descricao, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar descrição: %w", err)
	}
	if exists {
		return nil, errors.New("já existe uma marca de produto com esta descrição")
	}

	marca.Descricao = descricao
	marca.Situacao = req.Situacao
	marca.UpdatedBy = req.UpdatedBy

	if err := s.repo.Update(id, marca); err != nil {
		return nil, fmt.Errorf("erro ao atualizar marca de produto: %w", err)
	}

	return marca, nil
}

func (s *produtoMarcaService) Delete(id int) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return err
	}

	count, err := s.repo.CountByMarca(id)
	if err != nil {
		return fmt.Errorf("erro ao verificar uso da marca: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("marca está em uso por %d produto(s) e não pode ser excluída", count)
	}

	return s.repo.Delete(id)
}

func (s *produtoMarcaService) List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoMarca, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.List(limit, offset, filters)
}