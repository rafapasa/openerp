package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: DocumentoVendaItemImposto
// ============================================================

type DocumentoVendaItemImposto struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	DocumentoVendaID   int `gorm:"column:ddv_id;primaryKey" json:"documento_venda_id"`
	DocumentoVendaItem int `gorm:"column:dvi_item;primaryKey" json:"documento_venda_item"`
	ImpostoID          int `gorm:"column:imp_id;primaryKey" json:"imposto_id"`

	ReducaoBase        float64  `gorm:"column:dviimp_reducaobase;type:decimal(15,2);not null" json:"reducao_base"`
	ValorBase          float64  `gorm:"column:dviimp_valorbase;type:decimal(15,2);not null" json:"valor_base"`
	Aliquota           float64  `gorm:"column:dviimp_aliquota;type:decimal(5,2);not null" json:"aliquota"`
	ValorImposto       float64  `gorm:"column:dviimp_valorimposto;type:decimal(15,2);not null" json:"valor_imposto"`
	AliquotaMVA        *float64 `gorm:"column:dviimp_aliquotamva;type:decimal(5,2)" json:"aliquota_mva,omitempty"`
	AliquotaCreditosSN *float64 `gorm:"column:dviimp_aliquotacreditosn;type:decimal(5,2)" json:"aliquota_creditos_sn,omitempty"`
	ValorCreditosSN    *float64 `gorm:"column:dviimp_valorcreditosn;type:decimal(15,2)" json:"valor_creditos_sn,omitempty"`

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

func (DocumentoVendaItemImposto) TableName() string {
	return "documento_venda_item_imposto"
}

func (d *DocumentoVendaItemImposto) BeforeCreate(tx *gorm.DB) error {
	if d.CreatedBy == nil {
		d.CreatedBy = new(int)
		*d.CreatedBy = 0
	}
	if d.UpdatedBy == nil {
		d.UpdatedBy = new(int)
		*d.UpdatedBy = 0
	}
	return nil
}

func (d *DocumentoVendaItemImposto) BeforeUpdate(tx *gorm.DB) error {
	if d.UpdatedBy == nil {
		d.UpdatedBy = new(int)
		*d.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o imposto foi deletado logicamente
func (d *DocumentoVendaItemImposto) IsDeleted() bool {
	return d.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (d *DocumentoVendaItemImposto) SoftDelete() {
	now := time.Now()
	d.DeletedAt = &now
}

// HasAliquotaMVA verifica se possui aliquota MVA (Margem de Valor Agregado)
func (d *DocumentoVendaItemImposto) HasAliquotaMVA() bool {
	return d.AliquotaMVA != nil && *d.AliquotaMVA > 0
}

// HasCreditosSN verifica se possui créditos do Simples Nacional
func (d *DocumentoVendaItemImposto) HasCreditosSN() bool {
	return d.ValorCreditosSN != nil && *d.ValorCreditosSN > 0
}

// HasAliquotaCreditosSN verifica se possui aliquota de créditos do Simples Nacional
func (d *DocumentoVendaItemImposto) HasAliquotaCreditosSN() bool {
	return d.AliquotaCreditosSN != nil && *d.AliquotaCreditosSN > 0
}

// GetAliquotaMVAOrDefault retorna a aliquota MVA ou 0 se não definida
func (d *DocumentoVendaItemImposto) GetAliquotaMVAOrDefault() float64 {
	if d.HasAliquotaMVA() {
		return *d.AliquotaMVA
	}
	return 0
}

// GetValorCreditosSNDefault retorna o valor dos créditos SN ou 0 se não definido
func (d *DocumentoVendaItemImposto) GetValorCreditosSNDefault() float64 {
	if d.HasCreditosSN() {
		return *d.ValorCreditosSN
	}
	return 0
}

// GetValorBaseComReducao retorna o valor base com a redução aplicada
func (d *DocumentoVendaItemImposto) GetValorBaseComReducao() float64 {
	if d.ReducaoBase > 0 {
		return d.ValorBase * (1 - d.ReducaoBase/100)
	}
	return d.ValorBase
}

// CalculaValorImposto calcula o valor do imposto com base na aliquota e valor base
func (d *DocumentoVendaItemImposto) CalculaValorImposto() float64 {
	base := d.GetValorBaseComReducao()
	return base * (d.Aliquota / 100)
}

// IsICMS verifica se o imposto é ICMS (ID 1)
func (d *DocumentoVendaItemImposto) IsICMS() bool {
	return d.ImpostoID == 1
}

// IsICMSST verifica se o imposto é ICMS ST (ID 2)
func (d *DocumentoVendaItemImposto) IsICMSST() bool {
	return d.ImpostoID == 2
}

// IsIPI verifica se o imposto é IPI (ID 3)
func (d *DocumentoVendaItemImposto) IsIPI() bool {
	return d.ImpostoID == 3
}

// IsPIS verifica se o imposto é PIS (ID 4)
func (d *DocumentoVendaItemImposto) IsPIS() bool {
	return d.ImpostoID == 4
}

// IsCOFINS verifica se o imposto é COFINS (ID 5)
func (d *DocumentoVendaItemImposto) IsCOFINS() bool {
	return d.ImpostoID == 5
}

// GetTipoImposto retorna o nome do tipo de imposto
func (d *DocumentoVendaItemImposto) GetTipoImposto() string {
	switch d.ImpostoID {
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

// CalcularValorBase calcula o valor base do imposto
// Se houver redução de base, aplica o percentual de redução
func (d *DocumentoVendaItemImposto) CalcularValorBase(valorProduto float64) float64 {
	base := valorProduto
	if d.ReducaoBase > 0 {
		base = base * (1 - d.ReducaoBase/100)
	}
	return base
}

// CalcularValorComAliquotaMVA calcula o valor com MVA aplicada
func (d *DocumentoVendaItemImposto) CalcularValorComAliquotaMVA(valorBase float64) float64 {
	if d.HasAliquotaMVA() {
		return valorBase * (1 + *d.AliquotaMVA/100)
	}
	return valorBase
}

// CalcularValorST calcula o valor do ICMS ST
// Base ST = Base ICMS * (1 + MVA/100)
// Valor ST = (Base ST * Aliquota ST) - (Base ICMS * Aliquota ICMS)
func (d *DocumentoVendaItemImposto) CalcularValorST(valorBaseICMS float64, aliquotaICMS float64) float64 {
	if !d.IsICMSST() {
		return 0
	}

	baseST := d.CalcularValorComAliquotaMVA(valorBaseICMS)
	valorST := baseST * (d.Aliquota / 100)
	valorICMS := valorBaseICMS * (aliquotaICMS / 100)

	return valorST - valorICMS
}
