package models

import (
	"time"
)

// ============================================================
// MODEL: RequisicaoItem
// ============================================================

type RequisicaoItem struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	RequisicaoID int     `gorm:"column:req_id;primaryKey" json:"requisicao_id"`
	Item         int     `gorm:"column:rei_item;primaryKey" json:"item"`
	ProdutoID    int     `gorm:"column:pro_id;not null" json:"produto_id"`
	Quantidade   float64 `gorm:"column:rei_quantidade;type:decimal(15,4);not null" json:"quantidade"`
	Custo        float64 `gorm:"column:rei_custo;type:decimal(15,2);not null" json:"custo"`

	GradeID   *int `gorm:"column:grade_id" json:"grade_id,omitempty"`
	GradeItem *int `gorm:"column:grai_item" json:"grade_item,omitempty"`

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
	Requisicao *Requisicao           `gorm:"foreignKey:RequisicaoID;references:req_id" json:"requisicao,omitempty"`
	Produto    *Produto              `gorm:"foreignKey:ProdutoID;references:pro_id" json:"produto,omitempty"`
	Grade      *GradeProduto         `gorm:"foreignKey:GradeID;references:grade_id" json:"grade,omitempty"`
	Grades     []RequisicaoItemGrade `gorm:"foreignKey:RequisicaoID,RequisicaoItem;references:req_id,rei_item" json:"grades,omitempty"`
	Lotes      []RequisicaoItemLote  `gorm:"foreignKey:RequisicaoID,RequisicaoItem;references:req_id,rei_item" json:"lotes,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (RequisicaoItem) TableName() string {
	return "requisicao_item"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o item foi deletado logicamente
func (r *RequisicaoItem) IsDeleted() bool {
	return r.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (r *RequisicaoItem) SoftDelete() {
	now := time.Now()
	r.DeletedAt = &now
}

// HasGrade verifica se o item possui grade
func (r *RequisicaoItem) HasGrade() bool {
	return r.GradeID != nil && r.GradeItem != nil
}

// HasLotes verifica se o item possui lotes
func (r *RequisicaoItem) HasLotes() bool {
	return len(r.Lotes) > 0
}

// HasGrades verifica se o item possui grades
func (r *RequisicaoItem) HasGrades() bool {
	return len(r.Grades) > 0
}

// GetValorTotal retorna o valor total do item (quantidade * custo)
func (r *RequisicaoItem) GetValorTotal() float64 {
	return r.Quantidade * r.Custo
}

// GetQuantidadeTotalLotes retorna a soma das quantidades dos lotes
func (r *RequisicaoItem) GetQuantidadeTotalLotes() float64 {
	total := 0.0
	for _, lote := range r.Lotes {
		total += lote.Quantidade
	}
	return total
}

// GetQuantidadeTotalGrades retorna a soma das quantidades das grades
func (r *RequisicaoItem) GetQuantidadeTotalGrades() float64 {
	total := 0.0
	for _, grade := range r.Grades {
		total += grade.Quantidade
	}
	return total
}
