package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// ============================================================
// MAPPER - Conversão entre DTO e Model usando Reflection
// ============================================================

// MapToModel converte um DTO para um Model usando tags
func MapToModel(src interface{}, dst interface{}) error {
	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst)

	// Verificar se são ponteiros
	if srcVal.Kind() != reflect.Ptr || dstVal.Kind() != reflect.Ptr {
		return fmt.Errorf("src e dst devem ser ponteiros")
	}

	srcVal = srcVal.Elem()
	dstVal = dstVal.Elem()

	// Verificar se são structs
	if srcVal.Kind() != reflect.Struct || dstVal.Kind() != reflect.Struct {
		return fmt.Errorf("src e dst devem ser structs")
	}

	// Mapear campos do destino
	dstType := dstVal.Type()
	for i := 0; i < dstType.NumField(); i++ {
		dstField := dstType.Field(i)

		// Ignorar campos não exportados
		if !dstField.IsExported() {
			continue
		}

		// Ignorar campos de auditoria (já são tratados separadamente)
		if isAuditField(dstField.Name) {
			continue
		}

		// Buscar o nome do campo no DTO
		srcFieldName := getFieldName(dstField)
		if srcFieldName == "" {
			continue
		}

		// Buscar o campo no DTO
		srcField := srcVal.FieldByNameFunc(func(name string) bool {
			return strings.EqualFold(name, srcFieldName) ||
				strings.EqualFold(name, dstField.Name)
		})

		if !srcField.IsValid() {
			continue
		}

		// Verificar se o campo é preenchível
		if !srcField.CanInterface() {
			continue
		}

		// Converter e atribuir
		if err := setFieldValue(dstVal.Field(i), srcField); err != nil {
			return fmt.Errorf("erro ao converter campo %s: %w", dstField.Name, err)
		}
	}

	return nil
}

// MapToDTO converte um Model para um DTO usando tags
func MapToDTO(src interface{}, dst interface{}) error {
	return MapToModel(src, dst) // Mesma lógica, funciona nos dois sentidos
}

// ============================================================
// FUNÇÕES AUXILIARES
// ============================================================

// getFieldName obtém o nome do campo para busca no DTO
func getFieldName(field reflect.StructField) string {
	// Verificar tag json
	jsonTag := field.Tag.Get("json")
	if jsonTag != "" && jsonTag != "-" {
		parts := strings.Split(jsonTag, ",")
		return parts[0]
	}

	// Verificar tag mapstructure
	mapTag := field.Tag.Get("mapstructure")
	if mapTag != "" && mapTag != "-" {
		return mapTag
	}

	// Usar o nome do campo
	return field.Name
}

// setFieldValue atribui um valor a um campo
func setFieldValue(dstField reflect.Value, srcField reflect.Value) error {
	// Verificar se o destino pode receber valor
	if !dstField.CanSet() {
		return fmt.Errorf("campo destino não pode ser setado")
	}

	// Verificar se o destino é nulo
	if dstField.Kind() == reflect.Ptr && dstField.IsNil() {
		// Criar um novo ponteiro
		dstField.Set(reflect.New(dstField.Type().Elem()))
	}

	// Obter o valor efetivo do destino (se for ponteiro)
	dstValue := dstField
	if dstField.Kind() == reflect.Ptr {
		dstValue = dstField.Elem()
	}

	// Se o src for nil, deixar o destino como nil
	if srcField.Kind() == reflect.Ptr && srcField.IsNil() {
		if dstField.Kind() == reflect.Ptr {
			dstField.Set(reflect.Zero(dstField.Type()))
		}
		return nil
	}

	// Obter o valor efetivo do src (se for ponteiro)
	srcValue := srcField
	if srcField.Kind() == reflect.Ptr {
		srcValue = srcField.Elem()
	}

	// Verificar tipos
	if srcValue.Kind() != dstValue.Kind() {
		// Tentar converter
		converted, err := convertValue(srcValue, dstValue.Type())
		if err != nil {
			return fmt.Errorf("não foi possível converter %v para %v: %w",
				srcValue.Type(), dstValue.Type(), err)
		}
		dstValue.Set(converted)
		return nil
	}

	// Atribuir o valor
	dstValue.Set(srcValue)
	return nil
}

// convertValue converte um valor para um tipo específico
func convertValue(src reflect.Value, dstType reflect.Type) (reflect.Value, error) {
	// Converter entre tipos numéricos compatíveis
	if src.Kind() >= reflect.Int && src.Kind() <= reflect.Float64 {
		switch dstType.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return reflect.ValueOf(src.Int()).Convert(dstType), nil
		case reflect.Float32, reflect.Float64:
			return reflect.ValueOf(src.Float()).Convert(dstType), nil
		}
	}

	// Converter string para time.Time
	if src.Kind() == reflect.String && dstType == reflect.TypeOf(time.Time{}) {
		if src.String() == "" {
			return reflect.Zero(dstType), nil
		}
		t, err := time.Parse("2006-01-02", src.String())
		if err != nil {
			return reflect.Zero(dstType), err
		}
		return reflect.ValueOf(t), nil
	}

	return src, fmt.Errorf("conversão não suportada: %v para %v", src.Type(), dstType)
}

// isAuditField verifica se o campo é de auditoria
func isAuditField(fieldName string) bool {
	auditFields := []string{
		"CreatedAt", "UpdatedAt", "DeletedAt",
		"CreatedBy", "UpdatedBy",
	}
	for _, af := range auditFields {
		if fieldName == af {
			return true
		}
	}
	return false
}

// ============================================================
// MAPPER COM CONFIGURAÇÃO (mais flexível)
// ============================================================

// MapperConfig permite configurar o comportamento do mapper
type MapperConfig struct {
	IgnoreAudit  bool              // Ignorar campos de auditoria
	IgnoreFields []string          // Campos a serem ignorados
	MappingRules map[string]string // Regras de mapeamento customizadas
}

// MapToModelWithConfig converte com configuração personalizada
func MapToModelWithConfig(src interface{}, dst interface{}, config MapperConfig) error {
	// Implementação similar, mas respeitando as configurações
	// ...
	return nil
}
