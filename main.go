package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sauravkattel/agri/src/database"
	"sauravkattel/agri/src/handlers"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("env load err: ", err)
	}

	dbUserName := os.Getenv("DB_USERNAME")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := database.Connect(dbUserName, dbName, dbPassword)
	if err != nil {
		log.Fatal("error connecting to db: ", err)
	}

	routes := handlers.GetRoutes(db.DB, ctx)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "auth_token"},
		AllowCredentials: true,
		Debug:            false,
	})

	server := http.Server{
		Addr:    "localhost:4000",
		Handler: c.Handler(routes),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
