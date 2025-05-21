package port

import (
	"context"

	"github.com/nisibz/go-auth-tests/internal/core/domain"
)

type UserService interface {
	CreateUser(ctx context.Context, name, email, password string) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	ListUsers(ctx context.Context, limit, offset int64) ([]*domain.User, error)
	UpdateUser(ctx context.Context, id, name, email string) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
	CountUsers(ctx context.Context) (int64, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int64) ([]*domain.User, error)
	Count(ctx context.Context) (int64, error)
}
