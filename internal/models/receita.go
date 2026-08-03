package models

import (
	"time"

	"github.com/openerp/backend/internal/constants"
)

// ============================================================
// CONSTANTES
// ============================================================

// Espécie da receita
const (
	ReceitaEspecieSintetica = 1
	ReceitaEspecieAnalitica = 2
)

// ============================================================
// MODEL: Receitas
// ============================================================

type Receita struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID            int              `gorm:"column:rec_id;primaryKey;autoIncrement" json:"id"`
	PlanoContasID int              `gorm:"column:pcf_id;not null" json:"plano_contas_id"`
	Descricao     string           `gorm:"column:rec_descricao;type:varchar(255);not null" json:"descricao"`
	Situacao      constants.Status `gorm:"column:rec_situacao;not null;default:1" json:"situacao"`
	Mascara       *string          `gorm:"column:rec_mascara;type:varchar(20)" json:"mascara,omitempty"`
	Especie       *int             `gorm:"column:rec_especie" json:"especie,omitempty"` // 1-Sintética, 2-Analítica

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
	PlanoContas  *PlanoContasFinanceiro `gorm:"foreignKey:PlanoContasID;references:pcf_id" json:"plano_contas,omitempty"`
	Lancamentos  []LancamentoReceita    `gorm:"foreignKey:ReceitaID;references:ID" json:"lancamentos,omitempty"`
	Titulos      []Titulo               `gorm:"foreignKey:ReceitaID;references:ID" json:"titulos,omitempty"`
	Contratos    []ContratoEntidadeTipo `gorm:"foreignKey:ReceitaID;references:ID" json:"contratos,omitempty"`
	ProdutosTipo []ProdutoTipoProduto   `gorm:"foreignKey:ReceitaID;references:ID" json:"produtos_tipo,omitempty"`
	TituloBaixas []TituloBaixa          `gorm:"foreignKey:ReceitaID;references:ID" json:"titulo_baixas,omitempty"`
	Processos    []Processo             `gorm:"foreignKey:ReceitaID;references:ID" json:"processos,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (Receita) TableName() string {
	return "receitas"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsActive verifica se a receita está ativa
func (r *Receita) IsActive() bool {
	return r.Situacao == constants.StatusAtivo
}

// IsDeleted verifica se a receita foi deletada logicamente
func (r *Receita) IsDeleted() bool {
	return r.DeletedAt != nil
}

// IsSintetica verifica se a receita é sintética
func (r *Receita) IsSintetica() bool {
	return r.Especie != nil && *r.Especie == ReceitaEspecieSintetica
}

// IsAnalitica verifica se a receita é analítica
func (r *Receita) IsAnalitica() bool {
	return r.Especie != nil && *r.Especie == ReceitaEspecieAnalitica
}

// SoftDelete realiza a exclusão lógica
func (r *Receita) SoftDelete() {
	now := time.Now()
	r.DeletedAt = &now
	r.Situacao = constants.StatusInativo
}

// HasMascara verifica se a receita tem máscara definida
func (r *Receita) HasMascara() bool {
	return r.Mascara != nil && *r.Mascara != ""
}
