package user

import (
	"fmt"

	"github.com/jmonteiro/picpay-like/core/config"
	"github.com/jmonteiro/picpay-like/core/domain/auth"
	"github.com/jmonteiro/picpay-like/core/types"
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
	hashedPassword, err := auth.HashPassword(payload.Password)
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

// LoginUser contém a lógica de negócio para autenticar um usuário
func (s *UserService) LoginUser(payload types.LoginUserPayload) (string, error) {
	// Busca o usuário pelo email
	u, err := s.store.GetUserByEmail(payload.Email)
	if err != nil {
		return "", fmt.Errorf("invalid email or password")
	}

	// Compara as senhas
	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		return "", fmt.Errorf("invalid email or password")
	}

	// Cria o token JWT
	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, int(u.ID))
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	return token, nil
}
