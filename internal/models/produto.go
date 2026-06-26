package models

import (
	"time"
)

type Produto struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID                int      `gorm:"column:pro_id;primaryKey;autoIncrement" json:"id"`
	EmpresaFilialID   int      `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	TipoProdutoID     int      `gorm:"column:ptp_id;not null" json:"tipo_produto_id"`
	GrupoProdutoID    *int     `gorm:"column:prog_id" json:"grupo_produto_id,omitempty"`
	SubGrupoProdutoID *int     `gorm:"column:prosg_id" json:"sub_grupo_produto_id,omitempty"`
	MarcaID           *int     `gorm:"column:promar_id" json:"marca_id,omitempty"`
	NcmNumero         int      `gorm:"column:ncm_numero;not null" json:"ncm_numero"` // CORRIGIDO: int, não string
	Codigo            int      `gorm:"column:pro_codigo;not null;unique" json:"codigo"`
	Nome              string   `gorm:"column:pro_nome;type:varchar(80);not null" json:"nome"`     // CORRIGIDO: varchar(80)
	Resumo            string   `gorm:"column:pro_resumo;type:varchar(80);not null" json:"resumo"` // CORRIGIDO: varchar(80) e not null
	Situacao          int      `gorm:"column:pro_situacao;not null" json:"situacao"`
	CustoCompra       *float64 `gorm:"column:pro_custocompra;type:decimal(15,4)" json:"custo_compra,omitempty"`
	CodigoBarras      *string  `gorm:"column:pro_codigobarra;type:varchar(255)" json:"codigo_barras,omitempty"` // CORRIGIDO: nome e tamanho

	// ============================================================
	// CAMPOS DE AUDITORIA (VISÍVEIS NO JSON)
	// ============================================================
	CreatedAt time.Time  `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime;index" json:"deleted_at,omitempty"`
	CreatedBy *int       `gorm:"column:created_by" json:"created_by,omitempty"`
	UpdatedBy *int       `gorm:"column:updated_by" json:"updated_by,omitempty"`

	// ============================================================
	// RELACIONAMENTOS
	// ============================================================
	TipoProduto     *TipoProduto     `gorm:"foreignKey:TipoProdutoID;references:ptp_id" json:"tipo_produto,omitempty"`
	GrupoProduto    *ProdutoGrupo    `gorm:"foreignKey:GrupoProdutoID;references:prog_id" json:"grupo_produto,omitempty"`
	SubGrupoProduto *ProdutoSubgrupo `gorm:"foreignKey:SubGrupoProdutoID;references:prosg_id" json:"sub_grupo_produto,omitempty"`
	Marca           *ProdutoMarca    `gorm:"foreignKey:MarcaID;references:promar_id" json:"marca,omitempty"`
}

func (Produto) TableName() string {
	return "produto"
}

func (p *Produto) BeforeCreate() error {
	if p.CreatedBy == nil {
		p.CreatedBy = new(int)
		*p.CreatedBy = 0
	}
	if p.UpdatedBy == nil {
		p.UpdatedBy = new(int)
		*p.UpdatedBy = 0
	}
	return nil
}

func (p *Produto) BeforeUpdate() error {
	if p.UpdatedBy == nil {
		p.UpdatedBy = new(int)
		*p.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (p *Produto) IsActive() bool {
	return p.Situacao == 1
}

func (p *Produto) IsDeleted() bool {
	return p.DeletedAt != nil
}

func (p *Produto) SoftDelete() {
	now := time.Now()
	p.DeletedAt = &now
	p.Situacao = 0
}

// SafeToDelete verifica se o produto pode ser excluído
// Verifica se há registros relacionados (pedidos, estoque, etc.)
func (p *Produto) SafeToDelete() bool {
	// Aqui você pode adicionar lógica para verificar se o produto pode ser excluído
	// Por exemplo, verificar se há registros relacionados em outras tabelas
	// Por enquanto retorna true, mas em produção você deve verificar dependências
	return true
}
