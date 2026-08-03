package models

import (
	"time"
)

// ============================================================
// MODEL: ProdutoModelo
// ============================================================

type ProdutoModelo struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID        int    `gorm:"column:prom_id;primaryKey;autoIncrement" json:"id"`
	Descricao string `gorm:"column:prom_descricao;type:varchar(255);not null" json:"descricao"`
	Situacao  int    `gorm:"column:prom_situacao;not null;default:0" json:"situacao"`

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
	Produtos []Produto `gorm:"foreignKey:ModeloID;references:ID" json:"produtos,omitempty"`
}

func (ProdutoModelo) TableName() string {
	return "produto_modelo"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *ProdutoModelo) IsActive() bool {
	return m.Situacao == 1
}

func (m *ProdutoModelo) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *ProdutoModelo) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
	m.Situacao = 2
}
