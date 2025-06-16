package middleware

import (
	"github.com/gin-gonic/gin"
	"go-dianping/internal/base/user_holder"
)

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := user_holder.GetUser(ctx.Request.Context())
		if user == nil {
			ctx.AbortWithStatus(401)
			return
		}
		ctx.Next()
	}
}
