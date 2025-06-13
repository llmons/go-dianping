package service

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go-dianping/internal/model"
	"go-dianping/internal/repository"
	"go-dianping/pkg/constants"
	"go-dianping/pkg/helper/random"
	"go-dianping/pkg/helper/validator"
	"go.uber.org/zap"
)

type UserService interface {
	SendCode(ctx *gin.Context, phone string) error
	Login(ctx *gin.Context, phone, code, password string) error
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

	code := random.Number(6)

	session := sessions.Default(ctx)
	session.Set("code", code)
	err := session.Save()
	if err != nil {
		return err
	}

	s.logger.Info("send code success", zap.String("code", code))

	return nil
}

func (s *userService) Login(ctx *gin.Context, phone, code, _ string) error {
	if !validator.IsPhone(phone) {
		return errors.New("phone is invalidate")
	}

	session := sessions.Default(ctx)
	cacheCode := session.Get("code")
	if code != cacheCode {
		return errors.New("code is invalidate")
	}

	user, err := s.userRepository.GetUserByPhone(phone)
	if err != nil {
		return err
	}
	if user == nil {
		user, err = s.createUserWithPhone(phone)
		if err != nil {
			return err
		}
	}

	session.Set("userID", user.Id)
	err = session.Save()
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) createUserWithPhone(phone string) (*model.User, error) {
	var user model.User
	user.Phone = phone
	user.NickName = fmt.Sprintf("%s%s", constants.UserNickNamePrefix, random.String(10))
	err := s.userRepository.CreateUser(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
