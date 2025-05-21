package service

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/nisibz/go-auth-tests/internal/core/domain"
	"github.com/nisibz/go-auth-tests/internal/core/port"
)

type UserService struct {
	userRepo port.UserRepository
}

func NewUserService(userRepo port.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(ctx context.Context, name, email, password string) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return user, nil
}

func (s *UserService) ListUsers(ctx context.Context, limit, offset int64) ([]*domain.User, error) {
	users, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id, name, email string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user for update: %w", err)
	}

	existingUser, err := s.userRepo.GetByEmail(context.Background(), email)
	if err == nil && existingUser.ID != user.ID {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	user.Name = name
	user.Email = email

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (s *UserService) CountUsers(ctx context.Context) (int64, error) {
	return s.userRepo.Count(ctx)
}
