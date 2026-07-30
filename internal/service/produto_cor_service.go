package service

import (
	"fmt"

	"github.com/openerp/backend/internal/dto"
	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/utils"
)

// ProdutoCorService define os métodos públicos para o serviço de cores de produto.
type ProdutoCorService interface {
	Create(req *dto.ProdutoCorRequest) (*dto.ProdutoCorResponse, error)
	GetByID(id int) (*dto.ProdutoCorResponse, error)
	Update(id int, req *dto.ProdutoCorRequest) (*dto.ProdutoCorResponse, error)
	Delete(id int) error
	List(limit, offset int, filters map[string]interface{}) ([]dto.ProdutoCorResponse, int64, error)
	FindByID(id int) (*models.ProdutoCor, error) // Adicionado para uso interno por outros serviços
}

type produtoCorService struct {
	repo repository.ProdutoCorRepository
}

// NewProdutoCorService cria uma nova instância de ProdutoCorService.
func NewProdutoCorService(repo repository.ProdutoCorRepository) ProdutoCorService {
	return &produtoCorService{
		repo: repo,
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO (PRIVADOS)
// ============================================================

func (s *produtoCorService) validateProdutoCor(req *dto.ProdutoCorRequest, isUpdate bool) error {
	if err := req.Validate(); err != nil {
		return err
	}

	// Validar unicidade da sigla
	exists, err := s.repo.ExistsBySigla(req.Sigla, req.EmpresaFilialID, req.ID)
	if err != nil {
		return err
	}
	if exists {
		return apperrors.NewConflictError(fmt.Sprintf("Já existe uma cor com a sigla '%s' para esta filial.", req.Sigla))
	}

	// Validar unicidade do nome
	exists, err = s.repo.ExistsByNome(req.Nome, req.EmpresaFilialID, req.ID)
	if err != nil {
		return err
	}
	if exists {
		return apperrors.NewConflictError(fmt.Sprintf("Já existe uma cor com o nome '%s' para esta filial.", req.Nome))
	}

	return nil
}

// mapModelToResponse mapeia um modelo ProdutoCor para um DTO de resposta.
func (s *produtoCorService) mapModelToResponse(cor *models.ProdutoCor) (*dto.ProdutoCorResponse, error) {
	resp := &dto.ProdutoCorResponse{}
	if err := utils.MapToModel(cor, resp); err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear modelo para DTO de resposta.", err)
	}
	// TODO: Popular EmpresaFilialNome se necessário, injetando EmpresaFilialService
	return resp, nil
}

// ============================================================
// MÉTODOS PÚBLICOS (CRUD)
// ============================================================

// Create cria uma nova cor de produto.
func (s *produtoCorService) Create(req *dto.ProdutoCorRequest) (*dto.ProdutoCorResponse, error) {
	if err := s.validateProdutoCor(req, false); err != nil {
		return nil, err
	}

	cor, err := req.ToModel()
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear DTO para modelo.", err)
	}

	if err := s.repo.Create(cor); err != nil {
		return nil, err
	}

	return s.mapModelToResponse(cor)
}

// GetByID busca uma cor de produto pelo ID.
func (s *produtoCorService) GetByID(id int) (*dto.ProdutoCorResponse, error) {
	cor, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if cor.IsDeleted() {
		return nil, apperrors.NewNotFoundError(fmt.Sprintf("Cor de produto com ID %d não encontrada.", id))
	}
	return s.mapModelToResponse(cor)
}

// FindByID busca uma cor de produto pelo ID (retorna o modelo).
func (s *produtoCorService) FindByID(id int) (*models.ProdutoCor, error) {
	return s.repo.FindByID(id)
}

// Update atualiza uma cor de produto existente.
func (s *produtoCorService) Update(id int, req *dto.ProdutoCorRequest) (*dto.ProdutoCorResponse, error) {
	existingCor, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existingCor.IsDeleted() {
		return nil, apperrors.NewNotFoundError(fmt.Sprintf("Cor de produto com ID %d não encontrada para atualização.", id))
	}

	req.ID = id // Garante que o ID da requisição corresponde ao ID da URL
	if err := s.validateProdutoCor(req, true); err != nil {
		return nil, err
	}

	if err := utils.MapToModel(req, existingCor); err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear DTO para modelo existente.", err)
	}

	if err := s.repo.Update(existingCor); err != nil {
		return nil, err
	}

	return s.mapModelToResponse(existingCor)
}

// Delete realiza a exclusão lógica de uma cor de produto.
func (s *produtoCorService) Delete(id int) error {
	// TODO: Adicionar verificação de dependências (ex: se a cor está em uso por alguma variação de produto)
	return s.repo.Delete(id)
}

// List lista cores de produto com paginação e filtros.
func (s *produtoCorService) List(limit, offset int, filters map[string]interface{}) ([]dto.ProdutoCorResponse, int64, error) {
	cores, total, err := s.repo.List(limit, offset, filters)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.ProdutoCorResponse
	for _, cor := range cores {
		resp, err := s.mapModelToResponse(&cor)
		if err != nil {
			// Logar o erro, mas continuar processando as outras cores
			fmt.Printf("Erro ao mapear cor %d para resposta: %v\n", cor.ID, err)
			continue
		}
		responses = append(responses, *resp)
	}
	return responses, total, nil
}
