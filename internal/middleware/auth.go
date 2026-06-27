package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// MIDDLEWARE
// ============================================================

// AuthMiddleware é o middleware de autenticação
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Obter o token do header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "token_required",
				"message": "Token de autenticação é obrigatório",
			})
			return
		}

		// 2. Verificar o formato do token (Bearer)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_token_format",
				"message": "Formato do token inválido. Use: Bearer <token>",
			})
			return
		}

		tokenString := parts[1]

		// 3. Validar o token
		claims, err := utils.ValidateAccessToken(tokenString, jwtSecret)
		if err != nil {
			errorMessage := "Token inválido"
			errorCode := "invalid_token"

			if err == utils.ErrExpiredToken {
				errorMessage = "Token expirado"
				errorCode = "token_expired"
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   errorCode,
				"message": errorMessage,
			})
			return
		}

		// 4. Adicionar as claims ao contexto da requisição
		c.Set("user_id", claims.UserID)
		c.Set("login", claims.Login)
		c.Set("grupo_id", claims.GrupoID)
		c.Set("empresa_id", claims.EmpresaID)
		c.Set("claims", claims)

		// 5. Continuar com a requisição
		c.Next()
	}
}

// ============================================================
// HELPERS (para usar nos handlers)
// ============================================================

// GetUserID retorna o ID do usuário do contexto
func GetUserID(c *gin.Context) int {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(int)
}

// GetLogin retorna o login do usuário do contexto
func GetLogin(c *gin.Context) string {
	login, exists := c.Get("login")
	if !exists {
		return ""
	}
	return login.(string)
}

// GetGrupoID retorna o ID do grupo do usuário do contexto
func GetGrupoID(c *gin.Context) int {
	grupoID, exists := c.Get("grupo_id")
	if !exists {
		return 0
	}
	return grupoID.(int)
}

// GetEmpresaID retorna o ID da empresa do usuário do contexto
func GetEmpresaID(c *gin.Context) int {
	empresaID, exists := c.Get("empresa_id")
	if !exists {
		return 0
	}
	return empresaID.(int)
}

// GetClaims retorna as claims completas do contexto
func GetClaims(c *gin.Context) *utils.JWTClaims {
	claims, exists := c.Get("claims")
	if !exists {
		return nil
	}
	return claims.(*utils.JWTClaims)
}
