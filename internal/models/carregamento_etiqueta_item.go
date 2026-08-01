package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: CarregamentoEtiquetaItem
// ============================================================

type CarregamentoEtiquetaItem struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	EtiquetaID     int     `gorm:"column:caretq_id;primaryKey" json:"etiqueta_id"`
	ProdutoID      int     `gorm:"column:pro_id;primaryKey" json:"produto_id"`
	CarregamentoID int     `gorm:"column:car_id;not null" json:"carregamento_id"`
	QuantEmbalada  float64 `gorm:"column:careti_qunatembalada;type:decimal(15,4);not null" json:"quant_embalada"`
	Peso           float64 `gorm:"column:careti_peso;type:decimal(15,4);not null" json:"peso"`
	Cubagem        float64 `gorm:"column:careti_cubagem;type:decimal(15,4);not null" json:"cubagem"`

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
	Etiqueta     *CarregamentoEtiqueta `gorm:"foreignKey:EtiquetaID;references:caretq_id" json:"etiqueta,omitempty"`
	Produto      *Produto              `gorm:"foreignKey:ProdutoID;references:pro_id" json:"produto,omitempty"`
	Carregamento *Carregamento         `gorm:"foreignKey:CarregamentoID;references:car_id" json:"carregamento,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (CarregamentoEtiquetaItem) TableName() string {
	return "carregamento_etiqueta_item"
}

func (c *CarregamentoEtiquetaItem) BeforeCreate(tx *gorm.DB) error {
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

func (c *CarregamentoEtiquetaItem) BeforeUpdate(tx *gorm.DB) error {
	if c.UpdatedBy == nil {
		c.UpdatedBy = new(int)
		*c.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o item da etiqueta foi deletado logicamente
func (c *CarregamentoEtiquetaItem) IsDeleted() bool {
	return c.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (c *CarregamentoEtiquetaItem) SoftDelete() {
	now := time.Now()
	c.DeletedAt = &now
}

// IsActive verifica se o item da etiqueta está ativo (não deletado)
func (c *CarregamentoEtiquetaItem) IsActive() bool {
	return c.DeletedAt == nil
}

// GetTotalPeso retorna o peso total (quantidade * peso unitário)
func (c *CarregamentoEtiquetaItem) GetTotalPeso() float64 {
	return c.QuantEmbalada * c.Peso
}

// GetTotalCubagem retorna a cubagem total (quantidade * cubagem unitária)
func (c *CarregamentoEtiquetaItem) GetTotalCubagem() float64 {
	return c.QuantEmbalada * c.Cubagem
}
