package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/metgag/koda-weekly10/pkg"
)

// type accessResponse struct {
// 	success bool
// 	error   string
// }

// func newAccessResponse(success bool, err string) accessResponse {
// 	return accessResponse{success: success, error: err}
// }

func ValidateToken(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized,
			gin.H{
				"success": false,
				"error":   "need token to access this feature",
			})
		return
	}

	var claims pkg.Claims
	if err := claims.ValidateToken(token); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized,
			gin.H{
				"success": false,
				"error":   "token expired, sign in again to get a new valid token",
			})
		return
	}

	ctx.Set("claims", claims)
	ctx.Next()
}
