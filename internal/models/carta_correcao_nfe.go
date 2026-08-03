package models

import (
	"time"
)

// ============================================================
// MODEL: CartaCorrecaoNFe
// ============================================================

type CartaCorrecaoNFe struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID              int       `gorm:"column:ccnfe_id;primaryKey;autoIncrement" json:"id"`
	NotaFiscalID    int       `gorm:"column:ntf_id;not null" json:"nota_fiscal_id"`
	ChaveNFe        string    `gorm:"column:ccnfe_chavenfe;type:varchar(255);not null" json:"chave_nfe"`
	DataHoraEvento  time.Time `gorm:"column:ccnfe_datahoraevento;type:datetime;not null" json:"data_hora_evento"`
	NumeroSeq       int       `gorm:"column:ccnfe_numeroseq;not null" json:"numero_seq"`
	Texto           string    `gorm:"column:ccnfe_texto;type:text;not null" json:"texto"`
	XML             *string   `gorm:"column:ccnfe_xml;type:text" json:"xml,omitempty"`
	NumeroProtocolo *string   `gorm:"column:ccnfe_numeroprotocolo;type:varchar(255)" json:"numero_protocolo,omitempty"`
	CodigoStatus    *int      `gorm:"column:ccnfe_codigostatus" json:"codigo_status,omitempty"`
	MotivoStatus    *string   `gorm:"column:ccnfe_motivostatus;type:varchar(255)" json:"motivo_status,omitempty"`

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
	NotaFiscal *NotaFiscal `gorm:"foreignKey:NotaFiscalID;references:ntf_id" json:"nota_fiscal,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (CartaCorrecaoNFe) TableName() string {
	return "carta_correcao_nfe"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se a carta de correção foi deletada logicamente
func (c *CartaCorrecaoNFe) IsDeleted() bool {
	return c.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (c *CartaCorrecaoNFe) SoftDelete() {
	now := time.Now()
	c.DeletedAt = &now
}

// HasProtocolo verifica se possui protocolo
func (c *CartaCorrecaoNFe) HasProtocolo() bool {
	return c.NumeroProtocolo != nil && *c.NumeroProtocolo != ""
}

// HasXML verifica se possui XML
func (c *CartaCorrecaoNFe) HasXML() bool {
	return c.XML != nil && *c.XML != ""
}

// IsSucesso verifica se a carta de correção foi bem-sucedida
func (c *CartaCorrecaoNFe) IsSucesso() bool {
	return c.CodigoStatus != nil && *c.CodigoStatus == 100
}
