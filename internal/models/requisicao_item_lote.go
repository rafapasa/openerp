package models

import (
	"time"
)

// ============================================================
// MODEL: RequisicaoItemLote
// ============================================================

type RequisicaoItemLote struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	RequisicaoID    int       `gorm:"column:req_id;primaryKey" json:"requisicao_id"`
	RequisicaoItem  int       `gorm:"column:rei_item;primaryKey" json:"requisicao_item"`
	Item            int       `gorm:"column:reil_item;primaryKey" json:"item"`
	ProdutoID       int       `gorm:"column:pro_id;not null" json:"produto_id"`
	ProdutoLoteItem int       `gorm:"column:prol_item;not null" json:"produto_lote_item"`
	Quantidade      float64   `gorm:"column:reil_quantidade;type:decimal(15,4);not null" json:"quantidade"`
	Lote            string    `gorm:"column:reil_lote;type:varchar(30);not null" json:"lote"`
	DataFabricacao  time.Time `gorm:"column:reil_datafabricacao;type:datetime;not null" json:"data_fabricacao"`
	DataValidade    time.Time `gorm:"column:reil_datavalidade;type:datetime;not null" json:"data_validade"`

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
	Produto     *Produto     `gorm:"foreignKey:ProdutoID;references:pro_id" json:"produto,omitempty"`
	ProdutoLote *ProdutoLote `gorm:"foreignKey:ProdutoID,ProdutoLoteItem;references:pro_id,prol_item" json:"produto_lote,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (RequisicaoItemLote) TableName() string {
	return "requisicao_item_lote"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o registro foi deletado logicamente
func (r *RequisicaoItemLote) IsDeleted() bool {
	return r.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (r *RequisicaoItemLote) SoftDelete() {
	now := time.Now()
	r.DeletedAt = &now
}

// IsVencido verifica se o lote está vencido
func (r *RequisicaoItemLote) IsVencido() bool {
	return time.Now().After(r.DataValidade)
}

// IsProximoVencer verifica se o lote está próximo de vencer (30 dias)
func (r *RequisicaoItemLote) IsProximoVencer() bool {
	diasRestantes := time.Until(r.DataValidade).Hours() / 24
	return diasRestantes > 0 && diasRestantes <= 30
}
