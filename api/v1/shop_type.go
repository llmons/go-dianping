package v1

type (
	QueryTypeListRespDataItem struct {
		ID   int64   `json:"id" redis:"id"`
		Name *string `json:"name" redis:"name"`
		Icon *string `json:"icon" redis:"icon"`
		Sort *int32  `json:"sort"  redis:"sort"`
	}
	QueryTypeListResp struct {
		Response
		Data []QueryTypeListRespDataItem `json:"data"`
	}
)
