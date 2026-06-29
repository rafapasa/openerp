package utils

import (
	"time"

	"github.com/openerp/backend/internal/constants"
)

// ============================================================
// CONVERSÃO DE PONTEIROS PARA TIPOS BÁSICOS
// ============================================================

// StringValue retorna o valor de um ponteiro string ou string vazia
func StringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// StringPtr retorna um ponteiro para a string
func StringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

// IntValue retorna o valor de um ponteiro int ou 0
func IntValue(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}

// IntPtr retorna um ponteiro para o int
func IntPtr(value int) *int {
	if value == 0 {
		return nil
	}
	return &value
}

// Float64Value retorna o valor de um ponteiro float64 ou 0
func Float64Value(ptr *float64) float64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

// Float64Ptr retorna um ponteiro para o float64
func Float64Ptr(value float64) *float64 {
	if value == 0 {
		return nil
	}
	return &value
}

// BoolValue retorna o valor de um ponteiro bool ou false
func BoolValue(ptr *bool) bool {
	if ptr == nil {
		return false
	}
	return *ptr
}

// BoolPtr retorna um ponteiro para o bool
func BoolPtr(value bool) *bool {
	return &value
}

// TimeValue retorna o valor de um ponteiro time.Time ou time.Time vazio
func TimeValue(ptr *time.Time) time.Time {
	if ptr == nil {
		return time.Time{}
	}
	return *ptr
}

// TimePtr retorna um ponteiro para o time.Time
func TimePtr(value time.Time) *time.Time {
	if value.IsZero() {
		return nil
	}
	return &value
}

// ============================================================
// CONVERSÃO DE TIPOS ENUM
// ============================================================

// IntToStatus converte int para constants.Status
func IntToStatus(value int) constants.Status {
	return constants.Status(value)
}

// StatusToInt converte constants.Status para int
func StatusToInt(status constants.Status) int {
	return int(status)
}

// ============================================================
// FORMATAÇÃO DE DATAS
// ============================================================

// FormatDateTime formata time.Time para string no formato padrão
func FormatDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// FormatDate formata time.Time para string no formato de data
func FormatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

// ParseDateTime converte string para time.Time
func ParseDateTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02 15:04:05", value)
}

// ParseDate converte string para time.Time (apenas data)
func ParseDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02", value)
}

// ============================================================
// VALIDAÇÃO DE TIPOS (opcional)
// ============================================================

// IsZeroOrNil verifica se um valor é zero ou nil
func IsZeroOrNil(value interface{}) bool {
	if value == nil {
		return true
	}
	// Usar reflection para verificar valores zero
	// Por simplicidade, deixamos a cargo do chamador
	return false
}
