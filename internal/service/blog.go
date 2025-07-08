package service

import (
	"context"
	"go-dianping/internal/model"
)

type BlogService interface {
	GetBlog(ctx context.Context, id int64) (*model.Blog, error)
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

func (s *blogService) GetBlog(ctx context.Context, id int64) (*model.Blog, error) {
	return &model.Blog{}, nil
}
