package v1

type (
	QueryTypeListRespDataItem struct {
		ID   int64   `json:"id"`
		Name *string `json:"name"`
		Icon *string `json:"icon"`
		Sort *int32  `json:"sort"`
	}
	QueryTypeListRespData []*QueryTypeListRespDataItem
	QueryTypeListResp     struct {
		Response
		Data QueryTypeListRespData `json:"data"`
	}
)
