package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: TabelaPreco
// ============================================================

type TabelaPreco struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID                int        `gorm:"column:tbp_id;primaryKey;autoIncrement" json:"id"`
	Descricao         string     `gorm:"column:tbp_descricao;type:varchar(100);not null" json:"descricao"`
	DataInicio        time.Time  `gorm:"column:tbp_datainicio;type:date;not null" json:"data_inicio"`
	DataFim           *time.Time `gorm:"column:tbp_datafim;type:date" json:"data_fim,omitempty"`
	Tipo              *int       `gorm:"column:tbp_tipo" json:"tipo,omitempty"`
	Percentual        *float64   `gorm:"column:tbp_percentual;type:decimal(15,4)" json:"percentual,omitempty"`
	TipoServico       *int       `gorm:"column:tbp_tipo_servico" json:"tipo_servico,omitempty"`
	PercentualServico *float64   `gorm:"column:tbp_percentual_servico;type:decimal(15,4)" json:"percentual_servico,omitempty"`

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
	Produtos []TabelaPrecoProduto `gorm:"foreignKey:TabelaPrecoID;references:ID" json:"produtos,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (TabelaPreco) TableName() string {
	return "tabela_preco"
}

func (m *TabelaPreco) BeforeCreate(tx *gorm.DB) error {
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

func (m *TabelaPreco) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *TabelaPreco) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *TabelaPreco) IsActive() bool {
	now := time.Now()
	if m.DataFim == nil {
		return m.DataInicio.Before(now) || m.DataInicio.Equal(now)
	}
	return (m.DataInicio.Before(now) || m.DataInicio.Equal(now)) &&
		(m.DataFim.After(now) || m.DataFim.Equal(now))
}

func (m *TabelaPreco) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}

func (m *TabelaPreco) HasProdutos() bool {
	return len(m.Produtos) > 0
}

func (m *TabelaPreco) GetProdutosAtivos() []TabelaPrecoProduto {
	var ativos []TabelaPrecoProduto
	for _, p := range m.Produtos {
		if p.IsActive() {
			ativos = append(ativos, p)
		}
	}
	return ativos
}
