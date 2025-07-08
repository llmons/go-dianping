package service

import (
	"context"
	"go-dianping/internal/model"
)

type BlogCommentsService interface {
	GetBlogComments(ctx context.Context, id int64) (*model.BlogComments, error)
}

func NewBlogCommentsService(
	service *Service,
) BlogCommentsService {
	return &blogCommentsService{
		Service: service,
	}
}

type blogCommentsService struct {
	*Service
}

func (s *blogCommentsService) GetBlogComments(ctx context.Context, id int64) (*model.BlogComments, error) {
	return &model.BlogComments{}, nil
}
