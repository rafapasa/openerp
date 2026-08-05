// service/usuario_service.go
package service

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/constants"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// ============================================================
// INTERFACE DO SERVICE
// ============================================================

type UsuarioService interface {
	Create(ctx context.Context, req *dto.UsuarioRequest) (*models.Usuario, error)
	GetByID(ctx context.Context, id int) (*models.Usuario, error)
	GetByLogin(ctx context.Context, login string) (*models.Usuario, error)
	List(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]models.Usuario, int64, error)
	Update(ctx context.Context, id int, req *dto.UsuarioRequest) (*models.Usuario, error)
	Delete(ctx context.Context, id int) error
	Activate(ctx context.Context, id int) error
	Deactivate(ctx context.Context, id int) error
	ChangePassword(ctx context.Context, id int, novaSenha string) error
	ValidateCredentials(ctx context.Context, login, senha string) (*models.Usuario, error)
	AddFilial(ctx context.Context, usuarioID int, empresaFilialID int) error
	RemoveFilial(ctx context.Context, usuarioID int, empresaFilialID int) error
}

// ============================================================
// CONSTANTES DE VALIDAÇÃO
// ============================================================

const (
	minLengthSenha = 6
	maxLengthSenha = 100
	maxLengthNome  = 100
	maxLengthLogin = 20
)

// ============================================================
// IMPLEMENTAÇÃO
// ============================================================

type usuarioService struct {
	usuRepo              repository.UsuarioRepository
	usuarioFilialService UsuarioFilialService
	grupoUsuarioService  GrupoUsuarioService
	empresaFilialService EmpresaFilialService
}

// ============================================================
// CONSTRUTOR
// ============================================================

func NewUsuarioService(
	repo repository.UsuarioRepository,
	usuarioFilialService UsuarioFilialService,
	grupoUsuarioService GrupoUsuarioService,
	empresaFilialService EmpresaFilialService,
) UsuarioService {
	return &usuarioService{
		usuRepo:              repo,
		usuarioFilialService: usuarioFilialService,
		grupoUsuarioService:  grupoUsuarioService,
		empresaFilialService: empresaFilialService,
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO (PRIVADOS)
// ============================================================

// validateRequest realiza as validações básicas do DTO
func (s *usuarioService) validateRequest(req *dto.UsuarioRequest) error {
	if req == nil {
		return apperrors.NewValidationError("requisição não pode ser nula")
	}

	// Validar usando o método do DTO
	if err := req.Validate(); err != nil {
		return apperrors.NewValidationError(err.Error())
	}

	// Validações adicionais
	if len(req.Nome) > maxLengthNome {
		return apperrors.NewValidationError(fmt.Sprintf("nome deve ter no máximo %d caracteres", maxLengthNome))
	}

	if len(req.Login) > maxLengthLogin {
		return apperrors.NewValidationError(fmt.Sprintf("login deve ter no máximo %d caracteres", maxLengthLogin))
	}

	if len(req.Login) < 3 {
		return apperrors.NewValidationError("login deve ter pelo menos 3 caracteres")
	}

	return nil
}

// validateBusinessRules valida as regras de negócio
func (s *usuarioService) validateBusinessRules(ctx context.Context, req *dto.UsuarioRequest, excludeID int) error {
	// 1. Validar se o grupo de usuário existe e está ativo
	grupo, err := s.grupoUsuarioService.GetByID(ctx, req.GrupoUsuarioID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewNotFoundError(fmt.Sprintf("grupo de usuário %d não encontrado", req.GrupoUsuarioID))
		}
		return apperrors.NewInternalError("erro ao verificar grupo de usuário", err)
	}

	if grupo.Situacao != constants.StatusAtivo {
		return apperrors.NewValidationError(fmt.Sprintf("grupo de usuário %d está inativo", req.GrupoUsuarioID))
	}

	// 2. Validar se o login já existe
	usuario, err := s.usuRepo.FindByLogin(ctx, req.Login)
	if err != nil {
		return err
	}
	if usuario != nil {
		return apperrors.NewConflictError(fmt.Sprintf("login %s já está em uso", req.Login))
	}

	// 3. Validar filiais
	if len(req.UsuarioFiliais) > 0 {
		for _, uf := range req.UsuarioFiliais {
			// Validar se a filial existe
			_, err := s.empresaFilialService.GetByID(ctx, uf.EmpresaFilialID)
			if err != nil {
				return err
			}

		}
	}

	// 4. Validar senha (para criação ou alteração)
	if req.Senha != nil && len(*req.Senha) > 0 {
		if len(*req.Senha) < minLengthSenha {
			return apperrors.NewValidationError(fmt.Sprintf("senha deve ter no mínimo %d caracteres", minLengthSenha))
		}
		if len(*req.Senha) > maxLengthSenha {
			return apperrors.NewValidationError(fmt.Sprintf("senha deve ter no máximo %d caracteres", maxLengthSenha))
		}
	}

	return nil
}

// validateCreateValid valida dados para criação
func (s *usuarioService) validateCreateValid(ctx context.Context, req *dto.UsuarioRequest) error {
	// 1. Validações básicas do DTO
	if err := s.validateRequest(req); err != nil {
		return err
	}

	// 2. Senha é obrigatória na criação
	if req.Senha == nil || len(*req.Senha) == 0 {
		return apperrors.NewValidationError("senha é obrigatória para criação de usuário")
	}

	// 3. Regras de negócio (excludeID = 0)
	if err := s.validateBusinessRules(ctx, req, 0); err != nil {
		return err
	}

	return nil
}

// validateUpdateValid valida dados para atualização
func (s *usuarioService) validateUpdateValid(ctx context.Context, id int, req *dto.UsuarioRequest) error {
	if id <= 0 {
		return apperrors.NewValidationError("ID do usuário inválido")
	}

	// 1. Validações básicas do DTO
	if err := s.validateRequest(req); err != nil {
		return err
	}

	// 2. Regras de negócio (excludeID = id para atualização)
	if err := s.validateBusinessRules(ctx, req, id); err != nil {
		return err
	}

	return nil
}

// ============================================================
// MÉTODOS AUXILIARES (PRIVADOS)
// ============================================================

// hashPassword gera o hash da senha usando bcrypt
func (s *usuarioService) hashPassword(senha string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		return "", apperrors.NewInternalError("erro ao gerar hash da senha", err)
	}
	return string(hash), nil
}

