package service

import (
	"errors"
	"fmt"

	"github.com/openerp/backend/internal/constants"
	"github.com/openerp/backend/internal/dto"
	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/utils"
	"gorm.io/gorm"
)

type EntidadeService struct {
	entidadeRepo         *repository.EntidadeRepository
	entidadeEnderecoRepo *repository.EntidadeEnderecoRepository
}

// ============================================================
// CONSTANTES DE VALIDAÇÃO
// ============================================================

const (
	maxLengthRazaoSocial  = 100
	maxLengthNomeFantasia = 100
)

func NewEntidadeService(db *gorm.DB) *EntidadeService {
	return &EntidadeService{
		entidadeRepo:         repository.NewEntidadeRepository(db),
		entidadeEnderecoRepo: repository.NewEntidadeEnderecoRepository(db),
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO (PRIVADOS)
// ============================================================

// isDataValid realiza as validações básicas de uma entidade
func (s *EntidadeService) isDataValid(req *dto.EntidadeRequest) error {
	// 1. Validar campos obrigatórios
	if err := utils.ValidateMandatoryFields(req); err != nil {
		return err
	}

	// 2. Validar tamanhos dos campos
	if err := s.validateFieldLengths(req); err != nil {
		return err
	}

	// 3. Validar documento
	if err := s.validateDocument(req); err != nil {
		return err
	}

	// 4. Validar tipo de pessoa
	if err := s.validateTipoPessoa(req); err != nil {
		return err
	}

	return nil
}

// validateFieldLengths valida o tamanho dos campos
func (s *EntidadeService) validateFieldLengths(req *dto.EntidadeRequest) error {
	// Nome / Razão Social
	if len(req.RazaoSocial) > maxLengthRazaoSocial {
		return apperrors.NewValidationError(fmt.Sprintf("nome/razão social deve ter no máximo %d caracteres", maxLengthRazaoSocial))
	}

	// Nome Fantasia (se informado)
	if len(req.NomeFantasia) > maxLengthNomeFantasia {
		return apperrors.NewValidationError(fmt.Sprintf("nome fantasia deve ter no máximo %d caracteres", maxLengthNomeFantasia))
	}

	return nil
}

// validateDocument valida o documento (CPF/CNPJ)
func (s *EntidadeService) validateDocument(req *dto.EntidadeRequest) error {
	// Remover caracteres especiais para validação
	documentoLimpo := utils.LimparDocumento(req.InscricaoFederal)

	// Validar formato (CPF ou CNPJ)
	if !utils.IsValidDocumento(documentoLimpo) {
		return apperrors.NewValidationError("documento inválido, deve ser um CPF ou CNPJ válido")
	}

	return nil
}

// validateTipoPessoa valida o tipo de pessoa
func (s *EntidadeService) validateTipoPessoa(req *dto.EntidadeRequest) error {
	// Verificar se o tipo de pessoa é válido (1-Física, 2-Jurídica)
	if req.TipoPessoa != 1 && req.TipoPessoa != 2 {
		return apperrors.NewValidationError("tipo de pessoa inválido, deve ser 1 (Física) ou 2 (Jurídica)")
	}
	return nil
}

// validateUniqueDocument verifica se o documento já existe
func (s *EntidadeService) validateUniqueDocument(documento string, excludeID int) error {
	existe, err := s.entidadeRepo.ExistsByDocumento(documento, excludeID)
	if err != nil {
		return apperrors.NewInternalError("erro ao verificar duplicidade de documento", err)
	}
	if existe {
		return apperrors.NewConflictError(fmt.Sprintf("documento %s já está cadastrado", utils.LimparDocumento(documento)))
	}
	return nil
}

// isCreateValid valida dados para criação
func (s *EntidadeService) isCreateValid(req *dto.EntidadeRequest) error {
	// 1. Validações básicas
	if err := s.isDataValid(req); err != nil {
		return err
	}

	// 2. Validar duplicidade de documento
	if err := s.validateUniqueDocument(req.InscricaoFederal, 0); err != nil {
		return err
	}

	return nil
}

// isUpdateValid valida dados para atualização
func (s *EntidadeService) isUpdateValid(id int, req *dto.EntidadeRequest) error {
	// 1. Validações básicas
	if err := s.isDataValid(req); err != nil {
		return err
	}

	// 2. Validar duplicidade de documento (excluindo o próprio ID)
	if err := s.validateUniqueDocument(req.InscricaoFederal, id); err != nil {
		return err
	}

	return nil
}

func (s *EntidadeService) Create(rep *dto.EntidadeRequest) (*models.Entidade, error) {
	if err := s.isCreateValid(rep); err != nil {
		return nil, err
	}

	entidade := &models.Entidade{}
	if err := utils.MapToModel(rep, entidade); err != nil {
		return nil, fmt.Errorf("erro ao mapear dados da entidade: %w", err)
	}

	// 3. Definir campos que não podem ser mapeados automaticamente
	entidade.InscricaoFederal = utils.LimparDocumento(rep.InscricaoFederal)
	entidade.TipoPessoa = constants.TipoPessoa(rep.TipoPessoa)
	entidade.Situacao = constants.StatusAtivo

	result := s.entidadeRepo.Create(entidade)
	if result != nil {
		return nil, result
	}
	return entidade, nil
}

func (s *EntidadeService) GetByID(id int) (*models.Entidade, error) {
	// 1. Buscar a entidade pelo ID
	// A validaçãi de se a entidade existe e não foi deletada é feita no repositório
	entidade, err := s.entidadeRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return entidade, nil
}

func (s *EntidadeService) GetByDocumento(documento string) (*models.Entidade, error) {
	if !utils.IsValidDocumento(documento) {
		return nil, fmt.Errorf("Documento inválido, deve ser um CNPJ ou CPF válido")
	}

	entidade, err := s.entidadeRepo.FindByDocumento(documento)
	if err != nil {
		return nil, err
	}
	return entidade, nil
}

func (s *EntidadeService) Update(id int, req *dto.EntidadeRequest) (*models.Entidade, error) {
	if err := s.isUpdateValid(id, req); err != nil {
		return nil, err
	}
	// 1. Buscar a entidade pelo ID
	entidade, err := s.entidadeRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := utils.MapToModel(req, entidade); err != nil {
		return nil, fmt.Errorf("erro ao mapear dados da entidade: %w", err)
	}

	entidade.InscricaoFederal = utils.LimparDocumento(req.InscricaoFederal)
	entidade.TipoPessoa = constants.TipoPessoa(req.TipoPessoa)

	result := s.entidadeRepo.Update(id, entidade)
	if result != nil {
		return nil, result
	}
	return entidade, nil
}

func (s *EntidadeService) Delete(id int) error {
	// 1. Buscar a entidade pelo ID
	entidade, err := s.entidadeRepo.FindByID(id)
	if err != nil {
		return err
	}

	if err := s.checkDependencies(id); err != nil {
		return err
	}

	return s.entidadeRepo.Delete(entidade.ID)
}

func (s *EntidadeService) List(limit, offset int, filters map[string]interface{}) ([]models.Entidade, int64, error) {
	// 1. Listar entidades com filtros

	// Validar parâmetros de paginação
	if limit <= 0 {
		limit = 10 // valor padrão
	}
	if offset < 0 {
		offset = 0
	}
	entidades, total, err := s.entidadeRepo.List(limit, offset, filters)
	if err != nil {
		return nil, 0, err
	}
	return entidades, total, nil
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO (REGRAS DE NEGÓCIO)
// ============================================================

// checkDependencies verifica se a entidade tem dependências
// ESTE MÉTODO DEVE ESTAR NO SERVICE, NÃO NO REPOSITORY!
func (s *EntidadeService) checkDependencies(entidadeID int) error {
	// 1. Verificar se tem endereços
	enderecos, err := s.entidadeEnderecoRepo.FindByEntidadeID(entidadeID)
	if err != nil {
		return fmt.Errorf("erro ao verificar endereços: %w", err)
	}
	if len(enderecos) > 0 {
		return errors.New("não é possível excluir entidade com endereços cadastrados")
	}

	// 2. TODO: Verificar se tem contatos
	// contatos, err := s.entidadeContatoRepo.FindByEntidadeID(entidadeID)
	// if err != nil {
	//     return fmt.Errorf("erro ao verificar contatos: %w", err)
	// }
	// if len(contatos) > 0 {
	//     return errors.New("não é possível excluir entidade com contatos cadastrados")
	// }

	// 3. TODO: Verificar se tem documentos
	// documentos, err := s.entidadeDocumentoRepo.FindByEntidadeID(entidadeID)
	// if err != nil {
	//     return fmt.Errorf("erro ao verificar documentos: %w", err)
	// }
	// if len(documentos) > 0 {
	//     return errors.New("não é possível excluir entidade com documentos cadastrados")
	// }

	// 4. TODO: Verificar se tem pedidos (documento_venda)
	// pedidos, err := s.documentoVendaRepo.FindByEntidadeID(entidadeID)
	// if err != nil {
	//     return fmt.Errorf("erro ao verificar pedidos: %w", err)
	// }
	// if len(pedidos) > 0 {
	//     return errors.New("não é possível excluir entidade com pedidos associados")
	// }

	// 5. TODO: Verificar se tem notas fiscais
	// notas, err := s.notaFiscalRepo.FindByEntidadeID(entidadeID)
	// if err != nil {
	//     return fmt.Errorf("erro ao verificar notas fiscais: %w", err)
	// }
	// if len(notas) > 0 {
	//     return errors.New("não é possível excluir entidade com notas fiscais associadas")
	// }

	return nil
}
