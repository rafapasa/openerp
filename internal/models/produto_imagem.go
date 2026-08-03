package models

import (
	"time"
)

// ============================================================
// MODEL: ProdutoImagem
// ============================================================

type ProdutoImagem struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ProdutoID int     `gorm:"column:pro_id;primaryKey" json:"produto_id"`
	Item      int     `gorm:"column:proimg_item;primaryKey" json:"item"`
	Descricao *string `gorm:"column:proimg_descricao;type:text" json:"descricao,omitempty"`
	Imagem    []byte  `gorm:"column:proimg_imagem;type:longblob;not null" json:"-"` // Ocultar no JSON

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
	Produto *Produto `gorm:"foreignKey:ProdutoID;references:pro_id" json:"produto,omitempty"`
}

func (ProdutoImagem) TableName() string {
	return "produto_imagem"
}

// Buscar próximo item

func (m *ProdutoImagem) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *ProdutoImagem) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}