// comparePassword compara a senha com o hash
func (s *usuarioService) comparePassword(senha, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(senha)); err != nil {
		return apperrors.NewValidationError("senha inválida")
	}
	return nil
}

// ============================================================
// MÉTODOS PÚBLICOS (CRUD)
// ============================================================

// Create cria um novo usuário
func (s *usuarioService) Create(ctx context.Context, req *dto.UsuarioRequest) (*models.Usuario, error) {
	// Validar
	if err := s.validateCreateValid(ctx, req); err != nil {
		return nil, err
	}

	// Converter DTO para model
	usuario, err := req.ToModel()
	if err != nil {
		return nil, apperrors.NewInternalError("erro ao converter dados do usuário", err)
	}

	// Hash da senha
	hash, err := s.hashPassword(*req.Senha)
	if err != nil {
		return nil, err
	}
	usuario.Senha = hash

	// Definir situação padrão se não informada
	if usuario.Situacao == 0 {
		usuario.Situacao = constants.StatusAtivo
	}

	// Salvar usuário
	if err := s.usuRepo.Create(ctx, usuario); err != nil {
		return nil, err
	}

	// Buscar usuário completo com relacionamentos
	return s.usuRepo.FindByID(ctx, usuario.ID)
}

// GetByID busca um usuário pelo ID
func (s *usuarioService) GetByID(ctx context.Context, id int) (*models.Usuario, error) {
	if id <= 0 {
		return nil, apperrors.NewValidationError("ID do usuário inválido")
	}

	usuario, err := s.usuRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if usuario == nil {
		return nil, apperrors.NewNotFoundError(fmt.Sprintf("usuário %d não encontrado", id))
	}

	return usuario, nil
}

// GetByLogin busca um usuário pelo login
func (s *usuarioService) GetByLogin(ctx context.Context, login string) (*models.Usuario, error) {
	if login == "" {
		return nil, apperrors.NewValidationError("login não pode ser vazio")
	}

	usuario, err := s.usuRepo.FindByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	if usuario == nil {
		return nil, apperrors.NewNotFoundError(fmt.Sprintf("usuário %s não encontrado", login))
	}

	return usuario, nil
}

// List lista usuários com paginação e filtros
func (s *usuarioService) List(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]models.Usuario, int64, error) {

	usuarios, total, err := s.usuRepo.List(ctx, limit, offset, filters)
	if err != nil {
		return nil, 0, apperrors.NewInternalError("erro ao listar usuários", err)
	}

	return usuarios, total, nil
}

