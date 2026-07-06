package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: ProdutoEntidade
// ============================================================

type ProdutoEntidade struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ProdutoID        int      `gorm:"column:pro_id;primaryKey" json:"produto_id"`
	EntidadeID       int      `gorm:"column:ent_id;primaryKey" json:"entidade_id"`
	Codigo           string   `gorm:"column:proent_codigo;type:varchar(50);primaryKey" json:"codigo"`
	OperacaoFiscalID *int     `gorm:"column:opf_id" json:"operacao_fiscal_id,omitempty"`
	CSTIPIID         *int     `gorm:"column:cstipi_id" json:"cst_ipi_id,omitempty"`
	CSTICMSID        *int     `gorm:"column:csticms_id" json:"cst_icms_id,omitempty"`
	CSTPISCOFINSID   *int     `gorm:"column:cstpiscofins_id" json:"cst_pis_cofins_id,omitempty"`
	Descricao        *string  `gorm:"column:proent_descricao;type:varchar(1000)" json:"descricao,omitempty"`
	UndCompraID      *int     `gorm:"column:und_id_compra" json:"und_compra_id,omitempty"`
	UndVendaID       *int     `gorm:"column:und_id_venda" json:"und_venda_id,omitempty"`
	QuantCompra      *float64 `gorm:"column:proent_quant_compra;type:decimal(15,5)" json:"quant_compra,omitempty"`
	QuantEstoque     *float64 `gorm:"column:proent_quant_estoque;type:decimal(15,5)" json:"quant_estoque,omitempty"`
	MargemLucro      *float64 `gorm:"column:proent_margemlucro;type:decimal(15,2)" json:"margem_lucro,omitempty"`

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
	Produto        *Produto        `gorm:"foreignKey:ProdutoID;references:pro_id" json:"produto,omitempty"`
	Entidade       *Entidade       `gorm:"foreignKey:EntidadeID;references:ent_id" json:"entidade,omitempty"`
	OperacaoFiscal *OperacaoFiscal `gorm:"foreignKey:OperacaoFiscalID;references:opf_id" json:"operacao_fiscal,omitempty"`
	CSTIPI         *CSTIPI         `gorm:"foreignKey:CSTIPIID;references:cstipi_id" json:"cst_ipi,omitempty"`
	CSTICMS        *CSTICMS        `gorm:"foreignKey:CSTICMSID;references:csticms_id" json:"cst_icms,omitempty"`
	CSTPISCOFINS   *CSTPISCOFINS   `gorm:"foreignKey:CSTPISCOFINSID;references:cstpiscofins_id" json:"cst_pis_cofins,omitempty"`
}

func (ProdutoEntidade) TableName() string {
	return "produto_entidade"
}

func (m *ProdutoEntidade) BeforeCreate(tx *gorm.DB) error {
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

func (m *ProdutoEntidade) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

func (m *ProdutoEntidade) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *ProdutoEntidade) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}
