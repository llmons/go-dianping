package service

import (
	"context"
	"go-dianping/internal/model"
)

type BlogService interface {
	SaveBlog(ctx context.Context, blog *model.Blog) error
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
