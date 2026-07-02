package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: CarregamentoDocumentoItem
// ============================================================

type CarregamentoDocumentoItem struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	DocumentoVendaID int      `gorm:"column:ddv_id;primaryKey" json:"documento_venda_id"`
	CarregamentoID   int      `gorm:"column:car_id;primaryKey" json:"carregamento_id"`
	ProdutoID        int      `gorm:"column:pro_id;primaryKey" json:"produto_id"`
	QuantSolicitada  float64  `gorm:"column:cardci_quantsolicitada;type:decimal(15,4);not null" json:"quant_solicitada"`
	QuantLida        *float64 `gorm:"column:cardci_quantlida;type:decimal(15,4)" json:"quant_lida,omitempty"`
	Peso             *float64 `gorm:"column:cardci_peso;type:decimal(15,4)" json:"peso,omitempty"`
	Cubagem          *float64 `gorm:"column:cardci_cubagem;type:decimal(15,4)" json:"cubagem,omitempty"`

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
	CarregamentoDocumento *CarregamentoDocumento `gorm:"foreignKey:DocumentoVendaID,CarregamentoID;references:ddv_id,car_id" json:"carregamento_documento,omitempty"`
	Produto               *Produto               `gorm:"foreignKey:ProdutoID;references:pro_id" json:"produto,omitempty"`
	Etiquetas             []CarregamentoEtiqueta `gorm:"foreignKey:DocumentoVendaID,CarregamentoID,ProdutoID;references:ddv_id,car_id,pro_id" json:"etiquetas,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (CarregamentoDocumentoItem) TableName() string {
	return "carregamento_documento_item"
}

func (c *CarregamentoDocumentoItem) BeforeCreate(tx *gorm.DB) error {
	if c.CreatedBy == nil {
		c.CreatedBy = new(int)
		*c.CreatedBy = 0
	}
	if c.UpdatedBy == nil {
		c.UpdatedBy = new(int)
		*c.UpdatedBy = 0
	}
	return nil
}

func (c *CarregamentoDocumentoItem) BeforeUpdate(tx *gorm.DB) error {
	if c.UpdatedBy == nil {
		c.UpdatedBy = new(int)
		*c.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o item foi deletado logicamente
func (c *CarregamentoDocumentoItem) IsDeleted() bool {
	return c.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (c *CarregamentoDocumentoItem) SoftDelete() {
	now := time.Now()
	c.DeletedAt = &now
}

// HasQuantLida verifica se possui quantidade lida
func (c *CarregamentoDocumentoItem) HasQuantLida() bool {
	return c.QuantLida != nil && *c.QuantLida > 0
}

// HasPeso verifica se possui peso
func (c *CarregamentoDocumentoItem) HasPeso() bool {
	return c.Peso != nil && *c.Peso > 0
}

// HasCubagem verifica se possui cubagem
func (c *CarregamentoDocumentoItem) HasCubagem() bool {
	return c.Cubagem != nil && *c.Cubagem > 0
}

// GetQuantRestante retorna a quantidade restante a ser lida
func (c *CarregamentoDocumentoItem) GetQuantRestante() float64 {
	if c.HasQuantLida() {
		return c.QuantSolicitada - *c.QuantLida
	}
	return c.QuantSolicitada
}

// GetPercentualLido retorna o percentual lido
func (c *CarregamentoDocumentoItem) GetPercentualLido() float64 {
	if c.QuantSolicitada == 0 {
		return 0
	}
	if c.HasQuantLida() {
		return (*c.QuantLida / c.QuantSolicitada) * 100
	}
	return 0
}

// IsFinalizado verifica se o item foi totalmente lido
func (c *CarregamentoDocumentoItem) IsFinalizado() bool {
	if c.HasQuantLida() {
		return *c.QuantLida >= c.QuantSolicitada
	}
	return false
}

// HasEtiquetas verifica se possui etiquetas
func (c *CarregamentoDocumentoItem) HasEtiquetas() bool {
	return len(c.Etiquetas) > 0
}

// GetEtiquetasCount retorna a quantidade de etiquetas
func (c *CarregamentoDocumentoItem) GetEtiquetasCount() int {
	return len(c.Etiquetas)
}
