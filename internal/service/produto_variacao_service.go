package service

import (
	"context"
	"fmt"

	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/utils"
)

// ProdutoVariacaoService define os métodos públicos para o serviço de variações de produto.
type ProdutoVariacaoService interface {
	Create(ctx context.Context, req *dto.ProdutoVariacaoRequest) (*dto.ProdutoVariacaoResponse, error)
	GetByID(ctx context.Context, id int) (*dto.ProdutoVariacaoResponse, error)
	Update(ctx context.Context, id int, req *dto.ProdutoVariacaoRequest) (*dto.ProdutoVariacaoResponse, error)
	Delete(id int) error
	List(limit, offset int, filters map[string]interface{}) ([]dto.ProdutoVariacaoResponse, int64, error)
}

// ProdutoVariacaoService gerencia a lógica de negócios para variações de produto.
type produtoVariacaoService struct {
	provRepo       repository.ProdutoVariacaoRepository
	proService     ProdutoService
	corService     ProdutoCorService     // Renamed from procService
	tamanhoService ProdutoTamanhoService // Renamed from protService
}

// NewProdutoVariacaoService cria uma nova instância de ProdutoVariacaoService.
func NewProdutoVariacaoService(
	proService ProdutoService,
	procService ProdutoCorService,
	protService ProdutoTamanhoService,
	provRepo repository.ProdutoVariacaoRepository, // Changed parameter name to provRepo
) ProdutoVariacaoService {
	return &produtoVariacaoService{
		provRepo:       provRepo,
		proService:     proService,
		corService:     procService, // Assigned to new field name
		tamanhoService: protService, // Assigned to new field name
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO (PRIVADOS)
// ============================================================

func (s *produtoVariacaoService) validateProdutoVariacao(ctx context.Context, req *dto.ProdutoVariacaoRequest, isUpdate bool) error {
	if req.ProdutoID <= 0 {
		return apperrors.NewValidationError("O campo 'produto_id' é obrigatório e deve ser um valor positivo.")
	}
	if req.EmpresaFilialID <= 0 {
		return apperrors.NewValidationError("O campo 'empresa_filial_id' é obrigatório e deve ser um valor positivo.")
	}
	if req.SKU == "" {
		return apperrors.NewValidationError("O campo 'sku' é obrigatório.")
	}
	if len(req.SKU) > 50 {
		return apperrors.NewValidationError("O campo 'sku' não pode ter mais de 50 caracteres.")
	}
	if req.PrecoAdicional < 0 {
		return apperrors.NewValidationError("O campo 'preco_adicional' não pode ser negativo.")
	}
	if req.EstoqueAtual < 0 {
		return apperrors.NewValidationError("O campo 'estoque_atual' não pode ser negativo.")
	}

	// Verificar se o Produto existe
	produto, err := s.proService.FindById(req.ProdutoID)
	if err != nil {
		return err
	}
	if produto == nil {
		return apperrors.NewNotFoundError(fmt.Sprintf("Produto com ID %d não encontrado.", req.ProdutoID))
	}

	// Verificar se Cor existe (se informada)
	if req.CorID != nil && *req.CorID > 0 {
		_, err := s.corService.FindByID(ctx, *req.CorID)
		if err != nil {
			return err
		}

	}

	// Verificar se Tamanho existe (se informado)
	if req.TamanhoID != nil && *req.TamanhoID > 0 { // Changed from s.tamanhoRepo to s.tamanhoService
		_, err := s.tamanhoService.FindByID(ctx, *req.TamanhoID)
		if err != nil {
			return err
		}

	}

	// Verificar unicidade do SKU por EmpresaFilial
	existingVariacao, err := s.provRepo.FindBySKU(req.SKU)
	if err != nil {
		return apperrors.NewInternalError("Erro ao verificar SKU existente.", err)
	}
	if existingVariacao != nil {
		if isUpdate && existingVariacao.ID != req.ID { // Se for update, permite o mesmo SKU para o próprio registro
			return apperrors.NewConflictError(fmt.Sprintf("SKU '%s' já cadastrado para esta filial.", req.SKU))
		} else if !isUpdate { // Se for criação, não pode existir
			return apperrors.NewConflictError(fmt.Sprintf("SKU '%s' já cadastrado para esta filial.", req.SKU))
		}
	}

	return nil
}

// mapModelToResponse mapeia um modelo ProdutoVariacao para um DTO de resposta.
func (s *produtoVariacaoService) mapModelToResponse(variacao *models.ProdutoVariacao) (*dto.ProdutoVariacaoResponse, error) {
	resp := &dto.ProdutoVariacaoResponse{}
	if err := utils.MapToModel(variacao, resp); err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear modelo para DTO de resposta.", err)
	}

	if variacao.Produto != nil {
		resp.ProdutoNome = variacao.Produto.Nome
	}
	if variacao.EmpresaFilial != nil {
		resp.EmpresaFilialNome = variacao.EmpresaFilial.Nome
	}
	if variacao.Cor != nil {
		resp.CorNome = variacao.Cor.Nome
		resp.CorSigla = variacao.Cor.Sigla
	}
	if variacao.Tamanho != nil {
		resp.TamanhoNome = variacao.Tamanho.Nome
		resp.TamanhoSigla = variacao.Tamanho.Sigla
	}

	return resp, nil
}

// ============================================================
// MÉTODOS PÚBLICOS (CRUD)
// ============================================================

// Create cria uma nova variação de produto.
func (s *produtoVariacaoService) Create(ctx context.Context, req *dto.ProdutoVariacaoRequest) (*dto.ProdutoVariacaoResponse, error) {
	if err := s.validateProdutoVariacao(ctx, req, false); err != nil {
		return nil, err
	}

	variacao := &models.ProdutoVariacao{}
	if err := utils.MapToModel(req, variacao); err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear DTO para modelo.", err)
	}

	if err := s.provRepo.Create(variacao); err != nil {
		return nil, apperrors.NewInternalError("Erro ao criar variação de produto.", err)
	}

	createdVariacao, err := s.provRepo.FindByID(variacao.ID)
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao buscar variação de produto criada.", err)
	}

	return s.mapModelToResponse(createdVariacao)
}

