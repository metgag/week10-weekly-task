package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/metgag/koda-weekly10/internals/handlers"
	"github.com/metgag/koda-weekly10/internals/middlewares"
	"github.com/metgag/koda-weekly10/internals/repositories"
	"github.com/redis/go-redis/v9"
)

func InitAdminRouter(router *gin.Engine, dbpool *pgxpool.Pool, rdb *redis.Client) {
	or := repositories.NewOrderRepository(dbpool)
	oh := handlers.NewOrderHandler(or)

	mr := repositories.NewMovieRepository(dbpool, rdb)
	mh := handlers.NewMovieHandler(mr)

	adminGroup := router.Group("/admin")
	// adminGroup.Use(
	// 	middlewares.ValidateToken,
	// 	middlewares.Access("admin"),
	// )

	adminGroup.GET("/orders",
		middlewares.ValidateToken(rdb),
		middlewares.Access("admin"),
		oh.HandleGetOrderHistory,
	)

	movieGroup := adminGroup.Group("/movies")
	{
		movieGroup.GET("/", mh.HandleGetAllMovie)
		movieGroup.DELETE("/:id", mh.HandleDeleteMovie)
		movieGroup.PATCH("/:id", mh.HandleMovieUpdate)
		movieGroup.POST("/", mh.HandleCreateMovie)
	}
}
