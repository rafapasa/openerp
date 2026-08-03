// internal/container/container.go
package container

import (
	"sync"

	"gorm.io/gorm"
)

// Container gerencia todas as dependências da aplicação
type Container struct {
	db           *gorm.DB
	services     map[string]interface{}
	repositories map[string]interface{}
	handlers     map[string]interface{}
	mu           sync.RWMutex
}

func NewContainer(db *gorm.DB) *Container {
	return &Container{
		db:           db,
		services:     make(map[string]interface{}),
		repositories: make(map[string]interface{}),
		handlers:     make(map[string]interface{}),
	}
}

// Métodos para registrar e obter dependências
func (c *Container) RegisterService(name string, service interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services[name] = service
}

func (c *Container) GetService(name string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.services[name]
}

func (c *Container) RegisterRepository(name string, repo interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.repositories[name] = repo
}

func (c *Container) RegisterHandler(name string, handler interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[name] = handler
}

func (c *Container) GetHandler(name string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.handlers[name]
}

func (c *Container) GetDB() *gorm.DB {
	return c.db
}

// GetAllHandlers retorna todos os handlers registrados
func (c *Container) GetAllHandlers() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Retorna uma cópia para evitar modificações externas
	handlers := make(map[string]interface{})
	for k, v := range c.handlers {
		handlers[k] = v
	}
	return handlers
}
