package appvalidation

// Validator é a interface principal de validação
type Validator interface {
	ValidateStruct(s interface{}) error
	ValidateField(field interface{}, tag string) error
	GetMandatoryFields(structure interface{}) []string
	ValidateMandatoryFields(structure interface{}) error
}

// DocumentValidator interface para validação de documentos
type DocumentValidator interface {
	IsValidDocumento(documento string) bool
	IsValidCPF(documento string) bool
	IsValidCNPJ(documento string) bool
	FormatarDocumento(documento string) string
	LimparDocumento(documento string) string
}
