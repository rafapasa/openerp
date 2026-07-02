package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: Imposto
// ============================================================

type Imposto struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID              int    `gorm:"column:imp_id;primaryKey;autoIncrement" json:"id"`
	Codigo          string `gorm:"column:imp_codigo;type:varchar(10);not null" json:"codigo"`
	Descricao       string `gorm:"column:imp_descricao;type:varchar(100);not null" json:"descricao"`
	ModalidadeBase  *int   `gorm:"column:imp_modalidadebase" json:"modalidade_base,omitempty"`

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
	NotaFiscalItensImpostos []NotaFiscalItemImposto `gorm:"foreignKey:ImpostoID;references:imp_id" json:"nota_fiscal_itens_impostos,omitempty"`
	DocumentoVendaItensImpostos []DocumentoVendaItemImposto `gorm:"foreignKey:ImpostoID;references:imp_id" json:"documento_venda_itens_impostos,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (Imposto) TableName() string {
	return "imposto"
}

func (i *Imposto) BeforeCreate(tx *gorm.DB) error {
	if i.CreatedBy == nil {
		i.CreatedBy = new(int)
		*i.CreatedBy = 0
	}
	if i.UpdatedBy == nil {
		i.UpdatedBy = new(int)
		*i.UpdatedBy = 0
	}
	return nil
}

func (i *Imposto) BeforeUpdate(tx *gorm.DB) error {
	if i.UpdatedBy == nil {
		i.UpdatedBy = new(int)
		*i.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o imposto foi deletado logicamente
func (i *Imposto) IsDeleted() bool {
	return i.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (i *Imposto) SoftDelete() {
	now := time.Now()
	i.DeletedAt = &now
}

// HasModalidadeBase verifica se possui modalidade de base definida
func (i *Imposto) HasModalidadeBase() bool {
	return i.ModalidadeBase != nil && *i.ModalidadeBase > 0
}

// IsICMS verifica se é ICMS
func (i *Imposto) IsICMS() bool {
	return i.ID == 1 || i.Codigo == "ICMS"
}

// IsICMSST verifica se é ICMS ST
func (i *Imposto) IsICMSST() bool {
	return i.ID == 2 || i.Codigo == "ICMS_ST"
}

// IsIPI verifica se é IPI
func (i *Imposto) IsIPI() bool {
	return i.ID == 3 || i.Codigo == "IPI"
}

// IsPIS verifica se é PIS
func (i *Imposto) IsPIS() bool {
	return i.ID == 4 || i.Codigo == "PIS"
}

// IsCOFINS verifica se é COFINS
func (i *Imposto) IsCOFINS() bool {
	return i.ID == 5 || i.Codigo == "COFINS"
}