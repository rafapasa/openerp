package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: TipoProduto
// ============================================================

type TipoProduto struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID                int     `gorm:"column:ptp_id;primaryKey;autoIncrement" json:"id"`
	Codigo            string  `gorm:"column:tpp_codigo;type:varchar(2);not null" json:"codigo"`
	Descricao         string  `gorm:"column:ptp_descricao;type:varchar(60);not null" json:"descricao"`
	ValidarFornecedor int     `gorm:"column:ptp_validarfornecedor;not null" json:"validar_fornecedor"`
	MovimentaEstoque  int     `gorm:"column:ptp_movimentaestoque;not null" json:"movimenta_estoque"`
	Sigla             *string `gorm:"column:ptp_sigla;type:varchar(5)" json:"sigla,omitempty"`
	ReceitaID         *int    `gorm:"column:rec_id" json:"receita_id,omitempty"`
	DespesaID         *int    `gorm:"column:desp_id" json:"despesa_id,omitempty"`
	ProdutoPacote     int     `gorm:"column:ptp_produtopacote;not null;default:0" json:"produto_pacote"`
	Combustivel       int     `gorm:"column:ptp_combustivel;not null;default:0" json:"combustivel"`

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
	Produtos []Produto `gorm:"foreignKey:TipoProdutoID;references:ID" json:"produtos,omitempty"`
	// Receita  *Receita  `gorm:"foreignKey:ReceitaID;references:rec_id" json:"receita,omitempty"`
	// Despesa  *Despesa  `gorm:"foreignKey:DespesaID;references:desp_id" json:"despesa,omitempty"`
}

func (TipoProduto) TableName() string {
	return "produto_tipoproduto"
}

func (m *TipoProduto) BeforeCreate(tx *gorm.DB) error {
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

func (m *TipoProduto) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *TipoProduto) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *TipoProduto) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}

func (m *TipoProduto) IsProdutoPacote() bool {
	return m.ProdutoPacote == 1
}

func (m *TipoProduto) IsCombustivel() bool {
	return m.Combustivel == 1
}

func (m *TipoProduto) MovimentaEstoqueSim() bool {
	return m.MovimentaEstoque == 1
}
