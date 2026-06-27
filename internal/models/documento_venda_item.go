package models

import (
	"time"
)

// ============================================================
// MODEL: DocumentoVendaItem
// ============================================================

type DocumentoVendaItem struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID               int  `gorm:"column:dvi_item;primaryKey;autoIncrement" json:"id"`
	DocumentoVendaID int  `gorm:"column:ddv_id;not null" json:"documento_venda_id"`
	ProdutoID        int  `gorm:"column:pro_id;not null" json:"produto_id"`
	OperacaoFiscalID int  `gorm:"column:opf_id;not null" json:"operacao_fiscal_id"`
	CSTICMSID        int  `gorm:"column:csticms_id;not null" json:"cst_icms_id"`
	CSTIPIID         int  `gorm:"column:cstipi_id;not null" json:"cst_ipi_id"`
	CSTPISCOFINSID   int  `gorm:"column:cstpiscofins_id;not null" json:"cst_pis_cofins_id"`
	RotinaContabilID *int `gorm:"column:roc_id" json:"rotina_contabil_id,omitempty"`
	RepresentanteID  *int `gorm:"column:rep_ent_id" json:"representante_id,omitempty"`
	UsuarioID        *int `gorm:"column:usu_id" json:"usuario_id,omitempty"`
	GradeID          *int `gorm:"column:grade_id" json:"grade_id,omitempty"`
	GradeItem        *int `gorm:"column:grai_item" json:"grade_item,omitempty"`

	// ============================================================
	// CAMPOS DE QUANTIDADE E VALORES
	// ============================================================
	Quantidade         float64  `gorm:"column:dvi_quantidade;type:decimal(15,4);not null" json:"quantidade"`
	QuantidadeBruta    *float64 `gorm:"column:dvi_quantbruta;type:decimal(15,4)" json:"quantidade_bruta,omitempty"`
	QuantidadeTara     *float64 `gorm:"column:dvi_quanttara;type:decimal(15,4)" json:"quantidade_tara,omitempty"`
	ValorUnitario      float64  `gorm:"column:dvi_valorunitario;type:decimal(15,6);not null" json:"valor_unitario"`
	PercentualDesconto *float64 `gorm:"column:dvi_percentualdesconto;type:decimal(5,2)" json:"percentual_desconto,omitempty"`
	ValorDesconto      *float64 `gorm:"column:dvi_valordesconto;type:decimal(15,2)" json:"valor_desconto,omitempty"`
	ValorFrete         *float64 `gorm:"column:dvi_valorfrete;type:decimal(15,2)" json:"valor_frete,omitempty"`
	ValorAdicional     float64  `gorm:"column:dvi_valoradicional;type:decimal(15,4);not null;default:0.0000" json:"valor_adicional"`

	// ============================================================
	// CAMPOS DE TOTAIS
	// ============================================================
	TotalProdutos float64 `gorm:"column:dvi_totalprodutos;type:decimal(15,2);not null" json:"total_produtos"`
	TotalItem     float64 `gorm:"column:dvi_totalitem;type:decimal(15,2);not null" json:"total_item"`
	PesoBruto     float64 `gorm:"column:dvi_pesobruto;type:decimal(15,4);not null" json:"peso_bruto"`
	PesoLiquido   float64 `gorm:"column:dvi_pesoliquido;type:decimal(15,4);not null" json:"peso_liquido"`
	CustoUnitario float64 `gorm:"column:dvi_custounitario;type:decimal(15,4);not null;default:0.0000" json:"custo_unitario"`
	CustoTotal    float64 `gorm:"column:dvi_custototal;type:decimal(15,4);not null;default:0.0000" json:"custo_total"`

	// ============================================================
	// CAMPOS DE DESCRIÇÃO E CONTROLE
	// ============================================================
	DescricaoProduto *string    `gorm:"column:dvi_pro_descricao;type:varchar(255)" json:"descricao_produto,omitempty"`
	Impresso         *int       `gorm:"column:dvi_impresso;default:1" json:"impresso,omitempty"`
	HoraLancamento   *time.Time `gorm:"column:dvi_hora_lancamento;type:time" json:"hora_lancamento,omitempty"`

	// ============================================================
	// CAMPOS PARA COMODATO (empréstimo)
	// ============================================================
	ComodatoDocumentoID *int `gorm:"column:dvi_comodato_ddv_id" json:"comodato_documento_id,omitempty"`
	ComodatoItemID      *int `gorm:"column:dvi_comodato_dvi_item" json:"comodato_item_id,omitempty"`

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
	DocumentoVenda *DocumentoVenda `gorm:"foreignKey:DocumentoVendaID;references:ID" json:"documento_venda,omitempty"`
	Produto        *Produto        `gorm:"foreignKey:ProdutoID;references:pro_id" json:"produto,omitempty"`
	OperacaoFiscal *OperacaoFiscal `gorm:"foreignKey:OperacaoFiscalID;references:opf_id" json:"operacao_fiscal,omitempty"`
	CSTICMS        *CSTICMS        `gorm:"foreignKey:CSTICMSID;references:csticms_id" json:"cst_icms,omitempty"`
	CSTIPI         *CSTIPI         `gorm:"foreignKey:CSTIPIID;references:cstipi_id" json:"cst_ipi,omitempty"`
	CSTPISCOFINS   *CSTPISCOFINS   `gorm:"foreignKey:CSTPISCOFINSID;references:cstpiscofins_id" json:"cst_pis_cofins,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (DocumentoVendaItem) TableName() string {
	return "documento_venda_item"
}

func (d *DocumentoVendaItem) BeforeCreate() error {
	if d.CreatedBy == nil {
		d.CreatedBy = new(int)
		*d.CreatedBy = 0
	}
	if d.UpdatedBy == nil {
		d.UpdatedBy = new(int)
		*d.UpdatedBy = 0
	}
	return nil
}

func (d *DocumentoVendaItem) BeforeUpdate() error {
	if d.UpdatedBy == nil {
		d.UpdatedBy = new(int)
		*d.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// GetTotalBruto retorna o total bruto (quantidade * valor unitário)
func (d *DocumentoVendaItem) GetTotalBruto() float64 {
	return d.Quantidade * d.ValorUnitario
}

// GetTotalComDesconto retorna o total com desconto aplicado
func (d *DocumentoVendaItem) GetTotalComDesconto() float64 {
	if d.ValorDesconto != nil {
		return d.TotalProdutos - *d.ValorDesconto
	}
	return d.TotalProdutos
}

// GetMargemLucro calcula a margem de lucro do item
func (d *DocumentoVendaItem) GetMargemLucro() float64 {
	if d.CustoUnitario == 0 {
		return 0
	}
	return ((d.ValorUnitario - d.CustoUnitario) / d.ValorUnitario) * 100
}
