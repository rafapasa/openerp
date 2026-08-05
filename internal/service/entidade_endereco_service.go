package service

import (
	"fmt"
	"strings"

	"github.com/openerp/backend/internal/apperrors"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/utils"
)

// EntidadeEnderecoService define os métodos públicos para o serviço de endereço de entidade.
type EntidadeEnderecoService interface {
	Create(req *dto.EntidadeEnderecoRequest) (*models.EntidadeEndereco, error)
	GetByID(entidadeID, item int) (*models.EntidadeEndereco, error)
	GetByEntidadeID(entidadeID int) ([]models.EntidadeEndereco, error)
	GetByEntidadeIDAndTipo(entidadeID, tipo int) ([]models.EntidadeEndereco, error)
	Update(entidadeID, item int, req *dto.EntidadeEnderecoRequest) (*models.EntidadeEndereco, error)
	Delete(entidadeID, item int) error
	List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeEndereco, int64, error)
}

// ============================================================
// TYPES
// ============================================================

type entidadeEnderecoService struct {
	entidadeEnderecoRepo repository.EntidadeEnderecoRepository
	entidadeRepo         repository.EntidadeRepository // This should probably be EntidadeService, but keeping as is for now.
}

// ============================================================
// CONSTRUCTOR
// ============================================================

