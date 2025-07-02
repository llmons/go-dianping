package v1

import "go-dianping/internal/entity"

type (
	AddSeckillVoucherReq      entity.Voucher
	AddSeckillVoucherRespData struct {
		ID int64 `json:"id"`
	}
)
