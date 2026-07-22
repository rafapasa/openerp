package service

import (
	"errors"
	"fmt"

	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/utils"
	"gorm.io/gorm"
)

type TabelaPrecoProdutoService struct {
	repo        *repository.TabelaPrecoProdutoRepository
	tabelaRepo  *repository.TabelaPrecoRepository
	produtoRepo *repository.ProdutoRepository
}

func NewTabelaPrecoProdutoService(db *gorm.DB) *TabelaPrecoProdutoService {
	return &TabelaPrecoProdutoService{
		repo:        repository.NewTabelaPrecoProdutoRepository(db),
		tabelaRepo:  repository.NewTabelaPrecoRepository(db),
		produtoRepo: repository.NewProdutoRepository(db),
	}
}

func (s *TabelaPrecoProdutoService) validateDependencies(tabelaPrecoID, produtoID int) error {
	if _, err := s.tabelaRepo.FindByID(tabelaPrecoID); err != nil {
		return errors.New("tabela de preço não encontrada")
	}
	exists, err := s.produtoRepo.ExistsByID(produtoID)
	if err != nil {
		return fmt.Errorf("erro ao verificar produto: %w", err)
	}
	if !exists {
		return errors.New("produto não encontrado")
	}
	return nil
}

func (s *TabelaPrecoProdutoService) Create(req *dto.TabelaPrecoProdutoRequest) (*models.TabelaPrecoProduto, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if err := s.validateDependencies(req.TabelaPrecoID, req.ProdutoID); err != nil {
		return nil, err
	}

	exists, err := s.repo.ExistsByProduto(req.TabelaPrecoID, req.ProdutoID, 0)
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

	if err := s.repo.Create(createItem); err != nil {
		return nil, fmt.Errorf("erro ao adicionar produto à tabela de preço: %w", err)
	}

	return s.repo.FindByID(createItem.TabelaPrecoID, createItem.Item)
}

func (s *TabelaPrecoProdutoService) GetByID(id, item int) (*models.TabelaPrecoProduto, error) {
	return s.repo.FindByID(id, item)
}

func (s *TabelaPrecoProdutoService) Update(id, item int, req *dto.TabelaPrecoProdutoRequest) (*models.TabelaPrecoProduto, error) {
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

	exists, err := s.repo.ExistsByProduto(req.TabelaPrecoID, req.ProdutoID, id)
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

	if err := s.repo.Update(id, item, updatedItem); err != nil {
		return nil, fmt.Errorf("erro ao atualizar produto na tabela de preço: %w", err)
	}

	return s.repo.FindByID(id, item)
}

func (s *TabelaPrecoProdutoService) Delete(id, item int) error {
	if id == 0 {
		return errors.New("ID da tabela de preço inválido")
	}
	if item == 0 {
		return errors.New("item inválido")
	}
	return s.repo.Delete(id, item)
}

func (s *TabelaPrecoProdutoService) List(tabelaPrecoID, limit, offset int, filters map[string]interface{}) ([]models.TabelaPrecoProduto, int64, error) {
	if _, err := s.tabelaRepo.FindByID(tabelaPrecoID); err != nil {
		return nil, 0, errors.New("tabela de preço não encontrada")
	}
	if limit <= 0 {
		limit = 10
	}
	return s.repo.List(tabelaPrecoID, limit, offset, filters)
}

func (s *TabelaPrecoProdutoService) GetByProduto(tbpId, proId int) (*models.TabelaPrecoProduto, error) {
	return s.repo.FindByProduto(tbpId, proId)
}