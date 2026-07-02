package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: DocumentoVendaVeiculoMotorista
// ============================================================

type DocumentoVendaVeiculoMotorista struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	DocumentoVendaID   int      `gorm:"column:ddv_id;primaryKey" json:"documento_venda_id"`
	DocumentoVendaItem int      `gorm:"column:dvi_item;primaryKey" json:"documento_venda_item"`
	Item               int      `gorm:"column:dddvvm_item;primaryKey" json:"item"`
	VeiculoID          *int     `gorm:"column:vei_id" json:"veiculo_id,omitempty"`
	EntidadeID         int      `gorm:"column:ent_id;not null" json:"entidade_id"`
	Descricao          *string  `gorm:"column:ddvvm_descricao;type:varchar(2000)" json:"descricao,omitempty"`
	PercentualComissao *float64 `gorm:"column:ddvvm_percentualcomissao;type:decimal(15,4)" json:"percentual_comissao,omitempty"`

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
	Veiculo  *Veiculos `gorm:"foreignKey:VeiculoID;references:vei_id" json:"veiculo,omitempty"`
	Entidade *Entidade `gorm:"foreignKey:EntidadeID;references:ent_id" json:"entidade,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (DocumentoVendaVeiculoMotorista) TableName() string {
	return "documento_venda_veiculo_motorista"
}

func (d *DocumentoVendaVeiculoMotorista) BeforeCreate(tx *gorm.DB) error {
	if d.CreatedBy == nil {
		d.CreatedBy = new(int)
		*d.CreatedBy = 0
	}
	if d.UpdatedBy == nil {
		d.UpdatedBy = new(int)
		*d.UpdatedBy = 0
	}
	return nil
}

func (d *DocumentoVendaVeiculoMotorista) BeforeUpdate(tx *gorm.DB) error {
	if d.UpdatedBy == nil {
		d.UpdatedBy = new(int)
		*d.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o registro foi deletado logicamente
func (d *DocumentoVendaVeiculoMotorista) IsDeleted() bool {
	return d.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (d *DocumentoVendaVeiculoMotorista) SoftDelete() {
	now := time.Now()
	d.DeletedAt = &now
}

// HasVeiculo verifica se possui veículo associado
func (d *DocumentoVendaVeiculoMotorista) HasVeiculo() bool {
	return d.VeiculoID != nil && *d.VeiculoID > 0
}

// HasDescricao verifica se possui descrição
func (d *DocumentoVendaVeiculoMotorista) HasDescricao() bool {
	return d.Descricao != nil && *d.Descricao != ""
}

// HasPercentualComissao verifica se possui percentual de comissão
func (d *DocumentoVendaVeiculoMotorista) HasPercentualComissao() bool {
	return d.PercentualComissao != nil && *d.PercentualComissao > 0
}

// GetPercentualComissaoOrDefault retorna o percentual de comissão ou 0 se não definido
func (d *DocumentoVendaVeiculoMotorista) GetPercentualComissaoOrDefault() float64 {
	if d.HasPercentualComissao() {
		return *d.PercentualComissao
	}
	return 0
}

// GetDescricaoOrDefault retorna a descrição ou string vazia se não definida
func (d *DocumentoVendaVeiculoMotorista) GetDescricaoOrDefault() string {
	if d.HasDescricao() {
		return *d.Descricao
	}
	return ""
}

// CalcularComissao calcula o valor da comissão baseado no valor do item
func (d *DocumentoVendaVeiculoMotorista) CalcularComissao(valorItem float64) float64 {
	if d.HasPercentualComissao() {
		return valorItem * (*d.PercentualComissao / 100)
	}
	return 0
}

// GetNomeMotorista retorna o nome do motorista da entidade associada
func (d *DocumentoVendaVeiculoMotorista) GetNomeMotorista() string {
	if d.Entidade != nil {
		return d.Entidade.RazaoSocial
	}
	return ""
}

// GetPlacaVeiculo retorna a placa do veículo associado
func (d *DocumentoVendaVeiculoMotorista) GetPlacaVeiculo() string {
	if d.Veiculo != nil {
		return d.Veiculo.Placa
	}
	return ""
}

// GetNomeCompleto retorna o nome completo com motorista e veículo
func (d *DocumentoVendaVeiculoMotorista) GetNomeCompleto() string {
	nome := d.GetNomeMotorista()
	if d.HasVeiculo() {
		placa := d.GetPlacaVeiculo()
		if placa != "" {
			if nome != "" {
				return nome + " - " + placa
			}
			return placa
		}
	}
	return nome
}
