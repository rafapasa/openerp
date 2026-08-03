package models

import (
	"time"
)

// ============================================================
// MODEL: DocumentoVendaHistorico
// Registra as mudanças de situação (status) de um Documento de Venda.
// ============================================================

type DocumentoVendaHistorico struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID               int    `gorm:"column:ddvh_id;primaryKey;autoIncrement" json:"id"`
	DocumentoVendaID int    `gorm:"column:ddv_id;not null" json:"documento_venda_id"`
	SituacaoAnterior int    `gorm:"column:ddvh_situacaoanterior;not null" json:"situacao_anterior"` // ID da situação anterior
	SituacaoAtual    int    `gorm:"column:ddvh_situacaoatual;not null" json:"situacao_atual"`       // ID da situação atual
	Observacao       string `gorm:"column:ddvh_observacao;type:varchar(255)" json:"observacao"`

	// ============================================================
	// CAMPOS DE AUDITORIA
	// ============================================================
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy *int      `gorm:"column:created_by" json:"created_by,omitempty"`

	// ============================================================
	// RELACIONAMENTOS
	// ============================================================
	DocumentoVenda      *DocumentoVenda         `gorm:"foreignKey:DocumentoVendaID;references:ID" json:"documento_venda,omitempty"`
	SituacaoAnteriorRef *DocumentoVendaSituacao `gorm:"foreignKey:SituacaoAnterior;references:ID" json:"situacao_anterior_ref,omitempty"`
	SituacaoAtualRef    *DocumentoVendaSituacao `gorm:"foreignKey:SituacaoAtual;references:ID" json:"situacao_atual_ref,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (DocumentoVendaHistorico) TableName() string {
	return "documento_venda_historico"
}

// UpdatedBy e UpdatedAt não são aplicáveis para histórico de criação

// Não há SoftDelete para histórico, pois os registros devem ser imutáveis.
