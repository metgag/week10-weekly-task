package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/metgag/koda-weekly10/internals/handlers"
	"github.com/metgag/koda-weekly10/internals/repositories"
)

func InitCinemaRouter(router *gin.Engine, dbpool *pgxpool.Pool) {
	cr := repositories.NewCinemaRepository(dbpool)
	ch := handlers.NewCinemaHandler(cr)

	cinemaRouter := router.Group("/cinema")
	{
		cinemaRouter.GET("/schedule", ch.HandlerSchedule)
		cinemaRouter.GET("/seats", ch.HandlerSeats)
	}
}
