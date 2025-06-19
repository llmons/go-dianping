package api

type (
	GetShopByIdReq struct {
		Id string `uri:"id" binding:"required"`
	}
	GetShopByIdRespData struct {
		Id       string  `json:"id"`
		Name     string  `json:"name"`
		TypeId   string  `json:"type_id"`
		Images   string  `json:"images"`
		Area     string  `json:"area"`
		Address  string  `json:"address"`
		X        float64 `json:"x"`
		Y        float64 `json:"y"`
		AvgPrice uint    `json:"avg_price"`
		Sold     uint    `json:"sold"`
		Comments uint    `json:"comments"`
		Score    uint    `json:"score"`
	}
	GetShopByIdResp struct {
		response
		Data GetShopByIdRespData `json:"data"`
	}
)
