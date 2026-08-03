package models

import (
	"time"
)

// ============================================================
// MODEL: TituloOcorrencia
// ============================================================

type TituloOcorrencia struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	TituloID         int       `gorm:"column:tit_id;primaryKey" json:"titulo_id"`
	Item             int       `gorm:"column:tito_item;primaryKey;autoIncrement" json:"item"`
	OcorrenciaTipoID int       `gorm:"column:tot_id;not null" json:"ocorrencia_tipo_id"`
	Data             time.Time `gorm:"column:tito_data;type:date;not null" json:"data"`
	Observacao       *string   `gorm:"column:tito_observacao;type:varchar(255)" json:"observacao,omitempty"`

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
	Titulo         *Titulo               `gorm:"foreignKey:TituloID;references:tit_id" json:"titulo,omitempty"`
	OcorrenciaTipo *TituloOcorrenciaTipo `gorm:"foreignKey:OcorrenciaTipoID;references:tot_id" json:"ocorrencia_tipo,omitempty"`
}

func (TituloOcorrencia) TableName() string {
	return "titulo_ocorrencia"
}

func (t *TituloOcorrencia) IsDeleted() bool {
	return t.DeletedAt != nil
}

func (t *TituloOcorrencia) SoftDelete() {
	now := time.Now()
	t.DeletedAt = &now
}
