package models

import (
	"time"
)

// ============================================================
// MODEL: ProdutoGrupo
// ============================================================

type ProdutoGrupo struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID                 int      `gorm:"column:prog_id;primaryKey;autoIncrement" json:"id"`
	Descricao          string   `gorm:"column:prog_descricao;type:varchar(255);not null" json:"descricao"`
	Situacao           int      `gorm:"column:prog_situacao;not null;default:0" json:"situacao"`
	ImpressoraID       *int     `gorm:"column:print_id" json:"impressora_id,omitempty"`
	Cor                *string  `gorm:"column:prog_cor;type:varchar(200)" json:"cor,omitempty"`
	Imagem             []byte   `gorm:"column:prog_imagem;type:longblob" json:"-"`
	PercentualComissao *float64 `gorm:"column:prog_percentualcomissao;type:decimal(5,2)" json:"percentual_comissao,omitempty"`
	VisivelFrenteCaixa int      `gorm:"column:prog_visivelfrentecaixa;not null;default:0" json:"visivel_frente_caixa"`
	Agenda             *int     `gorm:"column:prog_agenda" json:"agenda,omitempty"`
	ControleLote       *int     `gorm:"column:prog_controlelote" json:"controle_lote,omitempty"`

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
	//SubGrupos  []ProdutoSubgrupo `gorm:"foreignKey:ProdutoGrupoID;references:ID" json:"sub_grupos,omitempty"`
	Produtos []Produto `gorm:"foreignKey:ProdutoGrupoID;references:ID" json:"produtos,omitempty"`
	// Impressora *Impressora       `gorm:"foreignKey:ImpressoraID;references:print_id" json:"impressora,omitempty"`
}

func (ProdutoGrupo) TableName() string {
	return "produto_grupo"
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsActive verifica se o grupo está ativo
func (m *ProdutoGrupo) IsActive() bool {
	return m.Situacao == 1
}

// IsDeleted verifica se foi deletado logicamente
func (m *ProdutoGrupo) IsDeleted() bool {
	return m.DeletedAt != nil
}

// SoftDelete realiza a exclusão lógica
func (m *ProdutoGrupo) SoftDelete() {
	now := time.Now()
	m.DeletedAt = &now
	m.Situacao = 2
}

// IsVisivelNoCaixa verifica se está visível no frente de caixa
func (m *ProdutoGrupo) IsVisivelNoCaixa() bool {
	return m.VisivelFrenteCaixa == 1
}

// HasComissao verifica se o grupo tem comissão
func (m *ProdutoGrupo) HasComissao() bool {
	return m.PercentualComissao != nil && *m.PercentualComissao > 0
}

// HasControleLote verifica se o grupo tem controle de lote
func (m *ProdutoGrupo) HasControleLote() bool {
	return m.ControleLote != nil && *m.ControleLote == 1
}

// HasAgenda verifica se o grupo tem agenda
func (m *ProdutoGrupo) HasAgenda() bool {
	return m.Agenda != nil && *m.Agenda == 1
}