// Create cria uma nova variação de produto.
// Este método já existe no arquivo, apenas garantindo que a interface o referencie.
// (No diff, ele já está presente, então não há mudança real aqui, apenas a menção para a interface)

// GetByID busca uma variação de produto pelo ID.
func (s *produtoVariacaoService) GetByID(ctx context.Context, id int) (*dto.ProdutoVariacaoResponse, error) {
	if id <= 0 {
		return nil, apperrors.NewValidationError("ID da variação de produto inválido.")
	}

	variacao, err := s.provRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if variacao == nil || variacao.IsDeleted() {
		return nil, apperrors.NewNotFoundError(fmt.Sprintf("Variação de produto com ID %d não encontrada.", id)) //
	}

	return s.mapModelToResponse(variacao)
}

// Update atualiza uma variação de produto existente.
func (s *produtoVariacaoService) Update(ctx context.Context, id int, req *dto.ProdutoVariacaoRequest) (*dto.ProdutoVariacaoResponse, error) {
	if id <= 0 {
		return nil, apperrors.NewValidationError("ID da variação de produto inválido.")
	}
	req.ID = id

	if err := s.validateProdutoVariacao(ctx, req, true); err != nil {
		return nil, err
	}

	existingVariacao, err := s.provRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existingVariacao == nil || existingVariacao.IsDeleted() {
		return nil, apperrors.NewNotFoundError(fmt.Sprintf("Variação de produto com ID %d não encontrada para atualização.", id)) //
	}

	if err := utils.MapToModel(req, existingVariacao); err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear DTO para modelo existente.", err)
	}

	if err := s.provRepo.Update(id, existingVariacao); err != nil {
		return nil, apperrors.NewInternalError("Erro ao atualizar variação de produto.", err)
	}

	updatedVariacao, err := s.provRepo.FindByID(existingVariacao.ID)
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao buscar variação de produto atualizada.", err)
	}

	return s.mapModelToResponse(updatedVariacao)
}

// Delete realiza a exclusão lógica de uma variação de produto.
func (s *produtoVariacaoService) Delete(id int) error {
	if id <= 0 {
		return apperrors.NewValidationError("ID da variação de produto inválido.")
	}

	existingVariacao, err := s.provRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existingVariacao == nil || existingVariacao.IsDeleted() {
		return apperrors.NewNotFoundError(fmt.Sprintf("Variação de produto com ID %d não encontrada para exclusão.", id)) //
	}

	if err := s.provRepo.Delete(id); err != nil {
		return apperrors.NewInternalError("Erro ao excluir variação de produto.", err)
	}
	return nil
}

// List lista variações de produto com paginação e filtros.
func (s *produtoVariacaoService) List(limit, offset int, filters map[string]interface{}) ([]dto.ProdutoVariacaoResponse, int64, error) {
	if limit < 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	variacoes, total, err := s.provRepo.List(limit, offset, filters)
	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro ao listar variações de produto.", err)
	}

	var responses []dto.ProdutoVariacaoResponse
	for _, variacao := range variacoes {
		if !variacao.IsDeleted() {
			resp, err := s.mapModelToResponse(&variacao)
			if err != nil {
				fmt.Printf("Erro ao mapear variação %d para resposta: %v\n", variacao.ID, err)
				continue
			}
			responses = append(responses, *resp)
		}
	}

	if len(responses) != int(total) && total > 0 {
		total = int64(len(responses))
	}

	return responses, total, nil
}
