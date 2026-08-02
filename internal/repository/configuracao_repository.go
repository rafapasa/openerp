package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/database"
	apperrors "github.com/openerp/backend/internal/erros"
	"github.com/openerp/backend/internal/models"
)

// ============================================================
// INTERFACE
// ============================================================

type ConfiguracaoRepository interface {
	// Operações básicas
	Create(ctx context.Context, config *models.Configuracao) error
	Update(ctx context.Context, config *models.Configuracao) error
	Delete(ctx context.Context, empresaFilialID int, configID int) error

	// Buscas
	Get(ctx context.Context, empresaFilialID int, configID int) (string, error)
	GetMany(ctx context.Context, empresaFilialID int, configIDs []int) (map[int]string, error)
	GetAll(ctx context.Context, empresaFilialID int) (map[int]string, error)

	// Cache
	LoadAllToCache(ctx context.Context, empresaFilialID int) error
	InvalidateCache(ctx context.Context, empresaFilialID int, configID int) error
	InvalidateAllCache(ctx context.Context, empresaFilialID int) error
}

// ============================================================
// TYPES
// ============================================================

type configuracaoRepository struct {
	db    *gorm.DB
	redis *database.Redis
}

// ============================================================
// CONSTRUCTOR
// ============================================================

func NewConfiguracaoRepository(db *gorm.DB, redis *database.Redis) ConfiguracaoRepository {
	return &configuracaoRepository{
		db:    db,
		redis: redis,
	}
}

// ============================================================
// MÉTODOS DE CACHE
// ============================================================

// LoadAllToCache carrega todas as configurações da empresa para o Redis
func (r *configuracaoRepository) LoadAllToCache(ctx context.Context, empresaFilialID int) error {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		log.Printf("⏱️ Cache carregado para empresa %d em %s", empresaFilialID, elapsed)
	}()

	log.Printf("🔄 Carregando configurações da empresa %d para cache...", empresaFilialID)

	// Buscar apenas ID e Valor (o que realmente importa)
	var configs []struct {
		ConfigID int    `gorm:"column:config_id"`
		Valor    string `gorm:"column:config_valor"`
	}

	err := r.db.WithContext(ctx).
		Table("configuracoes").
		Select("config_id, config_valor").
		Where("emf_id = ? AND deleted_at IS NULL", empresaFilialID).
		Find(&configs).Error

	if err != nil {
		return apperrors.NewInternalError("Erro ao buscar configurações", err)
	}

	if len(configs) == 0 {
		log.Printf("⚠️ Nenhuma configuração encontrada para empresa %d", empresaFilialID)
		return nil
	}

	// Usar pipeline para operações em lote
	pipe := r.redis.Client.Pipeline()
	ctxPipe := r.redis.WithContext(ctx)

	for _, config := range configs {
		key := r.getCacheKey(empresaFilialID, config.ConfigID)
		pipe.Set(ctxPipe.Client.Context(), key, config.Valor, models.ConfigCacheTTL)
	}

	// Executar todas as operações
	_, err = pipe.Exec(ctxPipe.Client.Context())
	if err != nil {
		return apperrors.NewInternalError("Erro ao carregar configurações no cache", err)
	}

	log.Printf("✅ %d configurações carregadas no cache para empresa %d", len(configs), empresaFilialID)
	return nil
}

// InvalidateCache invalida o cache de uma configuração específica
func (r *configuracaoRepository) InvalidateCache(ctx context.Context, empresaFilialID int, configID int) error {
	key := r.getCacheKey(empresaFilialID, configID)
	return r.redis.WithContext(ctx).Delete(key)
}

// InvalidateAllCache invalida todo o cache de uma empresa
func (r *configuracaoRepository) InvalidateAllCache(ctx context.Context, empresaFilialID int) error {
	pattern := r.getCachePattern(empresaFilialID)

	// Buscar todas as chaves com o padrão
	keys, err := r.redis.WithContext(ctx).GetClient().Keys(ctx, pattern).Result()
	if err != nil {
		return apperrors.NewInternalError("Erro ao buscar chaves para invalidar", err)
	}

	if len(keys) == 0 {
		return nil
	}

	// Deletar todas as chaves
	return r.redis.WithContext(ctx).Delete(keys...)
}

// ============================================================
// MÉTODOS CRUD
// ============================================================

