package v1

import "go-dianping/internal/entity"

type (
	QueryShopByIDReq struct {
		ID *uint64 `uri:"id" binding:"required"`
	}
	QueryShopByIDRespData entity.Shop
	QueryShopByIDResp     struct {
		Response
		Data QueryShopByIDRespData `json:"data"`
	}
)

type (
	UpdateShopReq struct {
		ID        *int64   `json:"id" example:"1"`
		Name      *string  `json:"name" example:"120茶餐厅"`
		TypeId    *int64   `json:"typeId" example:"1"`
		Images    *string  `json:"images" example:""`
		Area      *string  `json:"area" example:"大关"`
		Address   *string  `json:"address" example:"金华路锦昌文华苑29号"`
		X         *float64 `json:"x" example:"0"`
		Y         *float64 `json:"y" example:"0"`
		AvgPrice  *int64   `json:"avgPrice" example:"80"`
		Sold      *int32   `json:"sold" example:"4215"`
		Comments  *int32   `json:"comments" example:"3035"`
		Score     *int32   `json:"score" example:"37"`
		OpenHours *string  `json:"openHours" example:"10:00-22:00"`
	}
)
