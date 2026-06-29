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

// EntidadeRequest representa a requisição para criar/atualizar uma entidade
type EntidadeRequest struct {
	// ============================================================
	// DADOS PRINCIPAIS
	// ============================================================
	GrupoEntidadeID    *int   `json:"grupo_entidade_id,omitempty"`
	EmpresaFilialID    int    `json:"empresa_filial_id" binding:"required"`
	TabelaPrecoID      *int   `json:"tabela_preco_id,omitempty"`
	TabelaDescontoID   *int   `json:"tabela_desconto_id,omitempty"`
	HorarioID          *int   `json:"horario_id,omitempty"`
	Codigo             int    `json:"codigo,omitempty"`
	TipoPessoa         int    `json:"tipo_pessoa" binding:"required,min=1,max=2"`
	RazaoSocial        string `json:"razao_social" binding:"required"`
	NomeFantasia       string `json:"nome_fantasia,omitempty"`
	InscricaoFederal   string `json:"inscricao_federal" binding:"required"`
	InscricaoEstadual  string `json:"inscricao_estadual,omitempty"`
	InscricaoProdutor  string `json:"inscricao_produtor,omitempty"`
	InscricaoMunicipal string `json:"inscricao_municipal,omitempty"`
	Suframa            string `json:"suframa,omitempty"`
	Situacao           int    `json:"situacao,omitempty"`
	Classificacao      int    `json:"classificacao,omitempty"`
	Observacao         string `json:"observacao,omitempty"`

	// ============================================================
	// DADOS PESSOAIS
	// ============================================================
	DataNascimento string  `json:"data_nascimento,omitempty"`
	NomeDaMae      string  `json:"nome_da_mae,omitempty"`
	NomeDoPai      string  `json:"nome_do_pai,omitempty"`
	Sexo           int     `json:"sexo,omitempty"`
	CasaPropria    int     `json:"casa_propria,omitempty"`
	EstadoCivil    int     `json:"estado_civil,omitempty"`
	ConjujeNome    string  `json:"conjuje_nome,omitempty"`
	ConjujeCPF     string  `json:"conjuje_cpf,omitempty"`
	ConjujeRenda   float64 `json:"conjuje_renda,omitempty"`
	ConjujeRG      string  `json:"conjuje_rg,omitempty"`
	QuantFilhos    int     `json:"quant_filhos,omitempty"`

	// ============================================================
	// DADOS COMERCIAIS
	// ============================================================
	DataCadastro       string  `json:"data_cadastro,omitempty"`
	PercentualComissao float64 `json:"percentual_comissao,omitempty"`
	TaxaEntrega        float64 `json:"taxa_entrega,omitempty"`
	ArquivoImpDDV      string  `json:"arquivo_imp_ddv,omitempty"`

	// ============================================================
	// USUÁRIO (para auditoria)
	// ============================================================
	CreatedBy *int `json:"created_by,omitempty"`
	UpdatedBy *int `json:"updated_by,omitempty"`
}

// ============================================================
// RESPONSES
// ============================================================

