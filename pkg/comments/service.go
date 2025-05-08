package comments

import "gateway/pkg/models"

type Service interface {
	PostComment(newsId int, comment models.Comment) error
	Comments(newsId int) ([]models.Comment, error)
}

type ServiceImpl struct {
}

func NewService() *ServiceImpl {
	return &ServiceImpl{}
}

func (s ServiceImpl) PostComment(newsId int, comment models.Comment) error {
	return nil
}

func (s ServiceImpl) Comments(newsId int) ([]models.Comment, error) {
	comment := models.Comment{ID: 0, Text: "Test"}
	return []models.Comment{comment}, nil
}
