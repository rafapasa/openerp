package models

import (
	"time"

	"github.com/openerp/backend/internal/constants"
)

// ============================================================
// MODEL: Produto
// ============================================================

type Produto struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID                int              `gorm:"column:pro_id;primaryKey;autoIncrement" json:"id"`
	EmpresaFilialID   int              `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	TipoProdutoID     int              `gorm:"column:ptp_id;not null" json:"tipo_produto_id"`
	ProdutoGrupoID    *int             `gorm:"column:prog_id" json:"produto_grupo_id,omitempty"`
	ProdutoSubgrupoID *int             `gorm:"column:prosg_id" json:"produto_subgrupo_id,omitempty"` // CORRIGIDO: nome
	MarcaID           *int             `gorm:"column:promar_id" json:"marca_id,omitempty"`
	ModeloID          *int             `gorm:"column:prom_id" json:"modelo_id,omitempty"`
	SerieID           *int             `gorm:"column:pros_id" json:"serie_id,omitempty"`
	EspecieID         *int             `gorm:"column:proesp_id" json:"especie_id,omitempty"`
	NcmNumero         int              `gorm:"column:ncm_numero;not null" json:"ncm_numero"`
	Codigo            int              `gorm:"column:pro_codigo;not null;unique" json:"codigo"`
	Nome              string           `gorm:"column:pro_nome;type:varchar(80);not null" json:"nome"`
	Resumo            string           `gorm:"column:pro_resumo;type:varchar(80);not null" json:"resumo"`
	Situacao          constants.Status `gorm:"column:pro_situacao;not null;default:3" json:"situacao"` // CORRIGIDO: constants.Status
	CustoCompra       *float64         `gorm:"column:pro_custocompra;type:decimal(15,4)" json:"custo_compra,omitempty"`
	CodigoBarras      *string          `gorm:"column:pro_codigobarra;type:varchar(255)" json:"codigo_barras,omitempty"`
	PesoBruto         *float64         `gorm:"column:pro_pesobruto;type:decimal(15,4)" json:"peso_bruto,omitempty"`
	PesoLiquido       *float64         `gorm:"column:pro_pesoliquido;type:decimal(15,4)" json:"peso_liquido,omitempty"`
	Altura            *float64         `gorm:"column:pro_altura;type:decimal(15,4)" json:"altura,omitempty"`
	Largura           *float64         `gorm:"column:pro_largura;type:decimal(15,4)" json:"largura,omitempty"`
	Comprimento       *float64         `gorm:"column:pro_comprimento;type:decimal(15,4)" json:"comprimento,omitempty"`
	EstoqueMinimo     *float64         `gorm:"column:pro_estoqueminimo;type:decimal(15,4)" json:"estoque_minimo,omitempty"`
	LoteEconomico     *float64         `gorm:"column:pro_loteeconomico;type:decimal(15,4)" json:"lote_economico,omitempty"`
	Referencia        *string          `gorm:"column:pro_referencia;type:varchar(255)" json:"referencia,omitempty"`
	Referencia2       *string          `gorm:"column:pro_referencia2;type:varchar(255)" json:"referencia2,omitempty"`
	Referencia3       *string          `gorm:"column:pro_referencia3;type:varchar(255)" json:"referencia3,omitempty"`
	Referencia4       *string          `gorm:"column:pro_referencia4;type:varchar(255)" json:"referencia4,omitempty"`
	Descricao         *string          `gorm:"column:pro_descricao;type:text" json:"descricao,omitempty"`
	DataAlteracao     *time.Time       `gorm:"column:pro_dataalteracao;type:datetime" json:"data_alteracao,omitempty"`
	DescontoMaximo    *float64         `gorm:"column:pro_descontomaximo;type:decimal(15,4)" json:"desconto_maximo,omitempty"`

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
	TipoProduto     *TipoProduto         `gorm:"foreignKey:TipoProdutoID;references:ptp_id" json:"tipo_produto,omitempty"`
	ProdutoGrupo    *ProdutoGrupo        `gorm:"foreignKey:ProdutoGrupoID;references:prog_id" json:"produto_grupo,omitempty"`
	ProdutoSubgrupo *ProdutoSubgrupo     `gorm:"foreignKey:ProdutoSubgrupoID;references:prosg_id" json:"produto_subgrupo,omitempty"`
	Marca           *ProdutoMarca        `gorm:"foreignKey:MarcaID;references:promar_id" json:"marca,omitempty"`
	Modelo          *ProdutoModelo       `gorm:"foreignKey:ModeloID;references:prom_id" json:"modelo,omitempty"`
	Serie           *ProdutoSerie        `gorm:"foreignKey:SerieID;references:pros_id" json:"serie,omitempty"`
	Especie         *ProdutoEspecie      `gorm:"foreignKey:EspecieID;references:proesp_id" json:"especie,omitempty"`
	ItensPedido     []DocumentoVendaItem `gorm:"foreignKey:ProdutoID;references:pro_id" json:"itens_pedido,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (Produto) TableName() string {
	return "produto"
}

// CORRIGIDO: adicionado *gorm.DB

// CORRIGIDO: adicionado *gorm.DB

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsActive verifica se o produto está ativo
func (p *Produto) IsActive() bool {
	return p.Situacao == constants.StatusAtivo
}

// IsDeleted verifica se o produto foi deletado logicamente
func (p *Produto) IsDeleted() bool {
	return p.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (p *Produto) SoftDelete() {
	now := time.Now()
	p.DeletedAt = &now
	p.Situacao = constants.StatusInativo
}

// IsInactive verifica se o produto está inativo
func (p *Produto) IsInactive() bool {
	return p.Situacao == constants.StatusInativo
}

// HasEstoqueMinimo verifica se o produto tem estoque mínimo definido
func (p *Produto) HasEstoqueMinimo() bool {
	return p.EstoqueMinimo != nil && *p.EstoqueMinimo > 0
}

// HasCodigoBarras verifica se o produto tem código de barras
func (p *Produto) HasCodigoBarras() bool {
	return p.CodigoBarras != nil && *p.CodigoBarras != ""
}

// GetNomeCompleto retorna o nome completo do produto com código
func (p *Produto) GetNomeCompleto() string {
	return p.Nome + " (" + string(rune(p.Codigo)) + ")"
}

// SafeToDelete verifica se o produto pode ser excluído
// Verifica se há registros relacionados (pedidos, estoque, etc.)
func (p *Produto) SafeToDelete() bool {
	// TODO: Implementar verificação de dependências
	// - Verificar se tem itens em pedidos
	// - Verificar se tem movimentações de estoque
	// - Verificar se está em tabelas de preço
	// Por enquanto retorna true
	return true
}
