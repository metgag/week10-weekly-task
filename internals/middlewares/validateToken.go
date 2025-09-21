package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/metgag/koda-weekly10/internals/utils"
	"github.com/metgag/koda-weekly10/pkg"
	"github.com/redis/go-redis/v9"
)

type validateResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func newAccessResponse(success bool, err string) validateResponse {
	return validateResponse{Success: success, Error: err}
}

func ValidateToken(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := ctx.GetHeader("Authorization")
		if bearerToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				newAccessResponse(false, "need token to access this feature"))
			return
		}

		splitToken := strings.Split(bearerToken, " ")
		if len(splitToken) != 2 || !strings.EqualFold(splitToken[0], "Bearer") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				newAccessResponse(false, "token unrecognized"))
			return
		}
		token := splitToken[1]

		redisKey := fmt.Sprintf("archie:blacklist_%s", token)
		result, err := rdb.Get(ctx, redisKey).Result()
		if err != nil && err != redis.Nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				newAccessResponse(false, "authorization service unavailable"))
			return
		}
		if result == "1" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				newAccessResponse(false, "access token revoked"))
			return
		}

		var claims pkg.Claims
		if err := claims.ValidateToken(token); err != nil {
			utils.PrintError("ERROR VALIDATE TOKEN", 16, err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				newAccessResponse(false, "token expired, sign in again to get a new valid token"))
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}
