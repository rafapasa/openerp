package models

import (
	"time"
)

// ============================================================
// CONSTANTES
// ============================================================

const (
	LancamentoReceitaPrevista  = 0
	LancamentoReceitaRealizada = 1
)

// ============================================================
// MODEL: LancamentoReceita
// ============================================================

type LancamentoReceita struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID        int       `gorm:"column:lrec_id;primaryKey;autoIncrement" json:"id"`
	ReceitaID *int      `gorm:"column:rec_id" json:"receita_id,omitempty"`
	Historico string    `gorm:"column:lrec_historico;type:varchar(1000);not null" json:"historico"`
	Data      time.Time `gorm:"column:lrec_data;type:date;not null" json:"data"`
	Valor     float64   `gorm:"column:lrec_valor;type:decimal(15,2);not null" json:"valor"`
	TituloID  *int      `gorm:"column:tit_id" json:"titulo_id,omitempty"`
	BaixaItem *int      `gorm:"column:tbx_item" json:"baixa_item,omitempty"`
	Situacao  int       `gorm:"column:lrec_situacao;not null;default:0" json:"situacao"` // 0-Prevista, 1-Realizada

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
	Receita *Receita     `gorm:"foreignKey:ReceitaID;references:rec_id" json:"receita,omitempty"`
	Titulo  *Titulo      `gorm:"foreignKey:TituloID;references:tit_id" json:"titulo,omitempty"`
	Baixa   *TituloBaixa `gorm:"foreignKey:TituloID,BaixaItem;references:tit_id,tbx_item" json:"baixa,omitempty"`
}

func (LancamentoReceita) TableName() string {
	return "lancamento_receita"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (l *LancamentoReceita) IsDeleted() bool {
	return l.DeletedAt != nil
}

func (l *LancamentoReceita) IsPrevista() bool {
	return l.Situacao == LancamentoReceitaPrevista
}

func (l *LancamentoReceita) IsRealizada() bool {
	return l.Situacao == LancamentoReceitaRealizada
}

func (l *LancamentoReceita) SoftDelete() {
	now := time.Now()
	l.DeletedAt = &now
}
