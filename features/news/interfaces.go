package news

import (
	"raihpeduli/features/news/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page, size int) []News
	Insert(newNews News) int64
	SelectByID(newsID int) *News
	Update(news News) int64
	DeleteByID(newsID int) int64
}

type Usecase interface {
	FindAll(page, size int) []dtos.ResNews
	FindByID(newsID int) *dtos.ResNews
	Create(newNews dtos.InputNews) *dtos.ResNews
	Modify(newsData dtos.InputNews, newsID int) bool
	Remove(newsID int) bool
}

type Handler interface {
	GetNews() echo.HandlerFunc
	NewsDetails() echo.HandlerFunc
	CreateNews() echo.HandlerFunc
	UpdateNews() echo.HandlerFunc
	DeleteNews() echo.HandlerFunc
}
