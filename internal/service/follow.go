package service

import (
	"context"
	v1 "go-dianping/api/v1"
	"go-dianping/internal/base/user_holder"
	"go-dianping/internal/model"
	"strconv"
)

type FollowService interface {
	Follow(ctx context.Context, followUserID uint64, follow bool) error
	IsFollow(ctx context.Context, followUserID uint64) (bool, error)
	FollowCommons(ctx context.Context, id uint64) ([]*v1.SimpleUser, error)
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

	key := "follows:" + strconv.Itoa(int(*userID))
	//	判断是关注还是取关
	if follow {
		//	2. 关注，新增数据
		if err := s.query.Follow.Save(&model.Follow{UserID: *userID, FollowUserID: followUserID}); err != nil {
			return err
		}
		//	把关注用户的 ID 放入 redis set 集合
		if err := s.rdb.SAdd(ctx, key, followUserID).Err(); err != nil {
			return err
		}
	} else {
		//	3. 取关，删除数据
		if _, err := s.query.Follow.Delete(&model.Follow{UserID: *userID, FollowUserID: followUserID}); err != nil {
			return err
		}
		//	把关注用户的 ID 从 redis set 集合中移除
		if err := s.rdb.SRem(ctx, key, followUserID).Err(); err != nil {
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

func (s *followService) FollowCommons(ctx context.Context, id uint64) ([]*v1.SimpleUser, error) {
	userID := user_holder.GetUser(ctx).ID
	key := "follows:" + strconv.Itoa(int(*userID))
	key2 := "follows:" + strconv.Itoa(int(id))
	interset, err := s.rdb.SInter(ctx, key, key2).Result()
	if err != nil {
		return nil, err
	}
}
