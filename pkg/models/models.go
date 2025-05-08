package models

type ShortNews struct {
	ID      int    `json:"id"`       // идентификатор новости
	Title   string `json:"title"`    // заголовок публикации
	PubTime int64  `json:"pub_time"` // время публикации
	Link    string `json:"link"`     // ссылка на источник
}

type FullNews struct {
	ID      int    `json:"id"`       // идентификатор новости
	Title   string `json:"title"`    // заголовок новости
	Content string `json:"content"`  // содержание новости
	PubTime int64  `json:"pub_time"` // время публикации
	Link    string `json:"link"`     // ссылка на источни
}

type Comment struct {
	ID          int       `json:"id"`                    // идентификатор комментария
	Text        string    `json:"text"`                  // текст комментария
	Subcomments []Comment `json:"subcomments,omitempty"` // вложенные комментарии
}

type FullNewsWithComments struct {
	ID       int       `json:"id"`                 // идентификатор новости
	Title    string    `json:"title"`              // заголовок новости
	Content  string    `json:"content"`            // содержание новости
	PubTime  int64     `json:"pub_time"`           // время публикации
	Link     string    `json:"link"`               // ссылка на источни
	Comments []Comment `json:"comments,omitempty"` // комментарии к новости
}
