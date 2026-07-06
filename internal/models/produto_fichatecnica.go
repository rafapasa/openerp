package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: ProdutoFichaTecnica
// ============================================================

type ProdutoFichaTecnica struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ProdutoID           int      `gorm:"column:pro_id;primaryKey" json:"produto_id"`
	ProdutoComponenteID int      `gorm:"column:pro_pro_id;primaryKey" json:"produto_componente_id"`
	Quantidade          float64  `gorm:"column:prf_quantidade;type:decimal(15,4);not null" json:"quantidade"`
	Peso                float64  `gorm:"column:prf_peso;type:decimal(15,4);not null" json:"peso"`
	Custo               *float64 `gorm:"column:prf_custo;type:decimal(15,4)" json:"custo,omitempty"`

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
	Produto           *Produto `gorm:"foreignKey:ProdutoID;references:pro_id" json:"produto,omitempty"`
	ProdutoComponente *Produto `gorm:"foreignKey:ProdutoComponenteID;references:pro_id" json:"produto_componente,omitempty"`
}

func (ProdutoFichaTecnica) TableName() string {
	return "produto_fichatecnica"
}

func (m *ProdutoFichaTecnica) BeforeCreate(tx *gorm.DB) error {
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

func (m *ProdutoFichaTecnica) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

func (m *ProdutoFichaTecnica) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *ProdutoFichaTecnica) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}
