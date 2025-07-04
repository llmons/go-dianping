package v1

var (
	ErrInternalServerError = newError(500, "内部服务器错误")

	ErrPhoneIsInvalid = newError(1101, "手机号格式错误")
	ErrCodeIsInvalid  = newError(1102, "验证码错误")
	ErrCanNotGetUser  = newError(1103, "获取用户信息失败")

	ErrSeckillNotStart   = newError(1201, "秒杀尚未开始")
	ErrSeckillIsEnd      = newError(1202, "秒杀已经结束")
	ErrInsufficientStock = newError(1203, "库存不足")
	ErrAlreadySeckill    = newError(1204, "用户已经购买过一次")
)
