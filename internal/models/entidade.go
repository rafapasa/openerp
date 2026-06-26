package models

import "time"

type Entidade struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	ID                 int        `gorm:"column:ent_id;primaryKey;autoIncrement" json:"id"`
	GrupoEntidadeID    int        `gorm:"column:gpe_id;not null" json:"grupo_entidade_id"`
	EmpresaFilialID    int        `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	TabelaPrecoID      int        `gorm:"column:tbp_id;not null" json:"tabela_preco_id"`
	TabelaDescontoID   int        `gorm:"column:tdesc_id;not null" json:"tabela_desconto_id"`
	HorarioID          int        `gorm:"column:hor_id;not null" json:"horario_id"`
	Codigo             int        `gorm:"column:ent_codigo;type:varchar(20);not null;unique" json:"codigo"`
	TipoPessoa         int        `gorm:"column:ent_tipopessoa;not null" json:"tipo_pessoa"`
	RazaoSocial        string     `gorm:"column:ent_razaosocial;type:varchar(100);not null" json:"razao_social"`
	NomeFantasia       string     `gorm:"column:ent_nomefantasia;type:varchar(100);not null" json:"nome_fantasia"`
	InscricaoEstadual  string     `gorm:"column:ent_inscricaoestadual;type:varchar(20);not null" json:"inscricao_estadual"`
	InscricaoMunicipal string     `gorm:"column:ent_inscricaomunicipal;type:varchar(20);not null" json:"inscricao_municipal"`
	InscricaoFederal   string     `gorm:"column:ent_inscricaofederal;type:varchar(20);not null" json:"inscricao_federal"`
	Situacao           int        `gorm:"column:ent_situacao;not null" json:"situacao"`
	createdBy          *int       `gorm:"column:created_by" json:"created_by,omitempty"`
	updatedBy          *int       `gorm:"column:updated_by" json:"updated_by,omitempty"`
	createdAt          time.Time  `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	updatedAt          time.Time  `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt          *time.Time `gorm:"column:deleted_at;type:datetime;index" json:"deleted_at,omitempty"`
}

func (Entidade) TableName() string {
	return "entidade"
}

func (e *Entidade) BeforeCreate() error {
	if e.createdBy == nil {
		e.createdBy = new(int)
		*e.createdBy = 0
	}
	if e.updatedBy == nil {
		e.updatedBy = new(int)
		*e.updatedBy = 0
	}

	return nil
}

func (e *Entidade) BeforeUpdate() error {
	if e.updatedBy == nil {
		e.updatedBy = new(int)
		*e.updatedBy = 0
	}
	return nil
}

func (e *Entidade) IsActive() bool {
	return e.Situacao == 1
}

func (e *Entidade) IsDeleted() bool {
	return e.DeletedAt != nil
}

func (e *Entidade) SoftDelete() {
	now := time.Now()
	e.DeletedAt = &now
}
