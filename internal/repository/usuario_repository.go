package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/utils"
)

type UsuarioRepository interface {
	Create(ctx context.Context, usuario *models.Usuario) error
	Update(ctx context.Context, id int, usuario *models.Usuario) error
	Delete(ctx context.Context, id int) error
	FindByID(ctx context.Context, id int) (*models.Usuario, error)
	FindByLogin(ctx context.Context, login string) (*models.Usuario, error)
	FindByLoginWithGrupo(ctx context.Context, login string) (*models.Usuario, error)
	FindGrupoByID(ctx context.Context, id int) (*models.GrupoUsuario, error)
	FindUsuarioFiliais(ctx context.Context, usuarioID int) ([]models.UsuarioFilial, error)
	List(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]models.Usuario, int64, error)
}

type usuarioRepository struct {
	db *gorm.DB
}

func NewUsuarioRepository(db *gorm.DB) UsuarioRepository {
	return &usuarioRepository{db: db}
}

func (r *usuarioRepository) Create(ctx context.Context, usuario *models.Usuario) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Inserir usuário
		if err := tx.Create(usuario).Error; err != nil {
			return apperrors.NewInternalError("Erro ao criar usuário", err)
		}

		// 2. Inserir filiais (se houver)
		if len(usuario.UsuarioFiliais) > 0 {
			// Garantir que todas as filiais tenham o ID do usuário
			for i := range usuario.UsuarioFiliais {
				usuario.UsuarioFiliais[i].UsuarioID = usuario.ID
			}

			if err := tx.Create(&usuario.UsuarioFiliais).Error; err != nil {
				return apperrors.NewInternalError("Erro ao criar filiais do usuário", err)
			}
		}

		return nil
	})
}

func (r *usuarioRepository) Update(ctx context.Context, id int, usuario *models.Usuario) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // 1. Atualizar usuário
        if err := tx.Omit("UsuarioFiliais").Save(usuario).Error; err != nil {
            return err
        }

        // 2. Substituir todas as filiais
        if err := tx.Model(usuario).Association("UsuarioFiliais").Replace(usuario.UsuarioFiliais); err != nil {
            return err
        }

        return nil
    })
}


func (r *usuarioRepository) Delete(ctx context.Context, id int) error {
	usu, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if usu == nil {
		return apperrors.NewNotFoundError("usuário não encontrado")
	}
	usu.SoftDelete()
	return r.Update(ctx, usu.ID, usu)
}

// FindByID busca um usuário pelo ID - SEM PRELOAD para evitar erros
func (r *usuarioRepository) FindByID(ctx context.Context, id int) (*models.Usuario, error) {
	var usuario models.Usuario
	err := r.db.
		Where("usu_id = ? AND deleted_at IS NULL", id).
		First(&usuario).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("usuário não encontrado")
		}
		return nil, apperrors.NewInternalError("Erro ao buscar usuário.", err)
	}
	return &usuario, nil
}

// FindByLogin busca um usuário pelo login - SEM PRELOAD para evitar erros
func (r *usuarioRepository) FindByLogin(ctx context.Context, login string) (*models.Usuario, error) {
	var usuario models.Usuario
	err := r.db.
		Where("usu_login = ? AND deleted_at IS NULL", login).
		First(&usuario).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("usuário não encontrado")
		}
		return nil, apperrors.NewInternalError("Erro ao buscar usuário.", err)
	}
	return &usuario, nil
}

// FindByLoginWithGrupo busca um usuário com o grupo (apenas grupo)
func (r *usuarioRepository) FindByLoginWithGrupo(ctx context.Context, login string) (*models.Usuario, error) {
	var usuario models.Usuario
	err := r.db.
		Preload("GrupoUsuario").
		Where("usu_login = ? AND deleted_at IS NULL", login).
		First(&usuario).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("usuário não encontrado")
		}
		return nil, apperrors.NewInternalError("Erro ao buscar usuário.", err)
	}
	return &usuario, nil
}

// FindByIDWithGrupo busca um usuário com o grupo (apenas grupo)
func (r *usuarioRepository) FindByIDWithGrupo(ctx context.Context, id int) (*models.Usuario, error) {
	var usuario models.Usuario
	err := r.db.
		Preload("GrupoUsuario").
		Where("usu_id = ? AND deleted_at IS NULL", id).
		First(&usuario).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("usuário não encontrado")
		}
		return nil, apperrors.NewInternalError("Erro ao buscar usuário.", err)
	}
	return &usuario, nil
}

// FindGrupoByID busca um grupo de usuário pelo ID
func (r *usuarioRepository) FindGrupoByID(ctx context.Context, id int) (*models.GrupoUsuario, error) {
	var grupo models.GrupoUsuario
	err := r.db.
		Where("gpu_id = ? AND deleted_at IS NULL", id).
		First(&grupo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("grupo não encontrado")
		}
		return nil, apperrors.NewInternalError("Erro ao buscar grupo.", err)
	}
	return &grupo, nil
}

// FindUsuarioFiliais busca as filiais do usuário
func (r *usuarioRepository) FindUsuarioFiliais(ctx context.Context, usuarioID int) ([]models.UsuarioFilial, error) {
	var filiais []models.UsuarioFilial
	err := r.db.
		WithContext(ctx).
		Where("usu_id = ? AND deleted_at IS NULL", usuarioID).
		Find(&filiais).Error
	if err != nil {
		return nil, apperrors.NewInternalError("Erro ao buscar filiais do usuário.", err)
	}
	return filiais, nil
}

func (r *usuarioRepository) List(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]models.Usuario, int64, error) {
	var usuarios []models.Usuario
	var total int64

	query := r.db.
		WithContext(ctx).
		Model(&models.Usuario{}).
		Where("deleted_at IS NULL")
	query = utils.
		ApplyFilters(query, &models.Usuario{}, filters)
	// Aplicar paginação e buscar registros
	err := query.
		Limit(limit).
		Offset(offset).
		Find(&usuarios).
		Error
	if err != nil {
		return nil, 0, apperrors.NewInternalError("Erro ao buscar usuários.", err)
	}
	total = int64(len(usuarios)) // Contar total de registros encontrados
	return usuarios, total, nil
}
