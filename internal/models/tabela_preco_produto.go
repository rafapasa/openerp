package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: TabelaPrecoProduto
// ============================================================

type TabelaPrecoProduto struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	TabelaPrecoID       int      `gorm:"column:tbp_id;primaryKey;uniqueIndex:uk_tbp_pro,priority:1" json:"tabela_preco_id"`
	Item                int      `gorm:"column:tbpp_item;primaryKey" json:"item"`
	ProdutoID           int      `gorm:"column:pro_id;uniqueIndex:uk_tbp_pro,priority:2" json:"produto_id"`
	ValorMinimo         float64  `gorm:"column:tbpp_valorminimo;type:decimal(15,4);not null" json:"valor_minimo"`
	ValorPadrao         float64  `gorm:"column:tbpp_valorpadrao;type:decimal(15,4);not null" json:"valor_padrao"`
	Situacao            int      `gorm:"column:tbpp_situacao;not null;default:1" json:"situacao"`
	QuantAtacado        *float64 `gorm:"column:tbpp_quantatacado;type:decimal(15,4)" json:"quant_atacado,omitempty"`
	PercDescontoAtacado *float64 `gorm:"column:tbpp_percdescatacado;type:decimal(15,4)" json:"perc_desconto_atacado,omitempty"`
	ValorCusto          *float64 `gorm:"column:tbpp_valorcusto;type:decimal(15,4)" json:"valor_custo,omitempty"`
	MargemLucro         *float64 `gorm:"column:tbpp_margemlucro;type:decimal(15,2)" json:"margem_lucro,omitempty"`

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
	TabelaPreco *TabelaPreco `gorm:"foreignKey:TabelaPrecoID;references:ID" json:"tabela_preco,omitempty"`
	Produto     *Produto     `gorm:"foreignKey:ProdutoID;references:pro_id" json:"produto,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (TabelaPrecoProduto) TableName() string {
	return "tabela_preco_produto"
}

func (m *TabelaPrecoProduto) BeforeCreate(tx *gorm.DB) error {
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

func (m *TabelaPrecoProduto) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *TabelaPrecoProduto) IsActive() bool {
	return m.Situacao == 1
}

func (m *TabelaPrecoProduto) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *TabelaPrecoProduto) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
	m.Situacao = 2
}

func (m *TabelaPrecoProduto) HasAtacado() bool {
	return m.QuantAtacado != nil && *m.QuantAtacado > 0
}

func (m *TabelaPrecoProduto) GetValorComDescontoAtacado() float64 {
	if !m.HasAtacado() || m.PercDescontoAtacado == nil {
		return m.ValorPadrao
	}
	return m.ValorPadrao * (1 - *m.PercDescontoAtacado/100)
}

func (m *TabelaPrecoProduto) GetMargemLucro() float64 {
	if m.MargemLucro != nil {
		return *m.MargemLucro
	}
	if m.ValorCusto != nil && *m.ValorCusto > 0 {
		return ((m.ValorPadrao - *m.ValorCusto) / m.ValorPadrao) * 100
	}
	return 0
}
