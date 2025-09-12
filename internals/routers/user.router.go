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
	userGroup.Use(
		middlewares.ValidateToken,
		middlewares.Access("general"),
	)

	{
		userGroup.GET("/", uh.HandleUserinf)
		userGroup.PATCH("/", uh.HandleUpdateUserInf)
		userGroup.GET("/orders", uh.HandleUserOrderHistory)
		userGroup.PATCH("/password", uh.HandlePasswordEdit)
	}
}
