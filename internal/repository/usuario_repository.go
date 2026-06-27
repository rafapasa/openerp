package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/openerp/backend/internal/models"
)

type UsuarioRepository struct {
	db *gorm.DB
}

func NewUsuarioRepository(db *gorm.DB) *UsuarioRepository {
	return &UsuarioRepository{db: db}
}

// FindByID busca um usuário pelo ID - SEM PRELOAD para evitar erros
func (r *UsuarioRepository) FindByID(id int) (*models.Usuario, error) {
	var usuario models.Usuario
	err := r.db.
		Where("usu_id = ? AND deleted_at IS NULL", id).
		First(&usuario).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}
	return &usuario, nil
}

// FindByLogin busca um usuário pelo login - SEM PRELOAD para evitar erros
func (r *UsuarioRepository) FindByLogin(login string) (*models.Usuario, error) {
	var usuario models.Usuario
	err := r.db.
		Where("usu_login = ? AND deleted_at IS NULL", login).
		First(&usuario).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}
	return &usuario, nil
}

// FindByLoginWithGrupo busca um usuário com o grupo (apenas grupo)
func (r *UsuarioRepository) FindByLoginWithGrupo(login string) (*models.Usuario, error) {
	var usuario models.Usuario
	err := r.db.
		Preload("GrupoUsuario").
		Where("usu_login = ? AND deleted_at IS NULL", login).
		First(&usuario).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}
	return &usuario, nil
}

// FindByIDWithGrupo busca um usuário com o grupo (apenas grupo)
func (r *UsuarioRepository) FindByIDWithGrupo(id int) (*models.Usuario, error) {
	var usuario models.Usuario
	err := r.db.
		Preload("GrupoUsuario").
		Where("usu_id = ? AND deleted_at IS NULL", id).
		First(&usuario).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}
	return &usuario, nil
}

// FindGrupoByID busca um grupo de usuário pelo ID
func (r *UsuarioRepository) FindGrupoByID(id int) (*models.GrupoUsuario, error) {
	var grupo models.GrupoUsuario
	err := r.db.
		Where("gpu_id = ? AND deleted_at IS NULL", id).
		First(&grupo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("grupo não encontrado")
		}
		return nil, err
	}
	return &grupo, nil
}

// FindUsuarioFiliais busca as filiais do usuário
func (r *UsuarioRepository) FindUsuarioFiliais(usuarioID int) ([]models.UsuarioFilial, error) {
	var filiais []models.UsuarioFilial
	err := r.db.
		Where("usu_id = ? AND deleted_at IS NULL", usuarioID).
		Find(&filiais).Error
	if err != nil {
		return nil, err
	}
	return filiais, nil
}
