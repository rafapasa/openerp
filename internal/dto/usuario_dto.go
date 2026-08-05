// ============================================================
// FILE: usuario_dto.go
// PACKAGE: dto
// ============================================================

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

// UsuarioRequest representa a requisição para criar/atualizar um usuário
type UsuarioRequest struct {
	// ============================================================
	// DADOS PRINCIPAIS
	// ============================================================
	GrupoUsuarioID int              `json:"grupo_usuario_id" binding:"required"`
	Nome           string           `json:"nome" binding:"required"`
	Login          string           `json:"login" binding:"required"`
	Senha          *string          `json:"senha,omitempty"`
	Situacao       constants.Status `json:"situacao,omitempty"`
	Observacoes    *string          `json:"observacoes,omitempty"`

	// ============================================================
	// CONFIGURAÇÕES
	// ============================================================
	AlterarColGrid *constants.SimNao `json:"alterar_col_grid,omitempty"`

	// ============================================================
	// CONFIGURAÇÕES DE E-MAIL (SMTP)
	// ============================================================
	EmailSMTP     *string           `json:"email_smtp,omitempty"`
	PortaSMTP     *int              `json:"porta_smtp,omitempty"`
	ServidorSMTP  *string           `json:"servidor_smtp,omitempty"`
	UsarTLS       *constants.SimNao `json:"usar_tls,omitempty"`
	UsarSSL       *constants.SimNao `json:"usar_ssl,omitempty"`
	UsuarioSMTP   *string           `json:"usuario_smtp,omitempty"`
	ExigirSenhaDV *constants.SimNao `json:"exigir_senha_dv,omitempty"`

	// ============================================================
	// RELACIONAMENTOS
	// ============================================================
	UsuarioFiliais []UsuarioFilialRequest `json:"usuario_filiais,omitempty"`

	// ============================================================
	// USUÁRIO (para auditoria)
	// ============================================================
	CreatedBy *int `json:"created_by,omitempty"`
	UpdatedBy *int `json:"updated_by,omitempty"`
}

// UsuarioFilialRequest representa a relação usuário-filial no request
type UsuarioFilialRequest struct {
	EmpresaFilialID int `json:"empresa_filial_id" binding:"required"`
}

// ============================================================
// RESPONSES
// ============================================================

// UsuarioResponse representa a resposta de um usuário
type UsuarioResponse struct {
	// ============================================================
	// DADOS PRINCIPAIS
	// ============================================================
	ID               int              `json:"id"`
	GrupoUsuarioID   int              `json:"grupo_usuario_id"`
	GrupoUsuarioNome string           `json:"grupo_usuario_nome,omitempty"`
	Nome             string           `json:"nome"`
	Login            string           `json:"login"`
	Situacao         constants.Status `json:"situacao"`
	SituacaoLabel    string           `json:"situacao_label"`
	Observacoes      *string          `json:"observacoes,omitempty"`

	// ============================================================
	// CONFIGURAÇÕES
	// ============================================================
	AlterarColGrid      *constants.SimNao `json:"alterar_col_grid,omitempty"`
	AlterarColGridLabel string            `json:"alterar_col_grid_label,omitempty"`

	// ============================================================
	// CONFIGURAÇÕES DE E-MAIL (SMTP)
	// ============================================================
	EmailSMTP          *string           `json:"email_smtp,omitempty"`
	PortaSMTP          *int              `json:"porta_smtp,omitempty"`
	ServidorSMTP       *string           `json:"servidor_smtp,omitempty"`
	UsarTLS            *constants.SimNao `json:"usar_tls,omitempty"`
	UsarTLSLabel       string            `json:"usar_tls_label,omitempty"`
	UsarSSL            *constants.SimNao `json:"usar_ssl,omitempty"`
	UsarSSLLabel       string            `json:"usar_ssl_label,omitempty"`
	UsuarioSMTP        *string           `json:"usuario_smtp,omitempty"`
	ExigirSenhaDV      *constants.SimNao `json:"exigir_senha_dv,omitempty"`
	ExigirSenhaDVLabel string            `json:"exigir_senha_dv_label,omitempty"`

	// ============================================================
	// RELACIONAMENTOS
	// ============================================================
	UsuarioFiliais []UsuarioFilialResponse `json:"usuario_filiais,omitempty"`

	// ============================================================
	// AUDITORIA
	// ============================================================
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	CreatedBy *int   `json:"created_by,omitempty"`
	UpdatedBy *int   `json:"updated_by,omitempty"`
}

