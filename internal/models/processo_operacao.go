package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: ProcessoOperacao
// ============================================================

type ProcessoOperacao struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	OperacaoID int `gorm:"column:ope_id;primaryKey" json:"operacao_id"`
	ProcessoID int `gorm:"column:prc_id;primaryKey" json:"processo_id"`

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
	Operacao *Operacao `gorm:"foreignKey:OperacaoID;references:ope_id" json:"operacao,omitempty"`
	Processo *Processo `gorm:"foreignKey:ProcessoID;references:prc_id" json:"processo,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (ProcessoOperacao) TableName() string {
	return "processo_operacao"
}

func (p *ProcessoOperacao) BeforeCreate(tx *gorm.DB) error {
	if p.CreatedBy == nil {
		p.CreatedBy = new(int)
		*p.CreatedBy = 0
	}
	if p.UpdatedBy == nil {
		p.UpdatedBy = new(int)
		*p.UpdatedBy = 0
	}
	return nil
}

func (p *ProcessoOperacao) BeforeUpdate(tx *gorm.DB) error {
	if p.UpdatedBy == nil {
		p.UpdatedBy = new(int)
		*p.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o relacionamento foi deletado logicamente
func (p *ProcessoOperacao) IsDeleted() bool {
	return p.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (p *ProcessoOperacao) SoftDelete() {
	now := time.Now()
	p.DeletedAt = &now
}
