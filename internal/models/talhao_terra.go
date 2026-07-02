package models

import (
	"time"

	"gorm.io/gorm"
)

// ============================================================
// MODEL: TalhaoTerra
// ============================================================

type TalhaoTerra struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID          int      `gorm:"column:tat_id;primaryKey;autoIncrement" json:"id"`
	Descricao   string   `gorm:"column:tat_descricao;type:varchar(255);not null" json:"descricao"`
	Situacao    int      `gorm:"column:tat_situacao;not null" json:"situacao"` // 1-Ativo, 2-Inativo, 3-Bloqueado
	QuantHectar *float64 `gorm:"column:tat_quanthectar;type:decimal(15,4)" json:"quant_hectar,omitempty"`

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
	DocumentosVenda []DocumentoVenda `gorm:"foreignKey:TalhaoID;references:ID" json:"documentos_venda,omitempty"`
	ProdutosEstoque []ProdutoEstoque `gorm:"foreignKey:TalhaoID;references:ID" json:"produtos_estoque,omitempty"`
	Requisicoes     []Requisicao     `gorm:"foreignKey:TalhaoID;references:ID" json:"requisicoes,omitempty"`
}

func (TalhaoTerra) TableName() string {
	return "talhao_terra"
}

func (t *TalhaoTerra) BeforeCreate(tx *gorm.DB) error {
	if t.CreatedBy == nil {
		t.CreatedBy = new(int)
		*t.CreatedBy = 0
	}
	if t.UpdatedBy == nil {
		t.UpdatedBy = new(int)
		*t.UpdatedBy = 0
	}
	return nil
}

func (t *TalhaoTerra) BeforeUpdate(tx *gorm.DB) error {
	if t.UpdatedBy == nil {
		t.UpdatedBy = new(int)
		*t.UpdatedBy = 0
	}
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsActive verifica se o talhão está ativo
func (t *TalhaoTerra) IsActive() bool {
	return t.Situacao == 1
}

// IsDeleted verifica se o talhão foi deletado logicamente
func (t *TalhaoTerra) IsDeleted() bool {
	return t.DeletedAt != nil
}

// IsInactive verifica se o talhão está inativo
func (t *TalhaoTerra) IsInactive() bool {
	return t.Situacao == 2
}

// IsBlocked verifica se o talhão está bloqueado
func (t *TalhaoTerra) IsBlocked() bool {
	return t.Situacao == 3
}

// SoftDelete realiza a exclusão lógica
func (t *TalhaoTerra) SoftDelete() {
	now := time.Now()
	t.DeletedAt = &now
	t.Situacao = 2 // Marca como inativo
}

// HasQuantHectar verifica se o talhão tem quantidade de hectares definida
func (t *TalhaoTerra) HasQuantHectar() bool {
	return t.QuantHectar != nil && *t.QuantHectar > 0
}
