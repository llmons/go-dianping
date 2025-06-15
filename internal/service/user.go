package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go-dianping/internal/model"
	"go-dianping/internal/pkg/constants"
	"go-dianping/internal/pkg/dto"
	"go-dianping/internal/repository"
	"go-dianping/pkg/helper/random"
	"go-dianping/pkg/helper/uuid"
	"go-dianping/pkg/helper/validator"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type UserService interface {
	SendCode(ctx *gin.Context, phone string) error
	Login(ctx *gin.Context, form *dto.LoginForm) (string, error)
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

	err := s.rdb.Set(ctx, constants.RedisLoginCodeKey+phone, code, time.Minute*constants.RedisLoginCodeTTL).Err()
	if err != nil {
		return err
	}

	s.logger.Info("send code success", zap.String("code", code))

	return nil
}

func (s *userService) Login(ctx *gin.Context, form *dto.LoginForm) (string, error) {
	if !validator.IsPhone(form.Phone) {
		return "", errors.New("phone is invalidate")
	}

	cacheCode, err := s.rdb.Get(ctx, constants.RedisLoginCodeKey+form.Phone).Result()
	if err != nil {
		return "", err
	}
	if form.Code != cacheCode {
		return "", errors.New("code is invalidate")
	}

	user, err := s.userRepository.GetUserByPhone(form.Phone)
	if err != nil {
		return "", err
	}
	if user == nil {
		user, err = s.createUserWithPhone(form.Phone)
		if err != nil {
			return "", err
		}
	}

	token := uuid.GenUUID()
	err = s.rdb.HSet(ctx, constants.RedisLoginUserKey+token, map[string]string{
		"id":       strconv.Itoa(int(user.Id)),
		"nickname": user.NickName,
		"icon":     user.Icon,
	}).Err()
	if err != nil {
		return "", err
	}

	err = s.rdb.Expire(ctx, constants.RedisLoginUserKey+token, time.Minute*constants.RedisLoginUserTTL).Err()
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *userService) Me(ctx *gin.Context) (*dto.User, error) {
	result, err := s.rdb.HGetAll(ctx, constants.RedisLoginUserKey).Result()
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("user not found")
	}
	var user dto.User
	user.NickName = result["nickname"]
	user.Icon = result["icon"]
	return &user, nil
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
