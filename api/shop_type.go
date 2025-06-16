package api

type GetShopTypeListReq struct{}

type GetShopTypeListRespData []struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Sort int    `json:"sort"`
}

type GetShopTypeListResp struct {
	response
	Data GetShopTypeListRespData `json:"data"`
}
