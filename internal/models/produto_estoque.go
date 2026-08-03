package models

import (
	"time"
)

// ============================================================
// MODEL: ProdutoEstoque
// ============================================================

type ProdutoEstoque struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID               int       `gorm:"column:proest_id;primaryKey;autoIncrement" json:"id"`
	EmpresaFilialID  int       `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	ProdutoID        int       `gorm:"column:pro_id;not null" json:"produto_id"`
	IDDocGerador     int       `gorm:"column:proest_iddocgerador;not null" json:"id_doc_gerador"`
	Origem           int       `gorm:"column:proest_origem;not null" json:"origem"`
	DataMovimentacao time.Time `gorm:"column:proest_datamovimentacao;type:date;not null" json:"data_movimentacao"`
	CustoFiscal      *float64  `gorm:"column:proest_custofiscal;type:decimal(15,3)" json:"custo_fiscal,omitempty"`
	CustoMedio       *float64  `gorm:"column:proest_customedio;type:decimal(15,2)" json:"custo_medio,omitempty"`
	Quantidade       float64   `gorm:"column:proest_quantidade;type:decimal(15,4);not null" json:"quantidade"`
	TipoMovimentacao int       `gorm:"column:proest_tipomovimentacao;not null" json:"tipo_movimentacao"` // 1-entrada, -1-saida
	CorID            *int      `gorm:"column:cor_id" json:"cor_id,omitempty"`
	TamanhoID        *int      `gorm:"column:ptam_id" json:"tamanho_id,omitempty"`
	TalhaoID         *int      `gorm:"column:tat_id" json:"talhao_id,omitempty"`
	Condicional      *int      `gorm:"column:proest_condicional" json:"condicional,omitempty"`

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
	Produto       *Produto        `gorm:"foreignKey:ProdutoID;references:pro_id" json:"produto,omitempty"`
	EmpresaFilial *EmpresaFilial  `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
	Cor           *ProdutoCor     `gorm:"foreignKey:CorID;references:cor_id" json:"cor,omitempty"`
	Tamanho       *ProdutoTamanho `gorm:"foreignKey:TamanhoID;references:ptam_id" json:"tamanho,omitempty"`
	Talhao        *TalhaoTerra    `gorm:"foreignKey:TalhaoID;references:tat_id" json:"talhao,omitempty"`
}

func (ProdutoEstoque) TableName() string {
	return "produto_estoque"
}

func (m *ProdutoEstoque) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *ProdutoEstoque) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}

func (m *ProdutoEstoque) IsEntrada() bool {
	return m.TipoMovimentacao == 1
}

func (m *ProdutoEstoque) IsSaida() bool {
	return m.TipoMovimentacao == -1
}
