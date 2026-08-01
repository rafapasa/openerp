package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: GradeProduto
// ============================================================

type GradeProduto struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID              int    `gorm:"column:grade_id;primaryKey;autoIncrement" json:"id"`
	EmpresaFilialID int    `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	Nome            string `gorm:"column:grade_nome;type:varchar(255);not null" json:"nome"`

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
	EmpresaFilial *EmpresaFilial     `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
	Itens         []GradeProdutoItem `gorm:"foreignKey:GradeID;references:ID" json:"itens,omitempty"`
	Produtos      []Produto          `gorm:"foreignKey:GradeID;references:ID" json:"produtos,omitempty"`
}

func (GradeProduto) TableName() string {
	return "grade_produto"
}

func (g *GradeProduto) BeforeCreate(tx *gorm.DB) error {
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

func (g *GradeProduto) BeforeUpdate(tx *gorm.DB) error {
	if g.UpdatedBy == nil {
		g.UpdatedBy = new(int)
		*g.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (g *GradeProduto) IsDeleted() bool {
	return g.DeletedAt != nil
}

func (g *GradeProduto) SoftDelete() {
	now := time.Now()
	g.DeletedAt = &now
}

// IsActive verifica se a grade de produto está ativa (não deletada)
func (g *GradeProduto) IsActive() bool {
	return g.DeletedAt == nil
}

func (g *GradeProduto) HasItens() bool {
	return len(g.Itens) > 0
}
