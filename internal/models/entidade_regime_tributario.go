package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// CONSTANTES
// ============================================================
// (Adicionar no final do constants.go)

/*
const (
    RegimeTributarioSimplesNacional  = 1
    RegimeTributarioLucroPresumido   = 2
    RegimeTributarioLucroReal        = 3
    RegimeTributarioMEI              = 4
    RegimeTributarioIsento           = 5
)
*/

// ============================================================
// MODEL: EntidadeRegimeTributario
// ============================================================

type EntidadeRegimeTributario struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID              int       `gorm:"column:etr_item;primaryKey;autoIncrement" json:"id"`
	EntidadeID      int       `gorm:"column:ent_id;not null" json:"entidade_id"`
	Regime          int       `gorm:"column:etr_regime;not null" json:"regime"`
	Apuracao        *int      `gorm:"column:etr_apuracao" json:"apuracao,omitempty"`
	Data            time.Time `gorm:"column:etr_data;type:date;not null" json:"data"`
	RegimeEspecial  int       `gorm:"column:etr_regime_especial;not null;default:0" json:"regime_especial"`
	SituacaoTribISS *int      `gorm:"column:etr_situacao_trib_iss" json:"situacao_trib_iss,omitempty"`
	RegimeMunicipal *int      `gorm:"column:etr_regime_municipal" json:"regime_municipal,omitempty"`

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
	Entidade *Entidade `gorm:"foreignKey:EntidadeID;references:ent_id" json:"entidade,omitempty"`
}

func (EntidadeRegimeTributario) TableName() string {
	return "entidade_regime_tributario"
}

func (m *EntidadeRegimeTributario) BeforeCreate(tx *gorm.DB) error {
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

func (m *EntidadeRegimeTributario) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *EntidadeRegimeTributario) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *EntidadeRegimeTributario) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}

func (m *EntidadeRegimeTributario) IsSimplesNacional() bool {
	return m.Regime == 1
}

func (m *EntidadeRegimeTributario) IsLucroPresumido() bool {
	return m.Regime == 2
}

func (m *EntidadeRegimeTributario) IsLucroReal() bool {
	return m.Regime == 3
}

func (m *EntidadeRegimeTributario) IsMEI() bool {
	return m.Regime == 4
}

func (m *EntidadeRegimeTributario) IsIsento() bool {
	return m.Regime == 5
}
