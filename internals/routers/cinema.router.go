package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/metgag/koda-weekly10/internals/handlers"
	"github.com/metgag/koda-weekly10/internals/middlewares"
	"github.com/metgag/koda-weekly10/internals/repositories"
	"github.com/redis/go-redis/v9"
)

func InitCinemaRouter(router *gin.Engine, dbpool *pgxpool.Pool, rdb *redis.Client) {
	cr := repositories.NewCinemaRepository(dbpool)
	ch := handlers.NewCinemaHandler(cr)

	cinemaRouter := router.Group("/cinemas")
	cinemaRouter.Use(
		middlewares.ValidateToken(rdb),
		middlewares.Access("user", "admin"),
	)

	{
		cinemaRouter.GET("/schedules", ch.HandlerSchedule)
		cinemaRouter.GET("/:schedule_id/seats", ch.HandlerSeats)
		cinemaRouter.GET("/:schedule_id/selected", ch.HandlerCinemaNameAndTime)
	}
}
