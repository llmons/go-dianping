package v1

import "go-dianping/internal/model"

type (
	AddSeckillVoucherRespData struct {
		ID int64 `json:"id"`
	}
	AddSeckillVoucherResp struct {
		Response
		Data AddSeckillVoucherRespData `json:"data"`
	}
)

type (
	AddVoucherRespData struct {
		ID int64 `json:"id"`
	}
	AddVoucherResp struct {
		Response
		Data AddVoucherRespData `json:"data"`
	}
)

type (
	QueryVoucherOfShopReq struct {
		ShopId uint64 `uri:"shopId" binding:"required"`
	}
	QueryVoucherOfShopResp struct {
		Response
		Data []*model.Voucher `json:"data"`
	}
)
