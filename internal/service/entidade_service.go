package service

import (
	"fmt"

	"github.com/openerp/backend/internal/constants"
	"github.com/openerp/backend/internal/dto"
	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/utils"
)

// EntidadeService define os métodos públicos para o serviço de entidade.
type EntidadeService interface {
	Create(rep *dto.EntidadeRequest) (*models.Entidade, error)
	GetByID(id int) (*models.Entidade, error)
	GetByDocumento(documento string) (*models.Entidade, error)
	Update(id int, req *dto.EntidadeRequest) (*models.Entidade, error)
	Delete(id int) error
	List(limit, offset int, filters map[string]interface{}) ([]models.Entidade, int64, error)
}

type entidadeService struct {
	entidadeRepo            repository.EntidadeRepository
	entidadeEnderecoService EntidadeEnderecoService // Already an interface, just need to update constructor
}

// ============================================================
// CONSTANTES DE VALIDAÇÃO
// ============================================================

const (
	maxLengthRazaoSocial  = 100
	maxLengthNomeFantasia = 100
)

func NewEntidadeService(
	entRepo repository.EntidadeRepository,
	entEnderecoService EntidadeEnderecoService,
) EntidadeService {
	return &entidadeService{
		entidadeRepo:            entRepo,
		entidadeEnderecoService: entEnderecoService,
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO (PRIVADOS)
// ============================================================

// isDataValid realiza as validações básicas de uma entidade
func (s *entidadeService) isDataValid(req *dto.EntidadeRequest) error {
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
func (s *entidadeService) validateFieldLengths(req *dto.EntidadeRequest) error {
	// Nome / Razão Social
	if len(req.RazaoSocial) > maxLengthRazaoSocial {
		return apperrors.NewValidationError(fmt.Sprintf("Nome/razão social deve ter no máximo %d caracteres.", maxLengthRazaoSocial)) //
	}

	// Nome Fantasia (se informado)
	if len(req.NomeFantasia) > maxLengthNomeFantasia {
		return apperrors.NewValidationError(fmt.Sprintf("Nome fantasia deve ter no máximo %d caracteres.", maxLengthNomeFantasia)) //
	}

	return nil
}

// validateDocument valida o documento (CPF/CNPJ)
func (s *entidadeService) validateDocument(req *dto.EntidadeRequest) error {
	// Remover caracteres especiais para validação
	documentoLimpo := utils.LimparDocumento(req.InscricaoFederal)

	// Validar formato (CPF ou CNPJ)
	if !utils.IsValidDocumento(documentoLimpo) {
		return apperrors.NewValidationError("Documento inválido, deve ser um CPF ou CNPJ válido.") //
	}

	return nil
}

// validateTipoPessoa valida o tipo de pessoa
func (s *entidadeService) validateTipoPessoa(req *dto.EntidadeRequest) error {
	// Verificar se o tipo de pessoa é válido (1-Física, 2-Jurídica)
	if req.TipoPessoa != 1 && req.TipoPessoa != 2 {
		return apperrors.NewValidationError("Tipo de pessoa inválido, deve ser 1 (Física) ou 2 (Jurídica).") //
	}
	return nil
}

// validateUniqueDocument verifica se o documento já existe
func (s *entidadeService) validateUniqueDocument(documento string, excludeID int) error {
	existe, err := s.entidadeRepo.ExistsByDocumento(documento, excludeID)
	if err != nil {
		return apperrors.NewInternalError("Erro ao verificar duplicidade de documento.", err) //
	}
	if existe {
		return apperrors.NewConflictError(fmt.Sprintf("Documento %s já está cadastrado.", utils.LimparDocumento(documento))) //
	}
	return nil
}

// isCreateValid valida dados para criação
func (s *entidadeService) isCreateValid(req *dto.EntidadeRequest) error {
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
func (s *entidadeService) isUpdateValid(id int, req *dto.EntidadeRequest) error {
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

func (s *entidadeService) Create(rep *dto.EntidadeRequest) (*models.Entidade, error) {
	if err := s.isCreateValid(rep); err != nil {
		return nil, err
	}

	entidade := &models.Entidade{}
	if err := utils.MapToModel(rep, entidade); err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear dados da entidade.", err) //
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

func (s *entidadeService) GetByID(id int) (*models.Entidade, error) {
	// 1. Buscar a entidade pelo ID
	// A validaçãi de se a entidade existe e não foi deletada é feita no repositório
	entidade, err := s.entidadeRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return entidade, nil
}

func (s *entidadeService) GetByDocumento(documento string) (*models.Entidade, error) {
	if !utils.IsValidDocumento(documento) {
		return nil, apperrors.NewValidationError("Documento inválido, deve ser um CNPJ ou CPF válido.") //
	}

	entidade, err := s.entidadeRepo.FindByDocumento(documento)
	if err != nil {
		return nil, err //
	}
	return entidade, nil
}

func (s *entidadeService) Update(id int, req *dto.EntidadeRequest) (*models.Entidade, error) {
	if err := s.isUpdateValid(id, req); err != nil {
		return nil, err
	}
	// 1. Buscar a entidade pelo ID
	entidade, err := s.entidadeRepo.FindByID(id)
	if err != nil {
		return nil, err //
	}

	if err := utils.MapToModel(req, entidade); err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear dados da entidade.", err) //
	}

	entidade.InscricaoFederal = utils.LimparDocumento(req.InscricaoFederal)
	entidade.TipoPessoa = constants.TipoPessoa(req.TipoPessoa)

	result := s.entidadeRepo.Update(id, entidade)
	if result != nil {
		return nil, result
	}
	return entidade, nil
}

func (s *entidadeService) Delete(id int) error {
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

func (s *entidadeService) List(limit, offset int, filters map[string]interface{}) ([]models.Entidade, int64, error) {
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
func (s *entidadeService) checkDependencies(entidadeID int) error {
	// 1. Verificar se tem endereços
	enderecos, err := s.entidadeEnderecoService.GetByEntidadeID(entidadeID)
	if err != nil { //
		return apperrors.NewInternalError("Erro ao verificar endereços.", err) //
	} //
	if len(enderecos) > 0 { //
		return apperrors.NewConflictError("Não é possível excluir entidade com endereços cadastrados.") //
	}

	// 2. TODO: Verificar se tem contatos
	// contatos, err := s.entidadeContatoRepo.FindByEntidadeID(entidadeID)
	// if err != nil { //
	//     return apperrors.NewInternalError("Erro ao verificar contatos.", err) //
	// } //
	// if len(contatos) > 0 { //
	//     return apperrors.NewConflictError("Não é possível excluir entidade com contatos cadastrados.") //
	// }

	// 3. TODO: Verificar se tem documentos
	// documentos, err := s.entidadeDocumentoRepo.FindByEntidadeID(entidadeID)
	// if err != nil { //
	//     return apperrors.NewInternalError("Erro ao verificar documentos.", err) //
	// } //
	// if len(documentos) > 0 { //
	//     return apperrors.NewConflictError("Não é possível excluir entidade com documentos cadastrados.") //
	// }

	// 4. TODO: Verificar se tem pedidos (documento_venda)
	// pedidos, err := s.documentoVendaRepo.FindByEntidadeID(entidadeID)
	// if err != nil { //
	//     return apperrors.NewInternalError("Erro ao verificar pedidos.", err) //
	// } //
	// if len(pedidos) > 0 { //
	//     return apperrors.NewConflictError("Não é possível excluir entidade com pedidos associados.") //
	// }

	// 5. TODO: Verificar se tem notas fiscais
	// notas, err := s.notaFiscalRepo.FindByEntidadeID(entidadeID)
	// if err != nil { //
	//     return apperrors.NewInternalError("Erro ao verificar notas fiscais.", err) //
	// } //
	// if len(notas) > 0 { //
	//     return apperrors.NewConflictError("Não é possível excluir entidade com notas fiscais associadas.") //
	// }

	return nil
}
