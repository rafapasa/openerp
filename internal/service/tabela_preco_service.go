package service

import (
	"fmt"
	"strings"

	apperrors "github.com/openerp/backend/internal/erros"

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
		return nil, err
	}
	if exists {
		return nil, apperrors.NewConflictError("Já existe uma tabela de preço com esta descrição.") //
	}

	tabela, err := req.ToModel()
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao converter dados.", err) //
	}

	if err := s.tbpRepo.Create(tabela); err != nil {
		return nil, err
	}

	return tabela, nil
}

func (s *tabelaPrecoService) GetByID(id int) (*models.TabelaPreco, error) {
	return s.tbpRepo.FindByID(id) //
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
		return nil, err
	}
	if exists {
		return nil, apperrors.NewConflictError("Já existe uma tabela de preço com esta descrição.") //
	}

	tabela.Descricao = descricao
	tabela.UpdatedBy = req.UpdatedBy

	if err := s.tbpRepo.Update(id, tabela); err != nil {
		return nil, err
	}

	return tabela, nil
}

func (s *tabelaPrecoService) Delete(id int) error {
	if _, err := s.tbpRepo.FindByID(id); err != nil { //
		return err
	}

	count, err := s.tbpRepo.CountProdutosByTabela(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return apperrors.NewConflictError(fmt.Sprintf("Tabela de preço está em uso por %d documento(s) e não pode ser excluída.", count)) //
	}

	return s.tbpRepo.Delete(id)
}

func (s *tabelaPrecoService) List(limit, offset int, filters map[string]interface{}) ([]models.TabelaPreco, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.tbpRepo.List(limit, offset, filters)
}