// UsuarioFilialResponse representa a resposta da relação usuário-filial
type UsuarioFilialResponse struct {
	EmpresaFilialID   int    `json:"empresa_filial_id"`
	EmpresaFilialNome string `json:"empresa_filial_nome,omitempty"`
}

// ============================================================
// RESPONSE LISTA
// ============================================================

// UsuarioListResponse representa a resposta de listagem de usuários
type UsuarioListResponse struct {
	Items      []UsuarioResponse `json:"items"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE CONVERSÃO (TO MODEL)
// ============================================================

// ToModel converte UsuarioRequest para models.Usuario usando mapper
func (r *UsuarioRequest) ToModel() (*models.Usuario, error) {
	if r == nil {
		return nil, nil
	}

	usuario := &models.Usuario{}

	// 1. Usar o mapper para copiar campos (mapeamento automático)
	if err := utils.MapToModel(r, usuario); err != nil {
		return nil, err
	}

	// 2. Tratamentos especiais que o mapper não cobre
	// Se situação não informada, definir como ativo
	if r.Situacao == 0 {
		usuario.Situacao = constants.StatusAtivo
	}

	// Se alterar_col_grid não informado, definir como Sim (1)
	if r.AlterarColGrid == nil {
		defaultVal := constants.SimNaoSim
		usuario.AlterarColGrid = &defaultVal
	}

	return usuario, nil
}

// FromModel converte models.Usuario para UsuarioResponse usando mapper
func (r *UsuarioResponse) FromModel(usuario *models.Usuario) (*UsuarioResponse, error) {
	if usuario == nil {
		return nil, nil
	}

	// 1. Usar o mapper para copiar campos (mapeamento automático)
	if err := utils.MapToDTO(usuario, r); err != nil {
		// Se o mapper falhar, usar fallback manual
		return r.fromModelFallback(usuario), nil
	}

	// 2. Preencher campos calculados (labels)
	r.SituacaoLabel = usuario.Situacao.String()

	// Labels para campos do tipo constants.SimNao
	if usuario.AlterarColGrid != nil {
		r.AlterarColGridLabel = usuario.AlterarColGrid.String()
	}

	if usuario.UsarTLS != nil {
		r.UsarTLSLabel = usuario.UsarTLS.String()
	}

	if usuario.UsarSSL != nil {
		r.UsarSSLLabel = usuario.UsarSSL.String()
	}

	if usuario.ExigirSenhaDV != nil {
		r.ExigirSenhaDVLabel = usuario.ExigirSenhaDV.String()
	}

	// 3. Nome do grupo de usuário (se carregado)
	if usuario.GrupoUsuario != nil {
		r.GrupoUsuarioNome = usuario.GrupoUsuario.Descricao
	}

	// 4. Converter filiais
	if len(usuario.UsuarioFiliais) > 0 {
		filiaisResponse := make([]UsuarioFilialResponse, len(usuario.UsuarioFiliais))
		for i, uf := range usuario.UsuarioFiliais {
			filialResp := UsuarioFilialResponse{
				EmpresaFilialID: uf.EmpresaFilialID,
			}
			if uf.EmpresaFilial != nil {
				filialResp.EmpresaFilialNome = uf.EmpresaFilial.Nome
			}
			filiaisResponse[i] = filialResp
		}
		r.UsuarioFiliais = filiaisResponse
	}

	// 5. Formatar datas
	r.CreatedAt = utils.FormatDateTime(usuario.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(usuario.UpdatedAt)

	return r, nil
}

// ============================================================
// FALLBACK (caso o mapper falhe)
// ============================================================

// fromModelFallback é o fallback manual caso o mapper falhe
func (r *UsuarioResponse) fromModelFallback(usuario *models.Usuario) *UsuarioResponse {
	// Mapeamento manual campo por campo (seguro)
	r.ID = usuario.ID
	r.GrupoUsuarioID = usuario.GrupoUsuarioID
	r.Nome = usuario.Nome
	r.Login = usuario.Login
	r.Situacao = usuario.Situacao
	r.SituacaoLabel = usuario.Situacao.String()
	r.Observacoes = usuario.Observacoes
	r.AlterarColGrid = usuario.AlterarColGrid
	r.EmailSMTP = usuario.EmailSMTP
	r.PortaSMTP = usuario.PortaSMTP
	r.ServidorSMTP = usuario.ServidorSMTP
	r.UsarTLS = usuario.UsarTLS
	r.UsarSSL = usuario.UsarSSL
	r.UsuarioSMTP = usuario.UsuarioSMTP
	r.ExigirSenhaDV = usuario.ExigirSenhaDV
	r.CreatedBy = usuario.CreatedBy
	r.UpdatedBy = usuario.UpdatedBy

	// Labels para campos do tipo constants.SimNao
	if usuario.AlterarColGrid != nil {
		r.AlterarColGridLabel = usuario.AlterarColGrid.String()
	}
	if usuario.UsarTLS != nil {
		r.UsarTLSLabel = usuario.UsarTLS.String()
	}
	if usuario.UsarSSL != nil {
		r.UsarSSLLabel = usuario.UsarSSL.String()
	}
	if usuario.ExigirSenhaDV != nil {
		r.ExigirSenhaDVLabel = usuario.ExigirSenhaDV.String()
	}

	// Nome do grupo
	if usuario.GrupoUsuario != nil {
		r.GrupoUsuarioNome = usuario.GrupoUsuario.Descricao
	}

	// Filiais
	if len(usuario.UsuarioFiliais) > 0 {
		filiaisResponse := make([]UsuarioFilialResponse, len(usuario.UsuarioFiliais))
		for i, uf := range usuario.UsuarioFiliais {
			filialResp := UsuarioFilialResponse{
				EmpresaFilialID: uf.EmpresaFilialID,
			}
			if uf.EmpresaFilial != nil {
				filialResp.EmpresaFilialNome = uf.EmpresaFilial.Nome
			}
			filiaisResponse[i] = filialResp
		}
		r.UsuarioFiliais = filiaisResponse
	}

	// Datas
	r.CreatedAt = utils.FormatDateTime(usuario.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(usuario.UpdatedAt)

	return r
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o UsuarioRequest
func (r *UsuarioRequest) Validate() error {
	// 1. Validação com go-playground/validator
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}

	// 2. Validar campos obrigatórios usando utils
	if err := utils.ValidateMandatoryFields(r); err != nil {
		return err
	}

	// 3. Validar situação (se informada)
	if r.Situacao > 0 {
		if err := r.Situacao.IsValid(); err != nil {
			return err
		}
	}

	// 4. Validar campos do tipo constants.SimNao (se informados)
	if r.AlterarColGrid != nil {
		if err := r.AlterarColGrid.IsValid(); err != nil {
			return err
		}
	}

	if r.UsarTLS != nil {
		if err := r.UsarTLS.IsValid(); err != nil {
			return err
		}
	}

	if r.UsarSSL != nil {
		if err := r.UsarSSL.IsValid(); err != nil {
			return err
		}
	}

	if r.ExigirSenhaDV != nil {
		if err := r.ExigirSenhaDV.IsValid(); err != nil {
			return err
		}
	}

	// 5. Validar porta SMTP (se informada)
	if r.PortaSMTP != nil && (*r.PortaSMTP < 1 || *r.PortaSMTP > 65535) {
		return fmt.Errorf("porta_smtp deve estar entre 1 e 65535, recebido: %d", *r.PortaSMTP)
	}

	// 6. Validar filiais
	for i, uf := range r.UsuarioFiliais {
		if uf.EmpresaFilialID <= 0 {
			return fmt.Errorf("usuario_filiais[%d].empresa_filial_id deve ser maior que 0", i)
		}
	}

	// 7. Validar senha (se for criação de novo usuário)
	// Se ID for 0 (novo usuário) e senha não for informada
	if r.Senha != nil && len(*r.Senha) > 0 && len(*r.Senha) < 6 {
		return fmt.Errorf("senha deve ter no mínimo 6 caracteres")
	}

	return nil
}
