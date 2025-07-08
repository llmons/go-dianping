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
