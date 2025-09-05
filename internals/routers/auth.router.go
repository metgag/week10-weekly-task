package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/metgag/koda-weekly10/internals/handlers"
	"github.com/metgag/koda-weekly10/internals/repositories"
)

func InitAuthRouter(router *gin.Engine, dbpool *pgxpool.Pool) {
	ar := repositories.NewAuthRepository(dbpool)
	ah := handlers.NewAuthHandler(ar)

	authRouter := router.Group("/auth")
	{
		authRouter.POST("/register", ah.AddUser)
		authRouter.POST("/login", ah.Login)
	}
}
