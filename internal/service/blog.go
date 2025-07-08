package service

import (
	"context"
	"go-dianping/internal/base/user_holder"
	"go-dianping/internal/model"
)

type BlogService interface {
	SaveBlog(ctx context.Context, blog *model.Blog) error
	LikeBlog(ctx context.Context, id uint64) error
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
	//	1. 获取登录用户
	userId := user_holder.GetUser(ctx).ID
	//	2. 判断当前登录用户是否已经点赞
	// TODO implement like logic
	return nil
}
