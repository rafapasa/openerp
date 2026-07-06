package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: ContratoEntidadeProduto
// ============================================================

type ContratoEntidadeProduto struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ContratoEntidadeID int `gorm:"column:cent_id;primaryKey" json:"contrato_entidade_id"`
	Item               int `gorm:"column:cepro_item;primaryKey" json:"item"`
	ProdutoID          int `gorm:"column:pro_id;not null" json:"produto_id"`

	ValorUnitario float64 `gorm:"column:cepro_valorunitario;type:decimal(15,4);not null" json:"valor_unitario"`
	Quantidade    float64 `gorm:"column:cepro_quantidade;type:decimal(15,4);not null" json:"quantidade"`
	ValorTotal    float64 `gorm:"column:cepro_valortotal;type:decimal(15,4);not null" json:"valor_total"`

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
	Produto          *Produto          `gorm:"foreignKey:ProdutoID;references:pro_id" json:"produto,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (ContratoEntidadeProduto) TableName() string {
	return "contrato_entidade_produto"
}

func (c *ContratoEntidadeProduto) BeforeCreate(tx *gorm.DB) error {
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

func (c *ContratoEntidadeProduto) BeforeUpdate(tx *gorm.DB) error {
	if c.UpdatedBy == nil {
		c.UpdatedBy = new(int)
		*c.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o produto do contrato foi deletado logicamente
func (c *ContratoEntidadeProduto) IsDeleted() bool {
	return c.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (c *ContratoEntidadeProduto) SoftDelete() {
	now := time.Now()
	c.DeletedAt = &now
}

// CalcularValorTotal recalcula o valor total do item
func (c *ContratoEntidadeProduto) CalcularValorTotal() float64 {
	return c.ValorUnitario * c.Quantidade
}
