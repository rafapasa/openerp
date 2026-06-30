package service

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
)

// ============================================================
// TYPES
// ============================================================

// EntidadeContatoService é o serviço para contatos de entidade
type EntidadeContatoService struct {
	contatoRepo  *repository.EntidadeContatoRepository
	entidadeRepo *repository.EntidadeRepository
}

// ============================================================
// CONSTRUCTOR
// ============================================================

// NewEntidadeContatoService cria uma nova instância
func NewEntidadeContatoService(db *gorm.DB) *EntidadeContatoService {
	return &EntidadeContatoService{
		contatoRepo:  repository.NewEntidadeContatoRepository(db),
		entidadeRepo: repository.NewEntidadeRepository(db),
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO (PRIVADOS)
// ============================================================

// isDataValid realiza as validações básicas de um contato
func (s *EntidadeContatoService) isDataValid(req *dto.EntidadeContatoRequest) error {
	// 1. Validar campos obrigatórios
	if err := s.validateRequiredFields(req); err != nil {
		return err
	}

	// 2. Validar tipo de contato
	if err := s.validateTipoContato(req); err != nil {
		return err
	}

	return nil
}

// validateRequiredFields valida campos obrigatórios
func (s *EntidadeContatoService) validateRequiredFields(req *dto.EntidadeContatoRequest) error {
	if strings.TrimSpace(req.Informacao) == "" {
		return errors.New("informação do contato é obrigatória")
	}

	if req.FormaContatoID == 0 {
		return errors.New("forma de contato é obrigatória")
	}

	return nil
}

// validateTipoContato valida o tipo de contato
func (s *EntidadeContatoService) validateTipoContato(req *dto.EntidadeContatoRequest) error {
	if !isValidTipoContato(req.FormaContatoID) {
		return fmt.Errorf("tipo de contato inválido: %d", req.FormaContatoID)
	}
	return nil
}

// validateEntidadeExists verifica se a entidade existe
func (s *EntidadeContatoService) validateEntidadeExists(entidadeID int) error {
	_, err := s.entidadeRepo.FindByID(entidadeID)
	if err != nil {
		return fmt.Errorf("entidade não encontrada: %w", err)
	}
	return nil
}

// isCreateValid valida dados para criação
func (s *EntidadeContatoService) isCreateValid(req *dto.EntidadeContatoRequest) error {
	// 1. Validações básicas
	if err := s.isDataValid(req); err != nil {
		return err
	}

	// 2. Validar se a entidade existe
	if err := s.validateEntidadeExists(req.EntidadeID); err != nil {
		return err
	}

	// 3. Verificar se já existe um contato do mesmo tipo (opcional)
	// existe, err := s.contatoRepo.ExistsByEntidadeTipo(req.EntidadeID, req.FormaContatoID, 0)
	// if err != nil {
	//     return fmt.Errorf("erro ao verificar contato existente: %w", err)
	// }
	// if existe {
	//     return errors.New("já existe um contato deste tipo para esta entidade")
	// }

	return nil
}

// isUpdateValid valida dados para atualização
func (s *EntidadeContatoService) isUpdateValid(entidadeID, item int, req *dto.EntidadeContatoRequest) error {
	// 1. Validações básicas
	if err := s.isDataValid(req); err != nil {
		return err
	}

	// 2. Validar se a entidade existe
	if err := s.validateEntidadeExists(entidadeID); err != nil {
		return err
	}

	// 3. Validar se o contato existe
	if _, err := s.contatoRepo.FindByID(entidadeID, item); err != nil {
		return fmt.Errorf("contato não encontrado: %w", err)
	}

	return nil
}

// ============================================================
// FUNÇÕES AUXILIARES
// ============================================================

// isValidTipoContato valida se o tipo de contato é válido
func isValidTipoContato(tipo int) bool {
	switch tipo {
	case dto.TipoContatoTelefone,
		dto.TipoContatoCelular,
		dto.TipoContatoEmail,
		dto.TipoContatoWhatsApp,
		dto.TipoContatoSite,
		dto.TipoContatoFacebook,
		dto.TipoContatoInstagram:
		return true
	default:
		return false
	}
}

// ============================================================
// MÉTODOS PRINCIPAIS (CRUD)
// ============================================================

// Create cria um novo contato para uma entidade
func (s *EntidadeContatoService) Create(req *dto.EntidadeContatoRequest) (*models.EntidadeContato, error) {
	// 1. Validar dados
	if err := s.isCreateValid(req); err != nil {
		return nil, err
	}

	// 2. Converter DTO para Model
	contato, err := req.ToModel()
	if err != nil {
		return nil, fmt.Errorf("erro ao converter dados: %w", err)
	}

	// 3. Salvar (o repository cuida do sequencial do Item)
	if err := s.contatoRepo.Create(contato); err != nil {
		return nil, fmt.Errorf("erro ao criar contato: %w", err)
	}

	return contato, nil
}

// GetByID busca um contato específico
func (s *EntidadeContatoService) GetByID(entidadeID, item int) (*models.EntidadeContato, error) {
	contato, err := s.contatoRepo.FindByID(entidadeID, item)
	if err != nil {
		return nil, fmt.Errorf("contato não encontrado: %w", err)
	}
	return contato, nil
}

// GetByEntidadeID busca todos os contatos de uma entidade
func (s *EntidadeContatoService) GetByEntidadeID(entidadeID int) ([]models.EntidadeContato, error) {
	// Validar se a entidade existe
	if err := s.validateEntidadeExists(entidadeID); err != nil {
		return nil, err
	}

	contatos, err := s.contatoRepo.FindByEntidadeID(entidadeID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar contatos: %w", err)
	}

	return contatos, nil
}

// GetByEntidadeIDAndTipo busca contatos de uma entidade por tipo
func (s *EntidadeContatoService) GetByEntidadeIDAndTipo(entidadeID, formaContatoID int) ([]models.EntidadeContato, error) {
	// Validar se a entidade existe
	if err := s.validateEntidadeExists(entidadeID); err != nil {
		return nil, err
	}

	// Validar tipo
	if !isValidTipoContato(formaContatoID) {
		return nil, fmt.Errorf("tipo de contato inválido: %d", formaContatoID)
	}

	contatos, err := s.contatoRepo.FindByEntidadeIDAndTipo(entidadeID, formaContatoID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar contatos: %w", err)
	}

	return contatos, nil
}

// Update atualiza um contato existente
func (s *EntidadeContatoService) Update(entidadeID, item int, req *dto.EntidadeContatoRequest) (*models.EntidadeContato, error) {
	// 1. Validar dados
	if err := s.isUpdateValid(entidadeID, item, req); err != nil {
		return nil, err
	}

	// 2. Buscar contato existente
	contato, err := s.contatoRepo.FindByID(entidadeID, item)
	if err != nil {
		return nil, fmt.Errorf("contato não encontrado: %w", err)
	}

	// 3. Atualizar campos
	contato.FormaContatoID = req.FormaContatoID
	contato.Informacao = req.Informacao
	contato.Descricao = stringPtr(req.Descricao)

	// 4. Atualizar auditoria
	if req.UpdatedBy != nil {
		contato.UpdatedBy = req.UpdatedBy
	}

	// 5. Salvar
	if err := s.contatoRepo.Update(contato); err != nil {
		return nil, fmt.Errorf("erro ao atualizar contato: %w", err)
	}

	return contato, nil
}

// Delete exclui logicamente um contato
func (s *EntidadeContatoService) Delete(entidadeID, item int) error {
	// 1. Validar se o contato existe
	contato, err := s.contatoRepo.FindByID(entidadeID, item)
	if err != nil {
		return fmt.Errorf("contato não encontrado: %w", err)
	}

	// 2. Verificar se já foi deletado
	if contato.IsDeleted() {
		return errors.New("contato já foi deletado")
	}

	// 3. TODO: Verificar se o contato está sendo usado em algum lugar

	// 4. Excluir
	if err := s.contatoRepo.Delete(entidadeID, item); err != nil {
		return fmt.Errorf("erro ao excluir contato: %w", err)
	}

	return nil
}

// List lista contatos com paginação e filtros
func (s *EntidadeContatoService) List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeContato, int64, error) {
	// Validar parâmetros de paginação
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	contatos, total, err := s.contatoRepo.List(limit, offset, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao listar contatos: %w", err)
	}

	return contatos, total, nil
}

// ============================================================
// FUNÇÕES AUXILIARES
// ============================================================

// stringPtr converte string para ponteiro
func stringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
