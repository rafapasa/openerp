package models

import (
	"time"
)

// ============================================================
// MODEL: DespesaCentroDeCusto
// ============================================================

type DespesaCentroDeCusto struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	CentroCustoID int     `gorm:"column:cdc_id;primaryKey" json:"centro_custo_id"`
	DespesaID     int     `gorm:"column:desp_id;primaryKey" json:"despesa_id"`
	Percentual    float64 `gorm:"column:dcdc_percentual;type:decimal(5,2);not null" json:"percentual"`

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
	Despesa     *Despesa       `gorm:"foreignKey:DespesaID;references:desp_id" json:"despesa,omitempty"`
}

func (DespesaCentroDeCusto) TableName() string {
	return "despesa_centro_de_custo"
}

func (d *DespesaCentroDeCusto) IsDeleted() bool {
	return d.DeletedAt != nil
}

func (d *DespesaCentroDeCusto) SoftDelete() {
	now := time.Now()
	d.DeletedAt = &now
}

// IsActive verifica se o relacionamento despesa-centro de custo está ativo (não deletado)
func (d *DespesaCentroDeCusto) IsActive() bool {
	return d.DeletedAt == nil
}
