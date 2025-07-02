package v1

type (
	SeckillVoucherReq struct {
		ID uint64 `uri:"id" binding:"required"`
	}
)
