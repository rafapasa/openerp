package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: NotaFiscalItem
// ============================================================

type NotaFiscalItem struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	NotaFiscalID     int  `gorm:"column:ntf_id;primaryKey" json:"nota_fiscal_id"`
	Item             int  `gorm:"column:nfi_item;primaryKey" json:"item"`
	CSTICMSID        int  `gorm:"column:csticms_id;not null" json:"cst_icms_id"`
	ProdutoID        int  `gorm:"column:pro_id;not null" json:"produto_id"`
	CSTPISCOFINSID   int  `gorm:"column:cstpiscofins_id;not null" json:"cst_pis_cofins_id"`
	OperacaoFiscalID int  `gorm:"column:opf_id;not null" json:"operacao_fiscal_id"`
	CSTIPIID         int  `gorm:"column:cstipi_id;not null" json:"cst_ipi_id"`
	RotinaContabilID *int `gorm:"column:roc_id" json:"rotina_contabil_id,omitempty"`

	Quantidade     float64  `gorm:"column:nfi_quantidade;type:decimal(15,5);not null" json:"quantidade"`
	ValorUnitario  float64  `gorm:"column:nfi_valorunitario;type:decimal(15,4);not null" json:"valor_unitario"`
	ValorFrete     float64  `gorm:"column:nfi_valorfrete;type:decimal(15,2);not null" json:"valor_frete"`
	ValorSeguro    *float64 `gorm:"column:nfi_valorseguro;type:decimal(15,2)" json:"valor_seguro,omitempty"`
	ValorDesconto  *float64 `gorm:"column:nfi_valordesconto;type:decimal(15,2)" json:"valor_desconto,omitempty"`
	IndTotalizacao int      `gorm:"column:nfi_indtotalizacao;not null" json:"ind_totalizacao"` // 0 - compõe, 1 - não compõe
	ValorOutras    *float64 `gorm:"column:nfi_valoroutras;type:decimal(15,2)" json:"valor_outras,omitempty"`
	MovFisica      int      `gorm:"column:nfi_movfisica;not null" json:"mov_fisica"`
	TotalProduto   *float64 `gorm:"column:nfi_totalproduto;type:decimal(15,4)" json:"total_produto"`
	TotalItem      *float64 `gorm:"column:nfi_totalitem;type:decimal(15,4)" json:"total_item"`
	Descricao      *string  `gorm:"column:nfi_descricao;type:varchar(255)" json:"descricao,omitempty"`

	GradeID   *int `gorm:"column:grade_id" json:"grade_id,omitempty"`
	GradeItem *int `gorm:"column:grai_item" json:"grade_item,omitempty"`

	QuantidadeTrib float64  `gorm:"column:nfi_quantidade_trib;type:decimal(15,5);not null" json:"quantidade_trib"`
	PesoBruto      *float64 `gorm:"column:nfi_pesobruto;type:decimal(15,4)" json:"peso_bruto,omitempty"`
	PesoLiquido    *float64 `gorm:"column:nfi_pesoliquido;type:decimal(15,4)" json:"peso_liquido,omitempty"`

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
	NotaFiscal     *NotaFiscal     `gorm:"foreignKey:NotaFiscalID;references:ntf_id" json:"nota_fiscal,omitempty"`
	CSTICMS        *CSTICMS        `gorm:"foreignKey:CSTICMSID;references:csticms_id" json:"cst_icms,omitempty"`
	Produto        *Produto        `gorm:"foreignKey:ProdutoID;references:pro_id" json:"produto,omitempty"`
	CSTPISCOFINS   *CSTPISCOFINS   `gorm:"foreignKey:CSTPISCOFINSID;references:cstpiscofins_id" json:"cst_pis_cofins,omitempty"`
	OperacaoFiscal *OperacaoFiscal `gorm:"foreignKey:OperacaoFiscalID;references:opf_id" json:"operacao_fiscal,omitempty"`
	CSTIPI         *CSTIPI         `gorm:"foreignKey:CSTIPIID;references:cstipi_id" json:"cst_ipi,omitempty"`
	RotinaContabil *RotinaContabil `gorm:"foreignKey:RotinaContabilID;references:roc_id" json:"rotina_contabil,omitempty"`
	Grade          *GradeProduto   `gorm:"foreignKey:GradeID;references:grade_id" json:"grade,omitempty"`

	Impostos []NotaFiscalItemImposto `gorm:"foreignKey:NotaFiscalID,NotaFiscalItem;references:ntf_id,nfi_item" json:"impostos,omitempty"`
	Grades   []NotaFiscalItemGrade   `gorm:"foreignKey:NotaFiscalID,NotaFiscalItem;references:ntf_id,nfi_item" json:"grades,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (NotaFiscalItem) TableName() string {
	return "nota_fiscal_item"
}

func (n *NotaFiscalItem) BeforeCreate(tx *gorm.DB) error {
	if n.CreatedBy == nil {
		n.CreatedBy = new(int)
		*n.CreatedBy = 0
	}
	if n.UpdatedBy == nil {
		n.UpdatedBy = new(int)
		*n.UpdatedBy = 0
	}
	return nil
}

func (n *NotaFiscalItem) BeforeUpdate(tx *gorm.DB) error {
	if n.UpdatedBy == nil {
		n.UpdatedBy = new(int)
		*n.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o item foi deletado logicamente
func (n *NotaFiscalItem) IsDeleted() bool {
	return n.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (n *NotaFiscalItem) SoftDelete() {
	now := time.Now()
	n.DeletedAt = &now
}

// HasGrade verifica se o item possui grade
func (n *NotaFiscalItem) HasGrade() bool {
	return n.GradeID != nil && n.GradeItem != nil
}

// GetTotalProdutoCalculado recalcula o total do produto
func (n *NotaFiscalItem) GetTotalProdutoCalculado() float64 {
	total := n.Quantidade * n.ValorUnitario
	if n.ValorFrete > 0 {
		total += n.ValorFrete
	}
	if n.ValorSeguro != nil {
		total += *n.ValorSeguro
	}
	if n.ValorOutras != nil {
		total += *n.ValorOutras
	}
	if n.ValorDesconto != nil {
		total -= *n.ValorDesconto
	}
	return total
}

// CompoeTotal verifica se o item compõe o total da nota
func (n *NotaFiscalItem) CompoeTotal() bool {
	return n.IndTotalizacao == 0
}

// HasImpostos verifica se o item possui impostos
func (n *NotaFiscalItem) HasImpostos() bool {
	return len(n.Impostos) > 0
}

// HasGrades verifica se o item possui grades
func (n *NotaFiscalItem) HasGrades() bool {
	return len(n.Grades) > 0
}
