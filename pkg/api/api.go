package api

import (
	"encoding/json"
	"gateway/pkg/comments"
	"gateway/pkg/models"
	"gateway/pkg/news"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

// Программный интерфейс сервера
type API struct {
	newsService     news.Service
	commentsService comments.Service
	router          *mux.Router
}

// Конструктор объекта API
func NewApi(newsService news.Service, commentsService comments.Service) *API {
	api := API{newsService, commentsService, mux.NewRouter()}
	api.endpoints()
	return &api
}

// Маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.router
}

func (api *API) endpoints() {
	api.router.HandleFunc("/news/latest", api.latestNews).Methods(http.MethodGet)
	api.router.HandleFunc("/news/filter", api.filteredNews).Methods(http.MethodGet)
	api.router.HandleFunc("/news/{id}", api.news).Methods(http.MethodGet)
	api.router.HandleFunc("/news/{id}/comments", api.postComment).Methods(http.MethodPost)
}

func (api *API) latestNews(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		http.Error(w, "Page expected", http.StatusBadRequest)
		return
	}

	news, err := api.newsService.LatestNews(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(news)
}

func (api *API) filteredNews(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")

	fromStr := r.URL.Query().Get("from")
	from, _ := strconv.ParseInt(fromStr, 10, 64)

	toStr := r.URL.Query().Get("to")
	to, _ := strconv.ParseInt(toStr, 10, 64)

	countStr := r.URL.Query().Get("count")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		http.Error(w, "Count expected", http.StatusBadRequest)
		return
	}

	news, err := api.newsService.FilteredNews(title, from, to, count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(news)
}

func (api *API) news(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Id expected", http.StatusBadRequest)
		return
	}

	var news models.FullNews
	var newsErr error

	var comments []models.Comment
	var commentsErr error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		news, newsErr = api.newsService.News(id)
		wg.Done()
	}()

	go func() {
		comments, commentsErr = api.commentsService.Comments(id)
		wg.Done()
	}()

	wg.Wait()

	if newsErr != nil {
		http.Error(w, newsErr.Error(), http.StatusInternalServerError)
		return
	}

	if commentsErr != nil {
		http.Error(w, commentsErr.Error(), http.StatusInternalServerError)
		return
	}

	newsWithComments := models.FullNewsWithComments{
		ID:       news.ID,
		Title:    news.Title,
		Content:  news.Content,
		PubTime:  news.PubTime,
		Link:     news.Link,
		Comments: comments,
	}

	json.NewEncoder(w).Encode(newsWithComments)
}

func (api *API) postComment(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	newsId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Id expected", http.StatusBadRequest)
		return
	}

	var comment models.Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Comment decoding error", http.StatusBadRequest)
		return
	}

	err = api.commentsService.PostComment(newsId, comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