// EntidadeResponse representa a resposta de uma entidade
type EntidadeResponse struct {
	// ============================================================
	// DADOS PRINCIPAIS
	// ============================================================
	ID                 int    `json:"id"`
	GrupoEntidadeID    *int   `json:"grupo_entidade_id,omitempty"`
	EmpresaFilialID    int    `json:"empresa_filial_id"`
	TabelaPrecoID      *int   `json:"tabela_preco_id,omitempty"`
	TabelaDescontoID   *int   `json:"tabela_desconto_id,omitempty"`
	HorarioID          *int   `json:"horario_id,omitempty"`
	Codigo             int    `json:"codigo"`
	TipoPessoa         int    `json:"tipo_pessoa"`
	TipoPessoaLabel    string `json:"tipo_pessoa_label"`
	RazaoSocial        string `json:"razao_social"`
	NomeFantasia       string `json:"nome_fantasia,omitempty"`
	InscricaoFederal   string `json:"inscricao_federal"`
	InscricaoEstadual  string `json:"inscricao_estadual,omitempty"`
	InscricaoProdutor  string `json:"inscricao_produtor,omitempty"`
	InscricaoMunicipal string `json:"inscricao_municipal,omitempty"`
	Suframa            string `json:"suframa,omitempty"`
	Situacao           int    `json:"situacao"`
	SituacaoLabel      string `json:"situacao_label"`
	Classificacao      int    `json:"classificacao,omitempty"`
	Observacao         string `json:"observacao,omitempty"`

	// ============================================================
	// DADOS PESSOAIS
	// ============================================================
	DataNascimento   string  `json:"data_nascimento,omitempty"`
	NomeDaMae        string  `json:"nome_da_mae,omitempty"`
	NomeDoPai        string  `json:"nome_do_pai,omitempty"`
	Sexo             int     `json:"sexo,omitempty"`
	SexoLabel        string  `json:"sexo_label,omitempty"`
	CasaPropria      int     `json:"casa_propria,omitempty"`
	EstadoCivil      int     `json:"estado_civil,omitempty"`
	EstadoCivilLabel string  `json:"estado_civil_label,omitempty"`
	ConjujeNome      string  `json:"conjuje_nome,omitempty"`
	ConjujeCPF       string  `json:"conjuje_cpf,omitempty"`
	ConjujeRenda     float64 `json:"conjuje_renda,omitempty"`
	ConjujeRG        string  `json:"conjuje_rg,omitempty"`
	QuantFilhos      int     `json:"quant_filhos,omitempty"`

	// ============================================================
	// DADOS COMERCIAIS
	// ============================================================
	DataCadastro       string  `json:"data_cadastro,omitempty"`
	PercentualComissao float64 `json:"percentual_comissao,omitempty"`
	TaxaEntrega        float64 `json:"taxa_entrega,omitempty"`
	ArquivoImpDDV      string  `json:"arquivo_imp_ddv,omitempty"`

	// ============================================================
	// AUDITORIA
	// ============================================================
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	CreatedBy *int   `json:"created_by,omitempty"`
	UpdatedBy *int   `json:"updated_by,omitempty"`
}

// ============================================================
// RESPONSE LISTA
// ============================================================

