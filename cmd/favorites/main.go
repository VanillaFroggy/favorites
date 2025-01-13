package main

import (
	_ "favorites/docs"
	"favorites/internal/db"
	"favorites/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"os"
)

// @title			Favorites API
// @version		1.0
// @description	A favorites management service API in Go using Gin framework.
// @host			localhost:8080
// @BasePath		/favorites
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
	err = db.RunMigrations(dbConn, "file:///app/internal/db/migrations")
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	handlers.RegisterRoutes(dbConn, r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err = r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
