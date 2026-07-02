package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: FormatoProduto
// ============================================================

type FormatoProduto struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID        int    `gorm:"column:fpro_id;primaryKey;autoIncrement" json:"id"`
	Descricao string `gorm:"column:fpro_descricao;type:varchar(255);not null" json:"descricao"`
	Formula   string `gorm:"column:fpro_formula;type:text;not null" json:"formula"`
	Imagem    []byte `gorm:"column:fpro_imagem;type:longblob" json:"imagem,omitempty"`

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
	Medidas  []FormatoProdutoMedida  `gorm:"foreignKey:FormatoProdutoID;references:fpro_id" json:"medidas,omitempty"`
	Produtos []ProdutoFormatoProduto `gorm:"foreignKey:FormatoProdutoID;references:fpro_id" json:"produtos,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (FormatoProduto) TableName() string {
	return "formato_produto"
}

func (f *FormatoProduto) BeforeCreate(tx *gorm.DB) error {
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

func (f *FormatoProduto) BeforeUpdate(tx *gorm.DB) error {
	if f.UpdatedBy == nil {
		f.UpdatedBy = new(int)
		*f.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o formato foi deletado logicamente
func (f *FormatoProduto) IsDeleted() bool {
	return f.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (f *FormatoProduto) SoftDelete() {
	now := time.Now()
	f.DeletedAt = &now
}

// HasImage verifica se o formato possui imagem
func (f *FormatoProduto) HasImage() bool {
	return len(f.Imagem) > 0
}

// GetMedidasCount retorna a quantidade de medidas associadas
func (f *FormatoProduto) GetMedidasCount() int {
	return len(f.Medidas)
}

// GetProdutosCount retorna a quantidade de produtos associados
func (f *FormatoProduto) GetProdutosCount() int {
	return len(f.Produtos)
}
