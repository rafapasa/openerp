package models

import (
	"time"

	"gorm.io/gorm"
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

func (m *ProdutoImagem) BeforeCreate(tx *gorm.DB) error {
	// Buscar próximo item
	var maxItem int
	err := tx.Model(&ProdutoImagem{}).
		Where("pro_id = ?", m.ProdutoID).
		Select("COALESCE(MAX(proimg_item), 0) + 1").
		Scan(&maxItem).Error
	if err != nil {
		return err
	}
	m.Item = maxItem

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

func (m *ProdutoImagem) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

func (m *ProdutoImagem) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *ProdutoImagem) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}
