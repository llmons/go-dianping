package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go-dianping/api/v1"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/base/user_holder"
	"net/http"
	"strings"
)

func RefreshToken(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. 获取请求中的 token
		token := ctx.GetHeader("Authorization")
		token, found := strings.CutPrefix(token, "Bearer ")
		if !found {
			ctx.Next()
			return
		}
		// 2. 基于 token 获取 redis 中的用户
		key := constants.RedisLoginUserKey + token
		simpleUser := v1.SimpleUser{}
		if err := rdb.HGetAll(ctx, key).Scan(&simpleUser); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, v1.Response{
				Success:  false,
				ErrorMsg: err.Error(),
			})
			return
		}
		// 3. 判断用户是否存在
		if simpleUser.ID == nil {
			ctx.Next()
			return
		}
		// 6. 存在，保存用户信息到 context
		newCtx := user_holder.WithUser(ctx, &simpleUser)
		ctx.Request = ctx.Request.WithContext(newCtx)
		// 7. 刷新 token 有效期
		rdb.Expire(ctx, key, constants.RedisLoginUserTTL)
		// 8. 放行
		ctx.Next()
	}
}
