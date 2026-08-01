package handler

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

// ProdutoVariacaoHandler gerencia as requisições HTTP para variações de produto.
type ProdutoVariacaoHandler struct {
	service service.ProdutoVariacaoService
}

// NewProdutoVariacaoHandler cria uma nova instância de ProdutoVariacaoHandler.
func NewProdutoVariacaoHandler(s service.ProdutoVariacaoService) *ProdutoVariacaoHandler {
	return &ProdutoVariacaoHandler{service: s}
}

// @Summary Cria uma nova variação de produto
// @Description Cria uma nova variação de produto com base nos dados fornecidos.
// @Tags Produto - Variações
// @Accept json
// @Produce json
// @Param variacao body dto.ProdutoVariacaoRequest true "Dados para criar a variação de produto"
// @Success 201 {object} dto.ProdutoVariacaoResponse
// @Failure 400 {object} utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure 409 {object} utils.ErrorResponse "Conflito, por exemplo, SKU já existente"
// @Failure 500 {object} utils.ErrorResponse "Erro interno do servidor"
// @Router /produto-variacoes [post]
func (h *ProdutoVariacaoHandler) Create(c *gin.Context) {
	var req dto.ProdutoVariacaoRequest //
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	// TODO: Obter CreatedBy do contexto da requisição
	// if userID, exists := c.Get("user_id"); exists {
	// 	id := int(userID.(float64)) // ou int(userID.(int)) dependendo do tipo
	// 	req.CreatedBy = &id
	// }

	resp, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	utils.RespondWithCreated(c, resp)
}

// @Summary Busca uma variação de produto por ID
// @Description Retorna os detalhes de uma variação de produto específica.
// @Tags Produto - Variações
// @Accept json
// @Produce json
// @Param id path int true "ID da Variação de Produto"
// @Success 200 {object} dto.ProdutoVariacaoResponse
// @Failure 400 {object} utils.ErrorResponse "ID da variação de produto inválido"
// @Failure 404 {object} utils.ErrorResponse "Variação de produto não encontrada"
// @Failure 500 {object} utils.ErrorResponse "Erro interno do servidor"
// @Router /produto-variacoes/{id} [get]
func (h *ProdutoVariacaoHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondWithErrorAny(c, apperrors.NewValidationError(fmt.Sprintf("ID da variação de produto inválido: %v", err)))
		return
	}

	resp, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	utils.RespondWithOK(c, resp)
}

// @Summary Atualiza uma variação de produto
// @Description Atualiza os dados de uma variação de produto existente.
// @Tags Produto - Variações
// @Accept json
// @Produce json
// @Param id path int true "ID da Variação de Produto"
// @Param variacao body dto.ProdutoVariacaoRequest true "Dados para atualizar a variação de produto"
// @Success 200 {object} dto.ProdutoVariacaoResponse
// @Failure 400 {object} utils.ErrorResponse "Erro de validação ou dados inválidos"
// @Failure 404 {object} utils.ErrorResponse "Variação de produto não encontrada"
// @Failure 409 {object} utils.ErrorResponse "Conflito, por exemplo, SKU já existente"
// @Failure 500 {object} utils.ErrorResponse "Erro interno do servidor"
// @Router /produto-variacoes/{id} [put]
func (h *ProdutoVariacaoHandler) Update(c *gin.Context) {
	id := utils.GetQueryInt(c, "id", 0)
	if id <= 0 {
		utils.RespondWithErrorAny(c, apperrors.NewValidationError("ID da variação de produto inválido."))
		return
	}

	var req dto.ProdutoVariacaoRequest //
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}

	// TODO: Obter UpdatedBy do contexto da requisição
	// if userID, exists := c.Get("user_id"); exists {
	// 	id := int(userID.(float64))
	// 	req.UpdatedBy = &id
	// }

	resp, err := h.service.Update(id, &req)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	utils.RespondWithOK(c, resp)
}

// @Summary Exclui uma variação de produto
// @Description Realiza a exclusão lógica de uma variação de produto.
// @Tags Produto - Variações
// @Accept json
// @Produce json
// @Param id path int true "ID da Variação de Produto"
// @Success 204 "Nenhum conteúdo"
// @Failure 400 {object} utils.ErrorResponse "Erro ao excluir (ex: variação com estoque)"
// @Failure 404 {object} utils.ErrorResponse "Variação de produto não encontrada"
// @Failure 500 {object} utils.ErrorResponse "Erro interno do servidor"
// @Router /produto-variacoes/{id} [delete]
func (h *ProdutoVariacaoHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.RespondWithErrorAny(c, apperrors.NewValidationError(fmt.Sprintf("ID da variação de produto inválido: %v", err)))
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Lista as variações de produto
// @Description Retorna uma lista paginada de variações de produto, com suporte a filtros.
// @Tags Produto - Variações
// @Accept json
// @Produce json
// @Param limit query int false "Número de registros por página" default(10)
// @Param offset query int false "Offset para a paginação" default(0)
// @Param produto_id query int false "Filtrar por ID do Produto"
// @Param empresa_filial_id query int false "Filtrar por ID da Empresa Filial"
// @Param cor_id query int false "Filtrar por ID da Cor"
// @Param tamanho_id query int false "Filtrar por ID do Tamanho"
// @Param sku query string false "Filtrar por SKU"
// @Success 200 {object} dto.ProdutoVariacaoListResponse
// @Failure 500 {object} utils.ErrorResponse "Erro interno do servidor"
// @Router /produto-variacoes [get]
func (h *ProdutoVariacaoHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	filters := make(map[string]interface{})
	if produtoIDStr := c.Query("produto_id"); produtoIDStr != "" {
		if produtoID, err := strconv.Atoi(produtoIDStr); err == nil {
			filters["produto_id"] = produtoID
		}
	}
	if empresaFilialIDStr := c.Query("empresa_filial_id"); empresaFilialIDStr != "" {
		if empresaFilialID, err := strconv.Atoi(empresaFilialIDStr); err == nil {
			filters["empresa_filial_id"] = empresaFilialID
		}
	}
	if corIDStr := c.Query("cor_id"); corIDStr != "" {
		if corID, err := strconv.Atoi(corIDStr); err == nil {
			filters["cor_id"] = corID
		}
	}
	if tamanhoIDStr := c.Query("tamanho_id"); tamanhoIDStr != "" {
		if tamanhoID, err := strconv.Atoi(tamanhoIDStr); err == nil {
			filters["tamanho_id"] = tamanhoID
		}
	}
	if sku := c.Query("sku"); sku != "" {
		filters["sku"] = sku
	}

	responses, total, err := h.service.List(limit, offset, filters)
	if err != nil {
		utils.RespondWithErrorAny(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ProdutoVariacaoListResponse{
		Items:      responses,
		Total:      total,
		Limit:      limit,
		Page:       (offset / limit) + 1,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	})
}
