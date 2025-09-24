package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/metgag/koda-weekly10/internals/handlers"
	"github.com/metgag/koda-weekly10/internals/repositories"
	"github.com/redis/go-redis/v9"
)

func InitMovieRouter(router *gin.Engine, dbpool *pgxpool.Pool, rdb *redis.Client) {
	mr := repositories.NewMovieRepository(dbpool, rdb)
	mh := handlers.NewMovieHandler(mr)

	movieRouter := router.Group("api/movies")

	{
		movieRouter.GET("/upcoming", mh.GetUpcomingMovies)
		movieRouter.GET("/popular", mh.GetPopularMovies)
		movieRouter.GET("", mh.HandleMovieWithGenrePageSearch)
		movieRouter.GET("/:id", mh.GetMovieDetail)
		movieRouter.GET("/:id/schedules", mh.HandleGetMovieSchedule)
		movieRouter.GET("/:id/schedule", mh.HandleGetMovieScheduleFilter)
	}
}
