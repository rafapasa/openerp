// internal/appcontext/gorm_context.go
package appcontext

import (
	"context"
)

// Chaves para o contexto do GORM
type gormContextKey string

const (
	UserIDKey   gormContextKey = "gorm_user_id"
	UserEmailKey gormContextKey = "gorm_user_email"
	TraceIDKey   gormContextKey = "gorm_trace_id"
)

// WithUserID adiciona o userID ao contexto
func WithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserID do contexto
func GetUserID(ctx context.Context) int {
	if userID, ok := ctx.Value(UserIDKey).(int); ok {
		return userID
	}
	return 0
}

// WithUserEmail adiciona o email ao contexto
func WithUserEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, UserEmailKey, email)
}

// GetUserEmail do contexto
func GetUserEmail(ctx context.Context) string {
	if email, ok := ctx.Value(UserEmailKey).(string); ok {
		return email
	}
	return ""
}