// Create cria uma nova configuração
func (r *configuracaoRepository) Create(ctx context.Context, config *models.Configuracao) error {
	// Verificar se já existe
	exists, err := r.exists(ctx, config.EmpresaFilialID, config.ConfigID)
	if err != nil {
		return err
	}
	if exists {
		return apperrors.NewBadRequestError(
			fmt.Sprintf("Configuração %d já existe para empresa %d", config.ConfigID, config.EmpresaFilialID),
		)
	}

	// Usar transação
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Salvar no banco
		if err := tx.Create(config).Error; err != nil {
			return apperrors.NewInternalError("Erro ao criar configuração", err)
		}

		// Salvar no cache
		key := r.getCacheKey(config.EmpresaFilialID, config.ConfigID)
		if err := r.redis.WithContext(ctx).Set(key, config.Valor, models.ConfigCacheTTL); err != nil {
			log.Printf("⚠️ Erro ao salvar no cache: %v", err)
		}

		return nil
	})
}

// Update atualiza uma configuração existente
func (r *configuracaoRepository) Update(ctx context.Context, config *models.Configuracao) error {
	// Verificar se existe
	exists, err := r.exists(ctx, config.EmpresaFilialID, config.ConfigID)
	if err != nil {
		return err
	}
	if !exists {
		return apperrors.NewNotFoundError(
			fmt.Sprintf("Configuração %d não encontrada para empresa %d", config.ConfigID, config.EmpresaFilialID),
		)
	}

	// Usar transação
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Atualizar no banco
		result := tx.
			Model(&models.Configuracao{}).
			Where("emf_id = ? AND config_id = ? AND deleted_at IS NULL",
				config.EmpresaFilialID, config.ConfigID).
			Update("config_valor", config.Valor)

		if result.Error != nil {
			return apperrors.NewInternalError("Erro ao atualizar configuração", result.Error)
		}

		if result.RowsAffected == 0 {
			return apperrors.NewNotFoundError("Configuração não encontrada para atualização")
		}

		// Atualizar no cache
		key := r.getCacheKey(config.EmpresaFilialID, config.ConfigID)
		if err := r.redis.WithContext(ctx).Set(key, config.Valor, models.ConfigCacheTTL); err != nil {
			log.Printf("⚠️ Erro ao atualizar cache: %v", err)
		}

		return nil
	})
}

// Delete realiza exclusão lógica
func (r *configuracaoRepository) Delete(ctx context.Context, empresaFilialID int, configID int) error {
	// Verificar se existe
	exists, err := r.exists(ctx, empresaFilialID, configID)
	if err != nil {
		return err
	}
	if !exists {
		return apperrors.NewNotFoundError(
			fmt.Sprintf("Configuração %d não encontrada para empresa %d", configID, empresaFilialID),
		)
	}

	// Usar transação
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Soft delete
		result := tx.
			Model(&models.Configuracao{}).
			Where("emf_id = ? AND config_id = ?", empresaFilialID, configID).
			Update("deleted_at", time.Now())

		if result.Error != nil {
			return apperrors.NewInternalError("Erro ao deletar configuração", result.Error)
		}

		if result.RowsAffected == 0 {
			return apperrors.NewNotFoundError("Configuração não encontrada para deletar")
		}

		// Remover do cache
		key := r.getCacheKey(empresaFilialID, configID)
		if err := r.redis.WithContext(ctx).Delete(key); err != nil {
			log.Printf("⚠️ Erro ao remover do cache: %v", err)
		}

		return nil
	})
}

// ============================================================
// MÉTODOS DE LEITURA (OTIMIZADOS PARA CACHE)
// ============================================================

// Get busca o valor de uma configuração específica
// Primeiro tenta o cache, se falhar busca no banco e atualiza o cache
func (r *configuracaoRepository) Get(ctx context.Context, empresaFilialID int, configID int) (string, error) {
	// 1. Tentar cache
	key := r.getCacheKey(empresaFilialID, configID)
	valor, err := r.redis.WithContext(ctx).Get(key)

	if err == nil && valor != "" {
		return valor, nil
	}

	// 2. Buscar no banco
	var config models.Configuracao
	err = r.db.WithContext(ctx).
		Select("config_valor").
		Where("emf_id = ? AND config_id = ? AND deleted_at IS NULL", empresaFilialID, configID).
		First(&config).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperrors.NewNotFoundError(
				fmt.Sprintf("Configuração %d não encontrada para empresa %d", configID, empresaFilialID),
			)
		}
		return "", apperrors.NewInternalError("Erro ao buscar configuração", err)
	}

	// 3. Atualizar cache (em background)
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := r.redis.WithContext(cacheCtx).Set(key, config.Valor, models.ConfigCacheTTL); err != nil {
			log.Printf("⚠️ Erro ao atualizar cache: %v", err)
		}
	}()

	return config.Valor, nil
}

