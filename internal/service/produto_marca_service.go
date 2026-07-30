package service

import (
	"fmt"
	"strings"

	"github.com/openerp/backend/internal/dto"
	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
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
	promRepo repository.ProdutoMarcaRepository
}

func NewProdutoMarcaService(promRepo repository.ProdutoMarcaRepository) ProdutoMarcaService {
	return &produtoMarcaService{
		promRepo: promRepo,
	}
}

func (s *produtoMarcaService) Create(req *dto.ProdutoMarcaRequest) (*models.ProdutoMarca, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.promRepo.ExistsByDescricao(descricao, 0)
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao verificar descrição.", err) //
	}
	if exists {
		return nil, apperrors.NewConflictError("Já existe uma marca de produto com esta descrição.") //
	}

	marca, err := req.ToModel()
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao converter dados.", err) //
	}

	if err := s.promRepo.Create(marca); err != nil {
		return nil, apperrors.NewInternalError("Erro ao criar marca de produto.", err) //
	}

	return marca, nil
}

func (s *produtoMarcaService) GetByID(id int) (*models.ProdutoMarca, error) {
	return s.promRepo.FindByID(id) //
}

func (s *produtoMarcaService) Update(id int, req *dto.ProdutoMarcaRequest) (*models.ProdutoMarca, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	marca, err := s.promRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.promRepo.ExistsByDescricao(descricao, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperrors.NewConflictError("Já existe uma marca de produto com esta descrição.") //
	}

	marca.Descricao = descricao
	marca.Situacao = req.Situacao
	marca.UpdatedBy = req.UpdatedBy

	if err := s.promRepo.Update(id, marca); err != nil {
		return nil, err
	}

	return marca, nil
}

func (s *produtoMarcaService) Delete(id int) error {
	if _, err := s.promRepo.FindByID(id); err != nil {
		return err
	}

	count, err := s.promRepo.CountProdutosByMarca(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return apperrors.NewConflictError(fmt.Sprintf("Marca está em uso por %d produto(s) e não pode ser excluída.", count)) //
	}

	return s.promRepo.Delete(id)
}

func (s *produtoMarcaService) List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoMarca, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.promRepo.List(limit, offset, filters)
}
