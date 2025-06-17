package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go-dianping/api"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/base/user_holder"
	"strconv"
	"strings"
	"time"
)

func RefreshToken(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		token, found := strings.CutPrefix(token, "Bearer ")
		if found == false {
			ctx.Next()
			return
		}

		// ========== query user field from redis ==========
		key := constants.RedisLoginUserKey + token
		userField, err := rdb.HGetAll(ctx, key).Result()
		if err != nil {
			return
		}
		if len(userField) == 0 {
			ctx.Next()
			return
		}

		// ========== save user to user holder by context ==========
		idStr := userField["id"]
		idInt, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.AbortWithStatus(500)
		}
		newCtx := user_holder.WithUser(ctx, &api.SimpleUser{
			Id:       uint(idInt),
			NickName: userField["nickname"],
			Icon:     userField["icon"],
		})
		ctx.Request = ctx.Request.WithContext(newCtx)

		// ========== refresh token ==========
		key, ttl := constants.RedisLoginUserKey+token, time.Minute*constants.RedisLoginUserTTL
		rdb.Expire(ctx, key, ttl)
		ctx.Next()
	}
}
