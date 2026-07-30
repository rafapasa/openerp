package service

import (
	"github.com/openerp/backend/internal/dto"
	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/utils"
)

// TabelaPrecoProdutoService define os métodos públicos para o serviço de produtos em tabelas de preço.
type TabelaPrecoProdutoService interface {
	Create(req *dto.TabelaPrecoProdutoRequest) (*models.TabelaPrecoProduto, error)
	GetByID(id, item int) (*models.TabelaPrecoProduto, error)
	Update(id, item int, req *dto.TabelaPrecoProdutoRequest) (*models.TabelaPrecoProduto, error)
	Delete(id, item int) error
	List(tabelaPrecoID, limit, offset int, filters map[string]interface{}) ([]models.TabelaPrecoProduto, int64, error)
	GetByProduto(tbpId, proId int) (*models.TabelaPrecoProduto, error)
}

type tabelaPrecoProdutoService struct {
	tbppRepo   repository.TabelaPrecoProdutoRepository
	tbpService TabelaPrecoService
	proService ProdutoService
}

func NewTabelaPrecoProdutoService(
	tbppRepo repository.TabelaPrecoProdutoRepository,
	tbpService TabelaPrecoService,
	proService ProdutoService,
) TabelaPrecoProdutoService {
	return &tabelaPrecoProdutoService{
		tbppRepo:   tbppRepo,
		tbpService: tbpService,
		proService: proService,
	}
}

func (s *tabelaPrecoProdutoService) validateDependencies(tabelaPrecoID, produtoID int) error {
	if _, err := s.tbpService.GetByID(tabelaPrecoID); err != nil {
		return err
	}
	_, err := s.proService.GetByID(produtoID)
	if err != nil {
		return err
	}

	return nil
}

func (s *tabelaPrecoProdutoService) Create(req *dto.TabelaPrecoProdutoRequest) (*models.TabelaPrecoProduto, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if err := s.validateDependencies(req.TabelaPrecoID, req.ProdutoID); err != nil {
		return nil, err
	}

	exists, err := s.tbppRepo.ExistsByTabelaPrecoAndProduto(req.TabelaPrecoID, req.ProdutoID, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperrors.NewConflictError("Este produto já foi adicionado a esta tabela de preço.") //
	}

	createItem, err := req.ToModel()
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao converter dados.", err) //
	}

	if err := s.tbppRepo.Create(createItem); err != nil {
		return nil, err
	}

	return s.tbppRepo.FindByID(createItem.TabelaPrecoID, createItem.Item)
}

func (s *tabelaPrecoProdutoService) GetByID(id, item int) (*models.TabelaPrecoProduto, error) {
	return s.tbppRepo.FindByID(id, item)
}

func (s *tabelaPrecoProdutoService) Update(id, item int, req *dto.TabelaPrecoProdutoRequest) (*models.TabelaPrecoProduto, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	updatedItem, err := s.GetByID(id, item)
	if err != nil {
		return nil, err
	}

	if err := s.validateDependencies(req.TabelaPrecoID, req.ProdutoID); err != nil {
		return nil, err
	}

	exists, err := s.tbppRepo.ExistsByTabelaPrecoAndProduto(req.TabelaPrecoID, req.ProdutoID, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperrors.NewConflictError("Este produto já foi adicionado a esta tabela de preço.") //
	}

	err = utils.MapToModel(req, updatedItem)
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao converter dados.", err) //
	}

	updatedItem.TabelaPrecoID = id // Garantir que o ID não seja alterado
	updatedItem.Item = item        // Garantir que o ID não seja alterado

	if err := s.tbppRepo.Update(id, item, updatedItem); err != nil {
		return nil, err
	}

	return s.tbppRepo.FindByID(id, item)
}

func (s *tabelaPrecoProdutoService) Delete(id, item int) error {
	if id == 0 { //
		return apperrors.NewValidationError("ID da tabela de preço inválido.") //
	} //
	if item == 0 { //
		return apperrors.NewValidationError("Item inválido.") //
	} //
	return s.tbppRepo.Delete(id, item)
}

func (s *tabelaPrecoProdutoService) List(tabelaPrecoID, limit, offset int, filters map[string]interface{}) ([]models.TabelaPrecoProduto, int64, error) {
	if _, err := s.tbpService.GetByID(tabelaPrecoID); err != nil {
		return nil, 0, err
	}
	if limit <= 0 {
		limit = 10
	}
	return s.tbppRepo.List(tabelaPrecoID, limit, offset, filters)
}

func (s *tabelaPrecoProdutoService) GetByProduto(tbpId, proId int) (*models.TabelaPrecoProduto, error) {
	return s.tbppRepo.FindByTabelaPrecoAndProduto(tbpId, proId)
}
