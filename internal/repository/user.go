package repository

import (
	"github.com/pkg/errors"
	"go-dianping/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByPhone(phone string) (*model.User, error)
	CreateUser(*model.User) error
}
type userRepository struct {
	*Repository
}

func NewUserRepository(repository *Repository) UserRepository {
	return &userRepository{
		Repository: repository,
	}
}

func (r *userRepository) GetUserByPhone(phone string) (*model.User, error) {
	var user model.User
	err := r.db.Where(&model.User{Phone: phone}).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *model.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}
