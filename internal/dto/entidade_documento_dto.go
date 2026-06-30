package dto

import (
	"encoding/base64"
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// CONSTANTES
// ============================================================

const (
	TipoDocumentoCNH      = "CNH"
	TipoDocumentoRG       = "RG"
	TipoDocumentoCPF      = "CPF"
	TipoDocumentoCertidao = "CERTIDAO"
	TipoDocumentoContrato = "CONTRATO"
	TipoDocumentoNF       = "NF"
	TipoDocumentoOutros   = "OUTROS"
)

// ============================================================
// REQUESTS
// ============================================================

// EntidadeDocumentoRequest representa a requisição para criar/atualizar um documento
type EntidadeDocumentoRequest struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	EntidadeID   int    `json:"entidade_id" binding:"required"`
	Descricao    string `json:"descricao,omitempty"`
	DataInclusao string `json:"data_inclusao" binding:"required"` // Formato: 2006-01-02
	Arquivo      string `json:"arquivo" binding:"required"`       // Base64
	Tipo         string `json:"tipo,omitempty"`

	// ============================================================
	// AUDITORIA
	// ============================================================
	CreatedBy *int `json:"created_by,omitempty"`
	UpdatedBy *int `json:"updated_by,omitempty"`
}

// ============================================================
// RESPONSES
// ============================================================

// EntidadeDocumentoResponse representa a resposta de um documento
type EntidadeDocumentoResponse struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	EntidadeID   int    `json:"entidade_id"`
	Item         int    `json:"item"` // Sequencial por entidade
	Descricao    string `json:"descricao,omitempty"`
	DataInclusao string `json:"data_inclusao"`
	Arquivo      string `json:"arquivo,omitempty"` // Base64 (opcional)
	Tipo         string `json:"tipo,omitempty"`
	TipoLabel    string `json:"tipo_label,omitempty"`

	// ============================================================
	// METADADOS
	// ============================================================
	Tamanho     int64  `json:"tamanho,omitempty"`      // Tamanho do arquivo em bytes
	ContentType string `json:"content_type,omitempty"` // Tipo do arquivo (ex: application/pdf)

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

