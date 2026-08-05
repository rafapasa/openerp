package service

import (
	"fmt"
	"time"

	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/constants"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/utils"
	"gorm.io/gorm"
)

// DocumentoVendaService define os métodos públicos para o serviço de documento de venda.
type DocumentoVendaService interface {
	Create(req *dto.DocumentoVendaRequest) (*models.DocumentoVenda, error)
	FindByID(id int) (*models.DocumentoVenda, error)
	GetByID(id int) (*models.DocumentoVenda, error)
	Update(id int, req *dto.DocumentoVendaRequest) (*models.DocumentoVenda, error)
	List(limit, offset int, filters map[string]any) ([]models.DocumentoVenda, int64, error)
	Delete(id int) error
	ListItem(ddvId int) ([]models.DocumentoVendaItem, error)
	AddItem(reqItem *dto.DocumentoVendaItemRequest) (*models.DocumentoVenda, error)
	EditItem(ddvId, dviItem int, reqItem *dto.DocumentoVendaItemRequest) (*models.DocumentoVenda, error)
	DeleteItem(ddvId, dviItem int) (*models.DocumentoVenda, error)
	ListPagamento(ddvId int) ([]models.DocumentoVendaPagamento, error)
	AddPagamento(reqAddPagamento *dto.DocumentoVendaPagamentoRequest) (*models.DocumentoVenda, error)
	EditPagamento(ddvId, dvpItem int, reqPagamento *dto.DocumentoVendaPagamentoRequest) (*models.DocumentoVenda, error)
	DeletePagamento(ddvId, dvpItem int) (*models.DocumentoVenda, error)
}

// documentoVendaService é a implementação concreta de DocumentoVendaService.
type documentoVendaService struct {
	db         *gorm.DB
	ddvRepo    repository.DocumentoVendaRepository
	dviService DocumentoVendaItemService
	dvpService DocumentoVendaPagamentoService
	entService EntidadeService
	// proService        ProdutoService
	prcService        ProcessoService
	dviServiceFactory DocumentoVendaItemServiceFactory
}

func NewDocumentoVendaService(db *gorm.DB,
	ddvRepo repository.DocumentoVendaRepository,
	dviService DocumentoVendaItemService,
	dvpService DocumentoVendaPagamentoService,
	entService EntidadeService,
	// proService ProdutoService,
	prcService ProcessoService,
	dviServiceFactory DocumentoVendaItemServiceFactory) DocumentoVendaService {
	return &documentoVendaService{
		db:         db,
		ddvRepo:    ddvRepo,
		dviService: dviService,
		dvpService: dvpService,
		entService: entService,
		// proService:        proService,
		prcService:        prcService,
		dviServiceFactory: dviServiceFactory,
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO (PRIVADOS)
// ============================================================

func (s *documentoVendaService) isDataValid(req *dto.DocumentoVendaRequest, isUpdate bool) error {
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
			return err
		}
		if !entidade.IsActive() {
			return apperrors.NewValidationError(fmt.Sprintf("A entidade '%s' não está ativa.", entidade.GetNomeExibicao()))
		}
	}

	// 3. Validar Itens
	if len(req.Itens) == 0 && !isUpdate { // Na criação, pelo menos um item é obrigatório
		return apperrors.NewValidationError("O documento deve ter pelo menos um item.")
	}

	// for i, item := range req.Itens {
	// 	if err := item.Validate(); err != nil {
	// 		return fmt.Errorf("validação do item %d falhou: %w", i+1, err)
	// 	}

	// 	// Valida se o produto existe e está ativo
	// 	produto, err := s.proService.GetByID(item.ProdutoID)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if produto == nil {
	// 		return apperrors.NewNotFoundError(fmt.Sprintf("Produto com ID %d não encontrado para o item %d.", item.ProdutoID, i+1))
	// 	}
	// }

	return nil
}

