package v1

type (
	SeckillVoucherReq struct {
		VoucherId uint64 `uri:"id" binding:"required"`
	}
)
