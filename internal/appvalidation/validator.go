package appvalidation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/paemuri/brdoc"
)

// playValidator implementa a interface Validator usando validator/v10
type playValidator struct {
	validate *validator.Validate
}

// NewPlayValidator cria uma nova instância do validador com validator/v10
// Esta função será usada pelo Wire
func NewPlayValidator() Validator {
	v := validator.New()
	
	// Registra validações customizadas
	v.RegisterValidation("cpf", validateCPF)
	v.RegisterValidation("cnpj", validateCNPJ)
	v.RegisterValidation("documento", validateDocumento)
	
	return &playValidator{
		validate: v,
	}
}

// ValidateStruct implementa o método da interface Validator
func (v *playValidator) ValidateStruct(s interface{}) error {
	err := v.validate.Struct(s)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return v.formatValidationErrors(validationErrors)
		}
		return fmt.Errorf("erro de validação: %w", err)
	}
	return nil
}

// ValidateField valida um campo específico
func (v *playValidator) ValidateField(field interface{}, tag string) error {
	err := v.validate.Var(field, tag)
	if err != nil {
		return fmt.Errorf("validação do campo falhou: %w", err)
	}
	return nil
}

// GetMandatoryFields retorna campos obrigatórios baseado nas tags
func (v *playValidator) GetMandatoryFields(structure interface{}) []string {
	var fields []string
	val := reflect.ValueOf(structure)
	
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	
	if val.Kind() != reflect.Struct {
		return fields
	}
	
	typ := val.Type()
	
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)
		
		if v.isFieldMandatory(field, fieldValue) {
			fields = append(fields, field.Name)
		}
	}
	
	return fields
}

// ValidateMandatoryFields valida campos obrigatórios
func (v *playValidator) ValidateMandatoryFields(structure interface{}) error {
	mandatoryFields := v.GetMandatoryFields(structure)
	val := reflect.ValueOf(structure)
	
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	
	for _, fieldName := range mandatoryFields {
		field := val.FieldByName(fieldName)
		if !field.IsValid() {
			return fmt.Errorf("campo obrigatório %s não encontrado na estrutura", fieldName)
		}
		if isEmptyValue(field) {
			return fmt.Errorf("campo obrigatório %s está vazio", fieldName)
		}
	}
	return nil
}

// isFieldMandatory verifica se um campo é obrigatório
func (v *playValidator) isFieldMandatory(field reflect.StructField, fieldValue reflect.Value) bool {
	// Usa as tags do validator/v10 como fonte primária
	tag := field.Tag.Get("validate")
	if strings.Contains(tag, "required") {
		return true
	}
	
	// Compatibilidade com tags existentes
	if tag := field.Tag.Get("mandatory"); tag == "true" {
		return true
	}
	
	if tag := field.Tag.Get("binding"); strings.Contains(tag, "required") {
		return true
	}
	
	// Verifica pelo tipo (não ponteiro)
	return isFieldMandatoryByType(fieldValue)
}

// formatValidationErrors formata os erros de validação
func (v *playValidator) formatValidationErrors(errors validator.ValidationErrors) error {
	var errMsgs []string
	for _, err := range errors {
		switch err.Tag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("campo '%s' é obrigatório", err.Field()))
		case "email":
			errMsgs = append(errMsgs, fmt.Sprintf("campo '%s' deve ser um email válido", err.Field()))
		case "min":
			errMsgs = append(errMsgs, fmt.Sprintf("campo '%s' deve ter no mínimo %s caracteres", err.Field(), err.Param()))
		case "max":
			errMsgs = append(errMsgs, fmt.Sprintf("campo '%s' deve ter no máximo %s caracteres", err.Field(), err.Param()))
		case "cpf":
			errMsgs = append(errMsgs, fmt.Sprintf("campo '%s' deve ser um CPF válido", err.Field()))
		case "cnpj":
			errMsgs = append(errMsgs, fmt.Sprintf("campo '%s' deve ser um CNPJ válido", err.Field()))
		case "documento":
			errMsgs = append(errMsgs, fmt.Sprintf("campo '%s' deve ser um CPF ou CNPJ válido", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("campo '%s' falhou na validação '%s'", err.Field(), err.Tag()))
		}
	}
	return fmt.Errorf("erros de validação: %s", strings.Join(errMsgs, "; "))
}

// Validações customizadas
func validateCPF(fl validator.FieldLevel) bool {
	return brdoc.IsCPF(fl.Field().String())
}

func validateCNPJ(fl validator.FieldLevel) bool {
	return brdoc.IsCNPJ(fl.Field().String())
}

func validateDocumento(fl validator.FieldLevel) bool {
	doc := fl.Field().String()
	return brdoc.IsCPF(doc) || brdoc.IsCNPJ(doc)
}