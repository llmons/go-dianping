package service

import (
	"context"
	"github.com/duke-git/lancet/v2/random"
	"github.com/pkg/errors"
	"go-dianping/api/v1"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/base/user_holder"
	"go-dianping/internal/model"
	"go-dianping/internal/repository"
	"go-dianping/pkg/helper/validator"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type UserService interface {
	SendCode(ctx context.Context, req *v1.SendCodeReq) error
	Login(ctx context.Context, req *v1.LoginReq) (*v1.LoginRespData, error)
	GetMe(ctx context.Context) (*v1.GetMeRespData, error)
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

func (s *userService) SendCode(ctx context.Context, req *v1.SendCodeReq) error {
	if !validator.IsPhone(req.Phone) {
		return errors.New("phone is invalidate")
	}

	var code string
	if s.conf.Get("env") == "prod" {
		code = random.RandNumeral(6)
	} else {
		code = "123456"
	}

	key, ttl := constants.RedisLoginCodeKey+req.Phone, time.Minute*constants.RedisLoginCodeTTL
	err := s.rdb.Set(ctx, key, code, ttl).Err()
	if err != nil {
		return err
	}

	s.logger.Info("send code success", zap.String("code", code))

	return nil
}

func (s *userService) Login(ctx context.Context, req *v1.LoginReq) (*v1.LoginRespData, error) {
	if !validator.IsPhone(req.Phone) {
		return &v1.LoginRespData{}, errors.New("phone is invalidate")
	}

	key := constants.RedisLoginCodeKey + req.Phone
	cacheCode, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return &v1.LoginRespData{}, err
	}
	if req.Code != cacheCode {
		return &v1.LoginRespData{}, errors.New("code is invalidate")
	}

	user, err := s.userRepository.GetUserByPhone(ctx, req.Phone)
	if err != nil {
		return &v1.LoginRespData{}, err
	}
	if user == nil {
		user, err = s.createUserWithPhone(req.Phone)
		if err != nil {
			return &v1.LoginRespData{}, err
		}
	}

	token, err := random.UUIdV4()
	if err != nil {
		return nil, err
	}
	key = constants.RedisLoginUserKey + token
	err = s.rdb.HSet(ctx, key, map[string]string{
		"id":       strconv.Itoa(int(user.Id)),
		"nickname": user.NickName,
		"icon":     user.Icon,
	}).Err()
	if err != nil {
		return &v1.LoginRespData{}, err
	}

	ttl := time.Minute * constants.RedisLoginUserTTL
	err = s.rdb.Expire(ctx, key, ttl).Err()
	if err != nil {
		return &v1.LoginRespData{}, err
	}
	return &v1.LoginRespData{
		Token: token,
	}, nil
}

func (s *userService) GetMe(ctx context.Context) (*v1.GetMeRespData, error) {
	user := user_holder.GetUser(ctx)
	if user == nil {
		return nil, errors.New("user not found")
	}
	return (*v1.GetMeRespData)(user), nil
}

func (s *userService) createUserWithPhone(phone string) (*model.User, error) {
	user := model.User{
		Phone:    phone,
		NickName: constants.UserNickNamePrefix + random.RandString(10),
	}
	err := s.userRepository.CreateUser(nil, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
