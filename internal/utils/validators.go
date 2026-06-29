package utils

import (
	"regexp"
	"strings"

	"github.com/paemuri/brdoc"
)

// ============================================================
// VALIDAÇÃO DE DOCUMENTOS
// ============================================================

// IsValidDocumento verifica se o documento é um CPF ou CNPJ válido
func IsValidDocumento(documento string) bool {
	return brdoc.IsCPF(documento) || brdoc.IsCNPJ(documento)
}

// IsValidCPF verifica se é um CPF válido
func IsValidCPF(documento string) bool {
	return brdoc.IsCPF(documento)
}

// IsValidCNPJ verifica se é um CNPJ válido
func IsValidCNPJ(documento string) bool {
	return brdoc.IsCNPJ(documento)
}

// LimparDocumento remove caracteres especiais do documento
// Remove pontos, barras, traços e espaços
func LimparDocumento(documento string) string {
	// Remove espaços
	documento = strings.TrimSpace(documento)
	// Remove caracteres não numéricos
	re := regexp.MustCompile(`[^0-9]`)
	return re.ReplaceAllString(documento, "")
}

// FormatarDocumento formata CPF ou CNPJ automaticamente
func FormatarDocumento(documento string) string {
	documento = LimparDocumento(documento)

	if len(documento) == 11 {
		return formatCPF(documento)
	}
	if len(documento) == 14 {
		return formatCNPJ(documento)
	}
	return documento
}

// formatCPF formata um CPF válido
func formatCPF(cpf string) string {
	if len(cpf) != 11 {
		return cpf
	}
	return cpf[:3] + "." + cpf[3:6] + "." + cpf[6:9] + "-" + cpf[9:]
}

// formatCNPJ formata um CNPJ válido
func formatCNPJ(cnpj string) string {
	if len(cnpj) != 14 {
		return cnpj
	}
	return cnpj[:2] + "." + cnpj[2:5] + "." + cnpj[5:8] + "/" + cnpj[8:12] + "-" + cnpj[12:]
}

// ============================================================
// VALIDAÇÃO DE TAMANHO
// ============================================================

// IsCPFLength verifica se o CPF tem 11 dígitos
func IsCPFLength(documento string) bool {
	return len(LimparDocumento(documento)) == 11
}

// IsCNPJLength verifica se o CNPJ tem 14 dígitos
func IsCNPJLength(documento string) bool {
	return len(LimparDocumento(documento)) == 14
}

// ============================================================
// LIMPEZA DE STRINGS
// ============================================================

// CleanString remove espaços extras e caracteres especiais
func CleanString(value string) string {
	if value == "" {
		return ""
	}
	// Remove espaços extras
	value = strings.TrimSpace(value)
	// Remove múltiplos espaços
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(value, " ")
}
