package service

import (
	"fmt"

	"github.com/jmonteiro/picpay-like/core/domain/auth"
	domain "github.com/jmonteiro/picpay-like/core/domain/user"
)

type UserService struct {
	store domain.UserStore
}

func NewUserService(store domain.UserStore) *UserService {
	return &UserService{
		store: store,
	}
}

// RegisterUser contém a lógica de negócio para registrar um novo usuário
func (s *UserService) RegisterUser(payload domain.RegisterUserPayload) error {
	// Verifica se o usuário já existe
	_, err := s.store.GetUserByEmail(payload.Email)
	if err == nil {
		return fmt.Errorf("user with email %s already exists", payload.Email)
	}

	// Hash da senha
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Cria o usuário
	user := domain.User{
		Email:    payload.Email,
		Password: hashedPassword,
	}

	err = s.store.CreateUser(user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
