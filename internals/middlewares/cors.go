package middlewares

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(ctx *gin.Context) {
	// allowed pages
	whitelist := []string{
		"http://localhost:5173",
		"http://localhost:6096",
	}

	// curr origin
	origin := ctx.GetHeader("Origin")
	if slices.Contains(whitelist, origin) {
		ctx.Header("Access-Control-Allow-Origin", origin)
	}
	// else {
	// 	utils.PrintError("ORIGIN NOT WHITELISTED", 20, nil)
	// }

	ctx.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
	ctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
	// handler jika bertemu preflight
	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	ctx.Next()
}