// recalcularTotais calcula os totais do documento com base nos itens.
func (s *documentoVendaService) recalcularTotais(doc *models.DocumentoVenda) {
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
func (s *documentoVendaService) recalcularTotaisPagamentos(doc *models.DocumentoVenda) {
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

func (s *documentoVendaService) Create(req *dto.DocumentoVendaRequest) (*models.DocumentoVenda, error) {
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

// Create cria um novo documento de venda.
// Este método já existe no arquivo, apenas garantindo que a interface o referencie.
// (No diff, ele já está presente, então não há mudança real aqui, apenas a menção para a interface)

func (s *documentoVendaService) FindByID(id int) (*models.DocumentoVenda, error) {
	if id <= 0 {
		return nil, apperrors.NewValidationError("ID do documento inválido.")
	}
	return s.ddvRepo.FindByID(id)
}

func (s *documentoVendaService) GetByID(id int) (*models.DocumentoVenda, error) {
	if id <= 0 {
		return nil, apperrors.NewValidationError("ID do documento inválido.")
	}
	return s.ddvRepo.GetByID(id)
}

func (s *documentoVendaService) Update(id int, req *dto.DocumentoVendaRequest) (*models.DocumentoVenda, error) {
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
		return nil, apperrors.NewInternalError("Erro ao atualizar documento.", err) //
	}

	return s.GetByID(id)
}

func (s *documentoVendaService) List(limit, offset int, filters map[string]any) ([]models.DocumentoVenda, int64, error) {
	return s.ddvRepo.List(limit, offset, filters)
}

func (s *documentoVendaService) Delete(id int) error {
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
func (s *documentoVendaService) executarOperacaoItem(docVendaID int,
	itemOperation func(dviServiceTx DocumentoVendaItemService) error) (*models.DocumentoVenda, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, apperrors.NewInternalError("Erro ao iniciar transação.", tx.Error) //
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
	dviServiceTx := s.dviServiceFactory.CreateWithTx(tx) // Retorna a interface
	if err := itemOperation(dviServiceTx); err != nil {
		tx.Rollback()
		return nil, err // Retorna o erro original da operação do item.
	}

	// Instancia o repositório do documento com a transação.
	ddvRepoTx := repository.NewDocumentoVendaRepository(tx)

	// Busca o documento de venda completo para recalcular os totais.
	doc, err := ddvRepoTx.FindByID(docVendaID)
	if err != nil {
		tx.Rollback()                                                                                  //
		return nil, apperrors.NewInternalError("Erro ao buscar documento após operação no item.", err) //
	}

	// Recalcula e salva os totais atualizados.
	s.recalcularTotais(doc)
	if err := ddvRepoTx.Update(doc); err != nil {
		tx.Rollback()                                                                                         //
		return nil, apperrors.NewInternalError("Erro ao atualizar totais do documento via repositório.", err) //
	}

	// Commita a transação.
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()                                                              // Garante o rollback em caso de falha no commit. //
		return nil, apperrors.NewInternalError("Erro ao commitar transação.", err) //
	}

	return doc, nil
}

func (s *documentoVendaService) ListItem(ddvId int) ([]models.DocumentoVendaItem, error) {
	if ddvId <= 0 {
		return nil, apperrors.NewValidationError("ID do documento inválido.")
	}
	itens, _, err := s.dviService.List(ddvId)
	if err != nil {
		return nil, err
	}
	return itens, nil
}

func (s *documentoVendaService) AddItem(reqItem *dto.DocumentoVendaItemRequest) (*models.DocumentoVenda, error) {
	ddv, err := s.GetByID(reqItem.DocumentoVendaID)
	if err != nil {
		return nil, err
	}

	opInterna := true
	opSubTrib := true

	prcId := ddv.ProcessoID
	opf, err := s.prcService.GetOperacaoFiscal(prcId, opInterna, opSubTrib)
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao buscar operação fiscal.", err)
	}

	reqItem.OperacaoFiscalID = utils.IntPtr(opf.ID)
	reqItem.CstIcmsId = utils.IntPtr(*opf.CSTICMSID)
	reqItem.CstIpiId = utils.IntPtr(*opf.CSTIPIID)
	reqItem.CstPisCofinsId = utils.IntPtr(*opf.CSTPISCOFINSID)

	return s.executarOperacaoItem(reqItem.DocumentoVendaID, func(dviService DocumentoVendaItemService) error {
		// A função `Create` do dviService precisa ser criada ou ajustada para receber o DTO.
		// Assumindo que ela exista:
		return dviService.Create(reqItem)
	})
}

func (s *documentoVendaService) EditItem(ddvId, dviItem int, reqItem *dto.DocumentoVendaItemRequest) (*models.DocumentoVenda, error) {
	return s.executarOperacaoItem(ddvId, func(dviServicTx DocumentoVendaItemService) error {
		return dviServicTx.Update(ddvId, dviItem, reqItem)
	})
}

func (s *documentoVendaService) DeleteItem(ddvId, dviItem int) (*models.DocumentoVenda, error) {
	return s.executarOperacaoItem(ddvId, func(dviService DocumentoVendaItemService) error {
		// Delega a operação de exclusão para o serviço de item. (This comment is slightly inaccurate, it should be DocumentoVendaItemServiceInterface)
		return dviService.Delete(ddvId, dviItem)
	})
}

// executarOperacaoPagamento gerencia a transação para operações de pagamento.
func (s *documentoVendaService) executarOperacaoPagamento(
	docVendaID int,
	pagamentoOperation func(dvpService DocumentoVendaPagamentoService) error,
) (*models.DocumentoVenda, error) {
	tx := s.db.Begin()
	if tx.Error != nil { //
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
	dvpServiceTx := NewDocumentoVendaPagamentoService(tx, repository.NewDocumentoVendaPagamentoRepository(tx))
	if err := pagamentoOperation(dvpServiceTx); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Instancia o repositório do documento com a transação.
	ddvRepoTx := repository.NewDocumentoVendaRepository(tx)

	// Busca o documento de venda completo para recalcular os totais.
	doc, err := ddvRepoTx.FindByID(docVendaID)
	if err != nil {
		tx.Rollback()                                                                                       //
		return nil, apperrors.NewInternalError("Erro ao buscar documento após operação no pagamento.", err) //
	}

	// Recalcula e salva os totais de pagamento atualizados.
	s.recalcularTotaisPagamentos(doc)
	if err := ddvRepoTx.Update(doc); err != nil {
		tx.Rollback()                                                                                      //
		return nil, apperrors.NewInternalError("Erro ao atualizar totais de pagamento do documento.", err) //
	}

	// Commita a transação.
	if err := tx.Commit().Error; err != nil { //
		tx.Rollback()
		return nil, apperrors.NewInternalError("Erro ao commitar transação.", err)
	}

	return doc, nil
}

func (s *documentoVendaService) ListPagamento(ddvId int) ([]models.DocumentoVendaPagamento, error) {
	if ddvId <= 0 {
		return nil, apperrors.NewValidationError("ID do documento inválido.")
	}
	pagamentos, _, err := s.dvpService.ListByDocumentoVendaID(ddvId)
	if err != nil {
		return nil, err
	}
	return pagamentos, nil
}

func (s *documentoVendaService) AddPagamento(reqAddPagamento *dto.DocumentoVendaPagamentoRequest) (*models.DocumentoVenda, error) {
	return s.executarOperacaoPagamento(reqAddPagamento.DocumentoVendaID, func(dvpServiceTx DocumentoVendaPagamentoService) error {
		// The Create method of dvpServiceTx will handle the model conversion and date parsing
		// No need to manually convert dates here, as ToModel in DTO already does it.
		// However, the `dvpServiceTx.Create` expects a DTO, so we pass the original DTO.
		// The DTO's ToModel method will be called internally by the service.
		// If the service's Create method directly uses the DTO, then this is fine.
		// If the service's Create method expects a model, then we need to convert here.
		// Looking at `documento_venda_pagamento_service.go`, `Create` takes a DTO.

		// Ensure DocumentoVendaID is set in the request before passing to service
		// This is already done in the handler, but good to be explicit if this path changes.
		// reqAddPagamento.DocumentoVendaID = docVendaID // Assuming docVendaID is available here

		// The `dvpServiceTx.Create` method will handle the DTO to Model conversion and date parsing.
		// So, we just pass the request DTO directly.
		return dvpServiceTx.Create(reqAddPagamento)
	})
}

func (s *documentoVendaService) EditPagamento(ddvId, dvpItem int, reqPagamento *dto.DocumentoVendaPagamentoRequest) (*models.DocumentoVenda, error) {
	return s.executarOperacaoPagamento(ddvId,
		func(dvpServiceTx DocumentoVendaPagamentoService) error {
			return dvpServiceTx.Update(ddvId, dvpItem, reqPagamento)
		})
}

func (s *documentoVendaService) DeletePagamento(ddvId, dvpItem int) (*models.DocumentoVenda, error) {
	return s.executarOperacaoPagamento(ddvId,
		func(dvpServiceTx DocumentoVendaPagamentoService) error {
			return dvpServiceTx.Delete(ddvId, dvpItem)
		})
}
