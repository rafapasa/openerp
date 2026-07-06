package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: NotaFiscalHistorico
// ============================================================

type NotaFiscalHistorico struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	NotaFiscalID int       `gorm:"column:ntf_id;primaryKey" json:"nota_fiscal_id"`
	Item         int       `gorm:"column:nfh_item;primaryKey" json:"item"`
	Evento       string    `gorm:"column:nfh_evento;type:varchar(500);not null" json:"evento"`
	Descricao    *string   `gorm:"column:nfh_descricao;type:varchar(2000)" json:"descricao,omitempty"`
	DataHora     time.Time `gorm:"column:nfh_datahora;type:datetime;not null" json:"data_hora"`
	TipoEvento   int       `gorm:"column:nfh_tipoevento;not null" json:"tipo_evento"`

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
	NotaFiscal *NotaFiscal `gorm:"foreignKey:NotaFiscalID;references:ntf_id" json:"nota_fiscal,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (NotaFiscalHistorico) TableName() string {
	return "nota_fiscal_historico"
}

func (n *NotaFiscalHistorico) BeforeCreate(tx *gorm.DB) error {
	if n.CreatedBy == nil {
		n.CreatedBy = new(int)
		*n.CreatedBy = 0
	}
	if n.UpdatedBy == nil {
		n.UpdatedBy = new(int)
		*n.UpdatedBy = 0
	}
	return nil
}

func (n *NotaFiscalHistorico) BeforeUpdate(tx *gorm.DB) error {
	if n.UpdatedBy == nil {
		n.UpdatedBy = new(int)
		*n.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o histórico foi deletado logicamente
func (n *NotaFiscalHistorico) IsDeleted() bool {
	return n.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (n *NotaFiscalHistorico) SoftDelete() {
	now := time.Now()
	n.DeletedAt = &now
}
