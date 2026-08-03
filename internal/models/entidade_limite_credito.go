package models

import (
	"time"
)

// ============================================================
// MODEL: EntidadeLimiteCredito
// ============================================================

type EntidadeLimiteCredito struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	EntidadeID   int       `gorm:"column:ent_id;primaryKey" json:"entidade_id"`
	Item         int       `gorm:"column:elc_item;primaryKey" json:"item"`
	Data         time.Time `gorm:"column:elc_data;type:date;not null" json:"data"`
	Valor        float64   `gorm:"column:elc_valor;type:decimal(15,4);not null" json:"valor"`
	Descricao    *string   `gorm:"column:elc_descricao;type:text" json:"descricao,omitempty"`
	DiasBloqueio *int      `gorm:"column:elc_diasbloqueio" json:"dias_bloqueio,omitempty"`

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

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (EntidadeLimiteCredito) TableName() string {
	return "entidade_limitecredito"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o limite de crédito foi deletado logicamente
func (e *EntidadeLimiteCredito) IsDeleted() bool {
	return e.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (e *EntidadeLimiteCredito) SoftDelete() {
	now := time.Now()
	e.DeletedAt = &now
}

// HasDescricao verifica se possui descrição
func (e *EntidadeLimiteCredito) HasDescricao() bool {
	return e.Descricao != nil && *e.Descricao != ""
}

// HasDiasBloqueio verifica se possui dias de bloqueio
func (e *EntidadeLimiteCredito) HasDiasBloqueio() bool {
	return e.DiasBloqueio != nil && *e.DiasBloqueio > 0
}

// GetDescricaoOrDefault retorna a descrição ou string vazia
func (e *EntidadeLimiteCredito) GetDescricaoOrDefault() string {
	if e.HasDescricao() {
		return *e.Descricao
	}
	return ""
}

// GetDiasBloqueioOrDefault retorna os dias de bloqueio ou 0
func (e *EntidadeLimiteCredito) GetDiasBloqueioOrDefault() int {
	if e.HasDiasBloqueio() {
		return *e.DiasBloqueio
	}
	return 0
}

// IsValido verifica se o limite de crédito está válido
// Considera válido se a data de vigência não passou
func (e *EntidadeLimiteCredito) IsValido() bool {
	return time.Now().Before(e.Data) || time.Now().Equal(e.Data)
}

// IsVencido verifica se o limite de crédito está vencido
func (e *EntidadeLimiteCredito) IsVencido() bool {
	return time.Now().After(e.Data)
}

// GetDataVigencia retorna a data de vigência formatada
func (e *EntidadeLimiteCredito) GetDataVigencia() string {
	return e.Data.Format("2006-01-02")
}

// GetDiasRestantes retorna os dias restantes até o vencimento
func (e *EntidadeLimiteCredito) GetDiasRestantes() int {
	if e.IsVencido() {
		return 0
	}
	dias := int(time.Until(e.Data).Hours() / 24)
	if dias < 0 {
		return 0
	}
	return dias
}
