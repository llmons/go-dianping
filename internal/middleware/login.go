package middleware

import (
	"github.com/gin-gonic/gin"
	"go-dianping/internal/base/user_holder"
)

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. 判断是否需要拦截（context 中是否有用户）
		user := user_holder.GetUser(ctx.Request.Context())
		if user == nil {
			// 没有，需要拦截，设置状态码
			ctx.AbortWithStatus(401)
			// 拦截
			return
		}
		// 有用户，则放行
		ctx.Next()
	}
}
