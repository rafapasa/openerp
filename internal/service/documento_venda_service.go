package service

import (
	"fmt"
	"time"

	"github.com/openerp/backend/internal/constants"
	"github.com/openerp/backend/internal/dto"
	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/utils"
	"gorm.io/gorm"
)

type DocumentoVendaService struct {
	db         *gorm.DB
	ddvRepo    *repository.DocumentoVendaRepository
	dviService *DocumentoVendaItemService
	dvpService *DocumentoVendaPagamentoService
	entService *EntidadeService
	proService *ProdutoService
}

func NewDocumentoVendaService(db *gorm.DB) *DocumentoVendaService {
	return &DocumentoVendaService{
		db:         db,
		ddvRepo:    repository.NewDocumentoVendaRepository(db),
		dviService: NewDocumentoVendaItemService(db),
		dvpService: NewDocumentoVendaPagamentoService(db),
		entService: NewEntidadeService(db),
		proService: NewProdutoService(db),
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO (PRIVADOS)
// ============================================================

func (s *DocumentoVendaService) isDataValid(req *dto.DocumentoVendaRequest, isUpdate bool) error {
	// 1. Validar campos obrigatórios do cabeçalho
	if req.EmpresaFilialID == 0 {
		return apperrors.NewValidationError("O campo 'empresa_filial_id' é obrigatório.")
	}
	if req.CondicaoPagamentoID == 0 {
		return apperrors.NewValidationError("O campo 'condicao_pagamento_id' é obrigatório.")
	}
	if req.TabelaPrecoID == 0 {
		return apperrors.NewValidationError("O campo 'tabela_preco_id' é obrigatório.")
	}

	// 2. Validar Entidade (se informada)
	if req.EntidadeID != nil && *req.EntidadeID > 0 {
		entidade, err := s.entService.GetByID(*req.EntidadeID)
		if err != nil {
			return apperrors.NewNotFoundError(fmt.Sprintf("Entidade com ID %d não encontrada.", *req.EntidadeID))
		}
		if !entidade.IsActive() {
			return apperrors.NewValidationError(fmt.Sprintf("A entidade '%s' não está ativa.", entidade.GetNomeExibicao()))
		}
	}

	// 3. Validar Itens
	if len(req.Itens) == 0 && !isUpdate { // Na criação, pelo menos um item é obrigatório
		return apperrors.NewValidationError("O documento deve ter pelo menos um item.")
	}

	for i, item := range req.Itens {
		if item.ProdutoID <= 0 {
			return apperrors.NewValidationError(fmt.Sprintf("O 'produto_id' é obrigatório para o item %d.", i+1))
		}
		if item.Quantidade <= 0 {
			return apperrors.NewValidationError(fmt.Sprintf("A 'quantidade' deve ser maior que zero para o item %d.", i+1))
		}
		if item.ValorUnitario < 0 {
			return apperrors.NewValidationError(fmt.Sprintf("O 'valor_unitario' não pode ser negativo para o item %d.", i+1))
		}

		// Valida se o produto existe e está ativo
		produto, err := s.proService.GetByID(item.ProdutoID)
		if err != nil {
			return apperrors.NewInternalError("Erro ao verificar produto.", err)
		}
		if produto == nil {
			return apperrors.NewNotFoundError(fmt.Sprintf("Produto com ID %d não encontrado para o item %d.", item.ProdutoID, i+1))
		}
	}

	return nil
}

// recalcularTotais calcula os totais do documento com base nos itens.
func (s *DocumentoVendaService) recalcularTotais(doc *models.DocumentoVenda) {
	var totalProdutos, totalDescontos, totalPesoBruto, totalPesoLiquido float64

	for i := range doc.Itens {
		item := &doc.Itens[i]
		item.TotalProdutos = item.Quantidade * item.ValorUnitario
		item.TotalItem = item.TotalProdutos

		if item.ValorDesconto != nil {
			item.TotalItem -= *item.ValorDesconto
			totalDescontos += *item.ValorDesconto
		}

		totalProdutos += item.TotalProdutos
		totalPesoBruto += item.PesoBruto
		totalPesoLiquido += item.PesoLiquido
	}

	doc.TotalProdutos = totalProdutos
	doc.TotalDescontos = totalDescontos
	doc.TotalPesoBruto = totalPesoBruto
	doc.TotalPesoLiquido = totalPesoLiquido

	// Calcula o total líquido
	doc.TotalLiquido = doc.TotalProdutos - doc.TotalDescontos
	if doc.ValorFrete != nil {
		doc.TotalLiquido += *doc.ValorFrete
	}
}

// recalcularTotaisPagamentos calcula os totais de pagamento do documento.
func (s *DocumentoVendaService) recalcularTotaisPagamentos(doc *models.DocumentoVenda) {
	var totalPago float64

	for _, p := range doc.Pagamentos {
		if !p.IsDeleted() {
			totalPago += p.Valor
		}
	}

	doc.ValorPago = &totalPago
	// A lógica de troco pode ser mais complexa, mas por enquanto vamos zerar se não houver valor pago.
}

// ============================================================
// MÉTODOS PÚBLICOS (CRUD)
// ============================================================

func (s *DocumentoVendaService) Create(req *dto.DocumentoVendaRequest) (*models.DocumentoVenda, error) {
	if err := s.isDataValid(req, false); err != nil {
		return nil, err
	}

	doc := &models.DocumentoVenda{}
	if err := utils.MapToModel(req, doc); err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear dados do documento.", err)
	}

	// Definir valores padrão
	doc.DataDocumento = time.Now()
	doc.Situacao = constants.SituacaoPedidoAberto

	// Recalcular totais
	s.recalcularTotais(doc)

	// Usar transação para garantir a atomicidade
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, apperrors.NewInternalError("Erro ao iniciar transação.", tx.Error)
	}

	if err := s.ddvRepo.Create(doc); err != nil {
		tx.Rollback()
		return nil, apperrors.NewInternalError("Erro ao criar documento.", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, apperrors.NewInternalError("Erro ao commitar transação.", err)
	}

	return s.GetByID(doc.ID)
}

func (s *DocumentoVendaService) FindByID(id int) (*models.DocumentoVenda, error) {
	if id <= 0 {
		return nil, apperrors.NewValidationError("ID do documento inválido.")
	}
	return s.ddvRepo.FindByID(id)
}

func (s *DocumentoVendaService) GetByID(id int) (*models.DocumentoVenda, error) {
	if id <= 0 {
		return nil, apperrors.NewValidationError("ID do documento inválido.")
	}
	return s.ddvRepo.GetByID(id)
}

func (s *DocumentoVendaService) Update(id int, req *dto.DocumentoVendaRequest) (*models.DocumentoVenda, error) {
	if err := s.isDataValid(req, true); err != nil {
		return nil, err
	}

	doc, err := s.ddvRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if !doc.IsActive() {
		return nil, apperrors.NewValidationError("Apenas documentos com situação 'Aberto' ou 'Em Atividade' podem ser alterados.")
	}

	// Mapeia os dados do DTO para o modelo, preservando os itens existentes
	if err := utils.MapToModel(req, doc); err != nil {
		return nil, apperrors.NewInternalError("Erro ao mapear dados do documento.", err)
	}

	s.recalcularTotais(doc)

	if err := s.ddvRepo.Update(doc); err != nil {
		return nil, apperrors.NewInternalError("Erro ao atualizar documento.", err)
	}

	return s.GetByID(id)
}

func (s *DocumentoVendaService) List(limit, offset int, filters map[string]any) ([]models.DocumentoVenda, int64, error) {
	return s.ddvRepo.List(limit, offset, filters)
}

func (s *DocumentoVendaService) Delete(id int) error {
	doc, err := s.ddvRepo.FindByID(id)
	if err != nil {
		return err
	}

	if doc.IsFechado() {
		return apperrors.NewValidationError("Não é possível excluir um documento já fechado.")
	}

	return s.ddvRepo.Delete(id)
}

// executarOperacaoItem gerencia a transação para operações de item (Adicionar, Editar, Deletar).
func (s *DocumentoVendaService) executarOperacaoItem(
	docVendaID int,
	itemOperation func(dviService *DocumentoVendaItemService) error,
) (*models.DocumentoVenda, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, apperrors.NewInternalError("Erro ao iniciar transação.", tx.Error)
	}

	// Defer a rollback. Se o commit for bem-sucedido, o rollback não terá efeito.
	// Se ocorrer um pânico, o rollback será executado.
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // Re-lança o pânico
		}
	}()

	// Executa a operação do item (Create, Update, Delete) dentro da transação.
	dviServiceTx := NewDocumentoVendaItemService(tx)
	if err := itemOperation(dviServiceTx); err != nil {
		tx.Rollback()
		return nil, err // Retorna o erro original da operação do item.
	}

	// Instancia o repositório do documento com a transação.
	ddvRepoTx := repository.NewDocumentoVendaRepository(tx)

	// Busca o documento de venda completo para recalcular os totais.
	doc, err := ddvRepoTx.FindByID(docVendaID)
	if err != nil {
		tx.Rollback()
		return nil, apperrors.NewInternalError("Erro ao buscar documento após operação no item.", err)
	}

	// Recalcula e salva os totais atualizados.
	s.recalcularTotais(doc)
	if err := ddvRepoTx.Update(doc); err != nil {
		tx.Rollback()
		return nil, apperrors.NewInternalError("Erro ao atualizar totais do documento via repositório.", err)
	}

	// Commita a transação.
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // Garante o rollback em caso de falha no commit.
		return nil, apperrors.NewInternalError("Erro ao commitar transação.", err)
	}

	return doc, nil
}

