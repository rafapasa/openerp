package service

import (
	"context"

	"github.com/openerp/backend/internal/apperrors"
	"github.com/openerp/backend/internal/models"
	"github.com/openerp/backend/internal/repository"
)

type UsuarioFilialService interface {
	AddFilial(ctx context.Context, usuarioID int, empresaFilialID int) error
	RemoveFilial(ctx context.Context, usuarioID int, empresaFilialID int) error
	RemoveAllFiliais(ctx context.Context, usuarioID int) error
	GetFiliaisByUsuarioID(ctx context.Context, usuarioID int) ([]int, error)
	FindByUsuarioID(ctx context.Context, usuarioID int) ([]models.UsuarioFilial, error)
	HasFilial(ctx context.Context, usuarioID int, empresaFilialID int) (bool, error)
}

type usuarioFilialService struct {
	usuarioFilialRepo repository.UsuarioFilialRepository
}

func NewUsuarioFilialService(usuarioFilialRepo repository.UsuarioFilialRepository) UsuarioFilialService {
	return &usuarioFilialService{
		usuarioFilialRepo: usuarioFilialRepo,
	}
}

func (s *usuarioFilialService) AddFilial(ctx context.Context, usuarioID int, empresaFilialID int) error {
	if existing, err := s.usuarioFilialRepo.ExistByUsuarioIDAndFilialID(ctx, usuarioID, empresaFilialID); err != nil {
		return err
	} else if existing {
		return apperrors.NewConflictError("usuário já possui essa filial associada")
	}
	err := s.usuarioFilialRepo.Create(ctx, &models.UsuarioFilial{
		UsuarioID:       usuarioID,
		EmpresaFilialID: empresaFilialID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *usuarioFilialService) RemoveFilial(ctx context.Context, usuarioID int, empresaFilialID int) error {
	if existing, err := s.usuarioFilialRepo.ExistByUsuarioIDAndFilialID(ctx, usuarioID, empresaFilialID); err != nil {
		return err
	} else if !existing {
		return apperrors.NewNotFoundError("usuário não possui essa filial associada")
	}
	err := s.usuarioFilialRepo.Delete(ctx, usuarioID, empresaFilialID)
	if err != nil {
		return err
	}
	return nil
}

func (s *usuarioFilialService) RemoveAllFiliais(ctx context.Context, usuarioID int) error {
	if usuarioID <= 0 {
		return apperrors.NewValidationError("ID do usuário inválido")
	}
	return s.usuarioFilialRepo.DeleteByUsuarioID(ctx, usuarioID)
}

func (s *usuarioFilialService) GetFiliaisByUsuarioID(ctx context.Context, usuarioID int) ([]int, error) {
	if usuarioID <= 0 {
		return nil, apperrors.NewValidationError("ID do usuário inválido")
	}
	return s.usuarioFilialRepo.GetFiliaisByUsuarioID(ctx, usuarioID)
}

func (s *usuarioFilialService) HasFilial(ctx context.Context, usuarioID int, empresaFilialID int) (bool, error) {
	if usuarioID <= 0 {
		return false, apperrors.NewValidationError("ID do usuário inválido")
	}
	if empresaFilialID <= 0 {
		return false, apperrors.NewValidationError("ID da filial inválido")
	}
	existing, err := s.usuarioFilialRepo.ExistByUsuarioIDAndFilialID(ctx, usuarioID, empresaFilialID)
	if err != nil {
		return false, err
	}
	return existing, nil

}

func (s *usuarioFilialService) FindByUsuarioID(ctx context.Context, usuarioID int) ([]models.UsuarioFilial, error) {
	if usuarioID <= 0 {
		return nil, apperrors.NewValidationError("ID do usuário inválido")
	}
	return s.usuarioFilialRepo.FindByUsuarioID(ctx, usuarioID)
}
