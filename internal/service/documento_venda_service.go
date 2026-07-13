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
	db           *gorm.DB
	docRepo      *repository.DocumentoVendaRepository
	entidadeRepo *repository.EntidadeRepository
	produtoRepo  *repository.ProdutoRepository
}

func NewDocumentoVendaService(db *gorm.DB) *DocumentoVendaService {
	return &DocumentoVendaService{
		db:           db,
		docRepo:      repository.NewDocumentoVendaRepository(db),
		entidadeRepo: repository.NewEntidadeRepository(db),
		produtoRepo:  repository.NewProdutoRepository(db),
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
		entidade, err := s.entidadeRepo.GetByID(*req.EntidadeID)
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
		produto, err := s.produtoRepo.ExistsByID(item.ProdutoID)
		if err != nil {
			return apperrors.NewInternalError("Erro ao verificar produto.", err)
		}
		if !produto {
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

	if err := s.docRepo.Create(doc); err != nil {
		tx.Rollback()
		return nil, apperrors.NewInternalError("Erro ao criar documento.", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, apperrors.NewInternalError("Erro ao commitar transação.", err)
	}

	return s.GetByID(doc.ID)
}

func (s *DocumentoVendaService) GetByID(id int) (*models.DocumentoVenda, error) {
	if id <= 0 {
		return nil, apperrors.NewValidationError("ID do documento inválido.")
	}
	return s.docRepo.FindByID(id)
}

func (s *DocumentoVendaService) Update(id int, req *dto.DocumentoVendaRequest) (*models.DocumentoVenda, error) {
	if err := s.isDataValid(req, true); err != nil {
		return nil, err
	}

	doc, err := s.docRepo.FindByID(id)
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

	if err := s.docRepo.Update(doc); err != nil {
		return nil, apperrors.NewInternalError("Erro ao atualizar documento.", err)
	}

	return s.GetByID(id)
}

func (s *DocumentoVendaService) List(limit, offset int, filters map[string]interface{}) ([]models.DocumentoVenda, int64, error) {
	return s.docRepo.List(limit, offset, filters)
}

func (s *DocumentoVendaService) Delete(id int) error {
	doc, err := s.docRepo.FindByID(id)
	if err != nil {
		return err
	}

	if doc.IsFechado() {
		return apperrors.NewValidationError("Não é possível excluir um documento já fechado.")
	}

	return s.docRepo.Delete(id)
}
