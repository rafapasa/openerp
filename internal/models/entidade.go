package models

import (
	"time"
)

type Entidade struct {
    // ============================================================
    // CAMPOS PRINCIPAIS
    // ============================================================
    ID                 int    `gorm:"column:ent_id;primaryKey;autoIncrement" json:"id"`
    GrupoEntidadeID    int    `gorm:"column:gpe_id" json:"grupo_entidade_id"`                          // Pode ser NULL
    EmpresaFilialID    int    `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
    TabelaPrecoID      *int   `gorm:"column:tbp_id" json:"tabela_preco_id,omitempty"`                  // Pode ser NULL
    TabelaDescontoID   *int   `gorm:"column:tdesc_id" json:"tabela_desconto_id,omitempty"`            // Pode ser NULL
    HorarioID          *int   `gorm:"column:hor_id" json:"horario_id,omitempty"`                      // Pode ser NULL
    
    Codigo             int    `gorm:"column:ent_codigo;not null;unique" json:"codigo"`                // INT, não VARCHAR
    TipoPessoa         int    `gorm:"column:ent_tipopessoa;not null" json:"tipo_pessoa"`
    RazaoSocial        string `gorm:"column:ent_razaosocial;type:varchar(100);not null" json:"razao_social"`
    NomeFantasia       string `gorm:"column:ent_nomefantasia;type:varchar(100)" json:"nome_fantasia,omitempty"` // Pode ser NULL
    InscricaoFederal   string `gorm:"column:ent_inscricaofederal;type:varchar(20);not null" json:"inscricao_federal"`
    InscricaoEstadual  *string `gorm:"column:ent_inscricaoestadual;type:varchar(20)" json:"inscricao_estadual,omitempty"` // Pode ser NULL
    InscricaoMunicipal *string `gorm:"column:ent_inscricaomunicipal;type:varchar(20)" json:"inscricao_municipal,omitempty"` // Pode ser NULL
    Situacao           int    `gorm:"column:ent_situacao;not null" json:"situacao"`
    
    // Campos extras que podem existir
    InscricaoProdutor  *string `gorm:"column:ent_inscricaoprodutor;type:varchar(20)" json:"inscricao_produtor,omitempty"`
    Suframa            *string `gorm:"column:ent_suframa;type:varchar(9)" json:"suframa,omitempty"`
    Observacao         *string `gorm:"column:ent_observacao;type:text" json:"observacao,omitempty"`
    Classificacao      *int    `gorm:"column:ent_classificacao" json:"classificacao,omitempty"`
    DataNascimento     *time.Time `gorm:"column:ent_datanascimento;type:date" json:"data_nascimento,omitempty"`
    NomeDaMae          *string `gorm:"column:ent_nomedamae;type:varchar(255)" json:"nome_da_mae,omitempty"`
    NomeDoPai          *string `gorm:"column:ent_nomedopai;type:varchar(255)" json:"nome_do_pai,omitempty"`
    Sexo               *int    `gorm:"column:ent_sexo" json:"sexo,omitempty"` // 1-Masculino, 2-Feminino
    CasaPropria        *int    `gorm:"column:ent_casapropria" json:"casa_propria,omitempty"`
    EstadoCivil        *int    `gorm:"column:ent_estadocivil" json:"estado_civil,omitempty"` // 1-Solteiro, 2-Casado, 9-Outros
    ConjujeNome        *string `gorm:"column:ent_conjuje_nome;type:varchar(255)" json:"conjuje_nome,omitempty"`
    ConjujeCPF         *string `gorm:"column:ent_conjuje_cpf;type:varchar(20)" json:"conjuje_cpf,omitempty"`
    ConjujeRenda       *float64 `gorm:"column:ent_conjuje_renda;type:decimal(15,2)" json:"conjuje_renda,omitempty"`
    ConjujeRG          *string `gorm:"column:ent_conjuje_rg;type:varchar(20)" json:"conjuje_rg,omitempty"`
    QuantFilhos        *int    `gorm:"column:ent_quantfilhos" json:"quant_filhos,omitempty"`
    DataCadastro       *time.Time `gorm:"column:ent_datacadastro;type:datetime" json:"data_cadastro,omitempty"`
    PercentualComissao *float64 `gorm:"column:ent_percentualcomissao;type:decimal(5,2)" json:"percentual_comissao,omitempty"`
    TaxaEntrega        *float64 `gorm:"column:ent_taxaentrega;type:decimal(10,2)" json:"taxa_entrega,omitempty"`
    ArquivoImpDDV      *string `gorm:"column:ent_arq_imp_ddv;type:varchar(255)" json:"arquivo_imp_ddv,omitempty"`

    // ============================================================
    // CAMPOS DE AUDITORIA (COM LETRA MAIÚSCULA)
    // ============================================================
    CreatedAt time.Time  `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
    UpdatedAt time.Time  `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
    DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime;index" json:"deleted_at,omitempty"`
    CreatedBy *int       `gorm:"column:created_by" json:"created_by,omitempty"`
    UpdatedBy *int       `gorm:"column:updated_by" json:"updated_by,omitempty"`

    // ============================================================
    // RELACIONAMENTOS
    // ============================================================
    GrupoEntidade  *GrupoEntidade  `gorm:"foreignKey:GrupoEntidadeID;references:gpe_id" json:"grupo_entidade,omitempty"`
    EmpresaFilial  *EmpresaFilial  `gorm:"foreignKey:EmpresaFilialID;references:emf_id" json:"empresa_filial,omitempty"`
    TabelaPreco    *TabelaPreco    `gorm:"foreignKey:TabelaPrecoID;references:tbp_id" json:"tabela_preco,omitempty"`
    TabelaDesconto *TabelaDesconto `gorm:"foreignKey:TabelaDescontoID;references:tdesc_id" json:"tabela_desconto,omitempty"`
    Horario        *Horario        `gorm:"foreignKey:HorarioID;references:hor_id" json:"horario,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (Entidade) TableName() string {
    return "entidade"
}

func (e *Entidade) BeforeCreate() error {
    if e.CreatedBy == nil {
        e.CreatedBy = new(int)
        *e.CreatedBy = 0
    }
    if e.UpdatedBy == nil {
        e.UpdatedBy = new(int)
        *e.UpdatedBy = 0
    }
    return nil
}

func (e *Entidade) BeforeUpdate() error {
    if e.UpdatedBy == nil {
        e.UpdatedBy = new(int)
        *e.UpdatedBy = 0
    }
    return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

func (e *Entidade) IsActive() bool {
    return e.Situacao == 1
}

func (e *Entidade) IsDeleted() bool {
    return e.DeletedAt != nil
}

func (e *Entidade) SoftDelete() {
    now := time.Now()
    e.DeletedAt = &now
    e.Situacao = 0 // Marca como inativo
}

// Método para verificar se é pessoa física
func (e *Entidade) IsPessoaFisica() bool {
    return e.TipoPessoa == 1
}

// Método para verificar se é pessoa jurídica
func (e *Entidade) IsPessoaJuridica() bool {
    return e.TipoPessoa == 2
}

// Método para obter o documento (CPF ou CNPJ)
func (e *Entidade) GetDocumento() string {
    return e.InscricaoFederal
}