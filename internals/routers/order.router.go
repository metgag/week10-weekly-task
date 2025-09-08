package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/metgag/koda-weekly10/internals/handlers"
	"github.com/metgag/koda-weekly10/internals/middlewares"
	"github.com/metgag/koda-weekly10/internals/repositories"
)

func InitOrderRouter(router *gin.Engine, dbpool *pgxpool.Pool) {
	or := repositories.NewOrderRepository(dbpool)
	oh := handlers.NewOrderHandler(or)

	router.POST("/orders", middlewares.ValidateToken, oh.HandleCreateOrder)
	router.GET("/admin/orders",
		middlewares.ValidateToken,
		middlewares.Access("admin"),
		oh.HandleGetOrderHistory,
	)
}
