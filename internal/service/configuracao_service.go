package service

import "gorm.io/gorm"

// ConfiguracaoService define os métodos públicos para o serviço de configuração.
type ConfiguracaoService interface {
	GetConfig(configId int) (any, error)
}

type configuracaoService struct {
	db *gorm.DB
}

func NewConfiguracaoService(db *gorm.DB) ConfiguracaoService {
	return &configuracaoService{
		db: db,
	}
}

func (s *configuracaoService) GetConfig(configId int) (any, error) {
	return 1, nil // TODO: Criar busca de configurações no Reddis, em memoria
	// FIXME: Descrição do bug a ser corrigido
	// OPTIMIZE: Descrição da otimização
	// SECURITY: Descrição da melhoria de segurança
}
