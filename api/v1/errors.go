package v1

var (
	ErrPhoneIsInvalid = newError(1101, "手机号格式错误")
	ErrCodeIsInvalid  = newError(1102, "验证码错误")
	ErrCanNotGetUser  = newError(1103, "获取用户信息失败")

	ErrSeckillNotStart   = newError(1201, "秒杀尚未开始")
	ErrSeckillIsEnd      = newError(1202, "秒杀已经结束")
	ErrInsufficientStock = newError(1203, "库存不足")
	ErrAlreadySeckill    = newError(1204, "用户已经购买过一次")
	ErrNotAllowDoubleBuy = newError(1205, "不能重复下单")

	ErrIncorrectFilename = newError(1301, "错误的文件名称")
)
