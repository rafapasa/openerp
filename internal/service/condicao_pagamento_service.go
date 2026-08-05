package service

import (
	"fmt"
	"strings"

	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/constants"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
)

type CondicaoPagamentoService interface {
	Create(req *dto.CondicaoPagamentoRequest) (*models.CondicaoPagamento, error)
	GetByID(id int) (*models.CondicaoPagamento, error)
	Update(id int, req *dto.CondicaoPagamentoRequest) (*models.CondicaoPagamento, error)
	Delete(id int) error
	List(limit, offset int, filters map[string]interface{}) ([]models.CondicaoPagamento, int64, error)
}

type condicaoPagamentoService struct {
	cdpgtRepo repository.CondicaoPagamentoRepository
}

func NewCondicaoPagamentoService(cdpgtRepo repository.CondicaoPagamentoRepository) CondicaoPagamentoService {
	return &condicaoPagamentoService{
		cdpgtRepo: cdpgtRepo,
	}
}

// Create cria uma nova condição de pagamento.
// Este método já existe no arquivo, apenas garantindo que a interface o referencie.
// (No diff, ele já está presente, então não há mudança real aqui, apenas a menção para a interface)

// Create cria uma nova condição de pagamento.
func (s *condicaoPagamentoService) Create(req *dto.CondicaoPagamentoRequest) (*models.CondicaoPagamento, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.cdpgtRepo.ExistsByDescricao(descricao, 0)
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao verificar descrição da condição de pagamento.", err)
	}
	if exists {
		return nil, apperrors.NewConflictError("Já existe uma condição de pagamento com esta descrição.")
	}

	condicaoPagamento, err := req.ToModel()

	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear DTO para modelo de condição de pagamento.", err)
	}

	// Definir situação padrão se não informada
	if condicaoPagamento.Situacao == 0 {
		condicaoPagamento.Situacao = constants.ATIVO
	}

	if err := s.cdpgtRepo.Create(condicaoPagamento); err != nil {
		return nil, apperrors.NewInternalError("Erro ao criar condição de pagamento.", err)
	}

	return condicaoPagamento, nil
}

// GetByID busca uma condição de pagamento pelo ID.
func (s *condicaoPagamentoService) GetByID(id int) (*models.CondicaoPagamento, error) {
	if id <= 0 {
		return nil, apperrors.NewValidationError("ID da condição de pagamento inválido.")
	}
	condicaoPagamento, err := s.cdpgtRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return condicaoPagamento, nil
}

// Update atualiza uma condição de pagamento existente.
func (s *condicaoPagamentoService) Update(id int, req *dto.CondicaoPagamentoRequest) (*models.CondicaoPagamento, error) {
	if id <= 0 {
		return nil, apperrors.NewValidationError("ID da condição de pagamento inválido.")
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	_, err := s.cdpgtRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	descricao := strings.TrimSpace(req.Descricao)
	exists, err := s.cdpgtRepo.ExistsByDescricao(descricao, id)
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao verificar descrição da condição de pagamento.", err)
	}
	if exists {
		return nil, apperrors.NewConflictError("Já existe uma condição de pagamento com esta descrição.")
	}

	condicao, e := req.ToModel()
	if e != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear DTO para modelo de condição de pagamento.", err)
	}

	if err := s.cdpgtRepo.Update(id, condicao); err != nil {
		return nil, apperrors.NewInternalError("Erro ao atualizar condição de pagamento.", err)
	}

	return condicao, nil
}

// Delete exclui logicamente uma condição de pagamento.
func (s *condicaoPagamentoService) Delete(id int) error {
	if id <= 0 {
		return apperrors.NewValidationError("ID da condição de pagamento inválido.")
	}

	_, err := s.cdpgtRepo.FindByID(id)
	if err != nil {
		return err
	}

	// TODO: Adicionar verificação de dependências (ex: se a condição de pagamento está em uso por algum documento de venda)
	// count, err := s.documentoVendaRepo.CountByCondicaoPagamentoID(id)
	// if err != nil { return err }
	// if count > 0 { return apperrors.NewConflictError(fmt.Sprintf("Condição de pagamento em uso por %d documento(s) e não pode ser excluída.", count)) }

	if err := s.cdpgtRepo.Delete(id); err != nil {
		return apperrors.NewInternalError(fmt.Sprintf("Erro ao excluir condição de pagamento com ID %d.", id), err)
	}
	return nil
}

// List lista condições de pagamento com paginação e filtros.
func (s *condicaoPagamentoService) List(limit, offset int, filters map[string]interface{}) ([]models.CondicaoPagamento, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	condicoes, total, err := s.cdpgtRepo.List(limit, offset, filters)
	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro ao listar condições de pagamento.", err)
	}
	return condicoes, total, nil
}
