package utils

import (
	"github.com/openerp/backend/internal/appvalidation"
)

// Instâncias globais (ou use injeção de dependência)
var (
	validator    = appvalidation.NewPlayValidator()
	docValidator = appvalidation.NewDocumentValidator()
)

// ============================================================
// FUNÇÕES DE COMPATIBILIDADE (DELEGAM PARA A NOVA CAMADA)
// ============================================================

// IsValidDocumento verifica se o documento é um CPF ou CNPJ válido
func IsValidDocumento(documento string) bool {
	return docValidator.IsValidDocumento(documento)
}

// IsValidCPF verifica se é um CPF válido
func IsValidCPF(documento string) bool {
	return docValidator.IsValidCPF(documento)
}

// IsValidCNPJ verifica se é um CNPJ válido
func IsValidCNPJ(documento string) bool {
	return docValidator.IsValidCNPJ(documento)
}

// LimparDocumento remove caracteres especiais do documento
func LimparDocumento(documento string) string {
	return docValidator.LimparDocumento(documento)
}

// FormatarDocumento formata CPF ou CNPJ automaticamente
func FormatarDocumento(documento string) string {
	return docValidator.FormatarDocumento(documento)
}

// ValidateMandatoryFields valida campos obrigatórios
func ValidateMandatoryFields(structure any) error {
	return validator.ValidateMandatoryFields(structure)
}

// GetMandatoryFields retorna campos obrigatórios
func GetMandatoryFields(structure any) []string {
	return validator.GetMandatoryFields(structure)
}

// CleanString remove espaços extras e caracteres especiais
func CleanString(value string) string {
	return appvalidation.CleanString(value)
}

// ============================================================
// FUNÇÕES AUXILIARES ADICIONAIS
// ============================================================

// IsCPFLength verifica se o CPF tem 11 dígitos
func IsCPFLength(documento string) bool {
	return len(LimparDocumento(documento)) == 11
}

// IsCNPJLength verifica se o CNPJ tem 14 dígitos
func IsCNPJLength(documento string) bool {
	return len(LimparDocumento(documento)) == 14
}
