package api

import (
	"encoding/json"
	"gateway/pkg/models"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/gorilla/mux"
)

// Набор URL для API
type APIUrls struct {
	// Микросервис новостей
	News string

	// Микросервис комментариев
	Comments string
}

// Программный интерфейс сервера
type API struct {
	urls   APIUrls
	router *mux.Router
}

// Конструктор объекта API
func NewApi(urls APIUrls) *API {
	api := API{urls, mux.NewRouter()}
	api.endpoints()
	return &api
}

// Маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.router
}

func (api *API) endpoints() {
	api.router.Path("/news/latest").Methods(http.MethodGet).Handler(proxyHandler(api.urls.News))
	api.router.Path("/news/filtered").Methods(http.MethodGet).Handler(proxyHandler(api.urls.News))
	api.router.Path("/comments").Methods(http.MethodPost).Handler(proxyHandler(api.urls.Comments))
	api.router.HandleFunc("/news/{id}", api.news).Methods(http.MethodGet)
}

func (api *API) news(w http.ResponseWriter, r *http.Request) {
	newsId := mux.Vars(r)["id"]
	if newsId == "" {
		http.Error(w, "Id expected", http.StatusBadRequest)
		return
	}

	var news models.FullNews
	var newsErr error

	var comments models.NewsComments
	var commentsErr error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		url, err := makeUrl(api.urls.News, "/news/"+newsId, nil)
		if err != nil {
			commentsErr = err
			return
		}
		newsErr = getData(url, &news)
	}()

	go func() {
		defer wg.Done()
		url, err := makeUrl(api.urls.Comments, "/comments", map[string]string{"news_id": newsId})
		if err != nil {
			commentsErr = err
			return
		}

		commentsErr = getData(url, &comments)
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

	result := models.FullNewsWithComments{
		ID:          news.ID,
		Title:       news.Title,
		Content:     news.Content,
		PubTime:     news.PubTime,
		Link:        news.Link,
		Comments:    comments.Comments,
		Subcomments: comments.Subcomments,
	}

	json.NewEncoder(w).Encode(result)
}

func proxyHandler(target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetURL, err := url.Parse(target)
		if err != nil {
			http.Error(w, "Bad target URL", http.StatusInternalServerError)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		r.Host = targetURL.Host
		proxy.ServeHTTP(w, r)
	}
}

func getData[T any](url string, target *T) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	return json.NewDecoder(r.Body).Decode(target)
}

func makeUrl(base string, endpoint string, params map[string]string) (string, error) {
	req_url, err := url.Parse(base + endpoint)
	if err != nil {
		return "", err
	}
	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}

	req_url.RawQuery = values.Encode()
	return req_url.String(), nil
}
