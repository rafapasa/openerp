package models

import (
	"time"

	"github.com/openerp/backend/internal/constants"
)

// ============================================================
// MODEL: CentroDeCusto
// ============================================================

type CentroDeCusto struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID              int              `gorm:"column:cdc_id;primaryKey;autoIncrement" json:"id"`
	EmpresaFilialID int              `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	Descricao       string           `gorm:"column:cdc_descricao;type:varchar(100);not null" json:"descricao"`
	Situacao        constants.Status `gorm:"column:cdc_situacao;not null" json:"situacao"`

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
	EmpresaFilial              *EmpresaFilial         `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
	Despesas                   []DespesaCentroDeCusto `gorm:"foreignKey:CentroDeCustoID;references:cdc_id" json:"despesas,omitempty"`
	TituloCentroCusto          []TituloCentroDeCusto  `gorm:"foreignKey:CentroDeCustoID;references:cdc_id" json:"titulo_centro_custo,omitempty"`
	RotinaContabilItensDebito  []RotinaContabilItem   `gorm:"foreignKey:CentroDeCustoDebitoID;references:cdc_id" json:"rotina_contabil_itens_debito,omitempty"`
	RotinaContabilItensCredito []RotinaContabilItem   `gorm:"foreignKey:CentroDeCustoCreditoID;references:cdc_id" json:"rotina_contabil_itens_credito,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (CentroDeCusto) TableName() string {
	return "centro_de_custo"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsActive verifica se o centro de custo está ativo
func (c *CentroDeCusto) IsActive() bool {
	return c.Situacao == constants.StatusAtivo
}

// IsDeleted verifica se o centro de custo foi deletado logicamente
func (c *CentroDeCusto) IsDeleted() bool {
	return c.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (c *CentroDeCusto) SoftDelete() {
	now := time.Now()
	c.DeletedAt = &now
	c.Situacao = constants.StatusInativo
}

// IsInactive verifica se o centro de custo está inativo
func (c *CentroDeCusto) IsInactive() bool {
	return c.Situacao == constants.StatusInativo
}

// GetNomeCompleto retorna o nome completo do centro de custo
func (c *CentroDeCusto) GetNomeCompleto() string {
	return c.Descricao + " (ID: " + string(rune(c.ID)) + ")"
}

// HasDespesas verifica se o centro de custo possui despesas associadas
func (c *CentroDeCusto) HasDespesas() bool {
	return len(c.Despesas) > 0
}

// SafeToDelete verifica se o centro de custo pode ser excluído
func (c *CentroDeCusto) SafeToDelete() bool {
	// Verifica se há registros relacionados
	if c.HasDespesas() {
		return false
	}
	if len(c.TituloCentroCusto) > 0 {
		return false
	}
	if len(c.RotinaContabilItensDebito) > 0 {
		return false
	}
	if len(c.RotinaContabilItensCredito) > 0 {
		return false
	}
	return true
}

// GetTotalDespesas retorna a soma dos percentuais das despesas associadas
func (c *CentroDeCusto) GetTotalDespesas() float64 {
	total := 0.0
	for _, d := range c.Despesas {
		total += d.Percentual
	}
	return total
}
