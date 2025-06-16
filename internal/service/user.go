package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go-dianping/api"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/model"
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
	Login(ctx *gin.Context, params *api.LoginReq) (*api.LoginResp, error)
	GetMe(ctx *gin.Context) (*api.SimpleUser, error)
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

	key, ttl := constants.RedisLoginCodeKey+phone, time.Minute*constants.RedisLoginCodeTTL
	err := s.rdb.Set(ctx, key, code, ttl).Err()
	if err != nil {
		return err
	}

	s.logger.Info("send code success", zap.String("code", code))

	return nil
}

func (s *userService) Login(ctx *gin.Context, params *api.LoginReq) (*api.LoginResp, error) {
	if !validator.IsPhone(params.Phone) {
		return &api.LoginResp{}, errors.New("phone is invalidate")
	}

	key := constants.RedisLoginCodeKey + params.Phone
	cacheCode, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return &api.LoginResp{}, err
	}
	if params.Code != cacheCode {
		return &api.LoginResp{}, errors.New("code is invalidate")
	}

	user, err := s.userRepository.GetUserByPhone(params.Phone)
	if err != nil {
		return &api.LoginResp{}, err
	}
	if user == nil {
		user, err = s.createUserWithPhone(params.Phone)
		if err != nil {
			return &api.LoginResp{}, err
		}
	}

	token := uuid.GenUUID()
	key = constants.RedisLoginUserKey + token
	err = s.rdb.HSet(ctx, key, map[string]string{
		"id":       strconv.Itoa(int(user.Id)),
		"nickname": user.NickName,
		"icon":     user.Icon,
	}).Err()
	if err != nil {
		return &api.LoginResp{}, err
	}

	ttl := time.Minute * constants.RedisLoginUserTTL
	err = s.rdb.Expire(ctx, key, ttl).Err()
	if err != nil {
		return &api.LoginResp{}, err
	}
	return &api.LoginResp{
		Token: token,
	}, nil
}

func (s *userService) GetMe(ctx *gin.Context) (*api.SimpleUser, error) {
	result, err := s.rdb.HGetAll(ctx, constants.RedisLoginUserKey).Result()
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("user not found")
	}
	var user api.SimpleUser
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
