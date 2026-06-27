package models

import (
	"time"
)

// ============================================================
// MODEL: GrupoEntidade
// ============================================================

type GrupoEntidade struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID              int    `gorm:"column:gpe_id;primaryKey;autoIncrement" json:"id"`
	EmpresaFilialID int    `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	Descricao       string `gorm:"column:gpe_descricao;type:varchar(100);not null" json:"descricao"`
	AtivarLegenda   *int   `gorm:"column:gpe_ativarlegenda" json:"ativar_legenda,omitempty"`
	CorLegenda      *int   `gorm:"column:gpe_corlegenda" json:"cor_legenda,omitempty"`

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
	Entidades     []Entidade     `gorm:"foreignKey:GrupoEntidadeID;references:ID" json:"entidades,omitempty"`
}

func (GrupoEntidade) TableName() string {
	return "grupo_entidade"
}

func (m *GrupoEntidade) BeforeCreate() error {
	if m.CreatedBy == nil {
		m.CreatedBy = new(int)
		*m.CreatedBy = 0
	}
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

func (m *GrupoEntidade) BeforeUpdate() error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *GrupoEntidade) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *GrupoEntidade) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}

func (m *GrupoEntidade) HasLegenda() bool {
	return m.AtivarLegenda != nil && *m.AtivarLegenda == 1
}
