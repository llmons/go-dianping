package service

import (
	"context"
	"errors"
	v1 "go-dianping/api/v1"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/base/user_holder"
	"go-dianping/internal/model"
	"gorm.io/gorm"
	"strconv"
)

type BlogService interface {
	SaveBlog(ctx context.Context, blog *model.Blog) error
	LikeBlog(ctx context.Context, id uint64) error
	QueryMyBlog(ctx context.Context, current int) ([]*model.Blog, error)
	QueryHotBlog(ctx context.Context, current int) ([]*model.Blog, error)
	QueryBlogById(ctx context.Context, id uint64) (*model.Blog, error)
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

func (s *blogService) SaveBlog(_ context.Context, blog *model.Blog) error {
	return s.query.Blog.Save(blog)
}

func (s *blogService) LikeBlog(ctx context.Context, id uint64) error {
	// 1. 获取登陆用户
	userId := *user_holder.GetUser(ctx).ID
	// 2. 判断当前登录用户是否已经点赞
	key := "blog:liked:" + strconv.Itoa(int(id))
	isMember, err := s.rdb.SIsMember(ctx, key, strconv.Itoa(int(userId))).Result()
	if err != nil {
		return err
	}
	if !isMember {
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
			if err := s.rdb.SAdd(ctx, key, strconv.Itoa(int(userId))).Err(); err != nil {
				return err
			}
		}
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
			if err := s.rdb.SRem(ctx, key, strconv.Itoa(int(userId))).Err(); err != nil {
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
	userId := *user_holder.GetUser(ctx).ID
	// 2. 判断当前登录用户是否已经点赞
	key := "blog:liked:" + strconv.Itoa(int(blog.ID))
	isMember, err := s.rdb.SIsMember(ctx, key, strconv.Itoa(int(userId))).Result()
	blog.IsLike = isMember
	return err
}
