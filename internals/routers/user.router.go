package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/metgag/koda-weekly10/internals/handlers"
	"github.com/metgag/koda-weekly10/internals/repositories"
)

func InitUserRouter(r *gin.Engine, dbpool *pgxpool.Pool) {
	ur := repositories.NewUserRepository(dbpool)
	uh := handlers.NewUserHandler(ur)

	userGroup := r.Group("/userinf")
	{
		userGroup.GET("/:uid", uh.HandleUserinf)
		userGroup.PATCH("/:uid", uh.HandleUpdateUserInf)
		userGroup.GET("/:uid/order", uh.HandleUserOrderHistory)
	}
}
