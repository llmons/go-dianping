package service

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go-dianping/internal/repository"
	"go-dianping/pkg/helper/random"
	"go-dianping/pkg/helper/validator"
	"go.uber.org/zap"
)

type UserService interface {
	SendCode(ctx *gin.Context, phone string) error
}

type userService struct {
	*Service
	userRepository repository.UserRepository
}

func NewUserService(service *Service, userRepository repository.UserRepository) UserService {
	return &userService{
		Service:        service,
		userRepository: userRepository,
	}
}

func (s *userService) SendCode(ctx *gin.Context, phone string) error {
	if !validator.IsPhone(phone) {
		return errors.New("phone is invalidate")
	}

	code, err := random.Gen6DigitCode()
	if err != nil {
		return err
	}

	session := sessions.Default(ctx)
	session.Set("code", code)
	err = session.Save()
	if err != nil {
		return err
	}

	s.logger.Info("send code success", zap.String("code", code))

	return nil
}
