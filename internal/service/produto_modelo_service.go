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

type ProdutoModeloService struct {
	repo *repository.ProdutoModeloRepository
}

func NewProdutoModeloService(db *gorm.DB) *ProdutoModeloService {
	return &ProdutoModeloService{
		repo: repository.NewProdutoModeloRepository(db),
	}
}

func (s *ProdutoModeloService) Create(req *dto.ProdutoModeloRequest) (*models.ProdutoModelo, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.repo.ExistsByDescricao(descricao, 0)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar descrição: %w", err)
	}
	if exists {
		return nil, errors.New("já existe um modelo de produto com esta descrição")
	}

	modelo, err := req.ToModel()
	if err != nil {
		return nil, fmt.Errorf("erro ao converter dados: %w", err)
	}

	if err := s.repo.Create(modelo); err != nil {
		return nil, fmt.Errorf("erro ao criar modelo de produto: %w", err)
	}

	return modelo, nil
}

func (s *ProdutoModeloService) GetByID(id int) (*models.ProdutoModelo, error) {
	return s.repo.FindByID(id)
}

func (s *ProdutoModeloService) Update(id int, req *dto.ProdutoModeloRequest) (*models.ProdutoModelo, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	modelo, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.repo.ExistsByDescricao(descricao, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar descrição: %w", err)
	}
	if exists {
		return nil, errors.New("já existe um modelo de produto com esta descrição")
	}

	modelo.Descricao = descricao
	modelo.Situacao = req.Situacao
	modelo.UpdatedBy = req.UpdatedBy

	if err := s.repo.Update(id, modelo); err != nil {
		return nil, fmt.Errorf("erro ao atualizar modelo de produto: %w", err)
	}

	return modelo, nil
}

func (s *ProdutoModeloService) Delete(id int) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return err
	}

	count, err := s.repo.CountByModelo(id)
	if err != nil {
		return fmt.Errorf("erro ao verificar uso do modelo: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("modelo está em uso por %d produto(s) e não pode ser excluído", count)
	}

	return s.repo.Delete(id)
}

func (s *ProdutoModeloService) List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoModelo, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.List(limit, offset, filters)
}