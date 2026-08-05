package service

import (
	"fmt"

	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
)

// ProcessoService define os métodos públicos para o serviço de processo.
type ProcessoService interface {
	GetOperacaoFiscal(prcId int, opInterna, opSubTrib bool) (*models.OperacaoFiscal, error)
	Create(req *dto.ProcessoRequest) (*models.Processo, error)
	Update(id int, req *dto.ProcessoRequest) (*models.Processo, error)
	Delete(id int) error
	List(limit, offset int, filters map[string]interface{}) ([]models.Processo, int64, error)
	FindByID(id int) (*models.Processo, error)
	FindByCodigo(codigo int) (*models.Processo, error)
}

type processoService struct {
	prcRepo    repository.ProcessoRepository
	opfService OperacaoFiscalService
}

func NewProcessoService(prcRepo repository.ProcessoRepository, opfService OperacaoFiscalService) ProcessoService {
	return &processoService{
		prcRepo:    prcRepo,
		opfService: opfService,
	}
}

func (s *processoService) GetOperacaoFiscal(prcId int, opInterna, opSubTrib bool) (*models.OperacaoFiscal, error) {
	prc, err := s.prcRepo.FindByID(prcId)
	if err != nil {
		return nil, err
	}
	if opInterna && opSubTrib {
		return s.opfService.FindByID(*prc.OperacaoFiscalNoEstSTID)
	}
	if opInterna && !opSubTrib {
		return s.opfService.FindByID(*prc.OperacaoFiscalNoEstID)
	}
	if !opInterna && opSubTrib {
		return s.opfService.FindByID(*prc.OperacaoFiscalForaEstSTID)
	}
	if !opInterna && !opSubTrib { // Corrected: Changed condition
		return s.opfService.FindByID(*prc.OperacaoFiscalForaEstID)
	}
	return nil, apperrors.NewInternalError("Operação fiscal não encontrada.", nil)
}

func (s *processoService) FindByID(id int) (*models.Processo, error) { // Added this method to satisfy the interface
	if id <= 0 {
		return nil, apperrors.NewValidationError("ID do processo inválido.")
	}
	return s.prcRepo.FindByID(id)

}

func (s *processoService) FindByCodigo(codigo int) (*models.Processo, error) { // Added this method to satisfy the interface
	if codigo <= 0 {
		return nil, apperrors.NewValidationError("ID do processo inválido.")
	}
	return s.prcRepo.FindByCodigo(codigo)

}

// Create cria um novo processo.
func (s *processoService) Create(req *dto.ProcessoRequest) (*models.Processo, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Validar unicidade de código por filial
	exists, err := s.prcRepo.ExistsByCodigoAndFilial(req.Codigo, req.EmpresaFilialID, 0)
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao verificar código existente.", err)
	}
	if exists {
		return nil, apperrors.NewConflictError(fmt.Sprintf("Já existe um processo com o código '%d' para a filial %d.", req.Codigo, req.EmpresaFilialID))
	}

	processo, err := req.ToModel()
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear DTO para modelo de processo.", err)
	}

	if err := s.prcRepo.Create(processo); err != nil {
		return nil, apperrors.NewInternalError("Erro ao criar processo.", err)
	}

	return processo, nil
}

// Update atualiza um processo existente.
func (s *processoService) Update(id int, req *dto.ProcessoRequest) (*models.Processo, error) {
	if id <= 0 {
		return nil, apperrors.NewValidationError("ID do processo inválido.")
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}

	existingProcesso, err := s.prcRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Validar unicidade de código por filial, excluindo o próprio registro
	exists, err := s.prcRepo.ExistsByCodigoAndFilial(req.Codigo, req.EmpresaFilialID, id)
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao verificar código existente.", err)
	}
	if exists {
		return nil, apperrors.NewConflictError(fmt.Sprintf("Já existe um processo com o código '%d' para a filial %d.", req.Codigo, req.EmpresaFilialID))
	}

	updatedModel, err := req.ToModel()
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear DTO para modelo de processo.", err)
	}

	// Atualiza os campos do modelo existente
	existingProcesso.EmpresaFilialID = updatedModel.EmpresaFilialID
	existingProcesso.Codigo = updatedModel.Codigo
	existingProcesso.Descricao = updatedModel.Descricao
	existingProcesso.TipoOperacao = updatedModel.TipoOperacao
	existingProcesso.Situacao = updatedModel.Situacao
	existingProcesso.UpdatedBy = req.UpdatedBy

	if err := s.prcRepo.Update(id, existingProcesso); err != nil {
		return nil, apperrors.NewInternalError("Erro ao atualizar processo.", err)
	}

	return existingProcesso, nil
}

// Delete exclui logicamente um processo.
func (s *processoService) Delete(id int) error {
	if id <= 0 {
		return apperrors.NewValidationError("ID do processo inválido.")
	}
	return s.prcRepo.Delete(id)
}

// List lista processos com paginação e filtros.
func (s *processoService) List(limit, offset int, filters map[string]interface{}) ([]models.Processo, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	processos, total, err := s.prcRepo.List(limit, offset, filters)
	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro ao listar processos.", err)
	}
	return processos, total, nil
}

// GetByCodigo busca um processo pelo código.
func (s *processoService) GetByCodigo(codigo int) (*models.Processo, error) {
	if codigo <= 0 {
		return nil, apperrors.NewValidationError("Código do processo inválido.")
	}
	processo, err := s.prcRepo.FindByCodigo(codigo)
	if err != nil {
		return nil, err
	}
	return processo, nil
}