// GetMany busca múltiplas configurações de uma vez
// Retorna um map[configID]valor
func (r *configuracaoRepository) GetMany(ctx context.Context, empresaFilialID int, configIDs []int) (map[int]string, error) {
	if len(configIDs) == 0 {
		return make(map[int]string), nil
	}

	result := make(map[int]string)
	missingIDs := []int{}

	// 1. Buscar do cache
	for _, configID := range configIDs {
		key := r.getCacheKey(empresaFilialID, configID)
		valor, err := r.redis.WithContext(ctx).Get(key)
		if err == nil && valor != "" {
			result[configID] = valor
		} else {
			missingIDs = append(missingIDs, configID)
		}
	}

	// Se todos foram encontrados no cache
	if len(missingIDs) == 0 {
		return result, nil
	}

	// 2. Buscar os que faltam no banco
	var configs []models.Configuracao
	err := r.db.WithContext(ctx).
		Select("config_id, config_valor").
		Where("emf_id = ? AND config_id IN ? AND deleted_at IS NULL", empresaFilialID, missingIDs).
		Find(&configs).Error

	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao buscar configurações", err)
	}

	// 3. Atualizar cache e resultado
	pipe := r.redis.Client.Pipeline()
	ctxPipe := r.redis.WithContext(ctx)

	for _, config := range configs {
		result[config.ConfigID] = config.Valor
		key := r.getCacheKey(empresaFilialID, config.ConfigID)
		pipe.Set(ctxPipe.Client.Context(), key, config.Valor, models.ConfigCacheTTL)
	}

	// Executar pipeline em background
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := pipe.Exec(cacheCtx); err != nil {
			log.Printf("⚠️ Erro ao atualizar cache em lote: %v", err)
		}
	}()

	// Verificar se algum ID não foi encontrado
	for _, id := range missingIDs {
		if _, found := result[id]; !found {
			log.Printf("⚠️ Configuração %d não encontrada para empresa %d", id, empresaFilialID)
		}
	}

	return result, nil
}

// GetAll busca todas as configurações de uma empresa
// Retorna um map[configID]valor
func (r *configuracaoRepository) GetAll(ctx context.Context, empresaFilialID int) (map[int]string, error) {
	result := make(map[int]string)

	// 1. Buscar do cache (padrão de todas as configs da empresa)
	pattern := r.getCachePattern(empresaFilialID)
	keys, err := r.redis.WithContext(ctx).GetClient().Keys(ctx, pattern).Result()

	if err == nil && len(keys) > 0 {
		// Buscar valores em lote usando o método do Redis diretamente
		values, err := r.redis.WithContext(ctx).GetClient().MGet(ctx, keys...).Result()
		if err == nil {
			for i, key := range keys {
				configID := r.extractConfigID(key)
				if configID > 0 && i < len(values) && values[i] != nil {
					if valor, ok := values[i].(string); ok {
						result[configID] = valor
					}
				}
			}
			// Se encontrou todas, retorna
			if len(result) > 0 {
				return result, nil
			}
		}
	}

	// 2. Buscar do banco
	var configs []models.Configuracao
	err = r.db.WithContext(ctx).
		Select("config_id, config_valor").
		Where("emf_id = ? AND deleted_at IS NULL", empresaFilialID).
		Find(&configs).Error

	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao buscar todas as configurações", err)
	}

	// 3. Carregar no cache (em background)
	go func() {
		if err := r.LoadAllToCache(context.Background(), empresaFilialID); err != nil {
			log.Printf("⚠️ Erro ao carregar cache em background: %v", err)
		}
	}()

	// 4. Preencher resultado
	for _, config := range configs {
		result[config.ConfigID] = config.Valor
	}

	return result, nil
}

// ============================================================
// MÉTODOS AUXILIARES
// ============================================================

// exists verifica se uma configuração existe no banco
func (r *configuracaoRepository) exists(ctx context.Context, empresaFilialID int, configID int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Configuracao{}).
		Where("emf_id = ? AND config_id = ? AND deleted_at IS NULL", empresaFilialID, configID).
		Count(&count).Error

	if err != nil {
		return false, apperrors.NewInternalError("Erro ao verificar existência da configuração", err)
	}

	return count > 0, nil
}

// getCacheKey gera a chave do cache para uma configuração específica
func (r *configuracaoRepository) getCacheKey(empresaFilialID int, configID int) string {
	return fmt.Sprintf("config:%d:%d", empresaFilialID, configID)
}

// getCachePattern gera o padrão para buscar todas as chaves de uma empresa
func (r *configuracaoRepository) getCachePattern(empresaFilialID int) string {
	return fmt.Sprintf("config:%d:*", empresaFilialID)
}

// extractConfigID extrai o config_id de uma chave do cache
func (r *configuracaoRepository) extractConfigID(key string) int {
	var configID int
	// Formato: config:{empresaID}:{configID}
	if n, err := fmt.Sscanf(key, "config:%*d:%d", &configID); err == nil && n == 1 {
		return configID
	}
	return 0
}
