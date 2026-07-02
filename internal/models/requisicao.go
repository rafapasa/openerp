package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/constants"
)

// ============================================================
// MODEL: Requisicao
// ============================================================

type Requisicao struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID               int  `gorm:"column:req_id;primaryKey;autoIncrement" json:"id"`
	OrdemProducaoID  *int `gorm:"column:orp_id" json:"ordem_producao_id,omitempty"`
	EmpresaFilialID  int  `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	RotinaContabilID *int `gorm:"column:roc_id" json:"rotina_contabil_id,omitempty"`

	TipoOperacao  int              `gorm:"column:req_tipooperacao;not null" json:"tipo_operacao"` // 1 - entrada, 2 - saída, etc
	DataMovimento time.Time        `gorm:"column:req_datamovimento;type:date;not null" json:"data_movimento"`
	DataSistema   time.Time        `gorm:"column:req_data_sistema;type:date;not null" json:"data_sistema"`
	Observacao    *string          `gorm:"column:req_observacao;type:text" json:"observacao,omitempty"`
	Situacao      constants.Status `gorm:"column:req_situacao;not null" json:"situacao"`

	TipoRequisicaoID *int `gorm:"column:treq_id" json:"tipo_requisicao_id,omitempty"`
	TalhaoID         *int `gorm:"column:tat_id" json:"talhao_id,omitempty"`

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
	OrdemProducao  *OrdemProducao   `gorm:"foreignKey:OrdemProducaoID;references:orp_id" json:"ordem_producao,omitempty"`
	EmpresaFilial  *EmpresaFilial   `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
	RotinaContabil *RotinaContabil  `gorm:"foreignKey:RotinaContabilID;references:roc_id" json:"rotina_contabil,omitempty"`
	TipoRequisicao *TipoRequisicao  `gorm:"foreignKey:TipoRequisicaoID;references:treq_id" json:"tipo_requisicao,omitempty"`
	Talhao         *TalhaoTerra     `gorm:"foreignKey:TalhaoID;references:tat_id" json:"talhao,omitempty"`
	Itens          []RequisicaoItem `gorm:"foreignKey:RequisicaoID;references:req_id" json:"itens,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (Requisicao) TableName() string {
	return "requisicao"
}

func (r *Requisicao) BeforeCreate(tx *gorm.DB) error {
	if r.CreatedBy == nil {
		r.CreatedBy = new(int)
		*r.CreatedBy = 0
	}
	if r.UpdatedBy == nil {
		r.UpdatedBy = new(int)
		*r.UpdatedBy = 0
	}
	return nil
}

func (r *Requisicao) BeforeUpdate(tx *gorm.DB) error {
	if r.UpdatedBy == nil {
		r.UpdatedBy = new(int)
		*r.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsActive verifica se a requisição está ativa
func (r *Requisicao) IsActive() bool {
	return r.Situacao == constants.StatusAtivo
}

// IsDeleted verifica se a requisição foi deletada logicamente
func (r *Requisicao) IsDeleted() bool {
	return r.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (r *Requisicao) SoftDelete() {
	now := time.Now()
	r.DeletedAt = &now
	r.Situacao = constants.StatusInativo
}

// IsInactive verifica se a requisição está inativa
func (r *Requisicao) IsInactive() bool {
	return r.Situacao == constants.StatusInativo
}

// IsEntrada verifica se é uma requisição de entrada
func (r *Requisicao) IsEntrada() bool {
	return r.TipoOperacao == 1
}

// IsSaida verifica se é uma requisição de saída
func (r *Requisicao) IsSaida() bool {
	return r.TipoOperacao == 2
}

// IsAjuste verifica se é uma requisição de ajuste
func (r *Requisicao) IsAjuste() bool {
	return r.TipoOperacao == 3
}

// HasItens verifica se a requisição possui itens
func (r *Requisicao) HasItens() bool {
	return len(r.Itens) > 0
}

// GetTotalItens retorna a quantidade total de itens
func (r *Requisicao) GetTotalItens() int {
	return len(r.Itens)
}

// GetQuantidadeTotal retorna a soma das quantidades dos itens
func (r *Requisicao) GetQuantidadeTotal() float64 {
	total := 0.0
	for _, item := range r.Itens {
		total += item.Quantidade
	}
	return total
}
