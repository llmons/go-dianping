package api

type (
	GetShopByIdReq struct {
		Id string `path:"id" binding:"required"`
	}
	GetShopByIdRespData struct {
	}
	GetShopByIdResp struct{}
)
