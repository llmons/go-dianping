package v1

import "go-dianping/internal/model"

type (
	AddSeckillVoucherReq      model.Voucher
	AddSeckillVoucherRespData struct {
		ID int64 `json:"id"`
	}
)
