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

type TabelaPrecoService struct {
	repo *repository.TabelaPrecoRepository
}

func NewTabelaPrecoService(db *gorm.DB) *TabelaPrecoService {
	return &TabelaPrecoService{
		repo: repository.NewTabelaPrecoRepository(db),
	}
}

func (s *TabelaPrecoService) Create(req *dto.TabelaPrecoRequest) (*models.TabelaPreco, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.repo.ExistsByDescricao(descricao, 0)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar descrição: %w", err)
	}
	if exists {
		return nil, errors.New("já existe uma tabela de preço com esta descrição")
	}

	tabela, err := req.ToModel()
	if err != nil {
		return nil, fmt.Errorf("erro ao converter dados: %w", err)
	}

	if err := s.repo.Create(tabela); err != nil {
		return nil, fmt.Errorf("erro ao criar tabela de preço: %w", err)
	}

	return tabela, nil
}

func (s *TabelaPrecoService) GetByID(id int) (*models.TabelaPreco, error) {
	return s.repo.FindByID(id)
}

func (s *TabelaPrecoService) Update(id int, req *dto.TabelaPrecoRequest) (*models.TabelaPreco, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	tabela, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.repo.ExistsByDescricao(descricao, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar descrição: %w", err)
	}
	if exists {
		return nil, errors.New("já existe uma tabela de preço com esta descrição")
	}

	tabela.Descricao = descricao
	tabela.UpdatedBy = req.UpdatedBy

	if err := s.repo.Update(id, tabela); err != nil {
		return nil, fmt.Errorf("erro ao atualizar tabela de preço: %w", err)
	}

	return tabela, nil
}

func (s *TabelaPrecoService) Delete(id int) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return err
	}

	count, err := s.repo.CountByTabelaPreco(id)
	if err != nil {
		return fmt.Errorf("erro ao verificar uso da tabela de preço: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("tabela de preço está em uso por %d documento(s) e não pode ser excluída", count)
	}

	return s.repo.Delete(id)
}

func (s *TabelaPrecoService) List(limit, offset int, filters map[string]interface{}) ([]models.TabelaPreco, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.List(limit, offset, filters)
}
