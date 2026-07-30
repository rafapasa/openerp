// produto_service.go
package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/constants"
	"github.com/openerp/backend/internal/dto"
	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
	"github.com/openerp/backend/internal/utils"
)

// ProdutoService define os métodos públicos para o serviço de produto.
type ProdutoService interface {
	Create(req *dto.ProdutoRequest) (*models.Produto, error)
	GetByID(id int) (*models.Produto, error)
	GetByCodigo(codigo int) (*models.Produto, error)
	GetByCodigoBarras(codigoBarras string) (*models.Produto, error)
	Update(id int, req *dto.ProdutoRequest) (*models.Produto, error)
	Delete(id int) error
	List(limit, offset int, filters map[string]interface{}) ([]models.Produto, int64, error)
	Activate(id int) error
	Deactivate(id int) error
	FindById(id int) (*models.Produto, error)
	GetValorUnitario(tbpId, proId int) (float64, error)
}

type produtoService struct {
	proRepo     repository.ProdutoRepository
	tbppService TabelaPrecoProdutoService
	// estoqueRepo        *repository.EstoqueRepository
	// documentoVendaItemRepo *repository.DocumentoVendaItemRepository
}

// ============================================================
// CONSTANTES DE VALIDAÇÃO
// ============================================================

const (
	maxLengthProdutoNome   = 80
	maxLengthProdutoResumo = 80
	maxLengthCodigoBarras  = 255
	maxLengthReferencia    = 255
	minLengthProdutoNome   = 3
	minLengthProdutoResumo = 3
)

