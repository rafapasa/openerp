package models

import (
	"time"

	"github.com/openerp/backend/internal/constants"
)

// ============================================================
// MODEL: HistoricoContabil
// ============================================================

type HistoricoContabil struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID              int              `gorm:"column:hisctb_id;primaryKey;autoIncrement" json:"id"`
	EmpresaFilialID int              `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	Texto           string           `gorm:"column:hisctb_texto;type:varchar(255);not null" json:"texto"`
	Situacao        constants.Status `gorm:"column:hisctb_situacao;not null" json:"situacao"` // 1 - ativo, 2 - inativo

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
	EmpresaFilial       *EmpresaFilial       `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
	RotinaContabilItens []RotinaContabilItem `gorm:"foreignKey:HistoricoContabilID;references:hisctb_id" json:"rotina_contabil_itens,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (HistoricoContabil) TableName() string {
	return "historico_contabil"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsActive verifica se o histórico contábil está ativo
func (h *HistoricoContabil) IsActive() bool {
	return h.Situacao == constants.StatusAtivo
}

// IsDeleted verifica se o histórico contábil foi deletado logicamente
func (h *HistoricoContabil) IsDeleted() bool {
	return h.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (h *HistoricoContabil) SoftDelete() {
	now := time.Now()
	h.DeletedAt = &now
	h.Situacao = constants.StatusInativo
}

// IsInactive verifica se o histórico contábil está inativo
func (h *HistoricoContabil) IsInactive() bool {
	return h.Situacao == constants.StatusInativo
}

// HasRotinaContabilItens verifica se possui itens de rotina contábil
func (h *HistoricoContabil) HasRotinaContabilItens() bool {
	return len(h.RotinaContabilItens) > 0
}
