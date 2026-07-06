package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/openerp/backend/internal/dto"
	"github.com/openerp/backend/internal/middleware"
	"github.com/openerp/backend/internal/service"
	"github.com/openerp/backend/internal/utils"
)

type ProdutoHandler struct {
	service *service.ProdutoService
}

func NewProdutoHandler(service *service.ProdutoService) *ProdutoHandler {
	return &ProdutoHandler{
		service: service,
	}
}

func (h *ProdutoHandler) Create(c *gin.Context) {
	var req dto.ProdutoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}
	userID := middleware.GetUserID(c)
	req.CreatedBy = &userID
	req.UpdatedBy = &userID
	produto, err := h.service.Create(&req)
	if err != nil {
		utils.RespondWithValidationError(c, err.Error())
		return
	}
	var resp dto.ProdutoResponse
	resp.FromModel(produto)
	utils.RespondWithCreated(c, resp)
}

func (h *ProdutoHandler) GetByID(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}
	produto, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondWithNotFoundError(c, err.Error())
		return
	}
	var resp dto.ProdutoResponse
	resp.FromModel(produto)
	utils.RespondWithOK(c, resp)
}

func (h *ProdutoHandler) Update(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}
	var req dto.ProdutoRequest
	if !utils.BindAndValidateOrRespond(c, &req) {
		return
	}
	userID := middleware.GetUserID(c)
	req.UpdatedBy = &userID
	produto, err := h.service.Update(id, &req)
	if err != nil {
		utils.RespondWithNotFoundError(c, err.Error())
		return
	}
	var resp dto.ProdutoResponse
	resp.FromModel(produto)
	utils.RespondWithOK(c, resp)
}

func (h *ProdutoHandler) Delete(c *gin.Context) {
	id, ok := utils.ParseIDParam(c, "id")
	if !ok {
		return
	}
	err := h.service.Delete(id)
	if err != nil {
		utils.RespondWithNotFoundError(c, err.Error())
		return
	}
	utils.RespondWithNoContent(c)
}
