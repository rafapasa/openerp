package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: ContratoEntidadeTipo
// ============================================================

type ContratoEntidadeTipo struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID        int    `gorm:"column:cet_id;primaryKey;autoIncrement" json:"id"`
	Nome      string `gorm:"column:cet_nome;type:varchar(255);not null" json:"nome"`
	ReceitaID *int   `gorm:"column:rec_id" json:"receita_id,omitempty"`
	DespesaID *int   `gorm:"column:desp_id" json:"despesa_id,omitempty"`

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
	Receita   *Receita           `gorm:"foreignKey:ReceitaID;references:rec_id" json:"receita,omitempty"`
	Despesa   *Despesa           `gorm:"foreignKey:DespesaID;references:desp_id" json:"despesa,omitempty"`
	Contratos []ContratoEntidade `gorm:"foreignKey:ContratoEntidadeTipoID;references:cet_id" json:"contratos,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (ContratoEntidadeTipo) TableName() string {
	return "contrato_entidade_tipo"
}

func (c *ContratoEntidadeTipo) BeforeCreate(tx *gorm.DB) error {
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

func (c *ContratoEntidadeTipo) BeforeUpdate(tx *gorm.DB) error {
	if c.UpdatedBy == nil {
		c.UpdatedBy = new(int)
		*c.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o tipo de contrato foi deletado logicamente
func (c *ContratoEntidadeTipo) IsDeleted() bool {
	return c.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (c *ContratoEntidadeTipo) SoftDelete() {
	now := time.Now()
	c.DeletedAt = &now
}

// IsActive verifica se o tipo de contrato está ativo
// Um tipo de contrato é considerado ativo se não foi deletado
func (c *ContratoEntidadeTipo) IsActive() bool {
	return c.DeletedAt == nil
}

// HasReceita verifica se o tipo de contrato está associado a uma receita
func (c *ContratoEntidadeTipo) HasReceita() bool {
	return c.ReceitaID != nil && *c.ReceitaID > 0
}

// HasDespesa verifica se o tipo de contrato está associado a uma despesa
func (c *ContratoEntidadeTipo) HasDespesa() bool {
	return c.DespesaID != nil && *c.DespesaID > 0
}

// GetContratosCount retorna a quantidade de contratos associados
func (c *ContratoEntidadeTipo) GetContratosCount() int {
	return len(c.Contratos)
}
