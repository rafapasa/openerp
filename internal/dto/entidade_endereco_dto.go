package dto

import (
	"github.com/go-playground/validator/v10"

	"github.com/openerp/backend/internal/constants"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

// ============================================================
// REQUESTS
// ============================================================

// EntidadeEnderecoRequest representa a requisição para criar/atualizar um endereço
type EntidadeEnderecoRequest struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	EntidadeID  int      `json:"entidade_id" binding:"required"`
	PaisID      int      `json:"pais_id" binding:"required"`
	EstadoID    int      `json:"estado_id" binding:"required"`
	MunicipioID int      `json:"municipio_id" binding:"required"`
	Tipo        int      `json:"tipo" binding:"required,min=1,max=4"`
	CEP         string   `json:"cep" binding:"required"`
	Logradouro  string   `json:"logradouro" binding:"required"`
	Numero      string   `json:"numero" binding:"required"`
	Complemento string   `json:"complemento,omitempty"`
	Bairro      string   `json:"bairro" binding:"required"`
	Distancia   *float64 `json:"distancia,omitempty"`
	Observacao  string   `json:"observacao,omitempty"`
	Situacao    int      `json:"situacao,omitempty"` // 1-Ativo, 2-Inativo

	// ============================================================
	// DATAS
	// ============================================================
	DataIni string `json:"data_ini" binding:"required"` // Formato: 2006-01-02
	DataFim string `json:"data_fim,omitempty"`

	// ============================================================
	// AUDITORIA
	// ============================================================
	CreatedBy *int `json:"created_by,omitempty"`
	UpdatedBy *int `json:"updated_by,omitempty"`
}

// ============================================================
// RESPONSES
// ============================================================

// EntidadeEnderecoResponse representa a resposta de um endereço
type EntidadeEnderecoResponse struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	EntidadeID    int      `json:"entidade_id"`
	Item          int      `json:"item"` // Sequencial por entidade
	PaisID        int      `json:"pais_id"`
	PaisNome      string   `json:"pais_nome"`
	EstadoID      int      `json:"estado_id"`
	EstadoUF      string   `json:"estado_uf"`
	EstadoNome    string   `json:"estado_nome"`
	MunicipioID   int      `json:"municipio_id"`
	MunicipioNome string   `json:"municipio_nome"`
	Tipo          int      `json:"tipo"`
	TipoLabel     string   `json:"tipo_label"` // Cobrança, Entrega, Comercial, Residencial
	CEP           string   `json:"cep"`
	Logradouro    string   `json:"logradouro,omitempty"`
	Numero        string   `json:"numero"`
	Complemento   string   `json:"complemento,omitempty"`
	Bairro        string   `json:"bairro"`
	Distancia     *float64 `json:"distancia,omitempty"`
	Observacao    string   `json:"observacao,omitempty"`
	Situacao      int      `json:"situacao"`
	SituacaoLabel string   `json:"situacao_label"`

	// ============================================================
	// DATAS
	// ============================================================
	DataIni string `json:"data_ini"`
	DataFim string `json:"data_fim,omitempty"`

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

// EntidadeEnderecoListResponse representa a resposta de listagem de endereços
type EntidadeEnderecoListResponse struct {
	Items      []EntidadeEnderecoResponse `json:"items"`
	Total      int64                      `json:"total"`
	Page       int                        `json:"page"`
	Limit      int                        `json:"limit"`
	TotalPages int                        `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE CONVERSÃO
// ============================================================

// ToModel converte EntidadeEnderecoRequest para models.EntidadeEndereco
func (r *EntidadeEnderecoRequest) ToModel() (*models.EntidadeEndereco, error) {
	if r == nil {
		return nil, nil
	}

	endereco := &models.EntidadeEndereco{
		EntidadeID:  r.EntidadeID,
		PaisID:      r.PaisID,
		EstadoID:    r.EstadoID,
		MunicipioID: r.MunicipioID,
		Tipo:        r.Tipo,
		Numero:      r.Numero,
		Bairro:      r.Bairro,
		Complemento: utils.StringPtr(r.Complemento),
		Logradouro:  utils.StringPtr(r.Logradouro),
		Distancia:   r.Distancia,
		Observacao:  utils.StringPtr(r.Observacao),
		CreatedBy:   r.CreatedBy,
		UpdatedBy:   r.UpdatedBy,
	}

	// Converter CEP (remover caracteres especiais)
	cepLimpo := utils.LimparDocumento(r.CEP)
	if cepLimpo != "" {
		// Converter para int (se possível)
		if cepInt, err := utils.ParseInt(cepLimpo); err == nil {
			endereco.CEP = cepInt
		}
	}

	// Converter datas
	if r.DataIni != "" {
		if data, err := utils.ParseDate(r.DataIni); err == nil {
			endereco.DataIni = data
		}
	}

	if r.DataFim != "" {
		if data, err := utils.ParseDate(r.DataFim); err == nil {
			endereco.DataFim = &data
		}
	}

	// Converter situação
	if r.Situacao == 0 {
		endereco.Situacao = int(constants.StatusAtivo)
	} else {
		endereco.Situacao = r.Situacao
	}

	return endereco, nil
}

// FromModel converte models.EntidadeEndereco para EntidadeEnderecoResponse
func (r *EntidadeEnderecoResponse) FromModel(endereco *models.EntidadeEndereco) *EntidadeEnderecoResponse {
	if endereco == nil {
		return nil
	}

	r.EntidadeID = endereco.EntidadeID
	r.Item = endereco.Item
	r.PaisID = endereco.PaisID
	r.EstadoID = endereco.EstadoID
	r.MunicipioID = endereco.MunicipioID
	r.Tipo = endereco.Tipo
	r.TipoLabel = constants.TipoEndereco(endereco.Tipo).String()
	r.CEP = utils.FormatarCEP(endereco.CEP)
	r.Logradouro = utils.StringValue(endereco.Logradouro)
	r.Numero = endereco.Numero
	r.Complemento = utils.StringValue(endereco.Complemento)
	r.Bairro = endereco.Bairro
	r.Distancia = endereco.Distancia
	r.Observacao = utils.StringValue(endereco.Observacao)
	r.Situacao = int(endereco.Situacao)
	r.SituacaoLabel = constants.Status.String(constants.Status(endereco.Situacao))

	// Datas
	r.DataIni = endereco.DataIni.Format("2006-01-02")
	if endereco.DataFim != nil {
		r.DataFim = endereco.DataFim.Format("2006-01-02")
	}

	// Auditoria
	r.CreatedAt = endereco.CreatedAt.Format("2006-01-02 15:04:05")
	r.UpdatedAt = endereco.UpdatedAt.Format("2006-01-02 15:04:05")
	r.CreatedBy = endereco.CreatedBy
	r.UpdatedBy = endereco.UpdatedBy

	// Preencher nomes dos relacionamentos (se carregados)
	if endereco.Pais != nil {
		r.PaisNome = endereco.Pais.Nome
	}
	if endereco.Estado != nil {
		r.EstadoUF = endereco.Estado.UF
		r.EstadoNome = endereco.Estado.Nome
	}
	if endereco.Municipio != nil {
		r.MunicipioNome = endereco.Municipio.Nome
	}

	return r
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o EntidadeEnderecoRequest
func (r *EntidadeEnderecoRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
