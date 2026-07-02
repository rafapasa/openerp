package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: CarregamentoEtiqueta
// ============================================================

type CarregamentoEtiqueta struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID               int      `gorm:"column:caretq_id;primaryKey;autoIncrement" json:"id"`
	DocumentoVendaID int      `gorm:"column:ddv_id;not null" json:"documento_venda_id"`
	CarregamentoID   int      `gorm:"column:car_id;not null" json:"carregamento_id"`
	ProdutoID        int      `gorm:"column:pro_id;not null" json:"produto_id"`
	Descricao        *string  `gorm:"column:caretq_descricao;type:varchar(255)" json:"descricao,omitempty"`
	NumeroDaCaixa    *int     `gorm:"column:caretq_numerodacaixa" json:"numero_da_caixa,omitempty"`
	QuantDeCaixas    *int     `gorm:"column:caretq_quantdecaixas" json:"quant_de_caixas,omitempty"`
	Peso             *float64 `gorm:"column:caretq_peso;type:decimal(15,4)" json:"peso,omitempty"`
	Cubagem          *float64 `gorm:"column:caretq_cubagem;type:decimal(15,4)" json:"cubagem,omitempty"`

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
	CarregamentoDocumentoItem *CarregamentoDocumentoItem `gorm:"foreignKey:DocumentoVendaID,CarregamentoID,ProdutoID;references:ddv_id,car_id,pro_id" json:"carregamento_documento_item,omitempty"`
	Itens                     []CarregamentoEtiquetaItem `gorm:"foreignKey:EtiquetaID;references:caretq_id" json:"itens,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (CarregamentoEtiqueta) TableName() string {
	return "carregamento_etiqueta"
}

func (c *CarregamentoEtiqueta) BeforeCreate(tx *gorm.DB) error {
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

func (c *CarregamentoEtiqueta) BeforeUpdate(tx *gorm.DB) error {
	if c.UpdatedBy == nil {
		c.UpdatedBy = new(int)
		*c.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se a etiqueta foi deletada logicamente
func (c *CarregamentoEtiqueta) IsDeleted() bool {
	return c.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (c *CarregamentoEtiqueta) SoftDelete() {
	now := time.Now()
	c.DeletedAt = &now
}

// HasDescricao verifica se possui descrição
func (c *CarregamentoEtiqueta) HasDescricao() bool {
	return c.Descricao != nil && *c.Descricao != ""
}

// HasNumeroDaCaixa verifica se possui número da caixa
func (c *CarregamentoEtiqueta) HasNumeroDaCaixa() bool {
	return c.NumeroDaCaixa != nil && *c.NumeroDaCaixa > 0
}

// HasQuantDeCaixas verifica se possui quantidade de caixas
func (c *CarregamentoEtiqueta) HasQuantDeCaixas() bool {
	return c.QuantDeCaixas != nil && *c.QuantDeCaixas > 0
}

// HasPeso verifica se possui peso
func (c *CarregamentoEtiqueta) HasPeso() bool {
	return c.Peso != nil && *c.Peso > 0
}

// HasCubagem verifica se possui cubagem
func (c *CarregamentoEtiqueta) HasCubagem() bool {
	return c.Cubagem != nil && *c.Cubagem > 0
}

// HasItens verifica se possui itens
func (c *CarregamentoEtiqueta) HasItens() bool {
	return len(c.Itens) > 0
}

// GetItensCount retorna a quantidade de itens
func (c *CarregamentoEtiqueta) GetItensCount() int {
	return len(c.Itens)
}
