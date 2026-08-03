// internal/database/callbacks.go
package database

import (
	"gorm.io/gorm"

	"github.com/openerp/backend/internal/appcontext"
)

// RegisterAuditCallbacks registra callbacks globais para auditoria
func RegisterAuditCallbacks(db *gorm.DB) {
	// Callback para BeforeCreate
	db.Callback().Create().Before("gorm:create").Register("audit:before_create", func(tx *gorm.DB) {
		if tx.Statement.Schema != nil {
			// Verifica se o model tem campos de auditoria
			_, hasCreatedBy := tx.Statement.Schema.FieldsByName["CreatedBy"]
			_, hasUpdatedBy := tx.Statement.Schema.FieldsByName["UpdatedBy"]

			if hasCreatedBy && hasUpdatedBy {
				userID := appcontext.GetUserID(tx.Statement.Context)
				if userID != 0 {
					// ✅ CORRETO - SetColumn NÃO retorna valor
					tx.Statement.SetColumn("CreatedBy", userID)
					tx.Statement.SetColumn("UpdatedBy", userID)
				}
			}
		}
	})

	// Callback para BeforeUpdate
	db.Callback().Update().Before("gorm:update").Register("audit:before_update", func(tx *gorm.DB) {
		if tx.Statement.Schema != nil {
			_, hasUpdatedBy := tx.Statement.Schema.FieldsByName["UpdatedBy"]
			if hasUpdatedBy {
				userID := appcontext.GetUserID(tx.Statement.Context)
				if userID != 0 {
					// ✅ CORRETO - SetColumn NÃO retorna valor
					tx.Statement.SetColumn("UpdatedBy", userID)
				}
			}
		}
	})
}
