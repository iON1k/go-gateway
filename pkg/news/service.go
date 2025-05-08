package news

import "gateway/pkg/models"

type Service interface {
	LatestNews(page int) ([]models.ShortNews, error)
	FilteredNews(title string, from int64, to int64, count int) ([]models.ShortNews, error)
	News(id int) (models.FullNews, error)
}

type ServiceImpl struct {
}

func NewService() *ServiceImpl {
	return &ServiceImpl{}
}

func (s ServiceImpl) LatestNews(page int) ([]models.ShortNews, error) {
	news := models.ShortNews{ID: 0, Title: "Test", PubTime: 0, Link: "http://test.com"}
	return []models.ShortNews{news}, nil
}

func (s ServiceImpl) FilteredNews(title string, from int64, to int64, count int) ([]models.ShortNews, error) {
	news := models.ShortNews{ID: 0, Title: "Test", PubTime: 0, Link: "http://test.com"}
	return []models.ShortNews{news}, nil
}

func (s ServiceImpl) News(id int) (models.FullNews, error) {
	news := models.FullNews{ID: id, Title: "Test", Content: "Test Content", PubTime: 0, Link: "http://test.com"}
	return news, nil
}
