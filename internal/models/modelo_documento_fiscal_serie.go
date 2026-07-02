package models

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/constants"
)

// ============================================================
// MODEL: ModeloDocumentoFiscalSerie
// ============================================================

type ModeloDocumentoFiscalSerie struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	EmpresaFilialID             int              `gorm:"column:emf_id;primaryKey" json:"empresa_filial_id"`
	ModeloDocumentoFiscalCodigo string           `gorm:"column:mdf_codigo;type:varchar(3);primaryKey" json:"modelo_documento_fiscal_codigo"`
	Serie                       int              `gorm:"column:mdfs_serie;primaryKey" json:"serie"`
	UltimoNumero                int              `gorm:"column:mdfs_ultimonumero;not null;default:0" json:"ultimo_numero"`
	Situacao                    constants.Status `gorm:"column:mdfs_situacao;not null;default:1" json:"situacao"`

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
	ModeloDocumentoFiscal *ModeloDocumentoFiscal `gorm:"foreignKey:EmpresaFilialID,ModeloDocumentoFiscalCodigo;references:emf_id,mdf_codigo" json:"modelo_documento_fiscal,omitempty"`
	EmpresaFilial         *EmpresaFilial         `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (ModeloDocumentoFiscalSerie) TableName() string {
	return "modelo_documento_fiscal_serie"
}

func (m *ModeloDocumentoFiscalSerie) BeforeCreate(tx *gorm.DB) error {
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

func (m *ModeloDocumentoFiscalSerie) BeforeUpdate(tx *gorm.DB) error {
	if m.UpdatedBy == nil {
		m.UpdatedBy = new(int)
		*m.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsActive verifica se a série está ativa
func (m *ModeloDocumentoFiscalSerie) IsActive() bool {
	return m.Situacao == constants.StatusAtivo
}

// IsDeleted verifica se a série foi deletada logicamente
func (m *ModeloDocumentoFiscalSerie) IsDeleted() bool {
	return m.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (m *ModeloDocumentoFiscalSerie) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
	m.Situacao = constants.StatusInativo
}

// GetDescricao retorna a descrição do modelo pai
func (m *ModeloDocumentoFiscalSerie) GetDescricao() string {
	if m.ModeloDocumentoFiscal != nil {
		return m.ModeloDocumentoFiscal.Descricao
	}
	return ""
}

// GetSigla retorna a sigla do modelo pai
func (m *ModeloDocumentoFiscalSerie) GetSigla() string {
	if m.ModeloDocumentoFiscal != nil {
		return m.ModeloDocumentoFiscal.GetSigla()
	}
	return m.ModeloDocumentoFiscalCodigo
}

// GetNumeroCompleto retorna o número completo com modelo e série
func (m *ModeloDocumentoFiscalSerie) GetNumeroCompleto(numero int) string {
	return fmt.Sprintf("%s/%d/%d", m.ModeloDocumentoFiscalCodigo, m.Serie, numero)
}

// GetProximoNumero retorna o próximo número e incrementa
func (m *ModeloDocumentoFiscalSerie) GetProximoNumero() int {
	proximo := m.UltimoNumero + 1
	m.UltimoNumero = proximo
	return proximo
}

// GetUltimoNumero retorna o último número utilizado
func (m *ModeloDocumentoFiscalSerie) GetUltimoNumero() int {
	return m.UltimoNumero
}

// IncrementarNumero incrementa o último número em 1
func (m *ModeloDocumentoFiscalSerie) IncrementarNumero() {
	m.UltimoNumero++
}

// IncrementarNumeroBy incrementa o último número por um valor específico
func (m *ModeloDocumentoFiscalSerie) IncrementarNumeroBy(valor int) {
	if valor > 0 {
		m.UltimoNumero += valor
	}
}

// ResetNumero reseta o último número para 0
func (m *ModeloDocumentoFiscalSerie) ResetNumero() {
	m.UltimoNumero = 0
}

// SetUltimoNumero define o último número com validação
func (m *ModeloDocumentoFiscalSerie) SetUltimoNumero(numero int) error {
	if numero < 0 {
		return errors.New("número não pode ser negativo")
	}
	m.UltimoNumero = numero
	return nil
}

// IsNFe verifica se é série de NFe (Modelo 55)
func (m *ModeloDocumentoFiscalSerie) IsNFe() bool {
	return m.ModeloDocumentoFiscalCodigo == "55"
}

// IsNFCe verifica se é série de NFCe (Modelo 65)
func (m *ModeloDocumentoFiscalSerie) IsNFCe() bool {
	return m.ModeloDocumentoFiscalCodigo == "65"
}

// IsCTe verifica se é série de CTe (Modelo 57)
func (m *ModeloDocumentoFiscalSerie) IsCTe() bool {
	return m.ModeloDocumentoFiscalCodigo == "57"
}

// IsCTeOS verifica se é série de CTe OS (Modelo 67)
func (m *ModeloDocumentoFiscalSerie) IsCTeOS() bool {
	return m.ModeloDocumentoFiscalCodigo == "67"
}

// IsMDFe verifica se é série de MDFe (Modelo 58)
func (m *ModeloDocumentoFiscalSerie) IsMDFe() bool {
	return m.ModeloDocumentoFiscalCodigo == "58"
}

// IsCFe verifica se é série de CFe (Modelo 59)
func (m *ModeloDocumentoFiscalSerie) IsCFe() bool {
	return m.ModeloDocumentoFiscalCodigo == "59"
}

// IsBPe verifica se é série de BPe (Modelo 63)
func (m *ModeloDocumentoFiscalSerie) IsBPe() bool {
	return m.ModeloDocumentoFiscalCodigo == "63"
}

// IsEletronico verifica se o modelo da série é eletrônico
func (m *ModeloDocumentoFiscalSerie) IsEletronico() bool {
	eletronicos := map[string]bool{
		"55": true,
		"57": true,
		"58": true,
		"59": true,
		"63": true,
		"65": true,
		"67": true,
	}
	return eletronicos[m.ModeloDocumentoFiscalCodigo]
}

// GetNomeCompleto retorna o nome completo com modelo, série e descrição
func (m *ModeloDocumentoFiscalSerie) GetNomeCompleto() string {
	desc := m.GetDescricao()
	sigla := m.GetSigla()
	if desc != "" {
		return fmt.Sprintf("%s - Série %d", desc, m.Serie)
	}
	return fmt.Sprintf("%s/Série %d", sigla, m.Serie)
}

// GetDescricaoSerie retorna uma descrição amigável da série
func (m *ModeloDocumentoFiscalSerie) GetDescricaoSerie() string {
	sigla := m.GetSigla()
	return fmt.Sprintf("%s - Série %d (Último: %d)", sigla, m.Serie, m.UltimoNumero)
}
