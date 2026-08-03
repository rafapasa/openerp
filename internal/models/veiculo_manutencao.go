package models

import (
	"time"
)

// ============================================================
// MODEL: VeiculoManutencao
// ============================================================

type VeiculoManutencao struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	VeiculoID    int       `gorm:"column:vei_id;primaryKey" json:"veiculo_id"`
	ManutencaoID int       `gorm:"column:man_id;primaryKey" json:"manutencao_id"`
	Item         int       `gorm:"column:veiman_item;primaryKey" json:"item"`
	Data         time.Time `gorm:"column:veiman_data;type:date;not null" json:"data"`
	Hodometro    int       `gorm:"column:veiman_hodometro;not null" json:"hodometro"`

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
	Veiculo    *Veiculos   `gorm:"foreignKey:VeiculoID;references:vei_id" json:"veiculo,omitempty"`
	Manutencao *Manutencao `gorm:"foreignKey:ManutencaoID;references:man_id" json:"manutencao,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (VeiculoManutencao) TableName() string {
	return "veiculo_manutencao"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se a manutenção foi deletada logicamente
func (v *VeiculoManutencao) IsDeleted() bool {
	return v.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (v *VeiculoManutencao) SoftDelete() {
	now := time.Now()
	v.DeletedAt = &now
}
