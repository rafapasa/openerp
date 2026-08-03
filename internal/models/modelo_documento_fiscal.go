package models

import (
	"time"

	"github.com/openerp/backend/internal/constants"
)

// ============================================================
// MODEL: ModeloDocumentoFiscal
// ============================================================

type ModeloDocumentoFiscal struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	EmpresaFilialID int              `gorm:"column:emf_id;primaryKey" json:"empresa_filial_id"`
	Codigo          string           `gorm:"column:mdf_codigo;type:varchar(3);primaryKey" json:"codigo"`
	Descricao       string           `gorm:"column:mdf_descricao;type:varchar(255);not null" json:"descricao"`
	Situacao        constants.Status `gorm:"column:mdf_situacao;not null;default:1" json:"situacao"`

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
	EmpresaFilial *EmpresaFilial               `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
	Series        []ModeloDocumentoFiscalSerie `gorm:"foreignKey:EmpresaFilialID,ModeloDocumentoFiscalCodigo;references:emf_id,mdf_codigo" json:"series,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (ModeloDocumentoFiscal) TableName() string {
	return "modelo_documento_fiscal"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsActive verifica se o modelo está ativo
func (m *ModeloDocumentoFiscal) IsActive() bool {
	return m.Situacao == constants.StatusAtivo
}

// IsDeleted verifica se o modelo foi deletado logicamente
func (m *ModeloDocumentoFiscal) IsDeleted() bool {
	return m.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (m *ModeloDocumentoFiscal) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
	m.Situacao = constants.StatusInativo
}

// HasSeries verifica se o modelo possui séries associadas
func (m *ModeloDocumentoFiscal) HasSeries() bool {
	return len(m.Series) > 0
}

// GetSeriesCount retorna a quantidade de séries
func (m *ModeloDocumentoFiscal) GetSeriesCount() int {
	return len(m.Series)
}

// GetSerieByNumber retorna uma série pelo número
func (m *ModeloDocumentoFiscal) GetSerieByNumber(serie int) *ModeloDocumentoFiscalSerie {
	for i := range m.Series {
		if m.Series[i].Serie == serie {
			return &m.Series[i]
		}
	}
	return nil
}

// GetActiveSeries retorna apenas as séries ativas
func (m *ModeloDocumentoFiscal) GetActiveSeries() []ModeloDocumentoFiscalSerie {
	var active []ModeloDocumentoFiscalSerie
	for _, s := range m.Series {
		if s.IsActive() {
			active = append(active, s)
		}
	}
	return active
}

// IsEletronico verifica se o modelo é eletrônico
func (m *ModeloDocumentoFiscal) IsEletronico() bool {
	eletronicos := map[string]bool{
		"55": true,
		"57": true,
		"58": true,
		"59": true,
		"63": true,
		"65": true,
		"67": true,
	}
	return eletronicos[m.Codigo]
}

// IsPapel verifica se o modelo é em papel
func (m *ModeloDocumentoFiscal) IsPapel() bool {
	return !m.IsEletronico()
}

// IsNFe verifica se é NFe (Modelo 55)
func (m *ModeloDocumentoFiscal) IsNFe() bool {
	return m.Codigo == "55"
}

// IsNFCe verifica se é NFCe (Modelo 65)
func (m *ModeloDocumentoFiscal) IsNFCe() bool {
	return m.Codigo == "65"
}

// IsCTe verifica se é CTe (Modelo 57)
func (m *ModeloDocumentoFiscal) IsCTe() bool {
	return m.Codigo == "57"
}

// IsCTeOS verifica se é CTe OS (Modelo 67)
func (m *ModeloDocumentoFiscal) IsCTeOS() bool {
	return m.Codigo == "67"
}

// IsMDFe verifica se é MDFe (Modelo 58)
func (m *ModeloDocumentoFiscal) IsMDFe() bool {
	return m.Codigo == "58"
}

// IsCFe verifica se é CFe (Modelo 59)
func (m *ModeloDocumentoFiscal) IsCFe() bool {
	return m.Codigo == "59"
}

// IsBPe verifica se é BPe (Modelo 63)
func (m *ModeloDocumentoFiscal) IsBPe() bool {
	return m.Codigo == "63"
}

// GetTipoDocumento retorna o tipo de documento por extenso
func (m *ModeloDocumentoFiscal) GetTipoDocumento() string {
	switch m.Codigo {
	case "01":
		return "Nota Fiscal (Modelo 01)"
	case "04":
		return "Nota Fiscal de Produtor (Modelo 04)"
	case "06":
		return "Nota Fiscal de Energia Elétrica (Modelo 06)"
	case "07":
		return "Nota Fiscal de Serviço de Transporte (Modelo 07)"
	case "08":
		return "Conhecimento de Transporte Rodoviário de Cargas (Modelo 08)"
	case "55":
		return "NF-e (Nota Fiscal Eletrônica)"
	case "57":
		return "CT-e (Conhecimento de Transporte Eletrônico)"
	case "58":
		return "MDF-e (Manifesto Eletrônico de Documentos Fiscais)"
	case "59":
		return "CF-e SAT (Cupom Fiscal Eletrônico)"
	case "63":
		return "BP-e (Bilhete de Passagem Eletrônico)"
	case "65":
		return "NFC-e (Nota Fiscal de Consumidor Eletrônica)"
	case "67":
		return "CT-e OS (Conhecimento de Transporte Eletrônico para Outros Serviços)"
	default:
		return m.Descricao
	}
}

// GetSigla retorna a sigla do modelo
func (m *ModeloDocumentoFiscal) GetSigla() string {
	switch m.Codigo {
	case "55":
		return "NF-e"
	case "57":
		return "CT-e"
	case "58":
		return "MDF-e"
	case "59":
		return "CF-e"
	case "63":
		return "BP-e"
	case "65":
		return "NFC-e"
	case "67":
		return "CT-e OS"
	default:
		return m.Codigo
	}
}
