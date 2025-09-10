package middlewares

import (
	"log"
	"net/http"
	"strings"

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
	bearerToken := ctx.GetHeader("Authorization")
	if bearerToken == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized,
			gin.H{
				"success": false,
				"error":   "need token to access this feature",
			})
		return
	}

	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) != 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized,
			gin.H{
				"success": false,
				"error":   "token unrecognized",
			})
		return
	}

	token := splitToken[1]
	var claims pkg.Claims
	if err := claims.ValidateToken(token); err != nil {
		log.Printf("%s VALIDATE TOKEN ERROR %s", strings.Repeat("=", 8), strings.Repeat("=", 8))
		log.Println(err.Error())
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
