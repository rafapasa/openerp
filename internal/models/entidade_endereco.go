package models

import (
	"time"

	"github.com/openerp/backend/internal/constants"
	"gorm.io/gorm"
)

// ============================================================
// MODEL: EntidadeEndereco
// ============================================================

type EntidadeEndereco struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	EntidadeID  int        `gorm:"column:ent_id;primaryKey" json:"entidade_id"`
	Item        int        `gorm:"column:ete_item;primaryKey" json:"id"`
	PaisID      int        `gorm:"column:pai_id;not null" json:"pais_id"`
	MunicipioID int        `gorm:"column:mun_id;not null" json:"municipio_id"`
	EstadoID    int        `gorm:"column:est_id;not null" json:"estado_id"`
	Tipo        int        `gorm:"column:ete_tipo;not null" json:"tipo"`
	CEP         int        `gorm:"column:ete_cep;not null" json:"cep"`
	Numero      string     `gorm:"column:ete_numero;type:varchar(15);not null" json:"numero"`
	Bairro      string     `gorm:"column:ete_bairro;type:varchar(100);not null" json:"bairro"`
	Complemento *string    `gorm:"column:ete_complemento;type:varchar(100)" json:"complemento,omitempty"`
	Situacao    int        `gorm:"column:ete_situacao;not null" json:"situacao"`
	DataIni     time.Time  `gorm:"column:ete_dataini;type:date;not null" json:"data_ini"`
	DataFim     *time.Time `gorm:"column:ete_datafim;type:date" json:"data_fim,omitempty"`
	Logradouro  *string    `gorm:"column:ete_logradouro;type:varchar(100)" json:"logradouro,omitempty"`
	Distancia   *float64   `gorm:"column:ete_distancia;type:decimal(15,4)" json:"distancia,omitempty"`
	Observacao  *string    `gorm:"column:ete_observacao;type:varchar(255)" json:"observacao,omitempty"`

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
	// references aponta para o nome da COLUNA na tabela referenciada
	Entidade  *Entidade  `gorm:"foreignKey:EntidadeID;references:ent_id" json:"entidade,omitempty"`
	Pais      *Pais      `gorm:"foreignKey:PaisID;references:pai_id" json:"pais,omitempty"`
	Municipio *Municipio `gorm:"foreignKey:MunicipioID;references:mun_id" json:"municipio,omitempty"`
	Estado    *Estado    `gorm:"foreignKey:EstadoID;references:est_id" json:"estado,omitempty"` // ============================================================

}

func (EntidadeEndereco) TableName() string {
	return "entidade_endereco"
}

func (m *EntidadeEndereco) BeforeCreate(tx *gorm.DB) error {
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

func (m *EntidadeEndereco) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *EntidadeEndereco) IsActive() bool {
	return m.Situacao == 1
}

func (m *EntidadeEndereco) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *EntidadeEndereco) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
	m.Situacao = int(constants.StatusInativo)
}

func (m *EntidadeEndereco) IsCobranca() bool {
	return m.Tipo == 1
}

func (m *EntidadeEndereco) IsEntrega() bool {
	return m.Tipo == 2
}

func (m *EntidadeEndereco) GetEnderecoCompleto() string {
	endereco := ""
	if m.Logradouro != nil {
		endereco += *m.Logradouro
	}
	endereco += ", " + m.Numero
	if m.Complemento != nil && *m.Complemento != "" {
		endereco += " - " + *m.Complemento
	}
	endereco += " - " + m.Bairro
	return endereco
}
