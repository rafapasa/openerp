package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/service"
)

// ============================================================
// TYPES
// ============================================================

// AuthHandler representa o handler de autenticação
type AuthHandler struct {
	authService service.AuthService
}

// ============================================================
// CONSTRUCTOR
// ============================================================

// NewAuthHandler cria uma nova instância do AuthHandler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	panic("Implementar Refresh") // TODO: Implementar
}

// ============================================================
// HANDLERS
// ============================================================

// Login realiza o login do usuário
//
//	@Summary		Realiza o login do usuário
//	@Description	Autentica o usuário e retorna tokens JWT
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest	true	"Credenciais de login"
//	@Success		200		{object}	dto.LoginResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{object}	dto.ErrorResponse
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	// 1. Validar a requisição
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Dados inválidos: " + err.Error(),
		})
		return
	}

	// 2. Validar campos obrigatórios
	if req.Login == "" || req.Senha == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_fields",
			"message": "Login e senha são obrigatórios",
		})
		return
	}
	if req.EmpresaID == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_empresa_id",
			"message": "ID da empresa é obrigatório",
		})
		return
	}

	// 3. Realizar o login
	usuario, accessToken, refreshToken, err := h.authService.Login(c,
		req.Login,
		req.Senha,
		req.EmpresaID,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid_credentials",
			"message": err.Error(),
		})
		return
	}

	// 4. Preparar a resposta
	response := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(24 * time.Hour.Seconds()),
		Usuario: dto.UsuarioInfo{
			ID:          usuario.ID,
			Nome:        usuario.Nome,
			Login:       usuario.Login,
			GrupoID:     usuario.GrupoUsuarioID,
			GrupoNome:   usuario.GrupoUsuario.Descricao,
			EmpresaID:   1,                // TODO: Buscar empresa do usuário
			EmpresaNome: "Empresa Padrão", // TODO: Buscar nome da empresa
		},
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken renova o token de acesso
//
//	@Summary		Renova o token de acesso
//	@Description	Usa o refresh token para gerar um novo access token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RefreshTokenRequest	true	"Refresh token"
//	@Success		200		{object}	dto.RefreshTokenResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{object}	dto.ErrorResponse
//	@Router			/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest

	// 1. Validar a requisição
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Dados inválidos: " + err.Error(),
		})
		return
	}

	// 2. Validar o refresh token
	if req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "missing_refresh_token",
			"message": "Refresh token é obrigatório",
		})
		return
	}

	// 3. Gerar novo access token
	accessToken, err := h.authService.RefreshToken(c, req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid_refresh_token",
			"message": err.Error(),
		})
		return
	}

	// 4. Preparar a resposta
	response := dto.RefreshTokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(24 * time.Hour.Seconds()),
	}

	c.JSON(http.StatusOK, response)
}

// Logout realiza o logout do usuário
//
//	@Summary		Realiza o logout
//	@Description	Invalida a sessão do usuário
//	@Tags			Auth
//	@Security		BearerAuth
//	@Success		200	{object}	map[string]string
//	@Router			/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// TODO: Implementar blacklist de tokens (com Redis)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout realizado com sucesso",
	})
}

// GetMe retorna as informações do usuário autenticado
//
//	@Summary		Retorna informações do usuário
//	@Description	Retorna os dados do usuário autenticado
//	@Tags			Auth
//	@Security		BearerAuth
//	@Success		200	{object}	dto.UsuarioInfo
//	@Failure		401	{object}	dto.ErrorResponse
//	@Router			/auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	// Recuperar o usuário do contexto (setado pelo middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "Usuário não autenticado",
		})
		return
	}

	// TODO: Buscar dados completos do usuário no banco
	c.JSON(http.StatusOK, gin.H{
		"id":         userID,
		"login":      c.GetString("login"),
		"grupo_id":   c.GetInt("grupo_id"),
		"empresa_id": c.GetInt("empresa_id"),
	})
}
