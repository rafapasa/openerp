package models

import (
	"time"
)

// ============================================================
// CONSTANTES DO PLANO DE CONTAS
// ============================================================
// (Adicionar no final do constants.go)

/*
const (
    PlanoContasTipoBanco  = 1
    PlanoContasTipoOutro  = 9

    PlanoContasEspecieSintetica = 1
    PlanoContasEspecieAnalitica = 2

    PlanoContasSituacaoAtivo    = 1
    PlanoContasSituacaoInativo  = 2
    PlanoContasSituacaoBloqueado = 3
)
*/

// ============================================================
// MODEL: PlanoContasFinanceiro
// ============================================================

type PlanoContasFinanceiro struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID              int    `gorm:"column:pcf_id;primaryKey;autoIncrement" json:"id"`
	EmpresaFilialID int    `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	PlanoPaiID      *int   `gorm:"column:pcf_pcf_id" json:"plano_pai_id,omitempty"`
	ContaBancariaID *int   `gorm:"column:ctb_id" json:"conta_bancaria_id,omitempty"`
	Descricao       string `gorm:"column:pcf_descricao;type:varchar(255);not null" json:"descricao"`
	Tipo            int    `gorm:"column:pcf_tipo;not null" json:"tipo"`                   // 1-Banco, 9-Outro
	Especie         int    `gorm:"column:pcf_especie;not null" json:"especie"`             // 1-Sintética, 2-Analítica
	Situacao        int    `gorm:"column:pcf_situacao;not null;default:1" json:"situacao"` // 1-Ativo, 2-Inativo, 3-Bloqueado
	Mascara         string `gorm:"column:pcf_mascara;type:varchar(100);not null" json:"mascara"`

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
	EmpresaFilial *EmpresaFilial          `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
	PlanoPai      *PlanoContasFinanceiro  `gorm:"foreignKey:PlanoPaiID;references:pcf_id" json:"plano_pai,omitempty"`
	ContaBancaria *ContaBancaria          `gorm:"foreignKey:ContaBancariaID;references:ctb_id" json:"conta_bancaria,omitempty"`
	PlanosFilhos  []PlanoContasFinanceiro `gorm:"foreignKey:PlanoPaiID;references:ID" json:"planos_filhos,omitempty"`
}

func (PlanoContasFinanceiro) TableName() string {
	return "plano_contas_financeiro"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *PlanoContasFinanceiro) IsActive() bool {
	return m.Situacao == 1
}

func (m *PlanoContasFinanceiro) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *PlanoContasFinanceiro) IsSintetica() bool {
	return m.Especie == 1
}

func (m *PlanoContasFinanceiro) IsAnalitica() bool {
	return m.Especie == 2
}

func (m *PlanoContasFinanceiro) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
	m.Situacao = 2
}
