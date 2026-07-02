package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: NotaFiscalDocumentoVenda
// ============================================================

type NotaFiscalDocumentoVenda struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	NotaFiscalID     int `gorm:"column:ntf_id;primaryKey" json:"nota_fiscal_id"`
	DocumentoVendaID int `gorm:"column:ddv_id;primaryKey" json:"documento_venda_id"`

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
	NotaFiscal     *NotaFiscal     `gorm:"foreignKey:NotaFiscalID;references:ntf_id" json:"nota_fiscal,omitempty"`
	DocumentoVenda *DocumentoVenda `gorm:"foreignKey:DocumentoVendaID;references:ddv_id" json:"documento_venda,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (NotaFiscalDocumentoVenda) TableName() string {
	return "nota_fiscal_documento_venda"
}

func (n *NotaFiscalDocumentoVenda) BeforeCreate(tx *gorm.DB) error {
	if n.CreatedBy == nil {
		n.CreatedBy = new(int)
		*n.CreatedBy = 0
	}
	if n.UpdatedBy == nil {
		n.UpdatedBy = new(int)
		*n.UpdatedBy = 0
	}
	return nil
}

func (n *NotaFiscalDocumentoVenda) BeforeUpdate(tx *gorm.DB) error {
	if n.UpdatedBy == nil {
		n.UpdatedBy = new(int)
		*n.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o relacionamento foi deletado logicamente
func (n *NotaFiscalDocumentoVenda) IsDeleted() bool {
	return n.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (n *NotaFiscalDocumentoVenda) SoftDelete() {
	now := time.Now()
	n.DeletedAt = &now
}
