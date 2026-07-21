package models

import (
	"time"

	"github.com/openerp/backend/internal/constants"
	"gorm.io/gorm"
)

// ============================================================
// MODEL
// ============================================================

type DocumentoVenda struct {
	// ============================================================
	// CAMPOS PRINCIPAIS
	// ============================================================
	ID                  int  `gorm:"column:ddv_id;primaryKey;autoIncrement" json:"id"`
	EmpresaFilialID     int  `gorm:"column:emf_id;not null" json:"empresa_filial_id"`
	EntidadeID          *int `gorm:"column:ent_id" json:"entidade_id,omitempty"`             // CORRIGIDO: pode ser NULL
	CondicaoPagamentoID int  `gorm:"column:codpgt_id;not null" json:"condicao_pagamento_id"` // CORRIGIDO: codpgt_id
	ProcessoID          int  `gorm:"column:prc_id;not null" json:"processo_id"`
	FluxoID             int  `gorm:"column:flu_id;not null" json:"fluxo_id"`
	TabelaPrecoID       int  `gorm:"column:tbp_id;not null" json:"tabela_preco_id"`

	// ============================================================
	// CAMPOS DO DOCUMENTO
	// ============================================================
	Numero        int                          `gorm:"column:ddv_numero;not null" json:"numero"`
	TipoDocumento constants.TipoDocumentoVenda `gorm:"column:ddv_tipodocumento;not null;default:1" json:"tipo_documento"` //                     `gorm:"column:ddv_tipodocumeto;not null" json:"tipo_documento"` // 1-Orçamento, 2-Pedido
	TipoOperacao  constants.TipoOperacao       `gorm:"column:ddv_tipooperacao;not null" json:"tipo_operacao"`             // 0-Entrada, 1-Saída
	DataDocumento time.Time                    `gorm:"column:ddv_datadocumento;type:date;not null" json:"data_documento"`
	DataValidade  *time.Time                   `gorm:"column:ddv_datavalidade;type:date" json:"data_validade,omitempty"`
	DataPrevisao  *time.Time                   `gorm:"column:ddv_dataprevisao;type:date" json:"data_previsao,omitempty"`
	DataExpedicao *time.Time                   `gorm:"column:ddv_dataexpedicao;type:date" json:"data_expedicao,omitempty"`
	Situacao      constants.SituacaoPedido     `gorm:"column:ddv_situacao;not null;default:1" json:"situacao"`

	// ============================================================
	// CAMPOS DE VALORES (MONETÁRIOS)
	// ============================================================
	TotalProdutos  float64  `gorm:"column:ddv_totalprodutos;type:decimal(15,2);not null;default:0.00" json:"total_produtos"` // CORRIGIDO: decimal(15,2)
	TotalDescontos float64  `gorm:"column:ddv_totaldescontos;type:decimal(15,2);not null;default:0.00" json:"total_descontos"`
	TotalLiquido   float64  `gorm:"column:ddv_totalliquido;type:decimal(15,2);not null;default:0.00" json:"total_liquido"` // CORRIGIDO: decimal(15,2)
	ValorDesconto  *float64 `gorm:"column:ddv_valordesconto;type:decimal(15,2)" json:"valor_desconto,omitempty"`
	ValorFrete     *float64 `gorm:"column:ddv_valorfrete;type:decimal(15,2)" json:"valor_frete,omitempty"`

	// ============================================================
	// CAMPOS DE PESO
	// ============================================================
	TotalPesoBruto   float64 `gorm:"column:ddv_totalpesobruto;type:decimal(15,4);not null;default:0.0000" json:"total_peso_bruto"`
	TotalPesoLiquido float64 `gorm:"column:ddv_totalpesoliquido;type:decimal(15,4);not null;default:0.0000" json:"total_peso_liquido"`

	// ============================================================
	// CAMPOS DE DELIVERY
	// ============================================================
	TipoEntrega        string   `gorm:"column:ddv_tipo_entrega;type:enum('RETIRADA','ENTREGA','LOCAL');default:'LOCAL'" json:"tipo_entrega"`
	EnderecoEntrega    *string  `gorm:"column:ddv_endereco_entrega;type:text" json:"endereco_entrega,omitempty"`        // CORRIGIDO: nome da coluna
	TelefoneContato    *string  `gorm:"column:ddv_telefone_contato;type:varchar(20)" json:"telefone_contato,omitempty"` // CORRIGIDO: nome da coluna
	TaxaEntrega        *float64 `gorm:"column:ddv_taxaentrega;type:decimal(10,2)" json:"taxa_entrega,omitempty"`
	ObservacoesEntrega *string  `gorm:"column:ddv_observacoes_entrega;type:text" json:"observacoes_entrega,omitempty"`
	AceitaTroco        *bool    `gorm:"column:ddv_aceita_troco;default:false" json:"aceita_troco,omitempty"`
	TrocoPara          *float64 `gorm:"column:ddv_troco_para;type:decimal(10,2)" json:"troco_para,omitempty"`

	// ============================================================
	// CAMPOS DE ENDEREÇO (cópia para entrega)
	// ============================================================
	EndLogradouro  *string `gorm:"column:ddv_end_logradouro;type:varchar(100)" json:"end_logradouro,omitempty"`
	EndNumero      *string `gorm:"column:ddv_end_numero;type:varchar(20)" json:"end_numero,omitempty"`
	EndComplemento *string `gorm:"column:ddv_end_compl;type:varchar(255)" json:"end_complemento,omitempty"`
	EndBairro      *string `gorm:"column:ddv_end_bairro;type:varchar(255)" json:"end_bairro,omitempty"`
	EndCEP         *int    `gorm:"column:ddv_end_cep" json:"end_cep,omitempty"`

	// ============================================================
	// CAMPOS DE OBSERVAÇÕES
	// ============================================================
	ObservacoesInterna   *string `gorm:"column:ddv_observacoes_interna;type:text" json:"observacoes_interna,omitempty"`
	ObservacoesADM       *string `gorm:"column:ddv_observacoes_adm;type:text" json:"observacoes_adm,omitempty"`
	ObservacoesEspeciais *string `gorm:"column:ddv_observacoes_especiais;type:text" json:"observacoes_especiais,omitempty"`
	MotivoCancelamento   *string `gorm:"column:ddv_motivocancelamento;type:varchar(1000)" json:"motivo_cancelamento,omitempty"`

	// ============================================================
	// CAMPOS DE INFORMAÇÕES DO CLIENTE (CÓPIA)
	// ============================================================
	EntRazaoSocial       *string `gorm:"column:ddv_ent_razaosocial;type:varchar(255)" json:"ent_razao_social,omitempty"`
	EntInscricaoFederal  *string `gorm:"column:ddv_ent_inscricaofederal;type:varchar(20)" json:"ent_inscricao_federal,omitempty"`
	EntInscricaoEstadual *string `gorm:"column:ddv_ent_inscricaoestadual;type:varchar(20)" json:"ent_inscricao_estadual,omitempty"`

	// ============================================================
	// CAMPOS DE VEÍCULO (para entregas)
	// ============================================================
	VeiculoPlaca  *string `gorm:"column:ddv_veic_placa;type:varchar(20)" json:"veiculo_placa,omitempty"`
	VeiculoModelo *string `gorm:"column:ddv_veic_modelo;type:varchar(255)" json:"veiculo_modelo,omitempty"`
	VeiculoMarca  *string `gorm:"column:ddv_veic_marca;type:varchar(255)" json:"veiculo_marca,omitempty"`
	NomeMotorista *string `gorm:"column:ddv_nomemotorista;type:varchar(255)" json:"nome_motorista,omitempty"`

	// ============================================================
	// CAMPOS DE PAGAMENTO
	// ============================================================
	ValorPago  *float64 `gorm:"column:ddv_valorpago;type:decimal(15,2);default:0.00" json:"valor_pago,omitempty"`
	ValorTroco *float64 `gorm:"column:ddv_valortroco;type:decimal(15,2);default:0.00" json:"valor_troco,omitempty"`

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
	Entidade          *Entidade                 `gorm:"foreignKey:EntidadeID;references:ent_id" json:"entidade,omitempty"`
	CondicaoPagamento *CondicaoPagamento        `gorm:"foreignKey:CondicaoPagamentoID;references:codpgt_id" json:"condicao_pagamento,omitempty"`
	Itens             []DocumentoVendaItem      `gorm:"foreignKey:DocumentoVendaID;references:ID" json:"itens,omitempty"`
	Pagamentos        []DocumentoVendaPagamento `gorm:"foreignKey:DocumentoVendaID;references:ID" json:"pagamentos,omitempty"`
}

