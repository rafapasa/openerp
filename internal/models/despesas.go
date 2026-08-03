package models

import (
	"time"

	"github.com/openerp/backend/internal/constants"
)

// ============================================================
// CONSTANTES
// ============================================================

const (
	// Tipo de despesa
	DespesaFixa     = 1
	DespesaVariavel = 2

	// Espécie da despesa
	DespesaEspecieSintetica = 1
	DespesaEspecieAnalitica = 2
)

// ============================================================
// MODEL: Despesas
// ============================================================

type Despesa struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID            int              `gorm:"column:desp_id;primaryKey;autoIncrement" json:"id"`
	PlanoContasID int              `gorm:"column:pcf_id;not null" json:"plano_contas_id"`
	Descricao     string           `gorm:"column:desp_descricao;type:varchar(255);not null" json:"descricao"`
	Situacao      constants.Status `gorm:"column:desp_situacao;not null;default:1" json:"situacao"`
	FixaVariavel  int              `gorm:"column:desp_fixa_variavel;not null" json:"fixa_variavel"` // 1-Fixa, 2-Variável
	Mascara       *string          `gorm:"column:desp_mascara;type:varchar(20)" json:"mascara,omitempty"`
	Especie       *int             `gorm:"column:desp_especie" json:"especie,omitempty"` // 1-Sintética, 2-Analítica

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
	PlanoContas    *PlanoContasFinanceiro    `gorm:"foreignKey:PlanoContasID;references:pcf_id" json:"plano_contas,omitempty"`
	CentroDeCustos []DespesaCentroDeCusto    `gorm:"foreignKey:DespesaID;references:ID" json:"centro_de_custos,omitempty"`
	Lancamentos    []LancamentoDespesa       `gorm:"foreignKey:DespesaID;references:ID" json:"lancamentos,omitempty"`
	Produtos       []ProdutoFichaOperacional `gorm:"foreignKey:DespesaID;references:ID" json:"produtos,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (Despesa) TableName() string {
	return "despesas"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsActive verifica se a despesa está ativa
func (d *Despesa) IsActive() bool {
	return d.Situacao == constants.StatusAtivo
}

// IsDeleted verifica se a despesa foi deletada logicamente
func (d *Despesa) IsDeleted() bool {
	return d.DeletedAt != nil
}

// IsFixa verifica se a despesa é fixa
func (d *Despesa) IsFixa() bool {
	return d.FixaVariavel == DespesaFixa
}

// IsVariavel verifica se a despesa é variável
func (d *Despesa) IsVariavel() bool {
	return d.FixaVariavel == DespesaVariavel
}

// IsSintetica verifica se a despesa é sintética
func (d *Despesa) IsSintetica() bool {
	return d.Especie != nil && *d.Especie == DespesaEspecieSintetica
}

// IsAnalitica verifica se a despesa é analítica
func (d *Despesa) IsAnalitica() bool {
	return d.Especie != nil && *d.Especie == DespesaEspecieAnalitica
}

// SoftDelete realiza a exclusão lógica
func (d *Despesa) SoftDelete() {
	now := time.Now()
	d.DeletedAt = &now
	d.Situacao = constants.StatusInativo
}
