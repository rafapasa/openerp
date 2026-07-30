package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
)

// TabelaPrecoService define os métodos públicos para o serviço de tabela de preço.
type TabelaPrecoService interface {
	Create(req *dto.TabelaPrecoRequest) (*models.TabelaPreco, error)
	GetByID(id int) (*models.TabelaPreco, error)
	Update(id int, req *dto.TabelaPrecoRequest) (*models.TabelaPreco, error)
	Delete(id int) error
	List(limit, offset int, filters map[string]interface{}) ([]models.TabelaPreco, int64, error)
}

type tabelaPrecoService struct {
	tbpRepo repository.TabelaPrecoRepository
}

func NewTabelaPrecoService(tbptbpRepo repository.TabelaPrecoRepository) TabelaPrecoService {
	return &tabelaPrecoService{
		tbpRepo: tbptbpRepo,
	}
}

func (s *tabelaPrecoService) Create(req *dto.TabelaPrecoRequest) (*models.TabelaPreco, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.tbpRepo.ExistsByDescricao(descricao, 0)
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

	if err := s.tbpRepo.Create(tabela); err != nil {
		return nil, fmt.Errorf("erro ao criar tabela de preço: %w", err)
	}

	return tabela, nil
}

func (s *tabelaPrecoService) GetByID(id int) (*models.TabelaPreco, error) {
	return s.tbpRepo.FindByID(id)
}

func (s *tabelaPrecoService) Update(id int, req *dto.TabelaPrecoRequest) (*models.TabelaPreco, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	tabela, err := s.tbpRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.tbpRepo.ExistsByDescricao(descricao, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar descrição: %w", err)
	}
	if exists {
		return nil, errors.New("já existe uma tabela de preço com esta descrição")
	}

	tabela.Descricao = descricao
	tabela.UpdatedBy = req.UpdatedBy

	if err := s.tbpRepo.Update(id, tabela); err != nil {
		return nil, fmt.Errorf("erro ao atualizar tabela de preço: %w", err)
	}

	return tabela, nil
}

func (s *tabelaPrecoService) Delete(id int) error {
	if _, err := s.tbpRepo.FindByID(id); err != nil {
		return err
	}

	count, err := s.tbpRepo.CountProdutosByTabela(id)
	if err != nil {
		return fmt.Errorf("erro ao verificar uso da tabela de preço: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("tabela de preço está em uso por %d documento(s) e não pode ser excluída", count)
	}

	return s.tbpRepo.Delete(id)
}

func (s *tabelaPrecoService) List(limit, offset int, filters map[string]interface{}) ([]models.TabelaPreco, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.tbpRepo.List(limit, offset, filters)
}
