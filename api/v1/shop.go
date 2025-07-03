package v1

import "go-dianping/internal/model"

type (
	QueryShopByIDReq struct {
		ID *uint64 `uri:"id" binding:"required"`
	}
	QueryShopByIDResp struct {
		Response
		Data model.Shop `json:"data"`
	}
)
