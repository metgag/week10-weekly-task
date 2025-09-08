package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/metgag/koda-weekly10/internals/handlers"
	"github.com/metgag/koda-weekly10/internals/middlewares"
	"github.com/metgag/koda-weekly10/internals/repositories"
)

func InitUserRouter(r *gin.Engine, dbpool *pgxpool.Pool) {
	ur := repositories.NewUserRepository(dbpool)
	uh := handlers.NewUserHandler(ur)

	userGroup := r.Group("/users")
	{
		userGroup.GET("/", middlewares.ValidateToken, uh.HandleUserinf)
		userGroup.PATCH("/",
			middlewares.ValidateToken,
			uh.HandleUpdateUserInf,
		)
		userGroup.GET("/orders",
			middlewares.ValidateToken,
			uh.HandleUserOrderHistory,
		)
	}
}