func (s *DocumentoVendaService) AddItem(reqItem *dto.DocumentoVendaItemRequest) (*models.DocumentoVenda, error) {
	return s.executarOperacaoItem(reqItem.DocumentoVendaID, func(dviService *DocumentoVendaItemService) error {
		// A função `Create` do dviService precisa ser criada ou ajustada para receber o DTO.
		// Assumindo que ela exista:
		return dviService.Create(reqItem)
	})
}

func (s *DocumentoVendaService) EditItem(ddvId, dviItem int, reqItem *dto.DocumentoVendaItemRequest) (*models.DocumentoVenda, error) {
	return s.executarOperacaoItem(ddvId, func(dviService *DocumentoVendaItemService) error {
		// Delega a operação de atualização para o serviço de item.
		return dviService.Update(ddvId, dviItem, reqItem)
	})
}

func (s *DocumentoVendaService) DeleteItem(ddvId, dviItem int) (*models.DocumentoVenda, error) {
	return s.executarOperacaoItem(ddvId, func(dviService *DocumentoVendaItemService) error {
		// Delega a operação de exclusão para o serviço de item.
		return dviService.Delete(ddvId, dviItem)
	})
}

// executarOperacaoPagamento gerencia a transação para operações de pagamento.
func (s *DocumentoVendaService) executarOperacaoPagamento(
	docVendaID int,
	pagamentoOperation func(dvpService *DocumentoVendaPagamentoService) error,
) (*models.DocumentoVenda, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, apperrors.NewInternalError("Erro ao iniciar transação.", tx.Error)
	}
	defer tx.Rollback()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Executa a operação de pagamento dentro da transação.
	dvpServiceTx := NewDocumentoVendaPagamentoService(tx)
	if err := pagamentoOperation(dvpServiceTx); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Instancia o repositório do documento com a transação.
	ddvRepoTx := repository.NewDocumentoVendaRepository(tx)

	// Busca o documento de venda completo para recalcular os totais.
	doc, err := ddvRepoTx.FindByID(docVendaID)
	if err != nil {
		tx.Rollback()
		return nil, apperrors.NewInternalError("Erro ao buscar documento após operação no pagamento.", err)
	}

	// Recalcula e salva os totais de pagamento atualizados.
	s.recalcularTotaisPagamentos(doc)
	if err := ddvRepoTx.Update(doc); err != nil {
		tx.Rollback()
		return nil, apperrors.NewInternalError("Erro ao atualizar totais de pagamento do documento.", err)
	}

	// Commita a transação.
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, apperrors.NewInternalError("Erro ao commitar transação.", err)
	}

	return doc, nil
}

func (s *DocumentoVendaService) AddPagamento(reqAddPagamento *dto.DocumentoVendaPagamentoRequest) (*models.DocumentoVenda, error) {
	return s.executarOperacaoPagamento(
		reqAddPagamento.DocumentoVendaID,
		func(dvpService *DocumentoVendaPagamentoService) error {
			return dvpService.Create(reqAddPagamento)
		})
}

func (s *DocumentoVendaService) EditPagamento(ddvId, dvpItem int, reqPagamento *dto.DocumentoVendaPagamentoRequest) (*models.DocumentoVenda, error) {
	return s.executarOperacaoPagamento(ddvId,
		func(dvpService *DocumentoVendaPagamentoService) error {
			return dvpService.Update(ddvId, dvpItem, reqPagamento)
		})
}

func (s *DocumentoVendaService) DeletePagamento(ddvId, dvpItem int) (*models.DocumentoVenda, error) {
	return s.executarOperacaoPagamento(ddvId,
		func(dvpService *DocumentoVendaPagamentoService) error {
			return dvpService.Delete(ddvId, dvpItem)
		})
}
