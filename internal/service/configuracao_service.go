package service

import (
	"github.com/openerp/backend/internal/repository"
)

// ConfiguracaoService define os métodos públicos para o serviço de configuração.
type ConfiguracaoService interface {
	GetConfig(configId int) (any, error)
}

type configuracaoService struct {
	configRepo repository.ConfiguracaoRepository
}

func NewConfiguracaoService(configRepo repository.ConfiguracaoRepository) ConfiguracaoService {
	return &configuracaoService{
		configRepo: configRepo,
	}
}

func (s *configuracaoService) LoadConfig() error {
	panic("not implemented") // TODO: Implementar a função LoadConfig
}

func (s *configuracaoService) GetConfig(configId int) (any, error) {
	return 1, nil // TODO: Criar busca de configurações no Reddis, em memoria
	// FIXME: Descrição do bug a ser corrigido
	// OPTIMIZE: Descrição da otimização
	// SECURITY: Descrição da melhoria de segurança
}
