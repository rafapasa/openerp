package models

import (
	"time"
)

// ============================================================
// MODEL: OperacaoFiscal
// ============================================================

type OperacaoFiscal struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID               int        `gorm:"column:opf_id;primaryKey;autoIncrement" json:"id"`
	EmpresaFilialID  int        `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	CFOP             string     `gorm:"column:opf_cfop;type:varchar(8);not null" json:"cfop"`
	Descricao        string     `gorm:"column:opf_descricao;type:varchar(2000);not null" json:"descricao"`
	DataIni          time.Time  `gorm:"column:opf_dataini;type:date;not null" json:"data_ini"`
	DataFim          *time.Time `gorm:"column:opf_datafim;type:date" json:"data_fim,omitempty"`
	Carimbo          *string    `gorm:"column:opf_carimbo;type:text" json:"carimbo,omitempty"`
	ReducaoBCICMS_ST int        `gorm:"column:opf_reducao_bc_icms_st;not null;default:0" json:"reducao_bc_icms_st"`
	CSTIPIID         *int       `gorm:"column:cstipi_id" json:"cst_ipi_id,omitempty"`
	CSTPISCOFINSID   *int       `gorm:"column:cstpiscofins_id" json:"cst_pis_cofins_id,omitempty"`
	CSTICMSID        *int       `gorm:"column:csticms_id" json:"cst_icms_id,omitempty"`

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
	EmpresaFilial *EmpresaFilial `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
	CSTIPI        *CSTIPI        `gorm:"foreignKey:CSTIPIID;references:cstipi_id" json:"cst_ipi,omitempty"`
	CSTPISCOFINS  *CSTPISCOFINS  `gorm:"foreignKey:CSTPISCOFINSID;references:cstpiscofins_id" json:"cst_pis_cofins,omitempty"`
	CSTICMS       *CSTICMS       `gorm:"foreignKey:CSTICMSID;references:csticms_id" json:"cst_icms,omitempty"`
}

func (OperacaoFiscal) TableName() string {
	return "operacaofiscal"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (m *OperacaoFiscal) IsDeleted() bool {
	return m.DeletedAt != nil
}

func (m *OperacaoFiscal) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
}

func (m *OperacaoFiscal) IsActive() bool {
	now := time.Now()
	if m.DataFim == nil {
		return m.DataIni.Before(now) || m.DataIni.Equal(now)
	}
	return (m.DataIni.Before(now) || m.DataIni.Equal(now)) &&
		(m.DataFim.After(now) || m.DataFim.Equal(now))
}

func (m *OperacaoFiscal) HasReducaoBCICMS_ST() bool {
	return m.ReducaoBCICMS_ST == 1
}

func (m *OperacaoFiscal) GetTipoOperacao() string {
	// CFOP: 5.xxx = Entrada, 6.xxx = Saída
	if len(m.CFOP) > 0 {
		primeiroDigito := string(m.CFOP[0])
		switch primeiroDigito {
		case "5":
			return "Entrada"
		case "6":
			return "Saída"
		}
	}
	return "Desconhecido"
}
