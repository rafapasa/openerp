package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: TipoRequisicao
// ============================================================

type TipoRequisicao struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID        int    `gorm:"column:treq_id;primaryKey;autoIncrement" json:"id"`
	Descricao string `gorm:"column:treq_descricao;type:varchar(255);not null" json:"descricao"`

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
	Requisicoes []Requisicao `gorm:"foreignKey:TipoRequisicaoID;references:treq_id" json:"requisicoes,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (TipoRequisicao) TableName() string {
	return "tipo_requisicao"
}

func (t *TipoRequisicao) BeforeCreate(tx *gorm.DB) error {
	if t.CreatedBy == nil {
		t.CreatedBy = new(int)
		*t.CreatedBy = 0
	}
	if t.UpdatedBy == nil {
		t.UpdatedBy = new(int)
		*t.UpdatedBy = 0
	}
	return nil
}

func (t *TipoRequisicao) BeforeUpdate(tx *gorm.DB) error {
	if t.UpdatedBy == nil {
		t.UpdatedBy = new(int)
		*t.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o tipo de requisição foi deletado logicamente
func (t *TipoRequisicao) IsDeleted() bool {
	return t.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (t *TipoRequisicao) SoftDelete() {
	now := time.Now()
	t.DeletedAt = &now
}

// IsActive verifica se o tipo de requisição está ativo
func (t *TipoRequisicao) IsActive() bool {
	return t.DeletedAt == nil
}

// HasRequisicoes verifica se o tipo possui requisições associadas
func (t *TipoRequisicao) HasRequisicoes() bool {
	return len(t.Requisicoes) > 0
}