// ============================================================
// MÉTODOS DO MODEL
// ============================================================

func (DocumentoVenda) TableName() string {
	return "documento_venda"
}

func (d *DocumentoVenda) BeforeCreate(tx *gorm.DB) error {
	if d.CreatedBy == nil {
		d.CreatedBy = new(int)
		*d.CreatedBy = 0
	}
	if d.UpdatedBy == nil {
		d.UpdatedBy = new(int)
		*d.UpdatedBy = 0
	}
	// TODO: Gerar número do pedido automaticamente
	// TODO: Validar se a entidade existe e está ativa
	return nil
}

func (d *DocumentoVenda) BeforeUpdate(tx *gorm.DB) error {
	if d.UpdatedBy == nil {
		d.UpdatedBy = new(int)
		*d.UpdatedBy = 0
	}
	// TODO: Se situação mudar para "Fechado", verificar estoque
	// TODO: Se situação mudar para "Cancelado", liberar estoque
	return nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// IsActive verifica se o pedido está ativo (aberto ou em atividade)
func (d *DocumentoVenda) IsActive() bool {
	return d.Situacao == constants.SituacaoPedidoAberto || d.Situacao == constants.SituacaoPedidoEmAtividade
}

// IsDeleted verifica se o pedido foi deletado logicamente
func (d *DocumentoVenda) IsDeleted() bool {
	return d.DeletedAt != nil
}

// IsCancelado verifica se o pedido está cancelado
func (d *DocumentoVenda) IsCancelado() bool {
	return d.Situacao == constants.SituacaoPedidoCancelado
}

// IsFechado verifica se o pedido está fechado
func (d *DocumentoVenda) IsFechado() bool {
	return d.Situacao == constants.SituacaoPedidoFechado
}

// IsEmAtividade verifica se o pedido está em atividade
func (d *DocumentoVenda) IsEmAtividade() bool {
	return d.Situacao == constants.SituacaoPedidoEmAtividade
}

// SoftDelete realiza a exclusão lógica (CORRIGIDO: nome correto)
func (d *DocumentoVenda) SoftDelete() {
	now := time.Now()
	d.DeletedAt = &now
	d.Situacao = constants.SituacaoPedidoCancelado
}

// IsEntrega verifica se é entrega em domicílio
func (d *DocumentoVenda) IsEntrega() bool {
	return d.TipoEntrega == constants.TipoEntregaEntrega
}

// IsRetirada verifica se é retirada no local
func (d *DocumentoVenda) IsRetirada() bool {
	return d.TipoEntrega == constants.TipoEntregaRetirada
}

// IsLocal verifica se é consumo no local
func (d *DocumentoVenda) IsLocal() bool {
	return d.TipoEntrega == constants.TipoEntregaLocal
}

// GetTotalComDesconto retorna o total com desconto aplicado
func (d *DocumentoVenda) GetTotalComDesconto() float64 {
	return d.TotalLiquido - d.TotalDescontos
}
