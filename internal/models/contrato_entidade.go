package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: ContratoEntidade
// ============================================================

type ContratoEntidade struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID                     int  `gorm:"column:cent_id;primaryKey;autoIncrement" json:"id"`
	ContratoEntidadeTipoID int  `gorm:"column:cet_id;not null" json:"contrato_entidade_tipo_id"`
	EntidadeID             int  `gorm:"column:ent_id;not null" json:"entidade_id"`
	PortadorID             int  `gorm:"column:por_id;not null" json:"portador_id"`
	EmpresaFilialID        *int `gorm:"column:emf_id" json:"empresa_filial_id,omitempty"`

	DataInicial    time.Time  `gorm:"column:cent_datainicila;type:date;not null" json:"data_inicial"`
	DataFinal      *time.Time `gorm:"column:cent_datafinal;type:date" json:"data_final,omitempty"`
	DataVencimento time.Time  `gorm:"column:cent_datavencimento;type:date;not null" json:"data_vencimento"`

	Tipo          int        `gorm:"column:cent_tipo;not null" json:"tipo"` // 1 = pagamento / 2 = recebimento
	Valor         float64    `gorm:"column:cent_valor;type:decimal(15,2);not null;default:0" json:"valor"`
	ValorDesconto *float64   `gorm:"column:cent_valordesconto;type:decimal(15,2);default:0" json:"valor_desconto,omitempty"`
	DataDesconto  *time.Time `gorm:"column:cent_datadesconto;type:date" json:"data_desconto,omitempty"`

	Formato        int  `gorm:"column:cent_formato;not null;default:1" json:"formato"` // 1-mensal 2-periodo fixo
	NumeroParcelas *int `gorm:"column:cent_numeroparcelas" json:"numero_parcelas,omitempty"`

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
	ContratoEntidadeTipo *ContratoEntidadeTipo       `gorm:"foreignKey:ContratoEntidadeTipoID;references:cet_id" json:"contrato_entidade_tipo,omitempty"`
	Entidade             *Entidade                   `gorm:"foreignKey:EntidadeID;references:ent_id" json:"entidade,omitempty"`
	Portador             *Portador                   `gorm:"foreignKey:PortadorID;references:por_id" json:"portador,omitempty"`
	EmpresaFilial        *EmpresaFilial              `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
	Pagamentos           []ContratoEntidadePagamento `gorm:"foreignKey:ContratoEntidadeID;references:cent_id" json:"pagamentos,omitempty"`
	Produtos             []ContratoEntidadeProduto   `gorm:"foreignKey:ContratoEntidadeID;references:cent_id" json:"produtos,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (ContratoEntidade) TableName() string {
	return "contrato_entidade"
}

func (c *ContratoEntidade) BeforeCreate(tx *gorm.DB) error {
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

func (c *ContratoEntidade) BeforeUpdate(tx *gorm.DB) error {
	if c.UpdatedBy == nil {
		c.UpdatedBy = new(int)
		*c.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsDeleted verifica se o contrato foi deletado logicamente
func (c *ContratoEntidade) IsDeleted() bool {
	return c.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (c *ContratoEntidade) SoftDelete() {
	now := time.Now()
	c.DeletedAt = &now
}

// IsPagamento verifica se o contrato é do tipo pagamento
func (c *ContratoEntidade) IsPagamento() bool {
	return c.Tipo == 1
}

// IsRecebimento verifica se o contrato é do tipo recebimento
func (c *ContratoEntidade) IsRecebimento() bool {
	return c.Tipo == 2
}

// IsMensal verifica se o formato é mensal
func (c *ContratoEntidade) IsMensal() bool {
	return c.Formato == 1
}

// IsPeriodoFixo verifica se o formato é período fixo
func (c *ContratoEntidade) IsPeriodoFixo() bool {
	return c.Formato == 2
}

// HasDesconto verifica se o contrato possui desconto
func (c *ContratoEntidade) HasDesconto() bool {
	return c.ValorDesconto != nil && *c.ValorDesconto > 0
}

// GetValorLiquido retorna o valor líquido do contrato (valor - desconto)
func (c *ContratoEntidade) GetValorLiquido() float64 {
	valor := c.Valor
	if c.HasDesconto() {
		valor -= *c.ValorDesconto
	}
	return valor
}

// IsVencido verifica se o contrato está vencido
func (c *ContratoEntidade) IsVencido() bool {
	return time.Now().After(c.DataVencimento)
}

// GetPagamentosCount retorna a quantidade de pagamentos associados
func (c *ContratoEntidade) GetPagamentosCount() int {
	return len(c.Pagamentos)
}

// GetProdutosCount retorna a quantidade de produtos associados
func (c *ContratoEntidade) GetProdutosCount() int {
	return len(c.Produtos)
}
