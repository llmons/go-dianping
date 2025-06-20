package v1

type (
	GetShopByIdReq struct {
		Id string `uri:"id" binding:"required"`
	}
	GetShopByIdRespData struct {
		Id        uint    `json:"id"`
		Name      string  `json:"name"`
		TypeId    uint    `json:"typeId"`
		Images    string  `json:"images"`
		Area      string  `json:"area"`
		Address   string  `json:"address"`
		X         float64 `json:"x"`
		Y         float64 `json:"y"`
		AvgPrice  uint    `json:"avgPrice"`
		Sold      uint    `json:"sold"`
		Comments  uint    `json:"comments"`
		Score     uint    `json:"score"`
		OpenHours string  `json:"openHours"`
	}
	GetShopByIdResp struct {
		Response
		Data GetShopByIdRespData `json:"data"`
	}
)

type (
	UpdateShopReq struct {
		Id        uint     `json:"id" bind:"required"`
		Name      *string  `json:"name"`
		TypeId    *uint    `json:"type_id"`
		Images    *string  `json:"images"`
		Area      *string  `json:"area"`
		Address   *string  `json:"address"`
		X         *float64 `json:"x"`
		Y         *float64 `json:"y"`
		AvgPrice  *uint    `json:"avg_price"`
		Sold      *uint    `json:"sold"`
		Comments  *uint    `json:"comments"`
		Score     *uint    `json:"score"`
		OpenHours *string  `json:"open_hours"`
	}
)
