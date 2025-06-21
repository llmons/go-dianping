package v1

type (
	GetShopTypeListRespDataItem struct {
		ID   int64   `json:"id"`
		Name *string `json:"name"`
		Icon *string `json:"icon"`
		Sort *int32  `json:"sort"`
	}
	GetShopTypeListRespData []*GetShopTypeListRespDataItem
	GetShopTypeListResp     struct {
		Response
		Data GetShopTypeListRespData `json:"data"`
	}
)
