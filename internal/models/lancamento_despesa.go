package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// CONSTANTES
// ============================================================

const (
	LancamentoDespesaPrevista  = 0
	LancamentoDespesaRealizada = 1
)

// ============================================================
// MODEL: LancamentoDespesa
// ============================================================

type LancamentoDespesa struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID        int       `gorm:"column:ldesp_id;primaryKey;autoIncrement" json:"id"`
	DespesaID int       `gorm:"column:desp_id;not null" json:"despesa_id"`
	Historico string    `gorm:"column:ldesp_historico;type:varchar(1000);not null" json:"historico"`
	Data      time.Time `gorm:"column:ldesp_data;type:date;not null" json:"data"`
	Valor     *float64  `gorm:"column:ldesp_valor;type:decimal(15,2)" json:"valor,omitempty"`
	TituloID  *int      `gorm:"column:tit_id" json:"titulo_id,omitempty"`
	BaixaItem *int      `gorm:"column:tbx_item" json:"baixa_item,omitempty"`
	Situacao  int       `gorm:"column:ldesp_situacao;not null;default:0" json:"situacao"` // 0-Prevista, 1-Realizada

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
	Despesa *Despesa     `gorm:"foreignKey:DespesaID;references:desp_id" json:"despesa,omitempty"`
	Titulo  *Titulo      `gorm:"foreignKey:TituloID;references:tit_id" json:"titulo,omitempty"`
	Baixa   *TituloBaixa `gorm:"foreignKey:TituloID,BaixaItem;references:tit_id,tbx_item" json:"baixa,omitempty"`
}

func (LancamentoDespesa) TableName() string {
	return "lancamento_despesa"
}

func (l *LancamentoDespesa) BeforeCreate(tx *gorm.DB) error {
	if l.CreatedBy == nil {
		l.CreatedBy = new(int)
		*l.CreatedBy = 0
	}
	if l.UpdatedBy == nil {
		l.UpdatedBy = new(int)
		*l.UpdatedBy = 0
	}
	return nil
}

func (l *LancamentoDespesa) BeforeUpdate(tx *gorm.DB) error {
	if l.UpdatedBy == nil {
		l.UpdatedBy = new(int)
		*l.UpdatedBy = 0
	}
	return nil
}

func (l *LancamentoDespesa) IsDeleted() bool {
	return l.DeletedAt != nil
}

func (l *LancamentoDespesa) IsPrevista() bool {
	return l.Situacao == LancamentoDespesaPrevista
}

func (l *LancamentoDespesa) IsRealizada() bool {
	return l.Situacao == LancamentoDespesaRealizada
}

func (l *LancamentoDespesa) SoftDelete() {
	now := time.Now()
	l.DeletedAt = &now
}
