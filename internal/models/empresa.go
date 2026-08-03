// internal/models/empresa.go
package models

import (
	"time"
)

type Empresa struct {
	ID   int    `gorm:"column:emp_id;primaryKey;autoIncrement" json:"id"`
	Nome string `gorm:"column:emp_nome;type:varchar(100)" json:"nome"`

	CreatedAt time.Time  `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime;index" json:"deleted_at,omitempty"`
	CreatedBy *int       `gorm:"column:created_by" json:"created_by,omitempty"`
	UpdatedBy *int       `gorm:"column:updated_by" json:"updated_by,omitempty"`

	Filiais []EmpresaFilial `gorm:"foreignKey:EmpresaID;references:ID" json:"filiais,omitempty"`
}

func (Empresa) TableName() string {
	return "empresa"
}

func (e *Empresa) IsDeleted() bool {
	return e.DeletedAt != nil
}

func (e *Empresa) SoftDelete() {
	now := time.Now()
	e.DeletedAt = &now
}

// IsActive verifica se a empresa está ativa (não deletada)
func (e *Empresa) IsActive() bool {
	return e.DeletedAt == nil
}
