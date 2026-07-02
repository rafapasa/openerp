package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: FormatoProdutoMedida
// ============================================================

type FormatoProdutoMedida struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	FormatoProdutoID int     `gorm:"column:fpro_id;primaryKey" json:"formato_produto_id"`
	Item             int     `gorm:"column:fprom_item;primaryKey" json:"item"`
	Descricao        string  `gorm:"column:fprom_descricao;type:varchar(255);not null" json:"descricao"`
	Codigo           string  `gorm:"column:fprom_codigo;type:varchar(20);not null" json:"codigo"`
	ValorMinimo      float64 `gorm:"column:fprom_valorminimo;type:decimal(15,4);not null" json:"valor_minimo"`
	ValorMaximo      float64 `gorm:"column:fprom_valormaximo;type:decimal(15,4);not null" json:"valor_maximo"`

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
	FormatoProduto *FormatoProduto `gorm:"foreignKey:FormatoProdutoID;references:fpro_id" json:"formato_produto,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (FormatoProdutoMedida) TableName() string {
	return "formato_produto_medida"
}

func (f *FormatoProdutoMedida) BeforeCreate(tx *gorm.DB) error {
	if f.CreatedBy == nil {
		f.CreatedBy = new(int)
		*f.CreatedBy = 0
	}
	if f.UpdatedBy == nil {
		f.UpdatedBy = new(int)
		*f.UpdatedBy = 0
	}
	return nil
}

func (f *FormatoProdutoMedida) BeforeUpdate(tx *gorm.DB) error {
	if f.UpdatedBy == nil {
		f.UpdatedBy = new(int)
		*f.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se a medida foi deletada logicamente
func (f *FormatoProdutoMedida) IsDeleted() bool {
	return f.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (f *FormatoProdutoMedida) SoftDelete() {
	now := time.Now()
	f.DeletedAt = &now
}

// IsValidValue verifica se o valor está dentro do intervalo permitido
func (f *FormatoProdutoMedida) IsValidValue(valor float64) bool {
	return valor >= f.ValorMinimo && valor <= f.ValorMaximo
}
