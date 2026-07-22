package service

import "gorm.io/gorm"

type ConfiguracaoService struct {
	db *gorm.DB
}

func NewConfiguracaoService(db *gorm.DB) *ConfiguracaoService {
	return &ConfiguracaoService{
		db: db,
	}
}

func (s *ConfiguracaoService) GetConfig(configId int) (any, error) {
	return 1, nil // TODO: Criar busca de configurações no Reddis, em memoria
	// FIXME: Descrição do bug a ser corrigido
	// OPTIMIZE: Descrição da otimização
	// SECURITY: Descrição da melhoria de segurança
}
