package dto

import (
	"time"

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
	GrupoEntidadeID    *int                 `json:"grupo_entidade_id,omitempty"`
	EmpresaFilialID    int                  `json:"empresa_filial_id" binding:"required"`
	TabelaPrecoID      *int                 `json:"tabela_preco_id,omitempty"`
	TabelaDescontoID   *int                 `json:"tabela_desconto_id,omitempty"`
	HorarioID          *int                 `json:"horario_id,omitempty"`
	Codigo             *int                 `json:"codigo,omitempty"`
	TipoPessoa         constants.TipoPessoa `json:"tipo_pessoa" binding:"required,min=1,max=2"`
	RazaoSocial        string               `json:"razao_social" binding:"required"`
	NomeFantasia       string               `json:"nome_fantasia,omitempty"`
	InscricaoFederal   string               `json:"inscricao_federal" binding:"required"`
	InscricaoEstadual  *string              `json:"inscricao_estadual,omitempty"`
	InscricaoProdutor  *string              `json:"inscricao_produtor,omitempty"`
	InscricaoMunicipal *string              `json:"inscricao_municipal,omitempty"`
	Suframa            *string              `json:"suframa,omitempty"`
	Situacao           constants.Status     `json:"situacao,omitempty"`
	Classificacao      *int                 `json:"classificacao,omitempty"`
	Observacao         *string              `json:"observacao,omitempty"`

	// ============================================================
	// DADOS PESSOAIS
	// ============================================================
	DataNascimento string          `json:"data_nascimento,omitempty"`
	NomeDaMae      *string         `json:"nome_da_mae,omitempty"`
	NomeDoPai      *string         `json:"nome_do_pai,omitempty"`
	Sexo           *constants.Sexo `json:"sexo,omitempty"`
	CasaPropria    *int            `json:"casa_propria,omitempty"`
	EstadoCivil    *int            `json:"estado_civil,omitempty"`
	ConjujeNome    *string         `json:"conjuje_nome,omitempty"`
	ConjujeCPF     *string         `json:"conjuje_cpf,omitempty"`
	ConjujeRenda   *float64        `json:"conjuje_renda,omitempty"`
	ConjujeRG      *string         `json:"conjuje_rg,omitempty"`
	QuantFilhos    *int            `json:"quant_filhos,omitempty"`

	// ============================================================
	// DADOS COMERCIAIS
	// ============================================================
	DataCadastro       time.Time `json:"data_cadastro,omitempty"`
	PercentualComissao *float64  `json:"percentual_comissao,omitempty"`
	TaxaEntrega        *float64  `json:"taxa_entrega,omitempty"`
	ArquivoImpDDV      *string   `json:"arquivo_imp_ddv,omitempty"`

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
	ID                 int                  `json:"id"`
	GrupoEntidadeID    *int                 `json:"grupo_entidade_id,omitempty"`
	EmpresaFilialID    int                  `json:"empresa_filial_id"`
	TabelaPrecoID      *int                 `json:"tabela_preco_id,omitempty"`
	TabelaDescontoID   *int                 `json:"tabela_desconto_id,omitempty"`
	HorarioID          *int                 `json:"horario_id,omitempty"`
	Codigo             *int                 `json:"codigo"`
	TipoPessoa         constants.TipoPessoa `json:"tipo_pessoa"`
	TipoPessoaLabel    string               `json:"tipo_pessoa_label"`
	RazaoSocial        string               `json:"razao_social"`
	NomeFantasia       string               `json:"nome_fantasia,omitempty"`
	InscricaoFederal   string               `json:"inscricao_federal"`
	InscricaoEstadual  *string              `json:"inscricao_estadual,omitempty"`
	InscricaoProdutor  *string              `json:"inscricao_produtor,omitempty"`
	InscricaoMunicipal *string              `json:"inscricao_municipal,omitempty"`
	Suframa            *string              `json:"suframa,omitempty"`
	Situacao           constants.Status     `json:"situacao"`
	SituacaoLabel      string               `json:"situacao_label"`
	Classificacao      *int                 `json:"classificacao,omitempty"`
	Observacao         *string              `json:"observacao,omitempty"`

	// ============================================================
	// DADOS PESSOAIS
	// ============================================================
	DataNascimento   string          `json:"data_nascimento,omitempty"`
	NomeDaMae        *string         `json:"nome_da_mae,omitempty"`
	NomeDoPai        *string         `json:"nome_do_pai,omitempty"`
	Sexo             *constants.Sexo `json:"sexo,omitempty"`
	SexoLabel        *string         `json:"sexo_label,omitempty"`
	CasaPropria      *int            `json:"casa_propria,omitempty"`
	EstadoCivil      *int            `json:"estado_civil,omitempty"`
	EstadoCivilLabel *string         `json:"estado_civil_label,omitempty"`
	ConjujeNome      *string         `json:"conjuje_nome,omitempty"`
	ConjujeCPF       *string         `json:"conjuje_cpf,omitempty"`
	ConjujeRenda     *float64        `json:"conjuje_renda,omitempty"`
	ConjujeRG        *string         `json:"conjuje_rg,omitempty"`
	QuantFilhos      *int            `json:"quant_filhos,omitempty"`

	// ============================================================
	// DADOS COMERCIAIS
	// ============================================================
	DataCadastro       *string  `json:"data_cadastro,omitempty"`
	PercentualComissao *float64 `json:"percentual_comissao,omitempty"`
	TaxaEntrega        *float64 `json:"taxa_entrega,omitempty"`
	ArquivoImpDDV      *string  `json:"arquivo_imp_ddv,omitempty"`

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
// MÉTODOS DE CONVERSÃO (USANDO MAPPER)
// ============================================================

// ToModel converte EntidadeRequest para models.Entidade usando mapper
func (r *EntidadeRequest) ToModel() (*models.Entidade, error) {
	if r == nil {
		return nil, nil
	}

	entidade := &models.Entidade{}

	// 1. Usar o mapper para copiar campos (mapeamento automático)
	// 1. Usar o mapper para copiar campos (mapeamento automático)
	if err := utils.MapToModel(r, entidade); err != nil {
		return nil, err
	}

	// 2. Tratamentos especiais que o mapper não cobre
	// Converter TipoPessoa (int → constants.TipoPessoa)
	entidade.TipoPessoa = r.TipoPessoa

	// Limpar documento (remover pontos, traços, etc.)
	entidade.InscricaoFederal = utils.LimparDocumento(r.InscricaoFederal)

	// Converter Sexo (int → constants.Sexo)
	if r.Sexo != nil && *r.Sexo > 0 {
		entidade.Sexo = *r.Sexo
	}

	// Converter Situacao (se não informada, definir como ativo)
	if r.Situacao == 0 {
		entidade.Situacao = constants.StatusAtivo
	}

	entidade.DataNascimento = nil
	// Converter DataNascimento (se for uma data válida e não for o valor zero de time.Time)
	if r.DataNascimento != "" {
		if data, err := time.Parse("2006-01-02", r.DataNascimento); err == nil {
			entidade.DataNascimento = &data
		}
	}

	// Converter DataCadastro (se for uma data válida e não for o valor zero de time.Time)
	if !r.DataCadastro.IsZero() {
		entidade.DataCadastro = &r.DataCadastro

	}

	return entidade, nil
}

// FromModel converte models.Entidade para EntidadeResponse usando mapper
func (r *EntidadeResponse) FromModel(entidade *models.Entidade) *EntidadeResponse {
	if entidade == nil {
		return nil
	}

	// 1. Usar o mapper para copiar campos (mapeamento automático)
	// O mapper vai copiar todos os campos com nomes correspondentes
	if err := utils.MapToDTO(entidade, r); err != nil {
		// Se o mapper falhar, usar fallback manual
		return r.fromModelFallback(entidade)
	}

	if entidade.Codigo != nil {
		r.Codigo = entidade.Codigo
	}

	// 2. Preencher campos calculados (labels)
	r.TipoPessoaLabel = constants.TipoPessoa(r.TipoPessoa).String()
	r.SituacaoLabel = constants.Status(r.Situacao).String()

	if r.Sexo != nil && *r.Sexo > 0 {
		r.SexoLabel = utils.StringPtr((*r.Sexo).String())
	}

	if r.EstadoCivil != nil {
		r.EstadoCivilLabel = utils.StringPtr(getEstadoCivilLabel(*r.EstadoCivil))
	}

	// 3. Formatar datas (o mapper não faz isso)
	r.CreatedAt = utils.FormatDateTime(entidade.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(entidade.UpdatedAt)

	if !entidade.DataNascimento.IsZero() {
		r.DataNascimento = utils.FormatDate(*entidade.DataNascimento)
	}

	if entidade.DataCadastro != nil {
		r.DataCadastro = utils.StringPtr(utils.FormatDateTime(*entidade.DataCadastro))
	}

	return r
}

// ============================================================
// FALLBACK (caso o mapper falhe)
// ============================================================

// fromModelFallback é o fallback manual caso o mapper falhe
// Isso garante que o sistema continue funcionando mesmo com erro no mapper
func (r *EntidadeResponse) fromModelFallback(entidade *models.Entidade) *EntidadeResponse {
	// Mapeamento manual campo por campo (seguro)
	r.ID = entidade.ID
	r.GrupoEntidadeID = entidade.GrupoEntidadeID
	r.EmpresaFilialID = entidade.EmpresaFilialID
	r.TabelaPrecoID = entidade.TabelaPrecoID
	r.TabelaDescontoID = entidade.TabelaDescontoID
	r.HorarioID = entidade.HorarioID
	if entidade.Codigo != nil {
		r.Codigo = entidade.Codigo
	}
	r.TipoPessoa = entidade.TipoPessoa
	r.RazaoSocial = entidade.RazaoSocial
	r.NomeFantasia = utils.StringValue(entidade.NomeFantasia)
	r.InscricaoFederal = entidade.InscricaoFederal
	r.InscricaoEstadual = entidade.InscricaoEstadual
	r.InscricaoProdutor = entidade.InscricaoProdutor
	r.InscricaoMunicipal = entidade.InscricaoMunicipal
	r.Suframa = entidade.Suframa
	r.Situacao = entidade.Situacao
	r.Situacao = entidade.Situacao
	r.Classificacao = entidade.Classificacao
	r.Observacao = entidade.Observacao
	r.NomeDaMae = entidade.NomeDaMae
	r.NomeDoPai = entidade.NomeDoPai
	if entidade.Sexo > 0 {
		r.Sexo = &entidade.Sexo
	}
	r.CasaPropria = entidade.CasaPropria
	r.EstadoCivil = entidade.EstadoCivil
	r.ConjujeNome = entidade.ConjujeNome
	r.ConjujeCPF = entidade.ConjujeCPF
	r.ConjujeRenda = entidade.ConjujeRenda
	r.ConjujeRG = entidade.ConjujeRG
	r.QuantFilhos = entidade.QuantFilhos
	r.PercentualComissao = entidade.PercentualComissao
	r.TaxaEntrega = entidade.TaxaEntrega
	r.ArquivoImpDDV = entidade.ArquivoImpDDV
	r.CreatedBy = entidade.CreatedBy
	r.UpdatedBy = entidade.UpdatedBy

	// Labels
	r.TipoPessoaLabel = constants.TipoPessoa(entidade.TipoPessoa).String()
	r.SituacaoLabel = constants.Status(entidade.Situacao).String()
	if entidade.Sexo > 0 { // Check if Sexo is set
		r.SexoLabel = utils.StringPtr(constants.Sexo(entidade.Sexo).String())
	}
	if r.EstadoCivil != nil { // Check if EstadoCivil is set
		r.EstadoCivilLabel = utils.StringPtr(getEstadoCivilLabel(*r.EstadoCivil))
	}

	// Datas
	r.CreatedAt = utils.FormatDateTime(entidade.CreatedAt)
	r.UpdatedAt = utils.FormatDateTime(entidade.UpdatedAt)
	if !entidade.DataNascimento.IsZero() {
		r.DataNascimento = utils.FormatDate(*entidade.DataNascimento)
	}
	if entidade.DataCadastro != nil {
		r.DataCadastro = utils.StringPtr(utils.FormatDateTime(*entidade.DataCadastro))
	}

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
// MÉTODOS DE VALIDAÇÃO
// ============================================================

// Validate valida o EntidadeRequest
func (r *EntidadeRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
