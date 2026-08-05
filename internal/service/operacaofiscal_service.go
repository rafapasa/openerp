package service

import (
	"fmt"
	"strings"

	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
)

// OperacaoFiscalService define os métodos públicos para o serviço de operação fiscal.
type OperacaoFiscalService interface {
	FindByID(id int) (*models.OperacaoFiscal, error)
	Create(req *dto.OperacaoFiscalRequest) (*models.OperacaoFiscal, error)
	Update(id int, req *dto.OperacaoFiscalRequest) (*models.OperacaoFiscal, error)
	Delete(id int) error
	List(limit, offset int, filters map[string]interface{}) ([]models.OperacaoFiscal, int64, error)
	FindByCFOP(cfop string) ([]models.OperacaoFiscal, error)
	FindByEmpresaFilialID(filialID int) ([]models.OperacaoFiscal, error)
}

type operacaoFiscalService struct {
	opfRepo repository.OperacaoFiscalRepository
}

func NewOperacaoFiscalService(opfRepo repository.OperacaoFiscalRepository) OperacaoFiscalService {
	return &operacaoFiscalService{
		opfRepo: opfRepo,
	}
}

func (s *operacaoFiscalService) FindByID(id int) (*models.OperacaoFiscal, error) {
	if id <= 0 {
		return nil, apperrors.NewValidationError("ID da operação fiscal inválido.")
	}
	return s.opfRepo.FindByID(id)
}

// Create cria uma nova operação fiscal.
func (s *operacaoFiscalService) Create(req *dto.OperacaoFiscalRequest) (*models.OperacaoFiscal, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Validar unicidade de CFOP por filial
	exists, err := s.opfRepo.ExistsByCFOPAndFilial(req.CFOP, req.EmpresaFilialID, 0)
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao verificar CFOP existente.", err)
	}
	if exists {
		return nil, apperrors.NewConflictError(fmt.Sprintf("Já existe uma operação fiscal com o CFOP '%s' para a filial %d.", req.CFOP, req.EmpresaFilialID))
	}

	operacao, err := req.ToModel()
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear DTO para modelo de operação fiscal.", err)
	}

	if err := s.opfRepo.Create(operacao); err != nil {
		return nil, apperrors.NewInternalError("Erro ao criar operação fiscal.", err)
	}

	return operacao, nil
}

// Update atualiza uma operação fiscal existente.
func (s *operacaoFiscalService) Update(id int, req *dto.OperacaoFiscalRequest) (*models.OperacaoFiscal, error) {
	if id <= 0 {
		return nil, apperrors.NewValidationError("ID da operação fiscal inválido.")
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}

	existingOperacao, err := s.opfRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Validar unicidade de CFOP por filial, excluindo o próprio registro
	exists, err := s.opfRepo.ExistsByCFOPAndFilial(req.CFOP, req.EmpresaFilialID, id)
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao verificar CFOP existente.", err)
	}
	if exists {
		return nil, apperrors.NewConflictError(fmt.Sprintf("Já existe uma operação fiscal com o CFOP '%s' para a filial %d.", req.CFOP, req.EmpresaFilialID))
	}

	updatedModel, err := req.ToModel()
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear DTO para modelo de operação fiscal.", err)
	}

	// Atualiza os campos do modelo existente
	existingOperacao.EmpresaFilialID = updatedModel.EmpresaFilialID
	existingOperacao.CFOP = updatedModel.CFOP
	existingOperacao.Descricao = updatedModel.Descricao
	existingOperacao.DataIni = updatedModel.DataIni
	existingOperacao.DataFim = updatedModel.DataFim
	existingOperacao.CSTICMSID = updatedModel.CSTICMSID
	existingOperacao.CSTIPIID = updatedModel.CSTIPIID
	existingOperacao.CSTPISCOFINSID = updatedModel.CSTPISCOFINSID
	existingOperacao.UpdatedBy = req.UpdatedBy

	if err := s.opfRepo.Update(id, existingOperacao); err != nil {
		return nil, apperrors.NewInternalError("Erro ao atualizar operação fiscal.", err)
	}

	return existingOperacao, nil
}

// Delete exclui logicamente uma operação fiscal.
func (s *operacaoFiscalService) Delete(id int) error {
	if id <= 0 {
		return apperrors.NewValidationError("ID da operação fiscal inválido.")
	}
	return s.opfRepo.Delete(id)
}

// List lista operações fiscais com paginação e filtros.
func (s *operacaoFiscalService) List(limit, offset int, filters map[string]interface{}) ([]models.OperacaoFiscal, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	operacoes, total, err := s.opfRepo.List(limit, offset, filters)
	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro ao listar operações fiscais.", err)
	}
	return operacoes, total, nil
}

// FindByCFOP busca operações fiscais por CFOP.
func (s *operacaoFiscalService) FindByCFOP(cfop string) ([]models.OperacaoFiscal, error) {
	if strings.TrimSpace(cfop) == "" {
		return nil, apperrors.NewValidationError("CFOP não pode ser vazio.")
	}
	operacoes, err := s.opfRepo.FindByCFOP(cfop)
	if err != nil {
		return nil, err
	}
	return operacoes, nil
}

// FindByEmpresaFilialID busca operações fiscais por ID da filial.
func (s *operacaoFiscalService) FindByEmpresaFilialID(filialID int) ([]models.OperacaoFiscal, error) {
	if filialID <= 0 {
		return nil, apperrors.NewValidationError("ID da filial inválido.")
	}
	operacoes, err := s.opfRepo.FindByEmpresaFilialID(filialID)
	if err != nil {
		return nil, err
	}
	return operacoes, nil
}