// EntidadeDocumentoListResponse representa a resposta de listagem de documentos
type EntidadeDocumentoListResponse struct {
	Items      []EntidadeDocumentoResponse `json:"items"`
	Total      int64                       `json:"total"`
	Page       int                         `json:"page"`
	Limit      int                         `json:"limit"`
	TotalPages int                         `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE CONVERSÃO
// ============================================================

// ToModel converte EntidadeDocumentoRequest para models.EntidadeDocumento
func (r *EntidadeDocumentoRequest) ToModel() (*models.EntidadeDocumento, error) {
	if r == nil {
		return nil, nil
	}

	// 1. Validar se o arquivo é Base64 válido
	if r.Arquivo != "" {
		if _, err := base64.StdEncoding.DecodeString(r.Arquivo); err != nil {
			return nil, fmt.Errorf("arquivo inválido: %w", err)
		}
	}

	documento := &models.EntidadeDocumento{
		EntidadeID: r.EntidadeID,
		Descricao:  utils.StringPtr(r.Descricao),
		Tipo:       utils.StringPtr(r.Tipo),
		CreatedBy:  r.CreatedBy,
		UpdatedBy:  r.UpdatedBy,
	}

	// 2. Converter DataInclusao
	if r.DataInclusao != "" {
		if data, err := utils.ParseDate(r.DataInclusao); err == nil {
			documento.DataInclusao = data
		} else {
			return nil, fmt.Errorf("data de inclusão inválida: %w", err)
		}
	}

	// 3. Converter Arquivo (Base64 para []byte)
	if r.Arquivo != "" {
		data, err := base64.StdEncoding.DecodeString(r.Arquivo)
		if err != nil {
			return nil, fmt.Errorf("erro ao decodificar arquivo: %w", err)
		}
		documento.Arquivo = data
	}

	return documento, nil
}

// FromModel converte models.EntidadeDocumento para EntidadeDocumentoResponse
func (r *EntidadeDocumentoResponse) FromModel(documento *models.EntidadeDocumento) *EntidadeDocumentoResponse {
	if documento == nil {
		return nil
	}

	r.EntidadeID = documento.EntidadeID
	r.Item = documento.Item
	r.Descricao = utils.StringValue(documento.Descricao)
	r.DataInclusao = documento.DataInclusao.Format("2006-01-02")
	r.Tipo = utils.StringValue(documento.Tipo)
	r.TipoLabel = getTipoDocumentoLabel(utils.StringValue(documento.Tipo))

	// Tamanho do arquivo
	if len(documento.Arquivo) > 0 {
		r.Tamanho = int64(len(documento.Arquivo))
		r.ContentType = getContentType(documento.Arquivo)
		// Não retornamos o arquivo completo por padrão (performance)
		// Se necessário, pode ser retornado em um endpoint específico
		// r.Arquivo = base64.StdEncoding.EncodeToString(documento.Arquivo)
	}

	// Auditoria
	r.CreatedAt = documento.CreatedAt.Format("2006-01-02 15:04:05")
	r.UpdatedAt = documento.UpdatedAt.Format("2006-01-02 15:04:05")
	r.CreatedBy = documento.CreatedBy
	r.UpdatedBy = documento.UpdatedBy

	return r
}

// FromModelWithFile converte models.EntidadeDocumento para EntidadeDocumentoResponse com arquivo
func (r *EntidadeDocumentoResponse) FromModelWithFile(documento *models.EntidadeDocumento) *EntidadeDocumentoResponse {
	if documento == nil {
		return nil
	}

	r = r.FromModel(documento)

	// Incluir arquivo Base64
	if len(documento.Arquivo) > 0 {
		r.Arquivo = base64.StdEncoding.EncodeToString(documento.Arquivo)
	}

	return r
}

// ============================================================
// FUNÇÕES AUXILIARES
// ============================================================

// getTipoDocumentoLabel retorna o label do tipo de documento
func getTipoDocumentoLabel(tipo string) string {
	switch tipo {
	case TipoDocumentoCNH:
		return "CNH"
	case TipoDocumentoRG:
		return "RG"
	case TipoDocumentoCPF:
		return "CPF"
	case TipoDocumentoCertidao:
		return "Certidão"
	case TipoDocumentoContrato:
		return "Contrato"
	case TipoDocumentoNF:
		return "Nota Fiscal"
	case TipoDocumentoOutros:
		return "Outros"
	default:
		return "Não informado"
	}
}

// getContentType detecta o Content-Type baseado no cabeçalho do arquivo
func getContentType(data []byte) string {
	if len(data) == 0 {
		return "application/octet-stream"
	}

	// PDF
	if len(data) >= 4 && data[0] == 0x25 && data[1] == 0x50 && data[2] == 0x44 && data[3] == 0x46 {
		return "application/pdf"
	}

	// PNG
	if len(data) >= 8 && data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 &&
		data[4] == 0x0D && data[5] == 0x0A && data[6] == 0x1A && data[7] == 0x0A {
		return "image/png"
	}

	// JPEG
	if len(data) >= 2 && data[0] == 0xFF && data[1] == 0xD8 {
		return "image/jpeg"
	}

	// GIF
	if len(data) >= 6 && data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x38 &&
		(data[4] == 0x37 || data[4] == 0x39) && data[5] == 0x61 {
		return "image/gif"
	}

	// DOCX (ZIP)
	if len(data) >= 4 && data[0] == 0x50 && data[1] == 0x4B && data[2] == 0x03 && data[3] == 0x04 {
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	}

	return "application/octet-stream"
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o EntidadeDocumentoRequest
func (r *EntidadeDocumentoRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
