package service

import (
	"context"
	"go-dianping/internal/base/user_holder"
	"go-dianping/internal/model"
)

type FollowService interface {
	Follow(ctx context.Context, followUserID uint64, follow bool) error
	IsFollow(ctx context.Context, followUserID uint64) (bool, error)
	FollowCommons(ctx context.Context, id uint64) error
}

func NewFollowService(
	service *Service,
) FollowService {
	return &followService{
		Service: service,
	}
}

type followService struct {
	*Service
}

func (s *followService) Follow(ctx context.Context, followUserID uint64, follow bool) error {
	// 1. 获取登录用户
	userID := user_holder.GetUser(ctx).ID

	//	1. 判断是关注还是取关
	if follow {
		//	2. 关注，新增数据
		if err := s.query.Follow.Save(&model.Follow{UserID: *userID, FollowUserID: followUserID}); err != nil {
			return err
		}
	} else {
		//	3. 取关，删除数据
		if _, err := s.query.Follow.Delete(&model.Follow{UserID: *userID, FollowUserID: followUserID}); err != nil {
			return err
		}
	}
	return nil
}

func (s *followService) IsFollow(ctx context.Context, followUserID uint64) (bool, error) {
	// 1. 获取登录用户
	userID := user_holder.GetUser(ctx).ID
	//	2. 查询是否关注
	count, err := s.query.Follow.
		Where(s.query.Follow.UserID.Eq(*userID)).
		Where(s.query.Follow.FollowUserID.Eq(followUserID)).Count()
	if err != nil {
		return false, err
	}
	//	3. 判断
	return count > 0, nil
}

func (s *followService) FollowCommons(ctx context.Context, id uint64) error {
	//TODO implement me
	panic("implement me")
}
