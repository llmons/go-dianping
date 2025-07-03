package v1

type (
	AddSeckillVoucherRespData struct {
		ID int64 `json:"id"`
	}
	AddSeckillVoucherResp struct {
		Response
		Data AddSeckillVoucherRespData `json:"data"`
	}
)
