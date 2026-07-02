package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: OrdemServico
// ============================================================

type OrdemServico struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID               int        `gorm:"column:os_id;primaryKey;autoIncrement" json:"id"`
	EntidadeID       int        `gorm:"column:ent_id;not null" json:"entidade_id"`
	VeiculoID        *int       `gorm:"column:vei_id" json:"veiculo_id,omitempty"`
	DataCadastro     time.Time  `gorm:"column:os_datacadastro;type:date;not null" json:"data_cadastro"`
	DataPrevisao     *time.Time `gorm:"column:os_dataprevisao;type:date" json:"data_previsao,omitempty"`
	DataInicio       *time.Time `gorm:"column:os_datainicio;type:datetime" json:"data_inicio,omitempty"`
	DataTermino      *time.Time `gorm:"column:os_datatermino;type:datetime" json:"data_termino,omitempty"`
	Situacao         int        `gorm:"column:os_situacao;not null;default:0" json:"situacao"`
	DescricaoCliente *string    `gorm:"column:os_descricaocliente;type:text" json:"descricao_cliente,omitempty"`
	DescricaoTecnica *string    `gorm:"column:os_descricaotecnica;type:text" json:"descricao_tecnica,omitempty"`

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
	Entidade        *Entidade                    `gorm:"foreignKey:EntidadeID;references:ent_id" json:"entidade,omitempty"`
	Veiculo         *Veiculos                    `gorm:"foreignKey:VeiculoID;references:vei_id" json:"veiculo,omitempty"`
	DocumentosVenda []DocumentoVendaOrdemServico `gorm:"foreignKey:OrdemServicoID;references:os_id" json:"documentos_venda,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (OrdemServico) TableName() string {
	return "ordem_servico"
}

func (o *OrdemServico) BeforeCreate(tx *gorm.DB) error {
	if o.CreatedBy == nil {
		o.CreatedBy = new(int)
		*o.CreatedBy = 0
	}
	if o.UpdatedBy == nil {
		o.UpdatedBy = new(int)
		*o.UpdatedBy = 0
	}
	return nil
}

func (o *OrdemServico) BeforeUpdate(tx *gorm.DB) error {
	if o.UpdatedBy == nil {
		o.UpdatedBy = new(int)
		*o.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se a ordem de serviço foi deletada logicamente
func (o *OrdemServico) IsDeleted() bool {
	return o.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (o *OrdemServico) SoftDelete() {
	now := time.Now()
	o.DeletedAt = &now
}

// IsAberta verifica se a ordem está aberta (situação 0)
func (o *OrdemServico) IsAberta() bool {
	return o.Situacao == 0
}

// IsEmAndamento verifica se a ordem está em andamento (situação 1)
func (o *OrdemServico) IsEmAndamento() bool {
	return o.Situacao == 1
}

// IsConcluida verifica se a ordem está concluída (situação 2)
func (o *OrdemServico) IsConcluida() bool {
	return o.Situacao == 2
}

// IsCancelada verifica se a ordem está cancelada (situação 9)
func (o *OrdemServico) IsCancelada() bool {
	return o.Situacao == 9
}

// HasVeiculo verifica se possui veículo associado
func (o *OrdemServico) HasVeiculo() bool {
	return o.VeiculoID != nil && *o.VeiculoID > 0
}

// HasDataPrevisao verifica se possui data de previsão
func (o *OrdemServico) HasDataPrevisao() bool {
	return o.DataPrevisao != nil
}

// HasDataInicio verifica se possui data de início
func (o *OrdemServico) HasDataInicio() bool {
	return o.DataInicio != nil
}

// HasDataTermino verifica se possui data de término
func (o *OrdemServico) HasDataTermino() bool {
	return o.DataTermino != nil
}

// HasDescricaoCliente verifica se possui descrição do cliente
func (o *OrdemServico) HasDescricaoCliente() bool {
	return o.DescricaoCliente != nil && *o.DescricaoCliente != ""
}

// HasDescricaoTecnica verifica se possui descrição técnica
func (o *OrdemServico) HasDescricaoTecnica() bool {
	return o.DescricaoTecnica != nil && *o.DescricaoTecnica != ""
}

// HasDocumentosVenda verifica se possui documentos de venda associados
func (o *OrdemServico) HasDocumentosVenda() bool {
	return len(o.DocumentosVenda) > 0
}

// GetDocumentosVendaCount retorna a quantidade de documentos de venda
func (o *OrdemServico) GetDocumentosVendaCount() int {
	return len(o.DocumentosVenda)
}

// GetSituacaoDescricao retorna a descrição da situação
func (o *OrdemServico) GetSituacaoDescricao() string {
	switch o.Situacao {
	case 0:
		return "Aberta"
	case 1:
		return "Em Andamento"
	case 2:
		return "Concluída"
	case 9:
		return "Cancelada"
	default:
		return "Desconhecido"
	}
}

// GetStatusColor retorna a cor do status para UI
func (o *OrdemServico) GetStatusColor() string {
	switch o.Situacao {
	case 0:
		return "blue" // Aberta
	case 1:
		return "orange" // Em Andamento
	case 2:
		return "green" // Concluída
	case 9:
		return "red" // Cancelada
	default:
		return "gray"
	}
}

// GetTempoDecorrido retorna o tempo decorrido desde o início
func (o *OrdemServico) GetTempoDecorrido() *time.Duration {
	if o.DataInicio == nil {
		return nil
	}
	duracao := time.Since(*o.DataInicio)
	return &duracao
}

// GetTempoTotal retorna o tempo total entre início e término
func (o *OrdemServico) GetTempoTotal() *time.Duration {
	if o.DataInicio == nil || o.DataTermino == nil {
		return nil
	}
	duracao := o.DataTermino.Sub(*o.DataInicio)
	return &duracao
}

// IsVencida verifica se a ordem está vencida (previsão passou e não concluída)
func (o *OrdemServico) IsVencida() bool {
	if !o.HasDataPrevisao() || o.IsConcluida() || o.IsCancelada() {
		return false
	}
	return time.Now().After(*o.DataPrevisao)
}

// SafeToDelete verifica se a ordem pode ser excluída
func (o *OrdemServico) SafeToDelete() bool {
	if o.HasDocumentosVenda() {
		return false
	}
	if o.IsEmAndamento() || o.IsConcluida() {
		return false
	}
	return true
}
