package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: TituloOcorrenciaTipo
// ============================================================

type TituloOcorrenciaTipo struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID        int    `gorm:"column:tot_id;primaryKey;autoIncrement" json:"id"`
	Descricao string `gorm:"column:tot_descricao;type:varchar(255);not null" json:"descricao"`

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
	Ocorrencias []TituloOcorrencia `gorm:"foreignKey:OcorrenciaTipoID;references:ID" json:"ocorrencias,omitempty"`
}

func (TituloOcorrenciaTipo) TableName() string {
	return "titulo_ocorrencia_tipo"
}

func (t *TituloOcorrenciaTipo) BeforeCreate(tx *gorm.DB) error {
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

func (t *TituloOcorrenciaTipo) BeforeUpdate(tx *gorm.DB) error {
	if t.UpdatedBy == nil {
		t.UpdatedBy = new(int)
		*t.UpdatedBy = 0
	}
	return nil
}

func (t *TituloOcorrenciaTipo) IsDeleted() bool {
	return t.DeletedAt != nil
}

func (t *TituloOcorrenciaTipo) SoftDelete() {
	now := time.Now()
	t.DeletedAt = &now
}
