package models

import (
	"time"

	"github.com/openerp/backend/internal/constants"
)

// ============================================================
// MODEL: Usuario
// ============================================================

type Usuario struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID             int               `gorm:"column:usu_id;primaryKey;autoIncrement" json:"id"`
	GrupoUsuarioID int               `gorm:"column:gpu_id;not null" json:"grupo_usuario_id"`
	Nome           string            `gorm:"column:usu_nome;type:varchar(100);not null" json:"nome"`
	Login          string            `gorm:"column:usu_login;type:varchar(20);not null;unique" json:"login"`
	Senha          string            `gorm:"column:usu_senha;type:varchar(100);not null" json:"-"`
	Situacao       constants.Status  `gorm:"column:usu_situiacao;not null;default:1" json:"situacao"`
	Observacoes    *string           `gorm:"column:usu_observacoes;type:text" json:"observacoes,omitempty"`
	SenhaExclusao  *string           `gorm:"column:usu_senhaexclusao;type:varchar(100)" json:"-"`
	AlterarColGrid *constants.SimNao `gorm:"column:usu_alterarcolgrid;default:1" json:"alterar_col_grid,omitempty"`
	EmailSMTP      *string           `gorm:"column:usu_emailsmtp;type:varchar(200)" json:"email_smtp,omitempty"`
	PortaSMTP      *int              `gorm:"column:usu_portasmtp" json:"porta_smtp,omitempty"`
	SenhaSMTP      *string           `gorm:"column:usu_senhasmtp;type:varchar(20)" json:"-"`
	ServidorSMTP   *string           `gorm:"column:usu_servidorsmtp;type:varchar(200)" json:"servidor_smtp,omitempty"`
	UsarTLS        *constants.SimNao `gorm:"column:usu_usartls" json:"usar_tls,omitempty"`
	UsarSSL        *constants.SimNao `gorm:"column:usu_usarssl" json:"usar_ssl,omitempty"`
	UsuarioSMTP    *string           `gorm:"column:usu_usuariosmtp;type:varchar(200)" json:"usuario_smtp,omitempty"`
	ExigirSenhaDV  *constants.SimNao `gorm:"column:usu_exigirsenhadv" json:"exigir_senha_dv,omitempty"`

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
	GrupoUsuario   *GrupoUsuario   `gorm:"foreignKey:GrupoUsuarioID;references:gpu_id" json:"grupo_usuario,omitempty"`
	UsuarioFiliais []UsuarioFilial `gorm:"foreignKey:UsuarioID;references:ID" json:"usuario_filiais,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (Usuario) TableName() string {
	return "usuario"
}

// CORRIGIDO: adicionado *gorm.DB

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (u *Usuario) IsActive() bool {
	return u.Situacao == constants.StatusAtivo
}

func (u *Usuario) IsDeleted() bool {
	return u.DeletedAt != nil
}

func (u *Usuario) SoftDelete() {
	now := time.Now()
	u.DeletedAt = &now
	u.Situacao = constants.StatusInativo
}

func (u *Usuario) GetEmpresaFilialID() int {
	if len(u.UsuarioFiliais) > 0 {
		return u.UsuarioFiliais[0].EmpresaFilialID
	}
	return 0
}

func (u *Usuario) GetEmpresaFilialIDs() []int {
	var ids []int
	for _, uf := range u.UsuarioFiliais {
		ids = append(ids, uf.EmpresaFilialID)
	}
	return ids
}

func (u *Usuario) HasEmpresaFilial(emfID int) bool {
	for _, uf := range u.UsuarioFiliais {
		if uf.EmpresaFilialID == emfID {
			return true
		}
	}
	return false
}

func (u *Usuario) HasEmailSMTP() bool {
	return u.EmailSMTP != nil && *u.EmailSMTP != ""
}

func (u *Usuario) HasServidorSMTP() bool {
	return u.ServidorSMTP != nil && *u.ServidorSMTP != ""
}

func (u *Usuario) IsDeletavel() bool {
	// Verifica se o usuário pode ser deletado
	// TODO: Implementar regras de negócio
	return true
}
