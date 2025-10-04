package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	config "github.com/metgag/koda-weekly10/internals/configs"
	"github.com/metgag/koda-weekly10/internals/routers"
)

//	@title			LOKET TIKET
//	@version		1.0
//	@description	RESTful API of tixkitz ticket systeme

//	@host		localhost:6011
//	@basepath	/

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	// load dotenv
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	// init db psql
	dbpool, err := config.InitDB()
	if err != nil {
		log.Printf("unable to create connection pool: %s\n", err)
	}
	defer dbpool.Close()

	if err := config.PingDB(dbpool); err != nil {
		fmt.Fprintf(os.Stderr, "unable to connect to database: %v\n", err)
	} else {
		log.Printf("\\conninfo: connected to database \"%s\" as user \"%s\"", os.Getenv("DB_NAME_M"), os.Getenv("DB_USER_M"))
	}

	// init redis
	rdb := config.InitRedis()
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Printf("unable to connect to redis server")
	} else {
		log.Printf("\nREDIS> PING\n%sPONG", strings.Repeat(" ", 7))
	}
	defer rdb.Close()

	router := routers.InitRouter(dbpool, rdb)
	router.Run(":6011")
}
