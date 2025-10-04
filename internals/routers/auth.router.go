package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/metgag/koda-weekly10/internals/handlers"
	"github.com/metgag/koda-weekly10/internals/middlewares"
	"github.com/metgag/koda-weekly10/internals/repositories"
	"github.com/redis/go-redis/v9"
)

func InitAuthRouter(router *gin.Engine, dbpool *pgxpool.Pool, rdb *redis.Client) {
	ar := repositories.NewAuthRepository(dbpool, rdb)
	ah := handlers.NewAuthHandler(ar)

	authRouter := router.Group("auth")
	{
		authRouter.POST("/register", ah.HandleRegister)
		authRouter.POST("/login", ah.HandleLogin)
		authRouter.DELETE("/logout",
			middlewares.ValidateToken(rdb),
			ah.HandleLogout,
		)
	}
}