func NewEntidadeEnderecoService(db *gorm.DB) EntidadeEnderecoService {
	return &entidadeEnderecoService{
		entidadeEnderecoRepo: repository.NewEntidadeEnderecoRepository(db),
		entidadeRepo:         repository.NewEntidadeRepository(db),
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO (PRIVADOS)
// ============================================================

// isDataValid realiza as validações básicas de um endereço
func (s *entidadeEnderecoService) isDataValid(req *dto.EntidadeEnderecoRequest) error {
	// 1. Validar campos obrigatórios
	if err := s.validateRequiredFields(req); err != nil {
		return err
	}

	// 2. Validar tipo de endereço
	if err := s.validateTipoEndereco(req); err != nil {
		return err
	}

	// 3. Validar datas
	if err := s.validateDates(req); err != nil {
		return err
	}

	return nil
}

// validateRequiredFields valida campos obrigatórios
func (s *entidadeEnderecoService) validateRequiredFields(req *dto.EntidadeEnderecoRequest) error {
	if strings.TrimSpace(req.Logradouro) == "" {
		return apperrors.NewValidationError("Logradouro é obrigatório.") //
	}

	if strings.TrimSpace(req.Numero) == "" {
		return apperrors.NewValidationError("Número é obrigatório.") //
	}

	if strings.TrimSpace(req.Bairro) == "" {
		return apperrors.NewValidationError("Bairro é obrigatório.") //
	}

	if strings.TrimSpace(req.CEP) == "" {
		return apperrors.NewValidationError("CEP é obrigatório.") //
	}

	if strings.TrimSpace(req.DataIni) == "" {
		return apperrors.NewValidationError("Data inicial é obrigatória.") //
	}

	return nil
}

// validateTipoEndereco valida o tipo de endereço
func (s *entidadeEnderecoService) validateTipoEndereco(req *dto.EntidadeEnderecoRequest) error {
	if req.Tipo < 1 || req.Tipo > 4 {
		return apperrors.NewValidationError("Tipo de endereço inválido, deve ser 1 (Cobrança), 2 (Entrega), 3 (Comercial) ou 4 (Residencial).") //
	}
	return nil
}

// validateDates valida as datas
func (s *entidadeEnderecoService) validateDates(req *dto.EntidadeEnderecoRequest) error {
	// Validar formato da data inicial
	if req.DataIni != "" {
		if _, err := utils.ParseDate(req.DataIni); err != nil { //
			return apperrors.NewValidationError(fmt.Sprintf("Data inicial inválida: %v", err)) //
		}
	}

	// Validar formato da data final (se informada)
	if req.DataFim != "" {
		if _, err := utils.ParseDate(req.DataFim); err != nil { //
			return apperrors.NewValidationError(fmt.Sprintf("Data final inválida: %v", err)) //
		}
	}

	return nil
}

// validateEntidadeExists verifica se a entidade existe
func (s *entidadeEnderecoService) validateEntidadeExists(entidadeID int) error {
	_, err := s.entidadeRepo.FindByID(entidadeID)
	if err != nil { //
		return apperrors.NewNotFoundError(fmt.Sprintf("Entidade com ID %d não encontrada.", entidadeID)) //
	}
	return nil
}

// isCreateValid valida dados para criação
func (s *entidadeEnderecoService) isCreateValid(req *dto.EntidadeEnderecoRequest) error {
	// 1. Validações básicas
	if err := s.isDataValid(req); err != nil {
		return err
	}

	// 2. Validar se a entidade existe
	if err := s.validateEntidadeExists(req.EntidadeID); err != nil {
		return err
	}

	// 3. Verificar se já existe um endereço do mesmo tipo (opcional)
	// existe, err := s.entidadeEnderecoRepo.ExistsByEntidadeTipo(req.EntidadeID, req.Tipo, 0)
	// if err != nil {
	//     return fmt.Errorf("erro ao verificar endereço existente: %w", err)
	// }
	// if existe {
	//     return errors.New("já existe um endereço deste tipo para esta entidade")
	// }

	return nil
}

// isUpdateValid valida dados para atualização
func (s *entidadeEnderecoService) isUpdateValid(entidadeID, item int, req *dto.EntidadeEnderecoRequest) error {
	// 1. Validações básicas
	if err := s.isDataValid(req); err != nil {
		return err
	}

	// 2. Validar se a entidade existe
	if err := s.validateEntidadeExists(entidadeID); err != nil {
		return err
	} //

	// 3. Validar se o endereço existe
	if _, err := s.entidadeEnderecoRepo.FindByID(entidadeID, item); err != nil {
		return apperrors.NewNotFoundError(fmt.Sprintf("Endereço %d da entidade %d não encontrado.", item, entidadeID)) //
	}

	return nil
}

// ============================================================
// MÉTODOS PRINCIPAIS (CRUD)
// ============================================================

// Create cria um novo endereço para uma entidade
func (s *entidadeEnderecoService) Create(req *dto.EntidadeEnderecoRequest) (*models.EntidadeEndereco, error) {
	// 1. Validar dados
	if err := s.isCreateValid(req); err != nil {
		return nil, err
	}

	// 2. Converter DTO para Model
	endereco, err := req.ToModel()
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao converter dados.", err) //
	}

	// 3. Salvar
	if err := s.entidadeEnderecoRepo.Create(endereco); err != nil {
		return nil, apperrors.NewInternalError("Erro ao criar endereço.", err) //
	}

	return endereco, nil
}

// GetByID busca um endereço específico
func (s *entidadeEnderecoService) GetByID(entidadeID, item int) (*models.EntidadeEndereco, error) {
	endereco, err := s.entidadeEnderecoRepo.FindByID(entidadeID, item)
	if err != nil { //
		return nil, apperrors.NewNotFoundError(fmt.Sprintf("Endereço %d da entidade %d não encontrado.", item, entidadeID)) //
	}
	return endereco, nil
}

// GetByEntidadeID busca todos os endereços de uma entidade
func (s *entidadeEnderecoService) GetByEntidadeID(entidadeID int) ([]models.EntidadeEndereco, error) {
	// Validar se a entidade existe
	if err := s.validateEntidadeExists(entidadeID); err != nil {
		return nil, err
	}

	enderecos, err := s.entidadeEnderecoRepo.FindByEntidadeID(entidadeID)
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao buscar endereços.", err) //
	}

	return enderecos, nil
}

