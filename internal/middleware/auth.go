package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go-dianping/api"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/base/user_holder"
	"strconv"
	"time"
)

func Auth(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.AbortWithStatus(401)
			return
		}

		key := constants.RedisLoginUserKey + token
		userField, err := rdb.HGetAll(ctx, key).Result()
		if err != nil {
			return
		}
		if len(userField) == 0 {
			ctx.AbortWithStatus(401)
			return
		}

		idStr := userField["id"]
		idInt, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.AbortWithStatus(401)
		}
		newCtx := user_holder.WithUser(ctx, &api.SimpleUser{
			Id:       uint(idInt),
			NickName: userField["nickname"],
			Icon:     userField["icon"],
		})
		ctx.Request = ctx.Request.WithContext(newCtx)

		key, ttl := constants.RedisLoginUserKey+token, time.Minute*constants.RedisLoginUserTTL
		rdb.Expire(ctx, key, ttl)
		ctx.Next()
	}
}
