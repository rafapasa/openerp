package models

import (
	"time"

	"github.com/openerp/backend/internal/constants"
)

// ============================================================
// MODEL: TipoDocumento
// ============================================================

type TipoDocumento struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID              int                     `gorm:"column:tdoc_id;primaryKey;autoIncrement" json:"id"`
	EmpresaFilialID int                     `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	Nome            string                  `gorm:"column:tdoc_descricao;type:varchar(255);not null" json:"nome"`
	Comissao        int                     `gorm:"column:tdoc_comissoa;type:smallint;not null" json:"comissao"`
	Tipo            constants.TipoDocumento `gorm:"column:tdoc_tiponfe;not null;default:1" json:"tipo"` // CORRIGIDO: usando constantes
	Situacao        constants.Status        `gorm:"column:tdoc_situacao;not null;default:1" json:"situacao"`

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
	EmpresaFilial *EmpresaFilial `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (TipoDocumento) TableName() string {
	return "tipo_documento"
}

func (m *TipoDocumento) BeforeCreate() error {
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

func (m *TipoDocumento) BeforeUpdate() error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsActive verifica se o tipo de documento está ativo
func (m *TipoDocumento) IsActive() bool {
	return m.Situacao.IsActive()
}

// IsDeleted verifica se foi deletado logicamente
func (m *TipoDocumento) IsDeleted() bool {
	return m.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (m *TipoDocumento) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
	m.Situacao = constants.StatusInativo
}

// HasComissao verifica se o documento tem comissão
func (m *TipoDocumento) HasComissao() bool {
	return m.Comissao == 1
}

// IsBoleto verifica se é boleto
func (m *TipoDocumento) IsBoleto() bool {
	return m.Tipo == constants.TipoDocumentoBoleto
}

// IsDuplicata verifica se é duplicata
func (m *TipoDocumento) IsDuplicata() bool {
	return m.Tipo == constants.TipoDocumentoDuplicata
}

// IsCheque verifica se é cheque
func (m *TipoDocumento) IsCheque() bool {
	return m.Tipo == constants.TipoDocumentoCheque
}
