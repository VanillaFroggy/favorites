package main

import (
	"favorites/internal/db"
	"favorites/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"os"
)

func main() {
	dbConn, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}
	defer func(dbConn *sqlx.DB) {
		err := dbConn.Close()
		if err != nil {
			panic(err)
		}
	}(dbConn)
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		db.RunMigrations(dbConn)
	}
	r := gin.Default()
	handler.RegisterRoutes(dbConn, r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err = r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
