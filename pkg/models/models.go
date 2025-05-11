package models

// Полная модель новости
type FullNews struct {
	ID      int    `json:"id"`       // идентификатор публикации
	Title   string `json:"title"`    // заголовок публикации
	Content string `json:"content"`  // содержание публикации
	PubTime int64  `json:"pub_time"` // время публикации
	Link    string `json:"link"`     // ссылка на источник
}

// Модель комментария
type Comment struct {
	ID      int    `json:"id"`       // идентификатор комментария
	Content string `json:"content"`  // содержание комментария
	PubTime int64  `json:"pub_time"` // время комментария
}

// Модель коллекции комментариев к новости
type NewsComments struct {
	Comments    []Comment         `json:"сomments"`    // основные комментарии
	Subcomments map[int][]Comment `json:"subcomments"` // подкомментарии
}

// Полная модель новости с комментариями
type FullNewsWithComments struct {
	ID          int               `json:"id"`          // идентификатор новости
	Title       string            `json:"title"`       // заголовок новости
	Content     string            `json:"content"`     // содержание новости
	PubTime     int64             `json:"pub_time"`    // время публикации
	Link        string            `json:"link"`        // ссылка на источни
	Comments    []Comment         `json:"сomments"`    // основные комментарии
	Subcomments map[int][]Comment `json:"subcomments"` // подкомментарии
}
