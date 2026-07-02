package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: ProdutoProcesso
// ============================================================

type ProdutoProcesso struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ProdutoID             int      `gorm:"column:pro_id;primaryKey" json:"produto_id"`
	Sequencia             int      `gorm:"column:proepe_sequencia;primaryKey" json:"sequencia"`
	EquipamentoProcessoID int      `gorm:"column:epe_id;not null" json:"equipamento_processo_id"`
	Tempo                 *float64 `gorm:"column:proepe_tempo;type:decimal(15,4)" json:"tempo,omitempty"`
	UET                   *float64 `gorm:"column:proepe_uet;type:decimal(15,3)" json:"uet,omitempty"`
	TempoSetup            *float64 `gorm:"column:proepe_tempo_set;type:decimal(15,3)" json:"tempo_setup,omitempty"`
	QtdFuncionarios       *int     `gorm:"column:proepe_qtd_func" json:"qtd_funcionarios,omitempty"`
	Observacao            *string  `gorm:"column:proepe_observacao;type:text" json:"observacao,omitempty"`

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
	Produto *Produto `gorm:"foreignKey:ProdutoID;references:pro_id" json:"produto,omitempty"`
}

func (ProdutoProcesso) TableName() string {
	return "produto_processo"
}

func (m *ProdutoProcesso) BeforeCreate(tx *gorm.DB) error {
	if m.CreatedBy == nil {
		m.CreatedBy = new(int)
		*m.CreatedBy = 0
	}
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

func (m *ProdutoProcesso) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

func (m *ProdutoProcesso) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *ProdutoProcesso) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}
