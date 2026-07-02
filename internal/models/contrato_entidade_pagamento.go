package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: ContratoEntidadePagamento
// ============================================================

type ContratoEntidadePagamento struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ContratoEntidadeID int      `gorm:"column:cent_id;primaryKey" json:"contrato_entidade_id"`
	Item               int      `gorm:"column:cep_item;primaryKey" json:"item"`
	TipoDocumentoID    int      `gorm:"column:tdoc_id;not null" json:"tipo_documento_id"`
	PortadorID         int      `gorm:"column:por_id;not null" json:"portador_id"`
	
	DataVencimento     time.Time `gorm:"column:cep_datavencimento;type:date;not null" json:"data_vencimento"`
	Valor              float64   `gorm:"column:cep_valor;type:decimal(15,4)" json:"valor"`
	DataDesconto       *time.Time `gorm:"column:cep_datadesconto;type:date" json:"data_desconto,omitempty"`
	ValorDesconto      *float64   `gorm:"column:cep_valordesconto;type:decimal(15,4)" json:"valor_desconto,omitempty"`

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
	ContratoEntidade *ContratoEntidade `gorm:"foreignKey:ContratoEntidadeID;references:cent_id" json:"contrato_entidade,omitempty"`
	TipoDocumento    *TipoDocumento    `gorm:"foreignKey:TipoDocumentoID;references:tdoc_id" json:"tipo_documento,omitempty"`
	Portador         *Portador         `gorm:"foreignKey:PortadorID;references:por_id" json:"portador,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (ContratoEntidadePagamento) TableName() string {
	return "contrato_entidade_pagamentos"
}

func (c *ContratoEntidadePagamento) BeforeCreate(tx *gorm.DB) error {
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

func (c *ContratoEntidadePagamento) BeforeUpdate(tx *gorm.DB) error {
	if c.UpdatedBy == nil {
		c.UpdatedBy = new(int)
		*c.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o pagamento foi deletado logicamente
func (c *ContratoEntidadePagamento) IsDeleted() bool {
	return c.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (c *ContratoEntidadePagamento) SoftDelete() {
	now := time.Now()
	c.DeletedAt = &now
}

// HasDesconto verifica se o pagamento possui desconto
func (c *ContratoEntidadePagamento) HasDesconto() bool {
	return c.ValorDesconto != nil && *c.ValorDesconto > 0
}

// GetValorLiquido retorna o valor líquido do pagamento (valor - desconto)
func (c *ContratoEntidadePagamento) GetValorLiquido() float64 {
	valor := c.Valor
	if c.HasDesconto() {
		valor -= *c.ValorDesconto
	}
	return valor
}

// IsVencido verifica se o pagamento está vencido
func (c *ContratoEntidadePagamento) IsVencido() bool {
	return time.Now().After(c.DataVencimento)
}