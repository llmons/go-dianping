package service

import (
	"context"
	"errors"
	v1 "go-dianping/api/v1"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/base/user_holder"
	"go-dianping/internal/model"
	"gorm.io/gorm"
)

type BlogService interface {
	SaveBlog(ctx context.Context, blog *model.Blog) error
	LikeBlog(ctx context.Context, id uint64) error
	QueryMyBlog(ctx context.Context, current int) ([]*model.Blog, error)
	QueryHotBlog(ctx context.Context, current int) ([]*model.Blog, error)
	QueryById(ctx context.Context, id uint64) (*model.Blog, error)
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

func (s *blogService) LikeBlog(_ context.Context, id uint64) error {
	if _, err := s.query.Blog.
		Where(s.query.Blog.ID.Eq(id)).
		Update(s.query.Blog.Liked, s.query.Blog.Liked.Add(1)); err != nil {
		return err
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

func (s *blogService) QueryHotBlog(_ context.Context, current int) ([]*model.Blog, error) {
	result, _, err := s.query.Blog.
		Order(s.query.Blog.Liked.Desc()).
		FindByPage(current, constants.MaxPageSize)
	for _, blog := range result {
		if err := s.queryBlogUser(blog); err != nil {
			return nil, err
		}
	}
	return result, err
}

func (s *blogService) QueryById(ctx context.Context, id uint64) (*model.Blog, error) {
	//	1. 查询 blog
	blog, err := s.query.Blog.Where(s.query.Blog.ID.Eq(id)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, v1.ErrBlogNotFound
	} else if err != nil {
		return nil, err
	}
	//	2. 查询 blog 有关的用户
	return blog, s.queryBlogUser(blog)
}

func (s *blogService) queryBlogUser(blog *model.Blog) error {
	userID := blog.ID
	user, err := s.query.User.Where(s.query.User.ID.Eq(userID)).First()
	if err != nil {
		return err
	}
	blog.Name = user.NickName
	blog.Icon = user.Icon
	return nil
}
