package api

type (
	GetShopTypeListRespDataItem struct {
		Id   uint   `json:"id"`
		Name string `json:"name"`
		Icon string `json:"icon"`
		Sort uint   `json:"sort"`
	}
	GetShopTypeListRespData []*GetShopTypeListRespDataItem
	GetShopTypeListResp     struct {
		response
		Data GetShopTypeListRespData `json:"data"`
	}
)
