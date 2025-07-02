package v1

var (
	ErrShopIDIsNull = newError(1001, "商铺ID不能为空")

	ErrPhoneIsInvalid = newError(1101, "手机号格式错误")
	ErrCodeIsInvalid  = newError(1102, "验证码错误")

	ErrSeckillNotStart   = newError(1201, "秒杀尚未开始")
	ErrSeckillIsEnd      = newError(1202, "秒杀已经结束")
	ErrInsufficientStock = newError(1203, "库存不足")
)