// GetByEntidadeIDAndTipo busca endereços de uma entidade por tipo
func (s *entidadeEnderecoService) GetByEntidadeIDAndTipo(entidadeID, tipo int) ([]models.EntidadeEndereco, error) {
	// Validar se a entidade existe
	if err := s.validateEntidadeExists(entidadeID); err != nil {
		return nil, err
	}

	// Validar tipo
	if err := s.validateTipoEndereco(&dto.EntidadeEnderecoRequest{Tipo: tipo}); err != nil {
		return nil, err
	}

	enderecos, err := s.entidadeEnderecoRepo.FindByEntidadeIDAndTipo(entidadeID, tipo)
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao buscar endereços.", err) //
	}

	return enderecos, nil
}

// Update atualiza um endereço existente
func (s *entidadeEnderecoService) Update(entidadeID, item int, req *dto.EntidadeEnderecoRequest) (*models.EntidadeEndereco, error) {
	// 1. Validar dados
	if err := s.isUpdateValid(entidadeID, item, req); err != nil {
		return nil, err
	}

	// 2. Buscar endereço existente
	endereco, err := s.entidadeEnderecoRepo.FindByID(entidadeID, item)
	if err != nil {
		return nil, apperrors.NewNotFoundError(fmt.Sprintf("Endereço %d da entidade %d não encontrado.", item, entidadeID)) //
	}

	// 3. Atualizar campos
	endereco.PaisID = req.PaisID
	endereco.EstadoID = req.EstadoID
	endereco.MunicipioID = req.MunicipioID
	endereco.Tipo = req.Tipo
	endereco.CEP = utils.ParseIntOrDefault(req.CEP, 0)
	endereco.Logradouro = utils.StringPtr(req.Logradouro)
	endereco.Numero = req.Numero
	endereco.Complemento = utils.StringPtr(req.Complemento)
	endereco.Bairro = req.Bairro
	endereco.Distancia = req.Distancia
	endereco.Observacao = utils.StringPtr(req.Observacao)

	// 4. Atualizar situação (se informada)
	if req.Situacao > 0 {
		endereco.Situacao = req.Situacao
	}

	// 5. Atualizar datas
	if req.DataIni != "" {
		if data, err := utils.ParseDate(req.DataIni); err == nil {
			endereco.DataIni = data
		}
	}
	if req.DataFim != "" {
		if data, err := utils.ParseDate(req.DataFim); err == nil {
			endereco.DataFim = &data
		}
	}

	// 6. Atualizar auditoria
	if req.UpdatedBy != nil {
		endereco.UpdatedBy = req.UpdatedBy
	}

	// 7. Salvar
	if err := s.entidadeEnderecoRepo.Update(endereco); err != nil {
		return nil, apperrors.NewInternalError("Erro ao atualizar endereço.", err) //
	}

	return endereco, nil
}

// Delete exclui logicamente um endereço
func (s *entidadeEnderecoService) Delete(entidadeID, item int) error {
	// 1. Validar se o endereço existe
	endereco, err := s.entidadeEnderecoRepo.FindByID(entidadeID, item)
	if err != nil { //
		return apperrors.NewNotFoundError(fmt.Sprintf("Endereço %d da entidade %d não encontrado.", item, entidadeID)) //
	}

	// 2. Verificar se já foi deletado
	if endereco.IsDeleted() { //
		return apperrors.NewConflictError("Endereço já foi deletado.") //
	}

	// 3. TODO: Verificar se o endereço está sendo usado em algum lugar
	// Exemplo: verificar em documento_venda, nota_fiscal, etc.

	// 4. Excluir
	if err := s.entidadeEnderecoRepo.Delete(entidadeID, item); err != nil {
		return apperrors.NewInternalError("Erro ao excluir endereço.", err) //
	}

	return nil
}

// List lista endereços com paginação e filtros
func (s *entidadeEnderecoService) List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeEndereco, int64, error) {
	// Validar parâmetros de paginação
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	enderecos, total, err := s.entidadeEnderecoRepo.List(limit, offset, filters)
	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro ao listar endereços.", err) //
	}

	return enderecos, total, nil
}
