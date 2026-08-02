package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: DocumentoVendaSituacao
// Representa as possíveis situações (status) de um Documento de Venda.
// ============================================================

type DocumentoVendaSituacao struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID        int    `gorm:"column:ddvs_id;primaryKey;autoIncrement" json:"id"`
	Descricao string `gorm:"column:ddvs_descricao;type:varchar(100);not null" json:"descricao"`

	// ============================================================
	// CAMPOS DE AUDITORIA
	// ============================================================
	CreatedAt time.Time  `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime;index" json:"deleted_at,omitempty"`
	CreatedBy *int       `gorm:"column:created_by" json:"created_by,omitempty"`
	UpdatedBy *int       `gorm:"column:updated_by" json:"updated_by,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (DocumentoVendaSituacao) TableName() string {
	return "documento_venda_situacao"
}

func (d *DocumentoVendaSituacao) BeforeCreate(tx *gorm.DB) error {
	if d.CreatedBy == nil {
		d.CreatedBy = new(int)
		*d.CreatedBy = 0
	}
	if d.UpdatedBy == nil {
		d.UpdatedBy = new(int)
		*d.UpdatedBy = 0
	}
	return nil
}