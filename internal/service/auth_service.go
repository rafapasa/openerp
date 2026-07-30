package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/config"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/utils"
)

// authServiceInterface define o contrato para as operações de autenticação.
type AuthService interface {
	Login(login, senha string, empresaID *int) (*models.Usuario, string, string, error)
	RefreshToken(refreshToken string) (string, error)
	ValidateToken(token string) (*utils.JWTClaims, error)
}

// authService implementa authServiceInterface para gerenciar a autenticação de usuários.

type authService struct {
	userRepo   repository.UsuarioRepository
	jwtSecret  string
	expiresIn  time.Duration
	refreshExp time.Duration
}

func NewauthService(db *gorm.DB, cfg *config.Config) AuthService {
	return &authService{
		userRepo:   repository.NewUsuarioRepository(db),
		jwtSecret:  cfg.JWTSecret,
		expiresIn:  cfg.JWTExpiresIn,
		refreshExp: cfg.JWTRefreshExpiresIn,
	}
}

// Login realiza o login do usuário
func (s *authService) Login(login, senha string, empresaID *int) (*models.Usuario, string, string, error) {
	// 1. Buscar o usuário pelo login (COM GRUPO)
	usuario, err := s.userRepo.FindByLoginWithGrupo(login)
	if err != nil {
		return nil, "", "", errors.New("usuário ou senha inválidos")
	}

	// 2. Verificar se o usuário está ativo
	if !usuario.IsActive() {
		return nil, "", "", errors.New("usuário inativo")
	}

	// 3. Verificar se foi deletado
	if usuario.IsDeleted() {
		return nil, "", "", errors.New("usuário deletado")
	}

	// 4. Verificar a senha (MD5)
	hashedSenha := hashPassword(senha)
	if usuario.Senha != hashedSenha {
		return nil, "", "", errors.New("usuário ou senha inválidos")
	}

	// 5. Definir empresaID
	empresaIDFinal := s.getEmpresaID(usuario.ID, empresaID)

	// 6. Gerar tokens
	accessToken, err := utils.GenerateToken(
		usuario.ID,
		usuario.Login,
		usuario.GrupoUsuarioID,
		empresaIDFinal,
		utils.AccessToken,
		s.jwtSecret,
		s.expiresIn,
	)
	if err != nil {
		return nil, "", "", errors.New("erro ao gerar token de acesso")
	}

	refreshToken, err := utils.GenerateToken(
		usuario.ID,
		usuario.Login,
		usuario.GrupoUsuarioID,
		empresaIDFinal,
		utils.RefreshToken,
		s.jwtSecret,
		s.refreshExp,
	)
	if err != nil {
		return nil, "", "", errors.New("erro ao gerar token de refresh")
	}

	return usuario, accessToken, refreshToken, nil
}

// RefreshToken renova o token de acesso
func (s *authService) RefreshToken(refreshToken string) (string, error) {
	claims, err := utils.ValidateRefreshToken(refreshToken, s.jwtSecret)
	if err != nil {
		return "", errors.New("refresh token inválido")
	}

	usuario, err := s.userRepo.FindByIDWithGrupo(claims.UserID)
	if err != nil {
		return "", errors.New("usuário não encontrado")
	}

	if !usuario.IsActive() {
		return "", errors.New("usuário inativo")
	}

	accessToken, err := utils.GenerateToken(
		usuario.ID,
		usuario.Login,
		usuario.GrupoUsuarioID,
		claims.EmpresaID,
		utils.AccessToken,
		s.jwtSecret,
		s.expiresIn,
	)
	if err != nil {
		return "", errors.New("erro ao gerar novo token de acesso")
	}

	return accessToken, nil
}

func (s *authService) ValidateToken(token string) (*utils.JWTClaims, error) {
	return utils.ValidateToken(token, s.jwtSecret)
}

// getEmpresaID retorna o ID da empresa a ser usada no token
func (s *authService) getEmpresaID(usuarioID int, empresaID *int) int {
	if empresaID != nil && *empresaID > 0 {
		return *empresaID
	}

	filiais, err := s.userRepo.FindUsuarioFiliais(usuarioID)
	if err == nil && len(filiais) > 0 {
		return filiais[0].EmpresaFilialID
	}

	return 1
}

func hashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}
