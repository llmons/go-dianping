package repository

import (
	"context"
	"go-dianping/internal/entity"
)

type UserRepository interface {
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
}
type userRepository struct {
	*Repository
}

func NewUserRepository(repository *Repository) UserRepository {
	return &userRepository{
		Repository: repository,
	}
}

func (r *userRepository) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
	return r.query.WithContext(ctx).User.Where(r.query.User.Phone.Eq(phone)).First()
}

func (r *userRepository) Save(ctx context.Context, user *entity.User) error {
	return r.query.WithContext(ctx).User.Save(user)
}
