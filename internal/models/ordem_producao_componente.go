package models

import (
	"time"
)

// ============================================================
// MODEL: OrdemProducaoComponente
// ============================================================

type OrdemProducaoComponente struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	OrdemProducaoID int     `gorm:"column:orp_id;primaryKey" json:"ordem_producao_id"`
	ProdutoID       int     `gorm:"column:pro_id;primaryKey" json:"produto_id"`
	QuantSolicitada float64 `gorm:"column:opc_quantsolicitada;type:decimal(15,4);not null" json:"quant_solicitada"`
	QuantProduzida  float64 `gorm:"column:opc_quantproduzida;type:decimal(15,4);not null" json:"quant_produzida"`
	QuantBaixada    float64 `gorm:"column:opc_quantbaixada;type:decimal(15,4);not null" json:"quant_baixada"`
	QuantAlocada    float64 `gorm:"column:opc_quantalocada;type:decimal(15,4);not null" json:"quant_alocada"`

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

func (OrdemProducaoComponente) TableName() string {
	return "ordem_producao_componente"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o componente foi deletado logicamente
func (o *OrdemProducaoComponente) IsDeleted() bool {
	return o.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (o *OrdemProducaoComponente) SoftDelete() {
	now := time.Now()
	o.DeletedAt = &now
}

// GetQuantRestante retorna a quantidade restante do componente
func (o *OrdemProducaoComponente) GetQuantRestante() float64 {
	return o.QuantSolicitada - o.QuantBaixada
}

// GetPercentualBaixado retorna o percentual baixado
func (o *OrdemProducaoComponente) GetPercentualBaixado() float64 {
	if o.QuantSolicitada == 0 {
		return 0
	}
	return (o.QuantBaixada / o.QuantSolicitada) * 100
}

// IsFinalizado verifica se o componente foi totalmente baixado
func (o *OrdemProducaoComponente) IsFinalizado() bool {
	return o.QuantBaixada >= o.QuantSolicitada
}
