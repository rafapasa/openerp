package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// CONSTANTES
// ============================================================

const (
	TipoContatoTelefone  = 1
	TipoContatoCelular   = 2
	TipoContatoEmail     = 3
	TipoContatoWhatsApp  = 4
	TipoContatoSite      = 5
	TipoContatoFacebook  = 6
	TipoContatoInstagram = 7
)

// ============================================================
// REQUESTS
// ============================================================

// EntidadeContatoRequest representa a requisição para criar/atualizar um contato
type EntidadeContatoRequest struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	EntidadeID     int    `json:"entidade_id" binding:"required"`
	FormaContatoID int    `json:"forma_contato_id" binding:"required,min=1,max=7"`
	Informacao     string `json:"informacao" binding:"required"`
	Descricao      string `json:"descricao,omitempty"`

	// ============================================================
	// AUDITORIA
	// ============================================================
	CreatedBy *int `json:"created_by,omitempty"`
	UpdatedBy *int `json:"updated_by,omitempty"`
}

// ============================================================
// RESPONSES
// ============================================================

// EntidadeContatoResponse representa a resposta de um contato
type EntidadeContatoResponse struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	EntidadeID            int    `json:"entidade_id"`
	Item                  int    `json:"item"` // Sequencial por entidade
	FormaContatoID        int    `json:"forma_contato_id"`
	FormaContatoNome      string `json:"forma_contato_nome"` // Nome da forma de contato
	FormaContatoTipo      int    `json:"forma_contato_tipo"` // Tipo da forma de contato
	FormaContatoTipoLabel string `json:"forma_contato_tipo_label"`
	Informacao            string `json:"informacao"`
	Descricao             string `json:"descricao,omitempty"`

	// ============================================================
	// AUDITORIA
	// ============================================================
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	CreatedBy *int   `json:"created_by,omitempty"`
	UpdatedBy *int   `json:"updated_by,omitempty"`
}

// ============================================================
// LIST RESPONSE
// ============================================================

// EntidadeContatoListResponse representa a resposta de listagem de contatos
type EntidadeContatoListResponse struct {
	Items      []EntidadeContatoResponse `json:"items"`
	Total      int64                     `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE CONVERSÃO
// ============================================================

// ToModel converte EntidadeContatoRequest para models.EntidadeContato
func (r *EntidadeContatoRequest) ToModel() (*models.EntidadeContato, error) {
	if r == nil {
		return nil, nil
	}

	// Validar tipo de contato
	if !isValidTipoContato(r.FormaContatoID) {
		return nil, fmt.Errorf("tipo de contato inválido: %d", r.FormaContatoID)
	}

	contato := &models.EntidadeContato{
		EntidadeID:     r.EntidadeID,
		FormaContatoID: r.FormaContatoID,
		Informacao:     r.Informacao,
		Descricao:      utils.StringPtr(r.Descricao),
		CreatedBy:      r.CreatedBy,
		UpdatedBy:      r.UpdatedBy,
	}

	return contato, nil
}

// FromModel converte models.EntidadeContato para EntidadeContatoResponse
func (r *EntidadeContatoResponse) FromModel(contato *models.EntidadeContato) *EntidadeContatoResponse {
	if contato == nil {
		return nil
	}

	r.EntidadeID = contato.EntidadeID
	r.Item = contato.Item
	r.FormaContatoID = contato.FormaContatoID
	r.FormaContatoNome = getFormaContatoNome(contato.FormaContatoID)
	r.FormaContatoTipo = getFormaContatoTipo(contato.FormaContatoID)
	r.FormaContatoTipoLabel = getFormaContatoTipoLabel(contato.FormaContatoID)
	r.Informacao = contato.Informacao
	r.Descricao = utils.StringValue(contato.Descricao)

	// Auditoria
	r.CreatedAt = contato.CreatedAt.Format("2006-01-02 15:04:05")
	r.UpdatedAt = contato.UpdatedAt.Format("2006-01-02 15:04:05")
	r.CreatedBy = contato.CreatedBy
	r.UpdatedBy = contato.UpdatedBy

	// Preencher dados do relacionamento FormaContato (se carregado)
	if contato.FormaContato != nil {
		r.FormaContatoNome = contato.FormaContato.Descricao
		r.FormaContatoTipo = contato.FormaContato.Tipo
		r.FormaContatoTipoLabel = getTipoContatoLabel(contato.FormaContato.Tipo)
	}

	return r
}

// ============================================================
// FUNÇÕES AUXILIARES
// ============================================================

// isValidTipoContato valida se o tipo de contato é válido
func isValidTipoContato(tipo int) bool {
	switch tipo {
	case TipoContatoTelefone,
		TipoContatoCelular,
		TipoContatoEmail,
		TipoContatoWhatsApp,
		TipoContatoSite,
		TipoContatoFacebook,
		TipoContatoInstagram:
		return true
	default:
		return false
	}
}

// getFormaContatoNome retorna o nome da forma de contato
func getFormaContatoNome(id int) string {
	switch id {
	case TipoContatoTelefone:
		return "Telefone"
	case TipoContatoCelular:
		return "Celular"
	case TipoContatoEmail:
		return "E-mail"
	case TipoContatoWhatsApp:
		return "WhatsApp"
	case TipoContatoSite:
		return "Site"
	case TipoContatoFacebook:
		return "Facebook"
	case TipoContatoInstagram:
		return "Instagram"
	default:
		return "Desconhecido"
	}
}

// getFormaContatoTipo retorna o tipo da forma de contato
func getFormaContatoTipo(id int) int {
	switch id {
	case TipoContatoTelefone, TipoContatoCelular:
		return 1 // Telefone
	case TipoContatoEmail:
		return 2 // Email
	case TipoContatoWhatsApp:
		return 3 // WhatsApp
	case TipoContatoSite:
		return 4 // Site
	case TipoContatoFacebook, TipoContatoInstagram:
		return 5 // Rede Social
	default:
		return 0
	}
}

// getFormaContatoTipoLabel retorna o label do tipo da forma de contato
func getFormaContatoTipoLabel(id int) string {
	switch getFormaContatoTipo(id) {
	case 1:
		return "Telefone"
	case 2:
		return "E-mail"
	case 3:
		return "WhatsApp"
	case 4:
		return "Site"
	case 5:
		return "Rede Social"
	default:
		return "Desconhecido"
	}
}

// getTipoContatoLabel retorna o label do tipo de contato
func getTipoContatoLabel(tipo int) string {
	switch tipo {
	case 1:
		return "Telefone"
	case 2:
		return "E-mail"
	case 3:
		return "WhatsApp"
	case 4:
		return "Site"
	case 5:
		return "Rede Social"
	default:
		return "Desconhecido"
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o EntidadeContatoRequest
func (r *EntidadeContatoRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
