package service

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go-dianping/internal/dto"
	"go-dianping/internal/model"
	"go-dianping/internal/repository"
	"go-dianping/pkg/constants"
	"go-dianping/pkg/helper/random"
	"go-dianping/pkg/helper/validator"
	"go.uber.org/zap"
)

type UserService interface {
	SendCode(ctx *gin.Context, phone string) error
	Login(ctx *gin.Context, form *dto.LoginForm) error
	Me(ctx *gin.Context) (*dto.User, error)
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

func (s *userService) Login(ctx *gin.Context, form *dto.LoginForm) error {
	if !validator.IsPhone(form.Phone) {
		return errors.New("phone is invalidate")
	}

	session := sessions.Default(ctx)
	cacheCode := session.Get("code")
	if form.Code != cacheCode {
		return errors.New("code is invalidate")
	}

	user, err := s.userRepository.GetUserByPhone(form.Phone)
	if err != nil {
		return err
	}
	if user == nil {
		user, err = s.createUserWithPhone(form.Phone)
		if err != nil {
			return err
		}
	}

	var userDto = dto.User{
		Id:       user.Id,
		NickName: user.NickName,
		Icon:     user.Icon,
	}
	session.Set("user", userDto)
	err = session.Save()
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) Me(ctx *gin.Context) (*dto.User, error) {
	session := sessions.Default(ctx)
	val := session.Get("user")
	user, ok := val.(*dto.User)
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
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
