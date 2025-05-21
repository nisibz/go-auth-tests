package service

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/nisibz/go-auth-tests/internal/core/domain"
	"github.com/nisibz/go-auth-tests/internal/core/port"
	"github.com/nisibz/go-auth-tests/internal/core/util"
)

type AuthService struct {
	userRepo port.UserRepository
}

func NewAuthService(userRepo port.UserRepository) port.AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(name, email, password string) (string, error) {
	_, err := s.userRepo.GetByEmail(context.Background(), email)
	if err == nil {
		return "", fmt.Errorf("user with email %s already exists", email)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	err = s.userRepo.Create(context.Background(), user)
	if err != nil {
		return "", fmt.Errorf("failed to register user: %w", err)
	}

	if user.ID.IsZero() {
		return "", fmt.Errorf("failed to retrieve user ID after creation")
	}

	token, err := util.GenerateToken(user.ID.Hex())
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(context.Background(), email)
	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("login failed: invalid credentials")
	}

	token, err := util.GenerateToken(user.ID.Hex())
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