func NewProdutoService(proRepo repository.ProdutoRepository, tbppService TabelaPrecoProdutoService) ProdutoService {
	return &produtoService{
		proRepo:     proRepo,
		tbppService: tbppService,
	}
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO (PRIVADOS)
// ============================================================

// isDataValid realiza as validações básicas de um produto.
func (s *produtoService) isDataValid(req *dto.ProdutoRequest) error {
	// 1. Validar campos obrigatórios
	if err := s.validateRequiredFields(req); err != nil {
		return err
	}

	// 2. Validar tamanhos dos campos
	if err := s.validateFieldLengths(req); err != nil {
		return err
	}

	// 3. Validar campos opcionais com validação
	if err := s.validateOptionalFields(req); err != nil {
		return err
	}

	// 4. Validar relacionamentos (FKs)
	if err := s.validateRelationships(req); err != nil {
		return err
	}

	return nil
}

// validateRequiredFields valida campos obrigatórios
func (s *produtoService) validateRequiredFields(req *dto.ProdutoRequest) error {
	if strings.TrimSpace(req.Nome) == "" {
		return errors.New("nome do produto é obrigatório")
	}

	if strings.TrimSpace(req.Resumo) == "" {
		return errors.New("resumo do produto é obrigatório")
	}

	if req.EmpresaFilialID <= 0 {
		return errors.New("empresa/filial é obrigatória")
	}

	if req.TipoProdutoID <= 0 {
		return errors.New("tipo de produto é obrigatório")
	}

	if req.NcmNumero <= 0 {
		return errors.New("NCM é obrigatório")
	}

	if req.Codigo <= 0 {
		return errors.New("código do produto é obrigatório")
	}

	return nil
}

// validateFieldLengths valida o tamanho dos campos
func (s *produtoService) validateFieldLengths(req *dto.ProdutoRequest) error {
	// Nome
	if len(req.Nome) > maxLengthProdutoNome {
		return fmt.Errorf("nome do produto deve ter no máximo %d caracteres", maxLengthProdutoNome)
	}
	if len(req.Nome) < minLengthProdutoNome {
		return fmt.Errorf("nome do produto deve ter pelo menos %d caracteres", minLengthProdutoNome)
	}

	// Resumo
	if len(req.Resumo) > maxLengthProdutoResumo {
		return fmt.Errorf("resumo do produto deve ter no máximo %d caracteres", maxLengthProdutoResumo)
	}
	if len(req.Resumo) < minLengthProdutoResumo {
		return fmt.Errorf("resumo do produto deve ter pelo menos %d caracteres", minLengthProdutoResumo)
	}

	// Código de Barras (se informado)
	if req.CodigoBarras != "" && len(req.CodigoBarras) > maxLengthCodigoBarras {
		return fmt.Errorf("código de barras deve ter no máximo %d caracteres", maxLengthCodigoBarras)
	}

	// Referências (se informadas)
	referencias := map[string]*string{
		"referência":   &req.Referencia,
		"referência 2": &req.Referencia2,
		"referência 3": &req.Referencia3,
		"referência 4": &req.Referencia4,
	}
	for nome, ref := range referencias {
		if ref != nil && len(*ref) > maxLengthReferencia {
			return fmt.Errorf("%s deve ter no máximo %d caracteres", nome, maxLengthReferencia)
		}
	}

	return nil
}

// validateOptionalFields valida campos opcionais
func (s *produtoService) validateOptionalFields(req *dto.ProdutoRequest) error {
	// Validar se valores numéricos não são negativos
	valores := map[string]*float64{
		"custo de compra": req.CustoCompra,
		"peso bruto":      req.PesoBruto,
		"peso líquido":    req.PesoLiquido,
		"altura":          req.Altura,
		"largura":         req.Largura,
		"comprimento":     req.Comprimento,
		"estoque mínimo":  req.EstoqueMinimo,
		"lote econômico":  req.LoteEconomico,
		"desconto máximo": req.DescontoMaximo,
	}

	for nome, valor := range valores {
		if valor != nil && *valor < 0 {
			return fmt.Errorf("%s não pode ser negativo", nome)
		}
	}

	return nil
}

// validateRelationships valida as chaves estrangeiras
func (s *produtoService) validateRelationships(req *dto.ProdutoRequest) error {
	// Os repositórios de relacionamento seriam injetados para validação
	// Aqui faremos apenas validações básicas

	// Validar se os IDs são consistentes
	if req.ProdutoGrupoID != nil && *req.ProdutoGrupoID <= 0 {
		return errors.New("ID do grupo de produto inválido")
	}

	if req.ProdutoSubgrupoID != nil && *req.ProdutoSubgrupoID <= 0 {
		return errors.New("ID do subgrupo de produto inválido")
	}

	if req.MarcaID != nil && *req.MarcaID <= 0 {
		return errors.New("ID da marca inválido")
	}

	if req.ModeloID != nil && *req.ModeloID <= 0 {
		return errors.New("ID do modelo inválido")
	}

	if req.SerieID != nil && *req.SerieID <= 0 {
		return errors.New("ID da série inválido")
	}

	if req.EspecieID != nil && *req.EspecieID <= 0 {
		return errors.New("ID da espécie inválido")
	}

	return nil
}

// validateUniqueCodigo verifica se o código do produto já existe
func (s *produtoService) validateUniqueCodigo(codigo int, excludeID int) error {
	existe, err := s.proRepo.ExistsByCodigo(codigo, excludeID)
	if err != nil {
		return fmt.Errorf("erro ao verificar duplicidade de código: %w", err)
	}
	if existe {
		return fmt.Errorf("código %d já está cadastrado", codigo)
	}
	return nil
}

// validateUniqueCodigoBarras verifica se o código de barras já existe
func (s *produtoService) validateUniqueCodigoBarras(codigoBarras *string, excludeID int) error { // Corrected: Changed parameter type to *string
	if codigoBarras == nil || *codigoBarras == "" {
		return nil // Código de barras é opcional
	}

	existe, err := s.proRepo.ExistsByCodigoBarras(*codigoBarras, excludeID)
	if err != nil {
		return fmt.Errorf("erro ao verificar duplicidade de código de barras: %w", err)
	}
	if existe {
		return fmt.Errorf("código de barras %s já está cadastrado", *codigoBarras)
	}
	return nil
}

// isCreateValid valida dados para criação
func (s *produtoService) isCreateValid(req *dto.ProdutoRequest) error {
	// 1. Validações básicas
	if err := s.isDataValid(req); err != nil {
		return err
	}

	// 2. Validar duplicidade de código
	if err := s.validateUniqueCodigo(req.Codigo, 0); err != nil {
		return err
	}

	// 3. Validar duplicidade de código de barras
	if err := s.validateUniqueCodigoBarras(&req.CodigoBarras, 0); err != nil { // Corrected: Passed req.CodigoBarras directly
		return err
	}

	return nil
}

// isUpdateValid valida dados para atualização
func (s *produtoService) isUpdateValid(id int, req *dto.ProdutoRequest) error {
	// 1. Validações básicas
	if err := s.isDataValid(req); err != nil {
		return err
	}

	// 2. Validar duplicidade de código (excluindo o próprio ID)
	if err := s.validateUniqueCodigo(req.Codigo, id); err != nil {
		return err
	}

	// 3. Validar duplicidade de código de barras (excluindo o próprio ID)
	if err := s.validateUniqueCodigoBarras(&req.CodigoBarras, id); err != nil { // Corrected: Passed req.CodigoBarras directly
		return err
	}

	return nil
}

// ============================================================
// MÉTODOS PÚBLICOS (CRUD)
// ============================================================

// Create cria um novo produto
func (s *produtoService) Create(req *dto.ProdutoRequest) (*models.Produto, error) {
	if err := s.isCreateValid(req); err != nil {
		return nil, err
	}

	// Verificar se o código já existe no banco
	produtoExists, err := s.proRepo.FindByCodigo(req.Codigo)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("erro ao verificar código: %w", err)
	}
	if produtoExists != nil {
		return nil, fmt.Errorf("código %d já está cadastrado", req.Codigo)
	}

	produto := &models.Produto{}
	if err := utils.MapToModel(req, produto); err != nil {
		return nil, fmt.Errorf("erro ao mapear dados do produto: %w", err)
	}

	// Definir campos que não podem ser mapeados automaticamente
	produto.Situacao = constants.StatusAtivo
	produto.CreatedAt = time.Now()
	produto.UpdatedAt = time.Now()

	if err := s.proRepo.Create(produto); err != nil {
		return nil, fmt.Errorf("erro ao criar produto: %w", err)
	}

	return produto, nil
}

