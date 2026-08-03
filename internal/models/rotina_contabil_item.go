package models

import (
	"time"
)

// ============================================================
// MODEL: RotinaContabilItem
// ============================================================

type RotinaContabilItem struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	RotinaContabilID       int  `gorm:"column:roc_id;primaryKey" json:"rotina_contabil_id"`
	Item                   int  `gorm:"column:roci_item;primaryKey" json:"item"`
	CentroDeCustoCreditoID *int `gorm:"column:cdc_id_credito" json:"centro_de_custo_credito_id,omitempty"`
	HistoricoContabilID    int  `gorm:"column:hisctb_id;not null" json:"historico_contabil_id"`
	CentroDeCustoDebitoID  *int `gorm:"column:cdc_id_debito" json:"centro_de_custo_debito_id,omitempty"`
	PlanoContaCreditoID    int  `gorm:"column:pcc_id_credito;not null" json:"plano_conta_credito_id"`
	EventoNumero           int  `gorm:"column:eve_numero;not null" json:"evento_numero"`
	PlanoContaDebitoID     int  `gorm:"column:pcc_id_debito;not null" json:"plano_conta_debito_id"`

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
	RotinaContabil       *RotinaContabil    `gorm:"foreignKey:RotinaContabilID;references:roc_id" json:"rotina_contabil,omitempty"`
	CentroDeCustoCredito *CentroDeCusto     `gorm:"foreignKey:CentroDeCustoCreditoID;references:cdc_id" json:"centro_de_custo_credito,omitempty"`
	CentroDeCustoDebito  *CentroDeCusto     `gorm:"foreignKey:CentroDeCustoDebitoID;references:cdc_id" json:"centro_de_custo_debito,omitempty"`
	HistoricoContabil    *HistoricoContabil `gorm:"foreignKey:HistoricoContabilID;references:hisctb_id" json:"historico_contabil,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (RotinaContabilItem) TableName() string {
	return "rotina_contabil_itens"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o item foi deletado logicamente
func (r *RotinaContabilItem) IsDeleted() bool {
	return r.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (r *RotinaContabilItem) SoftDelete() {
	now := time.Now()
	r.DeletedAt = &now
}

// HasCentroDeCustoCredito verifica se possui centro de custo crédito
func (r *RotinaContabilItem) HasCentroDeCustoCredito() bool {
	return r.CentroDeCustoCreditoID != nil && *r.CentroDeCustoCreditoID > 0
}

// HasCentroDeCustoDebito verifica se possui centro de custo débito
func (r *RotinaContabilItem) HasCentroDeCustoDebito() bool {
	return r.CentroDeCustoDebitoID != nil && *r.CentroDeCustoDebitoID > 0
}
