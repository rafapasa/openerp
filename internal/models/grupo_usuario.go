package models

import (
	"time"

	"gorm.io/gorm"
)

type GrupoUsuario struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID              int    `gorm:"column:gpu_id;primaryKey;autoIncrement" json:"id"`
	Descricao       string `gorm:"column:gpu_descricao;type:varchar(100);not null" json:"descricao"`
	EmpresaFilialID int    `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	Situacao        int    `gorm:"column:gpu_situacao;not null" json:"situacao"`

	// ============================================================
	// CAMPOS DE AUDITORIA (sem prefixo gpu_)
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
	Usuarios      []Usuario      `gorm:"foreignKey:GrupoUsuarioID;references:ID" json:"usuarios,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (GrupoUsuario) TableName() string {
	return "grupo_usuario"
}

func (g *GrupoUsuario) BeforeCreate(tx *gorm.DB) error {
	if g.CreatedBy == nil {
		g.CreatedBy = new(int)
		*g.CreatedBy = 0
	}
	if g.UpdatedBy == nil {
		g.UpdatedBy = new(int)
		*g.UpdatedBy = 0
	}
	return nil
}

func (g *GrupoUsuario) BeforeUpdate(tx *gorm.DB) error {
	if g.UpdatedBy == nil {
		g.UpdatedBy = new(int)
		*g.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (g *GrupoUsuario) IsActive() bool {
	return g.Situacao == 1
}

// CORRIGIDO: IsDeleted com D maiúsculo
func (g *GrupoUsuario) IsDeleted() bool {
	return g.DeletedAt != nil
}

func (g *GrupoUsuario) SoftDelete() {
	now := time.Now()
	g.DeletedAt = &now
	g.Situacao = 0 // Marca como inativo também
}
