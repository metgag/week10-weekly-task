package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	docs "github.com/metgag/koda-weekly10/docs"
	"github.com/metgag/koda-weekly10/internals/middlewares"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(dbpool *pgxpool.Pool, rdb *redis.Client) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.CORSMiddleware)

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Static("backdrop", "public/backdrop")
	r.Static("poster", "public/poster")
	r.Static("user", "public/user")

	InitAuthRouter(r, dbpool, rdb)
	InitMovieRouter(r, dbpool, rdb)
	InitUserRouter(r, dbpool, rdb)
	InitCinemaRouter(r, dbpool, rdb)
	InitOrderRouter(r, dbpool, rdb)
	InitAdminRouter(r, dbpool, rdb)

	return r
}
