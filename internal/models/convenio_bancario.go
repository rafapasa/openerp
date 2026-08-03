package models

import (
	"time"
)

// ============================================================
// MODEL: ConvenioBancario
// ============================================================

type ConvenioBancario struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID                 int     `gorm:"column:cvb_id;primaryKey;autoIncrement" json:"id"`
	ContaBancariaID    int     `gorm:"column:ctb_id;not null" json:"conta_bancaria_id"`
	Descricao          string  `gorm:"column:cvb_descricao;type:varchar(255);not null" json:"descricao"`
	NumeroConvenio     string  `gorm:"column:cvb_numeroconvenio;type:varchar(20);not null" json:"numero_convenio"`
	Carteira           string  `gorm:"column:cvb_carteira;type:varchar(10);not null" json:"carteira"`
	BoletoInicial      *int64  `gorm:"column:cvb_boletoinicial;type:bigint" json:"boleto_inicial,omitempty"`
	BoletoFinal        *int64  `gorm:"column:cvb_boletofinal;type:bigint" json:"boleto_final,omitempty"`
	UltimoBoleto       int64   `gorm:"column:cvb_ultimoboleto;type:bigint;not null" json:"ultimo_boleto"`
	DiasProtesto       int     `gorm:"column:cvb_diasprotesto;not null" json:"dias_protesto"`
	TipoCobranca       *string `gorm:"column:cvb_tipocobranca;type:varchar(10)" json:"tipo_cobranca,omitempty"`
	TipoCarteira       *int    `gorm:"column:cvb_tipocarteira" json:"tipo_carteira,omitempty"`
	VariacaoCarteira   *string `gorm:"column:cvb_variacaocarteira;type:varchar(3)" json:"variacao_carteira,omitempty"`
	Layout             *int    `gorm:"column:cvb_layout" json:"layout,omitempty"`
	ResponsavelEmissao *int    `gorm:"column:cvb_responsavelemissao" json:"responsavel_emissao,omitempty"`

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
	ContaBancaria *ContaBancaria `gorm:"foreignKey:ContaBancariaID;references:ctb_id" json:"conta_bancaria,omitempty"`
	Portadores    []Portador     `gorm:"foreignKey:ConvenioBancarioID;references:ID" json:"portadores,omitempty"`
}

func (ConvenioBancario) TableName() string {
	return "convenio_bancario"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *ConvenioBancario) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *ConvenioBancario) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}

func (m *ConvenioBancario) GetProximoBoleto() int64 {
	m.UltimoBoleto++
	return m.UltimoBoleto
}

func (m *ConvenioBancario) HasLimiteBoleto() bool {
	return m.BoletoInicial != nil && m.BoletoFinal != nil
}

func (m *ConvenioBancario) GetProximoNumeroBoleto() (int64, error) {
	// TODO: Implementar validação de limite
	// if m.HasLimiteBoleto() && m.UltimoBoleto >= *m.BoletoFinal {
	//     return 0, errors.New("limite de boletos atingido")
	// }
	return m.GetProximoBoleto(), nil
}
