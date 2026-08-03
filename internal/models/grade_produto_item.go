package models

import (
	"time"
)

// ============================================================
// MODEL: GradeProdutoItem
// ============================================================

type GradeProdutoItem struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	GradeID   int `gorm:"column:grade_id;primaryKey" json:"grade_id"`
	Item      int `gorm:"column:grai_item;primaryKey;autoIncrement" json:"item"`
	CorID     int `gorm:"column:cor_id;not null" json:"cor_id"`
	TamanhoID int `gorm:"column:ptam_id;not null" json:"tamanho_id"`

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
	Grade   *GradeProduto   `gorm:"foreignKey:GradeID;references:grade_id" json:"grade,omitempty"`
	Cor     *ProdutoCor     `gorm:"foreignKey:CorID;references:cor_id" json:"cor,omitempty"`
	Tamanho *ProdutoTamanho `gorm:"foreignKey:TamanhoID;references:ptam_id" json:"tamanho,omitempty"`
}

func (GradeProdutoItem) TableName() string {
	return "grade_produto_item"
}

// Buscar próximo item para esta grade

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (g *GradeProdutoItem) IsDeleted() bool {
	return g.DeletedAt != nil
}

func (g *GradeProdutoItem) SoftDelete() {
	now := time.Now()
	g.DeletedAt = &now
}

// IsActive verifica se o item da grade está ativo (não deletado)
func (g *GradeProdutoItem) IsActive() bool {
	return g.DeletedAt == nil
}

// GetNomeCompleto retorna o nome completo da combinação (Cor + Tamanho)
func (g *GradeProdutoItem) GetNomeCompleto() string {
	corNome := ""
	tamanhoNome := ""

	if g.Cor != nil {
		corNome = g.Cor.Nome
	}
	if g.Tamanho != nil {
		tamanhoNome = g.Tamanho.Nome
	}

	if corNome != "" && tamanhoNome != "" {
		return corNome + " - " + tamanhoNome
	}
	if corNome != "" {
		return corNome
	}
	if tamanhoNome != "" {
		return tamanhoNome
	}
	return "Item " + string(rune(g.Item))
}
