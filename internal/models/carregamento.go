package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: Carregamento
// ============================================================

type Carregamento struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID                  int        `gorm:"column:car_id;primaryKey;autoIncrement" json:"id"`
	EntidadeID          *int       `gorm:"column:ent_id" json:"entidade_id,omitempty"` // Cliente/Destinatário
	EntidadeMotoristaID *int       `gorm:"column:ent_id_mot" json:"entidade_motorista_id,omitempty"`
	VeiculoID           *int       `gorm:"column:vei_id" json:"veiculo_id,omitempty"`
	EmpresaFilialID     int        `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	DataPrevisao        time.Time  `gorm:"column:car_dataprevisao;type:date;not null" json:"data_previsao"`
	DataCarregamento    *time.Time `gorm:"column:car_datacarregamento;type:date" json:"data_carregamento,omitempty"`
	Situacao            int        `gorm:"column:car_situacao;not null" json:"situacao"` // 1 - aberto, 2 - em atividade, 8 - finalizado, 9 - cancelado
	GrupoEntidadeID     *int       `gorm:"column:gpe_id" json:"grupo_entidade_id,omitempty"`
	DataIni             *time.Time `gorm:"column:car_dataini;type:date" json:"data_ini,omitempty"`
	DataFim             *time.Time `gorm:"column:car_datafim;type:date" json:"data_fim,omitempty"`

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
	Entidade          *Entidade      `gorm:"foreignKey:EntidadeID;references:ent_id" json:"entidade,omitempty"`
	EntidadeMotorista *Entidade      `gorm:"foreignKey:EntidadeMotoristaID;references:ent_id" json:"entidade_motorista,omitempty"`
	Veiculo           *Veiculos      `gorm:"foreignKey:VeiculoID;references:vei_id" json:"veiculo,omitempty"`
	EmpresaFilial     *EmpresaFilial `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
	GrupoEntidade     *GrupoEntidade `gorm:"foreignKey:GrupoEntidadeID;references:gpe_id" json:"grupo_entidade,omitempty"`

	Documentos []CarregamentoDocumento `gorm:"foreignKey:CarregamentoID;references:car_id" json:"documentos,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (Carregamento) TableName() string {
	return "carregamento"
}

func (c *Carregamento) BeforeCreate(tx *gorm.DB) error {
	if c.CreatedBy == nil {
		c.CreatedBy = new(int)
		*c.CreatedBy = 0
	}
	if c.UpdatedBy == nil {
		c.UpdatedBy = new(int)
		*c.UpdatedBy = 0
	}
	return nil
}

func (c *Carregamento) BeforeUpdate(tx *gorm.DB) error {
	if c.UpdatedBy == nil {
		c.UpdatedBy = new(int)
		*c.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o carregamento foi deletado logicamente
func (c *Carregamento) IsDeleted() bool {
	return c.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (c *Carregamento) SoftDelete() {
	now := time.Now()
	c.DeletedAt = &now
}

// IsAberto verifica se o carregamento está aberto
func (c *Carregamento) IsAberto() bool {
	return c.Situacao == 1
}

// IsEmAtividade verifica se o carregamento está em atividade
func (c *Carregamento) IsEmAtividade() bool {
	return c.Situacao == 2
}

// IsFinalizado verifica se o carregamento está finalizado
func (c *Carregamento) IsFinalizado() bool {
	return c.Situacao == 8
}

// IsCancelado verifica se o carregamento está cancelado
func (c *Carregamento) IsCancelado() bool {
	return c.Situacao == 9
}

// HasEntidade verifica se possui entidade (cliente)
func (c *Carregamento) HasEntidade() bool {
	return c.EntidadeID != nil && *c.EntidadeID > 0
}

// HasMotorista verifica se possui motorista
func (c *Carregamento) HasMotorista() bool {
	return c.EntidadeMotoristaID != nil && *c.EntidadeMotoristaID > 0
}

// HasVeiculo verifica se possui veículo
func (c *Carregamento) HasVeiculo() bool {
	return c.VeiculoID != nil && *c.VeiculoID > 0
}

// HasGrupoEntidade verifica se possui grupo de entidade
func (c *Carregamento) HasGrupoEntidade() bool {
	return c.GrupoEntidadeID != nil && *c.GrupoEntidadeID > 0
}

// HasDocumentos verifica se possui documentos associados
func (c *Carregamento) HasDocumentos() bool {
	return len(c.Documentos) > 0
}

// GetDocumentosCount retorna a quantidade de documentos
func (c *Carregamento) GetDocumentosCount() int {
	return len(c.Documentos)
}

// GetSituacaoDescricao retorna a descrição da situação
func (c *Carregamento) GetSituacaoDescricao() string {
	switch c.Situacao {
	case 1:
		return "Aberto"
	case 2:
		return "Em Atividade"
	case 8:
		return "Finalizado"
	case 9:
		return "Cancelado"
	default:
		return "Desconhecido"
	}
}

// IsFinalizadoOuCancelado verifica se está finalizado ou cancelado
func (c *Carregamento) IsFinalizadoOuCancelado() bool {
	return c.IsFinalizado() || c.IsCancelado()
}

// GetStatusColor retorna a cor do status para UI
func (c *Carregamento) GetStatusColor() string {
	switch c.Situacao {
	case 1:
		return "blue" // Aberto
	case 2:
		return "orange" // Em Atividade
	case 8:
		return "green" // Finalizado
	case 9:
		return "red" // Cancelado
	default:
		return "gray"
	}
}

// SafeToDelete verifica se o carregamento pode ser excluído
func (c *Carregamento) SafeToDelete() bool {
	if c.HasDocumentos() {
		return false
	}
	if c.IsFinalizado() || c.IsCancelado() {
		return false
	}
	return true
}
