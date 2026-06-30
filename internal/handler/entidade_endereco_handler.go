package handler

import "github.com/openerp/backend/internal/service"

type EntidadeEnderecoHandler struct {
	service *service.EntidadeEnderecoService
}

func NewEntidadeEnderecoHandler(service *service.EntidadeEnderecoService) EntidadeEnderecoHandler{
	return EntidadeEnderecoHandler{
		service: service,
	}
}

