package v1

import "go-dianping/internal/model"

type (
	QueryMyBlogResp struct {
		Response
		Data []model.Blog `json:"data"`
	}
)

type (
	QueryHotBlogResp struct {
		Response
		Data []model.Blog `json:"data"`
	}
)

type (
	QueryBlogByIDResp struct {
		Response
		Data *model.Blog `json:"data"`
	}
)

type (
	QueryBlogByUserIDResp struct {
		Response
		Data []*model.Blog `json:"data"`
	}
)

type (
	ScrollResult[T any] struct {
		List    []*T    `json:"list"`
		MinTime float64 `json:"min_time"`
		Offset  int     `json:"offset"`
	}
	QueryBlogOfFollowResp[T any] struct {
		Response
		Data ScrollResult[T] `json:"data"`
	}
)
