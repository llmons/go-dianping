package v1

type (
	SeckillVoucherReq struct {
		VoucherID uint64 `uri:"id" binding:"required"`
	}
)
