package v1

type (
	GetShopTypeListRespDataItem struct {
		Id   uint   `json:"id"`
		Name string `json:"name"`
		Icon string `json:"icon"`
		Sort uint   `json:"sort"`
	}
	GetShopTypeListRespData []*GetShopTypeListRespDataItem
	GetShopTypeListResp     struct {
		Response
		Data GetShopTypeListRespData `json:"data"`
	}
)
