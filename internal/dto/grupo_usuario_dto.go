package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/openerp/backend/internal/constants"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// REQUESTS
// ============================================================

// GrupoUsuarioRequest representa a requisição para criar/atualizar um grupo de usuário
type GrupoUsuarioRequest struct {
	Descricao       string `json:"descricao" binding:"required,max=100"`
	EmpresaFilialID int    `json:"empresa_filial_id" binding:"required"`
	Situacao        int    `json:"situacao,omitempty"`
	CreatedBy       *int   `json:"created_by,omitempty"`
	UpdatedBy       *int   `json:"updated_by,omitempty"`
}

// ============================================================
// RESPONSES
// ============================================================

// GrupoUsuarioResponse representa a resposta de um grupo de usuário
type GrupoUsuarioResponse struct {
	ID                int              `json:"id"`
	Descricao         string           `json:"descricao"`
	EmpresaFilialID   int              `json:"empresa_filial_id"`
	EmpresaFilialNome string           `json:"empresa_filial_nome,omitempty"`
	Situacao          constants.Status `json:"situacao"`
	SituacaoLabel     string           `json:"situacao_label"`
	CreatedAt         string           `json:"created_at"`
	UpdatedAt         string           `json:"updated_at"`
	CreatedBy         *int             `json:"created_by,omitempty"`
	UpdatedBy         *int             `json:"updated_by,omitempty"`
}

// ============================================================
// RESPONSE LISTA
// ============================================================

// GrupoUsuarioListResponse representa a resposta de listagem de grupos de usuário
type GrupoUsuarioListResponse struct {
	Items      []GrupoUsuarioResponse `json:"items"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	TotalPages int                    `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE CONVERSÃO
// ============================================================

// ToModel converte GrupoUsuarioRequest para models.GrupoUsuario.
func (r *GrupoUsuarioRequest) ToModel() (*models.GrupoUsuario, error) {
	if r == nil {
		return nil, nil
	}
	grupo := &models.GrupoUsuario{}
	if err := utils.MapToModel(r, grupo); err != nil {
		return nil, fmt.Errorf("erro ao mapear DTO para modelo de grupo de usuário: %w", err)
	}
	// Definir situação padrão se não informada
	if r.Situacao == 0 {
		grupo.Situacao = constants.StatusAtivo
	}
	return grupo, nil
}

// FromModel converte models.GrupoUsuario para GrupoUsuarioResponse.
func (r *GrupoUsuarioResponse) FromModel(grupo *models.GrupoUsuario) {
	if grupo == nil {
		return
	}
	// Tenta mapear diretamente, se falhar, usa fallback manual
	if err := utils.MapToDTO(grupo, r); err != nil {
		r.fromModelFallback(grupo)
	}

	r.SituacaoLabel = grupo.Situacao.String()

	// Preencher nome da filial se o relacionamento estiver carregado
	if grupo.EmpresaFilial != nil {
		r.EmpresaFilialNome = grupo.EmpresaFilial.Nome
	}

	r.CreatedAt = utils.FormatDateTime(grupo.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(grupo.UpdatedAt)
}

// fromModelFallback é o fallback manual caso o mapper falhe
func (r *GrupoUsuarioResponse) fromModelFallback(grupo *models.GrupoUsuario) {
	r.ID = grupo.ID
	r.Descricao = grupo.Descricao
	r.EmpresaFilialID = grupo.EmpresaFilialID
	r.Situacao = grupo.Situacao
	r.CreatedBy = grupo.CreatedBy
	r.UpdatedBy = grupo.UpdatedBy

	r.SituacaoLabel = constants.Status(grupo.Situacao).String()
	if grupo.EmpresaFilial != nil {
		r.EmpresaFilialNome = grupo.EmpresaFilial.Nome
	}
	r.CreatedAt = utils.FormatDateTime(grupo.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(grupo.UpdatedAt)
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o GrupoUsuarioRequest.
func (r *GrupoUsuarioRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
