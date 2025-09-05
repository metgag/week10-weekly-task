package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/metgag/koda-weekly10/internals/config"
	"github.com/metgag/koda-weekly10/internals/routers"
)

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
		log.Printf("\\conninfo: connected to database \"%s\" as user \"%s\"", os.Getenv("DB_NAME"), os.Getenv("DB_USER"))
	}

	router := routers.InitRouter(dbpool)
	router.Run("localhost:6011")
}
