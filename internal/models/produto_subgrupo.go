package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: ProdutoSubgrupo
// ============================================================

type ProdutoSubgrupo struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID        int    `gorm:"column:prosg_id;primaryKey;autoIncrement" json:"id"`
	Descricao string `gorm:"column:prosg_descricao;type:varchar(255);not null" json:"descricao"`
	Situacao  int    `gorm:"column:prosg_situacao;not null;default:0" json:"situacao"`

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
	Produtos []Produto `gorm:"foreignKey:ProdutoSubgrupoID;references:ID" json:"produtos,omitempty"`
	// Despesas     []Despesa     `gorm:"many2many:produto_subgrupo_despesa;foreignKey:prosg_id;joinForeignKey:prosg_id;References:desp_id;joinReferences:desp_id" json:"despesas,omitempty"`
}

func (ProdutoSubgrupo) TableName() string {
	return "produto_subgrupo"
}

func (m *ProdutoSubgrupo) BeforeCreate(tx *gorm.DB) error {
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

func (m *ProdutoSubgrupo) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsActive verifica se o subgrupo está ativo
func (m *ProdutoSubgrupo) IsActive() bool {
	return m.Situacao == 1
}

// IsDeleted verifica se foi deletado logicamente
func (m *ProdutoSubgrupo) IsDeleted() bool {
	return m.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (m *ProdutoSubgrupo) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
	m.Situacao = 2
}
