package models

import (
	"time"
)

// BaseModel contém campos comuns a todas as tabelas
type BaseModel struct {
	CreatedAt time.Time  `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime;index" json:"deleted_at,omitempty"`
	CreatedBy *int       `gorm:"column:created_by;type:int" json:"created_by,omitempty"`
	UpdatedBy *int       `gorm:"column:updated_by;type:int" json:"updated_by,omitempty"`
}

// SoftDelete implementa exclusão lógica
func (b *BaseModel) SoftDelete() {
	now := time.Now()
	b.DeletedAt = &now
}

// IsDeleted verifica se o registro foi excluído logicamente
func (b *BaseModel) IsDeleted() bool {
	return b.DeletedAt != nil
}

// SetCreatedBy define quem criou o registro
func (b *BaseModel) SetCreatedBy(userID int) {
	b.CreatedBy = &userID
}

// SetUpdatedBy define quem atualizou o registro
func (b *BaseModel) SetUpdatedBy(userID int) {
	b.UpdatedBy = &userID
}
