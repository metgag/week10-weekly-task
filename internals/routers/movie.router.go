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

	movieRouter := router.Group("/movies")
	{
		movieRouter.GET("/upcoming", mh.GetUpcomingMovies)
		movieRouter.GET("/popular", mh.GetPopularMovies)
		movieRouter.GET("/", mh.HandleMovieWithGenrePageSearch)
		movieRouter.GET("/:id", mh.GetMovieDetail)
	}

	// adminRouter := router.Group("/admin/movies")
	// {
	// 	adminRouter.GET("/",
	// 		middlewares.ValidateToken,
	// 		middlewares.Access("admin"),
	// 		mh.HandleGetAllMovie,
	// 	)
	// 	adminRouter.DELETE("/:id",
	// 		middlewares.ValidateToken,
	// 		middlewares.Access("admin"),
	// 		mh.HandleDeleteMovie,
	// 	)
	// 	adminRouter.PATCH("/:id",
	// 		middlewares.ValidateToken,
	// 		middlewares.Access("admin"),
	// 		mh.HandleMovieUpdate,
	// 	)
	// }
}
