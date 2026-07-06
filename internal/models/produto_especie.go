package models

import (
	"time"

	"github.com/openerp/backend/internal/constants"
)

// ============================================================
// MODEL: ProdutoEspecie
// ============================================================

type ProdutoEspecie struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID        int              `gorm:"column:proesp_id;primaryKey;autoIncrement" json:"id"`
	Descricao string           `gorm:"column:proesp_descricao;type:varchar(255);not null" json:"descricao"`
	Situacao  constants.Status `gorm:"column:proesp_situacao;not null" json:"situacao"`

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
	Produtos []Produto `gorm:"foreignKey:EspecieID;references:ID" json:"produtos,omitempty"`
}

func (ProdutoEspecie) TableName() string {
	return "produto_especie"
}

func (m *ProdutoEspecie) BeforeCreate() error {
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

func (m *ProdutoEspecie) BeforeUpdate() error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *ProdutoEspecie) IsActive() bool {
	return m.Situacao == 1
}

func (m *ProdutoEspecie) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *ProdutoEspecie) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
	m.Situacao = 2
}
