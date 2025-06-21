package v1

type (
	GetShopByIDReq struct {
		ID int64 `uri:"id" binding:"required"`
	}
	GetShopByIDRespData struct {
		ID        int64   `json:"id"`
		Name      string  `json:"name"`
		TypeId    int64   `json:"typeID"`
		Images    string  `json:"images"`
		Area      *string `json:"area"`
		Address   string  `json:"address"`
		X         float64 `json:"x"`
		Y         float64 `json:"y"`
		AvgPrice  *int64  `json:"avgPrice"`
		Sold      int32   `json:"sold"`
		Comments  int32   `json:"comments"`
		Score     int32   `json:"score"`
		OpenHours *string `json:"openHours"`
	}
	GetShopByIDResp struct {
		Response
		Data GetShopByIDRespData `json:"data"`
	}
)

type (
	UpdateShopReq struct {
		ID        int64   `json:"id"`
		Name      string  `json:"name"`
		TypeId    int64   `json:"typeId"`
		Images    string  `json:"images"`
		Area      *string `json:"area"`
		Address   string  `json:"address"`
		X         float64 `json:"x"`
		Y         float64 `json:"y"`
		AvgPrice  *int64  `json:"avgPrice"`
		Sold      int32   `json:"sold"`
		Comments  int32   `json:"comments"`
		Score     int32   `json:"score"`
		OpenHours *string `json:"openHours"`
	}
)
