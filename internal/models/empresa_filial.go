package models

import (
	"time"
)

// ============================================================
// CONSTANTES
// ============================================================

const (
	TipoFilialMatriz = 1
	TipoFilialFilial = 2
)

// ============================================================
// MODEL: EmpresaFilial
// ============================================================

type EmpresaFilial struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID         int    `gorm:"column:emf_id;primaryKey;autoIncrement" json:"id"`
	EmpresaID  int    `gorm:"column:emp_id;not null" json:"empresa_id"` // CORRIGIDO: NOT NULL
	EntidadeID *int   `gorm:"column:ent_id" json:"entidade_id,omitempty"`
	Numero     int    `gorm:"column:emf_numero;not null" json:"numero"`
	Nome       string `gorm:"column:emf_nome;type:varchar(100);not null" json:"nome"`

	// ============================================================
	// CAMPOS DE ALÍQUOTAS
	// ============================================================
	AliquotaPIS      float64  `gorm:"column:emf_aliquotapis;type:decimal(5,2);not null" json:"aliquota_pis"`
	AliquotaCOFINS   float64  `gorm:"column:emf_aliquotacofins;type:decimal(5,2);not null" json:"aliquota_cofins"`
	AliquotaISS      *float64 `gorm:"column:emf_aliquotaiss;type:decimal(5,2)" json:"aliquota_iss,omitempty"`
	AliquotaFunrural *float64 `gorm:"column:emf_aliquotafunrural;type:decimal(5,2)" json:"aliquota_funrural,omitempty"`

	// ============================================================
	// CAMPOS DE INFORMAÇÕES
	// ============================================================
	CodigoCNAE *string `gorm:"column:emf_codigocnae;type:varchar(20)" json:"codigo_cnae,omitempty"`
	Mei        int8    `gorm:"column:emf_mei;type:smallint" json:"mei"`

	// ============================================================
	// CAMPOS DE LOGOMARCA E WEB
	// ============================================================
	Logomarca    []byte  `gorm:"column:emf_logomarca;type:longblob" json:"-"`
	LogomarcaWeb *string `gorm:"column:emf_logomarcaweb;type:varchar(1000)" json:"logomarca_web,omitempty"`
	EnderecoWeb  *string `gorm:"column:emf_enderecoweb;type:varchar(1000)" json:"endereco_web,omitempty"`

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
	Empresa *Empresa `gorm:"foreignKey:EmpresaID;references:emp_id" json:"empresa,omitempty"`
	// Entidade *Entidade `gorm:"foreignKey:EntidadeID;references:ent_id" json:"entidade,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (EmpresaFilial) TableName() string {
	return "empresa_filial"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (e *EmpresaFilial) IsMEI() bool {
	return e.Mei == 1
}

func (e *EmpresaFilial) IsDeleted() bool {
	return e.DeletedAt != nil
}

func (e *EmpresaFilial) IsMatriz() bool {
	return e.Numero == TipoFilialMatriz
}

func (e *EmpresaFilial) IsFilial() bool {
	return e.Numero != TipoFilialMatriz
}

func (e *EmpresaFilial) SoftDelete() {
	now := time.Now()
	e.DeletedAt = &now
}
