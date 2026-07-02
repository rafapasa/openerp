package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: DocumentoVendaOrdemServico
// ============================================================

type DocumentoVendaOrdemServico struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	OrdemServicoID   int `gorm:"column:os_id;primaryKey" json:"ordem_servico_id"`
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
	OrdemServico   *OrdemServico   `gorm:"foreignKey:OrdemServicoID;references:os_id" json:"ordem_servico,omitempty"`
	DocumentoVenda *DocumentoVenda `gorm:"foreignKey:DocumentoVendaID;references:ddv_id" json:"documento_venda,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (DocumentoVendaOrdemServico) TableName() string {
	return "documento_venda_ordem_servico"
}

func (d *DocumentoVendaOrdemServico) BeforeCreate(tx *gorm.DB) error {
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

func (d *DocumentoVendaOrdemServico) BeforeUpdate(tx *gorm.DB) error {
	if d.UpdatedBy == nil {
		d.UpdatedBy = new(int)
		*d.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o relacionamento foi deletado logicamente
func (d *DocumentoVendaOrdemServico) IsDeleted() bool {
	return d.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (d *DocumentoVendaOrdemServico) SoftDelete() {
	now := time.Now()
	d.DeletedAt = &now
}