package main

import (
	"log"
	"os"
	"sauravkattel/agri/src/database"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("env load err: ", err)
	}

	dbUserName := os.Getenv("DB_USERNAME")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	_, err = database.Connect(dbUserName, dbName, dbPassword)
	if err != nil {
		log.Fatal("error connecting to db: ", err)
	}

}
