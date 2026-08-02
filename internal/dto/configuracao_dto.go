package dto

// ============================================================
// DTO: ConfiguracaoRequest
// ============================================================

type ConfiguracaoRequest struct {
	EmpresaFilialID int     `json:"empresa_filial_id" binding:"required"`
	ModuloID        int     `json:"modulo_id" binding:"required"`
	Nome            string  `json:"nome" binding:"required"`
	Valor           string  `json:"valor" binding:"required"`
	DataType        int     `json:"data_type" binding:"required,min=1,max=7"`
	Descricao       *string `json:"descricao,omitempty"`
	CreatedBy       *int    `json:"created_by,omitempty"`
	UpdatedBy       *int    `json:"updated_by,omitempty"`
}

// ============================================================
// DTO: ConfiguracaoResponse
// ============================================================

type ConfiguracaoResponse struct {
	EmpresaFilialID int     `json:"empresa_filial_id"`
	ConfigID        int     `json:"config_id"`
	ModuloID        int     `json:"modulo_id"`
	Nome            string  `json:"nome"`
	Valor           string  `json:"valor"`
	DataType        int     `json:"data_type"`
	DataTypeName    string  `json:"data_type_name"`
	Descricao       *string `json:"descricao,omitempty"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}
