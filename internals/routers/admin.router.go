package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/metgag/koda-weekly10/internals/handlers"
	"github.com/metgag/koda-weekly10/internals/middlewares"
	"github.com/metgag/koda-weekly10/internals/repositories"
)

func InitAdminRouter(router *gin.Engine, dbpool *pgxpool.Pool) {
	mr := repositories.NewMovieRepository(dbpool)
	mh := handlers.NewMovieHandler(mr)

	adminGroup := router.Group("/admin")

	movieGroup := adminGroup.Group("/movies")
	movieGroup.Use(
		middlewares.ValidateToken,
		middlewares.Access("admin"),
	)

	{
		movieGroup.GET("/", mh.HandleGetAllMovie)
		movieGroup.DELETE("/:id", mh.HandleDeleteMovie)
		movieGroup.PATCH("/:id", mh.HandleMovieUpdate)
	}
}
