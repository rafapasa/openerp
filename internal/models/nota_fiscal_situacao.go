package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: NotaFiscalSituacao
// ============================================================

type NotaFiscalSituacao struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID        int    `gorm:"column:nfsit_id;primaryKey;autoIncrement" json:"id"`
	Descricao string `gorm:"column:nfsit_descricao;type:varchar(255);not null" json:"descricao"`

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
	NotaFiscais []NotaFiscal `gorm:"foreignKey:SituacaoID;references:nfsit_id" json:"notas_fiscais,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (NotaFiscalSituacao) TableName() string {
	return "nota_fiscal_situacao"
}

func (n *NotaFiscalSituacao) BeforeCreate(tx *gorm.DB) error {
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

func (n *NotaFiscalSituacao) BeforeUpdate(tx *gorm.DB) error {
	if n.UpdatedBy == nil {
		n.UpdatedBy = new(int)
		*n.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se a situação foi deletada logicamente
func (n *NotaFiscalSituacao) IsDeleted() bool {
	return n.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (n *NotaFiscalSituacao) SoftDelete() {
	now := time.Now()
	n.DeletedAt = &now
}