// GetByID busca um produto pelo ID
func (s *produtoService) GetByID(id int) (*models.Produto, error) {
	produto, err := s.proRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("produto com ID %d não encontrado", id)
		}
		return nil, fmt.Errorf("erro ao buscar produto: %w", err)
	}
	return produto, nil
}

// GetByCodigo busca um produto pelo código
func (s *produtoService) GetByCodigo(codigo int) (*models.Produto, error) {
	if codigo <= 0 {
		return nil, errors.New("código inválido")
	}

	produto, err := s.proRepo.FindByCodigo(codigo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("produto com código %d não encontrado", codigo)
		}
		return nil, fmt.Errorf("erro ao buscar produto: %w", err)
	}
	return produto, nil
}

// GetByCodigoBarras busca um produto pelo código de barras
func (s *produtoService) GetByCodigoBarras(codigoBarras string) (*models.Produto, error) {
	if strings.TrimSpace(codigoBarras) == "" {
		return nil, errors.New("código de barras não pode ser vazio")
	}

	produto, err := s.proRepo.FindByCodigoBarras(codigoBarras)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("produto com código de barras %s não encontrado", codigoBarras)
		}
		return nil, fmt.Errorf("erro ao buscar produto: %w", err)
	}
	return produto, nil
}

// Update atualiza um produto existente
func (s *produtoService) Update(id int, req *dto.ProdutoRequest) (*models.Produto, error) {
	if err := s.isUpdateValid(id, req); err != nil {
		return nil, err
	}

	// Buscar o produto existente
	produto, err := s.proRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("produto com ID %d não encontrado", id)
		}
		return nil, fmt.Errorf("erro ao buscar produto: %w", err)
	}

	// Mapear os novos dados para o modelo existente
	if err := utils.MapToModel(req, produto); err != nil {
		return nil, fmt.Errorf("erro ao mapear dados do produto: %w", err)
	}

	// Atualizar timestamp
	produto.UpdatedAt = time.Now()

	if err := s.proRepo.Update(id, produto); err != nil {
		return nil, fmt.Errorf("erro ao atualizar produto: %w", err)
	}

	return produto, nil
}

// Delete realiza a exclusão lógica de um produto
func (s *produtoService) Delete(id int) error {
	// Buscar o produto
	existe, err := s.proRepo.ExistsByID(id)
	if err != nil {
		return fmt.Errorf("erro ao buscar produto: %w", err)
	}

	if !existe {
		return fmt.Errorf("produto com ID %d não encontrado", id)
	}

	// Verificar dependências
	if err := s.checkDependencies(id); err != nil {
		return err
	}

	// Realizar exclusão lógica
	return s.proRepo.Delete(id)
}

