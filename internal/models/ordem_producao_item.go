package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: OrdemProducaoItem
// ============================================================

type OrdemProducaoItem struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	OrdemProducaoID int     `gorm:"column:orp_id;primaryKey" json:"ordem_producao_id"`
	Item            int     `gorm:"column:opi_item;primaryKey" json:"item"`
	ProdutoID       int     `gorm:"column:pro_id;not null" json:"produto_id"`
	QuantSolicitada float64 `gorm:"column:opi_quantsolicitada;type:decimal(15,4);not null" json:"quant_solicitada"`
	QuantProduzida  float64 `gorm:"column:opi_quantproduzida;type:decimal(15,4);not null" json:"quant_produzida"`
	Lote            int     `gorm:"column:opi_lote;not null" json:"lote"`

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
	OrdemProducao *OrdemProducao `gorm:"foreignKey:OrdemProducaoID;references:orp_id" json:"ordem_producao,omitempty"`
	Produto       *Produto       `gorm:"foreignKey:ProdutoID;references:pro_id" json:"produto,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (OrdemProducaoItem) TableName() string {
	return "ordem_producao_item"
}

func (o *OrdemProducaoItem) BeforeCreate(tx *gorm.DB) error {
	if o.CreatedBy == nil {
		o.CreatedBy = new(int)
		*o.CreatedBy = 0
	}
	if o.UpdatedBy == nil {
		o.UpdatedBy = new(int)
		*o.UpdatedBy = 0
	}
	return nil
}

func (o *OrdemProducaoItem) BeforeUpdate(tx *gorm.DB) error {
	if o.UpdatedBy == nil {
		o.UpdatedBy = new(int)
		*o.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o item foi deletado logicamente
func (o *OrdemProducaoItem) IsDeleted() bool {
	return o.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (o *OrdemProducaoItem) SoftDelete() {
	now := time.Now()
	o.DeletedAt = &now
}

// GetQuantRestante retorna a quantidade restante a produzir
func (o *OrdemProducaoItem) GetQuantRestante() float64 {
	return o.QuantSolicitada - o.QuantProduzida
}

// IsFinalizado verifica se o item foi totalmente produzido
func (o *OrdemProducaoItem) IsFinalizado() bool {
	return o.QuantProduzida >= o.QuantSolicitada
}

// GetPercentualProduzido retorna o percentual produzido
func (o *OrdemProducaoItem) GetPercentualProduzido() float64 {
	if o.QuantSolicitada == 0 {
		return 0
	}
	return (o.QuantProduzida / o.QuantSolicitada) * 100
}

// AdicionarProducao adiciona quantidade produzida
func (o *OrdemProducaoItem) AdicionarProducao(quantidade float64) {
	o.QuantProduzida += quantidade
}
