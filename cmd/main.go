package main

import (
	"log"
	"net/http"
	"os"
	mongorep "rental-server/internal/repository/mongo"
	"rental-server/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	defaultDatabase := os.Getenv("MONGODB_DATABASE")
	uri := os.Getenv("MONGODB_URI")
	rep, err := mongorep.NewMongoDBRepository(uri, defaultDatabase)
	if err != nil {
		log.Fatal(err)
	}

	server := server.NewRentObjectServer(rep)

	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatal(err)
	}
}
