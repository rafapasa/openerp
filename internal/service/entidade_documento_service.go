package service

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/openerp/backend/internal/dto"
	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/utils"
	"gorm.io/gorm"
)

// ============================================================
// TYPES
// ============================================================

// EntidadeDocumentoService é o serviço para documentos de entidade
type EntidadeDocumentoService struct {
	documentoRepo *repository.EntidadeDocumentoRepository
	entidadeRepo  *repository.EntidadeRepository
}

// ============================================================
// CONSTRUCTOR
// ============================================================

// NewEntidadeDocumentoService cria uma nova instância
func NewEntidadeDocumentoService(db *gorm.DB) *EntidadeDocumentoService {
	return &EntidadeDocumentoService{
		documentoRepo: repository.NewEntidadeDocumentoRepository(db),
		entidadeRepo:  repository.NewEntidadeRepository(db),
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO (PRIVADOS)
// ============================================================

// isDataValid realiza as validações básicas de um documento
func (s *EntidadeDocumentoService) isDataValid(req *dto.EntidadeDocumentoRequest) error {
	// 1. Validar campos obrigatórios
	if err := s.validateRequiredFields(req); err != nil {
		return err
	}

	// 2. Validar arquivo (tamanho)
	if err := s.validateFile(req); err != nil {
		return err
	}

	return nil
}

// validateRequiredFields valida campos obrigatórios
func (s *EntidadeDocumentoService) validateRequiredFields(req *dto.EntidadeDocumentoRequest) error {
	if strings.TrimSpace(req.Arquivo) == "" {
		return apperrors.NewValidationError("arquivo é obrigatório")

	}

	if strings.TrimSpace(req.DataInclusao) == "" {
		return apperrors.NewValidationError("data de inclusão é obrigatória")
	}

	return nil
}

// validateFile valida o arquivo (tamanho)
func (s *EntidadeDocumentoService) validateFile(req *dto.EntidadeDocumentoRequest) error {
	// Validar tamanho máximo (ex: 10MB)
	// O tamanho em Base64 é ~33% maior que o original
	// 10MB * 1.33 ≈ 13.3MB em Base64
	const maxSize = 10 * 1024 * 1024 // 10MB

	if len(req.Arquivo) > maxSize {
		return apperrors.NewValidationError("arquivo muito grande, máximo permitido: 10MB")
	}

	return nil
}

// validateEntidadeExists verifica se a entidade existe
func (s *EntidadeDocumentoService) validateEntidadeExists(entidadeID int) error {
	_, err := s.entidadeRepo.FindByID(entidadeID)
	if err != nil {
		return fmt.Errorf("entidade não encontrada: %w", err)
	}
	return nil
}

// isCreateValid valida dados para criação
func (s *EntidadeDocumentoService) isCreateValid(req *dto.EntidadeDocumentoRequest) error {
	// 1. Validações básicas
	if err := s.isDataValid(req); err != nil {
		return err
	}

	// 2. Validar se a entidade existe
	if err := s.validateEntidadeExists(req.EntidadeID); err != nil {
		return err
	}

	return nil
}

// isUpdateValid valida dados para atualização
func (s *EntidadeDocumentoService) isUpdateValid(entidadeID, item int, req *dto.EntidadeDocumentoRequest) error {
	// 1. Validações básicas
	if err := s.isDataValid(req); err != nil {
		return err
	}

	// 2. Validar se a entidade existe
	if err := s.validateEntidadeExists(entidadeID); err != nil {
		return err
	}

	// 3. Validar se o documento existe
	if _, err := s.documentoRepo.FindByID(entidadeID, item); err != nil {
		return apperrors.NewValidationError(fmt.Sprintf("documento não encontrado: %w", err))
	}

	return nil
}

// ============================================================
// MÉTODOS PRINCIPAIS (CRUD)
// ============================================================

// Create cria um novo documento para uma entidade
func (s *EntidadeDocumentoService) Create(req *dto.EntidadeDocumentoRequest) (*models.EntidadeDocumento, error) {
	// 1. Validar dados
	if err := s.isCreateValid(req); err != nil {
		return nil, err
	}

	// 2. Converter DTO para Model
	documento, err := req.ToModel()
	if err != nil {
		return nil, apperrors.NewValidationError(fmt.Sprintf("erro ao converter dados: %w", err))
	}

	// 3. Salvar (o repository cuida do sequencial do Item)
	if err := s.documentoRepo.Create(documento); err != nil {
		return nil, apperrors.NewValidationError(fmt.Sprintf("erro ao criar documento: %w", err))
	}

	return documento, nil
}

// GetByID busca um documento específico
func (s *EntidadeDocumentoService) GetByID(entidadeID, item int) (*models.EntidadeDocumento, error) {
	documento, err := s.documentoRepo.FindByID(entidadeID, item)
	if err != nil {
		return nil, apperrors.NewValidationError(fmt.Sprintf("documento não encontrado: %w", err))
	}
	return documento, nil
}

// GetByEntidadeID busca todos os documentos de uma entidade
func (s *EntidadeDocumentoService) GetByEntidadeID(entidadeID int) ([]models.EntidadeDocumento, error) {
	// Validar se a entidade existe
	if err := s.validateEntidadeExists(entidadeID); err != nil {
		return nil, err
	}

	documentos, err := s.documentoRepo.FindByEntidadeID(entidadeID)
	if err != nil {
		return nil, apperrors.NewValidationError(fmt.Sprintf("erro ao buscar documentos: %w", err))
	}

	return documentos, nil
}

// Update atualiza um documento existente
func (s *EntidadeDocumentoService) Update(entidadeID, item int, req *dto.EntidadeDocumentoRequest) (*models.EntidadeDocumento, error) {
	// 1. Validar dados
	if err := s.isUpdateValid(entidadeID, item, req); err != nil {
		return nil, err
	}

	// 2. Buscar documento existente
	documento, err := s.documentoRepo.FindByID(entidadeID, item)
	if err != nil {
		return nil, apperrors.NewValidationError(fmt.Sprintf("documento não encontrado: %w", err))
	}

	// 3. Atualizar campos
	documento.Descricao = stringPtr(req.Descricao)
	documento.Tipo = stringPtr(req.Tipo)

	// 4. Atualizar arquivo (se fornecido)
	if req.Arquivo != "" {
		if err := s.validateFile(req); err != nil {
			return nil, err
		}
		// Converter Base64 para []byte
		if data, err := base64Decode(req.Arquivo); err == nil {
			documento.Arquivo = data
		} else {
			return nil, apperrors.NewValidationError(fmt.Sprintf("erro ao decodificar arquivo: %w", err))
		}
	}

	// 5. Atualizar data de inclusão (se fornecida)
	if req.DataInclusao != "" {
		if data, err := utils.ParseDate(req.DataInclusao); err == nil {
			documento.DataInclusao = data
		}
	}

	// 6. Atualizar auditoria
	if req.UpdatedBy != nil {
		documento.UpdatedBy = req.UpdatedBy
	}

	// 7. Salvar
	if err := s.documentoRepo.Update(documento); err != nil {
		return nil, apperrors.NewValidationError(fmt.Sprintf("erro ao atualizar documento: %w", err))
	}

	return documento, nil
}

// Delete exclui logicamente um documento
func (s *EntidadeDocumentoService) Delete(entidadeID, item int) error {
	// 1. Validar se o documento existe
	documento, err := s.documentoRepo.FindByID(entidadeID, item)
	if err != nil {
		return apperrors.NewValidationError(fmt.Sprintf("documento não encontrado: %w", err))
	}

	// 2. Verificar se já foi deletado
	if documento.IsDeleted() {
		return apperrors.NewConflictError("documento já foi deletado")
	}

	// 3. TODO: Verificar se o documento está sendo usado em algum lugar

	// 4. Excluir
	if err := s.documentoRepo.Delete(entidadeID, item); err != nil {
		return apperrors.NewValidationError(fmt.Sprintf("erro ao excluir documento: %w", err))
	}

	return nil
}

// List lista documentos com paginação e filtros
func (s *EntidadeDocumentoService) List(limit, offset int, filters map[string]interface{}) ([]models.EntidadeDocumento, int64, error) {
	// Validar parâmetros de paginação
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	documentos, total, err := s.documentoRepo.List(limit, offset, filters)
	if err != nil {
		return nil, 0, apperrors.NewValidationError(fmt.Sprintf("erro ao listar documentos: %w", err))
	}

	return documentos, total, nil
}

// ============================================================
// FUNÇÕES AUXILIARES
// ============================================================

// base64Decode decodifica Base64 para []byte
func base64Decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}
