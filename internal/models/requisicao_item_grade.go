package models

import (
	"time"
)

// ============================================================
// MODEL: RequisicaoItemGrade
// ============================================================

type RequisicaoItemGrade struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	RequisicaoID   int     `gorm:"column:req_id;primaryKey" json:"requisicao_id"`
	RequisicaoItem int     `gorm:"column:rei_item;primaryKey" json:"requisicao_item"`
	Item           int     `gorm:"column:reig_item;primaryKey" json:"item"`
	CorID          int     `gorm:"column:cor_id;not null" json:"cor_id"`
	TamanhoID      int     `gorm:"column:ptam_id;not null" json:"tamanho_id"`
	Quantidade     float64 `gorm:"column:reig_quantidade;type:decimal(15,4);not null;default:0" json:"quantidade"`

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
	Cor     *ProdutoCor     `gorm:"foreignKey:CorID;references:cor_id" json:"cor,omitempty"`
	Tamanho *ProdutoTamanho `gorm:"foreignKey:TamanhoID;references:ptam_id" json:"tamanho,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (RequisicaoItemGrade) TableName() string {
	return "requisicao_item_grade"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o registro foi deletado logicamente
func (r *RequisicaoItemGrade) IsDeleted() bool {
	return r.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (r *RequisicaoItemGrade) SoftDelete() {
	now := time.Now()
	r.DeletedAt = &now
}
