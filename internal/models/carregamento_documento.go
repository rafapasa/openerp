package models

import (
	"time"
)

// ============================================================
// MODEL: CarregamentoDocumento
// ============================================================

type CarregamentoDocumento struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	DocumentoVendaID int `gorm:"column:ddv_id;primaryKey" json:"documento_venda_id"`
	CarregamentoID   int `gorm:"column:car_id;primaryKey" json:"carregamento_id"`

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
	DocumentoVenda *DocumentoVenda             `gorm:"foreignKey:DocumentoVendaID;references:ddv_id" json:"documento_venda,omitempty"`
	Carregamento   *Carregamento               `gorm:"foreignKey:CarregamentoID;references:car_id" json:"carregamento,omitempty"`
	Itens          []CarregamentoDocumentoItem `gorm:"foreignKey:DocumentoVendaID,CarregamentoID;references:ddv_id,car_id" json:"itens,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (CarregamentoDocumento) TableName() string {
	return "carregamento_documento"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o documento foi deletado logicamente
func (c *CarregamentoDocumento) IsDeleted() bool {
	return c.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (c *CarregamentoDocumento) SoftDelete() {
	now := time.Now()
	c.DeletedAt = &now
}

// HasItens verifica se possui itens
func (c *CarregamentoDocumento) HasItens() bool {
	return len(c.Itens) > 0
}

// GetItensCount retorna a quantidade de itens
func (c *CarregamentoDocumento) GetItensCount() int {
	return len(c.Itens)
}
