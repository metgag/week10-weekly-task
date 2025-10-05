package utils

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/metgag/koda-weekly10/internals/models"
)

func LogCtxError(ctx *gin.Context, errHead, errClient string, errDev error, statusCode int) {
	log.Printf("%s\n\t%s", errHead, errDev.Error())
	ctx.JSON(statusCode, models.ErrorResponse{
		Success: false,
		Status:  statusCode,
		Error:   errClient,
	})
}

func PrintError(head string, rep int, err error) {
	log.Printf("%s %s %s", strings.Repeat("=", rep), head, strings.Repeat("=", rep))
	if err != nil {
		log.Println(err.Error())
	}
}
