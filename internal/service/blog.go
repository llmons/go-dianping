package service

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	v1 "go-dianping/api/v1"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/base/user_holder"
	"go-dianping/internal/model"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type BlogService interface {
	SaveBlog(ctx context.Context, blog *model.Blog) (uint64, error)
	LikeBlog(ctx context.Context, id uint64) error
	QueryMyBlog(ctx context.Context, current int) ([]*model.Blog, error)
	QueryHotBlog(ctx context.Context, current int) ([]*model.Blog, error)
	QueryBlogById(ctx context.Context, id uint64) (*model.Blog, error)
	QueryBlogLikes(ctx context.Context, id uint64) ([]*v1.SimpleUser, error)
	QueryBlogByUserID(ctx context.Context, id uint64, current int) ([]*model.Blog, error)
}

func NewBlogService(
	service *Service,
) BlogService {
	return &blogService{
		Service: service,
	}
}

type blogService struct {
	*Service
}

func (s *blogService) SaveBlog(ctx context.Context, blog *model.Blog) (uint64, error) {
	// 1. 获取登录用户
	user := user_holder.GetUser(ctx)
	blog.UserID = *user.ID
	//	2. 保存探店笔记
	if err := s.query.Blog.Save(blog); err != nil {
		return 0, err
	}
	// 3. 查询笔记作者的所有粉丝
	follows, err := s.query.Follow.Where(s.query.Follow.FollowUserID.Eq(blog.UserID)).Find()
	if err != nil {
		return 0, err
	}
	// 4. 推送笔记 id 给所有粉丝
	for _, follow := range follows {
		userID := follow.FollowUserID
		key := "feed:" + strconv.Itoa(int(userID))
		if err := s.rdb.ZAdd(ctx, key, redis.Z{
			Score:  float64(time.Now().UnixMilli()),
			Member: strconv.Itoa(int(blog.ID)),
		}).Err(); err != nil {
			return 0, err
		}
	}
	// 5. 返回 id
	return blog.ID, nil
}

func (s *blogService) LikeBlog(ctx context.Context, id uint64) error {
	// 1. 获取登陆用户
	userId := *user_holder.GetUser(ctx).ID
	// 2. 判断当前登录用户是否已经点赞
	key := constants.RedisBlogLikeKey + strconv.Itoa(int(id))
	_, err := s.rdb.ZScore(ctx, key, strconv.Itoa(int(userId))).Result()
	if errors.Is(err, redis.Nil) {
		// 3. 如果未点赞，可以点赞
		// 3.1. 数据库点赞数 +1
		info, err := s.query.Blog.
			Where(s.query.Blog.ID.Eq(id)).
			Update(s.query.Blog.Liked, s.query.Blog.Liked.Add(1))
		if err != nil {
			return err
		}
		// 3.2. 保存用户到 redis 的 set 集合
		if info.RowsAffected != 0 {
			if err := s.rdb.ZAdd(ctx, key, redis.Z{
				Score:  float64(time.Now().UnixMilli()),
				Member: strconv.Itoa(int(userId)),
			}).Err(); err != nil {
				return err
			}
		}
	} else if err != nil {
		return err
	} else {
		// 4. 如果已经点赞，取消点赞
		// 4.1. 数据库点赞数 -1
		info, err := s.query.Blog.
			Where(s.query.Blog.ID.Eq(id)).
			Update(s.query.Blog.Liked, s.query.Blog.Liked.Sub(1))
		if err != nil {
			return err
		}
		// 4.2. 把用户从 redis 的 set 集合中删除
		if info.RowsAffected != 0 {
			if err := s.rdb.ZRem(ctx, key, strconv.Itoa(int(userId))).Err(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *blogService) QueryMyBlog(ctx context.Context, current int) ([]*model.Blog, error) {
	user := user_holder.GetUser(ctx)
	result, _, err := s.query.Blog.
		Where(s.query.Blog.UserID.Eq(*user.ID)).
		FindByPage(current, constants.MaxPageSize)
	return result, err
}

func (s *blogService) QueryHotBlog(ctx context.Context, current int) ([]*model.Blog, error) {
	result, _, err := s.query.Blog.
		Order(s.query.Blog.Liked.Desc()).
		FindByPage(current, constants.MaxPageSize)
	for _, blog := range result {
		if err := s.queryBlogUser(blog); err != nil {
			return nil, err
		}
		if err := s.isBlogLiked(ctx, blog); err != nil {
			return nil, err
		}
	}
	return result, err
}

func (s *blogService) QueryBlogById(ctx context.Context, id uint64) (*model.Blog, error) {
	//	1. 查询 blog
	blog, err := s.query.Blog.Where(s.query.Blog.ID.Eq(id)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, v1.ErrBlogNotFound
	} else if err != nil {
		return nil, err
	}
	//	2. 查询 blog 有关的用户
	if err := s.queryBlogUser(blog); err != nil {
		return nil, err
	}
	// 3. 查询 blog 是否被点赞
	return blog, s.isBlogLiked(ctx, blog)
}

func (s *blogService) queryBlogUser(blog *model.Blog) error {
	userID := blog.UserID
	user, err := s.query.User.Where(s.query.User.ID.Eq(userID)).First()
	if err != nil {
		return err
	}
	blog.Name = user.NickName
	blog.Icon = user.Icon
	return nil
}

func (s *blogService) isBlogLiked(ctx context.Context, blog *model.Blog) error {
	// 1. 获取登陆用户
	user := user_holder.GetUser(ctx)
	if user == nil {
		// 如果没有登录用户，直接返回
		return nil
	}
	userId := *user.ID
	// 2. 判断当前登录用户是否已经点赞
	key := constants.RedisBlogLikeKey + strconv.Itoa(int(blog.ID))
	_, err := s.rdb.ZScore(ctx, key, strconv.Itoa(int(userId))).Result()
	blog.IsLike = !errors.Is(err, redis.Nil)
	return err
}

func (s *blogService) QueryBlogLikes(ctx context.Context, id uint64) ([]*v1.SimpleUser, error) {
	//	 1. 查询 top5 的点赞用户
	key := constants.RedisBlogLikeKey + strconv.Itoa(int(id))
	top5, err := s.rdb.ZRange(ctx, key, 0, 4).Result()
	if errors.Is(err, redis.Nil) {
		// 如果没有点赞用户，直接返回空
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	//	2. 解析出其中的用户 FollowUserId
	ids := make([]uint64, len(top5))
	for i, idStr := range top5 {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, err
		}
		ids[i] = uint64(id)
	}
	//	3. 根据用户 FollowUserId 查询用户
	result, err := s.query.User.Where(s.query.User.ID.In(ids...)).Order(s.query.User.ID.Field(ids...)).Find()
	if err != nil {
		return nil, err
	}
	users := make([]*v1.SimpleUser, len(ids))
	for i, user := range result {
		users[i] = &v1.SimpleUser{
			ID:       &user.ID,
			NickName: user.NickName,
			Icon:     user.Icon,
		}
	}
	//	4. 返回
	return users, nil
}

func (s *blogService) QueryBlogByUserID(ctx context.Context, id uint64, current int) ([]*model.Blog, error) {
	result, _, err := s.query.Blog.Where(s.query.Blog.UserID.Eq(id)).FindByPage(current, constants.MaxPageSize)
	return result, err
}
