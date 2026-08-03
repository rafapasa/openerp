package models

import (
	"time"
)

// ============================================================
// MODEL: ProdutoTamanho
// ============================================================

type ProdutoTamanho struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID              int    `gorm:"column:ptam_id;primaryKey;autoIncrement" json:"id"`
	EmpresaFilialID int    `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	Sigla           string `gorm:"column:ptam_sigla;type:varchar(20);not null" json:"sigla"`
	Nome            string `gorm:"column:ptam_nome;type:varchar(255);not null" json:"nome"`

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
}

func (ProdutoTamanho) TableName() string {
	return "produto_tamanho"
}

func (m *ProdutoTamanho) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *ProdutoTamanho) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}