// Update atualiza um usuário existente
func (s *usuarioService) Update(ctx context.Context, id int, req *dto.UsuarioRequest) (*models.Usuario, error) {
	// Validar
	if err := s.validateUpdateValid(ctx, id, req); err != nil {
		return nil, err
	}

	// Buscar usuário existente
	usuario, err := s.usuRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if usuario.IsDeleted() {
		return nil, apperrors.NewValidationError("não é possível atualizar um usuário deletado")
	}

	// Converter DTO para model (atualiza apenas campos permitidos)
	updatedUsuario, err := req.ToModel()
	if err != nil {
		return nil, apperrors.NewInternalError("erro ao converter dados do usuário", err)
	}

	// Atualizar campos permitidos (mantendo ID e campos de auditoria)
	usuario.GrupoUsuarioID = updatedUsuario.GrupoUsuarioID
	usuario.Nome = updatedUsuario.Nome
	usuario.Login = updatedUsuario.Login
	usuario.Situacao = updatedUsuario.Situacao
	usuario.Observacoes = updatedUsuario.Observacoes
	usuario.AlterarColGrid = updatedUsuario.AlterarColGrid
	usuario.EmailSMTP = updatedUsuario.EmailSMTP
	usuario.PortaSMTP = updatedUsuario.PortaSMTP
	usuario.ServidorSMTP = updatedUsuario.ServidorSMTP
	usuario.UsarTLS = updatedUsuario.UsarTLS
	usuario.UsarSSL = updatedUsuario.UsarSSL
	usuario.UsuarioSMTP = updatedUsuario.UsuarioSMTP
	usuario.ExigirSenhaDV = updatedUsuario.ExigirSenhaDV

	// Se senha foi fornecida, atualizar
	if req.Senha != nil && len(*req.Senha) > 0 {
		hash, err := s.hashPassword(*req.Senha)
		if err != nil {
			return nil, err
		}
		usuario.Senha = hash
	}

	// Salvar
	if err := s.usuRepo.Update(ctx, id, usuario); err != nil {
		return nil, err
	}
	// Buscar usuário atualizado com relacionamentos
	return s.usuRepo.FindByID(ctx, usuario.ID)
}

// Delete realiza a exclusão física de um usuário
func (s *usuarioService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return apperrors.NewValidationError("ID do usuário inválido")
	}

	// Verificar se o usuário existe
	usuario, err := s.usuRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if usuario == nil {
		return apperrors.NewNotFoundError(fmt.Sprintf("usuário %d não encontrado", id))
	}

	// Verificar se pode ser deletado
	if !usuario.IsDeletavel() {
		return apperrors.NewValidationError("usuário não pode ser deletado devido a dependências")
	}

	// Remover filiais primeiro
	if err := s.usuarioFilialService.RemoveAllFiliais(ctx, id); err != nil {
		return apperrors.NewInternalError("erro ao remover filiais do usuário", err)
	}

	// Realizar exclusão
	if err := s.usuRepo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

// ============================================================
// MÉTODOS DE ATUALIZAÇÃO DE ESTADO
// ============================================================

// Activate ativa um usuário
func (s *usuarioService) Activate(ctx context.Context, id int) error {
	if id <= 0 {
		return apperrors.NewValidationError("ID do usuário inválido")
	}

	usuario, err := s.usuRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if usuario.IsDeleted() {
		return apperrors.NewValidationError("não é possível ativar um usuário deletado")
	}

	if usuario.Situacao == constants.StatusAtivo {
		return apperrors.NewValidationError("usuário já está ativo")
	}

	usuario.Situacao = constants.StatusAtivo
	return s.usuRepo.Update(ctx, id, usuario)
}

// Deactivate desativa um usuário
func (s *usuarioService) Deactivate(ctx context.Context, id int) error {
	if id <= 0 {
		return apperrors.NewValidationError("ID do usuário inválido")
	}

	usuario, err := s.usuRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if usuario.IsDeleted() {
		return apperrors.NewValidationError("não é possível desativar um usuário deletado")
	}

	if usuario.Situacao == constants.StatusInativo {
		return apperrors.NewValidationError("usuário já está inativo")
	}

	usuario.Situacao = constants.StatusInativo
	return s.usuRepo.Update(ctx, id, usuario)
}

// ============================================================
// MÉTODOS DE AUTENTICAÇÃO
// ============================================================

