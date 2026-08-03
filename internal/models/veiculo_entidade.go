package models

import (
	"time"
)

// ============================================================
// MODEL: VeiculoEntidade
// ============================================================

type VeiculoEntidade struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	EntidadeID int        `gorm:"column:ent_id;primaryKey" json:"entidade_id"`
	VeiculoID  int        `gorm:"column:vei_id;primaryKey" json:"veiculo_id"`
	DataIni    time.Time  `gorm:"column:vee_dataini;type:date;not null" json:"data_ini"`
	DataFim    *time.Time `gorm:"column:vee_datafim;type:date" json:"data_fim,omitempty"`

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
	Entidade *Entidade `gorm:"foreignKey:EntidadeID;references:ent_id" json:"entidade,omitempty"`
	Veiculo  *Veiculos `gorm:"foreignKey:VeiculoID;references:vei_id" json:"veiculo,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (VeiculoEntidade) TableName() string {
	return "veiculo_entidade"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o relacionamento foi deletado logicamente
func (v *VeiculoEntidade) IsDeleted() bool {
	return v.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (v *VeiculoEntidade) SoftDelete() {
	now := time.Now()
	v.DeletedAt = &now
}

// IsAtivo verifica se o relacionamento está ativo
func (v *VeiculoEntidade) IsAtivo() bool {
	if v.DataFim == nil {
		return true
	}
	return time.Now().Before(*v.DataFim)
}

// IsVencido verifica se o relacionamento está vencido
func (v *VeiculoEntidade) IsVencido() bool {
	if v.DataFim == nil {
		return false
	}
	return time.Now().After(*v.DataFim)
}
