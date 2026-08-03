package models

import (
	"time"
)

// ============================================================
// MODEL: Municipio
// ============================================================

type Municipio struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID            int     `gorm:"column:mun_id;primaryKey;autoIncrement" json:"id"`
	EstadoID      int     `gorm:"column:est_id;not null" json:"estado_id"`
	Nome          string  `gorm:"column:mun_nome;type:varchar(100);not null" json:"nome"`
	CodigoFederal int     `gorm:"column:mun_codigofederal;not null" json:"codigo_federal"`
	AliquotaISS   float64 `gorm:"column:mun_aliquotaiss;type:decimal(5,2);not null" json:"aliquota_iss"`

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
	Enderecos []EntidadeEndereco `gorm:"foreignKey:MunicipioID;references:mun_id" json:"enderecos,omitempty"`
	// DocumentoVenda []DocumentoVenda   `gorm:"foreignKey:MunicipioID;references:mun_id" json:"documentos,omitempty"`
	// NotasFiscais []NotaFiscal       `gorm:"foreignKey:MunicipioID;references:ID" json:"notas_fiscais,omitempty"`
}

func (Municipio) TableName() string {
	return "municipio"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *Municipio) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *Municipio) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}

func (m *Municipio) GetAliquotaISS() float64 {
	return m.AliquotaISS
}