// ChangePassword altera a senha de um usuário
func (s *usuarioService) ChangePassword(ctx context.Context, id int, novaSenha string) error {
	if id <= 0 {
		return apperrors.NewValidationError("ID do usuário inválido")
	}

	if len(novaSenha) < minLengthSenha {
		return apperrors.NewValidationError(fmt.Sprintf("senha deve ter no mínimo %d caracteres", minLengthSenha))
	}
	if len(novaSenha) > maxLengthSenha {
		return apperrors.NewValidationError(fmt.Sprintf("senha deve ter no máximo %d caracteres", maxLengthSenha))
	}

	usuario, err := s.usuRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if !usuario.IsActive() {
		return apperrors.NewValidationError("usuário inativo não pode alterar senha")
	}

	hash, err := s.hashPassword(novaSenha)
	if err != nil {
		return err
	}
	usuario.Senha = hash

	return s.usuRepo.Update(ctx, id, usuario)
}

// ValidateCredentials valida as credenciais de um usuário
func (s *usuarioService) ValidateCredentials(ctx context.Context, login, senha string) (*models.Usuario, error) {
	if login == "" || senha == "" {
		return nil, apperrors.NewValidationError("login e senha são obrigatórios")
	}

	// Buscar usuário pelo login
	usuario, err := s.usuRepo.FindByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewValidationError("login ou senha inválidos")
		}
		return nil, apperrors.NewInternalError("erro ao buscar usuário", err)
	}

	if usuario == nil {
		return nil, apperrors.NewValidationError("login ou senha inválidos")
	}

	// Verificar se o usuário está ativo
	if !usuario.IsActive() {
		return nil, apperrors.NewValidationError("usuário inativo")
	}

	// Verificar se está deletado
	if usuario.IsDeleted() {
		return nil, apperrors.NewValidationError("usuário deletado")
	}

	// Comparar senha
	if err := s.comparePassword(senha, usuario.Senha); err != nil {
		return nil, apperrors.NewValidationError("login ou senha inválidos")
	}

	return usuario, nil
}

// ============================================================
// MÉTODOS DE FILIAIS
// ============================================================

// AddFilial adiciona uma filial a um usuário
func (s *usuarioService) AddFilial(ctx context.Context, usuarioID int, empresaFilialID int) error {
	if usuarioID <= 0 {
		return apperrors.NewValidationError("ID do usuário inválido")
	}
	if empresaFilialID <= 0 {
		return apperrors.NewValidationError("ID da empresa/filial inválido")
	}

	// Verificar se o usuário existe
	usuario, err := s.usuRepo.FindByID(ctx, usuarioID)
	if err != nil {
		return err
	}

	if !usuario.IsActive() {
		return apperrors.NewValidationError("usuário inativo não pode receber filiais")
	}

	// Verificar se a filial existe
	_, err = s.empresaFilialService.GetByID(ctx, empresaFilialID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewNotFoundError(fmt.Sprintf("empresa/filial %d não encontrada", empresaFilialID))
		}
		return err
	}
	// Verificar se já está associada
	exists, err := s.usuarioFilialService.HasFilial(ctx, usuarioID, empresaFilialID)
	if err != nil {
		return apperrors.NewInternalError("erro ao verificar associação", err)
	}
	if exists {
		return apperrors.NewValidationError("usuário já está associado a esta filial")
	}

	if err := s.usuarioFilialService.AddFilial(ctx, usuarioID, empresaFilialID); err != nil {
		return apperrors.NewInternalError("erro ao associar filial ao usuário", err)
	}

	return nil
}

// RemoveFilial remove uma filial de um usuário
func (s *usuarioService) RemoveFilial(ctx context.Context, usuarioID int, empresaFilialID int) error {
	if usuarioID <= 0 {
		return apperrors.NewValidationError("ID do usuário inválido")
	}
	if empresaFilialID <= 0 {
		return apperrors.NewValidationError("ID da empresa/filial inválido")
	}

	// Verificar se o usuário existe
	_, err := s.usuRepo.FindByID(ctx, usuarioID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewNotFoundError(fmt.Sprintf("usuário %d não encontrado", usuarioID))
		}
		return apperrors.NewInternalError("erro ao buscar usuário", err)
	}

	// Verificar se a associação existe
	exists, err := s.usuarioFilialService.HasFilial(ctx, usuarioID, empresaFilialID)
	if err != nil {
		return apperrors.NewInternalError("erro ao verificar associação", err)
	}
	if !exists {
		return apperrors.NewNotFoundError(fmt.Sprintf("associação usuário-filial %d-%d não encontrada", usuarioID, empresaFilialID))
	}

	// Remover associação
	if err := s.usuarioFilialService.RemoveFilial(ctx, usuarioID, empresaFilialID); err != nil {
		return apperrors.NewInternalError("erro ao remover filial do usuário", err)
	}

	return nil
}
