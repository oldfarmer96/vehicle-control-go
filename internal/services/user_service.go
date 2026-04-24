// Package services logica denecogcio
package services

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/oldfarmer96/vehicle-control-go/internal/models"
	"github.com/oldfarmer96/vehicle-control-go/internal/store"
)

type UserService struct {
	userStore *store.UserStore
}

func NewUserService(store *store.UserStore) *UserService {
	return &UserService{userStore: store}
}

func (s *UserService) CreateUser(ctx context.Context, payload models.CreateUserDTO) (*models.User, error) {
	_, err := s.userStore.FindByUsername(ctx, payload.Username)
	if err == nil {
		return nil, errors.New("el username ya está registrado")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error al procesar la contraseña")
	}

	rol := payload.Rol
	if rol == "" {
		rol = "CONSULTOR"
	}

	return s.userStore.Create(ctx, payload, string(hashedPassword))
}

func (s *UserService) ListUsers(ctx context.Context, page, limit int, search string) (*models.UserListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	users, total, err := s.userStore.List(ctx, page, limit, search)
	if err != nil {
		return nil, err
	}

	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	return &models.UserListResponse{
		Users:      users,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, payload models.UpdateUserDTO) (*models.User, error) {
	if payload.Password != nil && len(*payload.Password) > 0 {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*payload.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("error al procesar la contraseña")
		}
		hashedStr := string(hashed)
		payload.Password = &hashedStr
	}
	return s.userStore.Update(ctx, id, payload)
}

func (s *UserService) ToggleUserActive(ctx context.Context, id string) (*models.User, error) {
	return s.userStore.ToggleActive(ctx, id)
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.userStore.FindByID(ctx, id)
}
