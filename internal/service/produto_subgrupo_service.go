package service

import (
	"context"
	"strings"

	"github.com/openerp/backend/internal/apperrors"

	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
)

// ProdutoSubgrupoService define os métodos públicos para o serviço de subgrupo de produto.
type ProdutoSubgrupoService interface {
	Create(ctx context.Context, req *dto.ProdutoSubgrupoRequest) (*models.ProdutoSubgrupo, error)
	GetByID(ctx context.Context, id int) (*models.ProdutoSubgrupo, error)
	Update(ctx context.Context, id int, req *dto.ProdutoSubgrupoRequest) (*models.ProdutoSubgrupo, error)
	Delete(id int) error
	List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoSubgrupo, int64, error)
}

// ============================================================
// TYPES
// ============================================================

type produtoSubgrupoService struct {
	prosgRepo repository.ProdutoSubgrupoRepository
}

// ============================================================
// CONSTRUCTOR
// ============================================================

func NewProdutoSubgrupoService(prosgRepo repository.ProdutoSubgrupoRepository) ProdutoSubgrupoService {
	return &produtoSubgrupoService{
		prosgRepo: prosgRepo,
	}
}

// ============================================================
// MÉTODOS PRINCIPAIS (CRUD)
// ============================================================

// Create cria um novo subgrupo de produto.
func (s *produtoSubgrupoService) Create(ctx context.Context, req *dto.ProdutoSubgrupoRequest) (*models.ProdutoSubgrupo, error) {
	if err := req.Validate(); err != nil { // Context not used in Validate
		return nil, err
	}

	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.prosgRepo.ExistsByDescricao(descricao, 0)
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao verificar descrição.", err) //
	}
	if exists {
		return nil, apperrors.NewConflictError("Já existe um subgrupo de produto com esta descrição.") //
	}

	subgrupo, err := req.ToModel()
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao converter dados.", err) //
	}

	if err := s.prosgRepo.Create(subgrupo); err != nil {
		return nil, apperrors.NewInternalError("Erro ao criar subgrupo de produto.", err) //
	}

	return subgrupo, nil
}

// GetByID busca um subgrupo de produto por ID.
func (s *produtoSubgrupoService) GetByID(ctx context.Context, id int) (*models.ProdutoSubgrupo, error) {
	return s.prosgRepo.FindByID(id) // Context not used in FindByID
}

// Update atualiza um subgrupo de produto.
func (s *produtoSubgrupoService) Update(ctx context.Context, id int, req *dto.ProdutoSubgrupoRequest) (*models.ProdutoSubgrupo, error) {
	subgrupo, err := s.prosgRepo.FindByID(id) // Context not used in FindByID
	if err != nil {
		return nil, err
	}

	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.prosgRepo.ExistsByDescricao(descricao, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperrors.NewConflictError("Já existe um subgrupo de produto com esta descrição.") //
	}

	subgrupo.Descricao = descricao
	subgrupo.Situacao = req.Situacao
	subgrupo.UpdatedBy = req.UpdatedBy

	if err := s.prosgRepo.Update(id, subgrupo); err != nil {
		return nil, err
	}

	return subgrupo, nil
}

// Delete exclui um subgrupo de produto.
func (s *produtoSubgrupoService) Delete(id int) error {
	if _, err := s.prosgRepo.FindByID(id); err != nil {
		return err
	}

	// TODO: Adicionar verificação se o subgrupo está em uso por algum produto.
	// count, err := s.produtoRepo.CountBySubgrupo(id)
	// if err != nil { return err }
	// if count > 0 { return errors.New("subgrupo em uso e não pode ser excluído") }

	return s.prosgRepo.Delete(id)
}

// List lista todos os subgrupos de produto.
func (s *produtoSubgrupoService) List(limit, offset int, filters map[string]interface{}) ([]models.ProdutoSubgrupo, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.prosgRepo.List(limit, offset, filters)
}
