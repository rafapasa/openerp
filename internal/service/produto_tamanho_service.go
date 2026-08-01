package service

import (
	"fmt"

	"github.com/openerp/backend/internal/dto"
	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/utils"
)

// ProdutoTamanhoService define os métodos públicos para o serviço de tamanhos de produto.
type ProdutoTamanhoService interface {
	Create(req *dto.ProdutoTamanhoRequest) (*dto.ProdutoTamanhoResponse, error)
	GetByID(id int) (*dto.ProdutoTamanhoResponse, error)
	Update(id int, req *dto.ProdutoTamanhoRequest) (*dto.ProdutoTamanhoResponse, error)
	Delete(id int) error
	List(limit, offset int, filters map[string]interface{}) ([]dto.ProdutoTamanhoResponse, int64, error)
	FindByID(id int) (*models.ProdutoTamanho, error) // Adicionado para uso interno por outros serviços
}

type produtoTamanhoService struct {
	repo repository.ProdutoTamanhoRepository
}

// NewProdutoTamanhoService cria uma nova instância de ProdutoTamanhoService.
func NewProdutoTamanhoService(repo repository.ProdutoTamanhoRepository) ProdutoTamanhoService {
	return &produtoTamanhoService{
		repo: repo,
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO (PRIVADOS)
// ============================================================

func (s *produtoTamanhoService) validateProdutoTamanho(id int, req *dto.ProdutoTamanhoRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	// Validar unicidade da sigla
	exists, err := s.repo.ExistsBySigla(req.Sigla, req.EmpresaFilialID, id)
	if err != nil {
		return err
	}
	if exists {
		return apperrors.NewConflictError(fmt.Sprintf("Já existe um tamanho com a sigla '%s' para esta filial.", req.Sigla))
	}

	// Validar unicidade do nome
	exists, err = s.repo.ExistsByNome(req.Descricao, req.EmpresaFilialID, req.ID)
	if err != nil {
		return err
	}
	if exists {
		return apperrors.NewConflictError(fmt.Sprintf("Já existe um tamanho com o nome '%s' para esta filial.", req.Descricao))
	}

	return nil
}

// mapModelToResponse mapeia um modelo ProdutoTamanho para um DTO de resposta.
func (s *produtoTamanhoService) mapModelToResponse(tamanho *models.ProdutoTamanho) (*dto.ProdutoTamanhoResponse, error) {
	resp := &dto.ProdutoTamanhoResponse{}
	if err := utils.MapToModel(tamanho, resp); err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear modelo para DTO de resposta.", err)
	}
	// TODO: Popular EmpresaFilialNome se necessário, injetando EmpresaFilialService
	return resp, nil
}

// ============================================================
// MÉTODOS PÚBLICOS (CRUD)
// ============================================================

// Create cria um novo tamanho de produto.
func (s *produtoTamanhoService) Create(req *dto.ProdutoTamanhoRequest) (*dto.ProdutoTamanhoResponse, error) {
	if err := s.validateProdutoTamanho(0, req); err != nil {
		return nil, err
	}

	tamanho, err := req.ToModel()
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear DTO para modelo.", err)
	}

	if err := s.repo.Create(tamanho); err != nil {
		return nil, err
	}

	return s.mapModelToResponse(tamanho)
}

// GetByID busca um tamanho de produto pelo ID.
func (s *produtoTamanhoService) GetByID(id int) (*dto.ProdutoTamanhoResponse, error) {
	tamanho, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if tamanho.IsDeleted() {
		return nil, apperrors.NewNotFoundError(fmt.Sprintf("Tamanho de produto com ID %d não encontrado.", id))
	}
	return s.mapModelToResponse(tamanho)
}

// FindByID busca um tamanho de produto pelo ID (retorna o modelo).
func (s *produtoTamanhoService) FindByID(id int) (*models.ProdutoTamanho, error) {
	return s.repo.FindByID(id)
}

// Update atualiza um tamanho de produto existente.
func (s *produtoTamanhoService) Update(id int, req *dto.ProdutoTamanhoRequest) (*dto.ProdutoTamanhoResponse, error) {
	existingTamanho, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existingTamanho.IsDeleted() {
		return nil, apperrors.NewNotFoundError(fmt.Sprintf("Tamanho de produto com ID %d não encontrado para atualização.", id))
	}

	req.ID = id // Garante que o ID da requisição corresponde ao ID da URL
	if err := s.validateProdutoTamanho(0, req); err != nil {
		return nil, err
	}

	if err := utils.MapToModel(req, existingTamanho); err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear DTO para modelo existente.", err)
	}

	if err := s.repo.Update(existingTamanho); err != nil {
		return nil, err
	}

	return s.mapModelToResponse(existingTamanho)
}

// Delete realiza a exclusão lógica de um tamanho de produto.
func (s *produtoTamanhoService) Delete(id int) error {
	// TODO: Adicionar verificação de dependências (ex: se o tamanho está em uso por alguma variação de produto)
	return s.repo.Delete(id)
}

// List lista tamanhos de produto com paginação e filtros.
func (s *produtoTamanhoService) List(limit, offset int, filters map[string]interface{}) ([]dto.ProdutoTamanhoResponse, int64, error) {
	tamanhos, total, err := s.repo.List(limit, offset, filters)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.ProdutoTamanhoResponse
	for _, tamanho := range tamanhos {
		resp, err := s.mapModelToResponse(&tamanho)
		if err != nil {
			// Logar o erro, mas continuar processando os outros tamanhos
			fmt.Printf("Erro ao mapear tamanho %d para resposta: %v\n", tamanho.ID, err)
			continue
		}
		responses = append(responses, *resp)
	}
	return responses, total, nil
}
