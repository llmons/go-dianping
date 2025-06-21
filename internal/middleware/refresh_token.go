package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go-dianping/api/v1"

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
		if !found {
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
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.AbortWithStatus(500)
		}
		nickname, icon := userField["nickname"], userField["icon"]
		newCtx := user_holder.WithUser(ctx, &v1.SimpleUser{
			ID:       int64(id),
			NickName: &nickname,
			Icon:     &icon,
		})
		ctx.Request = ctx.Request.WithContext(newCtx)

		// ========== refresh token ==========
		key, ttl := constants.RedisLoginUserKey+token, time.Minute*constants.RedisLoginUserTTL
		rdb.Expire(ctx, key, ttl)
		ctx.Next()
	}
}
