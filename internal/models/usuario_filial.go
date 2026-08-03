package models

import (
	"time"
)

// ============================================================
// MODEL: UsuarioFilial
// ============================================================

type UsuarioFilial struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID              int `gorm:"column:ufi_id;primaryKey;autoIncrement" json:"id"`
	UsuarioID       int `gorm:"column:usu_id;not null" json:"usuario_id"`
	EmpresaFilialID int `gorm:"column:emf_id;not null" json:"empresa_filial_id"`

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
	Usuario       *Usuario       `gorm:"foreignKey:UsuarioID;references:usu_id" json:"usuario,omitempty"`
	EmpresaFilial *EmpresaFilial `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
}

func (UsuarioFilial) TableName() string {
	return "usuario_filial"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *UsuarioFilial) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *UsuarioFilial) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}
