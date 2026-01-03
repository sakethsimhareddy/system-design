package main

import (
	"log"
	"net/http"

	"url-shortener/api"
	"url-shortener/db"
	"url-shortener/internal/handler"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"
)

func main() {
	// Initialize dependencies
	database, err := db.ConnectMongo()
	if err != nil {
		log.Fatal(err)
	}
	redisDB := db.NewRedisDB()
	repo := repository.NewMongoRepository(database,redisDB)
	svc := service.NewURLService(repo)
	h := handler.NewHandler(svc)
	router := api.NewRouter(h)

	// Start server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
