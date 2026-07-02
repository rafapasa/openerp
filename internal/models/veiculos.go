package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: Veiculos
// ============================================================

type Veiculos struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID              int     `gorm:"column:vei_id;primaryKey;autoIncrement" json:"id"`
	EmpresaFilialID int     `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	Placa           string  `gorm:"column:vei_placa;type:varchar(8);not null" json:"placa"`
	PlacaUF         *string `gorm:"column:vei_placa_uf;type:char(2)" json:"placa_uf,omitempty"`
	Marca           *string `gorm:"column:vei_marca;type:varchar(100)" json:"marca,omitempty"`
	Modelo          *string `gorm:"column:vei_modelo;type:varchar(100)" json:"modelo,omitempty"`
	Renavam         *string `gorm:"column:vei_renava;type:varchar(100)" json:"renavam,omitempty"`
	RNTC            *string `gorm:"column:vei_rntc;type:varchar(255)" json:"rntc,omitempty"`
	Guincho         *int    `gorm:"column:vei_guincho" json:"guincho,omitempty"` // 1 - sim, 0 - não

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

	// Relacionamentos onde o veículo é referenciado
	Carregamentos       []Carregamento                   `gorm:"foreignKey:VeiculoID;references:vei_id" json:"carregamentos,omitempty"`
	DocumentoVendaItens []DocumentoVendaVeiculoMotorista `gorm:"foreignKey:VeiculoID;references:vei_id" json:"documento_venda_itens,omitempty"`
	NotaFiscais         []NotaFiscal                     `gorm:"foreignKey:VeiculoTransportadorID;references:vei_id" json:"notas_fiscais,omitempty"`
	Entidades           []VeiculoEntidade                `gorm:"foreignKey:VeiculoID;references:vei_id" json:"entidades,omitempty"`
	Manutencoes         []VeiculoManutencao              `gorm:"foreignKey:VeiculoID;references:vei_id" json:"manutencoes,omitempty"`
	OrdensServico       []OrdemServico                   `gorm:"foreignKey:VeiculoID;references:vei_id" json:"ordens_servico,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (Veiculos) TableName() string {
	return "veiculos"
}

func (v *Veiculos) BeforeCreate(tx *gorm.DB) error {
	if v.CreatedBy == nil {
		v.CreatedBy = new(int)
		*v.CreatedBy = 0
	}
	if v.UpdatedBy == nil {
		v.UpdatedBy = new(int)
		*v.UpdatedBy = 0
	}
	return nil
}

func (v *Veiculos) BeforeUpdate(tx *gorm.DB) error {
	if v.UpdatedBy == nil {
		v.UpdatedBy = new(int)
		*v.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o veículo foi deletado logicamente
func (v *Veiculos) IsDeleted() bool {
	return v.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (v *Veiculos) SoftDelete() {
	now := time.Now()
	v.DeletedAt = &now
}

// HasMarca verifica se o veículo possui marca definida
func (v *Veiculos) HasMarca() bool {
	return v.Marca != nil && *v.Marca != ""
}

// HasModelo verifica se o veículo possui modelo definido
func (v *Veiculos) HasModelo() bool {
	return v.Modelo != nil && *v.Modelo != ""
}

// HasRenavam verifica se o veículo possui Renavam definido
func (v *Veiculos) HasRenavam() bool {
	return v.Renavam != nil && *v.Renavam != ""
}

// HasRNTC verifica se o veículo possui RNTC definido
func (v *Veiculos) HasRNTC() bool {
	return v.RNTC != nil && *v.RNTC != ""
}

// IsGuincho verifica se o veículo é um guincho
func (v *Veiculos) IsGuincho() bool {
	return v.Guincho != nil && *v.Guincho == 1
}

// GetPlacaCompleta retorna a placa com UF
func (v *Veiculos) GetPlacaCompleta() string {
	if v.PlacaUF != nil && *v.PlacaUF != "" {
		return v.Placa + "-" + *v.PlacaUF
	}
	return v.Placa
}

// GetNomeCompleto retorna o nome completo do veículo
func (v *Veiculos) GetNomeCompleto() string {
	nome := v.Placa
	if v.Marca != nil && *v.Marca != "" {
		nome += " - " + *v.Marca
	}
	if v.Modelo != nil && *v.Modelo != "" {
		nome += " " + *v.Modelo
	}
	return nome
}

// HasCarregamentos verifica se o veículo possui carregamentos associados
func (v *Veiculos) HasCarregamentos() bool {
	return len(v.Carregamentos) > 0
}

// HasNotaFiscais verifica se o veículo possui notas fiscais associadas
func (v *Veiculos) HasNotaFiscais() bool {
	return len(v.NotaFiscais) > 0
}

// HasManutencoes verifica se o veículo possui manutenções associadas
func (v *Veiculos) HasManutencoes() bool {
	return len(v.Manutencoes) > 0
}

// HasEntidades verifica se o veículo possui entidades associadas
func (v *Veiculos) HasEntidades() bool {
	return len(v.Entidades) > 0
}

// SafeToDelete verifica se o veículo pode ser excluído
func (v *Veiculos) SafeToDelete() bool {
	if v.HasCarregamentos() {
		return false
	}
	if v.HasNotaFiscais() {
		return false
	}
	if v.HasManutencoes() {
		return false
	}
	if v.HasEntidades() {
		return false
	}
	if len(v.DocumentoVendaItens) > 0 {
		return false
	}
	if len(v.OrdensServico) > 0 {
		return false
	}

	return true
}
