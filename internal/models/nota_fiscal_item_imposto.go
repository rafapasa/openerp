package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: NotaFiscalItemImposto
// ============================================================

type NotaFiscalItemImposto struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	NotaFiscalID   int `gorm:"column:ntf_id;primaryKey" json:"nota_fiscal_id"`
	NotaFiscalItem int `gorm:"column:nfi_item;primaryKey" json:"nota_fiscal_item"`
	ImpostoID      int `gorm:"column:imp_id;primaryKey" json:"imposto_id"`

	ReducaoBase        float64  `gorm:"column:nfiimp_reducaobase;type:decimal(15,2);not null" json:"reducao_base"`
	ValorBase          float64  `gorm:"column:nfiimp_valorbase;type:decimal(15,2);not null" json:"valor_base"`
	Aliquota           *float64 `gorm:"column:nfiimp_aliquota;type:decimal(5,2)" json:"aliquota,omitempty"`
	Valor              float64  `gorm:"column:nfiimp_valor;type:decimal(15,2);not null" json:"valor"`
	AliquotaMVA        *float64 `gorm:"column:nfiimp_aliquotamva;type:decimal(5,2)" json:"aliquota_mva,omitempty"`
	AliquotaCreditosSN *float64 `gorm:"column:nfiimp_aliquotacreditosn;type:decimal(15,2)" json:"aliquota_creditos_sn,omitempty"`
	ValorCreditosSN    *float64 `gorm:"column:nfiimp_valorcreditosn;type:decimal(15,2)" json:"valor_creditos_sn,omitempty"`

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
	Imposto *Imposto `gorm:"foreignKey:ImpostoID;references:imp_id" json:"imposto,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (NotaFiscalItemImposto) TableName() string {
	return "nota_fiscal_item_imposto"
}

func (n *NotaFiscalItemImposto) BeforeCreate(tx *gorm.DB) error {
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

func (n *NotaFiscalItemImposto) BeforeUpdate(tx *gorm.DB) error {
	if n.UpdatedBy == nil {
		n.UpdatedBy = new(int)
		*n.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o imposto foi deletado logicamente
func (n *NotaFiscalItemImposto) IsDeleted() bool {
	return n.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (n *NotaFiscalItemImposto) SoftDelete() {
	now := time.Now()
	n.DeletedAt = &now
}

// HasAliquotaMVA verifica se possui aliquota MVA (Margem de Valor Agregado)
func (n *NotaFiscalItemImposto) HasAliquotaMVA() bool {
	return n.AliquotaMVA != nil && *n.AliquotaMVA > 0
}

// HasCreditosSN verifica se possui créditos do Simples Nacional
func (n *NotaFiscalItemImposto) HasCreditosSN() bool {
	return n.ValorCreditosSN != nil && *n.ValorCreditosSN > 0
}

// HasAliquotaCreditosSN verifica se possui aliquota de créditos do Simples Nacional
func (n *NotaFiscalItemImposto) HasAliquotaCreditosSN() bool {
	return n.AliquotaCreditosSN != nil && *n.AliquotaCreditosSN > 0
}

// HasAliquota verifica se possui aliquota definida
func (n *NotaFiscalItemImposto) HasAliquota() bool {
	return n.Aliquota != nil && *n.Aliquota > 0
}

// GetAliquotaOrDefault retorna a aliquota ou 0 se não definida
func (n *NotaFiscalItemImposto) GetAliquotaOrDefault() float64 {
	if n.HasAliquota() {
		return *n.Aliquota
	}
	return 0
}

// GetAliquotaMVAOrDefault retorna a aliquota MVA ou 0 se não definida
func (n *NotaFiscalItemImposto) GetAliquotaMVAOrDefault() float64 {
	if n.HasAliquotaMVA() {
		return *n.AliquotaMVA
	}
	return 0
}

// GetValorCreditosSNDefault retorna o valor dos créditos SN ou 0 se não definido
func (n *NotaFiscalItemImposto) GetValorCreditosSNDefault() float64 {
	if n.HasCreditosSN() {
		return *n.ValorCreditosSN
	}
	return 0
}

// GetValorBaseComReducao retorna o valor base com a redução aplicada
func (n *NotaFiscalItemImposto) GetValorBaseComReducao() float64 {
	if n.ReducaoBase > 0 {
		return n.ValorBase * (1 - n.ReducaoBase/100)
	}
	return n.ValorBase
}

// CalculaValorImposto calcula o valor do imposto com base na aliquota e valor base
func (n *NotaFiscalItemImposto) CalculaValorImposto() float64 {
	base := n.GetValorBaseComReducao()
	aliquota := n.GetAliquotaOrDefault()
	return base * (aliquota / 100)
}

// IsICMS verifica se o imposto é ICMS (ID 1)
func (n *NotaFiscalItemImposto) IsICMS() bool {
	return n.ImpostoID == 1
}

// IsICMSST verifica se o imposto é ICMS ST (ID 2)
func (n *NotaFiscalItemImposto) IsICMSST() bool {
	return n.ImpostoID == 2
}

// IsIPI verifica se o imposto é IPI (ID 3)
func (n *NotaFiscalItemImposto) IsIPI() bool {
	return n.ImpostoID == 3
}

// IsPIS verifica se o imposto é PIS (ID 4)
func (n *NotaFiscalItemImposto) IsPIS() bool {
	return n.ImpostoID == 4
}

// IsCOFINS verifica se o imposto é COFINS (ID 5)
func (n *NotaFiscalItemImposto) IsCOFINS() bool {
	return n.ImpostoID == 5
}

// GetTipoImposto retorna o nome do tipo de imposto
func (n *NotaFiscalItemImposto) GetTipoImposto() string {
	switch n.ImpostoID {
	case 1:
		return "ICMS"
	case 2:
		return "ICMS ST"
	case 3:
		return "IPI"
	case 4:
		return "PIS"
	case 5:
		return "COFINS"
	default:
		return "Outro"
	}
}
