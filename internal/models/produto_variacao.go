package models

import (
	"time"
)

// ============================================================
// MODEL: ProdutoVariacao
// Representa uma variação específica de um produto, combinando cor e tamanho.
// ============================================================

type ProdutoVariacao struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID              int     `gorm:"column:provar_id;primaryKey;autoIncrement" json:"id"`
	ProdutoID       int     `gorm:"column:pro_id;not null" json:"produto_id"`
	EmpresaFilialID int     `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	CorID           *int    `gorm:"column:cor_id" json:"cor_id,omitempty"`                                                 // Opcional, se nem toda variação tiver cor
	TamanhoID       *int    `gorm:"column:ptam_id" json:"tamanho_id,omitempty"`                                            // Opcional, se nem toda variação tiver tamanho
	SKU             string  `gorm:"column:provar_sku;type:varchar(50);uniqueIndex:idx_provar_sku_emf;not null" json:"sku"` // SKU único por filial
	PrecoAdicional  float64 `gorm:"column:provar_preco_adicional;type:decimal(15,4);default:0.00" json:"preco_adicional"`
	EstoqueAtual    float64 `gorm:"column:provar_estoque_atual;type:decimal(15,4);default:0.00" json:"estoque_atual"` // Estoque desta variação

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
}

// TableName define o nome da tabela no banco de dados
func (ProdutoVariacao) TableName() string {
	return "produto_variacao"
}

// BeforeCreate hook do GORM para preencher CreatedBy e UpdatedBy

// BeforeUpdate hook do GORM para preencher UpdatedBy

// IsDeleted verifica se o registro foi deletado logicamente
func (m *ProdutoVariacao) IsDeleted() bool {
	return m.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (m *ProdutoVariacao) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}
