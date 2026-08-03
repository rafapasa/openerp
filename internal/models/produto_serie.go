package models

import (
	"time"

	"github.com/openerp/backend/internal/constants"
)

// ============================================================
// MODEL: ProdutoSerie
// ============================================================

type ProdutoSerie struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID        int              `gorm:"column:pros_id;primaryKey;autoIncrement" json:"id"`
	Descricao string           `gorm:"column:pros_descricao;type:varchar(255);not null" json:"descricao"`
	Situacao  constants.Status `gorm:"column:pros_situacao;not null" json:"situacao"`

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
	Produtos []Produto `gorm:"foreignKey:SerieID;references:ID" json:"produtos,omitempty"`
}

func (ProdutoSerie) TableName() string {
	return "produto_serie"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *ProdutoSerie) IsActive() bool {
	return m.Situacao == 1
}

func (m *ProdutoSerie) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *ProdutoSerie) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
	m.Situacao = 2
}
