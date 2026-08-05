package appvalidation

import (
	"reflect"
	"regexp"
	"strings"
)

// isEmptyValue verifica se um valor está vazio
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

// isFieldMandatoryByType verifica se o campo é obrigatório baseado no tipo
func isFieldMandatoryByType(fieldValue reflect.Value) bool {
	if fieldValue.Kind() == reflect.Ptr {
		return false
	}

	if fieldValue.Kind() == reflect.Interface {
		return false
	}

	if fieldValue.Kind() == reflect.Slice ||
		fieldValue.Kind() == reflect.Map ||
		fieldValue.Kind() == reflect.Array {
		return false
	}

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

// CleanString remove espaços extras e caracteres especiais
func CleanString(value string) string {
	if value == "" {
		return ""
	}
	value = strings.TrimSpace(value)
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(value, " ")
}
