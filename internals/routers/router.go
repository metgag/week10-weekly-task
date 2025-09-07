package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	docs "github.com/metgag/koda-weekly10/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(dbpool *pgxpool.Pool) *gin.Engine {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	InitAuthRouter(r, dbpool)
	InitMovieRouter(r, dbpool)
	InitUserRouter(r, dbpool)
	InitCinemaRouter(r, dbpool)
	InitOrderRouter(r, dbpool)

	return r
}
