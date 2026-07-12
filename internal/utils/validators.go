package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/paemuri/brdoc"

	apperrors "github.com/openerp/backend/internal/erros"
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

func ValidateMandatoryFields(structure any) error {
	mandatoryFields := GetMandatoryFields(structure)
	v := reflect.ValueOf(structure)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for _, fieldName := range mandatoryFields {
		field := v.FieldByName(fieldName)
		if !field.IsValid() {
			return apperrors.NewValidationError(fmt.Sprintf("Campo obrigatório %s não encontrado na estrutura", fieldName))
		}
		if isEmptyValue(field) {
			return apperrors.NewValidationError(fmt.Sprintf("Campo obrigatório %s está vazio", fieldName))
		}
	}
	return nil
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Ptr:
		return v.IsNil()
	default:
		return false
	}
}

// GetMandatoryFields retorna uma lista de campos obrigatórios baseado no tipo (não ponteiro)
// e também respeita a tag `mandatory:"true"` ou `binding:"required"`
func GetMandatoryFields(structure any) []string {
	var fields []string
	v := reflect.ValueOf(structure)

	// Se for ponteiro, obtém o valor apontado
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Se não for struct, retorna vazio
	if v.Kind() != reflect.Struct {
		return fields
	}

	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Verifica se o campo é obrigatório
		if isFieldMandatory(field, fieldValue) {
			fields = append(fields, field.Name)
		}
	}

	return fields
}

func isFieldMandatory(field reflect.StructField, fieldValue reflect.Value) bool {
	// 1. Verifica tag `mandatory:"true"`
	if tag := field.Tag.Get("mandatory"); tag == "true" {
		return true
	}

	// 2. Verifica tag `binding:"required"`
	if tag := field.Tag.Get("binding"); strings.Contains(tag, "required") {
		return true
	}

	// 3. Verifica tag `validate:"required"`
	if tag := field.Tag.Get("validate"); strings.Contains(tag, "required") {
		return true
	}

	// 4. Verifica pelo tipo (não ponteiro)
	return isFieldMandatoryByType(fieldValue)
}

// isFieldMandatoryByType verifica se o campo é obrigatório baseado no tipo
func isFieldMandatoryByType(fieldValue reflect.Value) bool {
	// Campos que são ponteiros são opcionais
	if fieldValue.Kind() == reflect.Ptr {
		return false
	}

	// Campos que são interfaces são opcionais
	if fieldValue.Kind() == reflect.Interface {
		return false
	}

	// Campos que são slices, maps, arrays são opcionais (podem ser vazios)
	if fieldValue.Kind() == reflect.Slice ||
		fieldValue.Kind() == reflect.Map ||
		fieldValue.Kind() == reflect.Array {
		return false
	}

	// Tipos básicos são obrigatórios (não ponteiros)
	switch fieldValue.Kind() {
	case reflect.String,
		reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Struct:
		return true
	default:
		return false
	}
}
