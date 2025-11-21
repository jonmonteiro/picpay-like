package user

import (
	"fmt"

	"github.com/jmonteiro/picpay-like/core/types"
	"github.com/jmonteiro/picpay-like/core/utils"
)

type UserService struct {
	store types.UserStore
}

func NewUserService(store types.UserStore) *UserService {
	return &UserService{
		store: store,
	}
}

// RegisterUser contém a lógica de negócio para registrar um novo usuário
func (s *UserService) RegisterUser(payload types.RegisterUserPayload) error {
	// Verifica se o usuário já existe
	_, err := s.store.GetUserByEmail(payload.Email)
	if err == nil {
		return fmt.Errorf("user with email %s already exists", payload.Email)
	}

	// Hash da senha
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Cria o usuário
	user := types.User{
		Email:    payload.Email,
		Password: hashedPassword,
	}

	err = s.store.CreateUser(user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
