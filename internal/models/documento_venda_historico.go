package models

import (
	"time"
)

// ============================================================
// MODEL: DocumentoVendaHistorico
// ============================================================

type DocumentoVendaHistorico struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	DocumentoVendaID int       `gorm:"column:ddv_id;primaryKey" json:"documento_venda_id"`
	Item             int       `gorm:"column:ddvh_item;primaryKey" json:"item"`
	UsuarioID        int       `gorm:"column:usu_id;not null" json:"usuario_id"`
	FluxoID          int       `gorm:"column:flu_id;not null" json:"fluxo_id"`
	FluxoSequencia   int       `gorm:"column:fls_sequencia;not null" json:"fluxo_sequencia"`
	DataHistorico    time.Time `gorm:"column:ddvh_datahistorico;type:date;not null" json:"data_historico"`
	Descricao        string    `gorm:"column:ddvh_descricao;type:varchar(2000);not null" json:"descricao"`
	Motivo           *string   `gorm:"column:ddvh_motivo;type:varchar(2000)" json:"motivo,omitempty"`

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
	DocumentoVenda *DocumentoVenda `gorm:"foreignKey:DocumentoVendaID;references:ddv_id" json:"documento_venda,omitempty"`
	Usuario        *Usuario        `gorm:"foreignKey:UsuarioID;references:usu_id" json:"usuario,omitempty"`
	// FluxoSetor     *FluxoSetor     `gorm:"foreignKey:FluxoID,FluxoSequencia;references:flu_id,fls_sequencia" json:"fluxo_setor,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (DocumentoVendaHistorico) TableName() string {
	return "documento_venda_historico"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o histórico foi deletado logicamente
func (d *DocumentoVendaHistorico) IsDeleted() bool {
	return d.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (d *DocumentoVendaHistorico) SoftDelete() {
	now := time.Now()
	d.DeletedAt = &now
}

// HasMotivo verifica se possui motivo
func (d *DocumentoVendaHistorico) HasMotivo() bool {
	return d.Motivo != nil && *d.Motivo != ""
}

// GetMotivoOrDefault retorna o motivo ou string vazia
func (d *DocumentoVendaHistorico) GetMotivoOrDefault() string {
	if d.HasMotivo() {
		return *d.Motivo
	}
	return ""
}

// GetDataHistoricoFormatada retorna a data do histórico formatada
func (d *DocumentoVendaHistorico) GetDataHistoricoFormatada() string {
	return d.DataHistorico.Format("02/01/2006 15:04:05")
}

// GetDataHistoricoISO retorna a data do histórico no formato ISO
func (d *DocumentoVendaHistorico) GetDataHistoricoISO() string {
	return d.DataHistorico.Format("2006-01-02T15:04:05Z07:00")
}
