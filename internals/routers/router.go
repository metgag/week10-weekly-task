package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRouter(dbpool *pgxpool.Pool) *gin.Engine {
	r := gin.Default()

	InitAuthRouter(r, dbpool)
	InitMovieRouter(r, dbpool)
	InitUserRouter(r, dbpool)
	InitCinemaRouter(r, dbpool)
	InitOrderRouter(r, dbpool)

	return r
}
