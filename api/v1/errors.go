package v1

var (
	ErrPhoneIsInvalid = newError(1001, "手机号格式错误")
	ErrCodeIsInvalid  = newError(1002, "验证码错误")
	ErrShopIDIsNull   = newError(1003, "商铺ID不能为空")
)
