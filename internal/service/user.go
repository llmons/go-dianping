package service

import (
	"context"
	"github.com/duke-git/lancet/v2/random"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-dianping/api/v1"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/base/regex_utils"
	"go-dianping/internal/base/user_holder"
	"go-dianping/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService interface {
	SendCode(ctx context.Context, req *v1.SendCodeReq) error
	Login(ctx context.Context, req *v1.LoginReq) (*v1.LoginRespData, error)
	Me(ctx context.Context) (*v1.SimpleUser, error)
}

type userService struct {
	*Service
}

func NewUserService(service *Service) UserService {
	return &userService{
		Service: service,
	}
}

func (s *userService) SendCode(ctx context.Context, req *v1.SendCodeReq) error {
	// 1. 校验手机号
	if regex_utils.IsPhoneInvalid(req.Phone) {
		// 2. 如果不符合，返回错误信息
		return v1.ErrPhoneIsInvalid
	}
	// 3. 符合，生成验证码
	var code string
	//if s.conf.Get("env") == "prod" {
	//	code = random.RandNumeral(6)
	//} else {
	code = "123456"
	//}

	// 4. 保存验证码到 redis
	key := constants.RedisLoginCodeKey + req.Phone
	if err := s.rdb.Set(ctx, key, code, constants.RedisLoginCodeTTL).Err(); err != nil {
		return err
	}

	// 5. 发送验证码，接入第三方服务
	s.logger.Info("发送短信验证码成功", zap.String("验证码：", code))
	// 返回 ok
	return nil
}

func (s *userService) Login(ctx context.Context, req *v1.LoginReq) (*v1.LoginRespData, error) {
	// 1. 校验手机号
	if regex_utils.IsPhoneInvalid(req.Phone) {
		// 2. 如果不符合，返回错误信息
		return nil, v1.ErrPhoneIsInvalid
	}

	// 3. 从 redis 获取验证码并校验
	key := constants.RedisLoginCodeKey + req.Phone
	cacheCode, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if cacheCode != req.Code {
		// 不一致，报错
		return nil, v1.ErrCodeIsInvalid
	}

	// 4. 一致，根据手机号查询用户 select * from tb_user where phone = ?
	user, err := s.query.User.Where(s.query.User.Phone.Eq(req.Phone)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 5. 判断用户是否存在
	if user == nil {
		// 6. 不存在，创建新用户并保存
		user, err = s.createUserWithPhone(req.Phone)
		if err != nil {
			return nil, err
		}
	}

	// 7. 保存用户信息到 redis 中
	// 7.1. 随机生成 token，作为登录令牌
	token, err := random.UUIdV4()
	if err != nil {
		return nil, err
	}
	// 7.2. 存储 User 对象
	var simpleUser v1.SimpleUser
	if err := copier.Copy(&simpleUser, &user); err != nil {
		return nil, err
	}

	// 7.3. 存储
	key = constants.RedisLoginUserKey + token
	if err := s.rdb.HSet(ctx, key, simpleUser).Err(); err != nil {
		return nil, err
	}
	// 7.4. 设置 token 有效期
	if err := s.rdb.Expire(ctx, key, constants.RedisLoginUserTTL).Err(); err != nil {
		return nil, err
	}

	// 8. 返回 token
	return &v1.LoginRespData{
		Token: token,
	}, nil
}

func (s *userService) Me(ctx context.Context) (*v1.SimpleUser, error) {
	user := user_holder.GetUser(ctx)
	if user == nil {
		return nil, v1.ErrCanNotGetUser
	}
	return user, nil
}

func (s *userService) createUserWithPhone(phone string) (*model.User, error) {
	// 1. 创建用户
	nickname := constants.UserNickNamePrefix + random.RandString(10)
	user := model.User{
		Phone:    phone,
		NickName: &nickname,
	}
	// 2. 保存用户
	if err := s.query.User.Save(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
