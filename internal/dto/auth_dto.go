package dto

// ============================================================
// REQUESTS
// ============================================================

// LoginRequest representa a requisição de login
type LoginRequest struct {
	Login     string `json:"login" binding:"required"`
	Senha     string `json:"senha" binding:"required"`
	EmpresaID *int   `json:"empresa_id,omitempty"` // Opcional para multi-empresa
}

// RefreshTokenRequest representa a requisição de refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ============================================================
// RESPONSES
// ============================================================

// LoginResponse representa a resposta de login
type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	TokenType    string      `json:"token_type"`
	ExpiresIn    int64       `json:"expires_in"` // Segundos
	Usuario      UsuarioInfo `json:"usuario"`
}

// UsuarioInfo representa as informações do usuário na resposta
type UsuarioInfo struct {
	ID          int    `json:"id"`
	Nome        string `json:"nome"`
	Login       string `json:"login"`
	GrupoNome   string `json:"grupo_nome"`
	GrupoID     int    `json:"grupo_id"`
	EmpresaID   int    `json:"empresa_id"`
	EmpresaNome string `json:"empresa_nome"`
}

// RefreshTokenResponse representa a resposta de refresh
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

// ErrorResponse representa uma resposta de erro
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}
