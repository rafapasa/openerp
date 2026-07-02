package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: OrdemProducaoExplosao
// ============================================================

type OrdemProducaoExplosao struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	OrdemProducaoID int      `gorm:"column:orp_id;primaryKey" json:"ordem_producao_id"`
	ProdutoID       int      `gorm:"column:pro_id;primaryKey" json:"produto_id"`
	QuantNecessaria float64  `gorm:"column:ope_quantnecessaria;type:decimal(15,4);not null" json:"quant_necessaria"`
	QuantBaixada    *float64 `gorm:"column:ope_quantbaixada;type:decimal(15,4)" json:"quant_baixada,omitempty"`
	QuantAlocada    *float64 `gorm:"column:ope_quantalucada;type:decimal(15,4)" json:"quant_alocada,omitempty"`

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
	OrdemProducao *OrdemProducao `gorm:"foreignKey:OrdemProducaoID;references:orp_id" json:"ordem_producao,omitempty"`
	Produto       *Produto       `gorm:"foreignKey:ProdutoID;references:pro_id" json:"produto,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (OrdemProducaoExplosao) TableName() string {
	return "ordem_producao_explosao"
}

func (o *OrdemProducaoExplosao) BeforeCreate(tx *gorm.DB) error {
	if o.CreatedBy == nil {
		o.CreatedBy = new(int)
		*o.CreatedBy = 0
	}
	if o.UpdatedBy == nil {
		o.UpdatedBy = new(int)
		*o.UpdatedBy = 0
	}
	return nil
}

func (o *OrdemProducaoExplosao) BeforeUpdate(tx *gorm.DB) error {
	if o.UpdatedBy == nil {
		o.UpdatedBy = new(int)
		*o.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o registro foi deletado logicamente
func (o *OrdemProducaoExplosao) IsDeleted() bool {
	return o.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (o *OrdemProducaoExplosao) SoftDelete() {
	now := time.Now()
	o.DeletedAt = &now
}

// HasQuantBaixada verifica se possui quantidade baixada
func (o *OrdemProducaoExplosao) HasQuantBaixada() bool {
	return o.QuantBaixada != nil && *o.QuantBaixada > 0
}

// HasQuantAlocada verifica se possui quantidade alocada
func (o *OrdemProducaoExplosao) HasQuantAlocada() bool {
	return o.QuantAlocada != nil && *o.QuantAlocada > 0
}

// GetQuantRestante retorna a quantidade restante
func (o *OrdemProducaoExplosao) GetQuantRestante() float64 {
	if o.HasQuantBaixada() {
		return o.QuantNecessaria - *o.QuantBaixada
	}
	return o.QuantNecessaria
}
