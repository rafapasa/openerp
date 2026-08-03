package models

import (
	"time"
)

// ============================================================
// MODEL: ProdutoGrupoComissao
// ============================================================

type ProdutoGrupoComissao struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ProdutoGrupoID     int     `gorm:"column:prog_id;primaryKey" json:"produto_grupo_id"`
	EntidadeID         int     `gorm:"column:ent_id;primaryKey" json:"entidade_id"`
	PercentualComissao float64 `gorm:"column:progc_perccomissao;type:decimal(15,4);not null" json:"percentual_comissao"`

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
	ProdutoGrupo *ProdutoGrupo `gorm:"foreignKey:ProdutoGrupoID;references:prog_id" json:"produto_grupo,omitempty"`
	Entidade     *Entidade     `gorm:"foreignKey:EntidadeID;references:ent_id" json:"entidade,omitempty"`
}

func (ProdutoGrupoComissao) TableName() string {
	return "produto_grupo_comissao"
}

func (m *ProdutoGrupoComissao) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *ProdutoGrupoComissao) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}
