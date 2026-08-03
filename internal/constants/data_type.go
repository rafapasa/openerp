package constants

const (
	// Data Types para configurações
	DataTypeString   = 1 // string
	DataTypeInt      = 2 // int
	DataTypeFloat    = 3 // float
	DataTypeBool     = 4 // bool
	DataTypeDateTime = 5 // datetime
	DataTypeText     = 6 // text
	DataTypeDecimal  = 7 // decimal
)

// GetDataTypeName retorna o nome do tipo de dado
func GetDataTypeName(dataType int) string {
	switch dataType {
	case DataTypeString:
		return "string"
	case DataTypeInt:
		return "integer"
	case DataTypeFloat:
		return "float"
	case DataTypeBool:
		return "boolean"
	case DataTypeDateTime:
		return "datetime"
	case DataTypeText:
		return "text"
	case DataTypeDecimal:
		return "decimal"
	default:
		return "unknown"
	}
}