// List lista produtos com paginação e filtros
func (s *produtoService) List(limit, offset int, filters map[string]interface{}) ([]models.Produto, int64, error) {
	// Validar parâmetros de paginação
	if limit <= 0 {
		limit = 10 // valor padrão
	}
	if offset < 0 {
		offset = 0
	}

	produtos, total, err := s.proRepo.List(limit, offset, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao listar produtos: %w", err)
	}
	return produtos, total, nil
}

// ============================================================
// MÉTODOS DE VALIDAÇÃO (REGRAS DE NEGÓCIO)
// ============================================================

// checkDependencies verifica se o produto tem dependências
// ESTE MÉTODO DEVE ESTAR NO SERVICE, NÃO NO REPOSITORY! (Receiver changed)
func (s *produtoService) checkDependencies(produtoID int) error {
	// 1. TODO: Verificar se tem itens em pedidos
	// itens, err := s.documentoVendaItemRepo.FindByProdutoID(produtoID)
	// if err != nil {
	//     return fmt.Errorf("erro ao verificar itens de pedido: %w", err)
	// }
	// if len(itens) > 0 {
	//     return errors.New("não é possível excluir produto com itens de pedido associados")
	// }

	// 2. TODO: Verificar se tem movimentações de estoque
	// movimentacoes, err := s.estoqueMovimentoRepo.FindByProdutoID(produtoID)
	// if err != nil {
	//     return fmt.Errorf("erro ao verificar movimentações de estoque: %w", err)
	// }
	// if len(movimentacoes) > 0 {
	//     return errors.New("não é possível excluir produto com movimentações de estoque")
	// }

	// 3. TODO: Verificar se tem preços cadastrados
	// precos, err := s.produtoPrecoRepo.FindByProdutoID(produtoID)
	// if err != nil {
	//     return fmt.Errorf("erro ao verificar preços do produto: %w", err)
	// }
	// if len(precos) > 0 {
	//     return errors.New("não é possível excluir produto com preços cadastrados")
	// }

	// 4. TODO: Verificar se tem em tabelas de preço
	// tabelasPreco, err := s.tabelaPrecoItemRepo.FindByProdutoID(produtoID)
	// if err != nil {
	//     return fmt.Errorf("erro ao verificar tabelas de preço: %w", err)
	// }
	// if len(tabelasPreco) > 0 {
	//     return errors.New("não é possível excluir produto com associações em tabelas de preço")
	// }

	// 5. TODO: Verificar se tem composição (produtos compostos)
	// composicao, err := s.produtoComposicaoRepo.FindByProdutoID(produtoID)
	// if err != nil {
	//     return fmt.Errorf("erro ao verificar composição do produto: %w", err)
	// }
	// if len(composicao) > 0 {
	//     return errors.New("não é possível excluir produto com composição definida")
	// }

	return nil
}

// ============================================================
// MÉTODOS DE ATUALIZAÇÃO DE ESTADO
// ============================================================

// Activate ativa um produto
func (s *produtoService) Activate(id int) error {
	produto, err := s.proRepo.FindByID(id)
	if err != nil {
		return err
	}

	if produto.Situacao == constants.StatusAtivo {
		return nil // Produto já está ativo, não faz nada
	}

	produto.Situacao = constants.StatusAtivo
	return s.proRepo.Update(id, produto)
}

// Deactivate desativa um produto
func (s *produtoService) Deactivate(id int) error {
	produto, err := s.proRepo.FindByID(id)
	if err != nil {
		return err
	}

	if produto.Situacao == constants.StatusInativo {
		return nil // Produto já está inativo, não faz nada
	}

	produto.Situacao = constants.StatusInativo
	return s.proRepo.Update(id, produto)
}

func (s *produtoService) FindById(id int) (*models.Produto, error) {
	if id <= 0 {
		return nil, apperrors.NewValidationError("O 'id' é obrigatório.")
	}
	return s.proRepo.FindByID(id)
}

func (s *produtoService) GetValorUnitario(tbpId, proId int) (float64, error) {
	if tbpId <= 0 {
		return 0, apperrors.NewValidationError("O 'tabela_preco_id' é obrigatório.")
	}
	if proId <= 0 {
		return 0, apperrors.NewValidationError("O 'produto_id' é obrigatório.")
	}
	itemTabPreco, err := s.tbppService.GetByProduto(tbpId, proId)
	if err != nil {
		return 0, err
	}
	return itemTabPreco.ValorPadrao, nil
}
