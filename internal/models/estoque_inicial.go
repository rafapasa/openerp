package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: EstoqueInicial
// Registra o estoque inicial de um produto em uma filial específica.
// ============================================================

type EstoqueInicial struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID              int       `gorm:"column:esi_id;primaryKey;autoIncrement" json:"id"`
	ProdutoID       int       `gorm:"column:pro_id;not null" json:"produto_id"`
	EmpresaFilialID int       `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	Quantidade      float64   `gorm:"column:esi_quantidade;type:decimal(15,4);not null" json:"quantidade"`
	DataLancamento  time.Time `gorm:"column:esi_datalancamento;type:date;not null" json:"data_lancamento"`

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
	Produto       *Produto       `gorm:"foreignKey:ProdutoID;references:ID" json:"produto,omitempty"`
	EmpresaFilial *EmpresaFilial `gorm:"foreignKey:EmpresaFilialID;references:ID" json:"empresa_filial,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (EstoqueInicial) TableName() string {
	return "estoque_inicial"
}

func (e *EstoqueInicial) BeforeCreate(tx *gorm.DB) error {
	if e.CreatedBy == nil {
		e.CreatedBy = new(int)
		*e.CreatedBy = 0
	}
	return nil
}