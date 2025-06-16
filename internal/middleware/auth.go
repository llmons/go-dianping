package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go-dianping/internal/base/constants"
	"time"
)

func Auth(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.AbortWithStatus(401)
			return
		}

		userField, err := rdb.HGetAll(ctx, constants.RedisLoginUserKey+token).Result()
		if err != nil {
			return
		}
		if len(userField) == 0 {
			ctx.AbortWithStatus(401)
			return
		}

		rdb.Expire(ctx, constants.RedisLoginUserKey+token, time.Minute*constants.RedisLoginUserTTL)
		ctx.Next()
	}
}
