package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: ProdutoEmbalagemItem
// ============================================================

type ProdutoEmbalagemItem struct {
    // ============================================================
    // CAMPOS PRINCIPAIS
    // ============================================================
    ProdutoID     int     `gorm:"column:pro_id;primaryKey" json:"produto_id"`
    EmbalagemItem int     `gorm:"column:proemb_item;primaryKey" json:"embalagem_item"`
    Item          int     `gorm:"column:proebi_tem;primaryKey" json:"item"`
    ProdutoItemID int     `gorm:"column:proebi_pro_id;not null" json:"produto_item_id"`
    Quantidade    float64 `gorm:"column:proebi_quantidade;type:decimal(15,4);not null" json:"quantidade"`

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
    Produto     *Produto `gorm:"foreignKey:ProdutoID;references:pro_id" json:"produto,omitempty"`
    ProdutoItem *Produto `gorm:"foreignKey:ProdutoItemID;references:pro_id" json:"produto_item,omitempty"`
}

func (ProdutoEmbalagemItem) TableName() string {
    return "produto_embalagem_item"
}

func (m *ProdutoEmbalagemItem) BeforeCreate(tx *gorm.DB) error {
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

func (m *ProdutoEmbalagemItem) BeforeUpdate(tx *gorm.DB) error {
    if m.UpdatedBy == nil {
        m.UpdatedBy = new(int)
        *m.UpdatedBy = 0
    }
    return nil
}

func (m *ProdutoEmbalagemItem) IsDeleted() bool {
    return m.DeletedAt != nil
}

func (m *ProdutoEmbalagemItem) SoftDelete() {
    now := time.Now()
    m.DeletedAt = &now
}