// EntidadeListResponse representa a resposta de listagem de entidades
type EntidadeListResponse struct {
	Items      []EntidadeResponse `json:"items"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}

// ============================================================
// MÉTODOS DE CONVERSÃO
// ============================================================

// ToModel converte EntidadeRequest para models.Entidade
func (r *EntidadeRequest) ToModel() *models.Entidade {
	if r == nil {
		return nil
	}

	entidade := &models.Entidade{
		GrupoEntidadeID:    r.GrupoEntidadeID,
		EmpresaFilialID:    r.EmpresaFilialID,
		TabelaPrecoID:      r.TabelaPrecoID,
		TabelaDescontoID:   r.TabelaDescontoID,
		HorarioID:          r.HorarioID,
		Codigo:             r.Codigo,
		TipoPessoa:         constants.TipoPessoa(r.TipoPessoa),
		RazaoSocial:        r.RazaoSocial,
		NomeFantasia:       utils.StringPtr(r.NomeFantasia),
		InscricaoFederal:   r.InscricaoFederal,
		InscricaoEstadual:  utils.StringPtr(r.InscricaoEstadual),
		InscricaoProdutor:  utils.StringPtr(r.InscricaoProdutor),
		InscricaoMunicipal: utils.StringPtr(r.InscricaoMunicipal),
		Suframa:            utils.StringPtr(r.Suframa),
		Situacao:           constants.Status(r.Situacao),
		Classificacao:      utils.IntPtr(r.Classificacao),
		Observacao:         utils.StringPtr(r.Observacao),
		CasaPropria:        utils.IntPtr(r.CasaPropria),
		EstadoCivil:        utils.IntPtr(r.EstadoCivil),
		ConjujeNome:        utils.StringPtr(r.ConjujeNome),
		ConjujeCPF:         utils.StringPtr(r.ConjujeCPF),
		ConjujeRenda:       utils.Float64Ptr(r.ConjujeRenda),
		ConjujeRG:          utils.StringPtr(r.ConjujeRG),
		QuantFilhos:        utils.IntPtr(r.QuantFilhos),
		PercentualComissao: utils.Float64Ptr(r.PercentualComissao),
		TaxaEntrega:        utils.Float64Ptr(r.TaxaEntrega),
		ArquivoImpDDV:      utils.StringPtr(r.ArquivoImpDDV),
		CreatedBy:          r.CreatedBy,
		UpdatedBy:          r.UpdatedBy,
	}

	// Converter datas
	if r.DataNascimento != "" {
		if data, err := utils.ParseDate(r.DataNascimento); err == nil {
			entidade.DataNascimento = &data
		}
	}

	if r.DataCadastro != "" {
		if data, err := utils.ParseDateTime(r.DataCadastro); err == nil {
			entidade.DataCadastro = &data
		}
	}

	// Converter sexo
	if r.Sexo > 0 {
		sexo := constants.Sexo(r.Sexo)
		entidade.Sexo = &sexo
	}

	// Se situação não foi informada, definir como ativo
	if r.Situacao == 0 {
		entidade.Situacao = constants.StatusAtivo
	}

	return entidade
}

// FromModel converte models.Entidade para EntidadeResponse
func (r *EntidadeResponse) FromModel(entidade *models.Entidade) *EntidadeResponse {
	if entidade == nil {
		return nil
	}

	r.ID = entidade.ID
	r.GrupoEntidadeID = entidade.GrupoEntidadeID
	r.EmpresaFilialID = entidade.EmpresaFilialID
	r.TabelaPrecoID = entidade.TabelaPrecoID
	r.TabelaDescontoID = entidade.TabelaDescontoID
	r.HorarioID = entidade.HorarioID
	r.Codigo = entidade.Codigo
	r.TipoPessoa = int(entidade.TipoPessoa)
	r.TipoPessoaLabel = entidade.TipoPessoa.String()
	r.RazaoSocial = entidade.RazaoSocial
	r.NomeFantasia = utils.StringValue(entidade.NomeFantasia)
	r.InscricaoFederal = entidade.InscricaoFederal
	r.InscricaoEstadual = utils.StringValue(entidade.InscricaoEstadual)
	r.InscricaoProdutor = utils.StringValue(entidade.InscricaoProdutor)
	r.InscricaoMunicipal = utils.StringValue(entidade.InscricaoMunicipal)
	r.Suframa = utils.StringValue(entidade.Suframa)
	r.Situacao = int(entidade.Situacao)
	r.SituacaoLabel = entidade.Situacao.String()
	r.Classificacao = utils.IntValue(entidade.Classificacao)
	r.Observacao = utils.StringValue(entidade.Observacao)

	// Dados pessoais
	if entidade.DataNascimento != nil {
		r.DataNascimento = utils.FormatDate(*entidade.DataNascimento)
	}
	r.NomeDaMae = utils.StringValue(entidade.NomeDaMae)
	r.NomeDoPai = utils.StringValue(entidade.NomeDoPai)
	if entidade.Sexo != nil {
		r.Sexo = int(*entidade.Sexo)
		r.SexoLabel = entidade.Sexo.String()
	}
	r.CasaPropria = utils.IntValue(entidade.CasaPropria)
	r.EstadoCivil = utils.IntValue(entidade.EstadoCivil)
	r.EstadoCivilLabel = getEstadoCivilLabel(r.EstadoCivil)
	r.ConjujeNome = utils.StringValue(entidade.ConjujeNome)
	r.ConjujeCPF = utils.StringValue(entidade.ConjujeCPF)
	r.ConjujeRenda = utils.Float64Value(entidade.ConjujeRenda)
	r.ConjujeRG = utils.StringValue(entidade.ConjujeRG)
	r.QuantFilhos = utils.IntValue(entidade.QuantFilhos)

	// Dados comerciais
	if entidade.DataCadastro != nil {
		r.DataCadastro = utils.FormatDateTime(*entidade.DataCadastro)
	}
	r.PercentualComissao = utils.Float64Value(entidade.PercentualComissao)
	r.TaxaEntrega = utils.Float64Value(entidade.TaxaEntrega)
	r.ArquivoImpDDV = utils.StringValue(entidade.ArquivoImpDDV)

	// Auditoria
	r.CreatedAt = utils.FormatDateTime(entidade.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(entidade.UpdatedAt)
	r.CreatedBy = entidade.CreatedBy
	r.UpdatedBy = entidade.UpdatedBy

	return r
}

// ============================================================
// FUNÇÕES AUXILIARES
// ============================================================

// getEstadoCivilLabel retorna o label do estado civil
func getEstadoCivilLabel(valor int) string {
	switch valor {
	case 1:
		return "Solteiro(a)"
	case 2:
		return "Casado(a)"
	case 3:
		return "Divorciado(a)"
	case 4:
		return "Viúvo(a)"
	case 5:
		return "União Estável"
	default:
		return "Não informado"
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO (opcional)
// ============================================================

// Validate valida o EntidadeRequest
func (r *EntidadeRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
