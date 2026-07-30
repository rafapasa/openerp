package service

import (
	"fmt"
	"strings"

	"github.com/openerp/backend/internal/dto"
	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
)

// ProdutoModeloService define os métodos públicos para o serviço de modelo de produto.
type ProdutoModeloService interface {
	Create(req *dto.ProdutoModeloRequest) (*models.ProdutoModelo, error)
	GetByID(id int) (*models.ProdutoModelo, error)
	Update(id int, req *dto.ProdutoModeloRequest) (*models.ProdutoModelo, error)
	Delete(id int) error
	List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoModelo, int64, error)
}

type produtoModeloService struct {
	promRepo repository.ProdutoModeloRepository
}

func NewProdutoModeloService(promRepo repository.ProdutoModeloRepository) ProdutoModeloService {
	return &produtoModeloService{
		promRepo: promRepo,
	}
}

func (s *produtoModeloService) Create(req *dto.ProdutoModeloRequest) (*models.ProdutoModelo, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.promRepo.ExistsByDescricao(descricao, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperrors.NewConflictError("Já existe um modelo de produto com esta descrição.") //
	}

	modelo, err := req.ToModel()
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao converter dados.", err) //
	}

	if err := s.promRepo.Create(modelo); err != nil {
		return nil, err
	}

	return modelo, nil
}

func (s *produtoModeloService) GetByID(id int) (*models.ProdutoModelo, error) {
	return s.promRepo.FindByID(id) //
}

func (s *produtoModeloService) Update(id int, req *dto.ProdutoModeloRequest) (*models.ProdutoModelo, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	modelo, err := s.promRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.promRepo.ExistsByDescricao(descricao, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperrors.NewConflictError("Já existe um modelo de produto com esta descrição.") //
	}

	modelo.Descricao = descricao
	modelo.Situacao = req.Situacao
	modelo.UpdatedBy = req.UpdatedBy

	if err := s.promRepo.Update(id, modelo); err != nil {
		return nil, err
	}

	return modelo, nil
}

func (s *produtoModeloService) Delete(id int) error {
	if _, err := s.promRepo.FindByID(id); err != nil { //
		return err
	}

	count, err := s.promRepo.CountProdutosByModelo(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return apperrors.NewConflictError(fmt.Sprintf("Modelo está em uso por %d produto(s) e não pode ser excluído.", count)) //
	}

	return s.promRepo.Delete(id)
}

func (s *produtoModeloService) List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoModelo, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.promRepo.List(limit, offset, filters)
}
