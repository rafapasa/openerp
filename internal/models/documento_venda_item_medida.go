package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: DocumentoVendaItemMedida
// ============================================================

type DocumentoVendaItemMedida struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	DocumentoVendaID   int     `gorm:"column:ddv_id;primaryKey" json:"documento_venda_id"`
	DocumentoVendaItem int     `gorm:"column:dvi_item;primaryKey" json:"documento_venda_item"`
	FormatoProdutoID   int     `gorm:"column:fpro_id;primaryKey" json:"formato_produto_id"`
	FormatoMedidaItem  int     `gorm:"column:fprom_item;primaryKey" json:"formato_medida_item"`
	Valor              float64 `gorm:"column:dvim_valor;type:decimal(15,4);not null" json:"valor"`

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
	FormatoProdutoMedida *FormatoProdutoMedida `gorm:"foreignKey:FormatoProdutoID,FormatoMedidaItem;references:fpro_id,fprom_item" json:"formato_produto_medida,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (DocumentoVendaItemMedida) TableName() string {
	return "documento_venda_item_medida"
}

func (d *DocumentoVendaItemMedida) BeforeCreate(tx *gorm.DB) error {
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

func (d *DocumentoVendaItemMedida) BeforeUpdate(tx *gorm.DB) error {
	if d.UpdatedBy == nil {
		d.UpdatedBy = new(int)
		*d.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o registro foi deletado logicamente
func (d *DocumentoVendaItemMedida) IsDeleted() bool {
	return d.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (d *DocumentoVendaItemMedida) SoftDelete() {
	now := time.Now()
	d.DeletedAt = &now
}

// IsValidValue verifica se o valor está dentro do intervalo permitido da medida
func (d *DocumentoVendaItemMedida) IsValidValue(medida FormatoProdutoMedida) bool {
	return d.Valor >= medida.ValorMinimo && d.Valor <= medida.ValorMaximo
}
