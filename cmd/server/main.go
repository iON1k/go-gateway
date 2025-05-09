package main

import (
	"gateway/pkg/api"
	"gateway/pkg/comments"
	"gateway/pkg/news"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем файл окружения
	err := godotenv.Load()
	if err != nil {
		log.Fatal("No .env file found")
	}

	host := os.Getenv("HOST")
	if host == "" {
		log.Fatal("No environment for HOST")
	}

	// Запускаем API
	api := api.NewApi(news.NewService(), comments.NewService())
	http.ListenAndServe(host, api.Router())
}
