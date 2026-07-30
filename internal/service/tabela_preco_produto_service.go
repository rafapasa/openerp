package service

import (
	"errors"
	"fmt"

	"github.com/openerp/backend/internal/dto"
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

func NewTabelaPrecoProdutoService(tbppRepo repository.TabelaPrecoProdutoRepository, tbpService TabelaPrecoService, proService ProdutoService) TabelaPrecoProdutoService {
	return &tabelaPrecoProdutoService{
		tbppRepo:   tbppRepo,
		tbpService: tbpService,
		proService: proService,
	}
}

func (s *tabelaPrecoProdutoService) validateDependencies(tabelaPrecoID, produtoID int) error {
	if _, err := s.tbpService.GetByID(tabelaPrecoID); err != nil {
		return errors.New("tabela de preço não encontrada")
	}
	produto, err := s.proService.GetByID(produtoID)
	if err != nil {
		return fmt.Errorf("erro ao verificar produto: %w", err)
	}
	if produto == nil{
		return errors.New("produto não encontrado")
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
		return nil, fmt.Errorf("erro ao verificar duplicidade de produto: %w", err)
	}
	if exists {
		return nil, errors.New("este produto já foi adicionado a esta tabela de preço")
	}

	createItem, err := req.ToModel()
	if err != nil {
		return nil, fmt.Errorf("erro ao converter dados: %w", err)
	}

	if err := s.tbppRepo.Create(createItem); err != nil {
		return nil, fmt.Errorf("erro ao adicionar produto à tabela de preço: %w", err)
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
		return nil, fmt.Errorf("erro ao verificar duplicidade de produto: %w", err)
	}
	if exists {
		return nil, errors.New("este produto já foi adicionado a esta tabela de preço")
	}

	err = utils.MapToModel(req, updatedItem)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter dados: %w", err)
	}

	updatedItem.TabelaPrecoID = id // Garantir que o ID não seja alterado
	updatedItem.Item = item        // Garantir que o ID não seja alterado

	if err := s.tbppRepo.Update(id, item, updatedItem); err != nil {
		return nil, fmt.Errorf("erro ao atualizar produto na tabela de preço: %w", err)
	}

	return s.tbppRepo.FindByID(id, item)
}

func (s *tabelaPrecoProdutoService) Delete(id, item int) error {
	if id == 0 {
		return errors.New("ID da tabela de preço inválido")
	}
	if item == 0 {
		return errors.New("item inválido")
	}
	return s.tbppRepo.Delete(id, item)
}

func (s *tabelaPrecoProdutoService) List(tabelaPrecoID, limit, offset int, filters map[string]interface{}) ([]models.TabelaPrecoProduto, int64, error) {
	if _, err := s.tbpService.GetByID(tabelaPrecoID); err != nil {
		return nil, 0, errors.New("tabela de preço não encontrada")
	}
	if limit <= 0 {
		limit = 10
	}
	return s.tbppRepo.List(tabelaPrecoID, limit, offset, filters)
}

func (s *tabelaPrecoProdutoService) GetByProduto(tbpId, proId int) (*models.TabelaPrecoProduto, error) {
	return s.tbppRepo.FindByProduto(tbpId, proId)
}
