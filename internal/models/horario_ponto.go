package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: HorarioPonto
// ============================================================

type HorarioPonto struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID          int       `gorm:"column:horp_item;primaryKey;autoIncrement" json:"id"`
	HorarioID   int       `gorm:"column:hor_id;not null" json:"horario_id"`
	TurnoID     int       `gorm:"column:hort_id;not null" json:"turno_id"`
	DiaSemana   int       `gorm:"column:horp_diasemana;not null" json:"dia_semana"`
	HoraEntrada time.Time `gorm:"column:horp_horaentrada;type:time;not null" json:"hora_entrada"`
	HoraSaida   time.Time `gorm:"column:horp_horasaia;type:time;not null" json:"hora_saida"`
	DSR         *int      `gorm:"column:horp_dsr" json:"dsr,omitempty"`

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
	Horario *Horario      `gorm:"foreignKey:HorarioID;references:hor_id" json:"horario,omitempty"`
	Turno   *HorarioTurno `gorm:"foreignKey:TurnoID;references:hort_id" json:"turno,omitempty"`
}

func (HorarioPonto) TableName() string {
	return "horario_ponto"
}

func (m *HorarioPonto) BeforeCreate(tx *gorm.DB) error {
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

func (m *HorarioPonto) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *HorarioPonto) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *HorarioPonto) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}

func (m *HorarioPonto) IsDSR() bool {
	return m.DSR != nil && *m.DSR == 1
}

func (m *HorarioPonto) GetDiaSemana() string {
	dias := map[int]string{
		1: "Domingo",
		2: "Segunda-feira",
		3: "Terça-feira",
		4: "Quarta-feira",
		5: "Quinta-feira",
		6: "Sexta-feira",
		7: "Sábado",
	}
	if nome, ok := dias[m.DiaSemana]; ok {
		return nome
	}
	return "Desconhecido"
}
