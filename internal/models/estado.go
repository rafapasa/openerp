package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: Estado
// ============================================================

type Estado struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID                        int      `gorm:"column:est_id;primaryKey;autoIncrement" json:"id"`
	PaisID                    int      `gorm:"column:pai_id;not null" json:"pais_id"`
	UF                        string   `gorm:"column:est_uf;type:varchar(2);not null" json:"uf"`
	Nome                      string   `gorm:"column:est_nome;type:varchar(100);not null" json:"nome"`
	AliquotaICMS              float64  `gorm:"column:est_aliquotaicms;type:decimal(5,2);default:0.00" json:"aliquota_icms"`
	AliquotaICMS_ST           float64  `gorm:"column:est_aliquotaicms_st;type:decimal(5,2);default:0.00" json:"aliquota_icms_st"`
	AliquotaICMSInterestadual *float64 `gorm:"column:est_aliquotaicms_interest;type:decimal(5,2)" json:"aliquota_icms_interestadual,omitempty"`
	AliquotaFCP               *float64 `gorm:"column:est_aliquota_fcp;type:decimal(5,2)" json:"aliquota_fcp,omitempty"`

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
	Enderecos []EntidadeEndereco `gorm:"foreignKey:EstadoID;references:est_id" json:"enderecos,omitempty"`
	// DocumentoVenda []DocumentoVenda   `gorm:"foreignKey:EstadoID;references:est_id" json:"documentos,omitempty"`
	// NotasFiscais []NotaFiscal       `gorm:"foreignKey:EstadoID;references:est_id" json:"notas_fiscais,omitempty"`
	// NCMEstados   []NCMEstado        `gorm:"foreignKey:EstadoID;references:est_id" json:"ncm_estados,omitempty"`
}

func (Estado) TableName() string {
	return "estado"
}

func (m *Estado) BeforeCreate(tx *gorm.DB) error {
	if m.CreatedBy == nil {
		m.CreatedBy = new(int)
		*m.CreatedBy = 0
	}
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

func (m *Estado) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *Estado) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *Estado) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}

// IsActive verifica se o estado está ativo (não deletado)
func (m *Estado) IsActive() bool {
	return m.DeletedAt == nil
}

func (m *Estado) GetAliquotaICMS() float64 {
	return m.AliquotaICMS
}

func (m *Estado) GetAliquotaICMS_ST() float64 {
	return m.AliquotaICMS_ST
}

func (m *Estado) GetAliquotaICMSInterestadual() float64 {
	if m.AliquotaICMSInterestadual != nil {
		return *m.AliquotaICMSInterestadual
	}
	return 0
}

func (m *Estado) GetAliquotaFCP() float64 {
	if m.AliquotaFCP != nil {
		return *m.AliquotaFCP
	}
	return 0
}

func (m *Estado) GetNomeCompleto() string {
	return m.Nome + " (" + m.UF + ")"
}
