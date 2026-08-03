package models

import (
	"time"

	"github.com/openerp/backend/internal/constants"
)

// ============================================================
// MODEL: OrdemProducao
// ============================================================

type OrdemProducao struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID                   int              `gorm:"column:orp_id;primaryKey;autoIncrement" json:"id"`
	AssistenteProducaoID *int             `gorm:"column:asp_id" json:"assistente_producao_id,omitempty"`
	EmpresaFilialID      int              `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	OrdemProducaoPaiID   *int             `gorm:"column:orp_orp_id" json:"ordem_producao_pai_id,omitempty"`
	Tipo                 int              `gorm:"column:orp_tipo;not null" json:"tipo"`
	DataCadastro         time.Time        `gorm:"column:orp_datacadastro;type:date;not null" json:"data_cadastro"`
	DataPrevisaoTermino  *time.Time       `gorm:"column:orp_dataprevisaotermino;type:date" json:"data_previsao_termino,omitempty"`
	Situacao             constants.Status `gorm:"column:orp_situacao;not null" json:"situacao"`
	DataInicio           *time.Time       `gorm:"column:orp_datainicio;type:date" json:"data_inicio,omitempty"`
	DataFim              *time.Time       `gorm:"column:orp_datafim;type:date" json:"data_fim,omitempty"`
	Observacao           *string          `gorm:"column:orp_observacao;type:text" json:"observacao,omitempty"`
	Realizacao           int              `gorm:"column:orp_realizacao;not null" json:"realizacao"` // 1 - interna, 2 - externa

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
	AssistenteProducao *AssistenteProducao `gorm:"foreignKey:AssistenteProducaoID;references:asp_id" json:"assistente_producao,omitempty"`
	EmpresaFilial      *EmpresaFilial      `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
	OrdemProducaoPai   *OrdemProducao      `gorm:"foreignKey:OrdemProducaoPaiID;references:orp_id" json:"ordem_producao_pai,omitempty"`

	Itens        []OrdemProducaoItem       `gorm:"foreignKey:OrdemProducaoID;references:orp_id" json:"itens,omitempty"`
	Componentes  []OrdemProducaoComponente `gorm:"foreignKey:OrdemProducaoID;references:orp_id" json:"componentes,omitempty"`
	Explosao     []OrdemProducaoExplosao   `gorm:"foreignKey:OrdemProducaoID;references:orp_id" json:"explosao,omitempty"`
	Requisicoes  []Requisicao              `gorm:"foreignKey:OrdemProducaoID;references:orp_id" json:"requisicoes,omitempty"`
	OrdensFilhas []OrdemProducao           `gorm:"foreignKey:OrdemProducaoPaiID;references:orp_id" json:"ordens_filhas,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (OrdemProducao) TableName() string {
	return "ordem_producao"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsActive verifica se a ordem de produção está ativa
func (o *OrdemProducao) IsActive() bool {
	return o.Situacao == constants.StatusAtivo
}

// IsDeleted verifica se a ordem de produção foi deletada logicamente
func (o *OrdemProducao) IsDeleted() bool {
	return o.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (o *OrdemProducao) SoftDelete() {
	now := time.Now()
	o.DeletedAt = &now
	o.Situacao = constants.StatusInativo
}

// IsInactive verifica se a ordem de produção está inativa
func (o *OrdemProducao) IsInactive() bool {
	return o.Situacao == constants.StatusInativo
}

// IsInterna verifica se a produção é interna
func (o *OrdemProducao) IsInterna() bool {
	return o.Realizacao == 1
}

// IsExterna verifica se a produção é externa
func (o *OrdemProducao) IsExterna() bool {
	return o.Realizacao == 2
}

// HasAssistenteProducao verifica se possui assistente de produção
func (o *OrdemProducao) HasAssistenteProducao() bool {
	return o.AssistenteProducaoID != nil && *o.AssistenteProducaoID > 0
}

// HasOrdemPai verifica se possui ordem pai
func (o *OrdemProducao) HasOrdemPai() bool {
	return o.OrdemProducaoPaiID != nil && *o.OrdemProducaoPaiID > 0
}

// HasOrdensFilhas verifica se possui ordens filhas
func (o *OrdemProducao) HasOrdensFilhas() bool {
	return len(o.OrdensFilhas) > 0
}

// HasItens verifica se possui itens
func (o *OrdemProducao) HasItens() bool {
	return len(o.Itens) > 0
}

// HasComponentes verifica se possui componentes
func (o *OrdemProducao) HasComponentes() bool {
	return len(o.Componentes) > 0
}

// HasExplosao verifica se possui explosão
func (o *OrdemProducao) HasExplosao() bool {
	return len(o.Explosao) > 0
}

// HasRequisicoes verifica se possui requisições
func (o *OrdemProducao) HasRequisicoes() bool {
	return len(o.Requisicoes) > 0
}

// GetItensCount retorna a quantidade de itens
func (o *OrdemProducao) GetItensCount() int {
	return len(o.Itens)
}

// GetComponentesCount retorna a quantidade de componentes
func (o *OrdemProducao) GetComponentesCount() int {
	return len(o.Componentes)
}

// GetOrdensFilhasCount retorna a quantidade de ordens filhas
func (o *OrdemProducao) GetOrdensFilhasCount() int {
	return len(o.OrdensFilhas)
}

// IsFinalizada verifica se a ordem está finalizada
func (o *OrdemProducao) IsFinalizada() bool {
	return o.DataFim != nil
}

// IsEmAndamento verifica se a ordem está em andamento
func (o *OrdemProducao) IsEmAndamento() bool {
	return o.DataInicio != nil && o.DataFim == nil
}

// GetTempoDecorrido retorna o tempo decorrido desde o início
func (o *OrdemProducao) GetTempoDecorrido() *time.Duration {
	if o.DataInicio == nil {
		return nil
	}
	duracao := time.Since(*o.DataInicio)
	return &duracao
}

// SafeToDelete verifica se a ordem pode ser excluída
func (o *OrdemProducao) SafeToDelete() bool {
	if o.HasOrdensFilhas() {
		return false
	}
	if o.HasRequisicoes() {
		return false
	}
	if o.IsEmAndamento() {
		return false
	}
	return true
}
