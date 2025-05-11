package api

import (
	"bytes"
	"context"
	"encoding/json"
	"gateway/pkg/models"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ctxRequestIdKey struct{}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

// Набор URL для API
type APIUrls struct {
	// Микросервис новостей
	News string

	// Микросервис комментариев
	Comments string

	// Микросервис цензуры
	Censor string
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
	api.router.Use(requestIdValidator)
	api.router.Use(requestLogger)

	api.router.Methods(http.MethodGet).Path("/news").Handler(requestProxy(api.urls.News))
	api.router.Methods(http.MethodPost).Path("/comments").Handler(api.commentsValidator(requestProxy(api.urls.Comments)))
	api.router.Methods(http.MethodGet).Path("/news/{id}").HandlerFunc(api.newsDetails)
}

func (api *API) newsDetails(w http.ResponseWriter, r *http.Request) {
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
		url := api.urls.News + "/news/" + newsId
		newsErr = getData(url, &news)
	}()

	go func() {
		defer wg.Done()
		url := api.urls.Comments + "/comments"
		url, err := addQueryToString(url, map[string]string{"news_id": newsId})
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

func requestProxy(target string) http.HandlerFunc {
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

func requestIdValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r_key := "request_id"
		url := r.URL
		req_id := url.Query().Get(r_key)
		if req_id == "" {
			req_id = uuid.New().String()
			addQueryToUrl(url, map[string]string{r_key: req_id})
		}

		ctx := context.WithValue(r.Context(), ctxRequestIdKey{}, req_id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req_id := getRequestId(r.Context())
		d_format := "2006-01-02 15:04:05"

		log.Printf(
			"REQUEST - ID: %v IP: %v TIME: %v",
			req_id,
			r.RemoteAddr,
			time.Now().Format(d_format),
		)

		status_w := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(status_w, r)

		log.Printf(
			"RESPONSE - ID: %v STATIS: %v TIME %v",
			req_id,
			status_w.status,
			time.Now().Format(d_format),
		)
	})
}

func (api *API) commentsValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))

		v_resp, err := http.Post(api.urls.Censor+"/comments/validate", "application/json", bytes.NewReader(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if v_resp.StatusCode != http.StatusOK {
			http.Error(w, "Bad words were found", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getData[T any](url string, target *T) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	return json.NewDecoder(r.Body).Decode(target)
}

func addQueryToString(url_s string, params map[string]string) (string, error) {
	r_url, err := url.Parse(url_s)
	if err != nil {
		return "", err
	}

	addQueryToUrl(r_url, params)
	return r_url.String(), nil
}

func addQueryToUrl(url *url.URL, params map[string]string) {
	q := url.Query()
	for k, v := range params {
		q.Add(k, v)
	}

	url.RawQuery = q.Encode()
}

func getRequestId(ctx context.Context) string {
	id, ok := ctx.Value(ctxRequestIdKey{}).(string)
	if !ok {
		return ""
	}
	return id
}
