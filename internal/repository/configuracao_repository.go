package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/database"
	"github.com/openerp/backend/internal/models"
)

type ConfiguracaoRepository interface {
	Create(ctx context.Context, config *models.Configuracao) error
	FindByID(ctx context.Context, empresaFilialID int, configID int) (*models.Configuracao, error)
	FindByNome(ctx context.Context, empresaFilialID int, nome string) (*models.Configuracao, error)
	FindByModulo(ctx context.Context, empresaFilialID int, moduloID int) ([]models.Configuracao, error)
	FindAll(ctx context.Context, empresaFilialID int) ([]models.Configuracao, error)
	Update(ctx context.Context, config *models.Configuracao) error
	Delete(ctx context.Context, empresaFilialID int, configID int) error
	SoftDelete(ctx context.Context, empresaFilialID int, configID int) error
	LoadAllToCache(ctx context.Context, empresaFilialID int) error
}

type configuracaoRepository struct {
	db    *gorm.DB
	redis *database.Redis
}

func NewConfiguracaoRepository(db *gorm.DB, redis *database.Redis) ConfiguracaoRepository {
	return &configuracaoRepository{
		db:    db,
		redis: redis,
	}
}

func (r *configuracaoRepository) LoadAllToCache(ctx context.Context, empresaFilialID int) error {
	configs, err := r.FindAll(ctx, empresaFilialID)
	if err != nil {
		return err
	}

	for _, config := range configs {
		err := r.redis.Set(config.Key(), config.Valor, models.ConfigCacheTTL)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *configuracaoRepository) Create(ctx context.Context, config *models.Configuracao) error {
	return r.db.WithContext(ctx).Create(config).Error
}

func (r *configuracaoRepository) FindByID(ctx context.Context, empresaFilialID int, configID int) (*models.Configuracao, error) {
	var config models.Configuracao

	err := r.db.WithContext(ctx).
		Preload("EmpresaFilial").
		Where("emf_id = ? AND config_id = ? AND deleted_at IS NULL", empresaFilialID, configID).
		First(&config).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &config, nil
}

func (r *configuracaoRepository) FindByNome(ctx context.Context, empresaFilialID int, nome string) (*models.Configuracao, error) {
	var config models.Configuracao

	err := r.db.WithContext(ctx).
		Preload("Modulo").
		Where("emf_id = ? AND config_nome = ? AND deleted_at IS NULL", empresaFilialID, nome).
		First(&config).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &config, nil
}

func (r *configuracaoRepository) FindByModulo(ctx context.Context, empresaFilialID int, moduloID int) ([]models.Configuracao, error) {
	var configs []models.Configuracao

	err := r.db.WithContext(ctx).
		Where("emf_id = ? AND mod_id = ? AND deleted_at IS NULL", empresaFilialID, moduloID).
		Order("config_nome").
		Find(&configs).Error

	return configs, err
}

func (r *configuracaoRepository) FindAll(ctx context.Context, empresaFilialID int) ([]models.Configuracao, error) {
	var configs []models.Configuracao

	err := r.db.WithContext(ctx).
		Preload("EmpresaFilial").
		Where("emf_id = ? AND deleted_at IS NULL", empresaFilialID).
		Order("config_nome").
		Find(&configs).Error

	return configs, err
}

func (r *configuracaoRepository) Update(ctx context.Context, config *models.Configuracao) error {
	return r.db.WithContext(ctx).Save(config).Error
}

func (r *configuracaoRepository) Delete(ctx context.Context, empresaFilialID int, configID int) error {
	return r.db.WithContext(ctx).
		Where("emf_id = ? AND config_id = ?", empresaFilialID, configID).
		Delete(&models.Configuracao{}).Error
}

func (r *configuracaoRepository) SoftDelete(ctx context.Context, empresaFilialID int, configID int) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.Configuracao{}).
		Where("emf_id = ? AND config_id = ?", empresaFilialID, configID).
		Update("deleted_at", now).Error
}
