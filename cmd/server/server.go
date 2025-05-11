package main

import (
	"gateway/pkg/api"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем файл окружения
	godotenv.Load()

	news_host := os.Getenv("NEWS_HOST")
	if news_host == "" {
		log.Fatal("No environment for NEWS_HOST")
	}

	comments_host := os.Getenv("COMMENTS_HOST")
	if comments_host == "" {
		log.Fatal("No environment for COMMENTS_HOST")
	}

	censor_host := os.Getenv("CENSOR_HOST")
	if censor_host == "" {
		log.Fatal("No environment for CENSOR_HOST")
	}

	// Запускаем API
	apiUrls := api.APIUrls{News: news_host, Comments: comments_host, Censor: censor_host}
	api := api.NewApi(apiUrls)
	log.Print("Server is starting...")
	http.ListenAndServe(":8080", api.Router())
	log.Print("Server has been stopped.")
}
