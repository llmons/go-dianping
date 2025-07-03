package v1

import "go-dianping/internal/model"

type (
	QueryTypeListResp struct {
		Response
		Data []*model.User `json:"data"`
	}
)
