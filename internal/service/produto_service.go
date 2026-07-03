package service

import (
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
)

type ProdutoService struct {
	repo *repository.ProdutoRepository
}

func NewProdutoService(repo *repository.ProdutoRepository) *ProdutoService {
	return &ProdutoService{repo: repo}
}

func (s *ProdutoService) Create(produto *models.Produto) (models.Produto, error) {
	err := s.repo.Create(produto)
	if err != nil {
		return models.Produto{}, err
	}
	return *produto, nil
}

func (s *ProdutoService) Update(produto *models.Produto) (models.Produto, error) {
	err := s.repo.Update(produto)
	if err != nil {
		return models.Produto{}, err
	}
	return *produto, nil
}

func (s *ProdutoService) Delete(id int) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProdutoService) FindByID(id int) (*models.Produto, error) {
	produto, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return produto, nil
}

func (s *ProdutoService) FindByNome(nome string, limit int) ([]models.Produto, error) {
	produtos, err := s.repo.FindByNome(nome, limit)
	if err != nil {
		return nil, err
	}
	return produtos, nil
}

func (s *ProdutoService) FindByCodigo(codigo string, limit int) (models.Produto, error) {
	produto, err := s.repo.FindByCodigo(codigo)
	if err != nil {
		return models.Produto{}, err
	}
	return *produto, nil
}
