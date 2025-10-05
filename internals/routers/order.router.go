package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/metgag/koda-weekly10/internals/handlers"
	"github.com/metgag/koda-weekly10/internals/middlewares"
	"github.com/metgag/koda-weekly10/internals/repositories"
	"github.com/redis/go-redis/v9"
)

func InitOrderRouter(router *gin.Engine, dbpool *pgxpool.Pool, rdb *redis.Client) {
	or := repositories.NewOrderRepository(dbpool)
	oh := handlers.NewOrderHandler(or)

	router.POST("/orders",
		middlewares.ValidateToken(rdb),
		middlewares.Access("user"),
		oh.HandleCreateOrder,
	)
}
