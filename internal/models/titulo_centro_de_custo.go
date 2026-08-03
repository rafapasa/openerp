package models

import (
	"time"
)

// ============================================================
// MODEL: TituloCentroDeCusto
// ============================================================

type TituloCentroDeCusto struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	CentroCustoID int      `gorm:"column:cdc_id;primaryKey" json:"centro_custo_id"`
	TituloID      int      `gorm:"column:tit_id;primaryKey" json:"titulo_id"`
	Item          int      `gorm:"column:tcdc_item;primaryKey;autoIncrement" json:"item"`
	Percentual    *float64 `gorm:"column:tcdc_percentual;type:decimal(5,2)" json:"percentual,omitempty"`
	Valor         *float64 `gorm:"column:tcdc_valor;type:decimal(15,4)" json:"valor,omitempty"`

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
	CentroCusto *CentroDeCusto `gorm:"foreignKey:CentroCustoID;references:cdc_id" json:"centro_custo,omitempty"`
	Titulo      *Titulo        `gorm:"foreignKey:TituloID;references:tit_id" json:"titulo,omitempty"`
}

func (TituloCentroDeCusto) TableName() string {
	return "titulo_centro_de_custo"
}

func (t *TituloCentroDeCusto) IsDeleted() bool {
	return t.DeletedAt != nil
}

func (t *TituloCentroDeCusto) SoftDelete() {
	now := time.Now()
	t.DeletedAt = &now
}
