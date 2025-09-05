package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/metgag/koda-weekly10/internals/handlers"
	"github.com/metgag/koda-weekly10/internals/repositories"
)

func InitMovieRouter(router *gin.Engine, dbpool *pgxpool.Pool) {
	mr := repositories.NewMovieRepository(dbpool)
	mh := handlers.NewMovieHandler(mr)

	movieRouter := router.Group("/movie")
	{
		movieRouter.GET("/upcoming", mh.GetUpcomingMovies)
		movieRouter.GET("/popular", mh.GetPopularMovies)
		movieRouter.GET("/:id", mh.GetMovieDetail)
	}
}
