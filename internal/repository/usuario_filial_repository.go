package repository

import (
	"context"

	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
	"gorm.io/gorm"
)

type UsuarioFilialRepository interface {
	Create(ctx context.Context, usuarioFilial *models.UsuarioFilial) error
	Update(ctx context.Context, ufi_id int, usuarioFilial *models.UsuarioFilial) error
	Delete(ctx context.Context, usuarioID int, filialID int) error
	GetByUsuarioIDAndFilialID(ctx context.Context, usuarioID int, filialID int) (*models.UsuarioFilial, error)
	FindByUsuarioID(ctx context.Context, usuarioID int) ([]models.UsuarioFilial, error)
	FindByUsuarioIDAndFilialID(ctx context.Context, usuarioID int, filialID int) (*models.UsuarioFilial, error)
	ExistByUsuarioIDAndFilialID(ctx context.Context, usuarioID int, filialID int) (bool, error)
	List(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]models.UsuarioFilial, int64, error)
	DeleteByUsuarioID(ctx context.Context, usuarioID int) error
	GetFiliaisByUsuarioID(ctx context.Context, usuarioID int) ([]int, error)
}

type usuarioFilialRepository struct {
	db *gorm.DB
}

func NewUsuarioFilialRepository(db *gorm.DB) UsuarioFilialRepository {
	return &usuarioFilialRepository{db: db}
}

func (r *usuarioFilialRepository) Create(ctx context.Context, usuarioFilial *models.UsuarioFilial) error {
	err := r.db.WithContext(ctx).Create(usuarioFilial).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *usuarioFilialRepository) Update(ctx context.Context, ufi_id int, usuarioFilial *models.UsuarioFilial) error {
	err := r.db.
		WithContext(ctx).
		Model(&models.UsuarioFilial{}).
		Omit("Usuario", "Filial").
		Where("ufi_id = ?", ufi_id).
		Updates(usuarioFilial).
		Error
	if err != nil {
		return apperrors.NewInternalError("Erro atualizando usuario filial", err)
	}
	return nil
}

func (r *usuarioFilialRepository) Delete(ctx context.Context, usuarioID int, filialID int) error {
	usuarioFilial, err := r.GetByUsuarioIDAndFilialID(ctx, usuarioID, filialID)
	if err != nil {
		return err
	}
	if usuarioFilial == nil {
		return apperrors.NewNotFoundError("usuario filial não encontrado")
	}

	usuarioFilial.SoftDelete()

	return r.Update(ctx, usuarioFilial.ID, usuarioFilial)

}

func (r *usuarioFilialRepository) GetByUsuarioIDAndFilialID(ctx context.Context, usuarioID int, filialID int) (*models.UsuarioFilial, error) {
	var usuarioFilial models.UsuarioFilial
	err := r.db.
		WithContext(ctx).
		Where("deleted_at IS NULL").
		Where("usu_id = ? AND emf_id = ?", usuarioID, filialID).
		First(&usuarioFilial).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.NewNotFoundError("usuario filial não encontrado")
		}
		return nil, apperrors.NewInternalError("Erro ao buscar usuario filial", err)
	}
	return &usuarioFilial, nil
}

func (r *usuarioFilialRepository) FindByUsuarioID(ctx context.Context, usuarioID int) ([]models.UsuarioFilial, error) {
	var usuarioFiliais []models.UsuarioFilial
	err := r.db.
		WithContext(ctx).
		Preload("Usuario").
		Preload("EmpresaFilial").
		Where("deleted_at IS NULL").
		Where("usu_id = ?", usuarioID).
		Find(&usuarioFiliais).
		Error
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao buscar usuario filiais", err)
	}
	return usuarioFiliais, nil
}

func (r *usuarioFilialRepository) FindByUsuarioIDAndFilialID(ctx context.Context, usuarioID int, filialID int) (*models.UsuarioFilial, error) {
	var usuarioFiliais models.UsuarioFilial
	err := r.db.
		WithContext(ctx).
		Preload("Usuario").
		Preload("EmpresaFilial").
		Where("deleted_at IS NULL").
		Where("usu_id = ? and emf_id = ?", usuarioID, filialID).
		First(&usuarioFiliais).
		Error
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao buscar usuario filiais", err)
	}
	return &usuarioFiliais, nil
}
func (r *usuarioFilialRepository) ExistByUsuarioIDAndFilialID(ctx context.Context, usuarioID int, filialID int) (bool, error) {
	var total int64
	err := r.db.
		Model(models.UsuarioFilial{}).
		WithContext(ctx).
		Where("deleted_at IS NULL").
		Where("usu_id = ? and emf_id = ?", usuarioID, filialID).
		Count(&total).
		Error
	if err != nil {
		return false, apperrors.NewInternalError("Erro ao buscar usuario filiais", err)
	}
	return total > 0, nil
}

func (r *usuarioFilialRepository) List(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]models.UsuarioFilial, int64, error) {
	var usuarioFiliais []models.UsuarioFilial
	var total int64

	query := r.db.
		WithContext(ctx).
		Model(&models.UsuarioFilial{}).
		Where("deleted_at IS NULL")
	query = utils.
		ApplyFilters(query, &models.UsuarioFilial{}, filters)
	err := query.
		Limit(limit).
		Offset(offset).
		Find(&usuarioFiliais).
		Error
	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro ao buscar usuario filiais", err)
	}
	total = int64(len(usuarioFiliais)) // Contar total de registros encontrados
	return usuarioFiliais, total, nil
}

func (r *usuarioFilialRepository) DeleteByUsuarioID(ctx context.Context, usuarioID int) error {
	usuarioFiliais, err := r.FindByUsuarioID(ctx, usuarioID)
	if err != nil {
		return err
	}
	if len(usuarioFiliais) == 0 {
		return nil
	}
	for _, usuarioFilial := range usuarioFiliais {
		usuarioFilial.SoftDelete()
		err = r.Update(ctx, usuarioFilial.ID, &usuarioFilial)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *usuarioFilialRepository) GetFiliaisByUsuarioID(ctx context.Context, usuarioID int) ([]int, error) {
	var filiais []int
	err := r.db.
		WithContext(ctx).
		Model(&models.UsuarioFilial{}).
		Where("usu_id = ? AND deleted_at IS NULL", usuarioID).
		Pluck("emf_id", &filiais).
		Error
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao buscar filiais do usuário.", err)
	}
	return filiais, nil
}
