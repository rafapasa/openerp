package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: AssistenteProducao
// ============================================================

type AssistenteProducao struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID                int        `gorm:"column:asp_id;primaryKey;autoIncrement" json:"id"`
	EmpresaFilialID   int        `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	DataCriacao       time.Time  `gorm:"column:asp_datacriacao;type:date;not null" json:"data_criacao"`
	DataProcessamento *time.Time `gorm:"column:asp_dataprocessamento;type:date" json:"data_processamento,omitempty"`
	Observacao        *string    `gorm:"column:asp_observacao;type:text" json:"observacao,omitempty"`

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
	EmpresaFilial  *EmpresaFilial  `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
	OrdensProducao []OrdemProducao `gorm:"foreignKey:AssistenteProducaoID;references:asp_id" json:"ordens_producao,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (AssistenteProducao) TableName() string {
	return "assistente_producao"
}

func (a *AssistenteProducao) BeforeCreate(tx *gorm.DB) error {
	if a.CreatedBy == nil {
		a.CreatedBy = new(int)
		*a.CreatedBy = 0
	}
	if a.UpdatedBy == nil {
		a.UpdatedBy = new(int)
		*a.UpdatedBy = 0
	}
	return nil
}

func (a *AssistenteProducao) BeforeUpdate(tx *gorm.DB) error {
	if a.UpdatedBy == nil {
		a.UpdatedBy = new(int)
		*a.UpdatedBy = 0
	}
	return nil
}

// IsDeleted verifica se o assistente foi deletado logicamente
func (a *AssistenteProducao) IsDeleted() bool {
	return a.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (a *AssistenteProducao) SoftDelete() {
	now := time.Now()
	a.DeletedAt = &now
}

// IsProcessado verifica se o assistente foi processado
func (a *AssistenteProducao) IsProcessado() bool {
	return a.DataProcessamento != nil
}

// HasOrdensProducao verifica se possui ordens de produção
func (a *AssistenteProducao) HasOrdensProducao() bool {
	return len(a.OrdensProducao) > 0
}
