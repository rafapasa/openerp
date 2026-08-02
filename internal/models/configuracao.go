package models

import (
	"errors"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/constants"
)

// ============================================================
// MODEL: Configuracao
// ============================================================

const (
	ConfigCachePrefix = "config:"
	ConfigCacheTTL    = 24 * time.Hour
	ConfigCacheLimit  = 1000
)

type Configuracao struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	EmpresaFilialID int     `gorm:"column:emf_id;primaryKey" json:"empresa_filial_id"`
	ConfigID        int     `gorm:"column:config_id;primaryKey" json:"config_id"`
	ModuloID        int     `gorm:"column:mod_id;not null" json:"modulo_id"`
	Nome            string  `gorm:"column:config_nome;type:varchar(255);not null" json:"nome"`
	Valor           string  `gorm:"column:config_valor;type:varchar(255);not null" json:"valor"`
	DataType        int     `gorm:"column:config_datatype;not null" json:"data_type"`
	Descricao       *string `gorm:"column:config_descricao;type:text" json:"descricao,omitempty"`

	// ============================================================
	// CAMPOS DE AUDITORIA
	// ============================================================
	CreatedAt time.Time  `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime;index" json:"deleted_at,omitempty"`
	CreatedBy *int       `gorm:"column:created_by" json:"created_by,omitempty"`
	UpdatedBy *int       `gorm:"column:updated_by" json:"updated_by,omitempty"`

	// ============================================================
	// RELACIONAMENTOS
	// ============================================================
	EmpresaFilial *EmpresaFilial `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
	// Modulo        *Modulo        `gorm:"foreignKey:ModuloID;references:mod_id" json:"modulo,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (Configuracao) TableName() string {
	return "configuracao"
}

func (c *Configuracao) BeforeCreate(tx *gorm.DB) error {
	if c.CreatedBy == nil {
		c.CreatedBy = new(int)
		*c.CreatedBy = 0
	}
	if c.UpdatedBy == nil {
		c.UpdatedBy = new(int)
		*c.UpdatedBy = 0
	}
	return nil
}

func (c *Configuracao) BeforeUpdate(tx *gorm.DB) error {
	if c.UpdatedBy == nil {
		c.UpdatedBy = new(int)
		*c.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se a configuração foi deletada logicamente
func (c *Configuracao) IsDeleted() bool {
	return c.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (c *Configuracao) SoftDelete() {
	now := time.Now()
	c.DeletedAt = &now
}

// HasDescricao verifica se possui descrição
func (c *Configuracao) HasDescricao() bool {
	return c.Descricao != nil && *c.Descricao != ""
}

// GetDescricaoOrDefault retorna a descrição ou string vazia
func (c *Configuracao) GetDescricaoOrDefault() string {
	if c.HasDescricao() {
		return *c.Descricao
	}
	return ""
}

// GetValorAsInt retorna o valor como inteiro
func (c *Configuracao) GetValorAsInt() (int, error) {
	if c.DataType != constants.DataTypeInt {
		return 0, errors.New("configuração não é do tipo inteiro")
	}
	return strconv.Atoi(c.Valor)
}

// GetValorAsFloat retorna o valor como float64
func (c *Configuracao) GetValorAsFloat() (float64, error) {
	if c.DataType != constants.DataTypeFloat && c.DataType != constants.DataTypeDecimal {
		return 0, errors.New("configuração não é do tipo numérico")
	}
	return strconv.ParseFloat(c.Valor, 64)
}

// GetValorAsBool retorna o valor como booleano
func (c *Configuracao) GetValorAsBool() (bool, error) {
	if c.DataType != constants.DataTypeBool {
		return false, errors.New("configuração não é do tipo booleano")
	}
	return strconv.ParseBool(c.Valor)
}

// GetValorAsString retorna o valor como string
func (c *Configuracao) GetValorAsString() string {
	return c.Valor
}

// GetValorAsDateTime retorna o valor como time.Time
func (c *Configuracao) GetValorAsDateTime() (time.Time, error) {
	if c.DataType != constants.DataTypeDateTime {
		return time.Time{}, errors.New("configuração não é do tipo data/hora")
	}
	return time.Parse("2006-01-02 15:04:05", c.Valor)
}

// IsNumeric verifica se a configuração é numérica
func (c *Configuracao) IsNumeric() bool {
	return c.DataType == constants.DataTypeInt ||
		c.DataType == constants.DataTypeFloat ||
		c.DataType == constants.DataTypeDecimal
}

// IsBoolean verifica se a configuração é booleana
func (c *Configuracao) IsBoolean() bool {
	return c.DataType == constants.DataTypeBool
}

// IsString verifica se a configuração é string
func (c *Configuracao) IsString() bool {
	return c.DataType == constants.DataTypeString ||
		c.DataType == constants.DataTypeText
}

func (c *Configuracao) Key() string {
	return ConfigCachePrefix + strconv.Itoa(c.EmpresaFilialID) + ":" + strconv.Itoa(c.ConfigID)
}
