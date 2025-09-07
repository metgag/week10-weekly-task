package middlewares

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/metgag/koda-weekly10/pkg"
)

func Access(roles ...string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		claims, isExists := ctx.Get("claims")
		if !isExists {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "you dont have permission to access this feature",
			})
			return
		}

		user, ok := claims.(pkg.Claims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "server unable to cast claims",
			})
			return
		}

		if !slices.Contains(roles, user.Role) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "mismatch permission to access this feature",
			})
			return
		}

		ctx.Next()
	}
}
